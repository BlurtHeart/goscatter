// testhttp.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"io"
	"os"
	//	"strings"
)

type MyMux struct {
}

type Param struct {
	UrlLong string `json:"url_long"`
	ATest   int    `json:"A"`
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && r.Method == "POST" {
		test(w, r)
		return
	} else if r.URL.Path == "/upload" && r.Method == "POST" {
		upload(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
	fmt.Printf("body:%v", r.Body)
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	var p Param
	if err := json.Unmarshal(body, &p); err == nil {
		fmt.Println(p)
		p.ATest += 10
		fmt.Println(p)
		ret, _ := json.Marshal(p)
		fmt.Fprint(w, string(ret))
	} else {
		fmt.Println(err)
	}
}

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%s upload ok", handler.Filename)
	fmt.Printf("%s upload ok", handler.Filename)
	f, err := os.OpenFile("/home/active3/kwcs_buffer/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}

func main() {
	mux := &MyMux{}

	err := http.ListenAndServe(":18080", mux) // 设置监听端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
