package main

import (
	"context"
)

var ExportPlugin = Plugin{
	BasePlugin: &BasePlugin{
		Name:        "example",
		Description: "An example plugin",
		Version:     "v0.1.0",
	},
}

type BasePlugin struct {
	Name        string
	Description string
	Version     string
	Components  PluginComponents
}

type PluginComponents interface {
	GetLogger() any
	GetUtil() any
	GetComponents() any
}

func (p *BasePlugin) Load(components any) {
	p.Components = components.(PluginComponents)
}

type Plugin struct {
	*BasePlugin
}

func (p Plugin) GetName() string {
	return p.Name
}

func (p Plugin) GetDescription() string {
	return p.Description
}

type Logger interface {
	WarnCtx(ctx context.Context, format string, args ...any)
	ErrorCtx(ctx context.Context, format string, args ...any)
	InfoCtx(ctx context.Context, format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
	Info(format string, args ...any)
}

type Components interface {
	Get(name string) any
}

type Component interface {
	Name() string
	Service() any
}

type HttpContext interface {
	Query(key string) string
	JSON(code int, obj any)
}

type HttpRouter interface {
	ReplaceHandler(method, path string, handler func(ctx context.Context)) error
	GetHandler(method, path string) (func(ctx context.Context), error)
	GetHandlerName(method, path string) (string, error)
}

func (p Plugin) Run(args any) {
	ctx := args.(HttpContext)
	logger := p.Components.GetLogger().(Logger)
	logger.Info("Plugin %s is running, ctx %+v", p.Name, ctx)

	router := p.Components.GetComponents().(Components).Get("ginengine").(Component).Service().(HttpRouter)
	router.ReplaceHandler("GET", "/", func(ctx context.Context) {
		c := ctx.(HttpContext)
		c.JSON(200, map[string]string{"message": "Hello from plugin 2!!!"})
	})
	ctx.JSON(200, map[string]string{"message": "Plugin executed successfully"})
}

func (p Plugin) Methods() map[string]func(any) {
	return nil
}
