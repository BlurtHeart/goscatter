package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	dirname, err := ioutil.ReadDir("./")
	CheckErrorOnExit(err)
	for k, v := range dirname {
		fmt.Println(k, "=", v.Name())
		fmt.Println(v.IsDir())
		fmt.Println(v.ModTime())
		fmt.Println(v.Mode())
		fmt.Println(v.Size())
		fmt.Println(v.Sys())
	}

	byte1, err := ioutil.ReadFile("testpost.go")
	CheckErrorOnExit(err)
	fmt.Println(string(byte1))

	reader := strings.NewReader("hello world")
	byte2, err := ioutil.ReadAll(reader) //输入一个io.Reader元，返回的是一个[]byte
	CheckErrorOnExit(err)
	fmt.Println(string(byte2))
	// nopcloser
	reader = strings.NewReader("你好世界")
	nc := ioutil.NopCloser(reader) //读取一个io.Reader元，返回的是一个io.ReadClose接口，提供Close方法
	defer nc.Close()
	byte2, err = ioutil.ReadAll(nc)
	CheckErrorOnExit(err)
	fmt.Println(string(byte2))

	name, err := ioutil.TempDir("../", "temp") //读取一个目录，返回的是prefix+随机数字的临时目录，同时会创建这个目录
	CheckErrorOnExit(err)
	fmt.Println(name)
	os.Remove(name)

	err = ioutil.WriteFile("./test.txt", []byte("hello world"), 0644) //向一个文件写入数据，如果没有根据fileMode创建一个,清空文件后写入
	CheckErrorOnExit(err)
}

func CheckErrorOnExit(err error) {
	if err != nil {
		panic(err)
	}
}
