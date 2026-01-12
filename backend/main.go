package main

import (
	"auction-backend/blockchain"
	"auction-backend/config"
	"auction-backend/database"
	"auction-backend/routes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(config.AppConfig.GetDSN()); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动事件监听器
	listener, err := blockchain.NewEventListener()
	if err != nil {
		log.Fatalf("Failed to create event listener: %v", err)
	}

	go func() {
		if err := listener.StartListening(ctx); err != nil {
			log.Printf("Event listener error: %v", err)
		}
	}()

	// 设置 Gin 路由
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	})

	// 设置路由
	routes.SetupRoutes(r)

	// 启动 HTTP 服务器
	srv := &http.Server{
		Addr:    ":" + config.AppConfig.ServerPort,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s...", config.AppConfig.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 优雅关闭
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
