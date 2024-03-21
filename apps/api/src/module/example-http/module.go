package example_http

import (
	"example/src/module/example-http/controller"

	"github.com/sigmaott/gest/package/extension/echofx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("example",
		fx.Provide(
			echofx.AsRoute(controller.NewExampleController),
		),
		// fx.Provide(
		// 	service.NewExampleService,
		// ),
		// fx.Provide(
		// 	repository.NewExampleRepository,
		// ),
	)
}
