package auth

import (
	"testing"
	"time"
)

func TestManager_NewJWT(t *testing.T) {
	type fields struct {
		signingKey string
	}
	type args struct {
		userID string
		ttl    time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "test new jwt",
			fields: fields{signingKey: "secret_key"},
			args:   args{ttl: time.Duration(100), userID: "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				signingKey: tt.fields.signingKey,
			}
			_, err := m.NewJWT(tt.args.userID, tt.args.ttl)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestManager_NewRefreshToken(t *testing.T) {
	type fields struct {
		signingKey string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:   "test new refresh token",
			fields: fields{signingKey: "secret_key"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				signingKey: tt.fields.signingKey,
			}
			_, err := m.NewRefreshToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestManager_Parse(t *testing.T) {

	manager, _ := NewManager("secret_key")
	accessToken, _ := manager.NewJWT("1", 1000)

	type fields struct {
		signingKey string
	}
	type args struct {
		accessToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "manager parse",
			fields: fields{signingKey: "secret_key"},
			args:   args{accessToken: accessToken},
			want:   "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				signingKey: tt.fields.signingKey,
			}
			got, err := m.Parse(tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewManager(t *testing.T) {
	type args struct {
		signingKey string
	}
	tests := []struct {
		name    string
		args    args
		want    *Manager
		wantErr bool
	}{
		{
			name: "new manager",
			args: args{signingKey: "secret_key"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewManager(tt.args.signingKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
