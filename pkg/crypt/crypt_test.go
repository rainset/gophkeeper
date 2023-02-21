package crypt

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestDecodeBase64(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "TestDecodeBase64",
			args: args{s: "aGVsbG8="},
			want: []byte("hello"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeBase64(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeBase64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeBase64(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestDecodeBase64",
			args: args{b: []byte("hello")},
			want: "aGVsbG8=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeBase64(tt.args.b); got != tt.want {
				t.Errorf("EncodeBase64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	type args struct {
		data []byte
		key  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive test",
			args:    args{data: []byte("data to be encrypted"), key: []byte("O21X4YFsXTkPSs8ZBe3RN7MFmZPo64wZ")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.data, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		data []byte
		key  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "positive test",
			args:    args{data: []byte("data to be encrypted"), key: []byte("O21X4YFsXTkPSs8ZBe3RN7MFmZPo64wZ")},
			want:    []byte("data to be encrypted"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := Encrypt(tt.args.data, tt.args.key)
			assert.NoError(t, err)

			got, err := Decrypt(encrypted, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
