package app

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/intelops/go-common/logging"
	ginapiserver "github.com/kube-tarian/kad/capten/agent/gin-api-server"
	agentapi "github.com/kube-tarian/kad/capten/agent/internal/api"
	"github.com/kube-tarian/kad/capten/agent/internal/config"
	"github.com/kube-tarian/kad/capten/agent/internal/job"
	"github.com/kube-tarian/kad/capten/agent/internal/job/defaultplugindeployer"
	captenstore "github.com/kube-tarian/kad/capten/common-pkg/capten-store"
	"github.com/kube-tarian/kad/capten/common-pkg/crossplane"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/agentpb"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/clusterpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/pluginstorepb"
	pluginstore "github.com/kube-tarian/kad/capten/common-pkg/plugin-store"
	dbinit "github.com/kube-tarian/kad/capten/common-pkg/postgres/db-init"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	log         = logging.NewLogger()
	StrInterval = "@every %ss"
)

func Start() {
	log.Infof("Staring Agent")

	cfg, err := config.GetServiceConfig()
	if err != nil {
		log.Fatalf("service config reading failed, %v", err)
	}

	if err := configurePostgresDB(); err != nil {
		log.Fatalf("%v", err)
	}

	as, err := captenstore.NewStore(log)
	if err != nil {
		log.Errorf("failed to initialize store, %v", err)
		return
	}

	rpcapi, err := agentapi.NewAgent(log, cfg, as)
	if err != nil {
		log.Fatalf("Agent initialization failed, %v", err)
		return
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return
	}

	var grpcServer *grpc.Server
	if cfg.AuthEnabled {
		log.Info("Agent Authentication enabled")
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(rpcapi.AuthInterceptor))
	} else {
		log.Info("Agent Authentication disabled")
		grpcServer = grpc.NewServer()
	}
	agentpb.RegisterAgentServer(grpcServer, rpcapi)
	captenpluginspb.RegisterCaptenPluginsServer(grpcServer, rpcapi)
	clusterpluginspb.RegisterClusterPluginsServer(grpcServer, rpcapi)
	pluginstorepb.RegisterPluginStoreServer(grpcServer, rpcapi)

	log.Infof("Agent listening at %v", listener.Addr())
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Errorf("Failed to start agent : %v", err)
		}
	}()

	err = k8s.SetupCACertIssuser(cfg.ClusterCAIssuerName, log)
	if err != nil {
		log.Fatalf("Failed to setupt CA Cert Issuer in cert-manager %v", err)
	}

	go func() {
		err := ginapiserver.StartRestServer(rpcapi, cfg, log)
		if err != nil {
			log.Errorf("Failed to start REST server, %v", err)
		}
	}()

	err = registerK8SWatcher(as)
	if err != nil {
		log.Fatalf("Failed to initialize k8s watchers %v", err)
	}

	jobScheduler, err := initializeJobScheduler(cfg, as, rpcapi)
	if err != nil {
		log.Fatalf("Failed to create cron job: %v", err)
		return
	}

	jobScheduler.Start()
	defer jobScheduler.Stop()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	grpcServer.Stop()
	log.Debugf("Exiting Agent")
}

func configurePostgresDB() error {
	if err := dbinit.CreatedDatabase(log); err != nil {
		return errors.WithMessage(err, "error in creating postgres database")
	}

	if err := dbinit.RunMigrations(dbinit.UP); err != nil {
		return errors.WithMessage(err, "error in migrating postgres database")
	}
	return nil
}

func initializeJobScheduler(
	cfg *config.SericeConfig,
	as *captenstore.Store,
	handler pluginstore.PluginDeployHandler,
) (*job.Scheduler, error) {
	s := job.NewScheduler(log)
	if cfg.CrossplaneSyncJobEnabled {
		cs, err := job.NewCrossplaneResourcesSync(log, cfg.CrossplaneSyncJobInterval, as)
		if err != nil {
			log.Fatal("failed to init crossplane resources sync job", err)
		}
		err = s.AddJob("crossplane-resources-synch", cs)
		if err != nil {
			log.Fatal("failed to add crossplane resources sync job", err)
		}
	}

	// Add Default plugin deployer job
	addDefualtPluginsDeployerJob(s, as, handler)

	log.Info("successfully initialized job scheduler")
	return s, nil
}

func addDefualtPluginsDeployerJob(
	s *job.Scheduler,
	as *captenstore.Store,
	handler pluginstore.PluginDeployHandler,
) {
	dpd, err := defaultplugindeployer.NewDefaultPluginsDeployer(log, "@every 10m", as, handler)
	if err != nil {
		log.Fatal("failed to init default plugins deployer job", err)
	}
	err = s.AddJob("default-plugsin-deployer", dpd)
	if err != nil {
		log.Fatal("failed to add defualt plugins deployer job", err)
	}
}

func registerK8SWatcher(dbStore *captenstore.Store) error {
	if err := crossplane.RegisterK8SWatcher(log, dbStore); err != nil {
		return err
	}

	return nil
}
