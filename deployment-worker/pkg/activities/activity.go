package activities

import (
	"context"

	"go.temporal.io/sdk/activity"
)

type Activities struct {
}

func (a *Activities) Activity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", name)
	return "Hello " + name + "!", nil
}
