package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/test", getJsonTest)
	err := http.ListenAndServe(":8003", nil)
	if err != nil {
		log.Println("http server listen :", err)
	}
}

func getJsonTest(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal("parse form error ", err)
	}
	// 初始化请求变量结构
	formData := make(map[string]interface{})

	// 调用json包的解析，解析请求body
	json.NewDecoder(r.Body).Decode(&formData)
	for key, value := range formData {
		log.Println("key:", key, " => value :", value)
	}

	// 返回json字符串给客户端
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(formData)
}
