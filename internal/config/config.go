package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port           string `yaml:"port"`
	BaseURL        string `yaml:"base_url"`
	DataDir        string `yaml:"data_dir"`
	MongoURI       string `yaml:"mongo_uri"`
	MongoDB        string `yaml:"mongo_db"`
	LibreOfficeBin string `yaml:"libreoffice_bin"`
	MaxUploadSize  int64  `yaml:"max_upload_size"`
	ConversionDPI  int    `yaml:"conversion_dpi"`
	ThumbnailDPI   int    `yaml:"thumbnail_dpi"`
	SessionSecret  string `yaml:"session_secret"`
	APIKey         string `yaml:"api_key"`

	// APIKeyFile is the file holding the API key when none was configured
	// (empty if the key came from config/env, or could not be persisted).
	APIKeyFile string `yaml:"-"`
}

func Load() *Config {
	cfg := &Config{
		Port:          "8080",
		BaseURL:       "http://localhost:8080",
		DataDir:       "./data",
		MongoURI:      "",
		MongoDB:       "flipbook",
		MaxUploadSize: 104857600, // 100MB
		ConversionDPI: 300,
		ThumbnailDPI:  72,
	}

	// Try loading config file: FLIPBOOK_CONFIG env, then config.dev.yaml, then config.yaml
	configPaths := []string{
		os.Getenv("FLIPBOOK_CONFIG"),
		"config.dev.yaml",
		"config.yaml",
	}
	for _, path := range configPaths {
		if path == "" {
			continue
		}
		if data, err := os.ReadFile(path); err == nil {
			yaml.Unmarshal(data, cfg)
			break
		}
	}

	// Environment variables override config file values
	if v := os.Getenv("FLIPBOOK_PORT"); v != "" {
		cfg.Port = v
	}
	if v := os.Getenv("FLIPBOOK_BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if v := os.Getenv("FLIPBOOK_DATA_DIR"); v != "" {
		cfg.DataDir = v
	}
	if v := os.Getenv("FLIPBOOK_MONGO_URI"); v != "" {
		cfg.MongoURI = v
	}
	if v := os.Getenv("FLIPBOOK_MONGO_DB"); v != "" {
		cfg.MongoDB = v
	}
	if v := os.Getenv("FLIPBOOK_LIBREOFFICE_BIN"); v != "" {
		cfg.LibreOfficeBin = v
	}
	if v := os.Getenv("FLIPBOOK_MAX_UPLOAD_SIZE"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			cfg.MaxUploadSize = n
		}
	}
	if v := os.Getenv("FLIPBOOK_CONVERSION_DPI"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.ConversionDPI = n
		}
	}
	if v := os.Getenv("FLIPBOOK_THUMBNAIL_DPI"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.ThumbnailDPI = n
		}
	}
	if v := os.Getenv("FLIPBOOK_SESSION_SECRET"); v != "" {
		cfg.SessionSecret = v
	}
	if v := os.Getenv("FLIPBOOK_API_KEY"); v != "" {
		cfg.APIKey = v
	}

	if cfg.LibreOfficeBin == "" {
		cfg.LibreOfficeBin = findLibreOffice()
	}

	// Generate a random session secret if not set
	if cfg.SessionSecret == "" {
		b := make([]byte, 32)
		rand.Read(b)
		cfg.SessionSecret = hex.EncodeToString(b)
	}

	// No configured API key: reuse the persisted one, or generate and persist it.
	// The server and the `mcp` subprocess both call Load, so a shared file keeps
	// them on the same key, and the key survives restarts without being logged.
	if cfg.APIKey == "" {
		cfg.APIKeyFile = filepath.Join(cfg.DataDir, apiKeyFileName)
		key, err := loadOrCreateAPIKey(cfg.APIKeyFile)
		if err != nil {
			log.Printf("WARNING: could not persist API key to %s: %v", cfg.APIKeyFile, err)
			cfg.APIKeyFile = ""
		}
		cfg.APIKey = key
	}

	return cfg
}

const apiKeyFileName = "api_key"

// loadOrCreateAPIKey returns the key stored at path, creating it (mode 0600)
// with a random value if it does not exist. On a write failure it still
// returns a usable in-memory key along with the error.
func loadOrCreateAPIKey(path string) (string, error) {
	if data, err := os.ReadFile(path); err == nil {
		if key := strings.TrimSpace(string(data)); key != "" {
			return key, nil
		}
	}

	b := make([]byte, 32)
	rand.Read(b)
	key := hex.EncodeToString(b)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return key, err
	}
	if err := os.WriteFile(path, []byte(key+"\n"), 0o600); err != nil {
		return key, err
	}
	return key, nil
}

func findLibreOffice() string {
	paths := []string{
		"/Applications/LibreOffice.app/Contents/MacOS/soffice",
		"/usr/bin/soffice",
		"/usr/local/bin/soffice",
		"/snap/bin/soffice",
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "soffice"
}
