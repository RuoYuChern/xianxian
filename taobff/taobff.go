package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"strings"
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

func JwtAuthMiddleware(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.String(http.StatusForbidden, "authorization empty")
		c.Abort()
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if (len(parts) != 2) || (parts[0] != "Bearer") {
		c.String(http.StatusForbidden, "authorization error")
		c.Abort()
		return
	}

	err := common.VerifyToken(parts[1], c)
	if err != nil {
		c.String(http.StatusForbidden, "authorization faild")
		c.Abort()
		return
	}
	c.Next()
}

func hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello")
}

func auth(c *gin.Context) {
	c.String(http.StatusOK, "Auth ok")
}

func main() {
	conf := "../config/tao.yaml"
	baInfra := common.BaseInit(conf)
	baInfra.Logger.Info("Hello go")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	router := gin.New()
	router.MaxMultipartMemory = 8 << 20
	router.Use(gateFilter())
	router.Use(gin.Recovery())

	router.GET(fmt.Sprintf("%s/hello", baInfra.Conf.Http.Prefix), hello)
	router.GET(fmt.Sprintf("%s/auth", baInfra.Conf.Http.Prefix), JwtAuthMiddleware, auth)
	/**注册用户接口**/
	rest.RegisterUserRest(router, JwtAuthMiddleware)
	rest.RegisterBlogRest(router, JwtAuthMiddleware)

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
