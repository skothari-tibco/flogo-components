package imagebuilder

import (
	"github.com/project-flogo/core/data/coerce"
)

type Input struct {
	DockerPath string `md:"dockerpath"`
	ImageName  string `md:imagename`
}
type Output struct {
	ImageName string `md:"imagename"`
}

func (o *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"dockerpath": o.DockerPath,
		"imagename":  o.ImageName,
	}
}

func (o *Input) FromMap(values map[string]interface{}) error {

	var err error
	o.DockerPath, err = coerce.ToString(values["dockerpath"])
	if err != nil {
		return err
	}
	o.ImageName, err = coerce.ToString(values["imagename"])
	if err != nil {
		return err
	}

	return nil
}

func (r *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"imagename": r.ImageName,
	}
}

func (r *Output) FromMap(values map[string]interface{}) error {

	var err error
	r.ImageName, err = coerce.ToString(values["imagename"])
	if err != nil {
		return err
	}
	return nil
}
