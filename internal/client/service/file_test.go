package service

import (
	"io"
	"strings"
	"testing"
)

func TestFileService_Close(t *testing.T) {
	type fields struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "close",
			fields:  fields{path: "gophkeeper_files"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := FileService{
				path: tt.fields.path,
			}
			if err := repo.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileService_DeleteFile(t *testing.T) {
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
			name:    "delete",
			fields:  fields{path: "gophkeeper_files"},
			args:    args{filePath: "test_file.txt"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := FileService{
				path: tt.fields.path,
			}
			if err := repo.DeleteFile(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileService_GetFile(t *testing.T) {
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
			name:    "get file",
			fields:  fields{path: "gophkeeper_files"},
			args:    args{filePath: "no_file.txt"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := FileService{
				path: tt.fields.path,
			}
			_, err := repo.GetFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFileService_SaveFile(t *testing.T) {

	r := io.NopCloser(strings.NewReader("Hello, world!"))

	type fields struct {
		path string
	}
	type args struct {
		src io.ReadCloser
		ext string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantFilePath string
		wantErr      bool
	}{
		{
			name:    "save file",
			fields:  fields{path: "gophkeeper_files"},
			args:    args{src: r, ext: ".png"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := FileService{
				path: tt.fields.path,
			}
			_, err := repo.SaveFile(tt.args.src, tt.args.ext)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFileService_getPath(t *testing.T) {
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
			name:    "save file",
			fields:  fields{path: "gophkeeper_files"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := FileService{
				path: tt.fields.path,
			}
			_, err := repo.getPath()
			if (err != nil) != tt.wantErr {
				t.Errorf("getPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "new repo",
			args:    args{path: "gophkeeper_files"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.path)
			if err != nil {
				t.Error(err)
				return
			}

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
			name:    "generateRandomBytes",
			args:    args{n: 32},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := generateRandomBytes(tt.args.n)
			if err != nil {
				t.Error(err)
				return
			}
		})
	}
}
