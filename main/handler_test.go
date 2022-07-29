package main

import (
	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func BasicEngine() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/projecttemplate/api/v1/article/add", AddArticleHandler)
	mux.HandleFunc("/projecttemplate/api/v1/article/get", GetArticleHandler)
	mux.HandleFunc("/projecttemplate/api/v1/article/update", UpdateArticleHandler)
	return mux
}

// product

var (
	companyId   = "7316a95cec3541b297a5150a27bea986"
	operateUser = "admin"
	pageNum     = 1
	pageSize    = 10
	productId   = "2022072200003850271"
	apikeyId    = "2022072296008851125"
)

func TestAddArticleHandler(t *testing.T) {
	r := gofight.New()

	r.POST("/dnasnlicense/api/v1/article/add").
		SetHeader(gofight.H{
			"operateuser": operateUser,
		}).
		SetJSON(gofight.D{
			"tag":     "编程",
			"title":   "编程珠玑",
			"content": "文章内容",
		}).
		Run(BasicEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			status := int(gjson.Get(string(data), "status").Int())
			msg := gjson.Get(string(data), "msg").String()
			Debug(r.Body.String(), r.Code)
			assert.Equal(t, 0, status)
			assert.Equal(t, "success", msg)
			assert.Equal(t, "{\"status\":0,\"msg\":\"success\",\"data\":{\"articleid\":\"2022072960002451103\",\"tag\":\"编程\",\"title\":\"编程珠玑\",\"msgtext\":\"IuaWh+eroOWGheWuuSI=\",\"isdeleted\":0,\"operateuser\":\"1dsfjsbfjsb23fbjw\",\"createtime\":\"2022-07-29 23:37:32\",\"updatetime\":\"2022-07-29 23:37:32\"}}", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "application/json; charset=utf-8", rq.Header.Get("Content-Type"))
		})
}

func TestGetArticleHandler(t *testing.T) {
	r := gofight.New()

	r.POST("/dnasnlicense/api/v1/article/get").
		SetJSON(gofight.D{
			"tag":      "",
			"title":    "编程珠玑",
			"pagenum":  "1",
			"pagesize": "10",
		}).
		Run(BasicEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			status := int(gjson.Get(string(data), "status").Int())
			msg := gjson.Get(string(data), "msg").String()
			Debug(r.Body.String(), r.Code)
			assert.Equal(t, 0, status)
			assert.Equal(t, "success", msg)
			assert.Equal(t, "{\"status\":0,\"msg\":\"success\",\"data\":[{\"articleid\":\"2022072900003050647\",\"tag\":\"编程\",\"title\":\"编程珠玑2\",\"msgtext\":\"IuaWh+eroOWGheWuuTIi\",\"isdeleted\":0,\"operateuser\":\"1dsfjsbfjsb23fbjw\",\"createtime\":\"2022-07-29 23:39:48\",\"updatetime\":\"2022-07-29 23:39:48\"},{\"articleid\":\"2022072960002451103\",\"tag\":\"编程\",\"title\":\"编程珠玑\",\"msgtext\":\"IuaWh+eroOWGheWuuSI=\",\"isdeleted\":0,\"operateuser\":\"1dsfjsbfjsb23fbjw\",\"createtime\":\"2022-07-29 23:37:32\",\"updatetime\":\"2022-07-29 23:37:32\"}],\"total\":2}", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "application/json; charset=utf-8", rq.Header.Get("Content-Type"))
		})
}

func TestUpdateArticleHandler(t *testing.T) {
	r := gofight.New()

	r.POST("/dnasnlicense/api/v1/article/update").
		SetHeader(gofight.H{
			"operateuser": operateUser,
		}).
		SetJSON(gofight.D{
			"articleid": "xxxxx",
			"tag":       "编程",
			"title":     "编程珠玑",
			"content":   "文章内容",
		}).
		Run(BasicEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			status := int(gjson.Get(string(data), "status").Int())
			msg := gjson.Get(string(data), "msg").String()
			Debug(r.Body.String(), r.Code)
			assert.Equal(t, 0, status)
			assert.Equal(t, "success", msg)
			assert.Equal(t, "{\"status\":0,\"msg\":\"success\"}", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "application/json; charset=utf-8", rq.Header.Get("Content-Type"))
		})
}
