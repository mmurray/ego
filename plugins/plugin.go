package plugins

type Plugin struct {
	OnStart func()
	OnReady func()
	OnStop func()
}