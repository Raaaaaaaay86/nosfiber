package nosfiber

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/xerrors"
)

type FiberApp struct {
	config       fiber.Config
	listenConfig fiber.ListenConfig
	app          *fiber.App
}

func NewFiberApp(
	config fiber.Config,
	listenConfig fiber.ListenConfig,
) *FiberApp {
	return &FiberApp{
		config:       config,
		listenConfig: listenConfig,
		app:          fiber.New(config),
	}
}

func (f *FiberApp) Start(ctx context.Context, port int) error {
	if err := f.app.Listen(fmt.Sprintf(":%d", port), f.listenConfig); err != nil {
		return xerrors.Errorf("cannot listen port (%d)", port)
	}

	return nil
}

func (f *FiberApp) StartAsync(ctx context.Context, port int) error {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic: recovered from FiberApp.Start", "recover", r)
			}
		}()

		if err := f.Start(ctx, port); err != nil {
			slog.Error("failed to start FiberApp", "error", err)
		}
	}()

	return nil
}

func (f *FiberApp) Shutdown() error {
	return f.app.Shutdown()
}

func (f *FiberApp) App() *fiber.App {
	return f.app
}
