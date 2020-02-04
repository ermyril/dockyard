package main

import (

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"time"
	//"github.com/davecgh/go-spew/spew"
)

func (yard *Dockyard) Backup() {

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

	response, err := yard.dockerClient.ContainerExecAttach(
		context.Background(),
		rst.ID,
		types.ExecStartCheck{})

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
