package sdk

import (
	"net/http"
	"os"
	"testing"
)

var err error

// Print contents of object
// b, _ := json.Marshal(data)
// fmt.Println(string(b))

var basicClient = Client{HttpClient: &http.Client{}}

func mustEnv(t *testing.T, k string) string {
	t.Helper()

	if v, ok := os.LookupEnv(k); ok {
		return v
	}

	t.Fatalf("expected environment variable %q", k)
	return ""
}
