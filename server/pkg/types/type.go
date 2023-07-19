package types

const (
	ClientCertChainFileName = "cert-chain.pem"
	ClientCertFileName      = "client.crt"
	ClientKeyFileName       = "client.key"
	AgentPortCfgKey         = "agent.port"
	AgentTlsEnabledCfgKey   = "agent.tlsEnabled"
	ServerDbCfgKey          = "server.db"
)

type AgentInfo struct {
	Endpoint string
	CaPem    string
	Cert     string
	Key      string
}

type AgentConfiguration struct {
	Address    string `envconfig:"AGENT_ADDRESS" default:"localhost"`
	Port       int    `envconfig:"AGENT_PORT" default:"9091"`
	CaCert     string
	Cert       string
	Key        string
	TlsEnabled bool
}

type ClusterDetails struct {
	ClusterName string
	Endpoint    string
}
