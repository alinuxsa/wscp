package cmd

import (
	"io"
	"log"
	"net/http"

	"wscp/core"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var addr string

var serverCmd = &cobra.Command{
	Use:   "s",
	Short: "以服务端模式运行,默认监听 0.0.0.0:5467",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			addr = "0.0.0.0:5467"
		} else {
			addr = args[0]
		}
		log.Printf("服务端监听: %s\n", addr)
		http.HandleFunc("/wscp", server)
		log.Fatal(http.ListenAndServe(addr, nil))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func server(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	fileName := vars.Get("fname")
	var upgrader = websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	if len(fileName) == 0 {
		log.Println("未接收到文件名,关闭连接!")
		c.Close()
	} else {
		log.Printf("接收到文件名: %s\n", fileName)
		core.CheckFile(fileName)
	}
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				c.Close()
				// 非正常客户端关闭连接才报错
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Panic(err)
				}
				break
			}
		}
		core.WriteFile(fileName, message)

	}
}
