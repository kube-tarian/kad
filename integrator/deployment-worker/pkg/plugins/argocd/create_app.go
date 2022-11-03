package argocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/model"
)

func (a *ArgoCDCLient) Create(payload model.RequestPayload) error {
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

	body := `
	{
		"metadata":{"name":"{{.ReleaseName}}"},
		"spec":{
			"destination":{
				"name":"",
				"namespace":"{{.Namespace}}",
				"server":"https://kubernetes.default.svc"
			},
			"source":{
				"path":"{{.ChartName}}",
				"repoURL":"{{.RepoURL}}",
				"targetRevision":"HEAD"
			},
			"syncPolicy":{
				"automated":{
					"prune":false,
					"selfHeal":false
				}
			},
			"project":"default"
		}
	}
	`
	buf := &bytes.Buffer{}
	tmpl, err := template.New("payload").Parse(body)
	if err != nil {
		a.logger.Errorf("Failed to create template payload, %v", err)
		return err
	}
	err = tmpl.Execute(buf, req)
	if err != nil {
		a.logger.Errorf("Failed to execute template payload, %v", err)
		return err
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", a.conf.ServiceURL, ApplicationsAPIPath), buf)
	if err != nil {
		return fmt.Errorf("request preparation to create app %s failed, %v", req.ReleaseName, err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := a.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("create application %s failed, %v", req.ReleaseName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("create app %s failed, status code: %v, %v", req.ReleaseName, resp.StatusCode, err)
	}
	return nil
}
