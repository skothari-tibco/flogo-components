package shell

import (
	"os/exec"
	"strings"

	"github.com/project-flogo/core/activity"
)

func init() {
	activity.Register(&Activity{}, New)
}

type Activity struct {
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

var activityMd = activity.ToMetadata()

func New(ctx activity.InitContext) (activity.Activity, error) {
	act := &Activity{}
	return act, nil
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	output := &Output{}

	output = nil
	ctx.GetInputObject(input)

	cmds := strings.Split(input.Command, " ")

	result, err := exec.Command(cmds[0], cmds[1:]...).Output()
	output.Result = result
	output.Error = err

	ctx.SetOutputObject(output)

	return true, nil
}
