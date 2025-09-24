package di

import (
	"log"
	"net/http"

	"faizalmaulana/lsp/conf"
	handler "faizalmaulana/lsp/http/hanlder"
	"faizalmaulana/lsp/http/services"
	"faizalmaulana/lsp/middleware"
	"faizalmaulana/lsp/models/repo"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// Base / infrastructure providers
func ProvideEnvConfig() *conf.Config { return conf.NewEnvConfig() }

func ProvideRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	return r
}

func ProvideHTTPServer(env *conf.Config, router *gin.Engine) *http.Server {
	return &http.Server{Addr: ":" + env.Port, Handler: router}
}

func ProvideDB(env *conf.Config) *gorm.DB { return env.DB }

// Repository providers
func ProvideUsersRepo(db *gorm.DB) repo.UsersRepo       { return repo.NewGormUsersRepo(db) }
func ProvideProfilesRepo(db *gorm.DB) repo.ProfilesRepo { return repo.NewGormProfilesRepo(db) }
func ProvideSessionsRepo(db *gorm.DB) repo.SessionsRepo { return repo.NewGormSessionsRepo(db) }

// Services
// Service providers (currently only authentication is used)
func ProvideAuthenticationService(r repo.UsersRepo) services.AuthenticationService {
	return services.NewAuthenticationService(r)
}

func ProvideSessionService(r repo.SessionsRepo) services.SessionService {
	return services.NewSessionService(r)
}

func ProvideUsersService(r repo.UsersRepo) services.UsersService { return services.NewUsersService(r) }
func ProvideProfilesService(r repo.ProfilesRepo) services.ProfilesService {
	return services.NewProfilesService(r)
}

// Handlers
// Handler providers
func ProvideAuthenticationHandler(s services.AuthenticationService, sess services.SessionService, cfg *conf.Config) *handler.AuthenticationHandler {
	return handler.NewAuthenticationHandler(s, sess, cfg)
}

func ProvideUsersHandler(cfg *conf.Config, profile services.ProfilesService, users services.UsersService) *handler.UsersHandler {
	return handler.NewUsersHandler(cfg, profile, users)
}

// Register routes on the router
// ProvideRouterWithRoutes builds the router and registers routes (single *gin.Engine provider).
func ProvideRouterWithRoutes(ah *handler.AuthenticationHandler, uh *handler.UsersHandler) *gin.Engine {
	r := ProvideRouter()
	api := r.Group("/api")
	ah.Register(api)
	uh.Register(api)

	for _, rt := range r.Routes() {
		log.Printf("route: %s %s", rt.Method, rt.Path)
	}
	return r
}

// Wire provider sets (grouped for cleaner injector definitions)
var (
	ConfigSet  = wire.NewSet(ProvideEnvConfig, ProvideDB)
	RepoSet    = wire.NewSet(ProvideUsersRepo, ProvideProfilesRepo, ProvideSessionsRepo)
	ServiceSet = wire.NewSet(ProvideAuthenticationService, ProvideSessionService, ProvideUsersService, ProvideProfilesService)
	HandlerSet = wire.NewSet(ProvideAuthenticationHandler, ProvideUsersHandler)
	RouterSet  = wire.NewSet(ProvideRouterWithRoutes)
	ServerSet  = wire.NewSet(ProvideHTTPServer)
)
