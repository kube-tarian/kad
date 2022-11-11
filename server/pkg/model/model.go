package model

type DeployPayload struct {
	Operation string               `json:"operation"`
	Payload   DeployRequestPayload `json:"payload"`
}

type DeployResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
