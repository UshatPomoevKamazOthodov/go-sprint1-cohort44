package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock storage
var mockUrlStorage = make(map[string]string)

// Test helpers
func testStorage() {
	mockUrlStorage = make(map[string]string)
}

func GetOriginalUrl(encodedId string) string {
	if url, exists := mockUrlStorage[encodedId]; exists {
		return url
	}
	return "not_found"
}

func mockInsertUrl(url string) (string, error) {
	encoded := "http://localhost:8080/asdasd"
	mockUrlStorage[encoded] = url
	return encoded, nil
}

// Reusable handler functions for testing
func handleGetRedirect(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	if path == "" || strings.Contains(path, "/") {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	originalUrl := GetOriginalUrl(path)
	if originalUrl == "not_found" {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func handlePostShorten(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}

	bodyBytes := make([]byte, r.ContentLength)
	n, err := r.Body.Read(bodyBytes)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalUrl := strings.TrimSpace(string(bodyBytes[:n]))
	if originalUrl == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	shortenedId, err := mockInsertUrl(originalUrl)
	if err != nil {
		http.Error(w, "Error creating shortened URL", http.StatusBadRequest)
		return
	}

	shortenedUrl := "http://" + r.Host + "/" + shortenedId
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenedUrl))
}

// Test GET handler
func TestHandleGetRedirect(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		mockData       map[string]string
		expectedStatus int
		expectedHeader string
		expectedBody   string
	}{
		{
			name:           "successful redirect",
			path:           "/aHR0cDovL2xvY2FsaG9zdDo4MDgwL2FzZGFzZA==",
			mockData:       map[string]string{"aHR0cDovL2xvY2FsaG9zdDo4MDgwL2FzZGFzZA==": "http://localhost:8080/asdasd"},
			expectedStatus: http.StatusTemporaryRedirect,
			expectedHeader: "http://localhost:8080/asdasd",
		},
		{
			name:           "url not found",
			path:           "/nonexistent",
			mockData:       map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "URL not found",
		},
		{
			name:           "empty path",
			path:           "/",
			mockData:       map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid URL ID",
		},
		{
			name:           "nested path",
			path:           "/some/nested/path",
			mockData:       map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid URL ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStorage()
			mockUrlStorage = tt.mockData

			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			handleGetRedirect(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedHeader != "" {
				assert.Equal(t, tt.expectedHeader, w.Header().Get("Location"))
			}

			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}

// Test POST handler
func TestHandlePostShorten(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		body           string
		contentType    string
		host           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful creation",
			path:           "/",
			body:           "http://localhost:8080/asdasd",
			contentType:    "text/plain",
			host:           "localhost:8080",
			expectedStatus: http.StatusCreated,
			expectedBody:   "http://localhost:8080/http://localhost:8080/asdasd",
		},
		{
			name:           "wrong content type",
			path:           "/",
			body:           "http://localhost:8080/asdasd",
			contentType:    "application/json",
			host:           "localhost:8080",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Content-Type must be text/plain",
		},
		{
			name:           "empty body",
			path:           "/",
			body:           "",
			contentType:    "text/plain",
			host:           "localhost:8080",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "URL cannot be empty",
		},
		{
			name:           "wrong path",
			path:           "/wrong/path",
			body:           "http://localhost:8080/asdasd",
			contentType:    "text/plain",
			host:           "localhost:8080",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid path",
		},
		{
			name:           "whitespace only body",
			path:           "/",
			body:           "   \n\t   ",
			contentType:    "text/plain",
			host:           "localhost:8080",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "URL cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStorage()

			req := httptest.NewRequest("POST", tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			req.Host = tt.host
			w := httptest.NewRecorder()

			handlePostShorten(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
				assert.Equal(t, tt.expectedBody, w.Body.String())
			} else {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}
