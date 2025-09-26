//go:build wireinject
// +build wireinject

package di

import (
	"net/http"

	"faizalmaulana/lsp/models/repo"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type App struct {
	Server *http.Server
	Router *gin.Engine
}

type Repos struct {
	Users    repo.UsersRepo
	Profiles repo.ProfilesRepo
}

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

func InitializeRepos() *Repos {
	panic(wire.Build(
		ConfigSet,
		RepoSet,
		wire.Struct(new(Repos), "Users", "Profiles"),
	))
}

func InitializeServer() *http.Server {
	app := InitializeApp()
	return app.Server
}
