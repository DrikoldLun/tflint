package terraformrules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint/tflint"
	"testing"
)

func Test_AzurermArgsOrderRule(t *testing.T) {
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
		//		{
		//			Name: "comments",
		//			Content: `
		//#Configure the Azure provider
		//terraform {
		//  required_providers {
		//    azurerm = {
		//      source  = "hashicorp/azurerm"
		//      version = "~> 3.0.2"
		//    }
		//  }
		///*
		//  cloud {
		//    organization = "lunz-test"
		//    workspaces {
		//      name = "learn-terraform-azure"
		//    }
		//  }
		//*/
		//}
		//
		//provider "azurerm" {
		//  features {}
		//}
		//
		//resource "azurerm_resource_group" "rg" {
		//  #name     = "myTFResourceGroup"
		//  location = "westus2"
		//  tags = {
		//    Team = "DevOps"
		//    Environment = "Terraform Getting Started"
		//    #Team = "DevOps"
		//  }
		//}
		//
		//resource "azurerm_virtual_network" "vnet" {
		//  #name                = "myTFVnet"
		//  #address_space       = ["10.0.0.0/16"]
		//  location            = "westus2"
		//  resource_group_name = azurerm_resource_group.rg.name
		//}`,
		//			Expected: tflint.Issues{},
		//		},
		{
			Name: "comments",
			Content: `
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_container_group" "example" {
  name                = "example-continst"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"

  container {
    name   = "hello-world"
    image  = "mcr.microsoft.com/azuredocs/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "1.5"

    ports {
	  port     = 443
	  protocol = "TCP"
    }
  }

  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }

  tags = {
    environment = "testing"
  }
}`,
			Expected: tflint.Issues{
				{
					Rule: NewAzurermArgsOrderRule(),
					Message: `Arguments are not sorted in azurerm doc order, correct order is:
resource "azurerm_resource_group" "example" {
  location = "West Europe"
  name     = "example-resources"
}
`,
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 11,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 33,
						},
					},
				},
				{
					Rule: NewAzurermArgsOrderRule(),
					Message: `Arguments are not sorted in azurerm doc order, correct order is:
  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-helloworld:latest"
    memory = "1.5"
    name   = "hello-world"

    ports {
	  port     = 443
	  protocol = "TCP"
    }
  }
`,
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   15,
							Column: 3,
						},
						End: hcl.Pos{
							Line:   15,
							Column: 12,
						},
					},
				},
				{
					Rule: NewAzurermArgsOrderRule(),
					Message: `Arguments are not sorted in azurerm doc order, correct order is:
  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }
`,
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   27,
							Column: 3,
						},
						End: hcl.Pos{
							Line:   27,
							Column: 12,
						},
					},
				},
				{
					Rule: NewAzurermArgsOrderRule(),
					Message: `Arguments are not sorted in azurerm doc order, correct order is:
resource "azurerm_container_group" "example" {
  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-helloworld:latest"
    memory = "1.5"
    name   = "hello-world"

    ports {
	  port     = 443
	  protocol = "TCP"
    }
  }
  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }
  location            = azurerm_resource_group.example.location
  name                = "example-continst"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name

  dns_name_label      = "aci-label"
  ip_address_type     = "Public"
  tags = {
    environment = "testing"
  }
}
`,
					Range: hcl.Range{
						Filename: "config.tf",
						Start: hcl.Pos{
							Line:   7,
							Column: 11,
						},
						End: hcl.Pos{
							Line:   7,
							Column: 34,
						},
					},
				},
			},
		},
	}

	rule := NewAzurermArgsOrderRule()

	for _, tc := range cases {
		runner := tflint.TestRunner(t, map[string]string{"config.tf": tc.Content})

		//for _, file := range runner.Files() {
		//	tokens, _ := hclsyntax.LexConfig(file.Bytes, "config.tf", hcl.InitialPos)
		//	for _, token := range tokens {
		//		fmt.Println(token.Type)
		//		fmt.Println(string(token.Bytes))
		//	}
		//}

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		tflint.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
