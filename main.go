package main

import (
	"log"

	"mastergin/core"
	"mastergin/models"
	"mastergin/modules/auth"
	"mastergin/modules/frames"
	"mastergin/modules/kiosk"
	"mastergin/modules/sessions"
	"mastergin/modules/settings"
	"mastergin/modules/vouchers"
)

func main() {
	// 1. Load configuration from .env / environment.
	cfg := core.LoadConfig()

	// 2. Initialise logger.
	core.InitLogger(cfg)

	// 3. Initialise database and auto-migrate core & module models.
	db := core.InitDB(cfg,
		&models.User{},
		&settings.Setting{},
		&frames.Frame{},
		&sessions.Session{},
		&sessions.SessionPhoto{},
		&vouchers.Voucher{},
		&kiosk.PaymentOrder{},
	)

	// 4. Seed default admin user.
	if err := auth.SeedAdmin(db); err != nil {
		log.Printf("Warning: failed to seed admin user: %v", err)
	} else {
		log.Println("Admin user ready (admin@meridian.com / admin123)")
	}

	// 5. Seed default settings & frames.
	if err := settings.SeedSettings(db); err != nil {
		log.Printf("Warning: failed to seed settings: %v", err)
	}
	if err := frames.SeedDefaultFrames(db); err != nil {
		log.Printf("Warning: failed to seed default frames: %v", err)
	}

	// 6. Set up the Gin router with templates and middleware.
	r := core.SetupRouter(cfg)

	// 7. Register module routes.
	kiosk.RegisterRoutes(r, db, cfg)
	frames.RegisterRoutes(r, db, cfg)
	sessions.RegisterRoutes(r, db, cfg)
	settings.RegisterRoutes(r, db, cfg)
	vouchers.RegisterRoutes(r, db, cfg)
	auth.RegisterRoutes(r, db, cfg)

	// 8. Start the server.
	log.Printf("Photobooth Server starting on :%s (mode: %s)", cfg.Port, cfg.Env)
	log.Printf("Kiosk interface: http://localhost:%s/kiosk", cfg.Port)
	log.Printf("Admin panel: http://localhost:%s/admin/sessions", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
