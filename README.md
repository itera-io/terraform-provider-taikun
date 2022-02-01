# Terraform Provider for Taikun

- Website: https://www.terraform.io
- Forum: https://discuss.hashicorp.com
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://raw.githubusercontent.com/hashicorp/terraform-website/ff7a019259feb18b0a7b2f0ed7ce70b2e3e5d02f/content/source/assets/images/logo-terraform-main.svg" width="600px">

## Using the provider

### Requirements
-	[Terraform](https://www.terraform.io/downloads.html) >= 0.14.x

### Documentation
The provider's documentation is available on the [Terraform registry](https://registry.terraform.io/providers/itera-io/taikun/latest/docs).
See the section titled *USE PROVIDER* to start using it.

## Developing the provider
### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.14.x
-	[Go](https://golang.org/doc/install) >= 1.17

## <a name="build"></a> Building and installing the provider

### Building the provider locally
If you have the Go toolchain installed on your machine, you can build and install the provider with the following command.
```sh
make install
```
Otherwise, you can build the provider in a docker container.
```sh
make dockerinstall
```

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

### Running the locally built provider
To tell Terraform to retrieve the provider locally instead of fetching it from the registry, use the following terraform configuration block.
```tf
terraform {
  required_providers {
    taikun = {
      source  = "itera-io/dev/taikun"
    }
  }
}
```
