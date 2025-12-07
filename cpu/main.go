package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"

	"github.com/moby/moby/api/types/container"
)

const (
	MiB       = 1024 * 1024
	demoImage = "stress"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if _, err := os.Open("./name"); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cli, err := client.New(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	name := "cpu-test" + time.Now().Format("150405")
	res, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image: demoImage,
		Config: &container.Config{
			Tty: false,
			Cmd: []string{"sleep", "3600"},
		},
		HostConfig: &container.HostConfig{
			Resources: container.Resources{
				CPUPercent: 100000,
				CPUQuota:   100000,
			},
		},
		Name: name,
	})

	if err != nil {
		log.Fatalf("执行 CPU 限额探测失败: %v", err)
	}
	if _, err := cli.ContainerStart(context.Background(), res.ID, client.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	logRes, _ := cli.ContainerLogs(context.Background(), res.ID, client.ContainerLogsOptions{
		ShowStdout: true,
	})
	stdcopy.StdCopy(os.Stdout, os.Stderr, logRes)
}
