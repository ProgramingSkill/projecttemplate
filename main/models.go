package main

type ArticleStu struct {
	ArticleID   string `json:"articleid" xorm:"varchar(40) pk 'articleid'"`           // 文章ID
	Tag         string `json:"tag" xorm:"varchar(40) 'tag'"`                          // 标签
	Title       string `json:"title" xorm:"varchar(40) 'title'"`                      // 标题
	Content     string `json:"msgtext" xorm:"varchar(40) 'msgtext' NOT NULL"`         // 内容
	IsDeleted   int    `json:"isdeleted" xorm:"'isdeleted' int default 0"`            // 是否删除: 0,未删除;1,已删除
	OperateUser string `json:"operateuser" xorm:"varchar(40) NOT NULL 'operateuser'"` // 操作人
	CreateTime  string `json:"createtime" xorm:"'createtime' NOT NULL default ''"`    // 创建时间
	UpdateTime  string `json:"updatetime" xorm:"'updatetime' NOT NULL default ''"`    // 更新时间
}

func (ArticleStu) TableName() string {
	return "article"
}
