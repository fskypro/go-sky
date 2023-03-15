/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 实现带版本号的表
@author: fanky
@version: 1.0
@date: 2022-07-08
**/

// 带版本号的表
// 每个表中多增加一个列，用于表示表的版本，列名称格式为：_v_<本号>
// 版本号是一个大于 0 的 uint32

package fsmysql

import (
	"fmt"
	"strconv"

	"fsky.pro/fsmysql/fssql"
)

type F_TBUpper func(*S_Tx) *S_OPResult

// -----------------------------------------------------------------------------
// private
// -----------------------------------------------------------------------------
// 增加表版本列
func (this *S_Tx) addTableVersion(table *fssql.S_Table, version int) *S_OPExecResult {
	sqlInfo := table.AddColumnSQL("_v_"+strconv.Itoa(version), "BIT", "COMMENT 'table version' FIRST")
	return this.ExecSQLInfo(sqlInfo)
}

// 获取表格版本号
func (this *S_Tx) queryTBVersion(table *fssql.S_Table) *S_OPValueResult {
	rest := this.FetchColumns(table, "^_v_[:digit:]+")
	if rest.Err() != nil {
		return newOPValueResult(rest.sqlInfo.(*fssql.S_FetchInfo), rest.Err())
	}

	version := 0
	cols := rest.Value.([]string)
	if len(cols) > 0 {
		version, _ = strconv.Atoi(cols[0][3:])
	}
	return newOPValueResult2(rest.sqlInfo.(*fssql.S_FetchInfo), version)
}

// 修改表版本号
func (this *S_Tx) modTBVersion(table *fssql.S_Table, oldv, newv int) *S_OPResult {
	oldColumn := fmt.Sprintf("_v_%d", oldv)
	newColumn := fmt.Sprintf("_v_%d", newv)
	return this.RenameColumn(table, oldColumn, newColumn, "BIT", "COMMENT 'table version'")
}

// -----------------------------------------------------------------------------
// public
// -----------------------------------------------------------------------------
// 创建带版本号的表
// vstart 表示从哪个版本开始叠加
func (this *S_DB) CreateVersionTable(table *fssql.S_Table, vstart int, ups []F_TBUpper) *S_OPResult {
	tx, err := this.Begin()
	if err != nil {
		return newOPResult(table.CreateSQL(), err)
	}
	defer tx.Rollback()

	rest := tx.HasTable(table.Name())
	if rest.Err() != nil {
		return newOPResult(rest.sqlInfo, rest.Err())
	}

	if vstart < 1 {
		vstart = 1
	}
	newVersion := len(ups) + vstart

	// 表在数据库中还不存在
	if !rest.Value.(bool) {
		rest := tx.CreateTable(table)
		if rest.Err() != nil {
			return rest
		}
		if r := tx.addTableVersion(table, newVersion); r.Err() != nil {
			return newOPResult(r.sqlInfo, fmt.Errorf("add version(%d) to version table fail, %v", newVersion, r.Err()))
		}
		err = tx.Commit()
		if err != nil {
			return newOPResult(rest.sqlInfo, err)
		}
		return rest
	}

	// 表已经存在，查找旧版本号
	rest = tx.queryTBVersion(table)
	if rest.Err() != nil {
		return newOPResult(rest.sqlInfo, rest.Err())
	}
	oldVersion := rest.Value.(int)

	// 表格存在，但是版本号不存在
	if oldVersion == 0 {
		return newOPResult(table.CreateSQL(), fmt.Errorf("table %q exists, but has no version column", table.Name()))
	}
	// 旧版比新版还新
	if oldVersion > newVersion {
		return newOPResult(table.CreateSQL(), fmt.Errorf("old version(%d) in table %q is largs than the new version(%d)", oldVersion, table.Name(), newVersion))
	}
	if oldVersion == newVersion {
		return newOPResult(table.CreateSQL(), nil)
	}

	// 对所有未更新的执行一次更新操作
	for idx, up := range ups {
		v := idx + vstart + 1
		if v <= oldVersion {
			continue
		}
		if rest := up(tx); rest.Err() != nil {
			return rest
		}
	}

	if r := tx.modTBVersion(table, oldVersion, newVersion); r.Err() != nil {
		return r
	}
	err = tx.Commit()
	if err != nil {
		return newOPResult(table.CreateSQL(), err)
	}
	return newOPResult(table.CreateSQL(), nil)
}
