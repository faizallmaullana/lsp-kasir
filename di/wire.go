//go:build wireinject
// +build wireinject

package di

import (
	"net/http"

	"faizalmaulana/lsp/models/repo"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// App is an aggregate of application objects returned by the injector.
type App struct {
	Server *http.Server
	Router *gin.Engine
}

// Repos groups repository instances.
type Repos struct {
	Users    repo.UsersRepo
	Profiles repo.ProfilesRepo
}

// InitializeApp wires dependencies and returns an *App.
func InitializeApp() *App {
	panic(wire.Build(
		ConfigSet,
		RepoSet,
		ServiceSet,
		HandlerSet,
		RouterSet,
		ServerSet,
		wire.Struct(new(App), "Server", "Router"),
	))
}

// InitializeRepos wires and returns *Repos.
func InitializeRepos() *Repos {
	panic(wire.Build(
		ConfigSet,
		RepoSet,
		wire.Struct(new(Repos), "Users", "Profiles"),
	))
}

// InitializeServer is a convenience wrapper returning only the server.
func InitializeServer() *http.Server {
	app := InitializeApp()
	return app.Server
}
