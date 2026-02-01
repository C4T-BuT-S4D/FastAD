//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/c4t-but-s4d/fastad/internal/gameconfig"
	"github.com/c4t-but-s4d/fastad/internal/models"
	gamestateClient "github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	servicesClient "github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
	scoreboardpb "github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
)

const (
	baseURL      = "http://localhost:8080"
	grpcAddress  = "localhost:8080" // gRPC is served on same port as HTTP (multiplexed via h2c)
	testTimeout  = 5 * time.Minute
	pollInterval = 2 * time.Second
	dbDSN        = "postgres://fastad:fastad@localhost:5433/fastad?sslmode=disable"
)

type teamToken struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type SharedTestState struct {
	t          *testing.T
	rootDir    string
	preset     string
	db         *bun.DB
	grpcConn   *grpc.ClientConn
	gsClient   *gamestateClient.Client
	svcClient  *servicesClient.Client
	teams      []*teamspb.Team
	teamTokens []teamToken
	services   []*servicespb.Service
}

func NewSharedTestState(t *testing.T, preset string) *SharedTestState {
	s := &SharedTestState{
		t:      t,
		preset: preset,
	}
	s.rootDir = s.findRootDir()
	s.initGame()
	s.startGame()
	s.waitForHealthy()
	s.waitForGameRunning()
	s.connectDB()
	s.connectGRPC()
	s.loadTeams()
	s.loadTeamTokens()
	s.loadServices()
	return s
}

func (s *SharedTestState) TearDown() {
	s.t.Log("Tearing down test environment...")

	if s.grpcConn != nil {
		if err := s.grpcConn.Close(); err != nil {
			s.t.Logf("Warning: failed to close gRPC connection: %v", err)
		}
	}

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			s.t.Logf("Warning: failed to close DB: %v", err)
		}
	}

	cmd := exec.Command("go", "run", "./cmd/fastad", "reset")
	cmd.Dir = s.rootDir

	if output, err := cmd.CombinedOutput(); err != nil {
		s.t.Logf("Warning: failed to tear down: %v\nOutput: %s", err, output)
	}
}

func (s *SharedTestState) findRootDir() string {
	dir, err := os.Getwd()
	require.NoError(s.t, err)

	for {
		if _, err := os.Stat(filepath.Join(dir, ".fastad_root")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			s.t.Fatal("Could not find FastAD root directory")
		}
		dir = parent
	}
}

func (s *SharedTestState) initGame() {
	s.t.Log("Initializing game with preset:", s.preset)

	testConfigPath := filepath.Join(s.rootDir, "tests", "e2e", "fastad_test.yaml")

	cmd := exec.Command(
		"go", "run", "./cmd/fastad",
		"init",
		"--game-config", testConfigPath,
		"--preset", s.preset,
		"--force",
	)
	cmd.Dir = s.rootDir
	cmd.Env = append(os.Environ(), "FASTAD_LOG_LEVEL=debug")

	output, err := cmd.CombinedOutput()
	if err != nil {
		s.t.Fatalf("Failed to initialize game: %v\nOutput: %s", err, output)
	}
	s.t.Log("Game initialized successfully")
}

func (s *SharedTestState) startGame() {
	s.t.Log("Starting game services...")

	cmd := exec.Command(
		"go", "run", "./cmd/fastad",
		"run",
	)
	cmd.Dir = s.rootDir
	cmd.Env = append(os.Environ(), "FASTAD_LOG_LEVEL=debug")

	output, err := cmd.CombinedOutput()
	if err != nil {
		s.t.Fatalf("Failed to start game: %v\nOutput: %s", err, output)
	}
	s.t.Log("Game started successfully")
}

func (s *SharedTestState) waitForHealthy() {
	s.t.Log("Waiting for services to become healthy...")

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			s.t.Fatal("Timeout waiting for services to become healthy")
		default:
		}

		resp, err := http.Get(baseURL + "/healthcheck")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			s.t.Log("Services are healthy")
			return
		}
		if resp != nil {
			resp.Body.Close()
		}

		time.Sleep(pollInterval)
	}
}

func (s *SharedTestState) waitForGameRunning() {
	s.t.Log("Waiting for game to transition to RUNNING...")

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			s.t.Fatal("Timeout waiting for game to become RUNNING")
		default:
		}

		resp, err := http.Get(baseURL + "/api/game")
		if err == nil && resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			var gs gspb.GameState
			if err := protojson.Unmarshal(body, &gs); err == nil {
				if gs.GetStatus() == gspb.GameStatus_GAME_STATUS_RUNNING {
					s.t.Log("Game is RUNNING")
					return
				}
				s.t.Logf("Game status is %s, waiting...", gs.GetStatus())
			}
		}
		if resp != nil {
			resp.Body.Close()
		}

		time.Sleep(pollInterval)
	}
}

func (s *SharedTestState) connectDB() {
	s.t.Log("Connecting to database...")

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbDSN)))
	s.db = bun.NewDB(sqlDB, pgdialect.New())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		s.t.Fatalf("Failed to connect to database: %v", err)
	}

	s.t.Log("Database connected successfully")
}

func (s *SharedTestState) connectGRPC() {
	s.t.Log("Connecting to gRPC...")

	cfg, err := gameconfig.ReadGeneratedConfig()
	if err != nil {
		s.t.Fatalf("Failed to read generated config: %v", err)
	}

	conn, err := grpcext.Dial(
		grpcAddress,
		"e2e-tests",
		grpcext.AuthDialOptions(cfg.FastAD.IntercomToken)...,
	)
	if err != nil {
		s.t.Fatalf("Failed to connect to gRPC: %v", err)
	}
	s.grpcConn = conn

	s.gsClient = gamestateClient.NewClient(gspb.NewGameStateServiceClient(conn), "e2e-tests")
	s.svcClient = servicesClient.NewClient(servicespb.NewServicesServiceClient(conn), "e2e-tests")

	s.t.Log("gRPC connected successfully")
}

func (s *SharedTestState) loadTeams() {
	s.t.Log("Loading teams...")

	resp, err := http.Get(baseURL + "/api/teams")
	require.NoError(s.t, err)
	defer resp.Body.Close()

	require.Equal(s.t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(s.t, err)

	var teamsResp teamspb.Team_Batch
	err = protojson.Unmarshal(body, &teamsResp)
	require.NoError(s.t, err)

	s.teams = teamsResp.GetTeams()
	s.t.Logf("Loaded %d teams", len(s.teams))
}

func (s *SharedTestState) loadTeamTokens() {
	s.t.Log("Loading team tokens...")

	cmd := exec.Command("go", "run", "./cmd/fastad", "tokens", "--format", "json")
	cmd.Dir = s.rootDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		s.t.Logf("Warning: Failed to get team tokens: %v\nStderr: %s", err, stderr.String())
		return
	}

	if err := json.Unmarshal(stdout.Bytes(), &s.teamTokens); err != nil {
		s.t.Logf("Warning: Failed to parse tokens: %v\nStdout: %s", err, stdout.String())
		return
	}

	s.t.Logf("Loaded %d team tokens", len(s.teamTokens))
}

func (s *SharedTestState) loadServices() {
	s.t.Log("Loading services...")

	resp, err := http.Get(baseURL + "/api/services")
	require.NoError(s.t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(s.t, err)

	var servicesResp servicespb.Service_Batch
	err = protojson.Unmarshal(body, &servicesResp)
	require.NoError(s.t, err)

	s.services = servicesResp.GetServices()
	s.t.Logf("Loaded %d services", len(s.services))
}

type BaseSuite struct {
	suite.Suite
	sharedState *SharedTestState
}

func (s *BaseSuite) SetupSuite() {
	if s.sharedState == nil {
		s.T().Fatal("sharedState must be set before running suite")
	}
}

func (s *BaseSuite) TearDownSuite() {
}

func (s *BaseSuite) rootDir() string {
	return s.sharedState.rootDir
}

func (s *BaseSuite) db() *bun.DB {
	return s.sharedState.db
}

func (s *BaseSuite) gsClient() *gamestateClient.Client {
	return s.sharedState.gsClient
}

func (s *BaseSuite) svcClient() *servicesClient.Client {
	return s.sharedState.svcClient
}

func (s *BaseSuite) teams() []*teamspb.Team {
	return s.sharedState.teams
}

func (s *BaseSuite) teamTokens() []teamToken {
	return s.sharedState.teamTokens
}

func (s *BaseSuite) services() []*servicespb.Service {
	return s.sharedState.services
}

func (s *BaseSuite) GetGameState() *gspb.GameState {
	resp, err := http.Get(baseURL + "/api/game")
	s.Require().NoError(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var gs gspb.GameState
	err = protojson.Unmarshal(body, &gs)
	s.Require().NoError(err)

	return &gs
}

func (s *BaseSuite) GetCurrentRound() uint64 {
	return s.GetGameState().GetRunningRound()
}

func (s *BaseSuite) GetScoreboard() *scoreboardpb.Scoreboard {
	resp, err := http.Get(baseURL + "/api/scoreboard")
	s.Require().NoError(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var sb scoreboardpb.Scoreboard
	err = protojson.Unmarshal(body, &sb)
	s.Require().NoError(err)

	return &sb
}

func (s *BaseSuite) GetTeamServiceState(sb *scoreboardpb.Scoreboard, teamID, serviceID int64) *scoreboardpb.Scoreboard_TeamServiceState {
	for _, state := range sb.GetTeamServiceStates() {
		if state.GetTeamId() == teamID && state.GetServiceId() == serviceID {
			return state
		}
	}
	return nil
}

func (s *BaseSuite) SubmitFlags(token string, flags []string) *receiverpb.SubmitFlagsResponse {
	flagsJSON := `["` + strings.Join(flags, `","`) + `"]`
	body := []byte(`{"flags":` + flagsJSON + `}`)

	req, err := http.NewRequest(http.MethodPost, baseURL+"/flags", bytes.NewReader(body))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Team-Token", token)

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	if resp.StatusCode != http.StatusOK {
		s.T().Logf("Flag submission failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var flagsResp receiverpb.SubmitFlagsResponse
	err = protojson.Unmarshal(respBody, &flagsResp)
	s.Require().NoError(err, "Failed to unmarshal response: %s", string(respBody))

	return &flagsResp
}

func (s *BaseSuite) InsertTestFlag(teamID int, serviceID int, round uint64, putFinished bool) *models.Flag {
	ctx := context.Background()

	flag := &models.Flag{
		Flag:        fmt.Sprintf("TEST_%s_%d_%d_%d=", s.randomString(20), teamID, serviceID, round),
		TeamID:      teamID,
		ServiceID:   serviceID,
		Round:       round,
		PutFinished: putFinished,
		CreatedAt:   time.Now(),
	}

	_, err := s.db().NewInsert().Model(flag).Exec(ctx)
	s.Require().NoError(err, "Failed to insert test flag")

	s.T().Logf("Inserted test flag ID=%d for team %d, round %d: %s", flag.ID, teamID, round, flag.Flag)
	return flag
}

func (s *BaseSuite) InsertCheckerExecution(teamID, serviceID int, action checkerpb.Action, status checkerpb.Status, publicMsg string) *models.CheckerExecution {
	ctx := context.Background()

	exec := &models.CheckerExecution{
		ExecutionID: uuid.New().String(),
		TeamID:      teamID,
		ServiceID:   serviceID,
		Action:      action,
		Status:      status,
		Public:      publicMsg,
		CreatedAt:   time.Now(),
	}

	_, err := s.db().NewInsert().Model(exec).Exec(ctx)
	s.Require().NoError(err, "Failed to insert checker execution")

	s.T().Logf("Inserted checker execution ID=%d for team %d, service %d: %s -> %s",
		exec.ID, teamID, serviceID, action.String(), status.String())
	return exec
}

func (s *BaseSuite) WaitForRoundIncrement(initialRound uint64, timeout time.Duration) uint64 {
	var currentRound uint64
	require.Eventually(s.T(), func() bool {
		currentRound = s.GetCurrentRound()
		return currentRound > initialRound
	}, timeout, time.Second, "Round should increment from %d", initialRound)
	return currentRound
}

func (s *BaseSuite) WaitForCondition(condition func() bool, timeout time.Duration, interval time.Duration, msg string) {
	require.Eventually(s.T(), condition, timeout, interval, msg)
}

func (s *BaseSuite) randomString(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}

func (s *BaseSuite) SetGameStatus(status gspb.GameStatus) {
	ctx := context.Background()
	_, err := s.gsClient().Update(ctx, &gspb.UpdateRequest{
		Status: status,
	})
	s.Require().NoError(err, "Failed to set game status=%s", status)
	s.T().Logf("Set game status=%s via API", status)
}

func (s *BaseSuite) SetServiceDisabled(serviceID int, disabled bool) {
	ctx := context.Background()
	_, err := s.svcClient().Update(ctx, &servicespb.UpdateRequest{
		Id:       int64(serviceID),
		Disabled: proto.Bool(disabled),
	})
	s.Require().NoError(err, "Failed to set service %d disabled=%v via API", serviceID, disabled)
	s.T().Logf("Set service %d disabled=%v via API", serviceID, disabled)
}
