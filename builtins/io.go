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
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// DefaultHTTPTimeout is the default timeout for HTTP requests
const DefaultHTTPTimeout = 30 * time.Second

// FileHandle represents an open file or network resource
type FileHandle struct {
	ID        int
	Reader    io.Reader
	Writer    io.Writer
	Closer    io.Closer
	FilePath  string
	IsNetwork bool
	BufReader *bufio.Reader
}

// FileHandleManager manages open file and network handles.
// The mutex protects concurrent access to the handles map and nextID counter.
// Note: The FileHandleManager does not synchronize access to individual file handles.
// If multiple goroutines need to read/write the same handle concurrently,
// the caller is responsible for appropriate synchronization.
type FileHandleManager struct {
	mu         sync.RWMutex
	handles    map[int]*FileHandle
	nextID     int
	httpClient *http.Client
}

// NewFileHandleManager creates a new file handle manager
func NewFileHandleManager() *FileHandleManager {
	return &FileHandleManager{
		handles: make(map[int]*FileHandle),
		nextID:  1,
		httpClient: &http.Client{
			Timeout: DefaultHTTPTimeout,
		},
	}
}

// validatePath sanitizes and canonicalizes a file path.
// It converts the path to an absolute form and removes redundant separators.
// Note: This does NOT restrict access to sensitive directories - use validateFilePath for that.
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

	return absPath, nil
}

// Close closes all open file handles managed by FileHandleManager.
// It implements the io.Closer interface.
func (m *FileHandleManager) Close() error {
	var firstErr error
	for id, handle := range m.handles {
		if handle != nil && handle.Closer != nil {
			if err := handle.Closer.Close(); err != nil && firstErr == nil {
				firstErr = fmt.Errorf("error closing handle %d (%s): %w", id, handle.FilePath, err)
			}
		}
	}
	m.handles = make(map[int]*FileHandle)
	return firstErr
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

	m.mu.Lock()
	id := m.nextID
	m.nextID++

	handle := &FileHandle{
		ID:        id,
		Reader:    file,
		Writer:    file,
		Closer:    file,
		FilePath:  validPath,
		IsNetwork: false,
	}
	m.mu.Unlock()

	// Initialize buffered reader for read mode
	if strings.ToLower(mode) == "r" || strings.ToLower(mode) == "read" {
		handle.BufReader = bufio.NewReader(file)
	}

	m.handles[id] = handle

	return id, nil
}

// LapLine reads a single line from a file handle (lap = drink/read)
func (m *FileHandleManager) LapLine(handleID int) (string, error) {
	m.mu.RLock()
	handle, ok := m.handles[handleID]
	m.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("invalid file handle: %d", handleID)
	}

	// Use the buffered reader if available, otherwise create one
	var reader *bufio.Reader
	if handle.BufReader != nil {
		reader = handle.BufReader
	} else {
		reader = bufio.NewReader(handle.Reader)
		handle.BufReader = reader
	}

	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	line = strings.TrimSuffix(line, "\n")
	return line, nil
}

// DevourFile reads entire contents of a file (devour = eat entirely)
func (m *FileHandleManager) DevourFile(handleID int) (string, error) {
	m.mu.RLock()
	handle, ok := m.handles[handleID]
	m.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("invalid file handle: %d", handleID)
	}

	// If a buffered reader exists (from LapLine calls), use it to preserve buffered data
	var reader io.Reader
	if handle.BufReader != nil {
		reader = handle.BufReader
	} else {
		reader = handle.Reader
	}

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return string(bytes), nil
}

// ScratchLine writes a line to a file handle (scratch = write)
func (m *FileHandleManager) ScratchLine(handleID int, content string) error {
	m.mu.RLock()
	handle, ok := m.handles[handleID]
	m.mu.RUnlock()

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
	m.mu.RLock()
	handle, ok := m.handles[handleID]
	m.mu.RUnlock()

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
	m.mu.Lock()
	handle, ok := m.handles[handleID]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("invalid file handle: %d", handleID)
	}
	delete(m.handles, handleID)
	m.mu.Unlock()

	if handle.Closer != nil {
		if err := handle.Closer.Close(); err != nil {
			return fmt.Errorf("error closing file: %w", err)
		}
	}

	return nil
}

// validateURL checks if a URL is safe to access (prevents SSRF attacks)
// Only allows http and https schemes
func validateURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Only allow http and https schemes
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("invalid URL scheme '%s': only http and https are allowed", parsedURL.Scheme)
	}

	// Ensure host is present
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	return nil
}

// FetchURL performs a GET request to a URL (fetch = get)
func (m *FileHandleManager) FetchURL(targetURL string) (string, error) {
	// Validate URL to prevent SSRF attacks
	if err := validateURL(targetURL); err != nil {
		return "", err
	}

	resp, err := http.Get(targetURL)
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
func (m *FileHandleManager) CoughUpData(targetURL string, contentType string, data string) (string, error) {
	// Validate URL to prevent SSRF attacks
	if err := validateURL(targetURL); err != nil {
		return "", err
	}

	resp, err := http.Post(targetURL, contentType, strings.NewReader(data))
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

// validateFilePath checks if a file path is safe for deletion
// It prevents deletion of critical system files and restricts to safe directories
func validateFilePath(filePath string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// Clean the path to remove .. and .
	cleanPath := filepath.Clean(absPath)

	// List of protected directories that should never be deleted from
	// Unix/Linux protected directories
	protectedDirs := []string{
		"/etc",
		"/bin",
		"/sbin",
		"/usr/bin",
		"/usr/sbin",
		"/boot",
		"/sys",
		"/proc",
		"/dev",
		"/lib",
		"/lib64",
		"/var/lib",
		"/usr/lib",
		"/root",
	}

	// Windows protected directories
	if filepath.Separator == '\\' {
		protectedDirs = []string{
			"C:\\Windows",
			"C:\\Program Files",
			"C:\\Program Files (x86)",
			"C:\\ProgramData",
			"C:\\System Volume Information",
			"C:\\$Recycle.Bin",
		}
	}

	// Check if the file is in a protected directory
	for _, protected := range protectedDirs {
		// Check if path is within protected directory or is the protected directory itself
		if strings.HasPrefix(cleanPath, protected+string(filepath.Separator)) || cleanPath == protected {
			return fmt.Errorf("cannot delete files in protected directory: %s", protected)
		}
	}

	// Prevent deletion of home directory itself
	homeDir, err := os.UserHomeDir()
	if err == nil && cleanPath == homeDir {
		return fmt.Errorf("cannot delete home directory")
	}

	// Prevent deletion of current working directory
	cwd, err := os.Getwd()
	if err == nil && cleanPath == cwd {
		return fmt.Errorf("cannot delete current working directory")
	}

	return nil
}

// SwatFile deletes a file (swat = hit/delete)
// It validates the file path to prevent deletion of critical system files
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
