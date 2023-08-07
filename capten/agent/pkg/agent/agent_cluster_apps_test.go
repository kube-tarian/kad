package agent

import (
	"context"
	"os"
	"testing"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/config"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

type AgentTestSuite struct {
	suite.Suite
	agent  *Agent
	logger logging.Logger
}

func TestAgentTestSuite(t *testing.T) {
	setEnvVars()
	agentSuite := new(AgentTestSuite)
	agentSuite.logger = logging.NewLogger()

	/*if err := captenstore.Migrate(agentSuite.logger); err != nil {
		t.Fatal(err)
	}*/

	agent, err := NewAgent(agentSuite.logger, &config.SericeConfig{})
	if err != nil {
		t.Fatal(err)
	}

	agentSuite.agent = agent
	suite.Run(t, agentSuite)
}

func (suite *AgentTestSuite) SetupSuite() {
}

func (suite *AgentTestSuite) TearDownSuite() {
	/*if err := captenstore.MigratePurge(suite.logger); err != nil {
		suite.logger.Error(err.Error())
	}*/
}

func (suite *AgentTestSuite) Test_1_SyncApp() {
	request, req2 := &agentpb.SyncAppRequest{}, &agentpb.SyncAppRequest{}

	_, err := suite.agent.SyncApp(context.TODO(), request)
	suite.NotNil(err)

	request.Data = &agentpb.SyncAppData{Config: &agentpb.AppConfig{ReleaseName: "release"}}
	res, err := suite.agent.SyncApp(context.TODO(), request)

	suite.Nil(err)
	suite.NotNil(res)
	suite.Equal(agentpb.StatusCode(0), res.Status)
	suite.Equal("OK", res.StatusMessage)

	override, err := yaml.Marshal(overrideValues)
	suite.Nil(err)

	req2.Data = &agentpb.SyncAppData{
		Config: &agentpb.AppConfig{
			ReleaseName: "release2",
			Icon:        []byte{0x1, 0x2, 0x3, 0x4},
		},
		Values: &agentpb.AppValues{
			OverrideValues: override,
		},
	}
	res, err = suite.agent.SyncApp(context.TODO(), req2)

	suite.Nil(err)
	suite.NotNil(res)
	suite.Equal(agentpb.StatusCode(0), res.Status)
	suite.Equal("OK", res.StatusMessage)

}

func (suite *AgentTestSuite) Test_2_GetClusterApp() {
	req := &agentpb.GetClusterAppsRequest{}
	res, err := suite.agent.GetClusterApps(context.TODO(), req)

	suite.Nil(err)
	suite.Equal(2, len(res.GetAppData()))

	res2, err := suite.agent.GetClusterAppLaunches(context.TODO(), &agentpb.GetClusterAppLaunchesRequest{})

	suite.Nil(err)
	suite.Equal(2, len(res2.GetLaunchConfigList()))
}

func (suite *AgentTestSuite) Test_3_GetLaunchConfigList() {
	res, err := suite.agent.GetClusterAppLaunches(context.TODO(), &agentpb.GetClusterAppLaunchesRequest{})

	suite.Nil(err)
	suite.Equal(2, len(res.GetLaunchConfigList()))
}

func (suite *AgentTestSuite) Test_3_GetClusterAppConfig() {

	res, err := suite.agent.GetClusterAppConfig(context.TODO(), &agentpb.GetClusterAppConfigRequest{ReleaseName: "abc"})

	suite.NotNil(res)
	suite.Nil(err)
	suite.Equal("NOT_FOUND", res.GetStatusMessage())

	res, err = suite.agent.GetClusterAppConfig(context.TODO(), &agentpb.GetClusterAppConfigRequest{ReleaseName: "release2"})

	suite.NotNil(res)
	suite.Nil(err)

	suite.Equal("release2", res.GetAppConfig().GetReleaseName())
	suite.Equal([]byte{1, 2, 3, 4}, res.GetAppConfig().GetIcon())

}

func (suite *AgentTestSuite) Test_3_GetClusterAppValues() {

	res, err := suite.agent.GetClusterAppValues(context.TODO(), &agentpb.GetClusterAppValuesRequest{ReleaseName: "abc"})

	suite.NotNil(res)
	suite.Nil(err)
	suite.Equal("NOT_FOUND", res.GetStatusMessage())

	res, err = suite.agent.GetClusterAppValues(context.TODO(), &agentpb.GetClusterAppValuesRequest{ReleaseName: "release2"})

	suite.NotNil(res)
	suite.Nil(err)

	override, _ := yaml.Marshal(overrideValues)
	suite.Equal(override, res.GetValues().GetOverrideValues())

}

func setEnvVars() {

	os.Setenv("DB_ADDRESSES", "localhost:9042")
	os.Setenv("DB_ENTITY_NAME", "TEST_ENTITY")
	os.Setenv("DB_NAME", "apps")
	os.Setenv("DB_NAME", "apps")
	os.Setenv("DB_SERVICE_USERNAME", "apps_user")
	os.Setenv("DB_SERVICE_PASSWD", "apps_password")
	os.Setenv("SOURCE_URI", "file://../capten-store/test_migrations")
	os.Setenv("ENV", "LOCAL")

}

var overrideValues = `
clickhouse:
	password": admin
frontend:
    ingress:
      enabled": true
      hosts:
      - host: "signoz.{{.DomainName}}"
        paths:
        - path: /
          pathType: ImplementationSpecific
          port: 3301
      tls:
        hosts:
        - "signoz.{{.DomainName}}"
        secretName: cert-signoz
    annotations:
      cert-manager.io/cluster-issuer": letsencrypt-prod-cluster
      nginx.ingress.kubernetes.io/backend-protocol": HTTPS
      nginx.ingress.kubernetes.io/force-ssl-redirect": true
      nginx.ingress.kubernetes.io/ssl-redirect": true
`
