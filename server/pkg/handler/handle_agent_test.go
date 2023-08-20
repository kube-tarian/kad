package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/agent"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/stretchr/testify/require"
)

func TestAPIHandler_Close(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		customerId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.agentHandler.RemoveAgent(tt.args.customerId)
		})
	}
}

func TestAPIHandler_CloseAll(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.Close()
		})
	}
}

func TestAPIHandler_ConnectClient(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		customerId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			if _, err := a.agentHandler.GetAgent(tt.args.customerId); (err != nil) != tt.wantErr {
				t.Errorf("ConnectClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAPIHandler_DeleteAgentClimondeploy(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.DeleteAgentClimondeploy(tt.args.c)
		})
	}
}

func TestAPIHandler_DeleteAgentCluster(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.DeleteAgentCluster(tt.args.c)
		})
	}
}

func TestAPIHandler_DeleteAgentDeploy(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.DeleteAgentDeploy(tt.args.c)
		})
	}
}

func TestAPIHandler_DeleteAgentProject(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.DeleteAgentProject(tt.args.c)
		})
	}
}

func TestAPIHandler_DeleteAgentRepository(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.DeleteAgentRepository(tt.args.c)
		})
	}
}

func TestAPIHandler_GetAgentEndpoint(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.GetAgentEndpoint(tt.args.c)
		})
	}
}

func TestAPIHandler_GetApiDocs(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.GetApiDocs(tt.args.c)
		})
	}
}

func TestAPIHandler_GetClient(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		customerId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *agent.Agent
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			if got, err := a.agentHandler.GetAgent(tt.args.customerId); err != nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIHandler_GetStatus(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.GetStatus(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentApps(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}

	chartName := "argocd"
	name := "argocd"
	//override := ""
	releaseName := "test"
	repoName := "test"
	repoURL := "https://argocd.com"
	version := "v1.0.0"
	namespace := "capten"
	tools := []struct {
		ChartName   *string `json:"chartName,omitempty"`
		Name        *string `json:"name,omitempty"`
		Namespace   *string `json:"namespace,omitempty"`
		Override    *string `json:"override,omitempty"`
		ReleaseName *string `json:"releaseName,omitempty"`
		RepoName    *string `json:"repoName,omitempty"`
		RepoURL     *string `json:"repoURL,omitempty"`
		Version     *string `json:"version,omitempty"`
	}{
		{
			ChartName:   &chartName,
			Name:        &name,
			ReleaseName: &releaseName,
			RepoName:    &repoName,
			RepoURL:     &repoURL,
			Version:     &version,
			Namespace:   &namespace,
		},
	}

	apps := api.AgentAppsRequest{
		Apps: &tools,
	}

	jsonByte, err := json.Marshal(apps)
	require.NoError(t, err)
	fmt.Println(string(jsonByte))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("customer_id", "1")
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonByte))
	fmt.Println(c.Request.Body)
	agentConn, err := agent.NewAgent(logging.NewLogger(), &agent.Config{
		Address: "127.0.0.1",
	}, nil)

	require.NoError(t, err)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "post apps",
			fields: fields{agents: map[string]*agent.Agent{
				"1": agentConn,
			}},
			args: args{c: c},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.PostAgentApps(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentClimondeploy(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.PostAgentClimondeploy(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentCluster(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.PostAgentCluster(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentDeploy(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.PostAgentDeploy(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentEndpoint(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.PostAgentEndpoint(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentProject(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.PostAgentProject(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentRepository(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.PostAgentRepository(tt.args.c)
		})
	}
}

func TestAPIHandler_PostAgentSecret(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.PostAgentSecret(tt.args.c)
		})
	}
}

func TestAPIHandler_PutAgentClimondeploy(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.PutAgentClimondeploy(tt.args.c)
		})
	}
}

func TestAPIHandler_PutAgentDeploy(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.PutAgentDeploy(tt.args.c)
		})
	}
}

func TestAPIHandler_PutAgentEndpoint(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.PutAgentEndpoint(tt.args.c)
		})
	}
}

func TestAPIHandler_PutAgentProject(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.PutAgentProject(tt.args.c)
		})
	}
}

func TestAPIHandler_PutAgentRepository(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.PutAgentRepository(tt.args.c)
		})
	}
}

func TestAPIHandler_getFileContent(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c        *gin.Context
		fileInfo map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			got, err := a.getFileContent(tt.args.c, tt.args.fileInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFileContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIHandler_sendResponse(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c   *gin.Context
		msg string
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.sendResponse(tt.args.c, tt.args.msg, tt.args.err)
		})
	}
}

func TestAPIHandler_setFailedResponse(t *testing.T) {
	type fields struct {
		agents map[string]*agent.Agent
	}
	type args struct {
		c   *gin.Context
		msg string
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				agentHandler: agent.NewAgentHandler(logging.NewLogger(), nil, nil),
			}
			a.setFailedResponse(tt.args.c, tt.args.msg, tt.args.err)
		})
	}
}

func TestNewAPIHandler(t *testing.T) {
	tests := []struct {
		name    string
		want    *APIHandler
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAPIHandler(logging.NewLogger(), nil, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAPIHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPIHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toString(t *testing.T) {
	type args struct {
		resp *agentpb.JobResponse
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toString(tt.args.resp); got != tt.want {
				t.Errorf("toString() = %v, want %v", got, tt.want)
			}
		})
	}
}
