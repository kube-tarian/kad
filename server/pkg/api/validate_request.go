package api

import (
	"context"
	"fmt"
)

func validateRequest(ctx context.Context, args ...string) (string, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if orgId == "" {
		return orgId, fmt.Errorf("organizationid is missing")
	}

	for _, arg := range args {
		if len(arg) == 0 {
			return orgId, fmt.Errorf("mandatory argument is empty")
		}
	}
	return orgId, nil
}
