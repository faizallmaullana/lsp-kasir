package services

import (
	"encoding/base64"
	"errors"
	"faizalmaulana/lsp/helper"
	"faizalmaulana/lsp/models/entity"
	"faizalmaulana/lsp/models/repo"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type ImagesService interface {
	UploadBlob(fileName, contentType string, data []byte) (string, string, error) // returns (id, storedFilename)
	UploadBase64(fileName, contentType, b64 string) (string, string, error)       // returns (id, storedFilename)
	GetBlob(id string) (*entity.Images, error)
	GetBase64(id string) (string, string, string, error) // id, contentType, base64
	Delete(id string) error
}

type imagesService struct{ repo repo.ImagesRepo }

func NewImagesService(r repo.ImagesRepo) ImagesService { return &imagesService{repo: r} }

func (s *imagesService) UploadBlob(fileName, contentType string, data []byte) (string, string, error) {
	if len(data) == 0 {
		return "", "", errors.New("empty data")
	}
	storedName := generateFileName(fileName, contentType)
	fullPath, err := ensureStoragePath(storedName)
	if err != nil {
		return "", "", err
	}
	if err := os.WriteFile(fullPath, data, fs.FileMode(0644)); err != nil {
		return "", "", err
	}
	img := &entity.Images{IdImage: helper.Uuid(), FileName: storedName, ContentType: contentType, Size: int64(len(data))}
	if err := s.repo.Create(img); err != nil {
		return "", "", err
	}
	return img.IdImage, storedName, nil
}

func (s *imagesService) UploadBase64(fileName, contentType, b64 string) (string, string, error) {
	if b64 == "" {
		return "", "", errors.New("empty base64")
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", "", err
	}
	return s.UploadBlob(fileName, contentType, raw)
}

func (s *imagesService) GetBlob(id string) (*entity.Images, error) {
	fmt.Println("GetBlob id:", id)
	meta, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	fmt.Println(meta.FileName)

	fullPath := filepath.Join("storages", "images", id)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	meta.Size = int64(len(data))
	meta.Data = data
	return meta, nil
}

func (s *imagesService) GetBase64(id string) (string, string, string, error) {
	img, err := s.GetBlob(id)
	if err != nil {
		return "", "", "", err
	}
	b64 := base64.StdEncoding.EncodeToString(img.Data)
	return img.IdImage, img.ContentType, b64, nil
}

func (s *imagesService) Delete(id string) error { return s.repo.Delete(id) }

func ensureStoragePath(fileName string) (string, error) {
	base := filepath.Join("storages", "images")
	if err := os.MkdirAll(base, 0755); err != nil {
		return "", err
	}
	return filepath.Join(base, fileName), nil
}

func generateFileName(original, contentType string) string {
	ext := ""
	if contentType != "" {
		switch strings.ToLower(contentType) {
		case "image/png":
			ext = ".png"
		case "image/jpeg", "image/jpg":
			ext = ".jpg"
		case "image/gif":
			ext = ".gif"
		default:
			parts := strings.Split(contentType, "/")
			if len(parts) == 2 {
				ext = "." + parts[1]
			}
		}
	}
	if ext == "" {
		ext = filepath.Ext(original)
		if ext == "" {
			ext = ".bin"
		}
	}
	return fmt.Sprintf("%s%s", helper.Uuid(), ext)
}
