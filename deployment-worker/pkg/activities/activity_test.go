package activities

import (
	"testing"

	"github.com/kube-tarian/kad/deployment-worker/pkg/model"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Activity(t *testing.T) {
	a := &Activities{}
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(a)

	val, err := env.ExecuteActivity(a.DeploymentActivity, model.RequestPayload{SubAction: "World"})
	require.NoError(t, err)

	var res model.ResponsePayload
	require.NoError(t, val.Get(&res))
	require.Equal(t, "Success", res.Status)
}
