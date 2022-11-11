package helm

import (
	"context"
	"fmt"
	"log"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/pkg/errors"

	"github.com/kube-tarian/kad/climon/pkg/k8s"
	"github.com/kube-tarian/kad/climon/pkg/pb/climonpb"
)

const (
	DefaultTimeout = 15 * time.Minute
)

type helm struct{}

func NewHelm() *helm {
	return &helm{}
}

func (h *helm) Run(ctx context.Context, payload interface{}) error {
	settings := cli.New()
	request, ok := payload.(*climonpb.DeployRequest)
	if !ok {
		return errors.New("invalid payload")
	}

	repoEntry := &repo.Entry{
		Name: request.RepoName,
		URL:  request.RepoUrl,
	}

	k8sInfo := k8s.GetInfo()
	settings.KubeConfig = k8sInfo.GetConfigPath()
	settings.KubeToken = k8sInfo.GetToken()
	settings.KubeCaFile = k8sInfo.GetK8sCaFilePath()
	settings.KubeAPIServer = k8sInfo.GetK8sEndpoint()
	r, err := repo.NewChartRepository(repoEntry, getter.All(settings))
	if err != nil {
		return errors.Wrap(err, "failed to create new repo")
	}

	r.CachePath = settings.RepositoryCache
	_, err = r.DownloadIndexFile()
	if err != nil {
		return errors.Wrap(err, "unable to download the index file")
	}

	var repoFile repo.File
	repoFile.Update(repoEntry)
	err = repoFile.WriteFile(settings.RepositoryConfig, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write the helm-chart path")
	}

	actionConfig := new(action.Configuration)
	err = actionConfig.Init(settings.RESTClientGetter(), request.Namespace, "", debug)
	if err != nil {
		return errors.Wrap(err, "failed to setup actionConfig for helm")
	}

	client := action.NewInstall(actionConfig)
	client.Namespace = request.Namespace
	client.ReleaseName = request.ReleaseName
	client.Timeout = DefaultTimeout
	cp, err := client.ChartPathOptions.LocateChart(request.ChartName, settings)
	chartReq, err := loader.Load(cp)
	releaseInfo, err := client.Run(chartReq, nil)
	if err != nil {
		return errors.Wrap(err, "chart run error")
	}

	log.Println("release info", releaseInfo)
	fmt.Println("deploy method called", request)
	return nil
}

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v))
}

func (h *helm) Status() string {
	return "SUCCESS"
}
