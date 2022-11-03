package argocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

const TokenPath = "api/v1/session"

type TokenResponse struct {
	Token string `json:"token" required:"true"`
}

func (a *ArgoCDCLient) getToken() (string, error) {
	payloadTemplate := `{"username":"admin","password":"{{.Password}}"}`
	buf := &bytes.Buffer{}
	tmpl, err := template.New("token").Parse(payloadTemplate)
	if err != nil {
		a.logger.Errorf("Failed to create token template payload, %v", err)
		return "", err
	}
	err = tmpl.Execute(buf, a.conf)
	if err != nil {
		a.logger.Errorf("Failed to execute token template payload, %v", err)
		return "", err
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", a.conf.ServiceURL, TokenPath), buf)
	if err != nil {
		return "", fmt.Errorf("request preparation to token failed, %v", err)
	}

	resp, err := a.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("get token request failed, %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get token failed, status code:%v, %v", resp.StatusCode, err)
	}
	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to fetch token, %v", err)
	}

	token := &TokenResponse{}
	err = json.Unmarshal(msg, token)
	if err != nil {
		return "", fmt.Errorf("token not received in response, %v", err)
	}

	return token.Token, nil
}
