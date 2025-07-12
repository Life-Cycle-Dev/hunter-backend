package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"hunter-backend/di/config"
	"hunter-backend/di/database"
	"hunter-backend/entity"
	applicationsService "hunter-backend/service/applications"
	authService "hunter-backend/service/auth"
	healthCheckService "hunter-backend/service/health_check"
	"hunter-backend/service/middleware"
	permissionService "hunter-backend/service/permission"
)

func InitRouter(server *fiber.App) {
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	appConfig := config.GetConfig()
	healthCheck := healthCheckService.ProvideHealthCheckService(db, appConfig)
	server.Get("/_hc", healthCheck.HandlerGetRouter)

	auth := authService.ProvideAuthService(db, appConfig)
	server.Post("/auth/sign-up", auth.HandlerSignUp)
	server.Post("/auth/login", auth.HandlerLogin)
	server.Post("/auth/verify-email", auth.HandlerVerifyEmail)

	authProtected := server.Group("/auth/me", middleware.RequireAuth(db, appConfig, entity.JsonWebTokenAccessToken))
	authProtected.Get("/", auth.HandlerGetUserInfo)

	refreshProtected := server.Group("/auth/refresh", middleware.RequireAuth(db, appConfig, entity.JsonWebTokenRefreshToken))
	refreshProtected.Post("/", auth.HandlerRefreshAccessToken)

	application := applicationsService.ProvideApplicationsService(db, appConfig)
	applicationProtected := server.Group("/application", middleware.RequireAuth(db, appConfig, entity.JsonWebTokenAccessToken))
	applicationProtected.Post("/create", application.HandlerCreateApplication)
	applicationProtected.Get("/list", application.HandlerListApplication)
	applicationProtected.Get("/:id", application.HandlerGetApplicationById)
	applicationProtected.Put("/:id", application.HandlerUpdateApplicationById)

	permission := permissionService.ProvidePermissionService(db, appConfig)
	permissionProtected := server.Group("/permission", middleware.RequireAuth(db, appConfig, entity.JsonWebTokenAccessToken))
	permissionProtected.Post("/create", permission.HandlerCreatePermission)
	permissionProtected.Get("/list", permission.HandlerListPermission)
	permissionProtected.Get("/:id", permission.HandlerGetPermissionById)
	permissionProtected.Put("/:id", permission.HandlerUpdatePermission)

	roleProtected := server.Group("/role", middleware.RequireAuth(db, appConfig, entity.JsonWebTokenAccessToken))
	roleProtected.Post("/create", permission.HandlerCreateRole)
	roleProtected.Get("/list", permission.HandlerListRole)
	roleProtected.Get("/:id", permission.HandlerGetRoleById)
	roleProtected.Put("/:id", permission.HandlerUpdateRole)
}
