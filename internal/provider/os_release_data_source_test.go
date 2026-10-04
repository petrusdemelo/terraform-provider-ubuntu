package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOSReleaseDataSource(t *testing.T) {
	config := testAccConfig(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config.provider("") + `data "ubuntu_os_release" "this" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ubuntu_os_release.this", "id", "ubuntu"),
					resource.TestMatchResourceAttr("data.ubuntu_os_release.this", "version_id", regexp.MustCompile(`^\d{2}\.\d{2}$`)),
					resource.TestCheckResourceAttrSet("data.ubuntu_os_release.this", "version_codename"),
				),
			},
		},
	})
}

func TestAccProviderRejectsWrongHostKey(t *testing.T) {
	config := testAccConfig(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config.provider(`host_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILsl7Oz7lf/jQkmlMr/RqgZ4oXiT0ppAwDmqsdCjrRv1"`) +
					`data "ubuntu_os_release" "this" {}`,
				ExpectError: regexp.MustCompile(`host key\s+mismatch`),
			},
		},
	})
}

func TestAccProviderRequiresHostKeyDecision(t *testing.T) {
	config := testAccConfig(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config.provider("insecure_ignore_host_key = false") + `data "ubuntu_os_release" "this" {}`,
				ExpectError: regexp.MustCompile(`host_key must be set`),
			},
		},
	})
}

func TestAccProviderRejectsUnknownTarget(t *testing.T) {
	config := testAccConfig(t)
	config.host = "${terraform_data.host.output}"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config.provider("") + `
resource "terraform_data" "host" {
  input = "127.0.0.1"
}

data "ubuntu_os_release" "this" {}
`,
				ExpectError: regexp.MustCompile(`Unknown provider configuration`),
			},
		},
	})
}
