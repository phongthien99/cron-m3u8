package main

import (
	src "example/src"
	"example/src/config"
	"example/src/module/example-http/workflow"
	"log"

	_ "github.com/sigmaott/gest/package/technique/version"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func workerProssesor() {
	c, err := client.Dial(client.Options{
		HostPort: config.GetConfiguration().Temporal.HostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workerOptions := worker.Options{
		EnableSessionWorker: true, // Important for a worker to participate in the session
	}
	w := worker.New(c, "fileprocessing", workerOptions)

	w.RegisterWorkflow(workflow.SampleFileProcessingWorkflow)
	w.RegisterActivity(&workflow.Activities{})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
func main() {

	//rebuild
	go workerProssesor()
	app := src.NewApp(config.GetConfiguration())
	app.Run()

}
