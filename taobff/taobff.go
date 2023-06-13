package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"xiyu.com/common"
	"xiyu.com/taobff/rest"
)

func gateFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		status := c.Writer.Status()
		common.GlbBaInfa.Logger.Infof("perf stat: url=[%s], time used=[%d], status=[%d]", c.FullPath(), latency.Milliseconds(), status)
	}
}

func hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello")
}

func main() {
	conf := "../config/tao.yaml"
	baInfra := common.BaseInit(conf)
	baInfra.Logger.Info("Hello go")
	baInfra.Logger.Debug("Hello go")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	router := gin.New()
	router.MaxMultipartMemory = 8 << 20
	router.Use(gateFilter())
	router.Use(gin.Recovery())

	router.GET(fmt.Sprintf("%s/hello", baInfra.Conf.Http.Prefix), hello)
	/**注册用户接口**/
	rest.RegisterUserRest(router)
	rest.RegisterBlogRest(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", baInfra.Conf.Http.Port),
		Handler: router,
	}

	// Initializing the server
	go func() {
		/**进行连接**/
		baInfra.Logger.Infof("server listen on:%d", baInfra.Conf.Http.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			baInfra.Logger.Errorf("listen error:%s", err.Error())
		}
	}()
	<-ctx.Done()
	stop()
	baInfra.Logger.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		baInfra.Logger.Warn("Server shutdown:", err)
	}
	common.TcloseItems()
	baInfra.Logger.Info("Server exist")
}

func CreateRestHandler() {
	panic("unimplemented")
}

func RestHandler() {
	panic("unimplemented")
}
