package shell

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

	act := &Activity{}

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("command", "echo hello")
	//eval
	_, err := act.Eval(tc)

	//val := tc.GetOutput("result")
	assert.Nil(t, err)
}
