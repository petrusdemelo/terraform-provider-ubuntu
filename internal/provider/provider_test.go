package provider_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/petrusdemelo/terraform-provider-ubuntu/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"ubuntu": providerserver.NewProtocol6WithError(provider.New("acceptance")()),
}

func TestProviderSchema(t *testing.T) {
	var resp fwprovider.SchemaResponse

	provider.New("test")().Schema(t.Context(), fwprovider.SchemaRequest{}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(t.Context()); diags.HasError() {
		t.Fatalf("schema is invalid: %v", diags)
	}
}

type testConfig struct {
	host string
	port int64
	user string
	key  string
}

// provider renders a provider block for the acceptance host. extra is
// appended inside the ssh block, so a test can add host_key or override
// insecure_ignore_host_key.
func (c testConfig) provider(extra string) string {
	if extra == "" {
		extra = "insecure_ignore_host_key = true"
	}

	return fmt.Sprintf(`
provider "ubuntu" {
  ssh {
    user        = %q
    private_key = %q
    %s
  }

  default_target {
    target = %q
    port   = %d
  }
}
`, c.user, c.key, extra, c.host, c.port)
}

func testAccConfig(t *testing.T) testConfig {
	t.Helper()

	host := os.Getenv("UBUNTU_TEST_HOST")
	user := os.Getenv("UBUNTU_TEST_USER")
	keyPath := os.Getenv("UBUNTU_TEST_PRIVATE_KEY_PATH")

	if os.Getenv("TF_ACC") == "" || host == "" || user == "" || keyPath == "" {
		t.Skip("set TF_ACC, UBUNTU_TEST_HOST, UBUNTU_TEST_USER and UBUNTU_TEST_PRIVATE_KEY_PATH, or run make testenv-up testacc")
	}

	port, err := strconv.ParseInt(os.Getenv("UBUNTU_TEST_PORT"), 10, 64)
	if err != nil {
		t.Fatalf("UBUNTU_TEST_PORT: %v", err)
	}

	key, err := os.ReadFile(keyPath)
	if err != nil {
		t.Fatalf("read UBUNTU_TEST_PRIVATE_KEY_PATH: %v", err)
	}

	return testConfig{host: host, port: port, user: user, key: string(key)}
}
