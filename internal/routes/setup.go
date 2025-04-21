package routes

import (
	"github.com/StepOne-ai/pvz_avito/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/dummyLogin", DummyLogin)

	r.POST("/register", Register)

	r.POST("/login", Login)

	r.POST("/pvz",
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("PVZemployee"),
		PVZ_post)

	r.GET("/pvz",
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("PVZemployee", "Moderator"),
		PVZ_get)

	r.POST("/pvz/:pvzId/close_last_reception",
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("PVZemployee"),
		PVZ_close_last_reception)

	r.POST("/pvz/:pvzId/delete_last_product",
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("PVZemployee"),
		PVZ_delete_last_product)

	r.POST("/receptions",
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("PVZemployee"),
		Receptions)

	r.POST("/products",
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("PVZemployee"),
		Products)
}
