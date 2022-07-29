package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
	"xorm.io/core"
	"xorm.io/xorm"
)

type MysqlDB struct {
	sqlDB   *sql.DB
	sqlXorm *xorm.Engine
}

var mysqlDB = new(MysqlDB)

func mysqlInit() {
	mysqlDB.InitDB()
}

func (m *MysqlDB) InitDB() {
	mysqlPara := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&allowNativePasswords=true&tls=preferred",
		Config.MySql["1"].User,
		Config.MySql["1"].PassWord,
		Config.MySql["1"].IP,
		Config.MySql["1"].Port,
		Config.MySql["1"].DBname,
	)
	m.connect(mysqlPara)

	if m.syncTable() != Success {
		Error("err mysql sync table")
		os.Exit(-1)
	}

	return
}

func (m *MysqlDB) connect(mysqlAddr string) {
	var err error
	m.sqlXorm, err = xorm.NewEngine("mysql", mysqlAddr)
	if err != nil {
		Error("open mysql xorm failed:", err)
		os.Exit(-1)
	}
	for {
		if err := m.sqlXorm.Ping(); err != nil {
			Warning("mysql xorm failed:", err)
			time.Sleep(1 * time.Second)
			continue
		} else {
			Debug("mysql xorm ping ok")
		}
		break
	}
	m.sqlXorm.SetMapper(core.SameMapper{})

	m.sqlDB, err = sql.Open("mysql", mysqlAddr)
	if err != nil {
		Error("sql.Open:", err)
		os.Exit(-1)
	}

	m.sqlDB.SetMaxOpenConns(30)
	m.sqlDB.SetMaxIdleConns(30)

	for {
		if err := m.sqlDB.Ping(); err != nil {
			Error("mysql.Ping:", err)
			time.Sleep(1 * time.Second)
			continue
		} else {
			Debug("connect mysql server ok")
		}
		break
	}
}

func (m *MysqlDB) syncTable() int {
	var ret int
	err := m.sqlXorm.Sync2(new(ArticleStu))
	if err != nil {
		Error(err)
		return ErrInnerFault
	}
	ret = Success
	return ret
}

func (m *MysqlDB) insertArticleNew(article *ArticleStu) (ret int) {
	_, err := m.sqlXorm.Insert(article)
	if err != nil {
		Error("fail insert new article.", err)
		ret = ErrInnerFault
		return
	}
	return
}

func (m *MysqlDB) queryArticleList(tag, title string, pagesize, pagenum int) (articles []ArticleStu, total, ret int) {
	session := m.sqlXorm.Where("isdeleted=?", 0)
	// like
	if tag != "" {
		session.And("tag like ?", "%"+tag+"%")
	}
	if title != "" {
		session.And("title like ?", "%"+title+"%")
	}
	limit := pagesize
	offset := (pagenum - 1) * pagesize
	Debug("limit:", limit, "offset:", offset)
	articles = make([]ArticleStu, 0)
	// pgae limit offset
	total64, err := session.Asc("articleid").Limit(limit, offset).FindAndCount(&articles)
	total = int(total64)
	if err != nil {
		Error("fail query product list.", err)
		ret = ErrInnerFault
		return
	}
	Debug(articles)
	return
}

func (m *MysqlDB) UpdateArticle(session *xorm.Session, article *ArticleStu) (ret int) {
	exist, err := session.Where("articleid=? AND isdeleted=?", article.ArticleID, 0).Get(&ArticleStu{})
	if err != nil {
		session.Rollback()
		Error(err)
		return ErrInnerFault
	}
	if !exist {
		session.Rollback()
		return ErrNotExist
	}
	if _, err := m.sqlXorm.Where("articleid=? AND isdeleted=?", article.ArticleID, 0).Cols("tag", "title", "content", "updatetime").Update(article); err != nil {
		session.Rollback()
		Error("fail update article.", err)
		ret = ErrInnerFault
		return
	}
	return
}
