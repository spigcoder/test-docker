package scenario

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
)

const (
// MiB 提供统一的单位换算，方便在各模块中设置资源额度。
)

// NewDockerClient 复用环境变量并开启 API 版本协商。
func NewDockerClient() (*client.Client, error) {
	return client.New(client.FromEnv, client.WithAPIVersionNegotiation())
}

// PullImage 确保镜像可用，返回前会把拉取输出全部读取完毕。
func PullImage(ctx context.Context, cli *client.Client, ref string) error {
	log.Printf("拉取镜像 %s", ref)
	resp, err := cli.ImagePull(ctx, ref, client.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer resp.Close()
	_, err = io.Copy(io.Discard, resp)
	return err
}

// RunResult 记录一次容器运行的核心信息。
type RunResult struct {
	ContainerID string
	StatusCode  int64
	Logs        string
}

// RunControlledContainer 创建并等待容器退出，同时采集日志，方便各模块重用。
func RunControlledContainer(ctx context.Context, cli *client.Client, config *container.Config, hostConfig *container.HostConfig, namePrefix string) (RunResult, error) {
	containerName := fmt.Sprintf("%s-%d", namePrefix, time.Now().UnixNano())
	config.AttachStdout = false
	config.AttachStderr = false
	config.Tty = false

	resp, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config:     config,
		HostConfig: hostConfig,
		Name:       containerName,
	})
	if err != nil {
		return RunResult{}, fmt.Errorf("create container: %w", err)
	}

	cleanup := func(force bool) {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_, _ = cli.ContainerRemove(cleanupCtx, resp.ID, client.ContainerRemoveOptions{Force: force, RemoveVolumes: false})
	}
	defer cleanup(true)

	if _, err = cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		return RunResult{}, fmt.Errorf("start container: %w", err)
	}

	wait := cli.ContainerWait(ctx, resp.ID, client.ContainerWaitOptions{Condition: container.WaitConditionNotRunning})
	var status container.WaitResponse
	select {
	case err = <-wait.Error:
		if err != nil {
			return RunResult{}, fmt.Errorf("wait for container: %w", err)
		}
	case status = <-wait.Result:
	}

	logReader, err := cli.ContainerLogs(ctx, resp.ID, client.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return RunResult{}, fmt.Errorf("collect logs: %w", err)
	}
	defer logReader.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, logReader); err != nil {
		return RunResult{}, fmt.Errorf("demultiplex logs: %w", err)
	}

	logs := strings.Builder{}
	if stdoutBuf.Len() > 0 {
		fmt.Fprintf(&logs, "STDOUT:\n%s\n", strings.TrimSpace(stdoutBuf.String()))
	}
	if stderrBuf.Len() > 0 {
		fmt.Fprintf(&logs, "STDERR:\n%s\n", strings.TrimSpace(stderrBuf.String()))
	}

	return RunResult{
		ContainerID: resp.ID,
		StatusCode:  status.StatusCode,
		Logs:        strings.TrimSpace(logs.String()),
	}, nil
}

// LogRunResult 在统一格式下打印容器结论及日志片段。
func LogRunResult(label string, result RunResult) {
	log.Printf("%s container(%s) exit=%d", label, ShortID(result.ContainerID), result.StatusCode)
	if result.Logs != "" {
		log.Printf("%s logs:\n%s", label, result.Logs)
	}
}

// ShortID 返回 Docker ID 的前 12 位，便于日志阅读。
func ShortID(id string) string {
	if len(id) <= 12 {
		return id
	}
	return id[:12]
}

// HumanBytes 以 GiB/MiB/KiB 可读形式展示字节数。
func HumanBytes(v int64) string {
	const (
		kiB = 1024
		miB = kiB * 1024
		giB = miB * 1024
	)
	switch {
	case v >= giB:
		return fmt.Sprintf("%.2f GiB", float64(v)/float64(giB))
	case v >= miB:
		return fmt.Sprintf("%.2f MiB", float64(v)/float64(miB))
	case v >= kiB:
		return fmt.Sprintf("%.2f KiB", float64(v)/float64(kiB))
	default:
		return fmt.Sprintf("%d B", v)
	}
}

// BuildHostConfig 根据资源与挂载选项构造 HostConfig，统一设置根文件系统大小。
func BuildHostConfig(resources container.Resources, rootFsLimitBytes int64, mounts []mount.Mount) *container.HostConfig {

}

// RecreateTmpfsVolume 使用 tmpfs 驱动重新创建受限容量的 volume。
func RecreateTmpfsVolume(ctx context.Context, cli *client.Client, name string, sizeBytes int64) error {
	_, err := cli.VolumeRemove(ctx, name, client.VolumeRemoveOptions{Force: true})
	if err != nil && !cerrdefs.IsNotFound(err) {
		return err
	}
	opts := client.VolumeCreateOptions{
		Name:   name,
		Driver: "local",
		DriverOpts: map[string]string{
			"type":   "tmpfs",
			"device": "tmpfs",
			"o":      fmt.Sprintf("size=%dm", sizeBytes/MiB),
		},
		Labels: map[string]string{
			"scenario": "resource-limit",
		},
	}
	_, err = cli.VolumeCreate(ctx, opts)
	return err
}
