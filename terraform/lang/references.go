package lang

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint/terraform/addrs"
	"github.com/terraform-linters/tflint/terraform/tfdiags"
)

// References finds all of the references in the given set of traversals,
// returning diagnostics if any of the traversals cannot be interpreted as a
// reference.
//
// This function does not do any de-duplication of references, since references
// have source location information embedded in them and so any invalid
// references that are duplicated should have errors reported for each
// occurence.
//
// If the returned diagnostics contains errors then the result may be
// incomplete or invalid. Otherwise, the returned slice has one reference per
// given traversal, though it is not guaranteed that the references will
// appear in the same order as the given traversals.
func References(traversals []hcl.Traversal) ([]*addrs.Reference, tfdiags.Diagnostics) {
	if len(traversals) == 0 {
		return nil, nil
	}

	var diags tfdiags.Diagnostics
	refs := make([]*addrs.Reference, 0, len(traversals))

	for _, traversal := range traversals {
		ref, refDiags := addrs.ParseRef(traversal)
		diags = diags.Append(refDiags)
		if ref == nil {
			continue
		}
		refs = append(refs, ref)
	}

	return refs, diags
}

// ReferencesInExpr is a helper wrapper around References that first searches
// the given expression for traversals, before converting those traversals
// to references.
func ReferencesInExpr(expr hcl.Expression) ([]*addrs.Reference, tfdiags.Diagnostics) {
	if expr == nil {
		return nil, nil
	}
	traversals := expr.Variables()
	return References(traversals)
}
