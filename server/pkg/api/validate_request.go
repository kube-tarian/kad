package api

import (
	"context"
	"fmt"
)

func validateOrgWithArgs(ctx context.Context, args ...string) (string, error) {
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

func validateOrgClusterWithArgs(ctx context.Context, args ...string) (orgId string, clusteId string, err error) {
	metadataMap := metadataContextToMap(ctx)
	orgId = metadataMap[organizationIDAttribute]
	if orgId == "" {
		err = fmt.Errorf("organizationid is missing")
		return
	}

	clusteId = metadataMap[clusterIDAttribute]
	if orgId == "" {
		err = fmt.Errorf("clusteId is missing")
		return
	}

	if err = stringArrayEmptyCheck(args); err != nil {
		return
	}

	return
}

func validateArgs(args ...string) error {
	return stringArrayEmptyCheck(args)
}

func stringArrayEmptyCheck(args []string) error {
	for _, arg := range args {
		if len(arg) == 0 {
			return fmt.Errorf("mandatory argument is empty")
		}
	}
	return nil
}
