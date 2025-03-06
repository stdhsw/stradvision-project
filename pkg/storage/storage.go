package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Handler struct {
	currentFile *os.File
	name        string
	path        string

	currentCount int
	maxFileSize  int
	maxFileCount int
}

// NewHandler 새로운 Handler 생성
func NewHandler(name, path string, options ...Option) (*Handler, error) {
	config := fromOptions(options...)

	handler := &Handler{
		name:         name,
		path:         path,
		maxFileSize:  config.maxFileSize,
		maxFileCount: config.maxFileCount,
	}

	// 마지막 파일 번호 추출
	files, err := handler.GetSortFileList()
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		handler.currentCount = 0
	} else {
		lastCount := extractNumber(files[len(files)-1])
		handler.currentCount = lastCount
	}

	return handler, nil
}

// GetCurrentFile 현재 파일명을 반환
func (h *Handler) GetCurrentFile() string {
	return h.currentFile.Name()
}

// GetFileList 파일 목록을 반환
func (h *Handler) GetFileList() ([]string, error) {
	entries, err := os.ReadDir(h.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.Contains(entry.Name(), h.name) {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// GetSortFileList 파일 목록을 정렬해서 반환
func (h *Handler) GetSortFileList() ([]string, error) {
	files, err := h.GetFileList()
	if err != nil {
		return nil, err
	}

	SortByNumericSuffix(files)

	return files, nil
}

// RemoveFile 파일 삭제
func (h *Handler) RemoveFile(file string) error {
	if err := os.Remove(filepath.Join(h.path, file)); err != nil {
		return err
	}

	return nil
}

// RemoveFiles 파일 목록을 받아서 삭제
func (h *Handler) RemoveFiles(files []string) error {
	for _, file := range files {
		if err := h.RemoveFile(file); err != nil {
			return fmt.Errorf("failed to remove %s: %w", file, err)
		}
	}

	return nil
}

// CheckAndRemove 파일 개수가 maxFileCount를 넘으면 파일을 삭제
func (h *Handler) CheckAndRemove() error {
	files, err := h.GetSortFileList()
	if err != nil {
		return err
	}

	if len(files) >= h.maxFileCount {
		count := len(files) - h.maxFileCount + 1
		for i := 0; i < count; i++ {
			if err := h.RemoveFile(files[i]); err != nil {
				return fmt.Errorf("failed to remove file: %w", err)
			}
		}
	}

	return nil
}

// WriteData 파일에 데이터를 기록
func (h *Handler) WriteData(data []byte) error {
	// 파일이 없으면 새로 생성
	if h.currentFile == nil {
		file, err := os.Create(filepath.Join(h.path, fmt.Sprintf("%s_%d", h.name, h.currentCount)))
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		h.currentFile = file
	}

	// 파일 사이즈 체크
	// 파일의 최대 개수가 넘으면 파일을 닫고 새로 생성
	fInfo, _ := h.currentFile.Stat()
	if fInfo.Size() > int64(h.maxFileSize) {
		if err := h.CheckAndRemove(); err != nil {
			return fmt.Errorf("failed to check and remove: %w", err)
		}

		if err := h.currentFile.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}

		h.currentCount++
		file, err := os.Create(filepath.Join(h.path, fmt.Sprintf("%s_%d", h.name, h.currentCount)))
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		h.currentFile = file
	}

	// 데이터 기록
	if _, err := h.currentFile.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}
