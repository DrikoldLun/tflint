package terraformrules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint/tflint"
	"strings"
	"testing"
)

func Test_TerraformAlphabeticOrderRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected tflint.Issues
	}{
		//		{
		//			Name: "simple block",
		//			Content: `
		//resource "azurerm_resource_group" "rg" {
		//  name     = "myTFResourceGroup"
		//  location = "westus2"
		//}`,
		//			Expected: tflint.Issues{
		//				{
		//					Rule:    NewTerraformAlphabeticOrderRule(),
		//					Message: "Arguments `name` and `location` are not sorted in alphabetic order",
		//					Range: hcl.Range{
		//						Filename: "config.tf",
		//						Start: hcl.Pos{
		//							Line:   2,
		//							Column: 1,
		//						},
		//						End: hcl.Pos{
		//							Line:   2,
		//							Column: 9,
		//						},
		//					},
		//				},
		//			},
		//		},
		//		{
		//			Name: "inner block",
		//			Content: `
		//resource "azurerm_resource_group" "rg" {
		//  location = "westus2"
		//  name     = "myTFResourceGroup"
		//  tags = {
		//	Team = "DevOps"
		//    Environment = "Terraform Getting Started"
		//  }
		//}`,
		//			Expected: tflint.Issues{
		//				{
		//					Rule:    NewTerraformAlphabeticOrderRule(),
		//					Message: "Arguments `Team` and `Environment` are not sorted in alphabetic order",
		//					Range: hcl.Range{
		//						Filename: "config.tf",
		//						Start: hcl.Pos{
		//							Line:   5,
		//							Column: 3,
		//						},
		//						End: hcl.Pos{
		//							Line:   5,
		//							Column: 7,
		//						},
		//					},
		//				},
		//			},
		//		},
		//		{
		//			Name: "multiple blocks",
		//			Content: `
		//terraform {
		//  required_providers {
		//    azurerm = {
		//      source  = "hashicorp/azurerm"
		//      version = "~> 3.0.2"
		//    }
		//  }
		//}
		//
		//provider "azurerm" {
		//  features {}
		//}
		//
		//resource "azurerm_resource_group" "rg" {
		//  name     = "myTFResourceGroup"
		//  location = "westus2"
		//  tags = {
		//    Environment = "Terraform Getting Started"
		//    Team = "DevOps"
		//  }
		//}
		//
		//resource "azurerm_virtual_network" "vnet" {
		//  name                = "myTFVnet"
		//  address_space       = ["10.0.0.0/16"]
		//  location            = "westus2"
		//  resource_group_name = azurerm_resource_group.rg.name
		//}`,
		//			Expected: tflint.Issues{
		//				{
		//					Rule:    NewTerraformAlphabeticOrderRule(),
		//					Message: "Arguments `name` and `location` are not sorted in alphabetic order",
		//					Range: hcl.Range{
		//						Filename: "config.tf",
		//						Start: hcl.Pos{
		//							Line:   15,
		//							Column: 1,
		//						},
		//						End: hcl.Pos{
		//							Line:   15,
		//							Column: 9,
		//						},
		//					},
		//				},
		//			},
		//		},
		{
			Name: "comments",
			Content: `
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.2"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  #name     = "myTFResourceGroup"
  location = "westus2"
  tags = {
    Environment = "Terraform Getting Started"
    Team = "DevOps"
  }
}

resource "azurerm_virtual_network" "vnet" {
  #name                = "myTFVnet"
  #address_space       = ["10.0.0.0/16"]
  location            = "westus2"
  resource_group_name = azurerm_resource_group.rg.name
}`,
			Expected: tflint.Issues{},
		},
	}

	rule := NewTerraformAlphabeticOrderRule()

	for _, tc := range cases {
		runner := tflint.TestRunner(t, map[string]string{"config.tf": tc.Content})
		for _, file := range runner.Files() {
			tokens, _ := hclsyntax.LexConfig(file.Bytes, "config.tf", hcl.InitialPos)
			fmt.Println(tokens[len(tokens)-1].Range.End.Line)
			fmt.Println(len(strings.Split(string(file.Bytes), "\n")))
		}
		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		tflint.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
