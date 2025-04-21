<div align="center">
<h1 style="margin: 0; padding: 0">GoCaptcha Service SDK</h1>
<p style="margin: 0; padding: 0">Interface Development Kit for GoCaptcha Service</p>
<br/>
<a href="https://goreportcard.com/report/github.com/wenlng/go-captcha-service-sdk"><img src="https://goreportcard.com/badge/github.com/wenlng/go-captcha-service-sdk"/></a>
<a href="https://godoc.org/github.com/wenlng/go-captcha-service-sdk"><img src="https://godoc.org/github.com/wenlng/go-captcha-service-sdk?status.svg"/></a>
<a href="https://github.com/wenlng/go-captcha-service-sdk/releases"><img src="https://img.shields.io/github/v/release/wenlng/go-captcha-service-sdk.svg"/></a>
<a href="https://github.com/wenlng/go-captcha-service-sdk/blob/LICENSE"><img src="https://img.shields.io/badge/License-MIT-green.svg"/></a>
<a href="https://github.com/wenlng/go-captcha-service-sdk"><img src="https://img.shields.io/github/stars/wenlng/go-captcha-service-sdk.svg"/></a>
<a href="https://github.com/wenlng/go-captcha-service-sdk"><img src="https://img.shields.io/github/last-commit/wenlng/go-captcha-service-sdk.svg"/></a>
</div>

<br/>

`GoCaptcha Service SDK` is an interface development kit for the behavioral CAPTCHA service [GoCaptcha Service](https://github.com/wenlng/go-captcha-service), designed to quickly integrate CAPTCHA-related functionality into Go applications. It supports two communication protocols: HTTP (via the `resetapi` package) and gRPC (via the `grpcapi` package). The SDK is built to handle CAPTCHA data retrieval, validation, status checking, and resource management, with support for service discovery and load balancing through the `sdlb` component.

<br/>

> English | [中文](README_zh.md)
<p> ⭐️ If this helps you, please consider giving it a star!</p>

## Features

- **Multi-Protocol Support**: Provides client implementations for both HTTP and gRPC protocols to meet diverse use case requirements.
- **Core Functionality**:
  - Retrieve CAPTCHA data (e.g., images, dimensions).
  - Validate user-submitted CAPTCHA data.
  - Check CAPTCHA status.
  - Retrieve detailed CAPTCHA status information.
  - Delete CAPTCHA status information.
- **Resource Management (HTTP only)**:
  - Upload resource files.
  - Delete specific resources.
  - Retrieve resource file lists.
  - Retrieve and update CAPTCHA service configuration.
- **Service Discovery and Load Balancing**: Supports multi-instance deployment and service selection via the `sdlb` component.
- **Retry Mechanism**: Built-in configurable retry strategy with customizable retry counts and wait times.
- **Timeout Control**: Supports request timeout settings for high availability.
- **Authentication Support**: Uses API keys (`X-API-Key`) for authentication.

<br/>

---

## Install the SDK

```shell
$ go get -u github.com/wenlng/go-captcha-service-sdk@latest
```


## Usage

The SDK provides two client implementations: `resetapi` (HTTP-based) and `grpcapi` (gRPC-based). Both implement the same `Client` interface, supporting core CAPTCHA operations. Below are instructions for initializing and using each client.

### HTTP Client (`resetapi`)

The `resetapi` package communicates with the CAPTCHA service using HTTP, relying on the `go-resty` library for HTTP requests.

#### Initializing the HTTP Client

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

    // Use the client for operations
}
```

### gRPC Client (`grpcapi`)

The `grpcapi` package uses the gRPC protocol, relying on the `google.golang.org/grpc` library for remote calls.

#### Initializing the gRPC Client

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

    // Use the client for operations
}
```

---

## Service Discovery and Load Balancing

### SDLB Component Overview

The `SDLB` (Service Discovery and Load Balancing) component is provided by the SDK to dynamically select CAPTCHA service instances in a distributed environment. `SDLB` retrieves available instances from a service registry or configuration and selects an appropriate instance based on load balancing strategies, ensuring high availability and balanced load distribution.

#### Core Features

- **Service Discovery**: Retrieves a list of available instances from a service registry or static configuration.
- **Load Balancing**: Selects instances based on load balancing strategies (e.g., round-robin, random, or weighted selection).
- **Dynamic Updates**: Supports dynamic updates to the instance list, accommodating the addition or removal of service instances.
- **Address Selection**: Uses `SelectAddress` or `SelectAddressWithKey` methods to dynamically select a service address.

#### How It Works

1. **Instance Selection**: `SDLB` calls the `Select` method, which can select a service instance based on a provided `key` (e.g., user ID, file ID, or current hostname).
2. **Address Formatting**:
  - For HTTP clients, returns addresses in the format `http://<host>:<http_port>`.
  - For gRPC clients, returns addresses in the format `<host>:<grpc_port>`.
3. **Fallback Mechanism**: If `SDLB` is unavailable or fails to select an instance, the client falls back to the `BaseUrl` or `BaseAddress` specified in `ClientConfig`.
4. **Custom Address Formatting**: Users can customize address formats using the `FilterHost` function.

### Configuring SDLB

To use `SDLB`, initialize an `sdlb.SDLB` instance and pass it to the client constructor. The specific implementation of `SDLB` depends on external service registries (e.g., Consul, Etcd, Nacos, ZooKeeper).

### SDLB Usage Example

Below is a complete example of integrating `SDLB` with the HTTP client (the process is similar for gRPC):

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

// setupHttpClient initializes the HTTP client with SDLB
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
        APIKey: "your-api-key",
    }, sdlbInstance)

    return err
}

// closeSDLB closes the SDLB instance
func closeSDLB() {
    if sdlbInts != nil {
        sdlbInts.Close()
    }
}

func main() {
    err := setupHttpClient()
    if err != nil {
        fmt.Printf("Failed to initialize SDLB: %v\n", err)
        return
    }
    defer closeSDLB()

    // Retrieve CAPTCHA data
    ctx := context.Background()
    captchaData, err := restapiCli.GetData(ctx, "click-default-ch")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Captcha Key: %s\n", captchaData.CaptchaKey)
    fmt.Printf("Master Image (Base64): %s\n", captchaData.MasterImageBase64)
}
```

---

## Additional Example Code

Below are examples of common operations using the HTTP and gRPC clients.

### Retrieving CAPTCHA Data

Retrieve CAPTCHA data (e.g., images and dimensions) for frontend display.

#### HTTP Example

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

#### gRPC Example

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

### Validating CAPTCHA Data

Validate whether user-submitted CAPTCHA data is correct.

#### HTTP Example

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
    success, err := client.CheckData(ctx, "click-default-ch", "captcha-key-456", "[100,80,23,80]")
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

#### gRPC Example

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
    success, err := client.CheckData(ctx, "click-default-ch", "captcha-key-456", "[100,200]")
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

### Uploading Resource Files (HTTP only)

Upload images or other resource files to the CAPTCHA service.

#### HTTP Example

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

    // Open the file
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

## Configuration Details

The `ClientConfig` struct configures client behavior for both `resetapi` and `grpcapi`:

| Field               | Description                                                          | Default Value              |
|---------------------|----------------------------------------------------------------------|----------------------------|
| `BaseUrl` / `BaseAddress` | Service address (HTTP or gRPC)                                       | None (required)            |
| `APIKey`            | API key for authentication                                           | None (required)            |
| `RetryCount`        | Number of retries for failed requests                                | 3                          |
| `RetryWaitTime`     | Initial wait time for retries                                        | 500ms                      |
| `RetryMaxWaitTime`  | Maximum wait time for retries (exponential backoff)                  | 3s                         |
| `Timeout`           | Request timeout duration                                             | 10s                        |
| `FilterHost`        | Function to customize address format (e.g., host and port formatting) | nil                        |

### Configuration Example

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

---

## License

This project is licensed under the [MIT License](LICENSE). See the `LICENSE` file for details.