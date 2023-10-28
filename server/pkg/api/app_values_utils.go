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

func (s *Server) getClusterGlobalValues(orgId, clusterID string) (map[string]interface{}, error) {
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
	s.log.Debugf("cluster %s globalValues: %+v", clusterID, globalValues)
	return globalValues, nil
}

func (s *Server) deriveTemplateOverrideValues(overridedValues []byte,
	globlaValues map[string]interface{}) (transformedOverrideValues []byte, err error) {
	tmpl, err := template.New("templateVal").Parse(string(overridedValues))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, globlaValues)
	if err != nil {
		return
	}

	transformedOverrideValues = buf.Bytes()
	s.log.Debugf("cluster transformedOverrideValues: %+v", string(transformedOverrideValues))
	return
}
