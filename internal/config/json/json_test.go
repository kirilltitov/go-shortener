package json

import (
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    *Config
		wantErr bool
	}{
		{
			name: "Positive",
			input: []byte(`
				{
				  "server_address": "localhost:8080",
				  "base_url": "http://localhost",
				  "file_storage_path": "/path/to/file.db",
				  "database_dsn": "postgres://postgres:mysecretpassword@127.0.0.1:5432/postgres",
				  "enable_https": true
				}
			`),
			want: &Config{
				ServerAddress:   "localhost:8080",
				BaseURL:         "http://localhost",
				FileStoragePath: "/path/to/file.db",
				DatabaseDSN:     "postgres://postgres:mysecretpassword@127.0.0.1:5432/postgres",
				EnableHTTPS:     true,
			},
			wantErr: false,
		},
		{
			name:    "No input",
			input:   nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid input",
			input:   []byte("not a json"),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}
