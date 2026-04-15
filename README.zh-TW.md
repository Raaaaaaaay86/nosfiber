# nosfiber

`nosfiber` 封裝了 [Fiber v3](https://github.com/gofiber/fiber)，讓 HTTP server 的啟動、非同步啟動與優雅關機在 nos 生態系中更加直觀易用。

## 安裝

```bash
go get github.com/raaaaaaaay86/nosfiber
```

## 快速開始

### 同步啟動

```go
package main

import (
    "context"

    "github.com/gofiber/fiber/v3"
    "github.com/raaaaaaaay86/nosfiber"
)

func main() {
    app := nosfiber.NewFiberApp(
        fiber.Config{},
        fiber.ListenConfig{},
    )

    app.App().Get("/health", func(c fiber.Ctx) error {
        return c.SendString("ok")
    })

    // 阻塞直到 server 停止
    if err := app.Start(context.Background(), 8080); err != nil {
        panic(err)
    }
}
```

### 非同步啟動搭配優雅關機

搭配 `nosos.GracefulShutdown` 時建議使用此模式。

```go
package main

import (
    "context"
    "time"

    "github.com/gofiber/fiber/v3"
    "github.com/raaaaaaaay86/nosfiber"
    "github.com/raaaaaaaay86/nosos"
)

func main() {
    app := nosfiber.NewFiberApp(fiber.Config{}, fiber.ListenConfig{})

    app.App().Get("/ping", func(c fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "pong"})
    })

    // 在背景啟動，立即回傳
    if err := app.StartAsync(context.Background(), 8080); err != nil {
        panic(err)
    }

    // 阻塞至 SIGINT / SIGTERM，然後關機
    nosos.GracefulShutdown(context.Background(), nosos.GracefulShutdownSetup{
        ListenedSignals: nosos.DefaultShutdownSignals,
        MaxWait:         15 * time.Second,
        OnShutdown: func(ctx context.Context) error {
            return app.Shutdown()
        },
    })
}
```

## 實際範例（noschat）

```go
// component/fiber.go

func NewFiberApp(
    c *config.Application,
    members *httpapi.MemberController,
    chatroom *httpapi.ChatroomController,
) *nosfiber.FiberApp {
    fiberApp := nosfiber.NewFiberApp(fiber.Config{}, fiber.ListenConfig{})

    app := fiberApp.App()
    members.RegisterRoutes(app)
    chatroom.RegisterRoutes(app)

    return fiberApp
}

// main.go
if err := fiberApp.StartAsync(ctx, cfg.HTTP.Port); err != nil {
    return err
}

return nosos.GracefulShutdown(ctx, nosos.GracefulShutdownSetup{
    ListenedSignals: nosos.DefaultShutdownSignals,
    OnShutdown: func(ctx context.Context) error {
        return fiberApp.Shutdown()
    },
})
```

## API 說明

### `NewFiberApp`

```go
func NewFiberApp(config fiber.Config, listenConfig fiber.ListenConfig) *FiberApp
```

建立新的 `FiberApp`，兩個參數直接傳遞給 Fiber。

### `FiberApp` 方法

| 方法                    | 說明                                              |
|-------------------------|---------------------------------------------------|
| `Start(ctx, port)`      | 同步啟動 server（阻塞直到 server 停止）           |
| `StartAsync(ctx, port)` | 在 goroutine 中啟動 server（立即回傳）            |
| `Shutdown()`            | 優雅關閉 server                                   |
| `App()`                 | 回傳底層 `*fiber.App`，用於路由註冊               |
