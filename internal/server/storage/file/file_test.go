package file

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name     string
		args     args
		wantRepo *StorageFiles
		wantErr  bool
	}{
		{
			name: "new",
			args: args{path: "_file_storage"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.path)
			assert.NoError(t, err)
		})
	}
}

func TestStorageFiles_Close(t *testing.T) {
	type fields struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:   "close",
			fields: fields{path: "_file_storage"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := StorageFiles{
				path: tt.fields.path,
			}
			err := repo.Close()
			assert.NoError(t, err)
		})
	}
}

func TestStorageFiles_DeleteFile(t *testing.T) {
	type fields struct {
		path string
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "new",
			fields: fields{path: "_file_storage"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := StorageFiles{
				path: tt.fields.path,
			}
			err := repo.DeleteFile(tt.args.filePath)
			assert.NoError(t, err)
		})
	}
}

func TestStorageFiles_GetFile(t *testing.T) {
	type fields struct {
		path string
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantFileReader io.ReadCloser
		wantErr        bool
	}{
		{
			name:   "new",
			args:   args{filePath: "file.go"},
			fields: fields{path: "."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := StorageFiles{
				path: tt.fields.path,
			}
			_, err := repo.GetFile(tt.args.filePath)
			assert.NoError(t, err)

		})
	}
}

func TestStorageFiles_SaveFile(t *testing.T) {

	r := io.NopCloser(strings.NewReader("Hello, world!"))

	type fields struct {
		path string
	}
	type args struct {
		src io.ReadCloser
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantFilePath string
		wantErr      bool
	}{
		{
			name:   "save file",
			fields: fields{path: "_file_storage"},
			args:   args{src: r},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := StorageFiles{
				path: tt.fields.path,
			}
			_, err := repo.SaveFile(tt.args.src)
			assert.NoError(t, err)
		})
	}
}

func TestStorageFiles_getPath(t *testing.T) {
	type fields struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:   "get path",
			fields: fields{path: "_file_storage"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := StorageFiles{
				path: tt.fields.path,
			}
			_, err := repo.getPath()
			assert.NoError(t, err)
		})
	}
}

func Test_generateRandomBytes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "random generate",
			args: args{n: 32},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := generateRandomBytes(tt.args.n)
			assert.NoError(t, err)
		})
	}
}
