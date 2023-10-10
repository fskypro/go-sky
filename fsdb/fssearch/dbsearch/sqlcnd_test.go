package dbsearch

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"fsky.pro/fsserializer/jsonex"
	"fsky.pro/fsstr/fsfmt"
)

// 查询条件配置
const jstr = `
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

    // 可搜索字段定义
    "headers": {
        // 告警静态数据可搜索字段，可指定查询字段
        "get-alarms": [
			{"col": "alarmid", "key": "alarmID", "name": "故障码", "matchs": ["包含", "匹配"], "input": "textbox" },
			{"col": "lvname", "key": "lvName", "name": "告警等级", "matchs": ["包含", "匹配"], "input": "textbox" },
			{"col": "alarm_time", "key": "alarmTime", "name": "告警日期", "matchs": ["大于", "小于"], "input": "date" },
			{"col": "reason", "key": "reason", "name": "告警原因", "matchs": ["包含", "匹配"], "input": "textbox" },
			{"col": "native_reason", "key": "nativeReason", "name": "告警本地原因", "matchs": ["包含", "匹配"], "input": "textbox" },
			{"col": "tmf_reason", "key": "tmfReason", "name": "告警TMF标准原因", "matchs": ["包含", "匹配"], "input": "textbox" },
			{"col": "resolution", "key": "resolution", "name": "故障处理建议", "matchs": ["包含", "匹配"], "input": "textbox" },
			{"col": "nettype", "key": "netType", "name": "网络类型", "matchs": ["匹配"], "input": "textbox" }
        ],
	}
}
`

type S_HeaderItem struct {
	Col    string   `json:"col"`
	Key    string   `json:"key"`
	Name   string   `json:"name"`
	Matchs []string `json:"matchs"`
	Input  string   `json:"inputs"`
}

type T_Header []*S_HeaderItem

func (self T_Header) KeyToCol(key string) string {
	for _, item := range self {
		if item.Key == key {
			return item.Col
		}
	}
	return ""
}

type S_Config struct {
	Matchs  map[string]string `json:"matchs"`
	Headers map[string]T_Header
}

func (this *S_Config) GetMatchString(mName string) string {
	return this.Matchs[mName]
}

func (this *S_Config) GetHeader(reqPath string) T_Header {
	return this.Headers[reqPath]
}

var config *S_Config

func init() {
	config = new(S_Config)
	jsonex.Unmarshal([]byte(jstr), config)
}

func TestJsonSQLCondtion(t *testing.T) {
	jcnds := `[
		{"match":"匹配", "col": "alarmID", "value": "xxx"},
		["and",
			{"match":"小于", "col": "lvName", "value": "200"},
			{"match":"包含", "col": "reason", "value": "yyy"}
		],
		["or",
			{"match":"匹配", "col": "netType", "value": "zzz"},
			{"match":"大于", "col": "alarmTime", "value": "10"}
		]
	]`
	//jcnds = `[]`

	anyCnd := []any{}
	err := json.Unmarshal([]byte(jcnds), &anyCnd)
	if err != nil {
		panic(err)
	}
	cnd, err := NewCnds(anyCnd)
	if err != nil {
		panic(err)
	}

	var sb strings.Builder
	cnd.Tran(func(e E_CndElem, cnd *S_Cnd) error {
		switch e {
		case ECndLQuote:
			sb.WriteByte('(')
		case ECndRQuote:
			sb.WriteByte(')')
		case ECndAnd:
			sb.WriteString(" AND ")
		case ECndOr:
			sb.WriteString(" OR ")
		case ECndObj:
			if cnd.Type == CndEmpty {
				sb.WriteString("TRUE")
				return nil
			}
			if cnd.Type == CndStr {
				if cnd.String == "TRUE" || cnd.String == "FALSE" {
					sb.WriteString(cnd.String)
					return nil
				}
				return fmt.Errorf("error condition %q", cnd.String)
			}
			// 假设：
			//   cnd.Match = "等于"
			//   cnd.Col = "key"
			//   cnd.Value = "xxx"
			// 则：
			//   mstr := fsfmt.Smprintf(cnd.Match, map[string]any{"c": cnd.Col, "v": cnd.Value})
			header := config.GetHeader("get-alarms")
			mstr := fsfmt.Smprintf(config.GetMatchString(cnd.Match), map[string]any{"c": header.KeyToCol(cnd.Key), "v": cnd.Value})
			sb.WriteString(mstr)
		}
		return nil
	})
	fmt.Println(sb.String())
}
