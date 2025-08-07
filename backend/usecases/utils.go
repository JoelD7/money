package usecases

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

// PeriodHolder is an interface that describes entities that have a period.
//
// It was created with the purpose of abstracting away the process of setting period attributes to the entities that
// have a period.
type PeriodHolder interface {
	GetPeriodID() string
	SetPeriodName(name string)
}

func setEntitiesPeriods(ctx context.Context, username string, pm PeriodManager, entitites ...PeriodHolder) error {
	expensesPeriods := make([]string, 0, len(entitites))
	seen := make(map[string]struct{}, len(entitites))
	var periodID string

	for _, entity := range entitites {
		periodID = entity.GetPeriodID()

		if _, ok := seen[periodID]; ok || periodID == "" {
			continue
		}
		expensesPeriods = append(expensesPeriods, periodID)
		seen[periodID] = struct{}{}
	}

	periodResults, err := pm.BatchGetPeriods(ctx, username, expensesPeriods)
	if err != nil {
		return err
	}

	periodsByID := make(map[string]*models.Period, len(periodResults))
	for _, period := range periodResults {
		periodsByID[period.ID] = period
	}

	for _, entity := range entitites {
		p, ok := periodsByID[entity.GetPeriodID()]
		if !ok {
			continue
		}
		entity.SetPeriodName(*p.Name)
	}

	return nil
}
