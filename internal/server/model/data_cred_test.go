package model

import (
	"testing"
	"time"
)

func TestDataCred_Validate(t *testing.T) {
	type fields struct {
		ID        int
		UserID    int
		Title     string
		Username  string
		Password  string
		Meta      string
		UpdatedAt time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "cred model",
			fields: fields{
				ID:        0,
				UserID:    1,
				Title:     "title",
				Username:  "login",
				Password:  "12345",
				Meta:      "test",
				UpdatedAt: time.Time{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataCred{
				ID:        tt.fields.ID,
				UserID:    tt.fields.UserID,
				Title:     tt.fields.Title,
				Username:  tt.fields.Username,
				Password:  tt.fields.Password,
				Meta:      tt.fields.Meta,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			if err := d.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
