package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"

	"github.com/isucon/isucon11-qualify/jiaapi-mock/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

//go:embed ec256-private.pem
var privateKey []byte

func main() {
	// Controllers
	authController, err := controller.NewAuthController(privateKey)
	if err != nil {
		panic(err)
	}
	catalogController := controller.CatalogController{}
	activationController := controller.NewActivationController()

	// Echo instance
	e := echo.New()
	e.Debug = true
	e.Logger.SetLevel(log.DEBUG)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 動作確認用のログインページ
	e.Static("/", "./index.html")
	// APIs
	e.POST("/api/auth", authController.PostAuth)
	e.GET("/api/catalog/:catalog_id", catalogController.GetCatalog)
	e.POST("/api/activate", activationController.PostActivate)
	e.POST("/api/deactivate", activationController.PostDeactivate)
	e.POST("/api/die", func(ctx echo.Context) error {
		input := &struct {
			Password string `json:"password" validate:"required"`
		}{}
		err := ctx.Bind(input)
		if err != nil {
			ctx.Logger().Errorf("failed to bind: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		if input.Password != "U,YaCLe9tAnW8EdYphW)Wc/dN)5pPQ/3ue_af4rz" {
			return echo.NewHTTPError(http.StatusNotFound, "Not Found")
		}
		os.Exit(0)
		return nil
	})

	// Start server
	serverPort := fmt.Sprintf(":%v", getEnv("JIAAPI_SERVER_PORT", "5000"))
	e.Logger.Fatal(e.Start(serverPort))
}

func getEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultValue
}