package teams

import (
	"fmt"

	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

func (s *Service) validateCreateBatchRequest(req *teamspb.CreateBatchRequest) error {
	if len(req.Teams) == 0 {
		return fmt.Errorf("teams required")
	}
	for i, team := range req.Teams {
		if team.Name == "" {
			return fmt.Errorf("teams.%d: name required", i)
		}
		if team.Address == "" {
			return fmt.Errorf("teams.%d: address required", i)
		}
	}
	return nil
}
