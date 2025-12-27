package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			want: &Config{
				Server: ServerConfig{
					Port: 8080,
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "tsunagu",
					Password: "tsunagu_password",
					DBName:   "tsunagu_db",
					SSLMode:  "disable",
				},
				JWT: JWTConfig{
					Secret:      "your-secret-key-change-in-production",
					ExpiryHours: 24,
				},
				Zitadel: ZitadelConfig{
					URL:          "",
					ClientID:     "",
					ClientSecret: "",
				},
			},
			wantErr: false,
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"SERVER_PORT":           "3000",
				"DB_HOST":               "db.example.com",
				"DB_PORT":               "5433",
				"DB_USER":               "customuser",
				"DB_PASSWORD":           "custompass",
				"DB_NAME":               "customdb",
				"DB_SSLMODE":            "require",
				"JWT_SECRET":            "custom-secret",
				"JWT_EXPIRY_HOURS":      "48",
				"ZITADEL_URL":           "https://zitadel.example.com",
				"ZITADEL_CLIENT_ID":     "client123",
				"ZITADEL_CLIENT_SECRET": "secret123",
			},
			want: &Config{
				Server: ServerConfig{
					Port: 3000,
				},
				Database: DatabaseConfig{
					Host:     "db.example.com",
					Port:     5433,
					User:     "customuser",
					Password: "custompass",
					DBName:   "customdb",
					SSLMode:  "require",
				},
				JWT: JWTConfig{
					Secret:      "custom-secret",
					ExpiryHours: 48,
				},
				Zitadel: ZitadelConfig{
					URL:          "https://zitadel.example.com",
					ClientID:     "client123",
					ClientSecret: "secret123",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid server port",
			envVars: map[string]string{
				"SERVER_PORT": "invalid",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid db port",
			envVars: map[string]string{
				"DB_PORT": "invalid",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid jwt expiry hours",
			envVars: map[string]string{
				"JWT_EXPIRY_HOURS": "invalid",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			got, err := Load()

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDatabaseConfig_DSN(t *testing.T) {
	tests := []struct {
		name string
		cfg  DatabaseConfig
		want string
	}{
		{
			name: "default config",
			cfg: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "tsunagu",
				Password: "tsunagu_password",
				DBName:   "tsunagu_db",
				SSLMode:  "disable",
			},
			want: "host=localhost port=5432 user=tsunagu password=tsunagu_password dbname=tsunagu_db sslmode=disable",
		},
		{
			name: "custom config",
			cfg: DatabaseConfig{
				Host:     "db.example.com",
				Port:     5433,
				User:     "customuser",
				Password: "custompass",
				DBName:   "customdb",
				SSLMode:  "require",
			},
			want: "host=db.example.com port=5433 user=customuser password=custompass dbname=customdb sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.DSN()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		setEnv       bool
		want         string
	}{
		{
			name:         "use default value when env not set",
			key:          "TEST_KEY",
			defaultValue: "default",
			setEnv:       false,
			want:         "default",
		},
		{
			name:         "use env value when set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "custom",
			setEnv:       true,
			want:         "custom",
		},
		{
			name:         "use default when env is empty string",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "",
			setEnv:       true,
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.want, got)
		})
	}
}
