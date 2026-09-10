package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr              string
	DatabaseURL       string
	CORSOrigin        string
	AuthUsername      string
	AuthPassword      string
	CursorAgentBin    string
	CursorAgentMode   string
	CursorWorkspace   string
	CursorTimeoutSec  int
	MigrationsDir     string
}

func LoadDotEnv(path string) {
	_ = godotenv.Load(path)
}

func Load() Config {
	timeout := 300
	if v := strings.TrimSpace(os.Getenv("CURSOR_AGENT_TIMEOUT_SEC")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			timeout = n
		}
	}

	bin := strings.TrimSpace(os.Getenv("CURSOR_AGENT_BIN"))
	if bin == "" {
		bin = `C:\Users\armin\AppData\Local\cursor-agent\agent.cmd`
	}

	mode := strings.TrimSpace(os.Getenv("CURSOR_AGENT_MODE"))
	if mode == "" {
		mode = "ask"
	}

	addr := strings.TrimSpace(os.Getenv("ADDR"))
	if addr == "" {
		addr = ":8201"
	}

	user := strings.TrimSpace(os.Getenv("AUTH_USERNAME"))
	if user == "" {
		user = "armin"
	}
	pass := strings.TrimSpace(os.Getenv("AUTH_PASSWORD"))
	if pass == "" {
		pass = "dopadopa123"
	}

	mig := strings.TrimSpace(os.Getenv("MIGRATIONS_DIR"))
	if mig == "" {
		mig = "migrations"
	}

	cors := strings.TrimSpace(os.Getenv("CORS_ORIGIN"))
	if cors == "" {
		cors = "http://localhost:5201"
	}

	return Config{
		Addr:             addr,
		DatabaseURL:      strings.TrimSpace(os.Getenv("DATABASE_URL")),
		CORSOrigin:       cors,
		AuthUsername:     user,
		AuthPassword:     pass,
		CursorAgentBin:   bin,
		CursorAgentMode:  mode,
		CursorWorkspace:  strings.TrimSpace(os.Getenv("CURSOR_WORKSPACE")),
		CursorTimeoutSec: timeout,
		MigrationsDir:    mig,
	}
}
