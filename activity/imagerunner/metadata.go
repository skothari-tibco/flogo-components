package imagerunner

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Config        *container.Config         `md:"config"`
	Host          *container.HostConfig     `md:"host"`
	Networkconfig *network.NetworkingConfig `md:"networkconfig"`
}

type Input struct {
	Image string `md:"image"`
}

type Output struct {
	Logs string `md:"logs"`
	Code int64  `md:"code"`
}

func (o *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"image": o.Image,
	}
}

func (o *Input) FromMap(values map[string]interface{}) error {

	var err error
	o.Image, err = coerce.ToString(values["image"])
	if err != nil {
		return err
	}

	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"logs": o.Logs,
		"code": o.Code,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Logs, err = coerce.ToString(values["logs"])
	if err != nil {
		return err
	}
	o.Code, err = coerce.ToInt64(values["code"])
	if err != nil {
		return err
	}

	return nil
}
