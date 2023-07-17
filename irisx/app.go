package irisx

import (
	"fmt"
	"runtime/debug"

	"github.com/daqiancode/std"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/hero"
)

type NewAppConfig struct {
	IsDev             bool
	API_V1            string
	ViewDir           string
	ViewSuffix        string
	StaticDir         []string
	LocaleDir         string
	LocaleDefaultLang string
}

func NewApp(conf NewAppConfig) *iris.Application {
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
	if conf.API_V1 == "" {
		conf.API_V1 = "/v1"
	}

	context.WriteJSON = func(ctx *context.Context, v interface{}, options *context.JSON) error {
		bs, err := jsoniter.Marshal(v)
		if err != nil {
			return err
		}
		_, err = ctx.Write(bs)
		return err
	}
	app := iris.New()

	app.Validator = validator.New()
	app.Use(RecoverFilter)
	app.Get(conf.API_V1+"/healthz", func(ctx *context.Context) { ctx.WriteString("OK") })

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
	app.ConfigureContainer().UseResultHandler(func(next hero.ResultHandler) hero.ResultHandler {
		return func(ctx iris.Context, v interface{}) error {
			switch val := v.(type) {
			case std.Result, *std.Result:
				return next(ctx, val)
			}
			return next(ctx, std.Result{Data: v})
		}
	})
	return app

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
