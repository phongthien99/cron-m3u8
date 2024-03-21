package src

import (
	"example/src/config/types"
	example_http "example/src/module/example-http"

	"example/src/locales/loader"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/vi"
	"github.com/go-swagno/swagno"
	"github.com/labstack/echo/v4"
	"github.com/sigmaott/gest/package/extension/configfx"
	"github.com/sigmaott/gest/package/extension/echofx"
	"github.com/sigmaott/gest/package/extension/i18nfx"
	"github.com/sigmaott/gest/package/extension/logfx"
	"github.com/sigmaott/gest/package/extension/swaggerfx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewApp(config *types.EnvironmentVariable) *fx.App {
	return fx.New(
		//config
		configfx.ForRoot(configfx.Option[*types.EnvironmentVariable]{
			Variable: config,
		}),
		// http
		echofx.ForRoot(config.Http.Port),
		fx.Provide(
			fx.Annotate(
				SetGlobalPrefix,
				fx.ParamTags(`name:"platformEcho"`),
			)),
		fx.Invoke(func(*echo.Echo) {}),
		// log
		logfx.ForRoot(config.Log.Level),
		// i18n
		i18nfx.ForRoot(i18nfx.I18nModuleParams{
			FallbackLanguage: "en",
			Loader:           loader.NewI18nMemoryLoader(),
			Translators: []locales.Translator{
				en.New(),
				vi.New(),
			},
		}),
		// swagger
		swaggerfx.ForRoot(swagno.Config{
			Path: config.Http.BasePath,
		}),
		// helper
		fx.Invoke(EnableSwagger),

		// my module
		example_http.Module(),
	)

}

func SetGlobalPrefix(e *echo.Echo, config configfx.IConfig[*types.EnvironmentVariable]) *echo.Group {

	// i18nInterceptorInstance := i18nInterceptor.NewI18nInterceter(i18nService, universalTranslator.GetFallback().Locale())
	// e.Use(
	// 	interceptor.UseInterceptors(i18nInterceptorInstance.Interceptor()),
	// 	exceptionFilter.UseFilters(
	// 		i18nValidate.ValidateExceptionFilter(),
	// 		sigmaCommon.AllExceptionFilter,
	// 	),
	// )
	return e.Group(config.Get().Http.BasePath)
}

func EnableSwagger(e *echo.Group, logger *zap.SugaredLogger, sw *swagno.Swagger) {
	sw.Host = ""
	// e.GET("/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
	// 	f := func() []byte {
	// 		return sw.MustToJson()
	// 	}
	// 	c.DynamicLoad = &f
	// }))
}
