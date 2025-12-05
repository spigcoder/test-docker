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
	demoImage        = "docker.io/library/python:3.12-alpine"
	volumeName       = "volume-limit-demo"
	volumeMountPath  = "/demo-data"
	volumeLimitBytes = 32 * scenario.MiB
	cpuLimitNano     = 1_000_000_000
	memoryLimitBytes = 128 * scenario.MiB
	rootFsLimitBytes = 512 * scenario.MiB
	chunkMiB         = 4
)

const volumeFillScript = `set -euo pipefail
TARGET="%s"
CHUNK_MB=%d
TOTAL=0
rm -f "$TARGET/fillfile"
touch "$TARGET/fillfile"
while true; do
    if dd if=/dev/zero of="$TARGET/fillfile" bs=1M count="$CHUNK_MB" oflag=append conv=notrunc status=none; then
        TOTAL=$((TOTAL+CHUNK_MB))
        DF_LINE=$(df -m "$TARGET" | tail -1)
        USED=$(echo "$DF_LINE" | awk '{print $3}')
        AVAIL=$(echo "$DF_LINE" | awk '{print $4}')
        echo "累计写入=${TOTAL}MiB 已用=${USED}MiB 剩余=${AVAIL}MiB"
        sync
    else
        echo "写入失败：卷空间已耗尽" >&2
        df -m "$TARGET"
        exit 42
    fi
    sleep 0.1
done`

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

	if err := scenario.RecreateTmpfsVolume(ctx, cli, volumeName, volumeLimitBytes); err != nil {
		log.Fatalf("创建受限 volume 失败: %v", err)
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

	script := fmt.Sprintf(volumeFillScript, volumeMountPath, chunkMiB)
	result, err := scenario.RunControlledContainer(ctx, cli, &container.Config{
		Image: demoImage,
		Cmd:   []string{"sh", "-c", script},
	}, hostConfig, "volume-fill")
	if err != nil {
		log.Fatalf("执行写满测试失败: %v", err)
	}

	scenario.LogRunResult("volume fill", result)
	if result.StatusCode == 0 {
		log.Fatalf("期望写入失败以确认空间上限，但容器以 0 退出")
	}
}
