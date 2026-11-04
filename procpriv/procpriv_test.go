package procpriv

import "testing"

func TestWithin(t *testing.T) {
	for _, c := range []struct {
		root, path string
		want       bool
	}{
		{"/app/data", "/app/data/tasks/1", true},
		{"/app/data", "/app/data", true},
		{"/app/data", "/app/database", false},
		{"/app/data", "/app/data/../jwt.key", false},
		{"/app/data", "/", false},
	} {
		if got := within(c.root, c.path); got != c.want {
			t.Errorf("within(%q, %q) = %v, want %v", c.root, c.path, got, c.want)
		}
	}
}

func TestWithHome(t *testing.T) {
	got := withHome([]string{"A=1", "HOME=/root"}, "/home/agent")
	if len(got) != 2 || got[0] != "A=1" || got[1] != "HOME=/home/agent" {
		t.Fatalf("withHome = %v", got)
	}
}
