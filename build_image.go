package main

import (
	"archive/tar"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err, " :unable to init client")
	}

	destinationfile := "tarball.tar"
	sourcedir := "/home/shenshengkafei/swagger/src/IO.Swagger"

	dir, err := os.Open(sourcedir)
	defer dir.Close()

	// get list of files
	files, err := dir.Readdir(0)

	// create tar file
	tarfile, err := os.Create(destinationfile)
	defer tarfile.Close()

	var fileWriter io.WriteCloser = tarfile

	tarfileWriter := tar.NewWriter(fileWriter)
	defer tarfileWriter.Close()

	for _, fileInfo := range files {

		if fileInfo.IsDir() {
			continue
		}

		file, err := os.Open(dir.Name() + string(filepath.Separator) + fileInfo.Name())
		if err != nil {
			log.Fatal(err, " :unable to open file in source directory.")
		}
		defer file.Close()

		// prepare the tar header
		header := new(tar.Header)
		header.Name = fileInfo.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()

		err = tarfileWriter.WriteHeader(header)

		_, err = io.Copy(tarfileWriter, file)
	}

	var fileReader io.ReadCloser = tarfile

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		fileReader,
		types.ImageBuildOptions{
			Tags:       []string{"shenshengkafei/imagename"},
			Context:    fileReader,
			Dockerfile: destinationfile,
			Remove:     true})
	if err != nil {
		log.Fatal(err, " :unable to build docker image")
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}

	auth := types.AuthConfig{
		Username: "shenshengkafei",
		Password: "SudoPassword321",
	}
	authBytes, _ := json.Marshal(auth)
	authBase64 := base64.URLEncoding.EncodeToString(authBytes)

	imagePushResponse, err := cli.ImagePush(
		context.Background(),
		"shenshengkafei/imagename",
		types.ImagePushOptions{
			All:          true,
			RegistryAuth: authBase64,
		})
	if err != nil {
		log.Fatal(err, " :unable to push docker image")
	}
	defer imagePushResponse.Close()
	_, err = io.Copy(os.Stdout, imagePushResponse)
	if err != nil {
		log.Fatal(err, " :unable to read image push response")
	}
}
