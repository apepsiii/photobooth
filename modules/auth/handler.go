package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mastergin/core"
	"mastergin/middleware"
	"mastergin/modules/kiosk"
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
	h.db.Table("sessions").Where("printed_count > ?", 0).Count(&totalPrints)

	c.HTML(http.StatusOK, "landing", gin.H{
		"Config":         cfgMap,
		"RecentSessions": recentSessions,
		"ActiveFrames":   activeFrames,
		"TotalSessions":  totalSessions,
		"TotalPrints":    totalPrints,
	})
}

// AnalyticsPage renders the comprehensive financial & photobooth performance dashboard.
func (h *Handler) AnalyticsPage(c *gin.Context) {
	var totalSessions int64
	var totalPrints int64
	var totalDownloads int64
	h.db.Table("sessions").Count(&totalSessions)
	h.db.Table("sessions").Select("COALESCE(SUM(printed_count), 0)").Scan(&totalPrints)
	h.db.Table("sessions").Select("COALESCE(SUM(downloaded_count), 0)").Scan(&totalDownloads)

	// QRIS Revenue
	var qrisRevenue int64
	var qrisPaidOrders int64
	h.db.Table("payment_orders").Where("status = ?", "paid").Select("COALESCE(SUM(amount), 0)").Scan(&qrisRevenue)
	h.db.Table("payment_orders").Where("status = ?", "paid").Count(&qrisPaidOrders)

	// Voucher Stats
	var totalVouchers int64
	var usedVouchers int64
	h.db.Table("vouchers").Count(&totalVouchers)
	h.db.Table("vouchers").Where("is_used = ?", true).Count(&usedVouchers)

	cfgMap := settings.GetConfig(h.db, h.cfg.Port)
	voucherRevenue := int64(usedVouchers) * int64(cfgMap.PriceSingle)
	totalGrossRevenue := qrisRevenue + voucherRevenue

	// Paper & Consumables Estimations (assuming Rp 1,500 / photo print paper & ribbon ink cost)
	paperUsed := int64(cfgMap.PaperUsed)
	costPerSheet := int64(1500)
	estimatedCost := paperUsed * costPerSheet
	estimatedNetProfit := totalGrossRevenue - estimatedCost
	if estimatedNetProfit < 0 && totalGrossRevenue == 0 {
		estimatedNetProfit = 0
	}

	// Hourly Traffic Distribution (07 - 22)
	type HourlyStat struct {
		Hour  string `json:"hour"`
		Count int    `json:"count"`
	}
	var hourlyRaw []HourlyStat
	h.db.Raw("SELECT strftime('%H', created_at) as hour, count(*) as count FROM sessions WHERE created_at IS NOT NULL GROUP BY hour ORDER BY hour").Scan(&hourlyRaw)

	// Build 24-hour map
	hourlyMap := make(map[string]int)
	maxHourlyCount := 1
	for _, hr := range hourlyRaw {
		hourlyMap[hr.Hour] = hr.Count
		if hr.Count > maxHourlyCount {
			maxHourlyCount = hr.Count
		}
	}

	type HourBar struct {
		Label     string
		Count     int
		HeightPct int
	}
	var hourlyBars []HourBar
	for hIdx := 7; hIdx <= 22; hIdx++ {
		hStr := fmt.Sprintf("%02d", hIdx)
		cnt := hourlyMap[hStr]
		pct := int((float64(cnt) / float64(maxHourlyCount)) * 100)
		if cnt > 0 && pct < 12 {
			pct = 12
		}
		hourlyBars = append(hourlyBars, HourBar{
			Label:     hStr + ":00",
			Count:     cnt,
			HeightPct: pct,
		})
	}

	// Frame Leaderboard
	type FrameLeaderboard struct {
		FrameID    uint
		FrameName  string
		UsageCount int
	}
	var topFrames []FrameLeaderboard
	h.db.Raw(`
		SELECT f.id as frame_id, f.name as frame_name, count(s.id) as usage_count
		FROM frames f
		LEFT JOIN sessions s ON s.frame_id = f.id
		GROUP BY f.id
		ORDER BY usage_count DESC
		LIMIT 5
	`).Scan(&topFrames)

	// Filter Leaderboard
	type FilterLeaderboard struct {
		FilterApplied string
		Count         int
	}
	var topFilters []FilterLeaderboard
	h.db.Raw(`
		SELECT COALESCE(NULLIF(filter_applied, ''), 'normal') as filter_applied, count(*) as count
		FROM sessions
		GROUP BY filter_applied
		ORDER BY count DESC
		LIMIT 6
	`).Scan(&topFilters)

	// Recent QRIS Transactions
	var recentOrders []kiosk.PaymentOrder
	h.db.Order("id desc").Limit(8).Find(&recentOrders)

	c.HTML(http.StatusOK, "analytics", core.PageData{
		Title:    "Laporan Finansial & Analitik · Photobooth",
		UserName: middleware.CurrentUserName(c),
		Path:     "/admin/analytics",
		Data: gin.H{
			"TotalSessions":      totalSessions,
			"TotalPrints":        totalPrints,
			"TotalDownloads":     totalDownloads,
			"QRISRevenue":        qrisRevenue,
			"QRISPaidOrders":     qrisPaidOrders,
			"UsedVouchers":       usedVouchers,
			"VoucherRevenue":     voucherRevenue,
			"TotalGrossRevenue":  totalGrossRevenue,
			"PaperUsed":          paperUsed,
			"EstimatedCost":      estimatedCost,
			"EstimatedNetProfit": estimatedNetProfit,
			"Config":             cfgMap,
			"HourlyBars":         hourlyBars,
			"TopFrames":          topFrames,
			"TopFilters":         topFilters,
			"RecentOrders":       recentOrders,
		},
	})
}

// NotFound renders the 404 error page.
func (h *Handler) NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404", core.PageData{
		Title: "404 · Meridian",
	})
}
