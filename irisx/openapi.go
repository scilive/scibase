package irisx

import (
	"fmt"

	"github.com/flosch/pongo2"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/scilive/scibase/logs"
	"gopkg.in/yaml.v3"
)

type OpenAPIConfig struct {
	Prefix           string
	AuthorizationUrl string
	TokenUrl         string
	ViewDir          string
	OpenAPIDir       string
}

func SetupOpenAPIControllers(app *iris.Application, config OpenAPIConfig) {
	if config.Prefix == "" {
		config.Prefix = "/v1"
	}
	if config.ViewDir == "" {
		config.ViewDir = "views"
	}
	if config.OpenAPIDir == "" {
		config.OpenAPIDir = "openapi"
	}
	prefix := config.Prefix
	logs.Log.Info().Msg("setup openapi controllers. path: " + prefix + "/openapi")
	app.Get(prefix+"/openapi", func(ctx *context.Context) {
		data := map[string]any{
			"jsonUrl": prefix + "/openapi/json",
			"prefix":  prefix,
		}
		ctx.View(config.OpenAPIDir+"/openapi", data)
	})

	app.Get(prefix+"/openapi/json", func(ctx *context.Context) {
		tpl, err := pongo2.FromFile(fmt.Sprintf("%s/%s/openapi.yml", config.ViewDir, config.OpenAPIDir))
		if err != nil {
			logs.Log.Error().Err(err).Msg("failed to load openapi.yml")
			ctx.StopWithText(500, err.Error())
			return
		}
		data := map[string]any{
			"authorizationUrl": config.AuthorizationUrl,
			"tokenUrl":         config.TokenUrl,
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
		ctx.View(config.OpenAPIDir + "/oauth2-redirect")
	})
}
