package workflows

import (
	"testing"

	"github.com/kube-tarian/kad/integrator/config-worker/pkg/activities"
	"github.com/kube-tarian/kad/integrator/model"
	"go.temporal.io/sdk/worker"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_Workflow() {
	env := s.NewTestWorkflowEnvironment()
	env.SetWorkerOptions(worker.Options{
		EnableSessionWorker: true, // Important for a worker to participate in the session
	})
	var a *activities.Activities

	env.OnActivity(a.ConfigurationActivity, mock.Anything, model.RequestPayload{Action: "file1"}).Return(model.ResponsePayload{Status: "file2"}, nil)

	env.RegisterActivity(a)

	env.ExecuteWorkflow(Workflow, model.RequestPayload{Action: "file1"})

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	env.AssertExpectations(s.T())
}
