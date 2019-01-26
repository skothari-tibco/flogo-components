package twillio

import (
	"github.com/project-flogo/core/activity"
	"github.com/sfreiberg/gotwilio"
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
	out := &Output{}
	ctx.GetInputObject(input)

	accountSID := input.AccountSID
	authToken := input.AuthToken
	from := input.From
	to := input.To
	message := input.Message

	twilio := gotwilio.NewTwilioClient(accountSID, authToken)

	resp, _, err := twilio.SendSMS(from, to, message, "", "")

	if err != nil {
		ctx.Logger().Debugf("Error sending SMS:", err)
		return false, err
	}
	out.Result = resp.Status

	ctx.SetOutputObject(out)
	ctx.Logger().Infof("Response:", resp)

	return true, nil
}
