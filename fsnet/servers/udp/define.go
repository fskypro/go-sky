/**
@copyright: fantasysky 2016
@brief: udp server definations
@author: fanky
@version: 1.0
@date: 2019-12-21
**/

package udp

import "net"

type F_Receive func(*net.UDPAddr, []byte)
