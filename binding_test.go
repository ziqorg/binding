package binding

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

const yamlData = `
    input:
        abc: 12346561242130948129
        urn: 1234
        args1: 12345
        args2: true
        args3: arn:datadog_metrics
        nested_args:
            args1: ok
        nested_args_list:
        - args1:
            args2: 12345
    `

func TestBindings(t *testing.T) {

	var result map[string]interface{}
	err := yaml.Unmarshal([]byte(yamlData), &result)
	assert.NoError(t, err)

	binding := NewBinding(result)

	bind, err := binding.Get("input.urn")
	assert.Equal(t, bind, "1234")
	assert.NoError(t, err)

	bind, err = binding.Get("not_a_key.ggwp")
	assert.Empty(t, bind)
	assert.Error(t, err)

	bind, err = binding.Set("not_a_key.ggwp", true)
	assert.Equal(t, bind, "ok")
	assert.NoError(t, err)

	bind, err = binding.Set("not_a_key.ggwp", 1.234)
	assert.Equal(t, bind, "ok")
	assert.NoError(t, err)

	bind, err = binding.Set("not_a_key.ggwp", "now_is_a_key")
	assert.Equal(t, bind, "ok")
	assert.NoError(t, err)
}

func TestEvaluateBindingExpression(t *testing.T) {
	var result map[string]interface{}
	err := yaml.Unmarshal([]byte(yamlData), &result)
	assert.NoError(t, err)

	binding := NewBinding(result)

	ok, err := binding.Evaluate("input.urn == 1234 && input.nested_args.args1 == 'ok'")
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = binding.Evaluate("input.urn == 1234 && input.nested_args.args1 == 'o'")
	assert.NoError(t, err)
	assert.False(t, ok)

	_, _ = binding.Set("input.nested_args.args1", "o")
	ok, err = binding.Evaluate("input.urn == 1234 && input.nested_args.args1 == 'o'")
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = binding.Evaluate("input.nested_args_list[0].args1.args2 == 12345")
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = binding.Evaluate("input.nested_args_list[0.args1.args3 == 12345")
	assert.Error(t, err)
	assert.False(t, ok)

	ok, err = binding.Evaluate("input.nested_args_list[0].args1.args3 == 1234")
	assert.NoError(t, err)
	assert.False(t, ok)
}
