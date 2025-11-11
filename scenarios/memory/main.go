package main

import (
	"context"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
	"log"
	"os"
	"time"

	"github.com/moby/moby/api/types/container"
)

const (
	Mbyte            = 1024 * 1024
	demoImage        = "mem-test"
	memoryLimitBytes = 64 * Mbyte
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
	// 创建容器
	res, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image: demoImage,
		Config: &container.Config{
			Tty: false,
		},
		HostConfig: &container.HostConfig{
			Resources: container.Resources{
				Memory: memoryLimitBytes,
				// swap = 0
				MemorySwap: memoryLimitBytes,
			},
		},
		Name: "test-mem",
	})
	if err != nil {
		panic(err)
	}
	if _, err := cli.ContainerStart(context.Background(), res.ID, client.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	logRes, _ := cli.ContainerLogs(context.Background(), res.ID, client.ContainerLogsOptions{
		ShowStdout: true,
	})
	if _, err = stdcopy.StdCopy(os.Stdout, os.Stderr, logRes); err != nil {
		panic(err)
	}
}
