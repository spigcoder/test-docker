package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"

	"github.com/moby/moby/api/types/container"
)

const (
	demoImage        = "docker.io/library/alpine"
	rootFsLimitBytes = 128
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
	//拉去镜像
	resp, err := cli.ImagePull(ctx, demoImage, client.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, resp)
	fmt.Println("Pulled image successfully")
	defer resp.Close()
	// 设置 --storage-opt size=10M 限制可写层的大小
	hostConfig := &container.HostConfig{}
	hostConfig.StorageOpt = map[string]string{
		"size": fmt.Sprintf("%dM", rootFsLimitBytes),
	}
	// 创建容器
	res, err := cli.ContainerCreate(context.Background(), client.ContainerCreateOptions{
		Image: "alpine",
		Config: &container.Config{
			Tty: false,
			Cmd: []string{"sleep", "3600"},
		},
		HostConfig: hostConfig,
		Name:       "test-ds",
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
	stdcopy.StdCopy(os.Stdout, os.Stderr, logRes)
}
