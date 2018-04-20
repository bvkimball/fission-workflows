package builtin

import (
	"errors"

	"github.com/fission/fission-workflows/pkg/types"
	"github.com/fission/fission-workflows/pkg/types/typedvalues"
	"github.com/sirupsen/logrus"
)

const (
	Switch                 = "switch"
	SwitchInputCondition   = "switch"  // required
	SwitchInputCases       = "cases"   // optional
	SwitchInputDefaultCase = "default" // optional

	SwitchCaseKey   = "case"
	SwitchCaseValue = "action"
)

/*
Switch is very similar to how switch-constructs are implemented in most languages.
In this case the switch is limited to evaluating string keys.
The string-switch is matched to one of the cases, or - if none of those match - the default case.

**Specification**

**input**       | required | types             | description
----------------|----------|-------------------|--------------------------------------------------------
switch          | yes      | string            | The string to match to one of the cases.
cases           | no       | list              | List of cases to match to.
default         | no       | *                 | The default value if there is no matching case.

**output** (*) Either the value of the matching case, the default, or nothing (in case the default is not specified).

**Example**

```yaml
# ...
SwitchExample:
  run: switch
  inputs:
    switch: "{ param() }"
    cases:
    - case: foo
      action: bar
    - case: ac
      action: me
    default: 42
# ...
```

A complete example of this function can be found in the [switchwhale](../examples/whales/switchwhale.wf.yaml) example.
*/
type FunctionSwitch struct{}

func (fn *FunctionSwitch) Invoke(spec *types.TaskInvocationSpec) (*types.TypedValue, error) {
	switchVal, err := fn.getSwitch(spec.Inputs)
	if err != nil {
		return nil, err
	}
	cases, defaultCase, err := fn.getCases(spec.Inputs)
	if err != nil {
		return nil, err
	}

	// Evaluate
	logrus.Infof("Switch looking for %v in %v", switchVal, cases)
	if cases != nil {
		tv, ok := cases[switchVal]
		if ok {
			return tv, nil
		}
	}

	return defaultCase, nil
}

func (fn *FunctionSwitch) getSwitch(inputs map[string]*types.TypedValue) (string, error) {
	tv, err := ensureInput(inputs, SwitchInputCondition)
	if err != nil {
		return "", err
	}
	return typedvalues.FormatString(tv)
}

func (fn *FunctionSwitch) getCases(inputs map[string]*types.TypedValue) (map[string]*types.TypedValue,
	*types.TypedValue, error) {
	cases := map[string]*types.TypedValue{}
	defaultCase := inputs[SwitchInputDefaultCase]

	switchCases, ok := inputs[SwitchInputCases]
	if ok {
		ir, err := typedvalues.FormatArray(switchCases)
		if err != nil {
			return nil, nil, err
		}
		for _, c := range ir {
			m, ok := c.(map[string]interface{})
			if !ok {
				logrus.Warnf("Invalid case provided: %t", m)
				return nil, nil, errors.New("invalid case provided")
			}
			tva, err := typedvalues.Parse(m[SwitchCaseValue])
			if err != nil {
				return nil, nil, err
			}

			ic, ok := m[SwitchCaseKey]
			if !ok {
				return nil, nil, errors.New("case in switch does not have a key")
			}
			key, ok := ic.(string)
			if !ok {
				return nil, nil, errors.New("case key should be a string")
			}
			cases[key] = tva
		}
	}
	return cases, defaultCase, nil
}

func switchCase(key string, value interface{}) map[string]interface{} {

	return map[string]interface{}{
		SwitchCaseKey:   key,
		SwitchCaseValue: value,
	}
}
