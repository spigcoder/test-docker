此项目作为我们 Swan-Lab 中 go web 项目的通用模版，这里介绍一下各个部分的功能

## 项目结构

``` shell
.
├── api
│   └── v1
│       └── user
│           └── user.go
├── build
├── cmd
│   └── main.go
├── configs
├── docs
├── internal
│   ├── handler
│   │   └── user
│   │       ├── profile.go
│   │       └── profile_test.go
│   ├── pkg
│   │   ├── errs
│   │   ├── middleware
│   │   ├── mock
│   │   └── model
│   ├── repo
│   │   └── user
│   │       └── profile.go
│   └── router
│       ├── index.go
│       └── router.go
├── pkg
├── .env
├── .golangci.yaml
├── go.mod
└── README.md
```

* api：用来存放 http 中的 response request 的定义，或者 grpc 中的 proto 文件。
* cmd：可执行程序的入口，一般仅存放少量的初始化代码，如果有多个 main 文件应该用不同的目录分开，目录名为 app 名字。
* configs：一些公用的配置文件可以存储在这里。
* internal：不可被外部引用的代码，主要用来存放核心业务。
  * handler：直接操作 web 部分。
    * handler 这里我是按照路径来进行划分，比如说 url 中 https://host:port/api/user/profile 下面的操作就对应 /user/profile，这么做的好处是对于大型的项目来说，第一结构清晰，第二不至于单个文件过大。
    * 这里还有其他的划分方法。
      * 第一种：所有的 /user/profile、/user/auth 这些操作全部放到 /handler/user.go 中，这样做会导致这个文件随着项目的扩展越来越大，然后很容易 git 冲突，因为都改的是一个文件。
      * 第二种：/handler/user_profile 这种是 go 标准库常见的做法，但是全部放到 /handler 中我个人认为会导致这个文件很大，然后不好查找和扩展。
    * /handler/user/profile 有一个缺点就是我们 /user/profile.go 的包名是 user，/repo/user/profile 的包名也是 user，如果要同时使用就要重命名包名，所以如果对于小型项目，选择 /handler/user_profile 也是不错的选择。
  * repo：与数据库、缓存操作相关，repo 我这里的逻辑和 /handler 一样，如果同时含有数据库和缓存，可以在 repo/dao 中进行数据库操作，在 repo/cache 中进行缓存操作。
  * router：handler 初始化和路由注册操作。
    * index：进行路由的初始化。
    * router：定义各个路由的处理函数。
  * pkg: internal 中的 pkg 定义的是不想被外部使用，但是需要共享的代码。
    * errs：errs 中定义了 http code 和我们错误 code 的一些映射，因为我们日常使用的时候，database 出错可以是 500，redis 出错也可以是 500，所以这里进行映射的同时使用 middleware 中的中间件来进行错误打印，这样就不用每次都进行打印操作了。
* pkg：可以被外部引用的代码，如果有在 internal 中使用的公共代码，可以在 internal 下再创建一个 pkg 目录。
* test：外部测试，用来测试接口而非实现。

这里在使用的时候没有采用 controller + service 的结构，而是在 handler 直接对接 repo，我这么设计的主要目的是因为对于不是特别复杂的项目来说，如果使用 controller + service 的话，几乎所有的内容都是在 service 完成的，controller 就只是简单的进行函数的调用和一些校验的工作，所以这里可以根据实际的项目情况来进行考虑。

困惑的地方可以看一下其中的代码demo。

## .golangci.yaml

以下是当前 .golangci.yaml 的注释版本，可以根据实际项目的需要进行酌情更改

``` yaml
version: "2" # 配置文件版本，目前标准为 2

run:
  modules-download-mode: readonly # 运行时不修改 go.mod 文件，只读模式

linters:
  default: none # 禁用所有默认开启的 linter，改为白名单模式（手动指定 enable）
  enable:
    - bodyclose     # 检查 HTTP 响应体 (body) 是否被正确关闭
    - dogsled       # 检查是否使用了过多的空白标识符 (如: _, _, _ = func())
    - durationcheck # 检查两个不同类型的 time.Duration 是否在做乘法
    - errcheck      # 检查是否忽略了 error 返回值
    - goconst       # 检查代码中重复使用的字符串，建议提取为常量
    - gocyclo       # 检查函数的圈复杂度 (代码逻辑是否过于复杂)
    - govet         # Go 官方的静态分析工具
    - ineffassign   # 检查无效的赋值 (赋值后未使用，或被立刻覆盖)
    - lll           # 检查单行代码长度是否超标
    - misspell      # 检查单词拼写错误 (主要针对注释和文档)
    - mnd           # 检查代码中的“魔术数字” (Magic Numbers)
    - prealloc      # 检查切片 (slice) 是否可以预分配内存以优化性能
    - revive        # Golint 的替代品，检查代码风格和规范
    - staticcheck   # 高级的静态代码分析工具
    - unconvert     # 检查不必要的类型转换
    - unused        # 检查未使用的常量、变量、函数或类型
    - wastedassign  # 检查无用的赋值操作
    - whitespace    # 检查函数声明和调用时的首尾是否有不必要的空行
  
  settings: # 针对上述 linter 的具体配置参数
    gocyclo:
      min-complexity: 50 # 圈复杂度阈值，超过 50 报错 (默认通常是 30，50 比较宽容)
    govet:
      enable:
        - shadow # 开启变量遮蔽检查 (内部变量名覆盖了外部变量名)
    lll:
      line-length: 160 # 单行最大长度限制为 160 字符
    misspell:
      locale: US # 使用美式英语拼写规则
    mnd:
      checks: # 指定检查魔术数字的场景
        - case      # switch-case 语句中
        - condition # if/for 判断条件中
        - return    # return 返回值中
    whitespace:
      multi-func: true # 允许在多行函数签名中使用空行

  exclusions: # 排除规则配置
    generated: lax # 对生成的代码 (generated code) 放宽检查标准
    presets: # 排除某些预设类别的误报
      - comments             # 排除注释相关的特定问题
      - common-false-positives # 排除常见的误报
      - legacy               # 排除旧版遗留规则
      - std-error-handling   # 排除标准错误处理相关的特定模式
    rules: # 针对特定文件的特殊规则
      - linters:
          - goconst # 在测试文件中忽略字符串重复检查
        path: (.+)_test\.go # 正则匹配所有 _test.go 文件
    paths: # 全局忽略以下路径的检查
      - third_party$ # 忽略 third_party 目录
      - builtin$     # 忽略 builtin 目录
      - examples$    # 忽略 examples 目录

formatters: # 格式化工具配置
  enable: # 启用的格式化器
    - gofmt     # 标准的 Go 代码格式化
    - gofumpt   # 更严格的 Go 格式化工具 (gofmt 的超集)
    - goimports # 自动管理 import (添加缺失的，删除未用的)
  settings:
    goimports:
      local-prefixes: # 定义本地包前缀，用于 import 分组
        - github.com/SwanHubX # 这个前缀的包会被归类为“本地包”，与其他第三方包分开
  exclusions: # 格式化工具的忽略规则
    generated: lax # 对生成代码放宽格式化要求
    paths: # 以下目录不进行格式化
      - third_party$
      - builtin$
      - examples$
```

## 加载配置和环境变量

为了减少每次写 go 项目都要进行 viper 的配置和环境变量的加载，我们这里写了一个通用的配置加载函数来进行这些工作，函数的目录在 /argo/pkg/config/loader.go 代码并不复杂，这里主要讲如何使用。

首先推荐配置写在 /argo/config/config.yaml 中，如果你讲配置写在其他部分时请传递正确的配置路径。

暴露的函数如图所示：

```go
func Init(configPath string, configName string, envPrefix string) error 
```

这个函数会同时加载 configPath 路径下的 configName.yaml 和 .env 文件，所以 .env 文件使用的时候可以和配置文件放到同一个路径下,这样在运行的时候会自动加载。

配置的优先级如下：

> 环境变量  > .env > config.yaml

envPrefix 可以作为项目的名字，比如环境变量 ARGO_SERVER_PORT 可以使用 ARGO 作为 envprefix，然后我们可以直接使用 `viper.Get("server.port")` 来使用，ARGO_SERVER_PORT 会覆盖 yaml 里面的 server.port。

使用可以见 test 文件。

这里有一个规范，就是因为我们同时把环境变量和 config.yaml 中的配置都加载到了 viper 中，所以使用 viper 进行访问的时候，使用 

```go
viper.GetString("database.url")
```

用上面的代码来获得 ARGO_DATABASE_URL 的配置。或者用 DATABASE_URL 不要加前缀，viper 会自动帮你加。

使用

```go
func main() {
	err := config.Init("configs", "config", "ARGO")
	if err != nil {
		panic(err)
	}
}
```

这里的路径文件路径设置为 configs 是因为一般情况下我们都是在项目的根路径启动 main，而不是在 cmd 中。

## 日志设置

这里本来想给输出不同的日志级别设置不同的颜色的，但是如果要求输出为 json 格式的话，设置在使用会比较麻烦，而且我认为使用 json 作为日志的输出我们通常应该是不在服务器上查看日志的，而是有专门的日志查看工具，所以这里仅仅是把全局的 slog 改为 json 格式，同时使用 gin 中间件在打印日志的时候自动加上 TraceID。

使用

```go
import ""

func main() {
  log.Init("info")
}
```

gin 使用中间件来获取 X-Trace-ID，在 /pkg/middleware/trace.go 中捕获 X-Trace-ID

为了方便我们进行日志的追踪，我设置了一个 slog Handler 用来从 context 中拿去 traceID，所以每次使用 slog 进行打印的时候要使用 slog.ErrorContext(...) 这样日志会自动打印 traceID。

代码如下

```go
func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceId, ok := ctx.Value(middleware.TraceIdKey).(string); ok {
		r.AddAttrs(slog.String("trace_id", traceId))
	}
	return h.Handler.Handle(ctx, r)
}

// NewTraceHandler 创建带 Trace 能力的 Handler
func NewTraceHandler(h slog.Handler) *TraceHandler {
	return &TraceHandler{Handler: h}
}
```

中间件使用如下：

```go
	r := gin.Default()
  // 这里一定要设置，如果不设置这里的话，我们在 slog 的 Handler 中读取 trace_id 的时候就不会从 Request 里面去拿
	r.ContextWithFallback = true
	r.Use(middleware.ErrorHandler)
	r.Use(middleware.TraceMiddleware)
```

中间件打印结果大概如图所示

![image-20251204193520664](http://www.spigcoder.com/boke/image-20251204193520664.png)

我们自己打印的结果会有 time、level、source、和 trace_id 剩下的是我们自己填的。
