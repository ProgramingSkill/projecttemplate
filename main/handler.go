package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

//AddArticleHandler 添加文章
func AddArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var (
		reqData struct {
			Tag     string `json:"tag"`
			Title   string `json:"title"`
			Content string `json:"content"`
		}

		ret         int
		responseStr string

		respData struct {
			Status int         `json:"status"`
			Msg    string      `json:"msg"`
			Data   *ArticleStu `json:"data"`
		}
	)
	defer func() {
		if ret != Success {
			responseStr = fmt.Sprintf(`{"status":%v,"msg":"%v"}`, ret, errorMap[ret])
		}
		io.WriteString(w, responseStr)
	}()

	operateUser := r.Header.Get("operateuser")

	if ret = ParseRequest(r, &reqData); ret != Success {
		return
	}

	if reqData.Tag == "" || reqData.Title == "" || reqData.Content == "" {
		Error(ErrReqParam, errorMap[ErrReqParam])
		ret = ErrReqParam
		return
	}
	contentBytes, err := json.Marshal(reqData.Content)
	if err != nil {
		return
	}
	article := &ArticleStu{
		ArticleID:   GeneSignID(time.Now(), GetReplaceUUID()),
		Tag:         reqData.Tag,
		Title:       reqData.Title,
		Content:     Base64Encode(contentBytes),
		OperateUser: operateUser,
		CreateTime:  GetTimeNow(),
		UpdateTime:  GetTimeNow(),
	}

	if ret = mysqlDB.insertArticleNew(article); ret != Success {
		Error(ret, errorMap[ret])
		return
	}
	respData.Data = article
	respData.Msg = errorMap[ret]
	responseStr, ret = StructToJsonString(respData)

	return
}

func GetArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var (
		reqData struct {
			Tag      string `json:"tag"`
			Title    string `json:"title"`
			PageNum  int    `json:"pagenum"`
			PageSize int    `json:"pagesize"`
		}

		ret         int
		responseStr string

		respData struct {
			Status int          `json:"status"`
			Msg    string       `json:"msg"`
			Data   []ArticleStu `json:"data"`
			Total  int          `json:"total"`
		}
	)
	defer func() {
		if ret != Success {
			responseStr = fmt.Sprintf(`{"status":%v,"msg":"%v"}`, ret, errorMap[ret])
		}
		io.WriteString(w, responseStr)
	}()

	if ret = ParseRequest(r, &reqData); ret != Success {
		return
	}
	Debug("reqData:", reqData)

	if reqData.PageSize == 0 || reqData.PageNum == 0 {
		Error(ErrReqParam, errorMap[ErrReqParam])
		ret = ErrReqParam
		return
	}

	if respData.Data, respData.Total, ret = mysqlDB.queryArticleList(reqData.Tag, reqData.Title, reqData.PageSize, reqData.PageNum); ret != Success {
		Error(ret, errorMap[ret])
		return
	}

	respData.Msg = errorMap[ret]
	responseStr, ret = StructToJsonString(respData)
	return
}

//UpdateArticleHandler 更新
func UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var (
		reqData struct {
			ArticleID string `json:"articleid"`
			Tag       string `json:"tag"`
			Title     string `json:"title"`
			Content   string `json:"content"`
		}

		ret         int
		responseStr string

		respData struct {
			Status int    `json:"status"`
			Msg    string `json:"msg"`
		}
	)
	defer func() {
		if ret != Success {
			responseStr = fmt.Sprintf(`{"status":%v,"msg":"%v"}`, ret, errorMap[ret])
		}
		//traceKey := make(map[string]interface{})
		//bizhidata.UploadFormatData(r, "productadd", stime, body, responseStr, ret, errorMap[ret], Config.Upload.DataType, traceKey)
		io.WriteString(w, responseStr)
	}()

	operateUser := r.Header.Get("operateuser")

	if ret = ParseRequest(r, &reqData); ret != Success {
		return
	}

	if reqData.ArticleID == "" {
		Error(ErrReqParam, errorMap[ErrReqParam])
		ret = ErrReqParam
		return
	}

	session := mysqlDB.sqlXorm.NewSession()
	defer session.Close()
	session.Begin()

	contentBytes, err := json.Marshal(reqData.Content)
	if err != nil {
		return
	}
	article := ArticleStu{
		ArticleID:   reqData.ArticleID,
		Tag:         reqData.Tag,
		Title:       reqData.Title,
		Content:     Base64Encode(contentBytes),
		OperateUser: operateUser,
		IsDeleted:   0,
		UpdateTime:  GetTimeNow(),
	}

	if ret = mysqlDB.UpdateArticle(session, &article); ret != Success {
		Error(ret, errorMap[ret])
		return
	}
	session.Commit()
	respData.Msg = errorMap[ret]
	responseStr, ret = StructToJsonString(respData)
	return
}
