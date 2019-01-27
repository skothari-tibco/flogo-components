package imagerunner

import (
	"bytes"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"golang.org/x/net/context"
)

func init() {
	activity.Register(&Activity{}, New)
}

type Activity struct {
	settings *Settings
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	act := &Activity{settings: s}
	fmt.Println("Registering Settings")
	return act, nil
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	output := &Output{}
	ctx.GetInputObject(input)

	ctx.Logger().Info("Creating the container for image ", input.ImageName)
	if input.ImageName != "" {
		a.settings.Config = &container.Config{Image: input.ImageName}
	}

	bctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
		return true, err
	}

	resp, err := cli.ContainerCreate(bctx, a.settings.Config, a.settings.Host, a.settings.Networkconfig, "")
	if err != nil {
		return true, err
	}
	ctx.Logger().Info("Executing the container...")
	if err := cli.ContainerStart(bctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return true, err
	}
	statusCh, err := cli.ContainerWait(bctx, resp.ID)
	if err != nil {
		return true, err
	}
	output.Code, _ = coerce.ToInt(statusCh)
	out, err := cli.ContainerLogs(bctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})

	if err != nil {
		return true, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	s := buf.String()
	ctx.Logger().Debugf(s)
	output.Logs = s

	ctx.Logger().Info("Completed running the container with exit code ", statusCh)
	ctx.Logger().Info("Check debug logs for more information...")

	ctx.SetOutputObject(output)

	return true, nil
}
