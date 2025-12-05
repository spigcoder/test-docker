package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"

	"test-docker/internal/scenario"
)

const (
	demoImage           = "docker.io/library/python:3.12-alpine"
	volumeName          = "volume-limit-demo"
	volumeMountPath     = "/demo-data"
	expandedVolumeBytes = 96 * scenario.MiB
	cpuLimitNano        = 1_000_000_000
	memoryLimitBytes    = 128 * scenario.MiB
	rootFsLimitBytes    = 512 * scenario.MiB
	writeCountMiB       = 64
)

const expansionScript = `set -euo pipefail
TARGET="%s"
COUNT_MB=%d
rm -f "$TARGET/expanded.bin"
dd if=/dev/zero of="$TARGET/expanded.bin" bs=1M count="$COUNT_MB" status=none
sync
echo "完成 ${COUNT_MB}MiB 写入，卷可继续使用"`

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cli, err := scenario.NewDockerClient()
	if err != nil {
		log.Fatalf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	if err := scenario.PullImage(ctx, cli, demoImage); err != nil {
		log.Fatalf("拉取镜像失败: %v", err)
	}

	if err := scenario.RecreateTmpfsVolume(ctx, cli, volumeName, expandedVolumeBytes); err != nil {
		log.Fatalf("扩容 volume 失败: %v", err)
	}

	mounts := []mount.Mount{{
		Type:   mount.TypeVolume,
		Source: volumeName,
		Target: volumeMountPath,
	}}

	hostConfig := scenario.BuildHostConfig(container.Resources{
		NanoCPUs: cpuLimitNano,
		Memory:   memoryLimitBytes,
	}, rootFsLimitBytes, mounts)

	script := fmt.Sprintf(expansionScript, volumeMountPath, writeCountMiB)
	result, err := scenario.RunControlledContainer(ctx, cli, &container.Config{
		Image: demoImage,
		Cmd:   []string{"sh", "-c", script},
	}, hostConfig, "volume-expand")
	if err != nil {
		log.Fatalf("执行扩容验证失败: %v", err)
	}

	scenario.LogRunResult("volume expand", result)
	if result.StatusCode != 0 {
		log.Fatalf("扩容后的写入应成功，但容器退出码为 %d", result.StatusCode)
	}
}
