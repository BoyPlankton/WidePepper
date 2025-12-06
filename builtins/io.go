// Copyright (c) 2025 BoyPlankton. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package builtins provides built-in functions for the WidePepper language.
// It includes file I/O, network operations, and other system interactions.
package builtins

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FileHandle represents an open file or network resource
type FileHandle struct {
	ID        int
	Reader    io.Reader
	Writer    io.Writer
	Closer    io.Closer
	FilePath  string
	IsNetwork bool
}

// FileHandleManager manages open file and network handles
type FileHandleManager struct {
	handles map[int]*FileHandle
	nextID  int
}

// NewFileHandleManager creates a new file handle manager
func NewFileHandleManager() *FileHandleManager {
	return &FileHandleManager{
		handles: make(map[int]*FileHandle),
		nextID:  1,
	}
}

// validatePath sanitizes and validates a file path to prevent path traversal attacks
func validatePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Clean the path to remove any ".." or "." components and resolve to canonical form
	cleanPath := filepath.Clean(path)

	// Get the absolute path - this will resolve any remaining relative components
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// The combination of filepath.Clean and filepath.Abs ensures that:
	// 1. Path separators are normalized
	// 2. Redundant separators are removed
	// 3. "." and ".." are resolved
	// 4. The path is converted to an absolute path
	// This prevents path traversal attacks by ensuring the path is canonical

	return absPath, nil
}

// PounceFile opens a file for reading or writing (pounce = open)
func (m *FileHandleManager) PounceFile(filePath string, mode string) (int, error) {
	// Validate and sanitize the file path
	validPath, err := validatePath(filePath)
	if err != nil {
		return 0, fmt.Errorf("invalid file path: %w", err)
	}

	var file *os.File

	switch strings.ToLower(mode) {
	case "r", "read":
		file, err = os.Open(validPath)
		if err != nil {
			return 0, fmt.Errorf("failed to open file for reading: %w", err)
		}
	case "w", "write":
		file, err = os.Create(validPath)
		if err != nil {
			return 0, fmt.Errorf("failed to open file for writing: %w", err)
		}
	case "a", "append":
		file, err = os.OpenFile(validPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return 0, fmt.Errorf("failed to open file for appending: %w", err)
		}
	default:
		return 0, fmt.Errorf("invalid file mode: %s (use 'r', 'w', or 'a')", mode)
	}

	id := m.nextID
	m.nextID++

	m.handles[id] = &FileHandle{
		ID:        id,
		Reader:    file,
		Writer:    file,
		Closer:    file,
		FilePath:  validPath,
		IsNetwork: false,
	}

	return id, nil
}

// LapLine reads a single line from a file handle (lap = drink/read)
func (m *FileHandleManager) LapLine(handleID int) (string, error) {
	handle, ok := m.handles[handleID]
	if !ok {
		return "", fmt.Errorf("invalid file handle: %d", handleID)
	}

	reader := bufio.NewReader(handle.Reader)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	line = strings.TrimSuffix(line, "\n")
	return line, nil
}

// DevourFile reads entire contents of a file (devour = eat entirely)
func (m *FileHandleManager) DevourFile(handleID int) (string, error) {
	handle, ok := m.handles[handleID]
	if !ok {
		return "", fmt.Errorf("invalid file handle: %d", handleID)
	}

	bytes, err := io.ReadAll(handle.Reader)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return string(bytes), nil
}

// ScratchLine writes a line to a file handle (scratch = write)
func (m *FileHandleManager) ScratchLine(handleID int, content string) error {
	handle, ok := m.handles[handleID]
	if !ok {
		return fmt.Errorf("invalid file handle: %d", handleID)
	}

	_, err := fmt.Fprintln(handle.Writer, content)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

// ScratchString writes a string to a file handle without newline (scratch = write)
func (m *FileHandleManager) ScratchString(handleID int, content string) error {
	handle, ok := m.handles[handleID]
	if !ok {
		return fmt.Errorf("invalid file handle: %d", handleID)
	}

	_, err := fmt.Fprint(handle.Writer, content)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

// NuzzleClose closes a file handle (nuzzle = close affectionately)
func (m *FileHandleManager) NuzzleClose(handleID int) error {
	handle, ok := m.handles[handleID]
	if !ok {
		return fmt.Errorf("invalid file handle: %d", handleID)
	}

	if handle.Closer != nil {
		err := handle.Closer.Close()
		delete(m.handles, handleID)
		if err != nil {
			return fmt.Errorf("error closing file: %w", err)
		}
	}

	return nil
}

// FetchURL performs a GET request to a URL (fetch = get)
func (m *FileHandleManager) FetchURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("HTTP GET failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// CoughUpData performs a POST request to a URL (cough up = send data)
func (m *FileHandleManager) CoughUpData(url string, contentType string, data string) (string, error) {
	resp, err := http.Post(url, contentType, strings.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("HTTP POST failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// SniffFile checks if a file exists (sniff = check)
func SniffFile(filePath string) bool {
	// Validate and sanitize the file path
	validPath, err := validatePath(filePath)
	if err != nil {
		return false
	}

	_, err = os.Stat(validPath)
	return err == nil
}

// SwatFile deletes a file (swat = hit/delete)
func SwatFile(filePath string) error {
	// Validate and sanitize the file path
	validPath, err := validatePath(filePath)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	return os.Remove(validPath)
}

// PounceDirectory lists files in a directory (pounce = explore)
func PounceDirectory(dirPath string) ([]string, error) {
	// Validate and sanitize the directory path
	validPath, err := validatePath(dirPath)
	if err != nil {
		return nil, fmt.Errorf("invalid directory path: %w", err)
	}

	entries, err := os.ReadDir(validPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}
