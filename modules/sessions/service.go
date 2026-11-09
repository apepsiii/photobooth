package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"gorm.io/gorm"

	"mastergin/modules/settings"
)

// GenerateSessionCode generates a short, unique code like "PB-7X9K".
func GenerateSessionCode() string {
	const chars = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	result := make([]byte, 5)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[num.Int64()]
	}
	return "PB-" + string(result)
}

// SaveBase64File saves a data URI base64 string (image, video, etc.) to a file on disk.
func SaveBase64File(dataURI string, destPath string) error {
	idx := strings.Index(dataURI, ";base64,")
	rawB64 := dataURI
	if idx != -1 {
		rawB64 = dataURI[idx+8:]
	}

	data, err := base64.StdEncoding.DecodeString(rawB64)
	if err != nil {
		return err
	}

	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(destPath, data, 0644)
}

// SaveBase64Image saves a data URI base64 string to a file on disk.
func SaveBase64Image(dataURI string, destPath string) error {
	return SaveBase64File(dataURI, destPath)
}

// PrintPhoto triggers printing of the composite photo file on the target or default printer.
func PrintPhoto(db *gorm.DB, imagePath string, copies int) error {
	if copies <= 0 {
		copies = 1
	}

	absPath, err := filepath.Abs(strings.TrimPrefix(imagePath, "/"))
	if err != nil {
		return err
	}

	printerName := settings.Get(db, "printer_name")
	log.Printf("[PRINT] Spooling photo %s (copies: %d) to printer '%s'", absPath, copies, printerName)

	// Increment paper counter
	settings.IncrementPaperUsed(db, copies)

	if runtime.GOOS == "windows" {
		// On Windows, use PowerShell Start-Process with Print verb or mspaint /pt
		for i := 0; i < copies; i++ {
			var cmd *exec.Cmd
			if strings.TrimSpace(printerName) != "" {
				// Print to specific printer
				cmd = exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command",
					fmt.Sprintf(`Start-Process -FilePath "%s" -Verb PrintTo -ArgumentList '"%s"' -PassThru | Out-Null`, absPath, printerName))
			} else {
				// Print to default printer
				cmd = exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command",
					fmt.Sprintf(`Start-Process -FilePath "%s" -Verb Print -PassThru | Out-Null`, absPath))
			}
			if err := cmd.Start(); err != nil {
				log.Printf("[PRINT ERROR] Failed to trigger print: %v", err)
				return err
			}
		}
	} else {
		// Unix / macOS fallback (lpr command)
		cmd := exec.Command("lpr", absPath)
		_ = cmd.Run()
	}

	return nil
}

// GetByCode retrieves a session with frame and photos by its unique code.
func GetByCode(db *gorm.DB, code string) (*Session, error) {
	var s Session
	err := db.Preload("Frame").Preload("Photos").Where("session_code = ?", code).First(&s).Error
	return &s, err
}

// DeleteSession deletes a session record, associated raw photos, and files on disk.
func DeleteSession(db *gorm.DB, id uint) error {
	var s Session
	if err := db.Preload("Photos").First(&s, id).Error; err != nil {
		return err
	}

	// Delete composite file
	if s.CompositePath != "" {
		_ = os.Remove(filepath.FromSlash(strings.TrimPrefix(s.CompositePath, "/")))
	}

	// Delete live photo video file if exists
	if s.LivePhotoPath != "" {
		_ = os.Remove(filepath.FromSlash(strings.TrimPrefix(s.LivePhotoPath, "/")))
	}

	// Delete raw photos
	for _, p := range s.Photos {
		if p.PhotoPath != "" {
			_ = os.Remove(filepath.FromSlash(strings.TrimPrefix(p.PhotoPath, "/")))
		}
	}

	// Delete DB records safely within a transaction
	return db.Transaction(func(tx *gorm.DB) error {
		// Clean up any voucher session_code link
		if s.SessionCode != "" {
			_ = tx.Exec("UPDATE vouchers SET session_code = '' WHERE session_code = ?", s.SessionCode).Error
		}
		if err := tx.Where("session_id = ?", s.ID).Delete(&SessionPhoto{}).Error; err != nil {
			return err
		}
		return tx.Delete(&s).Error
	})
}

// ListRecentSessions returns sessions ordered by newest first.
func ListRecentSessions(db *gorm.DB, limit int) ([]Session, error) {
	var list []Session
	err := db.Preload("Frame").Order("id desc").Limit(limit).Find(&list).Error
	return list, err
}

func scaleAndQuantize(src image.Image, targetW, targetH int) *image.Paletted {
	rect := image.Rect(0, 0, targetW, targetH)
	pal := image.NewPaletted(rect, palette.Plan9)
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	for y := 0; y < targetH; y++ {
		srcY := srcBounds.Min.Y + (y * srcH) / targetH
		for x := 0; x < targetW; x++ {
			srcX := srcBounds.Min.X + (x * srcW) / targetW
			pal.Set(x, y, src.At(srcX, srcY))
		}
	}
	return pal
}

// GenerateBoomerangGIF creates an animated looping GIF from a session's raw photos.
func GenerateBoomerangGIF(photos []SessionPhoto, w io.Writer) error {
	if len(photos) == 0 {
		return fmt.Errorf("tidak ada foto untuk animasi GIF")
	}

	var decodedImages []image.Image
	for _, p := range photos {
		filePath := strings.TrimPrefix(p.PhotoPath, "/")
		file, err := os.Open(filePath)
		if err != nil {
			continue
		}
		img, _, err := image.Decode(file)
		file.Close()
		if err != nil {
			continue
		}
		decodedImages = append(decodedImages, img)
	}

	if len(decodedImages) == 0 {
		return fmt.Errorf("gagal membaca file foto mentah")
	}

	targetW := 480
	firstImg := decodedImages[0]
	origBounds := firstImg.Bounds()
	targetH := 360
	if origBounds.Dx() > 0 {
		targetH = (targetW * origBounds.Dy()) / origBounds.Dx()
	}

	var palettedImages []*image.Paletted
	for _, img := range decodedImages {
		palettedImages = append(palettedImages, scaleAndQuantize(img, targetW, targetH))
	}

	var sequence []*image.Paletted
	var delays []int

	// Forward sequence
	for _, p := range palettedImages {
		sequence = append(sequence, p)
		delays = append(delays, 32) // ~320ms per frame
	}
	// Reverse bounce sequence (Boomerang effect)
	if len(palettedImages) > 2 {
		for i := len(palettedImages) - 2; i > 0; i-- {
			sequence = append(sequence, palettedImages[i])
			delays = append(delays, 32)
		}
	}

	outGIF := &gif.GIF{
		Image:     sequence,
		Delay:     delays,
		LoopCount: 0,
	}

	return gif.EncodeAll(w, outGIF)
}

