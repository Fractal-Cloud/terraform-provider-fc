package secrets

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &SecretRefFunction{}

// SecretRefFunction builds a reference to an environment secret, for use as a
// component parameter or link-setting value in place of the raw secret.
type SecretRefFunction struct{}

func NewSecretRefFunction() function.Function {
	return &SecretRefFunction{}
}

func (f *SecretRefFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "secret_ref"
}

func (f *SecretRefFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "References an environment secret by short name",
		Description: "Returns a reference to a secret defined on the environment, for use as a parameter or link-setting " +
			"value. The agent resolves the reference from the environment secret store at reconciliation time, so the " +
			"raw secret never appears in the blueprint.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "short_name",
				Description: "Short name of the environment secret",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *SecretRefFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var shortName string
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &shortName))
	if resp.Error != nil {
		return
	}
	if strings.TrimSpace(shortName) == "" {
		resp.Error = function.NewArgumentFuncError(0, "short_name must not be empty")
		return
	}

	ref, err := SecretRef(shortName)
	if err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return
	}
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, ref))
}

// SecretRef returns the wire form of an environment-secret reference,
// {"$envSecret":"<shortName>"}, as the string a parameter map holds.
func SecretRef(shortName string) (string, error) {
	b, err := json.Marshal(map[string]string{"$envSecret": shortName})
	if err != nil {
		return "", err
	}
	return string(b), nil
}
