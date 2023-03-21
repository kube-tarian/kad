package model

import "github.com/kube-tarian/kad/server/api"

//	type DeployPayload struct {
//		Operation string               `json:"operation"`
//		Payload   DeployRequestPayload `json:"payload"`
//	}
type DeployPayload struct {
	Operation  string
	Plugin     string
	WorkerType string
	Payload    string
}
type DeployResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type AgentsResponse = []api.AgentRequest
