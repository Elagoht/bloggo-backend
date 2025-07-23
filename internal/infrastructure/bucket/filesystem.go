package bucket

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
)

type fileSystemBucket struct {
	RootDir string
	mutex   sync.RWMutex
}

func NewFileSystemBucket(subDir string) (Bucket, error) {
	// Get working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Set main storage as 'cwd/storage'
	baseBucketDir := filepath.Join(cwd, "storage")

	// Calculate final, absolute path
	finalRootDir := baseBucketDir
	if subDir != "" {
		finalRootDir = filepath.Join(baseBucketDir, subDir)
	}
	absRootDir, err := filepath.Abs(finalRootDir)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get absolute path for root directory %s: %w",
			finalRootDir,
			err,
		)
	}

	// Create all directories exists
	if err := os.MkdirAll(absRootDir, 0755); err != nil {
		return nil, fmt.Errorf(
			"failed to create root directory %s: %w",
			absRootDir,
			err,
		)
	}

	return &fileSystemBucket{
		RootDir: absRootDir,
	}, nil
}

// Save byte array to a file with given name
func (bucket *fileSystemBucket) Save(file []byte, name string) error {
	if name == "" {
		return errors.New("File name cannot be null")
	}

	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()

	if err := bucket.makeSureRootDirExists(); err != nil {
		return err
	}

	filePath := filepath.Join(bucket.RootDir, name)

	// Save file
	if err := os.WriteFile(filePath, file, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Save a file got by multipart form data
func (bucket *fileSystemBucket) SaveMultiPart(
	file multipart.File,
	name string,
) error {
	if name == "" {
		return errors.New("File name cannot be null")
	}

	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()

	if err := bucket.makeSureRootDirExists(); err != nil {
		return err
	}

	filePath := filepath.Join(bucket.RootDir, name)

	// Create file
	out, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Save it
	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Read the file with given name
func (bucket *fileSystemBucket) Get(name string) ([]byte, error) {
	if name == "" {
		return nil, errors.New("File name cannot be null")
	}

	bucket.mutex.RLock()
	defer bucket.mutex.RUnlock()

	filePath := filepath.Join(bucket.RootDir, name)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// Delete the file with given name
func (bucket *fileSystemBucket) Delete(name string) error {
	if name == "" {
		return errors.New("File name cannot be null")
	}
	filePath := filepath.Join(bucket.RootDir, name)

	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Delete files matching with pattern
// whilst keeping matching with blacklist pattern
func (bucket *fileSystemBucket) DeleteMatching(
	pattern string,
	blacklist string,
) error {
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()

	// Merge pattern with root directory
	patternPath := filepath.Join(bucket.RootDir, pattern)

	// Find matching files
	matches, err := filepath.Glob(patternPath)
	if err != nil {
		return fmt.Errorf("failed to match pattern: %w", err)
	}

	// Merge blacklist pattern with root directory if not empty
	var blacklistPath string
	if blacklist != "" {
		blacklistPath = filepath.Join(bucket.RootDir, blacklist)
	}

	// Delete matcing files
	for _, filePath := range matches {
		// Check if they are a file...
		info, err := os.Stat(filePath)
		if err != nil {
			return fmt.Errorf("failed to stat file %s: %w", filePath, err)
		}
		if info.IsDir() {
			continue // ...and not a directory
		}

		// Skip blacklisted files
		if blacklist != "" {
			matched, _ := filepath.Match(blacklistPath, filePath)
			if matched {
				continue
			}
		}

		// Delete files
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("failed to delete file %s: %w", filePath, err)
		}
	}

	return nil
}

func (bucket *fileSystemBucket) makeSureRootDirExists() error {
	if err := os.MkdirAll(bucket.RootDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}
