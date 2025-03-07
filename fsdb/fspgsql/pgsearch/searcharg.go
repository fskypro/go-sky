/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 搜索参数
@author: fanky
@version: 1.0
@date: 2023-10-01
**/

package pgsearch

import (
	"errors"
	"fmt"
	"strings"

	"fsky.pro/fsky"
	"fsky.pro/fssearch/dbsearch"
)

/* 定义配置示例
{
    // 匹配方式定义
    "matchs": {
        // %[c]s 表示数据库字段名称
        // %[v]s 表示传入值
        "包含": "%[c]s like '%%%[v]s%%'",
        "不包含": "%[c]s not like '%%%[v]s%%'",
        "匹配": "%[c]s='%[v]s'",
        "不匹配": "%[c]s<>'%[v]s'",

        "等于": "%[c]s=%[v]s",
        "大于": "%[c]s>%[v]s",
        "小于": "%[c]s<%[v]s",
        "大于等于": "%[c]s>=%[v]s",
        "小于等于": "%[c]s<=%[v]s",
        "不等于": "%[c]s<>%[v]s"
    },
    // 顶级字段对应请求 URL path
    "headers": {
        // 滚动告警列表
        "alarm/get-roll-alarms": [
            {"col": "net_type",    "key": "netType",     "name": "-",            "matchs": [], "input": "textbox" },
            {"col": "",            "key": "netTypeName", "name": "网络类型",     "matchs": [], "input": "textbox" },
            {"col": "equip_name",  "key": "equipName",   "name": "网元名称",     "matchs": [], "input": "textbox" },
            {"col": "alarm_id",    "key": "alarmID",     "name": "告警ID",       "matchs": [], "input": "textbox" },
            {"col": "serial",      "key": "serial",      "name": "告警流水号",   "matchs": [], "input": "textbox" },
            {"col": "level",       "key": "level",       "name": "告警等级",     "matchs": [], "input": "textbox" , "vmap": {"4": "紧急", "3": "重要", "2": "次要", "1": "提示", "0": "未知"} },
            {"col": "latest_time", "key": "latestTime",  "name": "最新告警时间", "matchs": [], "input": "date" },
            {"col": "reason",      "key": "reason",      "name": "告警原因",     "matchs": [], "input": "textbox" },
        ],
	}
}
*/

type I_Config interface {
	// 获取匹配字符串，对应上示例配置中的 matchs 的 value
	// key  ：条件表达式中的匹配字段
	// match：条件表达式中的匹配关键字（如：包含、匹配、不匹配,,,,）
	// value：要被搜索的字符串
	GetMatchString(key, match, value string) (string, error)

	// 通过 key 获取对应的数据库字段
	GetDBColumn(string) string

	// 获取默认分页大小
	DefaultPageSize() int
}

// -------------------------------------------------------------------
// SearchArg
// -------------------------------------------------------------------
// 通用搜索参数
type S_SearchArg struct {
	Cnds     []any  `json:"cnds"`     // 用户自定义的原始搜索条件
	Page     int    `json:"page"`     // 要查询的页码（第一页为 1）
	PageSize int    `json:"pageSize"` // 每页最大数量
	OrderBy  string `json:"orderBy"`  // 排序关键字
	Desc     int8   `json:"desc"`     // 是否倒序排列

	conf   I_Config
	cndExp string
	parsed bool
}

func NewSearchArg(conf I_Config) *S_SearchArg {
	return &S_SearchArg{
		Cnds:     make([]any, 0),
		Page:     1,
		PageSize: conf.DefaultPageSize(),
		conf:     conf,
	}
}

// 解释搜索条件
func (this *S_SearchArg) ParseCnd() error {
	this.parsed = true
	if this.Cnds == nil {
		return nil
	}
	cnds, err := dbsearch.NewCnds(this.Cnds)
	if err != nil {
		return fmt.Errorf("parse search condition(%#v) fail, %v", this.Cnds, err)
	}
	if cnds == nil {
		return errors.New("no conditions in search arg")
	}

	// 生成条件表达式
	var sb strings.Builder
	err = cnds.Tran(func(e dbsearch.E_CndElem, cnd *dbsearch.S_Cnd) error {
		switch e {
		case dbsearch.ECndLQuote:
			sb.WriteByte('(')
		case dbsearch.ECndRQuote:
			sb.WriteByte(')')
		case dbsearch.ECndAnd:
			sb.WriteString(" AND ")
		case dbsearch.ECndOr:
			sb.WriteString(" OR ")
		case dbsearch.ECndObj:
			if cnd.Type == dbsearch.CndEmpty {
				sb.WriteString("TRUE")
				return nil
			}
			if cnd.Type == dbsearch.CndObj {
				mstr, err := this.conf.GetMatchString(cnd.Key, cnd.Match, cnd.Value)
				if err != nil {
					return fmt.Errorf("format condition string fail, %v", err)
				}
				sb.WriteString(mstr)
				return nil
			}
			str := strings.ToLower(cnd.String)
			if str == "true" || str == "false" {
				sb.WriteString(strings.ToUpper(cnd.String))
				return nil
			}
			return fmt.Errorf("string %q is not a sql condition expression", cnd.String)
		}
		return nil
	})
	if err != nil {
		return err
	}
	this.cndExp = sb.String()
	if this.cndExp == "" {
		this.cndExp = "TRUE"
	}
	this.OrderBy = this.conf.GetDBColumn(this.OrderBy)
	return nil
}

func (this *S_SearchArg) PageIndex() int {
	if this.PageSize < 1 {
		return 0
	}
	return (this.Page - 1) * this.PageSize
}

func (this *S_SearchArg) Order() string {
	return fsky.IfElse(this.Desc > 0, "DESC", "ASC")
}

func (this *S_SearchArg) Where() string {
	if this.parsed {
		return this.cndExp
	}
	return "<unparsed search argument>"
}

func (this *S_SearchArg) WhereAndTail() string {
	if !this.parsed {
		return "<unparsed search argument>"
	}
	if this.PageSize < 1 {
		if this.OrderBy == "" {
			return this.cndExp
		}
		return fmt.Sprintf(`%s ORDER BY "%s" %s`, this.cndExp, this.OrderBy, this.Order())
	}
	if this.OrderBy == "" {
		return fmt.Sprintf("%s LIMIT %d OFFSET %d", this.cndExp, this.PageSize, this.PageIndex())
	}
	return fmt.Sprintf(`%s ORDER BY "%s" %s LIMIT %d OFFSET %d`, this.cndExp, this.OrderBy, this.Order(), this.PageSize, this.PageIndex())
}
