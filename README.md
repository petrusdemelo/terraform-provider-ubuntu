# terraform-provider-ubuntu

A Terraform provider that manages Ubuntu hosts over SSH.

It covers the work that usually ends up in a cloud-init `user_data` script or a
`remote-exec` provisioner — packages, files, systemd units, firewall rules,
hostnames — and turns it into resources with a plan, drift detection and a
destroy path.

The provider is agentless. It opens one SSH connection and runs plain shell
commands; nothing is copied to or installed on the host.

## Status

Pre-1.0. The resource API may still change between minor releases; every
breaking change is listed in the [changelog](CHANGELOG.md).

The provider was extracted from production use at Nuvian, where it manages
its Ubuntu hosts.

## Relationship to hashicorp/terraform-provider-ubuntu

HashiCorp publishes a provider with the same name and a much larger scope,
built around an executor it ships to the host. This project follows its
resource naming where the two overlap, but stays deliberately small: SSH and
the shell only, so that the whole provider can be read in an afternoon.

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
