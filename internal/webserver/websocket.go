package webserver

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"nova-panel/internal/store"
)

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// return r.Header.Get("Origin") == "http://your-trusted-domain.com"
		return true // 允许所有来源，生产环境需配置
	},
}

// 客户端连接管理
type wsHub struct {
	clients    map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.RWMutex
}

func newWsHub() *wsHub {
	return &wsHub{
		clients:    make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// 运行 WebSocket Hub，管理连接和广播
func (h *wsHub) run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()
			log.Printf("新 WebSocket 客户端连接，当前连接数: %d", len(h.clients))
		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
			}
			h.mu.Unlock()
			log.Printf("WebSocket 客户端断开，当前连接数: %d", len(h.clients))
		}
	}
}

// 广播 agent 状态给所有客户端
func (h *wsHub) broadcastStatus() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		agents := store.GetAllAgents()
		data, err := json.Marshal(agents)
		if err != nil {
			log.Printf("状态序列化错误: %v", err)
			continue
		}

		h.mu.RLock()
		for conn := range h.clients {
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("WebSocket 发送错误: %v", err)
				h.mu.RUnlock()
				h.unregister <- conn
				h.mu.RLock()
			}
		}
		h.mu.RUnlock()
	}
}

// 处理 WebSocket 连接
func handleWebSocket(h *wsHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket 升级失败: %v", err)
			return
		}

		h.register <- conn

		// 保持连接，监听客户端关闭
		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					h.unregister <- conn
					return
				}
			}
		}()
	}
}

// 初始化 WebSocket 服务并整合到 Gin 路由
func initWebSocket(r *gin.Engine) {
	hub := newWsHub()
	go hub.run()
	go hub.broadcastStatus()
	r.GET("/ws", handleWebSocket(hub))
}
