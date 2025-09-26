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
func ProvidePivotItemsToTransactionsRepo(db *gorm.DB) repo.PivotItemsToTransactionsRepo {
	return repo.NewGormPivotItemsToTransactionsRepo(db)
}
func ProvideImagesRepo(db *gorm.DB) repo.ImagesRepo { return repo.NewGormImagesRepo(db) }

// Services
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
func ProvideImagesService(r repo.ImagesRepo) services.ImagesService {
	return services.NewImagesService(r)
}

// Handlers
func ProvideAuthenticationHandler(s services.AuthenticationService, sess services.SessionService, cfg *conf.Config) *handler.AuthenticationHandler {
	return handler.NewAuthenticationHandler(s, sess, cfg)
}

func ProvideUsersHandler(cfg *conf.Config, profile services.ProfilesService, users services.UsersService) *handler.UsersHandler {
	return handler.NewUsersHandler(cfg, profile, users)
}

func ProvideItemsHandler(cfg *conf.Config, items services.ItemsService, images services.ImagesService) *handler.ItemsHandler {
	return handler.NewItemsHandler(cfg, items, images)
}

func ProvideReportHandler(cfg *conf.Config, tx services.TransactionsService, pivot repo.PivotItemsToTransactionsRepo, items repo.ItemsRepo) *handler.ReportHandler {
	return handler.NewReportHandler(cfg, tx, pivot, items)
}

func ProvideTransactionsHandler(cfg *conf.Config, tx services.TransactionsService, items repo.ItemsRepo, pivot repo.PivotItemsToTransactionsRepo) *handler.TransactionsHandler {
	return handler.NewTransactionsHandler(cfg, tx, items, pivot)
}

func ProvideImagesHandler(cfg *conf.Config, svc services.ImagesService) *handler.ImagesHandler {
	return handler.NewImagesHandler(cfg, svc)
}

func ProvideRouterWithRoutes(ah *handler.AuthenticationHandler, uh *handler.UsersHandler, ih *handler.ItemsHandler, th *handler.TransactionsHandler, rh *handler.ReportHandler, imh *handler.ImagesHandler) *gin.Engine {
	r := ProvideRouter()
	api := r.Group("/api")
	ah.Register(api)
	uh.Register(api)
	ih.Register(api)
	th.Register(api)
	rh.Register(api)
	imh.Register(api)

	for _, rt := range r.Routes() {
		log.Printf("route: %s %s", rt.Method, rt.Path)
	}
	return r
}

var (
	ConfigSet  = wire.NewSet(ProvideEnvConfig, ProvideDB)
	RepoSet    = wire.NewSet(ProvideUsersRepo, ProvideProfilesRepo, ProvideSessionsRepo, ProvideItemsRepo, ProvideTransactionsRepo, ProvidePivotItemsToTransactionsRepo, ProvideImagesRepo)
	ServiceSet = wire.NewSet(ProvideAuthenticationService, ProvideSessionService, ProvideUsersService, ProvideProfilesService, ProvideItemsService, ProvideTransactionsService, ProvideImagesService)
	HandlerSet = wire.NewSet(ProvideAuthenticationHandler, ProvideUsersHandler, ProvideItemsHandler, ProvideTransactionsHandler, ProvideReportHandler, ProvideImagesHandler)
	RouterSet  = wire.NewSet(ProvideRouterWithRoutes)
	ServerSet  = wire.NewSet(ProvideHTTPServer)
)
