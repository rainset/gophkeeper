package model

import (
	"testing"
	"time"
)

func TestDataFile_Validate(t *testing.T) {
	type fields struct {
		ID        int
		UserID    int
		Title     string
		Filename  string
		Path      string
		Meta      string
		UpdatedAt time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "file model",
			fields: fields{
				ID:        0,
				UserID:    1,
				Title:     "title",
				Filename:  "test.png",
				Path:      "/test.png",
				Meta:      "test",
				UpdatedAt: time.Time{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataFile{
				ID:        tt.fields.ID,
				UserID:    tt.fields.UserID,
				Title:     tt.fields.Title,
				Filename:  tt.fields.Filename,
				Path:      tt.fields.Path,
				Meta:      tt.fields.Meta,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			if err := d.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
