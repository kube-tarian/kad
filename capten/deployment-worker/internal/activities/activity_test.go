package activities

import (
	"testing"

	"github.com/kube-tarian/kad/capten/model"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Activity(t *testing.T) {
	a := &Activities{}
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(a)

	_, err := env.ExecuteActivity(a.DeploymentInstallActivity, model.RequestPayload{Action: "World"})
	require.Error(t, err)

	// var res model.ResponsePayload
	// require.NoError(t, val.Get(&res))
	// require.Equal(t, "Success", res.Status)
}
