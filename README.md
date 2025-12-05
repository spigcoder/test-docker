# Docker 资源控制模块化示例

本项目使用 `github.com/moby/moby` SDK，围绕“创建受限容器并观测资源耗尽行为”拆分成多个独立模块，每个模块位于 `scenarios/<资源类型>` 目录，并自带 Markdown 记录运行结果。核心能力：

| 模块目录 | 子命令 | 功能简介 |
| --- | --- | --- |
| `scenarios/volume` | `fill`, `expand` | 受限数据盘写满、扩容后再写入 |
| `scenarios/memory` | `pressure` | 分配内存直至 `MemoryError`/OOM |
| `scenarios/cpu` | `limit` | 读取 cgroup 配额并计算 1 vCPU 的实际使用率 |
| `scenarios/rootfs` | `fill` | 利用 `StorageOpt[\"size\"]` 写满系统盘（依赖驱动支持） |

公共逻辑（镜像拉取、容器运行、日志收集、Volume 复建等）被收敛到 `internal/scenario` 包，方便在不同模块之间复用。

## 环境要求

- 可访问的 Docker Engine，推荐 24.0+。
- 能拉取 `docker.io/library/python:3.12-alpine`。
- Go 1.21 及以上。
- 如果要验证系统盘限额，请使用支持 `StorageOpt[\"size\"]` 的存储驱动（devicemapper / btrfs / zfs 等）。

## 运行方式

所有命令都通过 `go run` 直接执行，示例：

```bash
# Volume 写满
GO111MODULE=on go run ./scenarios/volume/fill

# Volume 扩容验证
GO111MODULE=on go run ./scenarios/volume/expand

# 内存压测
GO111MODULE=on go run ./scenarios/memory/pressure

# CPU 限额探测
GO111MODULE=on go run ./scenarios/cpu/limit

# 系统盘写满
GO111MODULE=on go run ./scenarios/rootfs/fill
```

执行完毕后，请在对应模块目录的 `README.md` 中补充“结果记录”段落，形成可追溯的实验报告。

## 模块要点

- **Volume 模块**：`fill` 以 32 MiB `tmpfs` Volume 为例，循环写入并实时输出 `累计写入/已用/剩余`，观察满盘时的 `dd` 报错；`expand` 重建卷为 96 MiB，验证扩容后 64 MiB 写入可以成功完成。
- **Memory 模块**：容器内脚本每次分配 8 MiB，直到命中内存上限。日志中可看到最高分配的 MiB，退出码 23 或 137 均表示限制生效。
- **CPU 模块**：读取 cgroup 配额（`cpu.max` 或 `cpu.cfs_*`），跑 6 秒忙循环后根据 `cpu.stat` 计算平均 CPU 使用率，应接近 1.00 vCPU。
- **RootFS 模块**：通过 `StorageOpt["size"]=128m` 约束根文件系统，对 `/root/system-fill.bin` 进行写入。如果驱动支持，会在若干次写入后报错退出码 55；否则程序会提示未触发限额，需根据宿主机环境调整。

## 目录结构

```
internal/scenario/  # Docker 客户端、运行与日志采集的通用封装
scenarios/volume/   # 数据卷相关脚本 + README
scenarios/memory/   # 内存压测脚本 + README
scenarios/cpu/      # CPU 限额验证脚本 + README
scenarios/rootfs/   # 系统盘写满脚本 + README
```

每个 README 都包含“运行方式 / 预期现象 / 结果记录”，方便记录多次实验的对比结论。

## 注意事项

- Volume 场景使用 `tmpfs` 驱动，因此占用宿主机内存；如需真实磁盘，可换成具有 `size` 选项的驱动或外部块设备。
- RootFS 限额依赖存储驱动实现，若使用 `overlay2` 可能无法复现写满行为，此时请结合宿主机实际方案（例如 LVM loop 设备）。
- 所有模块默认限制 1 vCPU、128 MiB 左右内存，可在各自 `main.go` 的常量中调整。

通过这些模块，可以分别、清晰地验证 CPU、内存、系统盘、数据盘（Volume）的资源限制及扩容策略，为后续自动化或容量评估提供直接的脚本参考。
