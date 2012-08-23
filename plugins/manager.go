package plugins

import (

)

var plugins = make([]*Plugin, 0)

func Register(p *Plugin) *Plugin {
	plugins = append(plugins, p)
	return p
}

func All() []*Plugin {
	return plugins
}
