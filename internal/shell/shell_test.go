package shell

import (
	"os/exec"
	"testing"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", "''"},
		{"plain", "'plain'"},
		{"two words", "'two words'"},
		{"it's", `'it'\''s'`},
		{"$(rm -rf /)", "'$(rm -rf /)'"},
	}

	for _, tt := range tests {
		if got := Quote(tt.in); got != tt.want {
			t.Errorf("Quote(%q) = %s, want %s", tt.in, got, tt.want)
		}
	}
}

func TestQuoteRoundTripsThroughSh(t *testing.T) {
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("sh not found")
	}

	for _, value := range []string{"", "it's", `back\slash`, "$HOME `id`", "line\nbreak"} {
		out, err := exec.Command(sh, "-c", "printf '%s' "+Quote(value)).Output()
		if err != nil {
			t.Fatalf("sh: %v", err)
		}

		if string(out) != value {
			t.Errorf("round trip of %q gave %q", value, out)
		}
	}
}

func TestWrapRunsTheScriptUnchanged(t *testing.T) {
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("sh not found")
	}

	script := "x='a b'; printf '%s|%s' \"$x\" \"$(echo it\\'s)\""

	out, err := exec.Command(sh, "-c", Wrap(script, false)).Output()
	if err != nil {
		t.Fatalf("sh: %v", err)
	}

	if got, want := string(out), "a b|it's"; got != want {
		t.Errorf("Wrap() ran to %q, want %q", got, want)
	}
}
