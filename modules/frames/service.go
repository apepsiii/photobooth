package frames

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

// ListActive returns all active frames.
func ListActive(db *gorm.DB) ([]Frame, error) {
	var list []Frame
	err := db.Where("is_active = ?", true).Order("id asc").Find(&list).Error
	return list, err
}

// ListAll returns all frames.
func ListAll(db *gorm.DB) ([]Frame, error) {
	var list []Frame
	err := db.Order("id desc").Find(&list).Error
	return list, err
}

// GetByID finds a frame by ID.
func GetByID(db *gorm.DB, id uint) (*Frame, error) {
	var f Frame
	err := db.First(&f, id).Error
	return &f, err
}

// Create creates a new frame.
func Create(db *gorm.DB, f *Frame) error {
	return db.Create(f).Error
}

// Update updates an existing frame.
func Update(db *gorm.DB, f *Frame) error {
	return db.Save(f).Error
}

// Delete deletes a frame and safely unlinks associated sessions so foreign key constraint does not block deletion.
func Delete(db *gorm.DB, id uint) error {
	var f Frame
	if err := db.First(&f, id).Error; err != nil {
		return err
	}

	// Remove overlay image file from disk if exists
	if f.OverlayImage != "" {
		filePath := filepath.FromSlash(strings.TrimPrefix(f.OverlayImage, "/"))
		_ = os.Remove(filePath)
	}

	// Safely unlink foreign key references in sessions table first, then delete frame
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE sessions SET frame_id = NULL WHERE frame_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&f).Error
	})
}

// SeedDefaultFrames ensures standard default frame templates exist on disk and in database.
func SeedDefaultFrames(db *gorm.DB) error {
	framesDir := filepath.Join("storage", "frames")
	_ = os.MkdirAll(framesDir, 0755)

	// 1. Classic 2-Strip (Funcam Style)
	stripSlots := []PhotoSlot{
		{Index: 0, X: 40, Y: 70, Width: 520, Height: 420},
		{Index: 1, X: 40, Y: 520, Width: 520, Height: 420},
		{Index: 2, X: 40, Y: 970, Width: 520, Height: 420},
	}
	stripSlotsJSON, _ := json.Marshal(stripSlots)
	stripFile := filepath.Join(framesDir, "strip_classic.png")
	if _, err := os.Stat(stripFile); os.IsNotExist(err) {
		_ = generateStripPNG(stripFile, 600, 1800, stripSlots, "PHOTOBOOTH SEKOLAH")
	}

	var count int64
	db.Model(&Frame{}).Where("slug = ?", "strip-classic").Count(&count)
	if count == 0 {
		db.Create(&Frame{
			Name:         "Classic 2-Strip (Funcam Style)",
			Slug:         "strip-classic",
			Category:     "Classic",
			Type:         "strip_2x6",
			Description:  "Format strip klasik 3 foto vertikal (2x6 inci), sangat populer untuk bookmark & photobox.",
			TotalShots:   3,
			Width:        600,
			Height:       1800,
			OverlayImage: "/storage/frames/strip_classic.png",
			SlotsJSON:    string(stripSlotsJSON),
			IsActive:     true,
		})
	}

	// 2. Korean 4-Cut Grid
	gridSlots := []PhotoSlot{
		{Index: 0, X: 50, Y: 70, Width: 530, Height: 680},
		{Index: 1, X: 620, Y: 70, Width: 530, Height: 680},
		{Index: 2, X: 50, Y: 780, Width: 530, Height: 680},
		{Index: 3, X: 620, Y: 780, Width: 530, Height: 680},
	}
	gridSlotsJSON, _ := json.Marshal(gridSlots)
	gridFile := filepath.Join(framesDir, "grid_4cut.png")
	if _, err := os.Stat(gridFile); os.IsNotExist(err) {
		_ = generateGridPNG(gridFile, 1200, 1800, gridSlots, "PHOTOBOOTH SEKOLAH - CLASS OF 2026")
	}

	db.Model(&Frame{}).Where("slug = ?", "korean-4cut").Count(&count)
	if count == 0 {
		db.Create(&Frame{
			Name:         "Korean 4-Cut Grid (4R)",
			Slug:         "korean-4cut",
			Category:     "Pensi / Event",
			Type:         "grid_4cut",
			Description:  "Format 4-Cut ala photobooth Korea (4x6 inci / 4R), 4 foto dalam grid 2x2 elegan.",
			TotalShots:   4,
			Width:        1200,
			Height:       1800,
			OverlayImage: "/storage/frames/grid_4cut.png",
			SlotsJSON:    string(gridSlotsJSON),
			IsActive:     true,
		})
	}

	// 3. Polaroid Minimalist
	polaroidSlots := []PhotoSlot{
		{Index: 0, X: 80, Y: 80, Width: 1040, Height: 1040},
	}
	polaroidSlotsJSON, _ := json.Marshal(polaroidSlots)
	polaroidFile := filepath.Join(framesDir, "polaroid_single.png")
	if _, err := os.Stat(polaroidFile); os.IsNotExist(err) {
		_ = generatePolaroidPNG(polaroidFile, 1200, 1500, polaroidSlots, "BEST MEMORIES - OUR HIGH SCHOOL STORY")
	}

	db.Model(&Frame{}).Where("slug = ?", "polaroid-single").Count(&count)
	if count == 0 {
		db.Create(&Frame{
			Name:         "Polaroid Minimalist (Wide)",
			Slug:         "polaroid-single",
			Category:     "Graduation",
			Type:         "polaroid",
			Description:  "Format polaroid klasik 1 foto square besar dengan border putih luas untuk tanda tangan.",
			TotalShots:   1,
			Width:        1200,
			Height:       1500,
			OverlayImage: "/storage/frames/polaroid_single.png",
			SlotsJSON:    string(polaroidSlotsJSON),
			IsActive:     true,
		})
	}

	// Fill any empty category for legacy rows
	db.Model(&Frame{}).Where("category IS NULL OR category = ''").Update("category", "Umum")

	return nil
}

// generateStripPNG generates a clean minimalist PNG frame with transparent cutout slots.
func generateStripPNG(filename string, w, h int, slots []PhotoSlot, footerText string) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Clean white frame background
	bgColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Transparent color for photo cutouts
	transparent := color.RGBA{R: 0, G: 0, B: 0, A: 0}

	// Punch out transparent windows where photos will appear
	for _, s := range slots {
		slotRect := image.Rect(s.X, s.Y, s.X+s.Width, s.Y+s.Height)
		draw.Draw(img, slotRect, &image.Uniform{transparent}, image.Point{}, draw.Src)
	}

	// Save as PNG
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func generateGridPNG(filename string, w, h int, slots []PhotoSlot, footerText string) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	bgColor := color.RGBA{R: 250, G: 250, B: 250, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	transparent := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	for _, s := range slots {
		slotRect := image.Rect(s.X, s.Y, s.X+s.Width, s.Y+s.Height)
		draw.Draw(img, slotRect, &image.Uniform{transparent}, image.Point{}, draw.Src)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func generatePolaroidPNG(filename string, w, h int, slots []PhotoSlot, footerText string) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	bgColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	transparent := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	for _, s := range slots {
		slotRect := image.Rect(s.X, s.Y, s.X+s.Width, s.Y+s.Height)
		draw.Draw(img, slotRect, &image.Uniform{transparent}, image.Point{}, draw.Src)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
