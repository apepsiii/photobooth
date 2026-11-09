package sessions

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mastergin/core"
	"mastergin/middleware"
	"mastergin/modules/settings"
)

// Handler handles sessions HTTP requests.
type Handler struct {
	db  *gorm.DB
	cfg *core.Config
}

// NewHandler creates a new Handler.
func NewHandler(db *gorm.DB, cfg *core.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

// AdminSessionsPage renders the session gallery in Meridian admin layout.
func (h *Handler) AdminSessionsPage(c *gin.Context) {
	list, err := ListRecentSessions(h.db, 50)
	if err != nil {
		list = []Session{}
	}

	var totalSessions int64
	var totalPrints int64
	var totalDownloads int64

	h.db.Model(&Session{}).Count(&totalSessions)
	h.db.Model(&Session{}).Select("COALESCE(SUM(printed_count), 0)").Scan(&totalPrints)
	h.db.Model(&Session{}).Select("COALESCE(SUM(downloaded_count), 0)").Scan(&totalDownloads)

	c.HTML(http.StatusOK, "sessions", core.PageData{
		Title:    "Galeri Riwayat Pemotretan",
		UserName: middleware.CurrentUserName(c),
		Path:     "/admin/sessions",
		Data: gin.H{
			"Sessions":       list,
			"TotalSessions":  totalSessions,
			"TotalPrints":    totalPrints,
			"TotalDownloads": totalDownloads,
		},
	})
}

// ReprintSession triggers a reprint of a session.
func (h *Handler) ReprintSession(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var s Session
	if err := h.db.First(&s, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	copies := 1
	if c.PostForm("copies") != "" {
		copies, _ = strconv.Atoi(c.PostForm("copies"))
	}

	_ = PrintPhoto(h.db, s.CompositePath, copies)
	s.PrintedCount += copies
	h.db.Save(&s)

	c.Redirect(http.StatusFound, "/admin/sessions")
}

// DeleteSession handles deleting a session from admin.
func (h *Handler) DeleteSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err == nil {
		if err := DeleteSession(h.db, uint(id)); err != nil {
			log.Printf("[ERROR] DeleteSession failed for ID %d: %v", id, err)
		}
	}
	c.Redirect(http.StatusFound, "/admin/sessions")
}

// TestPrint triggers a printer connectivity test print.
func (h *Handler) TestPrint(c *gin.Context) {
	testFile := "storage/frames/strip_classic.png"
	_ = PrintPhoto(h.db, testFile, 1)
	c.Redirect(http.StatusFound, "/admin/settings?test_print=ok")
}

// DownloadSoftfilePage renders the mobile view for students scanning the QR code.
func (h *Handler) DownloadSoftfilePage(c *gin.Context) {
	code := c.Param("code")
	s, err := GetByCode(h.db, code)
	if err != nil {
		c.HTML(http.StatusNotFound, "404", core.PageData{
			Title: "Sesi Foto Tidak Ditemukan",
		})
		return
	}

	cfgMap := settings.GetConfig(h.db, h.cfg.Port)

	c.HTML(http.StatusOK, "download", gin.H{
		"Session": s,
		"Config":  cfgMap,
	})
}

// DownloadPhoto streams the final composite photo or individual shot to the student's phone.
func (h *Handler) DownloadPhoto(c *gin.Context) {
	code := c.Param("code")
	s, err := GetByCode(h.db, code)
	if err != nil {
		c.String(http.StatusNotFound, "Photo not found")
		return
	}

	// Increment download counter
	s.DownloadedCount++
	h.db.Save(s)

	// Local file path
	cleanPath := strings.TrimPrefix(s.CompositePath, "/")
	filename := fmt.Sprintf("%s_photobooth.jpg", code)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.File(cleanPath)
}

// DownloadGIF generates and streams an animated Boomerang GIF of the session.
func (h *Handler) DownloadGIF(c *gin.Context) {
	code := c.Param("code")
	s, err := GetByCode(h.db, code)
	if err != nil {
		c.String(http.StatusNotFound, "Photo not found")
		return
	}

	s.DownloadedCount++
	h.db.Save(s)

	filename := fmt.Sprintf("%s_boomerang.gif", code)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "image/gif")

	if err := GenerateBoomerangGIF(s.Photos, c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Gagal menghasilkan GIF: "+err.Error())
	}
}

// ExportAllZip packages all final composite photos into a single ZIP file for event documentation.
func (h *Handler) ExportAllZip(c *gin.Context) {
	var list []Session
	_ = h.db.Find(&list).Error

	zipFilename := fmt.Sprintf("Photobooth_Dokumentasi_%s.zip", time.Now().Format("20060102_150405"))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, zipFilename))
	c.Header("Content-Type", "application/zip")

	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	for _, sess := range list {
		filePath := strings.TrimPrefix(sess.CompositePath, "/")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}

		fileToZip, err := os.Open(filePath)
		if err != nil {
			continue
		}

		entryName := fmt.Sprintf("%s_%s.jpg", sess.SessionCode, sess.CreatedAt.Format("150405"))
		zipEntry, err := zipWriter.Create(entryName)
		if err == nil {
			_, _ = io.Copy(zipEntry, fileToZip)
		}
		fileToZip.Close()
	}
}

// ExportCSV streams a CSV file of all photobooth sessions for reporting and financial audits.
func (h *Handler) ExportCSV(c *gin.Context) {
	var list []Session
	_ = h.db.Preload("Frame").Order("id asc").Find(&list).Error

	csvFilename := fmt.Sprintf("Photobooth_Laporan_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, csvFilename))
	c.Header("Content-Type", "text/csv; charset=utf-8")

	// UTF-8 BOM for Microsoft Excel compatibility
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	_ = writer.Write([]string{"No", "Kode Sesi", "Template Frame", "Filter", "Waktu Sesi", "Jumlah Cetak", "Jumlah Unduhan Softfile"})

	for i, s := range list {
		frameName := "Default"
		if s.Frame.Name != "" {
			frameName = s.Frame.Name
		}
		_ = writer.Write([]string{
			strconv.Itoa(i + 1),
			s.SessionCode,
			frameName,
			s.FilterApplied,
			s.CreatedAt.Format("2006-01-02 15:04:05"),
			strconv.Itoa(s.PrintedCount),
			strconv.Itoa(s.DownloadedCount),
		})
	}
}

