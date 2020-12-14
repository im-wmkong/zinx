package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("Client 1 start ...")
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Printf("Client connect error: %s\n", err)
		return
	}
	for {
		dp := znet.NewDataPack()
		data, err := dp.Pack(znet.NewMessage(1, []byte("Hello Zinx")))
		if err != nil {
			fmt.Printf("Pack error: %s\n", err)
			break
		}

		if _, err = conn.Write(data); err != nil {
			fmt.Printf("Connect Write error: %s\n", err)
			break
		}

		head := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, head); err != nil {
			fmt.Printf("Client read head error: %s\n", err)
			break
		}

		imsg, err := dp.Unpack(head)
		if err != nil {
			fmt.Printf("Client uppack head error: %s\n", err)
			break
		}

		if imsg.GetMsgLen() > 0 {
			msg := imsg.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Printf("Client read data error: %s\n", err)
				break
			}
			fmt.Printf("Client read msg, id: %d, data: %s", msg.Id, msg.Data)
		}

		time.Sleep(1 * time.Second)
	}
}
