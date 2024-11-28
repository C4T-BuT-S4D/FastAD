package checkers

import (
	"context"

	"github.com/sirupsen/logrus"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type PutActivityParameters struct {
	FlagInfo *FQFlagInfo
}

type PutActivityResult struct {
	FlagInfo *FQFlagInfo
	Verdict  *Verdict
}

func (s *ActivityState) PutActivityDefinition(ctx context.Context, params *PutActivityParameters) (*PutActivityResult, error) {
	logger := logrus.WithFields(logrus.Fields{
		"team":    params.FlagInfo.Team.Name,
		"service": params.FlagInfo.Service.Name,
		"action":  checkerpb.Action_ACTION_PUT,
	})

	logger.Info("starting")
	verdict := RunPutAction(ctx, params)
	logger.Infof("finished: %v", verdict)

	return &PutActivityResult{
		FlagInfo: params.FlagInfo,
		Verdict:  verdict,
	}, nil
}
