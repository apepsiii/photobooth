package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mastergin/core"
	"mastergin/middleware"
)

// RegisterRoutes registers all authentication and dashboard routes on the
// given Gin engine. Public routes (/login, /register) are accessible without
// authentication; protected routes (/dashboard, /logout) require a valid JWT.
func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *core.Config) {
	svc := NewService(db, cfg)
	h := NewHandler(svc, cfg, db)

	// Public routes — no authentication required.
	public := r.Group("/")
	{
		public.GET("/login", h.LoginPage)
		public.POST("/login", h.Login)
		public.GET("/register", h.RegisterPage)
		public.POST("/register", h.Register)
	}

	// Protected routes — RequireAuth middleware enforces a valid JWT cookie.
	private := r.Group("/")
	private.Use(middleware.RequireAuth(cfg))
	{
		private.GET("/dashboard", h.Dashboard)
		private.GET("/orders", h.Orders)
		private.GET("/orders/new", h.OrderNew)
		private.GET("/orders/:id", h.OrderDetail)
		private.GET("/products", h.Products)
		private.GET("/products/new", h.ProductNew)
		private.GET("/customers", h.Customers)
		private.GET("/profile", h.Profile)
		private.GET("/reports", h.Reports)
		private.GET("/settings", h.Settings)
		private.GET("/403", h.Forbidden)
		private.POST("/logout", h.Logout)
	}

	// Public creative landing page
	r.GET("/", h.LandingPage)

	// 404 handler for unmatched routes.
	r.NoRoute(h.NotFound)
}
