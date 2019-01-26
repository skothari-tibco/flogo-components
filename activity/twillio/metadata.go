package twillio

import (
	"github.com/project-flogo/core/data/coerce"
)

type Input struct {
	AccountSID string `md:"accountSID"`
	AuthToken  string `md:"authtoken"`
	From       string `md:"from"`
	To         string `md:"to"`
	Message    string `md:"message"`
}
type Output struct {
	Result string `md:"result"`
}

func (o *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"accountSID": o.AccountSID,
		"authtoken":  o.AuthToken,
		"from":       o.From,
		"to":         o.To,
		"message":    o.Message,
	}
}

func (o *Input) FromMap(values map[string]interface{}) error {

	var err error
	o.AccountSID, err = coerce.ToString(values["accountSID"])
	if err != nil {
		return err
	}

	o.AuthToken, err = coerce.ToString(values["authtoken"])
	if err != nil {
		return err
	}

	o.From, err = coerce.ToString(values["from"])
	if err != nil {
		return err
	}

	o.To, err = coerce.ToString(values["to"])
	if err != nil {
		return err
	}

	o.Message, err = coerce.ToString(values["message"])
	if err != nil {
		return err
	}

	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"result": o.Result,
	}
}
func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Result, err = coerce.ToString(values["result"])
	if err != nil {
		return err
	}
	return nil
}
