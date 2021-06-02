package test

import "fmt"

type DD struct {
	Value int
}

type Obj struct {
	Aa  string
	bb  string
	Pcc *string
	pcc *string

	DD    DD
	dd    DD
	pdd   *DD
	nildd *DD
}

func (this *Obj) GetBB() string {
	return this.bb
}

func (this *Obj) GetPCC() *string {
	return this.pcc
}

func (this *Obj) GetDD() DD {
	return this.dd
}

func (this *Obj) GetPDD() *DD {
	return this.pdd
}

func (this *Obj) GetNilDD() *DD {
	return this.nildd
}

func (this *Obj) call(a int, b string) string {
	return fmt.Sprintf("%d: %s", a, b)
}

func NewObj() *Obj {
	Pcc := "3333"
	pcc := "4444"
	pdd := &DD{300}
	return &Obj{
		Aa:    "1111",
		bb:    "2222",
		Pcc:   &Pcc,
		pcc:   &pcc,
		DD:    DD{100},
		dd:    DD{200},
		pdd:   pdd,
		nildd: nil,
	}
}
