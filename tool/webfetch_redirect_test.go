package tool

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebFetchStopsCrossHostRedirect(t *testing.T) {
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Error("cross-host redirect was followed")
	}))
	defer other.Close()
	src := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/same" {
			http.Redirect(w, r, other.URL+"/x", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/same", http.StatusFound)
	}))
	defer src.Close()

	resp, err := newFetchClient(WebFetchConfig{}).Get(src.URL + "/start")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound || resp.Header.Get("Location") != other.URL+"/x" {
		t.Fatalf("got %d Location=%q, want the unfollowed cross-host 302", resp.StatusCode, resp.Header.Get("Location"))
	}
}
