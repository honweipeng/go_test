package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	getExample()
	postExample()
}
func getExample() {
	url := "https://jsonplaceholder.typicode.com/posts/1"

	// 发起GET请求
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("GET request failed:", err)
		return
	}
	//声明一个数组

	defer response.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	// 打印响应内容
	fmt.Println("GET Response:", string(body))
}

// 打开连接控制器
func postExample() {
	url := "https://jsonplaceholder.typicode.com/posts"
	//payload := []byte(`{"title":"foo","body":"bar","userId":1}`)

	// 发起POST请求 bytes.NewBuffer()增加参数到放回体
	response, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println("POST request failed:", err)
		return
	}
	defer response.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	// 打印响应内容
	fmt.Println("POST Response:", string(body))
}
