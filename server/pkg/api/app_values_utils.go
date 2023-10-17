package api

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func (s *Server) replaceGlobalValues(orgId, clusterID string, overridedValues []byte) ([]byte, error) {
	agent, err := s.agentHandeler.GetAgent(orgId, clusterID)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize agent for cluster %s", clusterID)
	}
	resp, err := agent.GetClient().GetClusterGlobalValues(context.TODO(), &agentpb.GetClusterGlobalValuesRequest{})
	if err != nil {
		return nil, err
	}
	if resp.Status != agentpb.StatusCode_OK {
		return nil, fmt.Errorf("failed to get global values for cluster %s", clusterID)
	}

	var globalValues map[string]interface{}
	err = yaml.Unmarshal(resp.GlobalValues, &globalValues)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to unmarshal cluster values")
	}

	var overrideValues map[string]interface{}
	err = yaml.Unmarshal(overridedValues, &overrideValues)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to unmarshal override values")
	}

	return replaceOverrideGlobalValues(overrideValues, globalValues)
}

func replaceOverrideGlobalValues(overrideValues map[string]interface{},
	globlaValues map[string]interface{}) (transformedData []byte, err error) {
	yamlData, err := yaml.Marshal(overrideValues)
	if err != nil {
		return
	}

	tmpl, err := template.New("templateVal").Parse(string(yamlData))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, globlaValues)
	if err != nil {
		return
	}

	transformedData = buf.Bytes()
	return
}
