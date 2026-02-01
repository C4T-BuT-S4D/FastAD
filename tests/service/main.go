package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()

	mu := sync.Mutex{}
	flags := map[string]string{}
	e.GET("/ping/", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/put/", func(c echo.Context) error {
		var req struct {
			FlagID string `json:"id"`
			Vuln   string `json:"vuln"`
			Flag   string `json:"flag"`
		}
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		mu.Lock()
		defer mu.Unlock()
		e.Logger.Infof("flag: %s, vuln: %s, id: %s", req.Flag, req.Vuln, req.FlagID)
		flags[fmt.Sprintf("%s:%s", req.FlagID, req.Vuln)] = req.Flag

		return c.NoContent(http.StatusOK)
	})

	e.GET("/get/", func(c echo.Context) error {
		var req struct {
			FlagID string `query:"id"`
			Vuln   string `query:"vuln"`
		}
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		mu.Lock()
		defer mu.Unlock()

		flag, ok := flags[fmt.Sprintf("%s:%s", req.FlagID, req.Vuln)]
		e.Logger.Infof("flag: %s, ok: %v, vuln: %s, id: %s", flag, ok, req.Vuln, req.FlagID)
		if !ok {
			return c.String(http.StatusNotFound, "not found")
		}

		return c.JSON(http.StatusOK, map[string]string{"flag": flag})
	})

	e.Logger.SetLevel(log.INFO)
	e.Logger.Fatal(e.Start(":1323"))
}
