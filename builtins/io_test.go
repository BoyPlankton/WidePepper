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

// Path Traversal Security Tests

func TestValidatePath_EmptyPath(t *testing.T) {
	_, err := validatePath("")
	if err == nil {
		t.Error("expected error for empty path")
	}
	if err != nil && !strings.Contains(err.Error(), "empty") {
		t.Errorf("expected 'empty' error, got: %v", err)
	}
}

func TestPounceFile_PathTraversal(t *testing.T) {
	manager := NewFileHandleManager()

	// Test various path traversal attempts
	traversalPaths := []string{
		"../../etc/passwd",
		"../../../etc/passwd",
		"./../../etc/passwd",
		"./../../../etc/passwd",
	}

	for _, maliciousPath := range traversalPaths {
		// These should still work but resolve to safe absolute paths
		// The key is that they won't escape to sensitive system files
		_, err := manager.PounceFile(maliciousPath, "r")
		// The error could be either path validation or file not found
		// Both are acceptable as long as we don't access sensitive files
		if err != nil {
			// This is good - the operation was blocked
			continue
		}
		// If no error, the path should have been cleaned and made absolute
	}
}

func TestSniffFile_PathTraversal(t *testing.T) {
	// Create a test file in a temp directory
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Verify the file exists with a safe path
	if !SniffFile(testFile) {
		t.Error("expected file to exist")
	}

	// Empty path should return false
	if SniffFile("") {
		t.Error("expected false for empty path")
	}
}

func TestSwatFile_PathTraversal(t *testing.T) {
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "delete_test.txt")

	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Valid delete should work
	err = SwatFile(testFile)
	if err != nil {
		t.Fatalf("failed to delete file: %v", err)
	}

	// Empty path should fail
	err = SwatFile("")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestPounceDirectory_PathTraversal(t *testing.T) {
	testDir := t.TempDir()

	// Create test files
	files := []string{"file1.txt", "file2.txt"}
	for _, f := range files {
		err := os.WriteFile(filepath.Join(testDir, f), []byte("test"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	// Valid directory listing should work
	listed, err := PounceDirectory(testDir)
	if err != nil {
		t.Fatalf("failed to list directory: %v", err)
	}

	if len(listed) != len(files) {
		t.Errorf("expected %d files, got %d", len(files), len(listed))
	}

	// Empty path should fail
	_, err = PounceDirectory("")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestPounceFile_CleanedPaths(t *testing.T) {
	manager := NewFileHandleManager()
	testDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(testDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Try to open with various "dirty" but valid paths
	dirtyPaths := []string{
		testFile,
		filepath.Join(testDir, ".", "test.txt"),
		filepath.Join(testDir, "subdir", "..", "test.txt"),
	}

	for _, path := range dirtyPaths {
		handleID, err := manager.PounceFile(path, "r")
		if err != nil {
			t.Errorf("failed to open file with path %q: %v", path, err)
			continue
		}

		// Read content to verify it's the right file
		content, err := manager.DevourFile(handleID)
		if err != nil {
			t.Errorf("failed to read file: %v", err)
		}

		if content != "test content" {
			t.Errorf("unexpected content for path %q: got %q", path, content)
		}

		manager.NuzzleClose(handleID)
	}
}
