package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func startServer() {
	http.HandleFunc("/projecttemplate/healthcheck", healthCheckHandler)
	http.HandleFunc("/projecttemplate/api/v1/article/add", AddArticleHandler)
	http.HandleFunc("/projecttemplate/api/v1/article/get", GetArticleHandler)
	http.HandleFunc("/projecttemplate/api/v1/article/update", UpdateArticleHandler)

	var s = &http.Server{
		Addr:           Config.Port,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if nil != err {
		Error("fail to listen and server :", err)
		os.Exit(-1)
	}
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		return
	}
	var responseStr string
	responseStr = fmt.Sprintf(`{"error":%d,"status":%d,"msg":"%s"}`, Success, Success, errorMap[Success])
	io.WriteString(w, responseStr)
	return
}
