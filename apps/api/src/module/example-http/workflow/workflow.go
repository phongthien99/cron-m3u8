package workflow

import (
	"time"

	"go.temporal.io/sdk/temporal"

	"go.temporal.io/sdk/workflow"
)

// SampleFileProcessingWorkflow workflow definition
func SampleFileProcessingWorkflow(ctx workflow.Context, fileName string, linkMediaFile string) (err error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		HeartbeatTimeout:    600 * time.Second, // such a short timeout to make sample fail over very fast
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	processFile(ctx, fileName, linkMediaFile)
	if err != nil {
		workflow.GetLogger(ctx).Error("Workflow failed.", "Error", err.Error())
	} else {
		workflow.GetLogger(ctx).Info("Workflow completed.")
	}
	return err
}

func processFile(ctx workflow.Context, fileName string, linkMediaFile string) (err error) {
	so := &workflow.SessionOptions{
		CreationTimeout:  time.Minute,
		ExecutionTimeout: time.Minute,
	}
	sessionCtx, err := workflow.CreateSession(ctx, so)
	if err != nil {
		return err
	}
	defer workflow.CompleteSession(sessionCtx)

	var downloadedName string
	var a = NewActivitie()
	err = workflow.ExecuteActivity(sessionCtx, a.DownloadFileActivity, fileName, linkMediaFile).Get(sessionCtx, &downloadedName)
	if err != nil {
		return err
	}

	var processedFileName string
	err = workflow.ExecuteActivity(sessionCtx, a.UploadActivity, downloadedName, fileName).Get(sessionCtx, &processedFileName)
	if err != nil {
		return err
	}

	return err
}
