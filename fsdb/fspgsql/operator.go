/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: db operator
@author: fanky
@version: 1.0
@date: 2023-02-01
**/

package fspgsql

import "database/sql"

// -----------------------------------------------------------------------------
// DB/Tx Wraper
// -----------------------------------------------------------------------------
type i_DBWrapper interface {
	Exec(string, ...any) (sql.Result, error)
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
	Prepare(string) (*sql.Stmt, error)
}

// -----------------------------------------------------------------------------
// Operator
// -----------------------------------------------------------------------------
type s_Operator struct {
	wrapper i_DBWrapper
}

func newOperator(wrapper i_DBWrapper) *s_Operator {
	return &s_Operator{wrapper}
}
