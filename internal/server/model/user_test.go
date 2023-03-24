package model

import "testing"

func TestUser_Validate(t *testing.T) {
	type fields struct {
		ID       int
		Login    string
		Password string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "user model",
			fields: fields{
				ID:       0,
				Login:    "login",
				Password: "12345",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:       tt.fields.ID,
				Login:    tt.fields.Login,
				Password: tt.fields.Password,
			}
			if err := u.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
