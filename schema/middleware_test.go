package schema

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
)

func TestRequestValidator(t *testing.T) {
	router := openapi3filter.NewRouter().WithSwaggerFromFile("./testdata/openapi.yaml")

	mux := http.NewServeMux()
	mux.Handle("/", RequestValidator(router, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))

	testCases := []struct{
		name string
		wantStatus int
		req *http.Request
	}{
		{
			name: "get ok",
			wantStatus: http.StatusOK,
			req: httptest.NewRequest(http.MethodGet, "/test", nil),
		},
		{
			name: "no route",
			wantStatus: http.StatusNotFound,
			req: httptest.NewRequest(http.MethodGet, "/not-found", nil),
		},
		{
			name: "invalid body",
			wantStatus: http.StatusBadRequest,
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"foo": "string", "bar": 3}`))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
		},
		{
			name: "post ok",
			wantStatus: http.StatusOK,
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"foo": "string", "bar": 2}`))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, tc.req)
			if rec.Code != tc.wantStatus {
				t.Fatalf("invalid status: want=%v, got=%v", rec.Code, tc.wantStatus)
			}
		})
	}

	t.Run("validationErrHook", func(t *testing.T) {
		mux := http.NewServeMux()
		hook := func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(http.StatusInternalServerError)
		}
		mux.Handle("/", RequestValidator(router, hook)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"foo": "string", "bar": 3}`))
		req.Header.Set("Content-Type", "application/json")
		mux.ServeHTTP(rec, req)
		wantStatus := http.StatusInternalServerError
		if rec.Code != wantStatus {
			t.Fatalf("invalid status: want=%v, got=%v", rec.Code, wantStatus)
		}
	})
}
