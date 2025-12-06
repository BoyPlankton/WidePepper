// Copyright (c) 2025 BoyPlankton. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package builtins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileHandleManager_OpenAndWriteFile(t *testing.T) {
	manager := NewFileHandleManager()
	testFile := filepath.Join(t.TempDir(), "test.txt")

	handleID, err := manager.PounceFile(testFile, "w")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	err = manager.ScratchLine(handleID, "Hello, WidePepper!")
	if err != nil {
		t.Fatalf("failed to write line: %v", err)
	}

	err = manager.NuzzleClose(handleID)
	if err != nil {
		t.Fatalf("failed to close file: %v", err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	expected := "Hello, WidePepper!\n"
	if string(content) != expected {
		t.Errorf("expected %q, got %q", expected, string(content))
	}
}

func TestFileHandleManager_ReadFile(t *testing.T) {
	manager := NewFileHandleManager()
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "read_test.txt")

	err := os.WriteFile(testFile, []byte("Line 1\nLine 2\nLine 3\n"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	handleID, err := manager.PounceFile(testFile, "r")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer manager.NuzzleClose(handleID)

	line, err := manager.LapLine(handleID)
	if err != nil {
		t.Fatalf("failed to read line: %v", err)
	}

	if line != "Line 1" {
		t.Errorf("expected 'Line 1', got %q", line)
	}
}

func TestFileHandleManager_ReadAll(t *testing.T) {
	manager := NewFileHandleManager()
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "readall_test.txt")

	expectedContent := "All content here\nMultiple lines\nTest data"
	err := os.WriteFile(testFile, []byte(expectedContent), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	handleID, err := manager.PounceFile(testFile, "r")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer manager.NuzzleClose(handleID)

	content, err := manager.DevourFile(handleID)
	if err != nil {
		t.Fatalf("failed to read all: %v", err)
	}

	if content != expectedContent {
		t.Errorf("expected %q, got %q", expectedContent, content)
	}
}

func TestFileHandleManager_AppendMode(t *testing.T) {
	manager := NewFileHandleManager()
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "append_test.txt")

	err := os.WriteFile(testFile, []byte("Initial content\n"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	handleID, err := manager.PounceFile(testFile, "a")
	if err != nil {
		t.Fatalf("failed to open file in append mode: %v", err)
	}

	err = manager.ScratchLine(handleID, "Appended line")
	if err != nil {
		t.Fatalf("failed to write append: %v", err)
	}

	err = manager.NuzzleClose(handleID)
	if err != nil {
		t.Fatalf("failed to close file: %v", err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	if !strings.Contains(string(content), "Initial content") {
		t.Error("initial content was lost")
	}

	if !strings.Contains(string(content), "Appended line") {
		t.Error("appended line not found")
	}
}

func TestFileHandleManager_InvalidHandle(t *testing.T) {
	manager := NewFileHandleManager()

	_, err := manager.LapLine(999)
	if err == nil {
		t.Error("expected error for invalid handle")
	}

	err = manager.ScratchLine(999, "test")
	if err == nil {
		t.Error("expected error for invalid handle")
	}
}

func TestFileHandleManager_InvalidMode(t *testing.T) {
	manager := NewFileHandleManager()
	testFile := filepath.Join(t.TempDir(), "test.txt")

	_, err := manager.PounceFile(testFile, "invalid")
	if err == nil {
		t.Error("expected error for invalid mode")
	}
}

func TestSniffFile(t *testing.T) {
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "exists_test.txt")

	if SniffFile(testFile) {
		t.Error("expected file to not exist")
	}

	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if !SniffFile(testFile) {
		t.Error("expected file to exist")
	}
}

func TestSwatFile(t *testing.T) {
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "delete_test.txt")

	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = SwatFile(testFile)
	if err != nil {
		t.Fatalf("failed to delete file: %v", err)
	}

	if SniffFile(testFile) {
		t.Error("expected file to be deleted")
	}
}

func TestPounceDirectory(t *testing.T) {
	testDir := t.TempDir()

	files := []string{"file1.txt", "file2.txt", "file3.wp"}
	for _, f := range files {
		err := os.WriteFile(filepath.Join(testDir, f), []byte("test"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	listed, err := PounceDirectory(testDir)
	if err != nil {
		t.Fatalf("failed to list files: %v", err)
	}

	if len(listed) != len(files) {
		t.Errorf("expected %d files, got %d", len(files), len(listed))
	}

	for _, f := range files {
		found := false
		for _, l := range listed {
			if l == f {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected file %s not found in listing", f)
		}
	}
}

func TestListFiles_NonExistentDirectory(t *testing.T) {
	_, err := PounceDirectory("/nonexistent/directory")
	if err == nil {
		t.Error("expected error for non-existent directory")
	}
}

func TestFileHandleManager_WriteString(t *testing.T) {
	manager := NewFileHandleManager()
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "writestring_test.txt")

	handleID, err := manager.PounceFile(testFile, "w")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	err = manager.ScratchString(handleID, "No newline")
	if err != nil {
		t.Fatalf("failed to write string: %v", err)
	}

	err = manager.ScratchString(handleID, " here")
	if err != nil {
		t.Fatalf("failed to write string: %v", err)
	}

	err = manager.NuzzleClose(handleID)
	if err != nil {
		t.Fatalf("failed to close file: %v", err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	expected := "No newline here"
	if string(content) != expected {
		t.Errorf("expected %q, got %q", expected, string(content))
	}
}

// TestValidateURL_ValidSchemes tests that http and https URLs are accepted
func TestValidateURL_ValidSchemes(t *testing.T) {
	validURLs := []string{
		"http://example.com",
		"https://example.com",
		"HTTP://example.com",
		"HTTPS://example.com",
		"http://example.com:8080/path",
		"https://api.example.com/v1/endpoint?param=value",
	}

	for _, testURL := range validURLs {
		err := validateURL(testURL)
		if err != nil {
			t.Errorf("expected URL %q to be valid, got error: %v", testURL, err)
		}
	}
}

// TestValidateURL_InvalidSchemes tests that non-http/https schemes are rejected
func TestValidateURL_InvalidSchemes(t *testing.T) {
	invalidURLs := []string{
		"file:///etc/passwd",
		"ftp://example.com",
		"javascript:alert(1)",
		"data:text/html,<script>alert(1)</script>",
		"gopher://example.com",
		"telnet://example.com",
		"ldap://example.com",
	}

	for _, testURL := range invalidURLs {
		err := validateURL(testURL)
		if err == nil {
			t.Errorf("expected URL %q to be invalid (SSRF risk), but it was accepted", testURL)
		}
	}
}

// TestValidateURL_MalformedURLs tests that malformed URLs are rejected
func TestValidateURL_MalformedURLs(t *testing.T) {
	malformedURLs := []string{
		"not a url",
		"://missing-scheme",
		"http://",
		"https://",
	}

	for _, testURL := range malformedURLs {
		err := validateURL(testURL)
		if err == nil {
			t.Errorf("expected malformed URL %q to be rejected, but it was accepted", testURL)
		}
	}
}

// TestFetchURL_InvalidScheme tests that FetchURL rejects invalid URL schemes
func TestFetchURL_InvalidScheme(t *testing.T) {
	manager := NewFileHandleManager()

	// Try to fetch a file:// URL (SSRF attack vector)
	_, err := manager.FetchURL("file:///etc/passwd")
	if err == nil {
		t.Error("expected FetchURL to reject file:// scheme")
	}

	// Try to fetch a javascript: URL
	_, err = manager.FetchURL("javascript:alert(1)")
	if err == nil {
		t.Error("expected FetchURL to reject javascript: scheme")
	}
}

// TestCoughUpData_InvalidScheme tests that CoughUpData rejects invalid URL schemes
func TestCoughUpData_InvalidScheme(t *testing.T) {
	manager := NewFileHandleManager()

	// Try to POST to a file:// URL (SSRF attack vector)
	_, err := manager.CoughUpData("file:///tmp/test", "text/plain", "data")
	if err == nil {
		t.Error("expected CoughUpData to reject file:// scheme")
	}

	// Try to POST to a ftp: URL
	_, err = manager.CoughUpData("ftp://example.com", "text/plain", "data")
	if err == nil {
		t.Error("expected CoughUpData to reject ftp: scheme")
	}
}
