module github.com/terraform-providers/terraform-provider-aws

go 1.15

require (
	github.com/aws/aws-sdk-go v1.36.19
	github.com/beevik/etree v1.1.0
	github.com/fatih/color v1.9.0 // indirect
	github.com/hashicorp/aws-sdk-go-base v0.7.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-version v1.2.1
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.0
	github.com/jen20/awspolicyequivalence v1.1.0
	github.com/keybase/go-crypto v0.0.0-20161004153544-93f5b35093ba
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mitchellh/copystructure v1.0.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-testing-interface v1.14.1
	github.com/pquerna/otp v1.3.0
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	golang.org/x/tools v0.0.0-20201230224404-63754364767c
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/hashicorp/terraform-plugin-test/v2 => github.com/Mongey/terraform-plugin-test/v2 v2.1.3-0.20201231030340-31624e2320cd
