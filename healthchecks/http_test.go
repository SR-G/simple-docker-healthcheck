package healthchecks

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type healthCheckFunc func() (bool, error)

func TestHTTPHealthChecks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/exact":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		case "/range":
			w.WriteHeader(http.StatusNoContent)
			_, _ = w.Write([]byte("range"))
		case "/text":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("hello world"))
		case "/json":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":"expected"}`))
		case "/invalid-json":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("not-json"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tests := []struct {
		name    string
		run     healthCheckFunc
		want    bool
		wantErr string
	}{
		{
			name: "exact status success",
			run: func() (bool, error) {
				return (&HTTPCodeHealthCheckExactStatusCode{URL: server.URL + "/exact", ExpectedStatusCode: 200}).Execute()
			},
			want: true,
		},
		{
			name: "exact status failure",
			run: func() (bool, error) {
				return (&HTTPCodeHealthCheckExactStatusCode{URL: server.URL + "/exact", ExpectedStatusCode: 204}).Execute()
			},
			want:    false,
			wantErr: "expected status code",
		},
		{
			name: "range status success",
			run: func() (bool, error) {
				return (&HTTPCodeHealthCheckRangeStatusCode{URL: server.URL + "/range", MinExpectedStatusCode: 200, MaxExpectedStatusCode: 299}).Execute()
			},
			want: true,
		},
		{
			name: "text content success",
			run: func() (bool, error) {
				return (&HTTPHealthCheckText{URL: server.URL + "/text", ExpectedText: "hello"}).Execute()
			},
			want: true,
		},
		{
			name: "jsonpath success",
			run: func() (bool, error) {
				return (&HTTPHealthCheckJSONPath{URL: server.URL + "/json", JSONPath: "$.data", ExpectedValue: "expected"}).Execute()
			},
			want: true,
		},
		{
			name: "jsonpath invalid json",
			run: func() (bool, error) {
				return (&HTTPHealthCheckJSONPath{URL: server.URL + "/invalid-json", JSONPath: "$.data", ExpectedValue: "expected"}).Execute()
			},
			want:    false,
			wantErr: "response is not a valid JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.run()
			if got != tt.want {
				t.Fatalf("want=%v got=%v", tt.want, got)
			}
			if tt.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
			}
		})
	}
}
