package core

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func CheckFile(fn string) {
	// 判断文件是否存在
	_, err := os.Stat(fn)
	if err == nil {
		os.Rename(fn, fn+"_bak")
	}
}

func WriteFile(fn string, message []byte) error {
	fileHandle, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fileHandle.Close()
	buf := bufio.NewWriter(fileHandle)
	buf.Write([]byte(message))
	err = buf.Flush()
	if err != nil {
		fmt.Println("flush error :", err)
		return err
	}
	return nil
}

func ReadFile(fn string, c *websocket.Conn) {
	_, err := os.Stat(fn)
	if err != nil {
		panic(err)
	}
	// 读文件写入c
	f, err := os.OpenFile(fn, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	buf := make([]byte, 2048)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
			} else {
				log.Panic(err)
			}
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			break
		}
		err = c.WriteMessage(websocket.BinaryMessage, buf[:n])
		if err != nil {
			panic(err)
		}
	}
}
