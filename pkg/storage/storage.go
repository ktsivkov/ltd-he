package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

func New() *Storage {
	return &Storage{}
}

type Storage struct {
}

func (s *Storage) WriteFile(payload []byte, pathSegments ...string) error {
	if len(pathSegments) == 0 {
		return fmt.Errorf("no path segments provided")
	}

	fp := filepath.Join(pathSegments...)

	dirPath := filepath.Dir(fp)
	ok, err := s.Exists(dirPath)
	if err != nil {
		return fmt.Errorf("could not check if file path exists: %v", err)
	}

	if !ok {
		if err := s.CreateDir(dirPath); err != nil {
			return fmt.Errorf("could not file path directory: %v", err)
		}
	}

	if err := os.WriteFile(fp, payload, os.ModePerm); err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

func (s *Storage) DeletePath(pathSegments ...string) error {
	if len(pathSegments) == 0 {
		return fmt.Errorf("no path segments provided")
	}

	fp := filepath.Join(pathSegments...)
	ok, err := s.Exists(fp)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	if err := os.RemoveAll(fp); err != nil {
		return fmt.Errorf("could not delete path: %w", err)
	}

	return nil
}

func (s *Storage) Exists(pathSegments ...string) (bool, error) {
	if len(pathSegments) == 0 {
		return false, fmt.Errorf("no path segments provided")
	}

	fp := filepath.Join(pathSegments...)

	_, err := os.Stat(fp)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, fmt.Errorf("could not check if path exists: %w", err)
}

func (s *Storage) CreateDir(pathSegments ...string) error {
	if len(pathSegments) == 0 {
		return fmt.Errorf("no path segments provided")
	}

	fp := filepath.Join(pathSegments...)

	if err := os.MkdirAll(fp, os.ModePerm); err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}

	return nil
}

func (s *Storage) ReadFile(pathSegments ...string) ([]byte, error) {
	if len(pathSegments) == 0 {
		return nil, fmt.Errorf("no path segments provided")
	}

	fp := filepath.Join(pathSegments...)

	bytes, err := os.ReadFile(fp)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	return bytes, nil
}

func (s *Storage) ReadDir(pathSegments ...string) ([]os.DirEntry, error) {
	if len(pathSegments) == 0 {
		return nil, fmt.Errorf("no path segments provided")
	}

	return os.ReadDir(filepath.Join(pathSegments...))
}
