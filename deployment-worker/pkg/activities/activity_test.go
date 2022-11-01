package activities

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Activity(t *testing.T) {
	a := &Activities{}
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(a)

	val, err := env.ExecuteActivity(a.Activity, "World")
	require.NoError(t, err)

	var res string
	require.NoError(t, val.Get(&res))
	require.Equal(t, "Hello World!", res)
}
