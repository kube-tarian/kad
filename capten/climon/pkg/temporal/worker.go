package temporal

import (
	"fmt"
)

type Worker interface {
	Start() error
	Stop() error
}

var workersInitFuncMap = map[string]func(address string) (Worker, error){
	"SyncWorker": NewSyncWorker,
}

type worker struct {
	temporalAddress string
	startedWorkers  map[string]Worker
}

func New(address string) *worker {
	return &worker{
		temporalAddress: address,
		startedWorkers:  make(map[string]Worker),
	}
}

func (w *worker) StartWorkers() error {
	for name, initFunc := range workersInitFuncMap {
		workerObj, err := initFunc(w.temporalAddress)
		if err != nil {
			return fmt.Errorf("failed to create worker: %s err: %v", name, err)
		}

		var workerError error
		go func() {
			if err := workerObj.Start(); err != nil {
				workerError = fmt.Errorf("failed to start worker: %s, err: %v", name, err)
			}
		}()

		if workerError != nil {
			return workerError
		}

		w.startedWorkers[name] = workerObj
	}

	return nil
}

func (w *worker) StopWorkers() {
	for name, workerObj := range w.startedWorkers {
		if err := workerObj.Stop(); err != nil {
			fmt.Println("error stopping worker:", name)
		}
	}
}
