# Terraform Provider for Taikun

- Website: https://www.terraform.io
- Forum: https://discuss.hashicorp.com
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://raw.githubusercontent.com/hashicorp/terraform-website/ff7a019259feb18b0a7b2f0ed7ce70b2e3e5d02f/content/source/assets/images/logo-terraform-main.svg" width="600px">

## Using the provider

### Requirements
-	[Terraform](https://www.terraform.io/downloads.html) >= 0.14.x

### QuickStart templates
Quickstart templates for creating VMs, K8s, Applications, Virtual clusters and more
are available in this repository at path [examples/quickstart-templates](https://github.com/itera-io/terraform-provider-taikun/tree/dev/examples/quickstart-templates)

### Documentation
The provider's documentation is available on the [Terraform registry](https://registry.terraform.io/providers/itera-io/taikun/latest/docs).
See the section titled *USE PROVIDER* to start using it.

## Developing the provider
### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.14.x
-	[Go](https://golang.org/doc/install) >= 1.19

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
### Documenting the provider

We use [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs) to generate documentation for the provider.

To generate or update documentation, run `go generate` locally or run the [generate documentation](https://github.com/itera-io/terraform-provider-taikun/actions/workflows/generate_documentation.yml) workflow with your target branch as an input.

This reads the templates in the [templates](./templates) directory, the Terraform configuration examples in the [examples](./examples) directory and finally the resource (or data source) schemas themselves to generate the documentation in the [docs](./docs) directory.

In other words, suppose you are creating a new resource `taikun_project`, you would need to add the following files before running `go generate`.
- A Terraform configuration example in `./examples/resources/taikun_project/resource.tf`
- A terraform import script in `./examples/resources/taikun_project/import.sh` (this is usually just `terraform import <resource type>.<name> <id>`)
- A template in `templates/resources/project.md.tmpl`

As mentioned previously, the documentation of provider releases is available on the [Terraform registry](https://registry.terraform.io/providers/itera-io/taikun/latest/docs).

The [Doc Preview Tool](https://registry.terraform.io/tools/doc-preview) by Hashicorp is also a useful way to preview the final look of the documentation.

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

## Running acceptance tests

### Prerequisites

Running the Taikun Terraform Provider's acceptance tests requires setting some
environment variables prior to launching the test suite.

The following environment variables are required to authenticate
the provider.
```
TAIKUN_EMAIL
TAIKUN_PASSWORD
```

In order to run tests that create resources linked to external services such as
AWS, Azure, GCP, OpenStack or Prometheus, set the following variables.
```sh
# AWS
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY

# Azure
AZURE_CLIENT_ID
AZURE_SECRET
AZURE_SUBSCRIPTION
AZURE_TENANT

# GCP
GCP_BILLING_ACCOUNT
GCP_CONFIG_FILE
GCP_FOLDER_ID
GCP_REGION

# OpenStack
OS_AUTH_URL
OS_INTERFACE
OS_PASSWORD
OS_PROJECT_NAME
OS_REGION_NAME
OS_USERNAME
OS_USER_DOMAIN_NAME

# Prometheus
PROMETHEUS_PASSWORD
PROMETHEUS_URL
PROMETHEUS_USERNAME
```

This list of environment variables can also be found in the
[provider_test.go](./taikun/provider_test.go) file, which defines the pre-check
functions for the acceptance tests.

### Running the full suite of tests

In order to run the full suite of Acceptance tests, run `make testacc`.

```sh
$ make testacc
```

*Note:* Acceptance tests create real resources, and often cost money to run.

*Note:* At the time of writing, running the full suite of Acceptance tests
takes on average two to three hours.

### Running specific tests

In order to run only some specific tests, set the `TESTARGS` environment
variable when calling `make testacc`.

The value of `TESTARGS` must be set to `-run <regexp>` where `regexp` is a
regular expression matching the identifiers of tests you wish to run.

For example, to run all tests related to the `taikun_showback_rules` data
source, run the following command.
```sh
TESTARGS='-run TestAccDataSourceTaikunShowbackRules' make testacc
```

At the moment of writing, this will run both of the following tests.
```
TestAccDataSourceTaikunShowbackRules
TestAccDataSourceTaikunShowbackRulesWithFilter
```

If you want to run only `TestAccDataSourceTaikunShowbackRules`, run the
following command.
```sh
TESTARGS='-run TestAccDataSourceTaikunShowbackRules$' make testacc
```
Notice we added a `$` sign at the end of the regular expression to match the
end of line. Thus, `TestAccDataSourceTaikunShowbackRulesWithFilter` will be
ignored.

To know more about the `-run <regexp>` test flag and other go test flags, see the
[go-testflag (7) man page](https://manpages.debian.org/testing/golang-go/go-testflag.7.en.html#run)

### Rigorous testing
For testing the prepared bundles of CI acceptance tests, you can use the ```make rtestacc``` command while uncommenting the correct line of tests in makefile.

The ```rtestacc``` command can be also used to specify how many parallel threads you want to run with ```retestacc1-4```.
You can also run the test bundle 4 times in sequence with increasing number of 1-4 threads with ```rtestaccrigorous```.

### Pipeline
The CI pipeline in GitHub Actions consists of 4 logical sequential steps.
- Build
  - Build - try to compile the code
  - Lint - Find any linting errors
- Test Taikun - Test Taikun functionality that does not create any instances, vms, servers.
  - 1alpha - A random half of the acceptance tests
  - 1beta - A second half of the acceptance tests
- Test creating reources
  - 2x - Test if Taikun can create Openstack resources 
  - 3x - Test if Taikun can create AWS resources
  - 4x - Test if Taikun can create Azure resources
- Test Kubernetes creation
  - Test if Taikun can create Kubernetes cluster in Openstack.

### Parallel testing
When running many tests in parallel, a race condition from Terraform can sometimes be observed.
If this happend during tests that do not create any resources (1alpha/1beta), the tests are simply rerun again, up to 3 times.

``` bash
resource_***_cloud_credential_azure_test.go:122: failed to create new working directory: unable to disable terraform-exec provider verification: fork/exec /tmp/plugintest-terraform1517105639/terraform: text file busy
```