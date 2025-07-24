package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"yamlflattener": providerserver.NewProtocol6WithError(New("test")()),
}

func TestProvider(t *testing.T) {
	t.Parallel()

	// Verify provider can be instantiated
	provider := New("test")()
	if provider == nil {
		t.Fatal("provider is nil")
	}
}

func TestProviderConfigure(t *testing.T) {
	t.Parallel()

	// Test basic provider configuration
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "yamlflattener" {
  max_depth = 50
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
				// No specific checks needed, just verifying the configuration is accepted
				),
			},
		},
	})
}

func TestProviderWithDataSource(t *testing.T) {
	t.Parallel()

	// Test provider with data source
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "yamlflattener" {}

data "yamlflattener_flatten" "test" {
  yaml_content = "key: value"
}

output "flattened" {
  value = data.yamlflattener_flatten.test.flattened
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("flattened", "{\"key\":\"value\"}"),
				),
			},
		},
	})
}

func TestProviderWithFunction(t *testing.T) {
	t.Parallel()

	// Test provider with function
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "yamlflattener" {}

output "flattened" {
  value = provider::yamlflattener::flatten("key: value")
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("flattened", "{\"key\":\"value\"}"),
				),
			},
		},
	})
}
