package cmd

import (
	"log"
	"net/url"
	"os"
	"path/filepath"

	"wscp/core"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var File string

var clientCmd = &cobra.Command{
	Use:   "c",
	Short: "以客户端模式运行",
	Run: func(cmd *cobra.Command, args []string) {
		File, _ := cmd.Flags().GetString("file")
		fileName := filepath.Base(File)
		var addr string
		if len(args) == 0 {
			log.Println("必须指定服务端地址 ip:port")
			os.Exit(1)
		} else {
			addr = args[0]
		}

		u := url.URL{
			Scheme: "ws",
			Host:   addr,
			Path:   "/wscp",
		}
		q := u.Query()
		q.Set("fname", fileName)
		u.RawQuery = q.Encode()
		log.Printf("连接到服务端: %s\n", u.String())
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			panic(err)
		}
		defer c.Close()
		core.ReadFile(File, c)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	// clientCmd.Flags().StringVar(&File, "f", "", "指定文件路径")
	clientCmd.Flags().StringP("file", "f", "", "指定文件路径")
}
