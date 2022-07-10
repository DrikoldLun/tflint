package terraformrules

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	sdk "github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint/tflint"
)

// TerraformModulePinnedSourceRule checks unpinned or default version module source
type TerraformModulePinnedSourceRule struct {
	attributeName string
}

type terraformModulePinnedSourceRuleConfig struct {
	Style           string   `hcl:"style,optional"`
	DefaultBranches []string `hcl:"default_branches,optional"`
}

type TerraformModulePinnedSourceRuleModuleCall struct {
	name       string
	source     string
	sourceAttr *hclext.Attribute
	block      *hclext.Block
}

// NewTerraformModulePinnedSourceRule returns new rule with default attributes
func NewTerraformModulePinnedSourceRule() *TerraformModulePinnedSourceRule {
	return &TerraformModulePinnedSourceRule{
		attributeName: "source",
	}
}

// Name returns the rule name
func (r *TerraformModulePinnedSourceRule) Name() string {
	return "terraform_module_pinned_source"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformModulePinnedSourceRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformModulePinnedSourceRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformModulePinnedSourceRule) Link() string {
	return tflint.ReferenceLink(r.Name())
}

// Check checks if module source version is pinned
// Note that this rule is valid only for Git or Mercurial source
func (r *TerraformModulePinnedSourceRule) Check(runner *tflint.Runner) error {
	if !runner.TFConfig.Path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	log.Printf("[TRACE] Check `%s` rule for `%s` runner", r.Name(), runner.TFConfigPath())

	config := terraformModulePinnedSourceRuleConfig{Style: "flexible"}
	config.DefaultBranches = append(config.DefaultBranches, "master", "main", "default", "develop")
	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		return err
	}

	body, diags := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "module",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{{Name: "source"}},
				},
			},
		},
	}, sdk.GetModuleContentOption{})
	if diags.HasErrors() {
		return diags
	}

	for _, block := range body.Blocks {
		sourceAttr, exists := block.Body.Attributes["source"]
		if !exists {
			continue
		}

		var source string
		if diags := gohcl.DecodeExpression(sourceAttr.Expr, nil, &source); diags.HasErrors() {
			return diags
		}

		module := &TerraformModulePinnedSourceRuleModuleCall{
			name:       block.Labels[0],
			source:     source,
			sourceAttr: sourceAttr,
			block:      block,
		}

		if err := r.checkModule(runner, module, config); err != nil {
			return err
		}
	}

	return nil
}

func (r *TerraformModulePinnedSourceRule) checkModule(runner *tflint.Runner, module *TerraformModulePinnedSourceRuleModuleCall, config terraformModulePinnedSourceRuleConfig) error {
	log.Printf("[DEBUG] Walk `%s` attribute", module.name+".source")

	source, err := getter.Detect(module.source, filepath.Dir(module.block.DefRange.Filename), []getter.Detector{
		// https://github.com/hashicorp/terraform/blob/51b0aee36cc2145f45f5b04051a01eb6eb7be8bf/internal/getmodules/getter.go#L30-L52
		new(getter.GitHubDetector),
		new(getter.GitDetector),
		new(getter.BitBucketDetector),
		new(getter.GCSDetector),
		new(getter.S3Detector),
		new(getter.FileDetector),
	})
	if err != nil {
		return err
	}

	u, err := url.Parse(source)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "git", "hg":
	default:
		return nil
	}

	if u.Opaque != "" {
		// for git:: or hg:: pseudo-URLs, Opaque is :https, but query will still be parsed
		query := u.RawQuery
		u, err = url.Parse(strings.TrimPrefix(u.Opaque, ":"))
		if err != nil {
			return err
		}

		u.RawQuery = query
	}

	if u.Hostname() == "" {
		runner.EmitIssue(
			r,
			fmt.Sprintf("Module source %q is not a valid URL", module.source),
			module.sourceAttr.Expr.Range(),
		)

		return nil
	}

	query := u.Query()

	if ref := query.Get("ref"); ref != "" {
		return r.checkRevision(runner, module, config, "ref", ref)
	}

	if rev := query.Get("rev"); rev != "" {
		return r.checkRevision(runner, module, config, "rev", rev)
	}

	runner.EmitIssue(
		r,
		fmt.Sprintf(`Module source "%s" is not pinned`, module.source),
		module.sourceAttr.Expr.Range(),
	)

	return nil
}

func (r *TerraformModulePinnedSourceRule) checkRevision(runner *tflint.Runner, module *TerraformModulePinnedSourceRuleModuleCall, config terraformModulePinnedSourceRuleConfig, key string, value string) error {
	switch config.Style {
	// The "flexible" style requires a revision that is not a default branch
	case "flexible":
		for _, branch := range config.DefaultBranches {
			if value == branch {
				runner.EmitIssue(
					r,
					fmt.Sprintf("Module source \"%s\" uses a default branch as %s (%s)", module.source, key, branch),
					module.sourceAttr.Expr.Range(),
				)

				return nil
			}
		}
	// The "semver" style requires a revision that is a semantic version
	case "semver":
		_, err := semver.NewVersion(value)
		if err != nil {
			runner.EmitIssue(
				r,
				fmt.Sprintf("Module source \"%s\" uses a %s which is not a semantic version string", module.source, key),
				module.sourceAttr.Expr.Range(),
			)
		}
	default:
		return fmt.Errorf("`%s` is invalid style", config.Style)
	}

	return nil
}
