package ssh

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	cryptossh "golang.org/x/crypto/ssh"
)

const ed25519HostKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILsl7Oz7lf/jQkmlMr/RqgZ4oXiT0ppAwDmqsdCjrRv1"

func TestHostKeyCallbackPinsAlgorithm(t *testing.T) {
	_, algorithms, err := hostKeyCallback(Config{HostKey: ed25519HostKey})
	if err != nil {
		t.Fatalf("hostKeyCallback() error = %v", err)
	}

	if len(algorithms) != 1 || algorithms[0] != cryptossh.KeyAlgoED25519 {
		t.Errorf("algorithms = %v, want [%s]", algorithms, cryptossh.KeyAlgoED25519)
	}
}

func TestHostKeyAlgorithmsAcceptsRSASignatureVariants(t *testing.T) {
	got := hostKeyAlgorithms(cryptossh.KeyAlgoRSA)

	for _, want := range []string{cryptossh.KeyAlgoRSASHA256, cryptossh.KeyAlgoRSASHA512, cryptossh.KeyAlgoRSA} {
		if !slices.Contains(got, want) {
			t.Errorf("hostKeyAlgorithms(ssh-rsa) = %v, missing %s", got, want)
		}
	}
}

func TestHostKeyCallbackRequiresExplicitOptOut(t *testing.T) {
	_, _, err := hostKeyCallback(Config{})
	if err == nil {
		t.Fatal("hostKeyCallback() with no host key and no opt-out should error")
	}

	if !strings.Contains(err.Error(), "insecure_ignore_host_key") {
		t.Errorf("error should name the opt-out, got: %v", err)
	}
}

func TestHostKeyCallbackAllowsExplicitOptOut(t *testing.T) {
	callback, algorithms, err := hostKeyCallback(Config{InsecureIgnoreHostKey: true})
	if err != nil {
		t.Fatalf("hostKeyCallback() error = %v", err)
	}

	if callback == nil {
		t.Error("callback should be set when host key checking is disabled")
	}

	if algorithms != nil {
		t.Errorf("algorithms should be unconstrained when not pinning, got %v", algorithms)
	}
}

func TestHostKeyCallbackRejectsGarbage(t *testing.T) {
	if _, _, err := hostKeyCallback(Config{HostKey: "not-a-key"}); err == nil {
		t.Fatal("hostKeyCallback() should reject an unparseable host key")
	}
}

func TestAuthMethodRequiresACredential(t *testing.T) {
	if _, err := authMethod(Config{}); err == nil {
		t.Fatal("authMethod() should error when neither private_key nor password is set")
	}
}

func TestDialHonoursCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := Dial(ctx, Config{Host: "192.0.2.1", Port: 22, Password: "x", InsecureIgnoreHostKey: true})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Dial() error = %v, want context.Canceled", err)
	}
}
