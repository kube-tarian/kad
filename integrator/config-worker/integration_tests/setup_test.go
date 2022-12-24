package integrationtests

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/config-worker/pkg/application"
)

type TestContextData struct {
}

func setupENV() {
}

func setup() *TestContextData {
	setupENV()
	cfg := &application.Configuration{}
	if err := envconfig.Process("", cfg); err != nil {
		logger.Fatalf("Could not parse env Config: %v", err)
	}

	return &TestContextData{}
}

func tearDown(t *TestContextData) {
}

func startMain() chan bool {
	stopCh := make(chan bool)

	// Start agent and client
	go startApplication(stopCh)
	time.Sleep(2 * time.Second)

	// Wait till Agent and Client healthy
	isApplicationHealthy := false
	for {
		select {
		// wait till 1min, after that exit 1
		case <-time.After(1 * time.Minute):
			logger.Fatalf("Configuration worker application not healthy")
		case <-time.After(2 * time.Second):
			// Check Agent health
			isApplicationHealthy = getHealth(http.MethodGet, "http://localhost:9080", "status", "agent")
		}
		if isApplicationHealthy {
			break
		}
	}
	return stopCh
}

func getHealth(method, url, path, serviceName string) bool {
	resp, err := callHTTPRequest(method, url, path, nil)
	if err != nil {
		logger.Errorf("%v health check call failed: %v", serviceName, err)
		return false
	}

	return checkResponse(resp, http.StatusOK)
}

func checkResponse(resp *http.Response, statusCode int) bool {
	return resp.StatusCode == statusCode
}

func startApplication(stop chan bool) {
	os.Setenv("PORT", "9080")
	log := logging.NewLogger()
	app := application.New(log)
	go app.Start()

	<-stop
}

func callHTTPRequest(method, url, path string, body []byte) (*http.Response, error) {
	finalURL := fmt.Sprintf("%s/%s", url, path)
	var req *http.Request
	if body != nil {
		req, _ = http.NewRequest(method, finalURL, bytes.NewBuffer(body))
	} else {
		req, _ = http.NewRequest(method, finalURL, nil)
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	return client.Do(req)
}
