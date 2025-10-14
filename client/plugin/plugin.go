package plugin

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
