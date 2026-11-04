package settings

import (
	"net"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

var defaultSettings = map[string]struct {
	Value string
	Label string
}{
	"app_name":          {Value: "School Photobooth", Label: "Nama Aplikasi Kiosk"},
	"event_name":        {Value: "Pentas Seni & Bazaar 2026", Label: "Nama Acara"},
	"school_name":       {Value: "SMA Negeri 1", Label: "Nama Sekolah / Institusi"},
	"countdown_seconds": {Value: "3", Label: "Durasi Countdown Kamera (detik)"},
	"print_enabled":     {Value: "true", Label: "Aktifkan Fitur Cetak"},
	"auto_print":        {Value: "false", Label: "Cetak Otomatis Tanpa Konfirmasi"},
	"copies_count":      {Value: "1", Label: "Jumlah Lembar Cetak Default"},
	"printer_name":      {Value: "", Label: "Nama Printer Target di Windows"},
	"kiosk_pin":             {Value: "1234", Label: "PIN Kiosk / Operator"},
	"local_server_ip":       {Value: "", Label: "IP Server Lokal (Kosongkan untuk Deteksi Otomatis)"},
	"voucher_mode_enabled":  {Value: "false", Label: "Aktifkan Sistem Tiket / Voucher Berbayar"},
	"payment_mode_enabled":  {Value: "false", Label: "Aktifkan Mode Pembayaran Mandiri (QRIS)"},
	"price_single":          {Value: "20000", Label: "Harga Paket 1 Lembar Cetak (Rp)"},
	"price_double":          {Value: "30000", Label: "Harga Paket 2 Lembar Cetak (Rp)"},
	"qris_static_code":      {Value: "", Label: "Kode / URL String QRIS Merchant"},
	"paper_capacity":        {Value: "100", Label: "Kapasitas Kertas / Ribbon Printer (Lembar)"},
	"paper_used":            {Value: "0", Label: "Jumlah Kertas Terpakai (Lembar)"},
}

// SeedSettings ensures all required settings exist in the database.
func SeedSettings(db *gorm.DB) error {
	for k, v := range defaultSettings {
		var existing Setting
		if err := db.Where("key = ?", k).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				setting := Setting{
					Key:   k,
					Value: v.Value,
					Label: v.Label,
				}
				if err := db.Create(&setting).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

// Get fetches a single setting by key.
func Get(db *gorm.DB, key string) string {
	var s Setting
	if err := db.Where("key = ?", key).First(&s).Error; err != nil {
		if d, ok := defaultSettings[key]; ok {
			return d.Value
		}
		return ""
	}
	return s.Value
}

// Set saves or updates a setting by key.
func Set(db *gorm.DB, key, value string) error {
	var s Setting
	if err := db.Where("key = ?", key).First(&s).Error; err != nil {
		label := key
		if d, ok := defaultSettings[key]; ok {
			label = d.Label
		}
		return db.Create(&Setting{Key: key, Value: value, Label: label}).Error
	}
	s.Value = value
	return db.Save(&s).Error
}

// DetectLocalIP returns the first non-loopback IPv4 address found on the system.
func DetectLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				// Prefer private network ranges 192.168.*, 10.*, 172.16-31.*
				if strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "172.") {
					return ip
				}
			}
		}
	}
	return "127.0.0.1"
}

// GetConfig returns a typed ConfigMap with all settings resolved.
func GetConfig(db *gorm.DB, port string) ConfigMap {
	countdown, _ := strconv.Atoi(Get(db, "countdown_seconds"))
	if countdown <= 0 {
		countdown = 3
	}

	copies, _ := strconv.Atoi(Get(db, "copies_count"))
	if copies <= 0 {
		copies = 1
	}

	localIP := Get(db, "local_server_ip")
	if strings.TrimSpace(localIP) == "" {
		localIP = DetectLocalIP()
	}

	priceSingle, _ := strconv.Atoi(Get(db, "price_single"))
	if priceSingle <= 0 {
		priceSingle = 20000
	}
	priceDouble, _ := strconv.Atoi(Get(db, "price_double"))
	if priceDouble <= 0 {
		priceDouble = 30000
	}

	paperCapacity, _ := strconv.Atoi(Get(db, "paper_capacity"))
	if paperCapacity <= 0 {
		paperCapacity = 100
	}
	paperUsed, _ := strconv.Atoi(Get(db, "paper_used"))
	paperRemaining := paperCapacity - paperUsed
	if paperRemaining < 0 {
		paperRemaining = 0
	}

	return ConfigMap{
		AppName:            Get(db, "app_name"),
		EventName:          Get(db, "event_name"),
		SchoolName:         Get(db, "school_name"),
		CountdownSeconds:   countdown,
		PrintEnabled:       Get(db, "print_enabled") == "true",
		AutoPrint:          Get(db, "auto_print") == "true",
		CopiesCount:        copies,
		PrinterName:        Get(db, "printer_name"),
		KioskPin:           Get(db, "kiosk_pin"),
		LocalServerIP:      localIP,
		Port:               port,
		VoucherModeEnabled: Get(db, "voucher_mode_enabled") == "true",
		PaymentModeEnabled: Get(db, "payment_mode_enabled") == "true",
		PriceSingle:        priceSingle,
		PriceDouble:        priceDouble,
		QRISStaticCode:     Get(db, "qris_static_code"),
		PaperCapacity:      paperCapacity,
		PaperUsed:          paperUsed,
		PaperRemaining:     paperRemaining,
	}
}

// IncrementPaperUsed adds count to paper_used counter.
func IncrementPaperUsed(db *gorm.DB, count int) int {
	if count <= 0 {
		count = 1
	}
	used, _ := strconv.Atoi(Get(db, "paper_used"))
	used += count
	_ = Set(db, "paper_used", strconv.Itoa(used))
	return used
}

// ResetPaperCounter resets paper_used to 0, optionally updating paper_capacity.
func ResetPaperCounter(db *gorm.DB, capacity int) error {
	if capacity > 0 {
		_ = Set(db, "paper_capacity", strconv.Itoa(capacity))
	}
	return Set(db, "paper_used", "0")
}


