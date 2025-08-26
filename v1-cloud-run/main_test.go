package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		wantBody   string
		wantStatus int
	}{
		{
			name:       "デフォルトターゲット",
			target:     "",
			wantBody:   "Hello, World!",
			wantStatus: http.StatusOK,
		},
		{
			name:       "カスタムターゲット",
			target:     "Cloud Run",
			wantBody:   "Hello, Cloud Run!",
			wantStatus: http.StatusOK,
		},
		{
			name:       "日本語ターゲット",
			target:     "世界",
			wantBody:   "Hello, 世界!",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.target != "" {
				os.Setenv("TARGET", tt.target)
				defer os.Unsetenv("TARGET")
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			handler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("handler() status = %v, want %v", rec.Code, tt.wantStatus)
			}

			got := rec.Body.String()
			if !strings.Contains(got, tt.wantBody) {
				t.Errorf("handler() body = %v, want %v", got, tt.wantBody)
			}
		})
	}
}

func TestHandlerMethods(t *testing.T) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)
			rec := httptest.NewRecorder()

			handler(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("handler() with method %s: status = %v, want %v", method, rec.Code, http.StatusOK)
			}
		})
	}
}

func TestHealthzHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	healthzHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("healthzHandler() status = %v, want %v", rec.Code, http.StatusOK)
	}

	got := rec.Body.String()
	want := "OK"
	if got != want {
		t.Errorf("healthzHandler() body = %v, want %v", got, want)
	}
}

func TestReadyzHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	readyzHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("readyzHandler() status = %v, want %v", rec.Code, http.StatusOK)
	}

	got := rec.Body.String()
	want := "READY"
	if got != want {
		t.Errorf("readyzHandler() body = %v, want %v", got, want)
	}
}

func BenchmarkHandler(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler(rec, req)
	}
}
