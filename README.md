# nosfiber

`nosfiber` wraps [Fiber v3](https://github.com/gofiber/fiber) to make HTTP server startup, async launch, and graceful shutdown more expressive within the nos ecosystem.

## Installation

```bash
go get github.com/raaaaaaaay86/nosfiber
```

## Quick Start

### Synchronous Start

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

    // Blocks until the server stops
    if err := app.Start(context.Background(), 8080); err != nil {
        panic(err)
    }
}
```

### Asynchronous Start + Graceful Shutdown

This is the recommended pattern when combined with `nosos.GracefulShutdown`.

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

    // Start in background – returns immediately
    if err := app.StartAsync(context.Background(), 8080); err != nil {
        panic(err)
    }

    // Block until SIGINT / SIGTERM, then shut down
    nosos.GracefulShutdown(context.Background(), nosos.GracefulShutdownSetup{
        ListenedSignals: nosos.DefaultShutdownSignals,
        MaxWait:         15 * time.Second,
        OnShutdown: func(ctx context.Context) error {
            return app.Shutdown()
        },
    })
}
```

## Real-World Example (noschat)

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

## API Reference

### `NewFiberApp`

```go
func NewFiberApp(config fiber.Config, listenConfig fiber.ListenConfig) *FiberApp
```

Creates a new `FiberApp`. Both arguments are passed directly to Fiber.

### `FiberApp` Methods

| Method                  | Description                                                     |
|-------------------------|-----------------------------------------------------------------|
| `Start(ctx, port)`      | Start the server synchronously (blocks until server stops)      |
| `StartAsync(ctx, port)` | Start the server in a goroutine (returns immediately)           |
| `Shutdown()`            | Gracefully shut down the server                                 |
| `App()`                 | Return the underlying `*fiber.App` for route registration       |
