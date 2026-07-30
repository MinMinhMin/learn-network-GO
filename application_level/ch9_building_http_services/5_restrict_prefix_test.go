package ch9buildinghttpservices

// Protecting Sensitive Files

import (
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
)

// middleware that protecting sensitive file with specific prefix
func RestrictPrefix(prefix string, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			for _, p := range strings.Split(path.Clean(r.URL.Path), "/") {
				if strings.HasPrefix(p, prefix) {
					http.Error(w, "Not Found", http.StatusNotFound)
					return
				}
			}
			next.ServeHTTP(w, r)
		},
	)
}

func TestRestrictPrefix(t *testing.T) {

	// http.FileServer is a handler that read and return file from ./files

	// our client will send request path include with endpoint "/static/", but http.FileServer read file from true path, so we need to remove the prefix from client
	handler := http.StripPrefix("/static/", RestrictPrefix(".", http.FileServer(http.Dir("./files/"))))

	testCases := []struct {
		path string
		code int
	}{
		// these cases will send "/static/" + "file"
		{"http://test/static/sage.svg", http.StatusOK},
		{"http://test/static/.secret", http.StatusNotFound},
		{"http://test/static/.dir/secret", http.StatusNotFound},
	}

	for i, c := range testCases {
		r := httptest.NewRequest(http.MethodGet, c.path, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		actual := w.Result().StatusCode
		if actual != c.code {
			t.Errorf("%d: expected %d; actual %d", i, c.code, actual)
		}
	}

}
