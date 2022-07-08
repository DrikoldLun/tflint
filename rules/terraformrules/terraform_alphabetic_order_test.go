package terraformrules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint/tflint"
	"testing"
)

func Test_TerraformAlphabeticOrderRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected tflint.Issues
	}{
		{
			Name: "",
			Content: `
resource "azurerm_resource_group" "rg" {
  name     = "myTFResourceGroup"
  location = "westus2"
}`,
			Expected: tflint.Issues{
				{
					Rule:    NewTerraformAlphabeticOrderRule(),
					Message: "Arguments `name` and `location` are not sorted in alphabetic order",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 9,
						},
					},
				},
			},
		},
		{
			Name: "inner block",
			Content: `
resource "azurerm_resource_group" "rg" {
  location = "westus2"
  name     = "myTFResourceGroup"
  tags = {
	Team = "DevOps"
    Environment = "Terraform Getting Started"
  }
}`,
			Expected: tflint.Issues{
				{
					Rule:    NewTerraformAlphabeticOrderRule(),
					Message: "Arguments `Team` and `Environment` are not sorted in alphabetic order",
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   5,
							Column: 3,
						},
						End: hcl.Pos{
							Line:   5,
							Column: 7,
						},
					},
				},
			},
		},
	}

	rule := NewTerraformAlphabeticOrderRule()

	for _, tc := range cases {
		runner := tflint.TestRunner(t, map[string]string{"config.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		tflint.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
