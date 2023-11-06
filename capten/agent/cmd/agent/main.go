package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/robfig/cron/v3"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/agent"
	"github.com/kube-tarian/kad/capten/agent/pkg/config"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	captenagentsync "github.com/kube-tarian/kad/capten/agent/pkg/sync"
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

	s, err := agent.NewAgent(log, cfg)
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

	go startClaimSync(cfg.CronInterval)

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

func startClaimSync(crontInterval string) {
	cronJob := cron.New()

	fetch, err := captenagentsync.NewFetch()
	if err != nil {
		log.Errorf("Failed to initialize the sync: %v", err)
	}

	_, jobErr := cronJob.AddJob(crontInterval, cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).Then(fetch))
	if jobErr != nil {
		log.Errorf("Failed to add cronJob for sync clusterClaim: %v", jobErr)
	}

	cronJob.Start()
	defer cronJob.Stop()

	log.Info("syncing clusterClaim started successfully...")
}
