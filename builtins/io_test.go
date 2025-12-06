// Copyright (c) 2025 BoyPlankton. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package builtins

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

// Network I/O tests

func TestFetchURL_Success(t *testing.T) {
	manager := NewFileHandleManager()
	expectedResponse := "Hello from test server"

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	// Test FetchURL
	response, err := manager.FetchURL(server.URL)
	if err != nil {
		t.Fatalf("FetchURL failed: %v", err)
	}

	if response != expectedResponse {
		t.Errorf("expected %q, got %q", expectedResponse, response)
	}
}

func TestFetchURL_HTTPError(t *testing.T) {
	manager := NewFileHandleManager()

	// Create a mock HTTP server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found")
	}))
	defer server.Close()

	// Test FetchURL with error response
	_, err := manager.FetchURL(server.URL)
	if err == nil {
		t.Error("expected error for HTTP 404, got nil")
	}

	if !strings.Contains(err.Error(), "404") {
		t.Errorf("expected error to contain '404', got %v", err)
	}
}

func TestFetchURL_NetworkError(t *testing.T) {
	manager := NewFileHandleManager()

	// Use an invalid URL to trigger a network error
	_, err := manager.FetchURL("http://invalid-host-that-does-not-exist-12345.com")
	if err == nil {
		t.Error("expected network error, got nil")
	}
}

func TestFetchURL_MultipleStatusCodes(t *testing.T) {
	manager := NewFileHandleManager()

	testCases := []struct {
		statusCode int
		expectErr  bool
	}{
		{http.StatusOK, false},
		{http.StatusBadRequest, true},
		{http.StatusUnauthorized, true},
		{http.StatusInternalServerError, true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Status_%d", tc.statusCode), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				fmt.Fprint(w, "test response")
			}))
			defer server.Close()

			_, err := manager.FetchURL(server.URL)
			if tc.expectErr && err == nil {
				t.Errorf("expected error for status %d, got nil", tc.statusCode)
			}
			if !tc.expectErr && err != nil {
				t.Errorf("expected no error for status %d, got %v", tc.statusCode, err)
			}
		})
	}
}

func TestCoughUpData_Success_StatusOK(t *testing.T) {
	manager := NewFileHandleManager()
	expectedResponse := "Data received"
	testData := "test payload"
	testContentType := "application/json"

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != testContentType {
			t.Errorf("expected Content-Type %s, got %s", testContentType, r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	// Test CoughUpData
	response, err := manager.CoughUpData(server.URL, testContentType, testData)
	if err != nil {
		t.Fatalf("CoughUpData failed: %v", err)
	}

	if response != expectedResponse {
		t.Errorf("expected %q, got %q", expectedResponse, response)
	}
}

func TestCoughUpData_Success_StatusCreated(t *testing.T) {
	manager := NewFileHandleManager()
	expectedResponse := "Resource created"
	testData := "create new resource"

	// Create a mock HTTP server that returns 201
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, expectedResponse)
	}))
	defer server.Close()

	// Test CoughUpData with 201 response
	response, err := manager.CoughUpData(server.URL, "text/plain", testData)
	if err != nil {
		t.Fatalf("CoughUpData failed: %v", err)
	}

	if response != expectedResponse {
		t.Errorf("expected %q, got %q", expectedResponse, response)
	}
}

func TestCoughUpData_HTTPError(t *testing.T) {
	manager := NewFileHandleManager()

	// Create a mock HTTP server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request")
	}))
	defer server.Close()

	// Test CoughUpData with error response
	_, err := manager.CoughUpData(server.URL, "text/plain", "test data")
	if err == nil {
		t.Error("expected error for HTTP 400, got nil")
	}

	if !strings.Contains(err.Error(), "400") {
		t.Errorf("expected error to contain '400', got %v", err)
	}
}

func TestCoughUpData_NetworkError(t *testing.T) {
	manager := NewFileHandleManager()

	// Use an invalid URL to trigger a network error
	_, err := manager.CoughUpData("http://invalid-host-that-does-not-exist-12345.com", "text/plain", "data")
	if err == nil {
		t.Error("expected network error, got nil")
	}
}

func TestCoughUpData_DifferentContentTypes(t *testing.T) {
	manager := NewFileHandleManager()

	contentTypes := []string{
		"application/json",
		"application/xml",
		"text/plain",
		"text/html",
	}

	for _, ct := range contentTypes {
		t.Run(ct, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Content-Type") != ct {
					t.Errorf("expected Content-Type %s, got %s", ct, r.Header.Get("Content-Type"))
				}
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "OK")
			}))
			defer server.Close()

			_, err := manager.CoughUpData(server.URL, ct, "test data")
			if err != nil {
				t.Errorf("CoughUpData failed for content type %s: %v", ct, err)
			}
		})
	}
}

func TestCoughUpData_MultipleStatusCodes(t *testing.T) {
	manager := NewFileHandleManager()

	testCases := []struct {
		statusCode int
		expectErr  bool
	}{
		{http.StatusOK, false},
		{http.StatusCreated, false},
		{http.StatusBadRequest, true},
		{http.StatusUnauthorized, true},
		{http.StatusNotFound, true},
		{http.StatusInternalServerError, true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Status_%d", tc.statusCode), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				fmt.Fprint(w, "test response")
			}))
			defer server.Close()

			_, err := manager.CoughUpData(server.URL, "text/plain", "test data")
			if tc.expectErr && err == nil {
				t.Errorf("expected error for status %d, got nil", tc.statusCode)
			}
			if !tc.expectErr && err != nil {
				t.Errorf("expected no error for status %d, got %v", tc.statusCode, err)
			}
		})
	}
}
