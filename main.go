package main

import (
	"github.com/ermyril/dockyard/config"

	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"os"
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


	//yard.Backup()
	yard.Restore()
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
