package main

import (
	"fmt"
	"math/rand"
	"time"
)

// iota用法：作用域内（其他作用域声明重新计0）变量声明计数类似index下标
const (
	a = iota
	b = iota
)
const (
	name  = "name"
	c     = iota
	d     = iota
	names = "names"

	f = iota
)

func main() {
	maxNum := 100
	rand.Seed(time.Now().UnixNano()) //  设置随机数种子
	secretNumber := rand.Intn(maxNum)
	fmt.Println("The secret number is ", secretNumber)
	//iotad的用法:ota 在 const 关键字出现时将被重置为0，const中每新增一行常量声明将使 iota 计数一次。
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
	fmt.Println(names)
	fmt.Println(f)
	//声明一个数组

	//var teacherNameArray = [1024]{1024}
}
