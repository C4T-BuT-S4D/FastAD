package main

import (
	"log"

	"github.com/c4t-but-s4d/fastad/internal/checkers"
	"github.com/c4t-but-s4d/fastad/internal/logging"
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	logging.Init()

	temporalClient, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
		Logger: logging.NewTemporalAdapter(
			logrus.WithFields(logrus.Fields{
				"component": "checkers_worker",
			}),
		),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporalClient.Close()

	checkersWorker := worker.New(temporalClient, "checkers", worker.Options{})
	checkersWorker.RegisterWorkflow(checkers.WorkflowDefinition)

	checkersWorker.RegisterActivity(checkers.CheckActivityDefinition)
	checkersWorker.RegisterActivity(checkers.PutActivityDefinition)
	checkersWorker.RegisterActivity(checkers.GetActivityDefinition)

	if err := checkersWorker.Run(worker.InterruptCh()); err != nil {
		logrus.Fatalf("Unable to start workers: %v", err)
	}
}
