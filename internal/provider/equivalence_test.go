package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEquivalence_DataSourceAndFunction(t *testing.T) {
	testCases := []struct {
		name        string
		yamlContent string
	}{
		{
			name: "nested_objects",
			yamlContent: `
key1: value1
key2:
  nested: value2
`,
		},
		{
			name: "arrays",
			yamlContent: `
items:
  - name: item1
    value: val1
  - name: item2
    value: val2
`,
		},
		{
			name: "mixed_types",
			yamlContent: `
string_val: "hello"
int_val: 42
float_val: 3.14
bool_val: true
null_val: null
`,
		},
		{
			name: "deeply_nested",
			yamlContent: `
a:
  b:
    c:
      d:
        e: "deep"
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
provider "yamlflattener" {}

locals {
  yaml_content = <<EOT
%s
EOT
}

data "yamlflattener_flatten" "ds" {
  yaml_content = local.yaml_content
}

output "are_equal" {
  value = jsonencode(data.yamlflattener_flatten.ds.flattened) == jsonencode(provider::yamlflattener::flatten(local.yaml_content))
}
`, tc.yamlContent),
						Check: resource.TestCheckOutput("are_equal", "true"),
					},
				},
			})
		})
	}
}

func TestEquivalence_ProviderMaxDepthConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "yamlflattener" {
  max_depth = 20
}

output "deep_value" {
  value = provider::yamlflattener::flatten(<<EOT
level1:
  level2:
    level3:
      level4:
        level5:
          value: "deep value"
EOT
)["level1.level2.level3.level4.level5.value"]
}
`,
				Check: resource.TestCheckOutput("deep_value", "deep value"),
			},
		},
	})
}
