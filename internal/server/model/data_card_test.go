package model

import (
	"testing"
	"time"
)

func TestDataCard_Validate(t *testing.T) {
	type fields struct {
		ID        int
		UserID    int
		Title     string
		Number    string
		Date      string
		Cvv       string
		Meta      string
		UpdatedAt time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "card model",
			fields: fields{
				ID:        0,
				UserID:    1,
				Title:     "title",
				Number:    "123123123",
				Date:      "11/23",
				Cvv:       "123",
				Meta:      "meta",
				UpdatedAt: time.Time{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataCard{
				ID:        tt.fields.ID,
				UserID:    tt.fields.UserID,
				Title:     tt.fields.Title,
				Number:    tt.fields.Number,
				Date:      tt.fields.Date,
				Cvv:       tt.fields.Cvv,
				Meta:      tt.fields.Meta,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			if err := d.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
