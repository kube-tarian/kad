package fetcher

import (
	"github.com/gocql/gocql"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
)

const (
	FetchPluginQuery = `select name, repo_name, repo_url, chart_name, namespace, release_name, version from tools where name = ?;`
)

type Configuration struct {
	ServiceURL   []string `envconfig:"CASSANDRA_SERVICE_URL" required:"true"`
	Username     string   `envconfig:"CASSANDRA_USERNAME" required:"true"`
	Password     string   `envconfig:"CASSANDRA_PASSWORD" required:"true"`
	KeyspaceName string   `envconfig:"CASSANDRA_KEYSPACE_NAME" required:"true"`
	TableName    string   `envconfig:"CASSANDRA_TABLE_NAME" required:"true"`
}

type Store struct {
	log logging.Logger

	conf    *Configuration
	session *gocql.Session
}

func NewStore(log logging.Logger) (*Store, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Errorf("Cassandra configuration detail missing, %v", err)
		return nil, err
	}

	// Create gocql client
	cluster := gocql.NewCluster(cfg.ServiceURL...)
	cluster.Keyspace = cfg.KeyspaceName
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Username,
		Password: cfg.Password,
	}
	// cluster.SslOpts = &gocql.SslOptions{
	// 	EnableHostVerification: false,
	// }

	session, err := cluster.CreateSession()
	if err != nil {
		log.Errorf("Cassandra session creation failed, %v", err)
		return nil, err
	}

	return &Store{
		log:     log,
		conf:    cfg,
		session: session,
	}, nil
}

func (s *Store) Close() {
	s.session.Close()
}

func (s *Store) FetchPluginDetails(pluginName string) (*PluginDetails, error) {
	pd := &PluginDetails{}
	// name, repo_name, repo_url, chart_name, namespace, release_name, version
	query := s.session.Query(FetchPluginQuery, pluginName)
	err := query.Scan(
		&pd.Name,
		&pd.RepoName,
		&pd.RepoURL,
		&pd.ChartName,
		&pd.Namespace,
		&pd.ReleaseName,
		&pd.Version,
	)

	if err != nil {
		s.log.Errorf("Fetch plugin details failed, %v", err)
		return nil, err
	}
	return pd, nil
}
