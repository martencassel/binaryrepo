package client

import "testing"

func TestPingable(t *testing.T) {
	testcases := map[string]struct {
		registry Registry
		expect   bool
	}{
		"Docker": {
			registry: Registry{URL: "https://index.docker.io"},
			expect:   true,
		},
		"GCR_global": {
			registry: Registry{URL: "https://gcr.io"},
			expect:   false,
		},
	}
	for label, testcatestcases := range testcases {
		if testcatestcases.registry.Pingable() != testcatestcases.expect {
			t.Fatalf("%s: expected %v, got %v", label, testcatestcases.expect, testcatestcases.registry.Pingable())
		}
	}
}
