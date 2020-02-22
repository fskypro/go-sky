1、internal/xml 包来自于 go-1.13.4 中的 src/encoding/xml 包
2、internal/xml/xml.go 复制自 src/encoding/xml/xml.go，并修改如下：
    1）删除了函数（第 1945 行）：
        func (p *printer) EscapeString(s string) {
            ....
        }

    2）结构 Name 增加了 SpaceName 字段（第 45 行）：
        type Name struct {
            Space, Local string
            Space        string   // 这是增加的字段
        }
