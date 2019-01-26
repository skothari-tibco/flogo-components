package twillio

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
	tc.SetInput("accountSID", "AC8d7eab76db20a7d02b586dae3d421419")
	tc.SetInput("authtoken", "d1978282544a1b4dbb8fe4bdabcac39c")
	tc.SetInput("from", "+12568576705")
	tc.SetInput("to", "+16693509593")
	tc.SetInput("message", "Hello ")
	//eval
	out, err := act.Eval(tc)

	fmt.Println(out)

	//val := tc.GetOutput("result")
	assert.Nil(t, err)
}
