package integrationtests

import (
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestIntegrationDeploymentEvent(t *testing.T) {
	data := setup()

	stop := startMain()

	log.Println("Sleeping now")
	time.Sleep(5 * time.Second)

	log.Println("Starting teardown")
	tearDown(data)
	stop <- true
}
