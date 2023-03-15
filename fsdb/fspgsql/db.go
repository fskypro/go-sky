/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief:
@author: fanky
@version: 1.0
@date: 2023-02-01
**/

package fspgsql

import "database/sql"

type S_DB struct {
	*s_Operator
	*sql.DB // 连接串
	DBInfo  *S_DBInfo
}
