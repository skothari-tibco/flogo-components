package shell

import (
	"github.com/project-flogo/core/data/coerce"
)

type Input struct {
	Command string `md:"command"`
}
type Output struct {
	Result interface{} `md:"result"`
	Error  interface{} `md:"error"`
}

func (o *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"command": o.Command,
	}
}

func (o *Input) FromMap(values map[string]interface{}) error {

	var err error
	o.Command, err = coerce.ToString(values["command"])
	if err != nil {
		return err
	}

	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"result": o.Result,
		"error":  o.Error,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	o.Result = values["result"]

	o.Error = values["error"]

	return nil
}
