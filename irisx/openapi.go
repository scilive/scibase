package irisx

import (
	"github.com/daqiancode/env"
	"github.com/flosch/pongo2"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/scilive/scibase/logs"
	"gopkg.in/yaml.v3"
)

func SetupOpenAPIControllers(app *iris.Application, prefix string) {
	logs.Log.Info().Msg("setup openapi controllers. path: " + prefix + "/openapi")
	app.Get(prefix+"/openapi", func(ctx *context.Context) {
		data := map[string]any{
			"jsonUrl": prefix + "/openapi/json",
			"prefix":  prefix,
		}
		ctx.View("openapi/openapi", data)
	})

	app.Get(prefix+"/openapi/json", func(ctx *context.Context) {
		tpl, err := pongo2.FromFile("views/openapi/openapi.yaml")
		if err != nil {
			logs.Log.Error().Err(err).Msg("failed to load openapi.yaml")
			ctx.StopWithText(500, err.Error())
			return
		}
		data := map[string]any{
			"authorizationUrl": env.Get("OAUTH_AUTHORIZE_ENDPOINT"),
			"tokenUrl":         env.Get("OAUTH_TOKEN_ENDPOINT"),
		}
		content, err := tpl.Execute(data)
		if err != nil {
			logs.Log.Error().Err(err).Msg("failed to render openapi")
			ctx.StopWithText(500, err.Error())
			return
		}
		var r any
		err = yaml.Unmarshal([]byte(content), &r)
		if err != nil {
			ctx.StopWithText(500, "Parse openapi.yml error. error info:"+err.Error())
			return
		}
		body, err := jsoniter.Marshal(r)
		if err != nil {
			ctx.StopWithText(500, "Marshal openapi.yml error. error info:"+err.Error())
			return
		}
		ctx.ContentType("application/json")
		_, err = ctx.Write(body)
		if err != nil {
			logs.Log.Error().Err(err).Msg("failed to render openapi")
		}
	})
	app.Get(prefix+"/openapi/oauth2-redirect", func(ctx *context.Context) {
		ctx.View("openapi/oauth2-redirect")
	})
}
