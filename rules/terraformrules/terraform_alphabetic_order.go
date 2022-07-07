package terraformrules

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint/tflint"
)

// TerraformAlphabeticOrderRule checks whether the arguments/attributes in a block are sorted in alphabetic order
type TerraformAlphabeticOrderRule struct{}

// NewTerraformAlphabeticOrderRule returns a new rule
func NewTerraformAlphabeticOrderRule() *TerraformAlphabeticOrderRule {
	return &TerraformAlphabeticOrderRule{}
}

// Name returns the rule name
func (r *TerraformAlphabeticOrderRule) Name() string {
	return "terraform_alphabetic_order"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformAlphabeticOrderRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformAlphabeticOrderRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformAlphabeticOrderRule) Link() string {
	return tflint.ReferenceLink(r.Name())
}

// Check checks whether the arguments/attributes in a block are sorted in alphabetic order
func (r *TerraformAlphabeticOrderRule) Check(runner *tflint.Runner) error {
	if !runner.TFConfig.Path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	log.Printf("[TRACE] Check `%s` rule for `%s` runner", r.Name(), runner.TFConfigPath())

	for name, file := range runner.Files() {
		if err := r.checkAlphaSeq(runner, name, file); err != nil {
			return err
		}
	}

	return nil
}

func (r *TerraformAlphabeticOrderRule) checkAlphaSeq(runner *tflint.Runner, filename string, file *hcl.File) error {
	if strings.HasSuffix(filename, ".json") {
		return nil
	}
	var lastIdentToken hclsyntax.Token
	tokens, diags := hclsyntax.LexConfig(file.Bytes, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return diags
	}
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == hclsyntax.TokenIdent {
			lastIdentToken = tokens[i]
		} else if tokens[i].Type == hclsyntax.TokenOBrace {
			i = r.checkAtrAlphaOrder(runner, lastIdentToken, tokens, i+1)
			if i == -1 {
				return nil
			}
		}
	}
	return nil
}

func (r *TerraformAlphabeticOrderRule) checkAtrAlphaOrder(runner *tflint.Runner,
	owner hclsyntax.Token, tokens hclsyntax.Tokens, startIndex int) int {
	var lastIdentifierToken hclsyntax.Token
	for i := startIndex; i < len(tokens); i++ {
		if r.isReceiver(tokens, i) {
			if lastIdentifierToken.Bytes != nil && string(lastIdentifierToken.Bytes) > string(tokens[i].Bytes) {
				runner.EmitIssue(
					r,
					fmt.Sprintf("Attributes `%s` and `%s` are not sorted in alphabetic order",
						string(lastIdentifierToken.Bytes), string(tokens[i].Bytes)),
					owner.Range,
				)
				return -1
			}
			lastIdentifierToken = tokens[i]
		} else if tokens[i].Type == hclsyntax.TokenOBrace {
			i = r.checkAtrAlphaOrder(runner, lastIdentifierToken, tokens, i+1)
			if i == -1 {
				return -1
			}
		} else if tokens[i].Type == hclsyntax.TokenCBrace {
			return i + 1
		}
	}
	return -1
}

func (r *TerraformAlphabeticOrderRule) isReceiver(tokens hclsyntax.Tokens, index int) bool {
	if tokens[index].Type != hclsyntax.TokenIdent {
		return false
	}
	isNewLine := false
	for i := index - 1; i >= 0; i-- {
		if tokens[i].Type == hclsyntax.TokenNewline {
			isNewLine = true
		} else if tokens[i].Type == hclsyntax.TokenEqual {
			if !isNewLine {
				return false
			}
		}
	}
	return true
}
