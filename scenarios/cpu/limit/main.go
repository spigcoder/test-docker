package main

import (
	"context"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
	"io"
	"log"
	"os"
	"time"

	"github.com/moby/moby/api/types/container"
)

const (
	MiB          = 1024 * 1024
	demoImage    = "docker.io/library/python:3.12-alpine"
	cpuLimitNano = 1_000_000_000 // 1 vCPU
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cli, err := client.New(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	resp, err := cli.ImagePull(ctx, demoImage, client.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer resp.Close()
	_, err = io.Copy(io.Discard, resp)

	res, err := cli.ContainerCreate(context.Background(), client.ContainerCreateOptions{
		Image: demoImage,
		Config: &container.Config{
			Tty: false,
			Cmd: []string{"sleep", "3600"},
		},
		HostConfig: &container.HostConfig{
			Resources: container.Resources{
				CPUPercent: 100000,
				CPUQuota:   200000,
			},
		},
		Name: "test-ds",
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
