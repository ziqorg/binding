package binding

import (
	"encoding/json"
	"errors"

	"github.com/cloudwego/gjson"
	"github.com/expr-lang/expr"
	"github.com/tidwall/sjson"
)

func marshal(i interface{}) string {
	s, _ := json.Marshal(i)
	return string(s)
}

func getRaw(raw gjson.Result) interface{} {
	return raw.Value()
}

type Binding struct {
	root gjson.Result
}

func NewBinding(raw interface{}) Binding {
	if s, ok := raw.(string); ok {
		return Binding{
			root: gjson.Parse(s),
		}
	}
	return Binding{
		root: gjson.Parse(marshal(raw)),
	}
}

// GetRoot
// Get root state -> json
func (bind *Binding) GetRoot() gjson.Result {
	return bind.root.Get("@this")
}

func (bind *Binding) GetRawRoot() interface{} {
	root := bind.GetRoot()
	return getRaw(root)
}

// Evaluate
// bind.Evaluate("input.abc == 'abc'") (true, nil)
func (bind *Binding) Evaluate(condition string) (bool, error) {
	state := bind.GetRawRoot()
	program, err := expr.Compile(condition, expr.Env(state))
	if err != nil {
		return false, err
	}
	res, err := expr.Run(program, state)
	if err != nil {
		return false, err
	}
	ok, parsable := res.(bool)
	if !parsable {
		return false, errors.New("unable to parse expression output to boolean")
	}
	return ok, nil
}

// Get
// Get("input.dataset.owners")
func (bind *Binding) Get(path string) gjson.Result {
	return bind.root.Get(path)
}

func (bind *Binding) GetRaw(path string) interface{} {
	return getRaw(bind.Get(path))
}

// Set
// Set("check_duplicate.duplicate", true)
func (bind *Binding) Set(path string, value interface{}) (string, error) {
	setter, err := sjson.Set(bind.root.String(), path, value)
	if err != nil {
		return "", err
	}
	bind.root = NewBinding(setter).root
	return setter, nil
}
