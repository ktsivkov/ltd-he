package history

import "os"

type StorageDriver interface {
	WriteFile(payload []byte, pathSegments ...string) error
	DeletePath(pathSegments ...string) error
	Exists(pathSegments ...string) (bool, error)
	CreateDir(pathSegments ...string) error
	ReadFile(pathSegments ...string) ([]byte, error)
	ReadDir(pathSegments ...string) ([]os.DirEntry, error)
}
