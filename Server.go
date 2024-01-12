// 模拟服务端
package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("接收数据错误:", err)
			//冒泡排序

			return
		}

		receivedData := string(buffer[:n])
		fmt.Printf("接收到客户端数据: %s\n", receivedData)

		// 原样返回数据给客户端

		_, err = conn.Write([]byte("服务器收到你的消息, " + receivedData))
		if err != nil {
			fmt.Println("发送响应错误:", err)
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("监听端口错误:", err)
		return
	}
	//123
	defer listener.Close()

	fmt.Println("服务器正在监听端口 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接受连接错误:", err)
			continue
		}

		fmt.Printf("接受到来自 %s 的连接\n", conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}
