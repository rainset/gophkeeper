package model

import (
	"testing"
	"time"
)

func TestDataText_Validate(t *testing.T) {
	type fields struct {
		ID        int
		UserID    int
		Title     string
		Text      string
		Meta      string
		UpdatedAt time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "text model",
			fields: fields{
				ID:        0,
				UserID:    1,
				Title:     "title",
				Text:      "text",
				Meta:      "test",
				UpdatedAt: time.Time{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataText{
				ID:        tt.fields.ID,
				UserID:    tt.fields.UserID,
				Title:     tt.fields.Title,
				Text:      tt.fields.Text,
				Meta:      tt.fields.Meta,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			if err := d.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
