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

type ClusterDetails struct {
	ClusterName string
	Endpoint    string
}
