package controller

import (
	"context"
	"example/src/config"
	"example/src/module/example-http/dto"
	"example/src/module/example-http/model"
	"example/src/module/example-http/workflow"
	"log"

	"github.com/akyoto/uuid"
	"github.com/gestgo/gest/package/extension/i18nfx"
	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/labstack/echo/v4"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

type IExampleController interface {
	Create()
	// FindOne()
	// Paginate()
	// UpdateOne()
	// DeleteOne()
}

type exampleController struct {
	router      *echo.Group
	i18nService i18nfx.II18nService
	swagger     *swagno.Swagger
	logger      *zap.SugaredLogger
}

func NewExampleController(router *echo.Group, swagger *swagno.Swagger, logger *zap.SugaredLogger) IExampleController {
	return &exampleController{
		router:  router.Group("/m3u8-crawl"),
		swagger: swagger,
		logger:  logger,
	}

}

func (e *exampleController) Create() {

	e.swagger.AddEndpoint(endpoint.New(
		endpoint.POST,
		"/m3u8-crawl",
		endpoint.WithTags("jobs"),
		endpoint.WithBody(dto.M3u8JobDto{}),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(model.Example{}, "OK", "200")}),
		endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
		endpoint.WithConsume([]mime.MIME{mime.JSON}),
	))
	e.router.POST("", func(c echo.Context) error {
		payload := new(dto.M3u8JobDto)
		err := c.Bind(payload)
		if err != nil {
			e.logger.Error(err)
			return err
		}
		con, err := client.Dial(client.Options{
			HostPort: config.GetConfiguration().Temporal.HostPort,
		})
		if err != nil {
			e.logger.Error(err)
			return err
		}
		defer con.Close()

		fileID := uuid.New()
		nameSlug := payload.Name
		workflowOptions := client.StartWorkflowOptions{
			ID:        "fileprocessing_" + nameSlug + fileID.String(),
			TaskQueue: "fileprocessing",
		}

		we, err := con.ExecuteWorkflow(context.Background(), workflowOptions, workflow.SampleFileProcessingWorkflow, nameSlug, payload.LinkMediaPlaylist)
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}
		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

		return c.JSON(201, model.Example{
			Name:  nameSlug,
			RunID: we.GetRunID(),
		})
	})

}
