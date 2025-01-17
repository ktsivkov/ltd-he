package backup

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ktsivkov/ltd-he/pkg/player"
)

const backupTimestampFormat = "02-01-2006-150405"
const backupFolder = "backups"

func NewService(appFolder string) *Service {
	return &Service{
		appFolder: appFolder,
		mu:        &sync.Mutex{},
	}
}

type Service struct {
	appFolder string
	mu        *sync.Mutex
}

func (s *Service) BackupFolder(p *player.Player) string {
	return filepath.Join(s.appFolder, backupFolder, p.BattleTag)
}

func (s *Service) generateBackupFilename() string {
	return time.Now().Format(backupTimestampFormat)
}

func (s *Service) Backup(_ context.Context, p *player.Player) (*Backup, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	backupsPath := s.BackupFolder(p)
	_, err := os.Stat(backupsPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("could not check backup folder exists: %w", err)
		}
		if err := os.MkdirAll(backupsPath, os.ModePerm); err != nil {
			return nil, fmt.Errorf("could not create backup folder: %w", err)
		}
	}

	backupFile := filepath.Join(backupsPath, fmt.Sprintf("%s.zip", s.generateBackupFilename()))
	file, err := os.Create(backupFile)
	if err != nil {
		return nil, fmt.Errorf("could not create backup archive: %w", err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	if err := s.addDirToArchive(w, p.ReportFilePathAbsolute, p.ReportFilePathRelative); err != nil {
		return nil, fmt.Errorf("could not add report file path to archive: %w", err)
	}

	if err := s.addDirToArchive(w, p.LogsPathAbsolute, p.LogsPathRelative); err != nil {
		return nil, fmt.Errorf("could not add logs file path to archive: %w", err)
	}

	return &Backup{
		File: backupFile,
	}, nil
}

func (s *Service) addDirToArchive(w *zip.Writer, sourcePath string, targetPath string) error {
	return filepath.Walk(sourcePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}

		return s.addFileToArchive(w, path, filepath.Join(targetPath, relPath))
	})
}

func (s *Service) addFileToArchive(w *zip.Writer, sourcePath string, targetPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	target, err := w.Create(targetPath)
	if err != nil {
		return fmt.Errorf("could not create file in archive: %w", err)
	}

	_, err = io.Copy(target, source)
	if err != nil {
		return fmt.Errorf("could not copy file to archive: %w", err)
	}

	return nil
}
