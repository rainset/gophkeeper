package service

import (
	"crypto/rand"
	"errors"
	"io"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"

	"github.com/rainset/gophkeeper/pkg/logger"
)

type FileService struct {
	path string
}

func New(path string) (repo *FileService, err error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, 0750)
		if err != nil {
			return nil, err
		}
	}

	repo = &FileService{path: path}

	return repo, nil
}

// SaveFile сохраняет файл пользователя

func (repo FileService) SaveFile(src io.ReadCloser, ext string) (filePath string, err error) {
	filePath, err = repo.getPath()
	if err != nil {
		logger.Error("getPath() ", err)

		return "", err
	}

	filePath += ext

	filePathAbs, err := filepath.Abs(filePath)
	if err != nil {
		logger.Error("filepath.Abs ", err)

		return filePath, err
	}

	dst, err := os.Create(filePathAbs)
	if err != nil {
		logger.Error("os.Create ", err)

		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		logger.Error("io.Copy ", err)

		return "", err
	}

	return filePathAbs, err
}

// DeleteFile - удаляет файл пользователя

func (repo FileService) DeleteFile(filePath string) (err error) {
	if filePath == "" {
		return nil
	}

	err = os.Remove(filePath)
	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	return nil
}

// GetFile - возвращает файл пользователя

func (repo FileService) GetFile(filePath string) (fileReader io.ReadCloser, err error) {
	return os.Open(filePath)
}

// Close - закрывает файловый репозиторий

func (repo FileService) Close() error {
	return nil
}

func (repo FileService) getPath() (string, error) {
	for {
		bytes, err := generateRandomBytes(32)
		if err != nil {
			return "", err
		}

		targetDir := filepath.Join(repo.path, string(bytes[0:2]), string(bytes[2:4]), string(bytes[4:6]), string(bytes[6:]))
		targetFile := filepath.Join(targetDir, string(bytes[6:]))

		if _, err := os.Stat(targetFile); errors.Is(err, os.ErrNotExist) {
			if _, err := os.Stat(targetDir); errors.Is(err, os.ErrNotExist) {
				err := os.MkdirAll(targetDir, 0750)
				if err != nil {
					return "", nil
				}
			}

			return targetFile, nil
		}
	}
}

func generateRandomBytes(n int) ([]byte, error) {
	const letters = "0123456789abcdef"

	result := make([]byte, n)

	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return nil, err
		}

		result[i] = letters[num.Int64()]
	}

	return result, nil
}
