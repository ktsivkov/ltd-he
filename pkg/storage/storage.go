package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func New() *Storage {
	return &Storage{
		mu: sync.Mutex{},
	}
}

type Storage struct {
	mu sync.Mutex
}

func (s *Storage) WriteFile(payload []byte, pathSegments ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.writeFile(payload, pathSegments...)
}

func (s *Storage) DeletePath(pathSegments ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.deletePath(pathSegments...)
}

func (s *Storage) Exists(pathSegments ...string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.exists(pathSegments...)
}

func (s *Storage) CreateDir(pathSegments ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.createDir(pathSegments...)
}

func (s *Storage) ReadFile(pathSegments ...string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.readFile(pathSegments...)
}

func (s *Storage) ReadDir(pathSegments ...string) ([]os.DirEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.readDir(pathSegments...)
}

func (s *Storage) writeFile(payload []byte, pathSegments ...string) error {
	if len(pathSegments) == 0 {
		return fmt.Errorf("no path segments provided")
	}

	fp := filepath.Join(pathSegments...)

	dirPath := filepath.Dir(fp)
	ok, err := s.exists(dirPath)
	if err != nil {
		return fmt.Errorf("could not check if file path exists: %v", err)
	}

	if !ok {
		if err := s.createDir(dirPath); err != nil {
			return fmt.Errorf("could not file path directory: %v", err)
		}
	}

	if err := os.WriteFile(fp, payload, os.ModePerm); err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

func (s *Storage) deletePath(pathSegments ...string) error {
	if len(pathSegments) == 0 {
		return fmt.Errorf("no path segments provided")
	}

	fp := filepath.Join(pathSegments...)
	ok, err := s.exists(fp)
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

func (s *Storage) exists(pathSegments ...string) (bool, error) {
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

func (s *Storage) createDir(pathSegments ...string) error {
	if len(pathSegments) == 0 {
		return fmt.Errorf("no path segments provided")
	}

	fp := filepath.Join(pathSegments...)

	if err := os.MkdirAll(fp, os.ModePerm); err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}

	return nil
}

func (s *Storage) readFile(pathSegments ...string) ([]byte, error) {
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

func (s *Storage) readDir(pathSegments ...string) ([]os.DirEntry, error) {
	if len(pathSegments) == 0 {
		return nil, fmt.Errorf("no path segments provided")
	}

	return os.ReadDir(filepath.Join(pathSegments...))
}
