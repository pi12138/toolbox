package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Addrs []string `json:"addrs"`
}

func server() {
	cfgData, err := os.ReadFile(configPath)
	if err != nil {
		slog.Error("os.ReadFile error.", "error", err)
		return
	}

	var cfg Config
	if err := json.Unmarshal(cfgData, &cfg); err != nil {
		slog.Error("json.Unmarshal error.", "error", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	r := gin.Default()
	router(r, cancel)

	srvCh := make(chan *http.Server, len(cfg.Addrs))
	for i := 0; i < len(cfg.Addrs); i++ {
		go func(addr string) {
			srv := &http.Server{
				Addr:    addr,
				Handler: r.Handler(),
			}
			srvCh <- srv
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("srv.ListenAndServe error.", "error", err)
			}
		}(cfg.Addrs[i])
	}

	<-ctx.Done()
	slog.Info("server stoping")
	close(srvCh)
	for srv := range srvCh {
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("srv.Shutdown error.", "error", err, "addr", srv.Addr)
		} else {
			slog.Info("srv.Shutdown success.", "addr", srv.Addr)
		}
	}
	slog.Info("server stop end.")
}

var pingCost int64

func router(r *gin.Engine, cancel context.CancelFunc) {
	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(time.Duration(pingCost) * time.Second)
		c.String(http.StatusOK, "pong")
	})
	r.POST("/all_close", func(c *gin.Context) {
		cancel()
		c.String(http.StatusOK, "")
	})
	r.GET("/ping_cost", func(c *gin.Context) {
		v := c.Query("cost")
		pingCost, _ = strconv.ParseInt(v, 10, 0)
		slog.Info("set pingCost.", "value", pingCost)
		c.String(http.StatusOK, "")
	})
}
