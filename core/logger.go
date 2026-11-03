package core

import (
	"io"
	"log"
	"os"
)

// InitLogger configures the standard library logger output.
// In debug mode logs go to stdout; in release mode they go to a file.
func InitLogger(cfg *Config) {
	var writer io.Writer = os.Stdout

	if !cfg.IsDebug() {
		f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Could not open log file, falling back to stdout:", err)
		} else {
			writer = f
		}
	}

	log.SetOutput(writer)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
