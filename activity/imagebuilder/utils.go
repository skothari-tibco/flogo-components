package imagebuilder

import (
	"archive/tar"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

func getDockerTar(path string) (*bytes.Reader, error) {

	dockerfile, err := getDockerFile(path)

	if err != nil {
		return nil, err
	}

	dockerBuildContext, err := os.Open(dockerfile)
	defer dockerBuildContext.Close()

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	readDockerFile, err := ioutil.ReadAll(dockerBuildContext)
	if err != nil {
		return nil, err
	}

	tarHeader := &tar.Header{
		Name: "Dockerfile",
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return nil, err
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

func getDockerFile(path string) (string, error) {

	file, err := os.Stat(path)

	if err != nil {
		return "", err
	}
	switch mode := file.Mode(); {

	case mode.IsDir():

		files, err := ioutil.ReadDir(path)
		if err != nil {
			return "", nil
		}

		for _, f := range files {
			if f.Name() == "Dockerfile" {
				return filepath.Join(path, f.Name()), nil
			}
		}

	case mode.IsRegular():
		return path, nil
	}
	return "", nil
}
