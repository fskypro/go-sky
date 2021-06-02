/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: rpc service interface
@author: fanky
@version: 1.0
@date: 2021-03-20
**/

package fsrpcs

// -------------------------------------------------------------------
// rpc service must implement this interface
// -------------------------------------------------------------------
type I_Service interface {
	Name() string
}
