# AGENTS.md - AI Developer Guide for mastergin

## Quick Reference

**Stack:** Go + Gin + GORM + SQLite (pure-Go driver) + Meridian HTML template
**Run:** `go run .` (ensure `.env` exists — copy from `.env.example`)
**Default login:** admin@meridian.com / admin123
**Port:** 8090 (configurable via `APP_PORT` in `.env`)

## Architecture

```
main.go              # Entry point: init config → logger → DB → seed → router → modules → serve
core/                # Framework layer (do NOT add business logic here)
  config.go          # Loads .env into Config struct
  database.go        # GORM + SQLite singleton, AutoMigrate, PRAGMA tuning
  logger.go          # Log output configuration
  router.go          # Gin engine setup, middleware, static assets, health check
  template.go        # Custom multi-template renderer (per-page template sets)
middleware/
  auth.go            # RequireAuth (JWT cookie) + CurrentUserName helper
models/
  user.go            # User GORM model (ID, Name, Email, PasswordHash, Role)
modules/             # Feature modules (AI works here)
  auth/              # Built-in auth module
    service.go       # Register, Login (bcrypt + JWT), SeedAdmin
    handler.go       # HTTP handlers (LoginPage, Login, RegisterPage, Register, Logout, Dashboard, NotFound)
    routes.go        # RegisterRoutes(r, db, cfg) — public + private route groups
template/meridian/
  layouts/master.html   # Dashboard master layout (sidebar + navbar + {{block "content"}})
  pages/*.html          # Page templates (each defines its name + optional "content" block)
  assets/               # CSS, JS (served at /assets/)
```

## Template System

The boilerplate uses a **custom multi-template renderer** (`core/template.go`) that
gives each page its own isolated `*template.Template` set. This enables true Go
template inheritance (`{{block "content"}}`) without block-name conflicts across pages.

### Adding a new dashboard page

1. Create `template/meridian/pages/myfeature.html`:
```html
{{define "myfeature"}}
{{template "master" .}}
{{end}}

{{define "content"}}
<!-- your page content here -->
{{end}}
```

2. Render it from a handler: `c.HTML(http.StatusOK, "myfeature", core.PageData{Title: "My Feature", UserName: middleware.CurrentUserName(c)})`

### Standalone pages (no sidebar)

Pages that don't use the master layout (like login/register/404) simply omit the
`{{template "master" .}}` call and are full HTML documents wrapped in `{{define "pagename"}}...{{end}}`.

## Adding a New Module

Follow the pattern in `modules/auth/`:

1. **`modules/yourmodule/model.go`** — GORM struct
2. **`modules/yourmodule/service.go`** — business logic, accepts `*gorm.DB`
3. **`modules/yourmodule/handler.go`** — Gin handlers
4. **`modules/yourmodule/routes.go`** — `func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *core.Config)`
5. Add model to `AutoMigrate` call in `main.go`
6. Call `yourmodule.RegisterRoutes(r, db, cfg)` in `main.go`

## Key Decisions

- **Pure-Go SQLite** (`github.com/glebarez/sqlite`) — no CGO/GCC required. The DEVGUIDE mentions `gorm.io/driver/sqlite` but that needs CGO; this project uses the drop-in pure-Go alternative.
- **JWT in HTTP-only cookie** — server-side rendered HTML app, so cookies are more secure than localStorage tokens.
- **Static assets** served at `/assets/` from `template/meridian/assets/`.
