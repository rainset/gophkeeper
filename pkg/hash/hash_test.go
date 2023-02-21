package hash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateRandom(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "positive test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandom(tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}

func TestMd5(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test",
			args: args{text: "hello"},
			want: "5d41402abc4b2a76b9719d911017c592",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Md5(tt.args.text); got != tt.want {
				t.Errorf("Md5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSha256(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test",
			args: args{text: "hello"},
			want: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sha256(tt.args.text); got != tt.want {
				t.Errorf("Sha256() = %v, want %v", got, tt.want)
			}
		})
	}
}
