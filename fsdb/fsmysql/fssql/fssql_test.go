package fssql

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

type S_User struct {
	UserID int64  `db:"uid"`
	Name   string `name:"name"`
}
type S_Order struct {
	OrderID string `db:"order_id" tbtd:"varchar(64) not null default ''"`
	UserID  int64  `db:"uid" dbtd:"bigint not null default '0'"`
	Value   int    `db:"value" dbtd:"int not null default '0'"`
}

type S_UserOrder struct {
	UserID  int64  `db:"user.uid"`
	Name    string `db:"user.name"`
	OrderID int64  `db:"order.order_id"`
	Value   int    `db:"order.value"`
}

var tbUser *S_Table
var tbOrder *S_Table
var tbUOrder *S_Table

func init() {
	var err error
	tbUser, err = NewTable("user", new(S_User))
	fmt.Println("error: ", err)
	tbOrder, err = NewTable("order", new(S_Order))
	fmt.Println("error: ", err)
	tbUOrder, err = NewLinkTable(new(S_UserOrder),
		"#[1] JOIN #[2] ON $[3]=$[4]", tbUser, tbOrder, tbUser.M("UserID"), tbOrder.M("UserID"))
	fmt.Println("error: ", err)
}

func TestCreateTable(t *testing.T) {
	fstest.PrintTestBegin("CreateTable")
	defer fstest.PrintTestEnd()
	tbOrder.AddSchemes("PRIMARY KEY (`UserID`)", "UNIQUE KEY `UserID` (`UserID`) USING BTREE")
	fmt.Println(tbOrder.CreateTableSQLInfo().FmtSQLText("  "))
}

func TestSelect(t *testing.T) {
	fstest.PrintTestBegin("Select")
	defer fstest.PrintTestEnd()

	// select members
	sqlInfo := Select("OrderID").
		From(tbOrder).
		Where("$[1]>?[2]", "Value", 100).
		End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
		fmt.Println("--------------------")
	}

	// to object
	sqlInfo = SelectBesides("UserID").
		From(tbOrder).
		Where("$[1]=$[2] and $[1] in ?[3]", tbOrder.M("UserID"), tbUser.M("UserID"), []string{"11111", "22222"}).
		View("limit ?[1]", 10).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
		fmt.Println("--------------------")
	}

	// order join user
	sqlInfo = SelectAll().From(tbUOrder).
		Where("$[1]=$[2]", tbOrder.M("UserID"), tbUser.M("UserID")).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
		fmt.Println("--------------------")
	}

	// join
	sqlInfo = SelectAll().
		From(tbUOrder).
		Where("$[1] in ?[2]", "OrderID", []string{"11111", "22222"}).
		View("limit ?[1]", 10).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
		fmt.Println("--------------------")
	}

	sqlInfo = SelectExp("max($[1])+?[2]", "Value", 100).
		From(tbOrder).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}
}

func TestInsert(t *testing.T) {
	fstest.PrintTestBegin("Insert")
	defer fstest.PrintTestEnd()

	sqlInfo := Insert(tbOrder, "OrderID", "UserID", "Value").
		Values([]interface{}{"xxxx", 1000, 2000}).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}

	fmt.Println("--------------------")
	order := &S_Order{"xxxx", 10000, 20000}
	sqlInfo = Insert(tbOrder, "OrderID", "Value").
		Objects(order).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}

	fmt.Println("--------------------")
	sqlInfo = InsertAll(tbOrder).
		Objects(order).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}

	fmt.Println("--------------------")
	sqlInfo = InsertBesides(tbOrder, "UserID").
		Objects(order).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}

	fmt.Println("--------------------")
	sqlInfo = InsertIgnore(tbOrder, "UserID").
		Objects(order).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}
}

func TestUpdate(t *testing.T) {
	fstest.PrintTestBegin("Update")
	defer fstest.PrintTestEnd()

	sqlInfo := Update(tbOrder, "OrderID", "UserID", "Value").
		Set("xxxx", 1000, 2000).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}

	fmt.Println("1 --------------------")
	order := &S_Order{"xxxx", 10000, 20000}
	sqlInfo = Update(tbOrder, "OrderID", "Value").
		SetObject(order).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}

	fmt.Println("2 --------------------")
	sqlInfo = UpdateAll(tbOrder).
		SetObject(order).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}

	fmt.Println("3 --------------------")
	sqlInfo = UpdateBesides(tbOrder, "UserID").
		SetObject(order).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}

	fmt.Println("4 --------------------")
	sqlInfo = Update(tbOrder, "Value").
		SetExp("$[1]=$[1]+?[2]", "Value", 1000).
		Where("$[1] like ?[2]", "OrderID", "%xxxxxxx%").End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}

	fmt.Println("5 --------------------")
	sqlInfo = UpdateTables(tbOrder, tbUser).
		SetExp("$[1]=$[2]", tbOrder.M("UserID"), tbUser.M("UserID")).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}
}

func TestDelete(t *testing.T) {
	fstest.PrintTestBegin("Delete")
	defer fstest.PrintTestEnd()

	sqlInfo := Delete(tbOrder).
		Where("$[1]=?[2]", "OrderID", "xxxx").End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}
}

func TestInsertUpdate(t *testing.T) {
	fstest.PrintTestBegin("InsertUpdate")
	defer fstest.PrintTestEnd()

	order := &S_Order{"YYYY", 10000, 5000}
	sqlInfo := InsertOrUpdate(tbOrder, "OrderID", "Value").
		Values("xxxx", 1000).
		OrUpdate("OrderID", "UserID", "Value").
		WithObject(order).End()
	if sqlInfo.Err() != nil {
		fmt.Println(sqlInfo.Err())
	} else {
		fmt.Println(sqlInfo.SQLText())
		fmt.Println("sqltx:\n", sqlInfo.FmtSQLText("  "))
	}
}
