package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManagerReloadAndSnapshot(t *testing.T) {
	tmpDir := t.TempDir()
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})

	writeConfigFiles(t, tmpDir,
		`port = 8080
mode = "release"
`,
		`path = "storage/db/app.sqlite"
max_open_conns = 5
`,
		`level = "info"
file = "storage/tmp/app.log"
enable_console = true
`)

	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	snap := mgr.Snapshot()
	if snap.Server.Mode != "release" || snap.Server.Port != 8080 {
		t.Fatalf("unexpected server config: %+v", snap.Server)
	}
	if snap.Database.Path != "storage/db/app.db" || snap.Database.MaxOpenConns != 5 {
		t.Fatalf("unexpected database config: %+v", snap.Database)
	}
	if snap.Logging.Level != "info" || snap.Logging.File != "storage/tmp/app.log" {
		t.Fatalf("unexpected logging config: %+v", snap.Logging)
	}

	snap.Server.Mode = "debug"
	newSnap := mgr.Snapshot()
	if newSnap.Server.Mode != "release" {
		t.Fatalf("Snapshot should return cloned config, got: %s", newSnap.Server.Mode)
	}
}

func TestManagerReloadErrorOnInvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})

	// 缺少 server.port（为 0），应触发校验错误。
	writeConfigFiles(t, tmpDir,
		`port = 0
mode = "release"
`,
		`path = "storage/db/app.db"
max_open_conns = 5
`,
		`level = "info"
file = "storage/tmp/app.log"
enable_console = true
`)

	if _, err := NewManager(); err == nil {
		t.Fatalf("expected NewManager() to fail on invalid config")
	}
}

func writeConfigFiles(t *testing.T, baseDir, serverContent, databaseContent, loggingContent string) {
	t.Helper()
	configDir := filepath.Join(baseDir, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "server.toml"), []byte(serverContent), 0o644); err != nil {
		t.Fatalf("write server.toml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "database.toml"), []byte(databaseContent), 0o644); err != nil {
		t.Fatalf("write database.toml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "logging.toml"), []byte(loggingContent), 0o644); err != nil {
		t.Fatalf("write logging.toml: %v", err)
	}
}
