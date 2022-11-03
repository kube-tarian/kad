package argocd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/model"
)

func (a *ArgoCDCLient) Delete(payload model.RequestPayload) error {
	req := &model.Request{}
	err := json.Unmarshal(payload.Data, req)
	if err != nil {
		a.logger.Errorf("payload unmarshal failed, %v", err)
		return err
	}

	token, err := a.getToken()
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s", a.conf.ServiceURL, ApplicationsAPIPath, req.ReleaseName), nil)
	if err != nil {
		return fmt.Errorf("request preparation to delete app %s failed, %v", req.ReleaseName, err)
	}
	request.Header["Authorization"] = []string{fmt.Sprintf("Bearer %s", token)}

	resp, err := a.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("delete application %s failed, %v", req.ReleaseName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete app %s failed, %v", req.ReleaseName, err)
	}
	return nil
}
