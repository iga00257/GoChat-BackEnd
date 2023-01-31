package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//	Upgrader
//
// ReadBufferSize 和 WriteBufferSize
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// 可以用来检查连接的来源
	// 这将允许从我们的 React 服务向这里发出请求。
	// 现在，我们可以不需要检查并运行任何连接
	CheckOrigin: func(r *http.Request) bool { return true },
}

// 定义一个 reader listen message from WebSocket
func reader(conn *websocket.Conn) {
	for {
		// read message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print message
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

// WebSocket handler
func serveWs(c *gin.Context) {
	w := c.Writer
	r := c.Request
	fmt.Println(w)
	fmt.Println(r)

	// switch to websocket
	// ws, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	log.Println(err)
	// }

	// listen WebSocket message
	// reader(ws)
}

// Socket handler
func SocketHandler(c *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		closeSocketErr := ws.Close()
		if closeSocketErr != nil {
			panic(err)
		}
	}()

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Message Type: %d, Message: %s\n", msgType, string(msg))
		ws.WriteMessage(msgType, msg)
		// err = ws.WriteJSON(struct {
		// 	Reply string `json:"reply"`
		// }{
		// 	Reply: "Echo...",
		// })
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	router := gin.New()
	router.Use(gin.Recovery())
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}
	//GET Query
	router.GET("/", func(c *gin.Context) {
		firstname := "Liu"
		lastname := "Derrick" // c.Request.URL.Query().Get("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	router.GET("/ws", SocketHandler)
	router.Run(":3000")
}
