package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func TestCreate(t *testing.T) {
	cli, err := client.New(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	resp, err := cli.ImagePull(context.Background(), "docker.io/library/alpine", client.ImagePullOptions{})

	if err != nil {
		panic(err)

	}
	defer resp.Close()
	io.Copy(os.Stdout, resp)

	fmt.Println("Pulled image successfully")
	res, err := cli.ContainerCreate(context.Background(), client.ContainerCreateOptions{
		Image: "alpine",
		Config: &container.Config{
			Cmd: []string{"echo", "hello world"},
			Tty: false,
		},
	})
	fmt.Printf("Created container with ID: %s\n", res.ID)
	if err != nil {
		panic(err)
	}
	if _, err := cli.ContainerStart(context.Background(), res.ID, client.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	waitResult := cli.ContainerWait(context.Background(), res.ID, client.ContainerWaitOptions{
		Condition: container.WaitConditionNotRunning,
	})
	select {
	case err := <-waitResult.Error:
		if err != nil {
			panic(err)
		}
	case status := <-waitResult.Result:
		fmt.Printf("Container exited with status code: %d\n", status.StatusCode)
	}
	logRes, _ := cli.ContainerLogs(context.Background(), res.ID, client.ContainerLogsOptions{
		ShowStdout: true,
	})
	stdcopy.StdCopy(os.Stdout, os.Stderr, logRes)

}
