// 处理telnet端口的长连接
// 能以字节发送数据
// 能控制关闭
// 解决并发连接断开程序直接关闭
// 2023-11-27 处理客户端字节发送，模拟控制器应用，发字节
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

// sendDataToServer 函数用于向服务器发送数据
func sendDataToServer(conn net.Conn, data []byte) error {
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println("发送数据错误:", err)
		return err
	}
	return nil
}

// receiveDataFromServer 函数用于异步接收服务器的数据
func receiveDataFromServer(conn net.Conn, wg *sync.WaitGroup, reconnect chan struct{}) {
	defer wg.Done()

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("IO流断开，请检查网络！:", err)
			close(reconnect) // 通知主线程进行重新连接
			return
		}

		//quotedData := strconv.QuoteToASCII(string(buffer[:n]))
		quotedData := string(buffer[:n])
		// 使用 Unquote 进行解除转义
		//unquotedData, err := strconv.Unquote(quotedData)
		//if err != nil {
		//	fmt.Println("解除转义错误:", err)
		//	continue
		//}

		fmt.Printf("从服务器接收到ASCII格式数据: %s\n", quotedData)
	}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// 定义服务器地址
	serverAddr := "127.0.0.1:8080"

	for {
		// 连接服务器
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			fmt.Println("连接服务器错误:", err)
			fmt.Println("尝试重新连接...")
			time.Sleep(5 * time.Second) // 休眠5秒后尝试重新连接
			continue
		}
		//是localhost
		fmt.Println("已连接到服务器:", conn.RemoteAddr().String())

		// 使用 WaitGroup 等待 goroutine 完成
		var wg sync.WaitGroup
		wg.Add(1)

		// 启动 goroutine 异步接收服务器的响应
		reconnect := make(chan struct{})
		go func() {
			receiveDataFromServer(conn, &wg, reconnect)
		}()

		// 启动 goroutine 读取控制台输入并发送到服务器
		go func() {
			defer wg.Done()
			for {
				fmt.Print("请输入要发送的消息: ")
				reader := bufio.NewReader(os.Stdin)
				input, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("读取输入错误:", err)
					continue
				}

				// 准备要发送的数据
				message := []byte(input)

				// 发送数据到服务器
				err = sendDataToServer(conn, message)
				if err != nil {
					conn.Close()
					fmt.Println("尝试重新连接...")
					time.Sleep(5 * time.Second) // 休眠5秒后尝试重新连接
					reconnect <- struct{}{}     // 发送重新连接信号
					return
				}
			}
		}()

		select {
		case <-interrupt:
			fmt.Println("中断。正在关闭连接.")
			err := conn.Close()
			if err != nil {
				fmt.Println("关闭连接时发生错误:", err)
			}
			wg.Wait() // 等待接收数据的 goroutine 完成
			return
		case <-reconnect:
			fmt.Println("连接断开，尝试重新连接...")
			// 关闭连接后等待接收数据的 goroutine 完成
			conn.Close()
			wg.Wait()
			time.Sleep(5 * time.Second) // 休眠5秒后尝试重新连接
		}
	}
}
