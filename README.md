<div align="center">
    <h1>Terraform Provider for <a href="https://github.com/ctfer-io/ctfd-chall-manager">CTFd Chall-Manager plugin</a></h1>
    <p><b>Time for CTF(d) as Code, with Chall-Manager</b><p>
    <a href="https://pkg.go.dev/github.com/ctfer-io/terraform-provider-ctfdcm"><img src="https://shields.io/badge/-reference-blue?logo=go&style=for-the-badge" alt="reference"></a>
	<a href="https://goreportcard.com/report/github.com/ctfer-io/terraform-provider-ctfdcm"><img src="https://goreportcard.com/badge/github.com/ctfer-io/terraform-provider-ctfdcm?style=for-the-badge" alt="go report"></a>
	<a href="https://coveralls.io/github/ctfer-io/terraform-provider-ctfdcm?branch=main"><img src="https://img.shields.io/coverallsCoverage/github/ctfer-io/terraform-provider-ctfdcm?style=for-the-badge" alt="Coverage Status"></a>
	<br>
	<a href=""><img src="https://img.shields.io/github/license/ctfer-io/terraform-provider-ctfdcm?style=for-the-badge" alt="License"></a>
	<a href="https://github.com/ctfer-io/terraform-provider-ctfdcm/actions?query=workflow%3Aci+"><img src="https://img.shields.io/github/actions/workflow/status/ctfer-io/terraform-provider-ctfdcm/ci.yaml?style=for-the-badge&label=CI" alt="CI"></a>
	<a href="https://github.com/ctfer-io/terraform-provider-ctfdcm/actions/workflows/codeql-analysis.yaml"><img src="https://img.shields.io/github/actions/workflow/status/ctfer-io/terraform-provider-ctfdcm/codeql-analysis.yaml?style=for-the-badge&label=CodeQL" alt="CodeQL"></a>
    <br>
    <a href="https://securityscorecards.dev/viewer/?uri=github.com/ctfer-io/terraform-provider-ctfdcm"><img src="https://img.shields.io/ossf-scorecard/github.com/ctfer-io/terraform-provider-ctfdcm?label=openssf%20scorecard&style=for-the-badge" alt="OpenSSF Scoreboard"></a>
</div>

## Why this ?

To manipulate our [Terraform Provider for CTFd](https://github.com/ctfer-io/terraform-provider-ctfd) along with the [Chall-Manager plugin](https://github.com/ctfer-io/ctfd-chall-manager).
This enable reusing the configuration thus integrate seamlessly.

## How to use it ?

Install the **Terraform Provider for CTFd** by setting the following in your `main.tf file`.
```hcl
terraform {
    required_providers {
        ctfdcm = {
            source = "registry.terraform.io/ctfer-io/ctfdcm"
        }
    }
}

provider "ctfdcm" {
    url = "https://my-ctfd.lan"
}
```

We recommend setting the environment variable `CTFD_API_KEY` to enable the provider to communicate with your CTFd instance.

Then, you could use a `ctfdcm_challenge_dynamiciac` resource to setup your CTFd challenge, with for instance the following configuration.
```hcl
resource "ctfdcm_challenge_dynamiciac" "my_challenge" {
    name        = "My Challenge"
    category    = "Some category"
    description = <<-EOT
        My superb description !

        And it's multiline :o
    EOT
    state       = "visible"
    value       = 500

    shared          = true
    destroy_on_flag = true
    mana_cost       = 1
    scenario        = "localhost:5000/some/scenario:v0.1.0"
    timeout         = 600
}
```

## OpenTelemetry support

Understanding what is going on under the hood or what could fail throughout the CTF lifecycle remains an important concern, even with such provider. For better understandability, we ship support for OpenTelemetry.

You can configure it using [the SDK environment variables](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/).

Note that CTFd **does not support it natively**, you may want to use our [instrumented and packaged CTFd](https://github.com/ctfer-io/ctfd-packaged) or proceed similarly for auto-instrumentation.

Also, the provider uses the `always` sampler hence we recommend you use a [Collector probability sampler](https://opentelemetry.io/docs/specs/otel/trace/tracestate-probability-sampling/). An example follows, with arbitrary values.
```yaml
processors:
  probabilistic_sampler:
    hash_seed: 22
    sampling_percentage: 22

service:
  pipelines:
    traces:
      receivers: [...]
      processors: [probabilistic_sampler, ...]
      exporters: [...]
```

A more complete example is [available here](./examples/opentelemetry).
