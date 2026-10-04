// Package ssh runs shell commands on a remote host over a single SSH connection.
package ssh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	cryptossh "golang.org/x/crypto/ssh"
)

const handshakeTimeout = 30 * time.Second

type Config struct {
	Host                  string
	Port                  int64
	User                  string
	PrivateKey            string
	Password              string
	HostKey               string
	InsecureIgnoreHostKey bool
}

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type Client struct {
	conn *cryptossh.Client
}

func Dial(ctx context.Context, cfg Config) (*Client, error) {
	auth, err := authMethod(cfg)
	if err != nil {
		return nil, err
	}

	callback, algorithms, err := hostKeyCallback(cfg)
	if err != nil {
		return nil, err
	}

	address := net.JoinHostPort(cfg.Host, strconv.FormatInt(cfg.Port, 10))

	var dialer net.Dialer

	tcp, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", address, err)
	}

	conn, chans, reqs, err := cryptossh.NewClientConn(tcp, address, &cryptossh.ClientConfig{
		User:              cfg.User,
		Auth:              []cryptossh.AuthMethod{auth},
		HostKeyCallback:   callback,
		HostKeyAlgorithms: algorithms,
		Timeout:           handshakeTimeout,
	})
	if err != nil {
		_ = tcp.Close()

		return nil, fmt.Errorf("ssh handshake with %s: %w", address, err)
	}

	return &Client{conn: cryptossh.NewClient(conn, chans, reqs)}, nil
}

// Run feeds cmd to sh on its standard input, so the command line never passes
// through the login shell's own parsing. A non-zero exit status is reported in
// Result.ExitCode, not as an error.
func (c *Client) Run(ctx context.Context, cmd string) (Result, error) {
	session, err := c.conn.NewSession()
	if err != nil {
		return Result{}, fmt.Errorf("open ssh session: %w", err)
	}
	// Close reports io.EOF once Wait has already torn the session down.
	defer func() { _ = session.Close() }()

	var stdout, stderr bytes.Buffer
	session.Stdin = strings.NewReader(cmd)
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Start("sh"); err != nil {
		return Result{}, fmt.Errorf("start remote command: %w", err)
	}

	waited := make(chan error, 1)
	go func() { waited <- session.Wait() }()

	select {
	case <-ctx.Done():
		_ = session.Signal(cryptossh.SIGKILL)

		return Result{}, ctx.Err()
	case err := <-waited:
		result := Result{Stdout: stdout.String(), Stderr: stderr.String()}

		var exitErr *cryptossh.ExitError
		if errors.As(err, &exitErr) {
			result.ExitCode = exitErr.ExitStatus()

			return result, nil
		}

		if err != nil {
			return Result{}, fmt.Errorf("remote command: %w", err)
		}

		return result, nil
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func authMethod(cfg Config) (cryptossh.AuthMethod, error) {
	switch {
	case cfg.PrivateKey != "":
		signer, err := cryptossh.ParsePrivateKey([]byte(cfg.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("parse private_key: %w", err)
		}

		return cryptossh.PublicKeys(signer), nil
	case cfg.Password != "":
		return cryptossh.Password(cfg.Password), nil
	default:
		return nil, errors.New("either private_key or password must be set")
	}
}

// The host key algorithm is pinned to the configured key's type. Otherwise a
// server holding several host keys may present another type, and the handshake
// fails as a mismatch instead of succeeding.
func hostKeyCallback(cfg Config) (cryptossh.HostKeyCallback, []string, error) {
	if cfg.HostKey != "" {
		key, _, _, _, err := cryptossh.ParseAuthorizedKey([]byte(cfg.HostKey))
		if err != nil {
			return nil, nil, fmt.Errorf("parse host_key: %w", err)
		}

		return cryptossh.FixedHostKey(key), hostKeyAlgorithms(key.Type()), nil
	}

	if !cfg.InsecureIgnoreHostKey {
		return nil, nil, errors.New("host_key must be set, or insecure_ignore_host_key explicitly enabled")
	}

	return cryptossh.InsecureIgnoreHostKey(), nil, nil
}

// An RSA host key is announced as ssh-rsa, but current servers sign with the
// SHA-2 variants.
func hostKeyAlgorithms(keyType string) []string {
	if keyType == cryptossh.KeyAlgoRSA {
		return []string{cryptossh.KeyAlgoRSASHA256, cryptossh.KeyAlgoRSASHA512, cryptossh.KeyAlgoRSA}
	}

	return []string{keyType}
}
