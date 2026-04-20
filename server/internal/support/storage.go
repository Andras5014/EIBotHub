package support

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ObjectStorage interface {
	SaveUploadedFile(resourceType string, fileHeader *multipart.FileHeader) (string, string, error)
	WriteSeedFile(resourceType, fileName, content string) (string, error)
	WriteGeneratedFile(resourceType, fileName, content string) (string, error)
	ResolvePath(relativePath string) string
	SplitFile(sourceRelativePath, targetDir, baseName string, parts int) ([]string, error)
}

type LocalStorage struct {
	root string
}

func NewLocalStorage(root string) *LocalStorage {
	return &LocalStorage{root: root}
}

func (s *LocalStorage) SaveUploadedFile(resourceType string, fileHeader *multipart.FileHeader) (string, string, error) {
	if fileHeader == nil {
		return "", "", nil
	}

	dir := filepath.Join(s.root, resourceType)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", "", err
	}

	safeName := strings.ReplaceAll(fileHeader.Filename, " ", "_")
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), safeName)
	relativePath := filepath.Join(resourceType, fileName)
	fullPath := filepath.Join(s.root, relativePath)

	src, err := fileHeader.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", "", err
	}

	return filepath.ToSlash(relativePath), fileHeader.Filename, nil
}

func (s *LocalStorage) WriteSeedFile(resourceType, fileName, content string) (string, error) {
	return s.WriteGeneratedFile(resourceType, fileName, content)
}

func (s *LocalStorage) WriteGeneratedFile(resourceType, fileName, content string) (string, error) {
	dir := filepath.Join(s.root, resourceType)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	relativePath := filepath.Join(resourceType, fileName)
	if err := os.WriteFile(filepath.Join(s.root, relativePath), []byte(content), 0o644); err != nil {
		return "", err
	}

	return filepath.ToSlash(relativePath), nil
}

func (s *LocalStorage) ResolvePath(relativePath string) string {
	return filepath.Join(s.root, filepath.FromSlash(relativePath))
}

func (s *LocalStorage) SplitFile(sourceRelativePath, targetDir, baseName string, parts int) ([]string, error) {
	if parts <= 0 {
		parts = 1
	}
	sourcePath := s.ResolvePath(sourceRelativePath)
	file, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if info.Size() == 0 {
		return nil, fmt.Errorf("source file is empty")
	}

	chunkSize := info.Size() / int64(parts)
	if info.Size()%int64(parts) != 0 {
		chunkSize++
	}

	dir := filepath.Join(s.root, targetDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	ext := filepath.Ext(baseName)
	name := strings.TrimSuffix(baseName, ext)
	relativePaths := make([]string, 0, parts)
	for index := 0; index < parts; index++ {
		partName := fmt.Sprintf("%s.part%03d%s", name, index+1, ext)
		relativePath := filepath.Join(targetDir, partName)
		fullPath := filepath.Join(s.root, relativePath)

		dst, err := os.Create(fullPath)
		if err != nil {
			return nil, err
		}

		written, copyErr := io.CopyN(dst, file, chunkSize)
		closeErr := dst.Close()
		if closeErr != nil {
			return nil, closeErr
		}
		if copyErr != nil && copyErr != io.EOF {
			return nil, copyErr
		}
		if written == 0 {
			_ = os.Remove(fullPath)
			break
		}

		relativePaths = append(relativePaths, filepath.ToSlash(relativePath))
		if copyErr == io.EOF {
			break
		}
	}

	return relativePaths, nil
}
