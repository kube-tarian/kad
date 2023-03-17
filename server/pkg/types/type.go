package types

const (
	ClientCertChainFileName = "cert-chain.pem"
	ClientCertFileName      = "client.crt"
	ClientKeyFileName       = "client.key"
)

type AgentInfo struct {
	Endpoint string
	CaPem    string
	Cert     string
	Key      string
}

type AgentConfiguration struct {
	Address string `envconfig:"AGENT_ADDRESS" default:"localhost"`
	Port    int    `envconfig:"AGENT_PORT" default:"9091"`
	CaCert  string
	Cert    string
	Key     string
}
