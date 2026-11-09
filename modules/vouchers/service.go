package vouchers

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GenerateRandomCode generates a 6-digit numeric voucher code like "829401".
func GenerateRandomCode() string {
	const firstDigits = "123456789"
	const digits = "0123456789"
	result := make([]byte, 6)

	first, _ := rand.Int(rand.Reader, big.NewInt(int64(len(firstDigits))))
	result[0] = firstDigits[first.Int64()]

	for i := 1; i < 6; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		result[i] = digits[num.Int64()]
	}
	return string(result)
}

// GenerateBatch creates multiple unique vouchers in bulk for the ticketing desk.
func GenerateBatch(db *gorm.DB, count int) ([]Voucher, error) {
	if count <= 0 {
		count = 10
	}
	if count > 200 {
		count = 200
	}

	var vouchers []Voucher
	for i := 0; i < count; i++ {
		code := GenerateRandomCode()
		for {
			var exists int64
			db.Model(&Voucher{}).Where("code = ?", code).Count(&exists)
			if exists == 0 {
				break
			}
			code = GenerateRandomCode()
		}

		v := Voucher{
			Code:      code,
			IsUsed:    false,
			CreatedAt: time.Now(),
		}
		if err := db.Create(&v).Error; err == nil {
			vouchers = append(vouchers, v)
		}
	}
	return vouchers, nil
}

// ValidateAndRedeem checks if a voucher is valid and marks it as used.
func ValidateAndRedeem(db *gorm.DB, code string, sessionCode string) error {
	cleanCode := strings.ToUpper(strings.TrimSpace(code))
	var v Voucher
	if err := db.Where("code = ?", cleanCode).First(&v).Error; err != nil {
		return errors.New("Kode voucher tidak ditemukan.")
	}

	if v.IsUsed {
		return errors.New("Voucher ini sudah pernah digunakan sebelumnya.")
	}

	now := time.Now()
	v.IsUsed = true
	v.UsedAt = &now
	v.SessionCode = sessionCode
	return db.Save(&v).Error
}

// CheckValidity checks if a voucher code is currently valid without redeeming it yet.
func CheckValidity(db *gorm.DB, code string) (*Voucher, error) {
	cleanCode := strings.ToUpper(strings.TrimSpace(code))
	var v Voucher
	if err := db.Where("code = ?", cleanCode).First(&v).Error; err != nil {
		return nil, errors.New("Kode voucher tidak valid.")
	}
	if v.IsUsed {
		return nil, errors.New("Kode voucher sudah digunakan.")
	}
	return &v, nil
}

// ListAll returns all vouchers ordered by newest first.
func ListAll(db *gorm.DB) ([]Voucher, error) {
	var list []Voucher
	err := db.Order("id desc").Find(&list).Error
	return list, err
}

// Delete removes a voucher.
func Delete(db *gorm.DB, id uint) error {
	return db.Delete(&Voucher{}, id).Error
}

