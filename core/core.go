package core

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
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
	finfo, _ := f.Stat()
	total := finfo.Size()
	// bar := progressbar.Default(total)

	bar := progressbar.NewOptions(int(total),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		// progressbar.OptionSetWidth(100),
		progressbar.OptionSetDescription("[cyan][1/1][reset] uploading file..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
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
		bar.Add(n)
	}
}
