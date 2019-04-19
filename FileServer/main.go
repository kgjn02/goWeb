package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("src/upload.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseMultipartForm(32 << 20) // 指定缓存大小，若文件大小超过缓存，则保存在系统临时文件
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile(sharePath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

var curPath = ""
var sharePath = "share/"

func init() {
	err := os.Mkdir("share/", 0777)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) // os.Args[0]表示当前文件绝对路径
	if err != nil {
		log.Fatal(err)
	}
	curPath = dir
}

// 证书文件和公钥文件路径
var (
	cert = ""
	key  = ""
)

func main() {
	fh := http.FileServer(http.Dir("share"))
	http.Handle("/", fh)
	http.HandleFunc("/upload", upload)
	//err := http.ListenAndServe(":9090", nil) // 不使用 https
	err := http.ListenAndServeTLS(":9090", cert, key, nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
