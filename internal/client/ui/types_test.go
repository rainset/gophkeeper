package ui

import "testing"

func TestDataType_String(t *testing.T) {
	tests := []struct {
		name string
		t    DataType
		want string
	}{
		{
			name: "datatype",
			t:    TypeCard,
			want: "Карта",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTabName_String(t *testing.T) {
	tests := []struct {
		name string
		t    TabName
		want string
	}{
		{
			name: "tabname",
			t:    TabCard,
			want: "Карты",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
