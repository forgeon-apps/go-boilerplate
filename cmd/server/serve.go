package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/byeblogs/go-boilerplate/pkg/config"
	"github.com/byeblogs/go-boilerplate/pkg/middleware"
	"github.com/byeblogs/go-boilerplate/pkg/route"
	"github.com/byeblogs/go-boilerplate/platform/database"
	"github.com/byeblogs/go-boilerplate/platform/logger"
	"github.com/gofiber/fiber/v2"
)

// Serve ..
func Serve() {
	appCfg := config.AppCfg()

	logger.SetUpLogger()
	logr := logger.GetLogger()

	// connect to DB
	if err := database.ConnectDB(); err != nil {
		logr.Panicf("failed database setup. error: %v", err)
	}

	// Define Fiber config & app.
	fiberCfg := config.FiberConfig()
	app := fiber.New(fiberCfg)

	// Attach Middlewares.
	middleware.FiberMiddleware(app)

	// Routes.
	route.GeneralRoute(app)
	route.SwaggerRoute(app)
	route.PublicRoutes(app)
	route.PrivateRoutes(app)
	route.NotFoundRoute(app)

	// signal channel to capture system calls
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// start shutdown goroutine
	go func() {
		// capture sigterm and other system call here
		<-sigCh
		logr.Infoln("Shutting down server...")
		_ = app.Shutdown()
	}()

	// start http server
	serverAddr := fmt.Sprintf("%s:%d", appCfg.Host, appCfg.Port)
	if err := app.Listen(serverAddr); err != nil {
		logr.Errorf("Oops... server is not running! error: %v", err)
	}

}
