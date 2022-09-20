module github.com/worldremit/terraform-provider-slack

go 1.16

require (
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/terraform-plugin-docs v0.10.1
	github.com/hashicorp/terraform-plugin-log v0.7.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.23.0
	github.com/jmatsu/terraform-provider-slack v0.9.0
	github.com/slack-go/slack v0.11.0
	gopkg.in/djherbis/times.v1 v1.3.0
)

replace github.com/jmatsu/terraform-provider-slack => ./
