# Terraform Provider for Taikun

- Website: https://www.terraform.io
- Forum: https://discuss.hashicorp.com
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://raw.githubusercontent.com/hashicorp/terraform-website/ff7a019259feb18b0a7b2f0ed7ce70b2e3e5d02f/content/source/assets/images/logo-terraform-main.svg" width="600px">

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.14.x
-	[Go](https://golang.org/doc/install) >= 1.17

## <a name="build"></a> Building and installing the provider

Clone the repo to Go's `src` directory.
```sh
mkdir -p $GOPATH/src/github.com/itera-io
git clone git@github.com:itera-io/terraform-provider-taikun $GOPATH/src/github.com/itera-io/terraform-provider-taikun
cd $GOPATH/src/github.com/itera-io/terraform-provider-taikun
```

- If you have the Go toolchain installed on your machine, you can build and install the provider with the following command.
```sh
make install
```
- Otherwise, you can build the provider in a docker container.
```sh
make dockerinstall
```

Either way, the provider will be installed in `~/.terraform.d/plugins/itera-io/dev/taikun/`.

## Using the provider

Until the provider is listed on Terraform's plugin [registry](https://registry.terraform.io/browse/providers), the provider must be installed locally (see [Building and installing the provider](#build)).
To tell Terraform to retrieve the provider locally, use the following terraform configuration block.
```tf
terraform {
  required_providers {
    taikun = {
      source  = "itera-io/dev/taikun"
    }
  }
}
```

<!---
 
  TODO add link to documentation on https://registry.terraform.io, once up.
  See the [Taikun Provider documentation](https://registry.terraform.io/providers/itera-io/taikun/latest/docs) to get started using the Taikun provider.

-->

## Developing the provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
