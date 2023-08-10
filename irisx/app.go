package irisx

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/daqiancode/std"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/kataras/iris/v12/view"
	"github.com/scilive/scibase/logs"
	"github.com/scilive/scibase/utils/strs"
)

type NewAppConfig struct {
	IsDev             bool
	Prefix            string
	ViewDir           string
	ViewSuffix        string
	StaticDir         []string
	LocaleDir         string
	LocaleDefaultLang string
}

func NewApp(conf NewAppConfig) (*iris.Application, *view.DjangoEngine) {
	if conf.ViewDir == "" {
		conf.ViewDir = "./views"
	}
	if conf.ViewSuffix == "" {
		conf.ViewSuffix = ".html"
	}
	if conf.LocaleDir == "" {
		conf.LocaleDir = "./locales/*/*"
	}
	if conf.LocaleDefaultLang == "" {
		conf.LocaleDefaultLang = "en-US"
	}
	logs.Log.Info().Msg("set json to jsoniter")
	context.WriteJSON = func(ctx *context.Context, v interface{}, options *context.JSON) error {
		bs, err := jsoniter.Marshal(v)
		if err != nil {
			return err
		}
		_, err = ctx.Write(bs)
		return err
	}
	app := iris.New()
	logs.Log.Info().Msg("set access log")
	al := accesslog.New(os.Stdout)
	al.Async = true
	al.IP = false
	al.RequestBody = false
	al.ResponseBody = false
	al.BytesReceived = false
	al.BytesSent = false
	al.BytesReceivedBody = false
	al.BytesSentBody = false
	al.AddFields(func(ctx *context.Context, f *accesslog.Fields) {
		ip := GetRealIP(ctx)
		f.Set("IP", ip)
	})
	app.UseRouter(al.Handler)
	app.Validator = validator.New()
	app.Use(RecoverFilter)

	logs.Log.Info().Msg("mount health check:" + conf.Prefix + "/health")
	app.Get(conf.Prefix+"/health", func(ctx *context.Context) { ctx.WriteString("OK") })

	tmpl := iris.Django(conf.ViewDir, conf.ViewSuffix)
	tmpl.Reload(conf.IsDev)
	app.RegisterView(tmpl)

	app.I18n.DefaultMessageFunc = func(langInput, langMatched, key string, args ...interface{}) string {
		msg := fmt.Sprintf("user language input: %s: matched as: %s: not found key: %s: args: %v", langInput, langMatched, key, args)
		app.Logger().Warn(msg)
		return "Invalid value"
	}
	err := app.I18n.Load(conf.LocaleDir, conf.LocaleDefaultLang)
	if err != nil {
		panic(err)
	}
	app.I18n.SetDefault(conf.LocaleDefaultLang)
	app.ConfigureContainer().OnError(func(ctx iris.Context, err error) {
		if err == nil || ctx.IsStopped() {
			return
		}
		err = transformValidationErrors(ctx, err)
		if e, ok := err.(*std.Error); ok {
			if e.HttpStatus == 0 {
				ctx.StatusCode(500)
			} else {
				ctx.StatusCode(e.HttpStatus)
			}
			ctx.JSON(e.ToResult(ctx.Tr))
			return
		}
		ctx.JSON(std.WrapError(err, "").ToResult(ctx.Tr))
	})
	// app.ConfigureContainer().UseResultHandler(func(next hero.ResultHandler) hero.ResultHandler {
	// 	return func(ctx iris.Context, v interface{}) error {
	// 		switch val := v.(type) {
	// 		case std.Result, *std.Result:
	// 			return next(ctx, val)
	// 		}
	// 		return next(ctx, std.Result{Data: v})
	// 	}
	// })
	return app, tmpl

}

func RecoverFilter(ctx iris.Context) {
	defer func() {
		if err := recover(); err != nil {
			stack := string(debug.Stack())
			fmt.Println(err)
			fmt.Println(stack)
			if e, ok := err.(error); ok {
				if ctx.IsStopped() { // handled by other middleware.
					return
				} else {
					ctx.StopWithError(iris.StatusInternalServerError, e)
				}
			}
			ctx.StopExecution()
		}
	}()

	ctx.Next()
}

var DefaultValidationErrorMessage = "Invalid request parameters"

func transformValidationErrors(ctx iris.Context, err error) error {
	if err != nil {

		if es, ok := err.(validator.ValidationErrors); ok {
			e := std.Err(ctx.Tr("error.invalid_param")).SetHttpStatus(iris.StatusBadRequest)
			if e.Message == "" {
				e.Message = DefaultValidationErrorMessage
			}
			for _, v := range es {
				field := strs.Decapitalize(v.Field())
				e.AddFieldError(field, "invalid value", "error."+v.ActualTag(), v.Param())
			}
			return e
		}
		return err
	}
	return nil
}
