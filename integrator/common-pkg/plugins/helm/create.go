package helm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	helmclient "github.com/kube-tarian/kad/integrator/common-pkg/plugins/helm/go-helm-client"
	"github.com/kube-tarian/kad/integrator/model"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/repo"
)

const (
	in_built_cluster = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvakNDQWVhZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJek1ERXhNVEl6TWpFME1Wb1hEVE16TURFd09ESXpNakUwTVZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTEE1CkpJZ3ZrKytldGsxbytBeEs4bzFyeE55MkE5eGg1L1F2Z1NxTFBRcU9odEh0TTFHUEhOdG5nN0RWMHlHbTJxR1kKVFVGWEE4b25hbjlvVUs3TTcyS0ZDUDh3Q3dzRFpWeGRBZTBPQWplNkh5OGhqYS9GcHVZMFF6c1VJOXAzNzdnNQpQeHYwc0sydVNNQXZQSXNtcXg1VFVyTmdVYWVMUmlPNGlDdXBwMDkvdXp6TGZYcFM1cWtwbWpHRjEzdjlvcHZSCnVUUjVidndJSGRvTHJYbVZoOURXUlU5bXZwalM5NHhTV3V5RDRTOTVDOFRYNStSVnhHQ0ozNlRSelhtTXJ0dGsKeUEvTzNpVFIxa0hac2dtcE9VeHhLaHRjd2tTRTdaVW9Uc2xScmlYb2piMkErRzhWTEhPb2Zsb3Vxb1RZRHNHTgphMC9EYkVRblp1RVdUNEd0czRNQ0F3RUFBYU5aTUZjd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZCZmxMZ3pKMTlNSWNVcUVRWloxWjlJY2dZa09NQlVHQTFVZEVRUU8KTUF5Q0NtdDFZbVZ5Ym1WMFpYTXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSTZkcWtXZ2pGcHZuNS9ieTJHeQpvRmtJQjJKR2d4V1IreG1MNmFhRHFiOTRFdGMvYnFKeXU2eVNZM01SQnV2c05BWUlGalp2M0YxT1R4ay9VaWVNCm03YVRCSEt1WUh3QlNhK25UUkdWRGlGaVVOWloxSGpUTEFleGJpRzNDUHVyWDhDRlhQTXJuSG5GdTBwYXlEOUwKUXUwbkxhRVR3MG1Ub1QreWZVYU4rL3gzNUd3L3ZvanhuT0dSQWtjb2xIK3NUcmFscnZLN1plUTJxTzhGNjBPdwpPOU1qdGRTWC91M3Z3RFRPNVMyWWxOanBJZEtCbmJyK2lzMkhtR3IxYVQxb01rbDRTeWxodDlCNUlFbC92UnlqCk9waWlaalhyV0UzWGtPVFZiV3h3Ty9BN3d3a0ZoR1prWE5kSjFGaENIWTFPRFMvTjI2dWJvOHdFM2dQay9IT1gKUWtnPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    server: https://192.168.73.3:8443
  name: kind-dev-cluster
contexts:
- context:
    cluster: kind-dev-cluster
    user: kind-dev-cluster
  name: kind-dev-cluster
current-context: kind-dev-cluster
kind: Config
preferences: {}
users:
- name: kind-dev-cluster
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURJVENDQWdtZ0F3SUJBZ0lJV1Z6K2hZbG56clF3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TXpBeE1URXlNekl4TkRGYUZ3MHlOREF4TVRFeU16SXhOREphTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQTBwWUNxNjg5bmk0NEc0V1UKbUt1eVBxK3YzeHNTNjRmNFZsT3VyZzN2SFlVWlJBY003djhlV253R1ByWGZvclgrZHhSQlBaR0NzRGIrU2NKOQpTOXZ2d2l5Z3pURUtXUXZubzY1a3FWdkhNT09veGcvdURyK1QwTW5zQ1dzVitJTmNCMFVmejFMbXRHOFl1VEZ0ClRiZEM4U3gxMEtLVEY4cDdUanhNUm9zdXhaU29JS3dQd2ZZYzVvbkw3VzQzSW9UTlROVi9zd2RLWElDVkE4aVAKL1ZTdk9jTkFGamlBUnBwa0pqV0Rsd2xCK20yemhib2FCK2ZPS1B3YUZDRHBuaktGMzk3UHFtenE0dlRjdUtRKwo4TDZkM1ljV1dLNWltN2dRSUIyTkJJaVQxYk1yZUxKK0NHdkJwNDJQUWlyNlFkLzJXQ1FzNXJFYnR2b0hoNEROClIrSFhIUUlEQVFBQm8xWXdWREFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RBWURWUjBUQVFIL0JBSXdBREFmQmdOVkhTTUVHREFXZ0JRWDVTNE15ZGZUQ0hGS2hFR1dkV2ZTSElHSgpEakFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBZVBUT09FSVJPLzFKRWtyYWFCR04va0RDOGZMTmkwSENlUlhSCkVLVUI0cDR6T3p3YTBLdE9waHA2TlRBbkc2akwrYUVNZ05qRFppeTdTRnZqUG9oYTNqaHBXNkh6amMzUzlEcGgKZnYwUnhQVk9CY1J1ZkRDcnJPemw3V3paVnJWcmVGekxMSnhiQ3VOOFlWeEVHOTlXSExsaVRIc2xRelpjZFMyUgpyNEZnei9DSkMzQVJKNStLV2daMzlzUzBBRi9ST3UwMXVEb1pBU0R0Tmp3TnIvMmRTV25vQXFZdGo0STRUbHhpCjVCWWkzNUlFTXA0am1haTRaNFNDc2tIOVBWeTdLMUdnU0o0eHd2VFV6RTZBQWh1VEMzL2RuUHd0TDhQQkRyT3gKUlJpWFJUbGtmbW9HTkFHejNhTlE5MUM3S1ZjWjc5RjZXYW5GQUthREVROHFETi9oSGc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb2dJQkFBS0NBUUVBMHBZQ3E2ODluaTQ0RzRXVW1LdXlQcSt2M3hzUzY0ZjRWbE91cmczdkhZVVpSQWNNCjd2OGVXbndHUHJYZm9yWCtkeFJCUFpHQ3NEYitTY0o5Uzl2dndpeWd6VEVLV1F2bm82NWtxVnZITU9Pb3hnL3UKRHIrVDBNbnNDV3NWK0lOY0IwVWZ6MUxtdEc4WXVURnRUYmRDOFN4MTBLS1RGOHA3VGp4TVJvc3V4WlNvSUt3UAp3ZlljNW9uTDdXNDNJb1ROVE5WL3N3ZEtYSUNWQThpUC9WU3ZPY05BRmppQVJwcGtKaldEbHdsQittMnpoYm9hCkIrZk9LUHdhRkNEcG5qS0YzOTdQcW16cTR2VGN1S1ErOEw2ZDNZY1dXSzVpbTdnUUlCMk5CSWlUMWJNcmVMSisKQ0d2QnA0MlBRaXI2UWQvMldDUXM1ckVidHZvSGg0RE5SK0hYSFFJREFRQUJBb0lCQUc2Uzdod1FEQjYremg5RgphTjB4YW9xWDNaVWN0amFPVXN1aGJSdGZuYXEyZEtuUHVlN1VicSs4WjlzTnpMdTNMRUtDbEM4cjlKOXFnT05pCkNFQ0kzNy9waHhXM0ptUFRhSEg5NUVVNU44Sm9CL3JYNm53OEEvV2gwUnF3Ni94dG5Ta0VGc3ZhRCtHMlpCajUKNXhiam4zYmJqWkZiakRqMXpRRXJrREdLYTZpNmpRanlmT2ZhV3RPWFJreDVHUURHMDkrSGNtYlQyYXlOOFdSNQpvSjZnNEU5UVVQTUduR3B2M2U2YzM5SHVIMU9OelhpQUp3M0JEekxhNmhhY0pWY3g0RGRUUFAzaFNsdHltZVJLCjVKWVBub0F5UXFoQ3dCT0wyTUtZNmp1b2JGbDl1aFc3d2t6ZTBCTXJzd3YrZDViTDNpU2Z5MkI5WXcrWk9FeVQKVldFcVR0RUNnWUVBMHZhZXY0SWNBQWhkZWZnbnB3L0NwSFcybkhGYmFTRHRzaWF1b2ZGVjdYR3ZBYVV0ZzVQZwpJV3ZOdkZVQWhrMkdUdVBHdlp6aFZiRElsNVRjSWU5OXFmb0JoUFBlY3pLcnZ2NHluRkx1SEpsaldCWmJJaEZOCjZxM3VlNTQ5eXdQc3pqV0QzSnd3dFlQZjZQcm1QSzlneElyY3dTV05NR3ZvVFVlQmM3Qm1KZ2NDZ1lFQS80ckUKR09uanJMTHpRd1dkNDRjanVtVHlwNFJtVlhYMXdTR0ZJdXhhTGJXQmlpWC90cHpMWmVIdWdBN2xvK01WUlZjLwprbmp6WTZQbUM0djg4bU9FSnFzZ3gwd3dFNHowOUFSaVpLQWFyMlFYMGJPUjhwNmV5TU90RFoyOVViOUh3RXhHCkt3dGtsTXg5NXhkd0NrYk1ockh5bFdIcGVzWkM2cTVlZjU0WGNMc0NnWUJhMllnTjB2czU3R0JOQ1ZnU010QlEKd0x5dWJJYkFKRVVZeGwzSU1jVWVaeW5Gbkp1WUlWT1JNUHE5a3lHUnRNc1ZLRFJMTGNkQWZzd3pzeENGc0x3KwpPZ0x6Zlk0YnNBT1VVYVg3K2g2K3hET3JHSjJRYzBGSndqT0VtdVhqaXNJdEg1QzByYkt3U0tWaGtNTWIrUzdFCkZVVHlETGpiMUd5Szh6TkZYZjd2ZXdLQmdBZWVlWTVNbXU4eFByT0czVmhGVlRsZmZTU2xlKytjWHNGdFlHelUKSXpRdHJ6a1JQUGlTNERXZmNONzhrcmc2TXc0b05jc0dOQ3VLWFhlR3F2b0hJWStObHFLYWtPeGtUWUZoQ0JYNworQSsycWtja1ZYdW9ZdytWVmZtTDlITVZndXdtMmdpNmhEc3poYVY0TzJ6ekEzSVlxQ1R3RUdnS3RVQU9CdDlECk5XdTFBb0dBVFQvOWhpSzVnY2hodGoybzZLTTlEL3FjNEgyMDdSWE42VytzUFE1Ui9XZ2huWFA0QXVKbzZIaGgKOXp0R2tUMWVoMC9OSmphN29HWDlFRi93RFJPVHBhYWNOZlRxOVVrL0svY1VaQThQdDBBQy8wR3p0emp0YTNjOApKNVl0aHdnTHlXcFlFaVJwUTJQcUVMQzgxdm5qdnh4R2dNVlNQQnBCeTB0MDFIamVJMnM9Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
`
)

func (h *HelmCLient) Create(payload model.RequestPayload) (json.RawMessage, error) {
	h.logger.Infof("Helm client Install invoke started")

	req := &model.Request{}
	err := json.Unmarshal(payload.Data, req)
	if err != nil {
		h.logger.Errorf("payload unmarshal failed, %v", err)
		return nil, err
	}

	helmClient, err := h.getHelmClient(req)
	if err != nil {
		h.logger.Errorf("helm client initialization failed, %v", err)
		return nil, err
	}

	err = h.addOrUpdate(helmClient, req)
	if err != nil {
		h.logger.Errorf("helm repo add failed, %v", err)
		return nil, err
	}

	// Use an unpacked chart directory.
	chartSpec := helmclient.ChartSpec{
		ReleaseName: req.ReleaseName,
		ChartName:   fmt.Sprintf("%s/%s", req.RepoName, req.ChartName),
		Namespace:   req.Namespace,
		Wait:        true,
		Timeout:     time.Duration(req.Timeout) * time.Minute,
	}

	// Use the default rollback strategy offer by HelmClient (revert to the previous version).
	rel, err := helmClient.InstallOrUpgradeChart(
		context.Background(),
		&chartSpec,
		&helmclient.GenericHelmOptions{
			RollBack:              helmClient,
			InsecureSkipTLSverify: true,
		})
	if err != nil {
		h.logger.Errorf("helm install or update for request %+v failed, %v", req, err)
		return nil, err
	}

	h.logger.Infof("helm install of app %s successful in namespace: %v, status: %v", rel.Name, rel.Info.Status, rel.Namespace)
	h.logger.Infof("Helm client Install invoke finished")
	return json.RawMessage(fmt.Sprintf("{\"status\": \"Application %s install successful\"}", rel.Name)), nil
}

func (h *HelmCLient) getHelmClient(req *model.Request) (helmclient.Client, error) {
	// Change this to the namespace you wish the client to operate in.
	// helmClient, err := helmclient.New(opt)

	opt := &helmclient.Options{
		Namespace:        req.Namespace,
		RepositoryCache:  "/tmp/.helmcache",
		RepositoryConfig: "/tmp/.helmrepo",
		Debug:            true,
		Linting:          true,
		DebugLog:         h.logger.Debugf,
	}

	var yamlKubeConfig interface{}
	err := yaml.Unmarshal([]byte(in_built_cluster), &yamlKubeConfig)
	if err != nil {
		h.logger.Errorf("kubeconfig yaml unmarshal failed, error: %v", err)
		return nil, err
	}

	jsonKubeConfig, err := jsoniter.Marshal(yamlKubeConfig)
	if err != nil {
		h.logger.Errorf("json Marhsal of kubeconfig failed, err: json Mashal: %v", err)
		return nil, err
	}

	return helmclient.NewClientFromKubeConf(
		&helmclient.KubeConfClientOptions{
			Options:     opt,
			KubeContext: "cluster-1",
			KubeConfig:  jsonKubeConfig,
		},
	)
}

func (h *HelmCLient) addOrUpdate(client helmclient.Client, req *model.Request) error {
	// Define a public chart repository.
	chartRepo := repo.Entry{
		Name:                  req.RepoName,
		URL:                   req.RepoURL,
		InsecureSkipTLSverify: true,
	}

	// Add a chart-repository to the client.
	if err := client.AddOrUpdateChartRepo(chartRepo); err != nil {
		h.logger.Errorf("helm repo add failed, %v", err)
		return err
	}
	return nil
}
