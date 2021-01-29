/**
@copyright: fantasysky 2016
@brief: implement xml attribute
@author: fanky
@version: 1.0
@date: 2020-02-20
**/

package fsxml

type S_Attr struct {
	s_NameText
}

// 新建一个属性
// 注意：如果属性名称不合法，则返回 nil
func NewAttr(name, value string) *S_Attr {
	if !isValidName(name) {
		return nil
	}
	return &S_Attr{
		s_NameText: s_NameText{name, value},
	}
}

func (this *S_Attr) clone() *S_Attr {
	return &S_Attr{
		s_NameText: s_NameText{this.name, this.text},
	}
}
