<div align="center">
<h1 style="margin: 0; padding: 0">GoCaptcha Service SDK</h1>
<p style="margin: 0; padding: 0">GoCaptcha Servce 服务的接口开发工具包</p>
<br/>
<a href="https://goreportcard.com/report/github.com/wenlng/go-captcha-service-sdk"><img src="https://goreportcard.com/badge/github.com/wenlng/go-captcha-service-sdk"/></a>
<a href="https://godoc.org/github.com/wenlng/go-captcha-service-sdk"><img src="https://godoc.org/github.com/wenlng/go-captcha-service-sdk?status.svg"/></a>
<a href="https://github.com/wenlng/go-captcha-service-sdk/releases"><img src="https://img.shields.io/github/v/release/wenlng/go-captcha-service-sdk.svg"/></a>
<a href="https://github.com/wenlng/go-captcha-service-sdk/blob/LICENSE"><img src="https://img.shields.io/badge/License-MIT-green.svg"/></a>
<a href="https://github.com/wenlng/go-captcha-service-sdk"><img src="https://img.shields.io/github/stars/wenlng/go-captcha-service-sdk.svg"/></a>
<a href="https://github.com/wenlng/go-captcha-service-sdk"><img src="https://img.shields.io/github/last-commit/wenlng/go-captcha-service-sdk.svg"/></a>
</div>

<br/>

`GoCaptcha Service SDK` 是行为验证码服务 [GoCaptcha Service](https://github.com/wenlng/go-captcha-service) 的接口开发工具包，用于在 Go 应用程序中快速集成验证码相关功能。它支持两种通信协议：HTTP（通过 `resetapi` 包）和 gRPC（通过 `grpcapi` 包）。该 SDK 设计用于处理验证码数据的获取、验证、状态检查以及资源的管理，并支持服务发现功能，通过 `sdlb` 组件支持服务发现和负载均衡。


<br/>

> [English](README.md) | 中文
<p> ⭐️ 如果能帮助到你，请随手给点一个star</p>


## 功能特性

- **多协议支持**：提供 HTTP 和 gRPC 两种协议的客户端实现，满足不同场景需求。
- **核心功能**：
    - 获取验证码数据（如图片、尺寸等）。
    - 验证用户提交的验证码数据。
    - 检查验证码状态。
    - 获取验证码状态详细信息。
    - 删除验证码状态信息。
- **资源管理 (仅 HTTP)**：
    - 上传资源文件。
    - 删除指定资源。
    - 获取资源文件列表。
    - 获取和更新验证码配置。
- **可支持服务发现与负载均衡**：通过 `sdlb` 组件支持多实例部署和服务选择。
- **重试机制**：内置可配置的重试策略，支持自定义重试次数和等待时间。
- **超时控制**：支持请求超时设置，确保高可用性。
- **认证支持**：通过 API 密钥 (`X-API-Key`) 进行身份验证。

<br/>

---

## 安装

### 设置Go代理
- Window
```shell
$ set GO111MODULE=on
$ set GOPROXY=https://goproxy.io,direct

### The Golang 1.13+ can be executed directly
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.io,direct
```
- Linux or Mac
```shell
$ export GO111MODULE=on
$ export GOPROXY=https://goproxy.io,direct

### or
$ echo "export GO111MODULE=on" >> ~/.profile
$ echo "export GOPROXY=https://goproxy.cn,direct" >> ~/.profile
$ source ~/.profile
```

### 安装
```shell
$ go get -u github.com/wenlng/go-captcha-service-sdk@latest
```

---

## 使用方法

SDK 提供了两个客户端实现：`resetapi`（基于 HTTP）和 `grpcapi`（基于 gRPC）。两者实现了相同的 `Client` 接口，支持核心验证码操作。以下分别介绍两种客户端的初始化和使用方法。

### HTTP 客户端 (`resetapi`)

`resetapi` 包使用 HTTP 协议与验证码服务通信，依赖 `go-resty` 进行 HTTP 请求。

#### 初始化 HTTP 客户端

```go
package main

import (
    "context"
    "time"
    "github.com/wenlng/go-captcha-service-sdk/golang/resetapi"
)

func main() {
    config := resetapi.ClientConfig{
        BaseUrl:          "http://127.0.0.1:8080",
        APIKey:           "your-api-key",
        RetryCount:       3,
        RetryWaitTime:    500 * time.Millisecond,
        RetryMaxWaitTime: 3 * time.Second,
        Timeout:          10 * time.Second,
    }

    client, err := resetapi.NewHTTPClient(config, nil)
    if err != nil {
        panic(err)
    }

    // 使用 client 进行操作
}
```

### gRPC 客户端 (`grpcapi`)

`grpcapi` 包使用 gRPC 协议，依赖 `google.golang.org/grpc` 进行远程调用。

#### 初始化 gRPC 客户端

```go
package main

import (
    "context"
    "time"
    "github.com/wenlng/go-captcha-service-sdk/golang/grpcapi"
)

func main() {
    config := grpcapi.ClientConfig{
        BaseAddress:      "127.0.0.1:8080",
        APIKey:           "your-api-key",
        RetryCount:       3,
        RetryWaitTime:    500 * time.Millisecond,
        RetryMaxWaitTime: 3 * time.Second,
        Timeout:          10 * time.Second,
    }

    client, err := grpcapi.NewGRPCClient(config, nil)
    if err != nil {
        panic(err)
    }

    // 使用 client 进行操作
}
```

---

## 服务发现与负载均衡

### SDLB 组件简介

`SDLB`（Service Discovery and Load Balancing）是 SDK 提供的服务发现与负载均衡组件，用于在分布式环境中动态选择验证码服务实例。`SDLB` 通过服务注册信息（如主机和端口）选择合适的实例，支持高可用性和负载均衡。

#### 核心功能

- **服务发现**：从服务注册中心或配置中获取可用实例列表。
- **负载均衡**：根据负载均衡策略（如轮询、随机或基于权重的选择）选择服务实例。
- **动态更新**：支持实例列表的动态更新，适应服务实例的添加或移除。
- **地址选择**：通过 `SelectAddress` 或 `SelectAddressWithKey` 方法动态选择服务地址。

#### 工作原理

1. **实例选择**：`SDLB` 调用 `Select` 方法，也可以根据提供的 `key`（通常是用户ID、文件ID、当前主机名等）计算与选择一个服务实例。
2. **地址格式化**：
    - 对于 HTTP 客户端，返回格式为 `http://<host>:<http_port>`。
    - 对于 gRPC 客户端，返回格式为 `<host>:<grpc_port>`。
3. **回退机制**：如果 `SDLB` 不可用或选择失败，客户端会回退到 `ClientConfig` 中的 `BaseUrl` 或 `BaseAddress`。
4. **自定义地址**：通过 `FilterHost` 函数，用户可以自定义地址格式。

### 配置 SDLB

要使用 `SDLB`，需要初始化 `sdlb.SDLB` 实例并将其传递给客户端构造函数。`SDLB` 的具体实现依赖于外部服务注册中心（如 Consul、Etcd、Nacos、ZooKeeper）。


### SDLB 使用示例

以下是将 `SDLB` 集成到 HTTP 客户端的完整示例（gRPC 同理）：

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/wenlng/go-captcha-service-sdk/golang/resetapi"
    "github.com/wenlng/go-captcha-service-sdk/golang/sdlb"
)

var restapiCli resetapi.Client
var sdlbInts *sdlb.SDLB

// setupHttpClient .
func setupHttpClient() error {
	sdlbInstance, err := sdlb.NewServiceDiscoveryLB(sdlb.ClientConfig{
		ServiceDiscoveryType: sdlb.ServiceDiscoveryTypeNacos,
		Addrs:                "localhost:8848",
		LoadBalancerType:     sdlb.LoadBalancerTypeRandom,
		ServiceName:          "go-captcha-service",
		Username:             "nacos",
		Password:             "nacos",
	})

	if err != nil {
		return fmt.Errorf("failed to new sdlb: %v", err)
	}

	sdlbInts = sdlbInstance
	restapiCli, err = resetapi.NewHTTPClient(resetapi.ClientConfig{
		APIKey: ApiKey,
	}, sdlbInstance)

	return err
}

func closeSDLB() {
	if sdlbInts != nil {
		sdlbInts.Close()
	}
}

func main() {
	err := setupHttpClient()
	if err != nil {
		fmt.Printf("Failed to new sdlb: %v\n", err)
		return
	}
	defer closeSDLB()
	
    // 获取验证码数据
    ctx := context.Background()
    captchaData, err := restapiCli.GetData(ctx, "click-default-ch")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Captcha Key: %s\n", captchaData.CaptchaKey)
    fmt.Printf("Selected Address: %s\n", captchaData.MasterImageBase64)
}
```

---

## 更多示例代码

以下是通过 HTTP 和 gRPC 客户端实现常见操作的示例代码。

### 获取验证码数据

获取验证码数据（如图片和尺寸信息）用于前端展示。

#### HTTP 示例

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/wenlng/go-captcha-service-sdk/golang/resetapi"
)

func main() {
    config := resetapi.ClientConfig{
        BaseUrl: "http://127.0.0.1:8080",
        APIKey:  "your-api-key",
    }

    client, err := resetapi.NewHTTPClient(config, nil)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    captchaData, err := client.GetData(ctx, "click-default-ch")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Captcha Key: %s\n", captchaData.CaptchaKey)
    fmt.Printf("Master Image (Base64): %s\n", captchaData.MasterImageBase64)
}
```

#### gRPC 示例

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/wenlng/go-captcha-service-sdk/golang/grpcapi"
)

func main() {
    config := grpcapi.ClientConfig{
        BaseAddress: "127.0.0.1:8080",
        APIKey:      "your-api-key",
    }

    client, err := grpcapi.NewGRPCClient(config, nil)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    captchaData, err := client.GetData(ctx, "click-default-ch")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Captcha Key: %s\n", captchaData.CaptchaKey)
    fmt.Printf("Master Image (Base64): %s\n", captchaData.MasterImageBase64)
}
```

### 验证验证码数据

验证用户提交的验证码数据是否正确。

#### HTTP 示例

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/wenlng/go-captcha-service-sdk/golang/resetapi"
)

func main() {
    config := resetapi.ClientConfig{
        BaseUrl: "http://127.0.0.1:8080",
        APIKey:  "your-api-key",
    }

    client, err := resetapi.NewHTTPClient(config, nil)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    success, err := client.CheckData(ctx, "click-default-ch", "captcha-key-456", "100,80,23,80}")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    if success {
        fmt.Println("Verification passed!")
    } else {
        fmt.Println("Verification failed.")
    }
}
```

#### gRPC 示例

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/wenlng/go-captcha-service-sdk/golang/grpcapi"
)

func main() {
    config := grpcapi.ClientConfig{
        BaseAddress: "127.0.0.1:8080",
        APIKey:      "your-api-key",
    }

    client, err := grpcapi.NewGRPCClient(config, nil)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    success, err := client.CheckData(ctx, "click-default-ch", "captcha-key-456", "100,200")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    if success {
        fmt.Println("Verification passed!")
    } else {
        fmt.Println("Verification failed.")
    }
}
```

### 上传资源文件 (仅 HTTP)

上传图片或其他资源文件到验证码服务。

#### HTTP 示例

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    "github.com/wenlng/go-captcha-service-sdk/golang/resetapi"
)

func main() {
    config := resetapi.ClientConfig{
        BaseUrl: "http://127.0.0.1:8080",
        APIKey:  "your-api-key",
    }

    client, err := resetapi.NewHTTPClient(config, nil)
    if err != nil {
        panic(err)
    }

    // 打开文件
    file, err := os.Open("captcha_image.png")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    ctx := context.Background()
    success, partialSuccess, err := client.UploadResource(ctx, "captcha_images", []*os.File{file})
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    if success {
        fmt.Println("Upload successful!")
        if partialSuccess {
            fmt.Println("Some files were uploaded successfully.")
        }
    } else {
        fmt.Println("Upload failed.")
    }
}
```

---

## 配置说明

`ClientConfig` 结构体用于配置客户端行为，适用于 `resetapi` 和 `grpcapi`：

| 字段                | 描述                                                                 | 默认值                     |
|---------------------|----------------------------------------------------------------------|----------------------------|
| `BaseUrl` / `BaseAddress` | 服务地址（HTTP 或 gRPC）                                             | 无（必填）                 |
| `APIKey`            | API 密钥，用于身份验证                                               | 无（必填）                 |
| `RetryCount`        | 请求失败时的重试次数                                                 | 3                          |
| `RetryWaitTime`     | 初次重试的等待时间                                                   | 500ms                      |
| `RetryMaxWaitTime`  | 最大重试等待时间（指数退避）                                         | 3s                         |
| `Timeout`           | 请求超时时间                                                         | 10s                        |
| `FilterHost`        | 自定义地址格式的函数（如将主机和端口转换为特定格式）                 | nil                        |

### 配置示例

```go
config := resetapi.ClientConfig{
    BaseUrl:          "http://127.0.0.1:8080",
    APIKey:           "your-api-key",
    RetryCount:       5,
    RetryWaitTime:    1 * time.Second,
    RetryMaxWaitTime: 5 * time.Second,
    Timeout:          15 * time.Second,
    FilterHost: func(host, port string) string {
        return fmt.Sprintf("http://%s:%s", host, port)
    },
}
```

## 许可证

本项目采用 [MIT 许可证](LICENSE)。详情请参阅 `LICENSE` 文件。