package frames

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mastergin/core"
	"mastergin/middleware"
)

// Handler handles frame admin requests and API endpoints.
type Handler struct {
	db  *gorm.DB
	cfg *core.Config
}

// NewHandler creates a new frame handler.
func NewHandler(db *gorm.DB, cfg *core.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

// AdminFramesPage renders the frames management page.
func (h *Handler) AdminFramesPage(c *gin.Context) {
	frames, err := ListAll(h.db)
	if err != nil {
		frames = []Frame{}
	}

	// Collect unique categories for filter tabs
	categoryMap := make(map[string]bool)
	for _, f := range frames {
		if f.Category != "" {
			categoryMap[f.Category] = true
		}
	}
	var categories []string
	for cat := range categoryMap {
		categories = append(categories, cat)
	}

	c.HTML(http.StatusOK, "frames", core.PageData{
		Title:    "Manajemen Frame Photobooth",
		UserName: middleware.CurrentUserName(c),
		Path:     "/admin/frames",
		Data: gin.H{
			"Frames":     frames,
			"Categories": categories,
		},
	})
}

// CreateFrame handles uploading and adding a custom frame.
func (h *Handler) CreateFrame(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	category := strings.TrimSpace(c.PostForm("category"))
	if category == "" {
		category = "Umum"
	}
	frameType := c.PostForm("type")
	if frameType == "" {
		frameType = "custom"
	}
	desc := c.PostForm("description")
	totalShots, _ := strconv.Atoi(c.PostForm("total_shots"))
	width, _ := strconv.Atoi(c.PostForm("width"))
	height, _ := strconv.Atoi(c.PostForm("height"))
	slotsJSON := c.PostForm("slots_json")
	layerOrder := c.PostForm("layer_order")
	if layerOrder == "" {
		layerOrder = "photo_on_top"
	}

	// Munculkan atau tidak (Active status)
	isActiveVal := c.PostForm("is_active")
	isActive := isActiveVal == "1" || isActiveVal == "true" || isActiveVal == "on"

	if width <= 0 {
		width = 600
	}
	if height <= 0 {
		height = 1800
	}
	if totalShots <= 0 {
		totalShots = 3
	}

	overlayPath := ""
	file, err := c.FormFile("overlay_file")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("frame_%d%s", time.Now().UnixNano(), ext)
		dst := filepath.Join("storage", "frames", filename)
		if err := c.SaveUploadedFile(file, dst); err == nil {
			overlayPath = "/storage/frames/" + filename
		}
	}

	slug := fmt.Sprintf("frame-%d", time.Now().Unix())
	newFrame := Frame{
		Name:         name,
		Slug:         slug,
		Category:     category,
		Type:         frameType,
		Description:  desc,
		TotalShots:   totalShots,
		Width:        width,
		Height:       height,
		OverlayImage: overlayPath,
		SlotsJSON:    slotsJSON,
		LayerOrder:   layerOrder,
		IsActive:     isActive,
	}

	_ = Create(h.db, &newFrame)
	c.Redirect(http.StatusFound, "/admin/frames")
}

// UpdateFrame handles editing an existing frame.
func (h *Handler) UpdateFrame(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	f, err := GetByID(h.db, uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/frames")
		return
	}

	name := strings.TrimSpace(c.PostForm("name"))
	if name != "" {
		f.Name = name
	}
	category := strings.TrimSpace(c.PostForm("category"))
	if category != "" {
		f.Category = category
	}
	if t := c.PostForm("type"); t != "" {
		f.Type = t
	}
	f.Description = c.PostForm("description")

	if shots, err := strconv.Atoi(c.PostForm("total_shots")); err == nil && shots > 0 {
		f.TotalShots = shots
	}
	if w, err := strconv.Atoi(c.PostForm("width")); err == nil && w > 0 {
		f.Width = w
	}
	if hgt, err := strconv.Atoi(c.PostForm("height")); err == nil && hgt > 0 {
		f.Height = hgt
	}
	if sJSON := c.PostForm("slots_json"); sJSON != "" {
		f.SlotsJSON = sJSON
	}
	if lo := c.PostForm("layer_order"); lo != "" {
		f.LayerOrder = lo
	}

	// Munculkan atau tidak
	isActiveVal := c.PostForm("is_active")
	f.IsActive = isActiveVal == "1" || isActiveVal == "true" || isActiveVal == "on"

	// Check if new overlay image file uploaded
	file, err := c.FormFile("overlay_file")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("frame_%d%s", time.Now().UnixNano(), ext)
		dst := filepath.Join("storage", "frames", filename)
		if err := c.SaveUploadedFile(file, dst); err == nil {
			f.OverlayImage = "/storage/frames/" + filename
		}
	}

	_ = Update(h.db, f)
	c.Redirect(http.StatusFound, "/admin/frames")
}

// ToggleFrame toggles active state.
func (h *Handler) ToggleFrame(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	f, err := GetByID(h.db, uint(id))
	if err == nil {
		f.IsActive = !f.IsActive
		_ = Update(h.db, f)
	}
	c.Redirect(http.StatusFound, "/admin/frames")
}

// DeleteFrame deletes a frame.
func (h *Handler) DeleteFrame(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err == nil {
		if err := Delete(h.db, uint(id)); err != nil {
			log.Printf("[ERROR] DeleteFrame failed for ID %d: %v", id, err)
		}
	}
	c.Redirect(http.StatusFound, "/admin/frames")
}

// APIGetFrame returns single frame JSON.
func (h *Handler) APIGetFrame(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	f, err := GetByID(h.db, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Frame not found"})
		return
	}
	c.JSON(http.StatusOK, f)
}

// APIListActive returns JSON list of active frames for Kiosk.
func (h *Handler) APIListActive(c *gin.Context) {
	frames, err := ListActive(h.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, frames)
}

