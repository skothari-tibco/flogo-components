package imagerunner

import (
	"fmt"
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}
func TestSettings(t *testing.T) {
	fmt.Println("Starting Test")
	settings := &Settings{Config: nil, Host: nil, Networkconfig: nil}
	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	fmt.Println(err)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("imagename", "my-app:latest")
	//eval
	_, err = act.Eval(tc)
	fmt.Println(tc.GetOutput("code"))
	//val := tc.GetOutput("result")
	assert.Equal(t, 0, tc.GetOutput("code"))
}
