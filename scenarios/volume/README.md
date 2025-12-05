# Volume 模块记录

该模块包含两个阶段：

1. `fill/`：创建带有 32 MiB `tmpfs` 限额的 Volume，并持续向 `/demo-data` 写入数据，实时输出“累计写入/已用/剩余”。当卷空间耗尽时，容器会以退出码 `42` 结束，并打印 `df` 结果，用来观察满盘后的行为。
2. `expand/`：重新创建同名 Volume，将容量扩展到 96 MiB，再次写入 64 MiB 数据，确认扩容后写入可成功完成。

## 运行方式

```bash
# 写满并观察剩余空间
GO111MODULE=on go run ./scenarios/volume/fill

# 扩容并验证可继续写入
GO111MODULE=on go run ./scenarios/volume/expand
```

## 预期现象

- `fill` 的日志会不断打印 `累计写入=<N>MiB 已用=<X>MiB 剩余=<Y>MiB`，最终出现 `写入失败：卷空间已耗尽`，对应的容器退出码非 0（预期 42）。
- `expand` 会输出 `完成 64MiB 写入，卷可继续使用`，容器退出码为 0，证明扩容后的卷能正常工作。

## 结果记录

- 最近一次 `fill`：**待运行**（运行后请把关键日志粘贴在这里，方便回溯）。
- 最近一次 `expand`：**待运行**。
