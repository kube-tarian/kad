package argocd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	JSONContentType     = "application/json"
	ApplicationsAPIPath = "api/v1/applications"
)

func (a *ArgoCDCLient) List() (json.RawMessage, error) {
	token, err := a.getToken()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", a.conf.ServiceURL, ApplicationsAPIPath), nil)
	if err != nil {
		return nil, fmt.Errorf("request preparation to list applications failed, %v", err)
	}
	request.Header["Authorization"] = []string{fmt.Sprintf("Bearer %s", token)}
	request.Header["Content-Type"] = []string{JSONContentType}

	resp, err := a.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("list applications failed, %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list applications failed, status code: %v", resp.StatusCode)
	}
	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the application list information, %v", err)
	}
	return msg, nil
}
