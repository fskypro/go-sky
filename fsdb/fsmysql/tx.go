/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: tx
@author: fanky
@version: 1.0
@date: 2021-12-31
**/

package fsmysql

import "database/sql"

type S_Tx struct {
	*s_Operator
	*sql.Tx
}

func newTx(tx *sql.Tx) *S_Tx {
	return &S_Tx{
		s_Operator: newOperator(tx),
		Tx:         tx,
	}
}
