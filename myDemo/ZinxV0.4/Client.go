package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client start ...")
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Printf("Client connect error: %s\n", err)
		return
	}
	for {
		_, err := conn.Write([]byte("Hello Zinx V0.4"))
		if err != nil {
			fmt.Printf("Connect Write error: %s\n", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Connect read error: %s\n", err)
			return
		}
		fmt.Printf("Connect read massage: %s, length: %d\n", buf, cnt)

		time.Sleep(1 * time.Second)
	}
}
