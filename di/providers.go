package di

import (
	"log"
	"net/http"

	"faizalmaulana/lsp/conf"
	handler "faizalmaulana/lsp/http/hanlder"
	"faizalmaulana/lsp/http/middleware"
	"faizalmaulana/lsp/http/services"
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
func ProvideItemsRepo(db *gorm.DB) repo.ItemsRepo       { return repo.NewGormItemsRepo(db) }
func ProvideTransactionsRepo(db *gorm.DB) repo.TransactionsRepo {
	return repo.NewGormTransactionsRepo(db)
}

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
func ProvideItemsService(r repo.ItemsRepo) services.ItemsService {
	return services.NewItemsService(r)
}
func ProvideTransactionsService(r repo.TransactionsRepo) services.TransactionsService {
	return services.NewTransactionsService(r)
}

// Handlers
// Handler providers
func ProvideAuthenticationHandler(s services.AuthenticationService, sess services.SessionService, cfg *conf.Config) *handler.AuthenticationHandler {
	return handler.NewAuthenticationHandler(s, sess, cfg)
}

func ProvideUsersHandler(cfg *conf.Config, profile services.ProfilesService, users services.UsersService) *handler.UsersHandler {
	return handler.NewUsersHandler(cfg, profile, users)
}

func ProvideItemsHandler(cfg *conf.Config, items services.ItemsService) *handler.ItemsHandler {
	return handler.NewItemsHandler(cfg, items)
}

func ProvideReportHandler(cfg *conf.Config, tx services.TransactionsService) *handler.ReportHandler {
	return handler.NewReportHandler(cfg, tx)
}

// Register routes on the router
// ProvideRouterWithRoutes builds the router and registers routes (single *gin.Engine provider).
func ProvideRouterWithRoutes(ah *handler.AuthenticationHandler, uh *handler.UsersHandler, ih *handler.ItemsHandler, rh *handler.ReportHandler) *gin.Engine {
	r := ProvideRouter()
	api := r.Group("/api")
	ah.Register(api)
	uh.Register(api)
	ih.Register(api)
	rh.Register(api)

	for _, rt := range r.Routes() {
		log.Printf("route: %s %s", rt.Method, rt.Path)
	}
	return r
}

// Wire provider sets (grouped for cleaner injector definitions)
var (
	ConfigSet  = wire.NewSet(ProvideEnvConfig, ProvideDB)
	RepoSet    = wire.NewSet(ProvideUsersRepo, ProvideProfilesRepo, ProvideSessionsRepo, ProvideItemsRepo, ProvideTransactionsRepo)
	ServiceSet = wire.NewSet(ProvideAuthenticationService, ProvideSessionService, ProvideUsersService, ProvideProfilesService, ProvideItemsService, ProvideTransactionsService)
	HandlerSet = wire.NewSet(ProvideAuthenticationHandler, ProvideUsersHandler, ProvideItemsHandler, ProvideReportHandler)
	RouterSet  = wire.NewSet(ProvideRouterWithRoutes)
	ServerSet  = wire.NewSet(ProvideHTTPServer)
)
