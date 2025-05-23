package binding

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/bytedance/sonic/ast"
	"github.com/expr-lang/expr"
)

// PP - Pretty Print Anything
func pp(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

type Binding struct {
	root ast.Node
}

func NewBinding(raw interface{}) Binding {
	s, parsable := raw.(string)
	if parsable {
		return Binding{root: ast.NewRaw(string(s))}
	}
	data, _ := json.Marshal(raw)
	return Binding{root: ast.NewRaw(string(data))}
}

// GetRoot
// Get root state -> json
func (bind *Binding) GetRoot() (map[string]interface{}, error) {
	r, err := bind.root.Raw()

	if err != nil {
		return make(map[string]interface{}), err
	}
	var res map[string]interface{}
	err = json.Unmarshal([]byte(r), &res)
	if err != nil {
		return make(map[string]interface{}), err
	}
	return res, nil
}

// Flatten
// Recursively flatten AST into dot-separated map
func (bind *Binding) Flatten() (map[string]interface{}, error) {
	root, err := bind.GetRoot()
	if err != nil {
		return make(map[string]interface{}), nil
	}
	res := flattenMap(root, "")
	return res, nil
}

// Evaluate
// bind.Evaluate("input.abc == 'abc'") (true, nil)
func (bind *Binding) Evaluate(condition string) (bool, error) {
	state, err := bind.GetRoot()
	if err != nil {
		return false, err
	}
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
func (bind *Binding) Get(path string) (interface{}, error) {
	node := &bind.root
	if path != "." {
		splitPath := strings.Split(path, ".")
		for _, p := range splitPath {
			node = node.Get(p)
			if node.Check() != nil {
				return "", errors.New("no key named: " + p)
			}
		}
	}
	raw, err := node.Raw()
	if err != nil {
		return "", err
	}
	raw = strings.TrimPrefix(raw, "\"")
	raw = strings.TrimSuffix(raw, "\"")
	return autoParse(raw), nil
}

func (bind *Binding) GetMap(path string) (map[string]interface{}, error) {
	b, err := bind.Get(path)
	if err != nil {
		return make(map[string]interface{}), err
	}
	var res map[string]interface{}
	_ = json.Unmarshal([]byte((b).(string)), &res)
	return res, nil
}

// Set
// Set("check_duplicate.duplicate", true)
func (bind *Binding) Set(path string, value interface{}) (string, error) {
	var newValueNode = ast.NewRaw(pp(value))
	if newValueNode.Check() != nil {
		return "", errors.New("invalid value")
	}

	splitPath := strings.Split(path, ".")
	splitPathLen := len(splitPath)
	node := &bind.root
	for i, p := range splitPath {
		if i == splitPathLen-1 {
			_, err := node.Set(p, newValueNode)
			if err != nil {
				return "", err
			}
			break
		}
		if node.Get(p).Check() != nil {
			_, err := node.Set(p, ast.NewObject([]ast.Pair{}))
			if err != nil {
				return "", err
			}
		}
		node = node.Get(p)
	}
	return "ok", nil
}

func (bind *Binding) MergeMap(mp map[string]interface{}) (string, error) {
	newBinding := NewBinding(mp)
	rootMap, err := newBinding.GetRoot()
	if err != nil {
		return "", err
	}
	for k, v := range rootMap {
		_, _ = bind.Set(k, v)
	}
	return "ok", nil
}
