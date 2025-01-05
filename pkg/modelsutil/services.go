package modelsutil

import (
	"time"

	"github.com/samber/lo"

	"github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	"github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
)

const defaultRunCount = 1

func ServiceCheckerTimeout(service *services.Service, action checker.Action) time.Duration {
	actionCfg, ok := lo.Find(service.GetChecker().GetActions(), func(item *services.Service_Checker_Action) bool {
		return item.GetAction() == action
	})
	if !ok || actionCfg.GetTimeout().AsDuration() == 0 {
		return service.GetChecker().GetDefaultTimeout().AsDuration()
	}
	return actionCfg.GetTimeout().AsDuration()
}

func ServiceCheckerRunCount(service *services.Service, action checker.Action) int {
	actionCfg, ok := lo.Find(service.GetChecker().GetActions(), func(item *services.Service_Checker_Action) bool {
		return item.GetAction() == action
	})
	if !ok || actionCfg.GetRunCount() == 0 {
		return defaultRunCount
	}
	return int(actionCfg.GetRunCount())
}
