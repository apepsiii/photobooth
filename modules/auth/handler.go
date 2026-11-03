package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mastergin/core"
	"mastergin/middleware"
	"mastergin/modules/settings"
)

// Handler exposes HTTP handlers for authentication and the dashboard.
type Handler struct {
	svc *Service
	cfg *core.Config
	db  *gorm.DB
}

// NewHandler creates a new auth Handler.
func NewHandler(svc *Service, cfg *core.Config, db *gorm.DB) *Handler {
	return &Handler{svc: svc, cfg: cfg, db: db}
}

// LoginPage renders the login form.
func (h *Handler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login", core.PageData{
		Title: "Sign in · Meridian",
	})
}

// Login processes the login form submission. On success it sets an HTTP-only
// cookie with the JWT and redirects to /dashboard.
func (h *Handler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	token, err := h.svc.Login(email, password)
	if err != nil {
		c.HTML(http.StatusOK, "login", core.PageData{
			Title:     "Sign in · Meridian",
			Flash:     err.Error(),
			FlashType: "error",
		})
		return
	}

	c.SetCookie("auth_token", token, 72*3600, "/", "", false, true)
	c.Redirect(http.StatusFound, "/admin/sessions")
}

// RegisterPage renders the registration form.
func (h *Handler) RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register", core.PageData{
		Title: "Create account · Meridian",
	})
}

// Register processes the registration form submission.
func (h *Handler) Register(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	_, err := h.svc.Register(name, email, password)
	if err != nil {
		c.HTML(http.StatusOK, "register", core.PageData{
			Title:     "Create account · Meridian",
			Flash:     err.Error(),
			FlashType: "error",
		})
		return
	}

	c.Redirect(http.StatusFound, "/login")
}

// Logout clears the auth cookie and redirects to /login.
func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

// Dashboard renders the main photobooth overview dashboard.
func (h *Handler) Dashboard(c *gin.Context) {
	var totalSessions int64
	var totalPrints int64
	var totalDownloads int64
	var activeFrames int64
	var availableVouchers int64

	h.db.Table("sessions").Count(&totalSessions)
	h.db.Table("sessions").Select("COALESCE(SUM(printed_count), 0)").Scan(&totalPrints)
	h.db.Table("sessions").Select("COALESCE(SUM(downloaded_count), 0)").Scan(&totalDownloads)
	h.db.Table("frames").Where("is_active = ?", true).Count(&activeFrames)
	h.db.Table("vouchers").Where("is_used = ?", false).Count(&availableVouchers)

	type RecentSession struct {
		ID            uint
		SessionCode   string
		CompositePath string
		FilterApplied string
		PrintedCount  int
		CreatedAt     time.Time
	}
	var recentSessions []RecentSession
	h.db.Table("sessions").Order("id desc").Limit(6).Scan(&recentSessions)

	cfgMap := settings.GetConfig(h.db, h.cfg.Port)

	c.HTML(http.StatusOK, "dashboard", core.PageData{
		Title:    "Dashboard Ringkasan · Photobooth",
		UserName: middleware.CurrentUserName(c),
		Path:     "/dashboard",
		Data: gin.H{
			"TotalSessions":     totalSessions,
			"TotalPrints":       totalPrints,
			"TotalDownloads":    totalDownloads,
			"ActiveFrames":      activeFrames,
			"AvailableVouchers": availableVouchers,
			"RecentSessions":    recentSessions,
			"Config":            cfgMap,
			"LocalIP":           settings.DetectLocalIP(),
		},
	})
}

// Orders renders the orders list page (protected route).
func (h *Handler) Orders(c *gin.Context) {
	c.HTML(http.StatusOK, "orders", core.PageData{
		Title:    "Orders · Meridian",
		UserName: middleware.CurrentUserName(c),
		Path:     "/orders",
	})
}

// Products renders the products list page (protected route).
func (h *Handler) Products(c *gin.Context) {
	c.HTML(http.StatusOK, "products", core.PageData{
		Title:    "Products · Meridian",
		UserName: middleware.CurrentUserName(c),
		Path:     "/products",
	})
}

// Customers renders the customers list page (protected route).
func (h *Handler) Customers(c *gin.Context) {
	c.HTML(http.StatusOK, "customers", core.PageData{
		Title:    "Customers · Meridian",
		UserName: middleware.CurrentUserName(c),
		Path:     "/customers",
	})
}

// Profile renders the user profile page (protected route).
func (h *Handler) Profile(c *gin.Context) {
	c.HTML(http.StatusOK, "profile", core.PageData{
		Title:    "Profile · Meridian",
		UserName: middleware.CurrentUserName(c),
		Path:     "/profile",
	})
}

// Reports renders the reports page (protected route).
func (h *Handler) Reports(c *gin.Context) {
	c.HTML(http.StatusOK, "reports", core.PageData{
		Title:    "Reports · Meridian",
		UserName: middleware.CurrentUserName(c),
		Path:     "/reports",
	})
}

// Settings renders the settings page (protected route).
func (h *Handler) Settings(c *gin.Context) {
	c.HTML(http.StatusOK, "settings", core.PageData{
		Title:    "Settings · Meridian",
		UserName: middleware.CurrentUserName(c),
		Path:     "/settings",
	})
}

// OrderNew renders the new order form (protected route).
func (h *Handler) OrderNew(c *gin.Context) {
	c.HTML(http.StatusOK, "order-new", core.PageData{
		Title:    "New Order · Meridian",
		UserName: middleware.CurrentUserName(c),
		Path:     "/orders",
	})
}

// OrderDetail renders the order detail page (protected route).
func (h *Handler) OrderDetail(c *gin.Context) {
	c.HTML(http.StatusOK, "order-detail", core.PageData{
		Title:    "Order #10428 · Meridian",
		UserName: middleware.CurrentUserName(c),
		Path:     "/orders",
	})
}

// ProductNew renders the new product form (protected route).
func (h *Handler) ProductNew(c *gin.Context) {
	c.HTML(http.StatusOK, "product-new", core.PageData{
		Title:    "New Product · Meridian",
		UserName: middleware.CurrentUserName(c),
		Path:     "/products",
	})
}

// Forbidden renders the 403 error page.
func (h *Handler) Forbidden(c *gin.Context) {
	c.HTML(http.StatusForbidden, "403", core.PageData{
		Title: "403 · Meridian",
	})
}

// LandingPage renders the creative neobrutalism landing page for guests and students.
func (h *Handler) LandingPage(c *gin.Context) {
	cfgMap := settings.GetConfig(h.db, h.cfg.Port)

	// Fetch recent completed photos for showcase
	var recentSessions []struct {
		SessionCode   string
		CompositePath string
		CreatedAt     time.Time
	}
	h.db.Table("sessions").
		Where("status = ? AND composite_path != ''", "completed").
		Order("id desc").
		Limit(8).
		Select("session_code, composite_path, created_at").
		Scan(&recentSessions)

	var activeFrames []struct {
		ID           uint
		Name         string
		Category     string
		OverlayImage string
		TotalShots   int
		Type         string
	}
	h.db.Table("frames").
		Where("is_active = ?", true).
		Order("id asc").
		Select("id, name, category, overlay_image, total_shots, type").
		Scan(&activeFrames)

	var totalSessions int64
	var totalPrints int64
	h.db.Table("sessions").Where("status = ?", "completed").Count(&totalSessions)
	h.db.Table("sessions").Where("print_dispatched = ?", true).Count(&totalPrints)

	c.HTML(http.StatusOK, "landing", gin.H{
		"Config":         cfgMap,
		"RecentSessions": recentSessions,
		"ActiveFrames":   activeFrames,
		"TotalSessions":  totalSessions,
		"TotalPrints":    totalPrints,
	})
}

// NotFound renders the 404 error page.
func (h *Handler) NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404", core.PageData{
		Title: "404 · Meridian",
	})
}
