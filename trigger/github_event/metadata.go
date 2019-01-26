package github_event

import (
	"github.com/project-flogo/core/data/coerce"
)

type Output struct {
	Content interface{} `md:"content"`
}
type Settings struct {
	Port string `md:"port"`
	Url  string `md:"url"`
}

type Reply struct {
	Code int         `md:"code"`
	Data interface{} `md:"data"`
}

type HandlerSettings struct {
	GetAll bool `md:"get_all"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"content": o.Content,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	o.Content = values["content"]

	return nil
}
func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code": r.Code,
		"data": r.Data,
	}
}

func (r *Reply) FromMap(values map[string]interface{}) error {
	var err error
	r.Code, err = coerce.ToInt(values["code"])
	if err != nil {
		return err
	}
	r.Data = values["data"]

	return nil
}
