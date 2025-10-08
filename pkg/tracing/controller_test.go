package tracing

import "testing"

func TestIsTruthy(t *testing.T) {
	cases := map[string]bool{
		"true":  true,
		"TRUE":  true,
		"1":     true,
		"yes":   true,
		"on":    true,
		"false": false,
		"0":     false,
		"":      false,
		"no":    false,
	}

	for input, expected := range cases {
		if actual := isTruthy(input); actual != expected {
			t.Fatalf("isTruthy(%q) = %t, expected %t", input, actual, expected)
		}
	}
}

func TestIsTruthyWithDefault(t *testing.T) {
	if !isTruthyWithDefault("", true) {
		t.Fatalf("expected default true when value empty")
	}
	if isTruthyWithDefault("", false) {
		t.Fatalf("expected default false when value empty")
	}
	if !isTruthyWithDefault("true", false) {
		t.Fatalf("expected override true when value explicit")
	}
	if isTruthyWithDefault("no", true) {
		t.Fatalf("expected override false when value explicit")
	}
}
