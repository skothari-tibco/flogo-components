package imagebuilder

import (
	"context"
	"errors"

	"bytes"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/project-flogo/core/activity"
)

func init() {
	activity.Register(&Activity{})
}

type Activity struct {
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	ctx.GetInputObject(input)

	//dockerFileTarReader, err := getDockerTar(input.DockerPath)

	ctx.Logger().Info(" Starting to build the image")

	reader, err := archive.TarWithOptions(input.DockerPath, &archive.TarOptions{})

	cli, _ := client.NewEnvClient()

	if input.ImageName == "" {
		ctx.Logger().Info("Image name not specified", input)
		return true, errors.New("Image name not specified")
	}

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		//Context:    dockerFileTarReader,
		Remove: true,
		Tags:   []string{input.ImageName},
	}

	buildResp, err := cli.ImageBuild(context.Background(), reader, buildOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer buildResp.Body.Close()
	ctx.Logger().Info("Building the Image...")
	buf := new(bytes.Buffer)

	buf.ReadFrom(buildResp.Body)

	ctx.Logger().Debug(buf.String())
	result := buf.String()
	streams := strings.Split(result, "\n")

	if strings.Contains(streams[len(streams)-2], "Successfully") {
		stream := strings.Split(streams[len(streams)-2], " ")

		imageID := strings.Replace(stream[len(stream)-1], "\\n\"}", "", -1)

		ctx.Logger().Info("Image id is ", imageID)

		ctx.SetOutput("imagename", input.ImageName)
		return true, nil
	}

	ctx.Logger().Info("Error in building image")

	return true, errors.New("Error in building image refer logs")

	/*
		_, err = io.Copy(os.Stdout, buildResp.Body)
		if err != nil {
			log.Fatal(err, " :unable to read image build response")
		}*/

}
