package config

import "testing"

func TestValidateServer(t *testing.T) {
	tests := []struct {
		name    string
		cfg     ServerConfig
		wantErr bool
	}{
		{
			name: "valid",
			cfg: ServerConfig{
				Port: 8080,
				Mode: "release",
			},
		},
		{
			name: "invalid port",
			cfg: ServerConfig{
				Port: 0,
				Mode: "release",
			},
			wantErr: true,
		},
		{
			name: "missing mode",
			cfg: ServerConfig{
				Port: 8080,
				Mode: " ",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateServer(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDatabase(t *testing.T) {
	tests := []struct {
		name    string
		cfg     DatabaseConfig
		wantErr bool
	}{
		{
			name: "valid",
			cfg: DatabaseConfig{
				Path:         "storage/db/app.sqlite",
				MaxOpenConns: 10,
			},
		},
		{
			name: "missing path",
			cfg: DatabaseConfig{
				Path:         " ",
				MaxOpenConns: 10,
			},
			wantErr: true,
		},
		{
			name: "invalid max_open_conns",
			cfg: DatabaseConfig{
				Path:         "storage/db/app.sqlite",
				MaxOpenConns: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDatabase(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLogging(t *testing.T) {
	tests := []struct {
		name    string
		cfg     LoggingConfig
		wantErr bool
	}{
		{
			name: "valid",
			cfg: LoggingConfig{
				Level:         "info",
				File:          "storage/tmp/app.log",
				EnableConsole: true,
			},
		},
		{
			name: "missing level",
			cfg: LoggingConfig{
				Level:         " ",
				File:          "storage/tmp/app.log",
				EnableConsole: true,
			},
			wantErr: true,
		},
		{
			name: "missing file",
			cfg: LoggingConfig{
				Level:         "info",
				File:          " ",
				EnableConsole: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLogging(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateLogging() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
