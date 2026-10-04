# terraform-provider-ubuntu

A Terraform provider that manages Ubuntu hosts over SSH.

It covers the work that usually ends up in a cloud-init `user_data` script or a
`remote-exec` provisioner — packages, files, systemd units, firewall rules,
hostnames — and turns it into resources with a plan, drift detection and a
destroy path.

The provider is agentless. It opens one SSH connection and runs plain shell
commands; nothing is copied to or installed on the host.

## Motivation

Host configuration in Terraform usually ends up in one of two places, and
neither is managed state:

- **`user_data`** runs once, at first boot. Changing it either does nothing on
  a running host or replaces the instance.
- **`remote-exec` provisioners** run once, at create time. Terraform does not
  know what they changed, cannot show a diff for them, and cannot undo them.

Configuration management tools solve this, but add a second tool, a second
state model and often an agent on every host. For a handful of hosts already
described in Terraform, that is more machinery than the problem needs.

HashiCorp's [`hashicorp/terraform-provider-ubuntu`](https://github.com/hashicorp/terraform-provider-ubuntu)
has the right resource model, and this project follows its API. That provider
describes itself as alpha and not officially supported, and it works by
pushing an executor binary to every host. This project takes the API in a
smaller, independent form instead: SSH and the shell only, no binary on the
host, and a codebase that can be read in an afternoon.

## Usage

### Connect to a host and install a package

```hcl
terraform {
  required_providers {
    ubuntu = {
      source  = "petrusdemelo/ubuntu"
      version = "~> 0.1"
    }
  }
}

provider "ubuntu" {
  ssh {
    user        = "terraform"
    private_key = var.ssh_private_key_pem
    host_key    = var.host_public_key
  }

  default_target {
    target = var.host_address
    port   = 22
  }
}

resource "ubuntu_package" "nginx" {
  name = "nginx"
}
```

`host_key` pins the server's public key, in `authorized_keys` format. Skipping
verification requires `insecure_ignore_host_key = true`.

### Configure nginx as a reverse proxy

```hcl
resource "ubuntu_file" "nginx_route" {
  path    = "/etc/nginx/conf.d/internal-api.conf"
  content = templatefile("${path.module}/internal-api.conf.tftpl", {
    upstream = "http://10.0.2.15:8080"
  })
  owner = "root"
  group = "root"
  mode  = "0644"

  depends_on = [ubuntu_package.nginx]
}

resource "ubuntu_systemd_unit" "nginx" {
  name    = "nginx.service"
  enabled = true
  state   = "started"

  reload_triggers = [ubuntu_file.nginx_route.content]
}

resource "ubuntu_ufw_rule" "http" {
  port     = "80"
  protocol = "tcp"
}
```

Changing the template rewrites the file and reloads nginx in the same apply.
Editing the file by hand on the host shows up as drift in the next plan.

## Status

Pre-1.0. The resources above land during the
[v0.1.0 milestone](https://github.com/petrusdemelo/terraform-provider-ubuntu/milestone/1);
the resource API may still change between minor releases, and every breaking
change is listed in the [changelog](CHANGELOG.md).

The provider was extracted from production use at Nuvian, where it manages
its Ubuntu hosts.

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
