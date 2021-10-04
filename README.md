# Terraform Provider for Taikun

- Website: https://www.terraform.io
- Forum: https://discuss.hashicorp.com
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/itera-io/terraform-provider-taikun`

```sh
$ mkdir -p $GOPATH/src/github.com/itera-io; cd $GOPATH/src/github.com/itera-io
$ git clone git@github.com:itera-io/terraform-provider-taikun
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/itera-io/terraform-provider-taikun
$ go install
```

## Using the provider

[//] # (TODO add link to documentation on https://registry.terraform.io, once up)

See the [Taikun Provider documentation](https://registry.terraform.io/providers/) to get started using the Taikun provider.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
