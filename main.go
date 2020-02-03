package main

import (
	"github.com/ermyril/dockyard/config"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"time"
	//"github.com/davecgh/go-spew/spew"
)

type Dockyard struct {
	config config.Config
	workdir string
	dbContainerID string
	dockerClient *client.Client
}

func main() {

	yard := GetDockyardClient(config.GetConfig())

	//spew.Printf("yard: %v", yard)


	yard.backup()

}

func GetDockyardClient(config config.Config) Dockyard {

	dockyard := Dockyard{}

	dockyard.config = config
	dockyard.workdir = getCurrentDirectory()
	dockyard.dockerClient = getDockerClient()

	dockyard.fetchDatabaseContainer()

	return dockyard
}


func getDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return cli
}


func (yard *Dockyard) backup() {

	optionsCreate := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: false,
		Cmd:          []string{"mysqldump", "-u", yard.config.User, "-p"+yard.config.Password, yard.config.Database},
	}

	rst, err := yard.dockerClient.ContainerExecCreate(
		context.Background(),
		yard.dbContainerID,
		optionsCreate)
	if err != nil {
		panic(err)
	}

	response, err := yard.dockerClient.ContainerExecAttach(context.Background(), rst.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}
	defer response.Close()

	fmt.Println(response)

	data, _ := ioutil.ReadAll(response.Reader)

	if (len(data) == 0) {
		fmt.Println("Error: no output from the database")
	}


	fmt.Println(time.Now().Format("02-01-2006_15:04"))
	fmt.Println(yard.config.Directory)



	dumpName := yard.config.Database + time.Now().Format("_02-01-06_1504.sql")
	dumpPath := filepath.Join(yard.config.Directory, dumpName)

	fmt.Println(dumpPath)

	dump, err := os.Create(dumpPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	length, err := dump.Write(data)
	if err != nil {
		fmt.Println(err)
		dump.Close()
		return
	}
	fmt.Println(length, "bytes written successfully")
	err = dump.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (yard *Dockyard) fetchDatabaseContainer() {
	yard.dbContainerID = getDatabaseContainerId(getContainerList(), yard.workdir, yard.config.Host)
}

func getContainerList() []types.Container {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	//containers, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		panic(err)
	}

	return containers
}

func getDatabaseContainerId(containers []types.Container, workdir string, service string) string {

	for _, container := range containers {
		if container.Labels["com.docker.compose.project.working_dir"] == workdir &&
			container.Labels["com.docker.compose.service"] == service {
			return container.ID
		}
	}

	panic("container not found")
}

func getCurrentDirectory() string {
	workdir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	return workdir
}
