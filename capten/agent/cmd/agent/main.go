package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/agent"
	captenstore "github.com/kube-tarian/kad/capten/agent/pkg/capten-store"
	"github.com/kube-tarian/kad/capten/agent/pkg/config"
	"github.com/kube-tarian/kad/capten/agent/pkg/job"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/agent/pkg/util"
	dbinit "github.com/kube-tarian/kad/capten/common-pkg/cassandra/db-init"
	dbmigrate "github.com/kube-tarian/kad/capten/common-pkg/cassandra/db-migrate"
	"github.com/pkg/errors"
	"google.golang.org/grpc/reflection"
)

var (
	log         = logging.NewLogger()
	StrInterval = "@every %ss"
)

func main() {
	log.Infof("Staring Agent")

	cfg, err := config.GetServiceConfig()
	if err != nil {
		log.Fatalf("service config reading failed, %v", err)
	}

	if err := configureDB(); err != nil {
		log.Fatalf("%v", err)
	}

	as, err := captenstore.NewStore(log)
	if err != nil {
		// ignoring store failure until DB user creation working
		// return nil, err
		log.Errorf("failed to initialize store, %v", err)
	}

	s, err := agent.NewAgent(log, cfg, as)
	if err != nil {
		log.Fatalf("Agent initialization failed, %v", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var grpcServer *grpc.Server
	if cfg.AuthEnabled {
		log.Info("Agent Authentication enabled")
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(s.AuthInterceptor))
	} else {
		log.Info("Agent Authentication disabled")
		grpcServer = grpc.NewServer()
	}
	agentpb.RegisterAgentServer(grpcServer, s)
	captenpluginspb.RegisterCaptenPluginsServer(grpcServer, s)

	log.Infof("Agent listening at %v", listener.Addr())
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Errorf("Failed to start agent : %v", err)
		}
	}()

	jobScheduler, err := initializeJobScheduler(cfg, as)
	if err != nil {
		log.Fatalf("Failed to create cron job: %v", err)
	}

	jobScheduler.Start()
	defer jobScheduler.Stop()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	grpcServer.Stop()
	log.Debugf("Exiting Agent")
}

func configureDB() error {
	if err := util.SyncCassandraAdminSecret(log); err != nil {
		return errors.WithMessage(err, "error in update cassandra secret to vault")
	}

	if err := dbinit.CreatedDatabase(log); err != nil {
		return errors.WithMessage(err, "error creating database")
	}

	if err := dbmigrate.RunMigrations(log, dbmigrate.UP); err != nil {
		return errors.WithMessage(err, "error in migrating cassandra DB")
	}
	return nil
}

func initializeJobScheduler(cfg *config.SericeConfig, as *captenstore.Store) (*job.Scheduler, error) {
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

	log.Info("successfully initialized job scheduler")
	return s, nil
}
