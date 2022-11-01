package integrationtests

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/deployment-worker/pkg/application"
)

type TestContextData struct {
}

func setupENV() {
	os.Setenv("NATS_TOKEN", "UfmrJOYwYCCsgQvxvcfJ3BdI6c8WBbnD")
	os.Setenv("NATS_ADDRESS", "nats://localhost:4222")
	os.Setenv("STREAM_NAME", "CONTAINERMETRICS")
	os.Setenv("DB_ADDRESS", "localhost:9000")
}

func setup() *TestContextData {
	setupENV()
	cfg := &application.Configuration{}
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("Could not parse env Config: %v", err)
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
	isAgentHealthy := false
	isClientHealthy := false
	for {
		select {
		// wait till 1min, after that exit 1
		case <-time.After(1 * time.Minute):
			log.Fatalf("Agent/Client not healthy")
		case <-time.After(2 * time.Second):
			// Check Agent health
			isAgentHealthy = getHealth(http.MethodGet, "http://localhost:8090", "status", "agent")
			// Check Client health
			isClientHealthy = getHealth(http.MethodGet, "http://localhost:8091", "status", "client")
		}
		if isAgentHealthy && isClientHealthy {
			break
		}
	}
	return stopCh
}

func getHealth(method, url, path, serviceName string) bool {
	resp, err := callHTTPRequest(method, url, path, nil)
	if err != nil {
		log.Printf("%v health check call failed: %v", serviceName, err)
		return false
	}

	return checkResponse(resp, http.StatusOK)
}

func checkResponse(resp *http.Response, statusCode int) bool {
	return resp.StatusCode == statusCode
}

func startApplication(stop chan bool) {
	os.Setenv("PORT", "8090")
	app := application.New()
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
