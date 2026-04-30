package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFlattenDataSource_YAMLContent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFlattenDataSourceConfigYAMLContent(`
key1: value1
key2:
  nested: value2
array:
  - name: item1
  - name: item2
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.key1", "value1"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.key2.nested", "value2"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.array[0].name", "item1"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.array[1].name", "item2"),
				),
			},
		},
	})
}

func TestAccFlattenDataSource_YAMLFile(t *testing.T) {
	tempDir := t.TempDir()
	yamlFilePath := filepath.Join(tempDir, "test.yaml")
	yamlContent := `
key1: value1
key2:
  nested: value2
array:
  - name: item1
  - name: item2
`
	err := os.WriteFile(yamlFilePath, []byte(yamlContent), 0600)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFlattenDataSourceConfigYAMLFile(yamlFilePath),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.key1", "value1"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.key2.nested", "value2"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.array[0].name", "item1"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.array[1].name", "item2"),
				),
			},
		},
	})
}

func TestAccFlattenDataSource_ErrorHandling(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccFlattenDataSourceConfigYAMLContent(`invalid: yaml: : content`),
				ExpectError: regexp.MustCompile(`Invalid YAML Syntax`),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccFlattenDataSourceConfigYAMLFile("/path/to/nonexistent/file.yaml"),
				ExpectError: regexp.MustCompile(`File Access Error`),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccFlattenDataSourceConfigBothInputs(`key: value`, "/path/to/file.yaml"),
				ExpectError: regexp.MustCompile(`Conflicting Inputs`),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccFlattenDataSourceConfigNoInputs(),
				ExpectError: regexp.MustCompile(`Missing Required Input`),
			},
		},
	})
}

func testAccFlattenDataSourceConfigYAMLContent(yamlContent string) string {
	return fmt.Sprintf(`
data "yamlflattener_flatten" "test" {
  yaml_content = <<EOT
%s
EOT
}
`, yamlContent)
}

func testAccFlattenDataSourceConfigYAMLFile(filePath string) string {
	return fmt.Sprintf(`
data "yamlflattener_flatten" "test" {
  yaml_file = %q
}
`, filePath)
}

func testAccFlattenDataSourceConfigBothInputs(yamlContent, filePath string) string {
	return fmt.Sprintf(`
data "yamlflattener_flatten" "test" {
  yaml_content = <<EOT
%s
EOT
  yaml_file = %q
}
`, yamlContent, filePath)
}

func testAccFlattenDataSourceConfigNoInputs() string {
	return `
data "yamlflattener_flatten" "test" {
}
`
}
