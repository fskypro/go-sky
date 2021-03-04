/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: aync print log
@author: fanky
@version: 1.0
@date: 2021-02-28
**/

package fslog

func ADebug(vs ...interface{}) {
	go _logger.print(1, "debug", vs)
}

func ADebugf(format string, vs ...interface{}) {
	go _logger.printf(1, "debug", format, vs)
}

func AInfo(vs ...interface{}) {
	go _logger.print(1, "info", vs)
}

func AInfof(format string, vs ...interface{}) {
	go _logger.printf(1, "info", format, vs)
}

func AWarn(vs ...interface{}) {
	go _logger.print(1, "warn", vs)
}

func AWarnf(format string, vs ...interface{}) {
	go _logger.printf(1, "warn", format, vs)
}

func AError(vs ...interface{}) {
	go _logger.print(1, "error", vs)
}

func AErrorf(format string, vs ...interface{}) {
	go _logger.printf(1, "error", format, vs)
}

func AHack(vs ...interface{}) {
	go _logger.print(1, "hack", vs)
}

func AHackf(format string, vs ...interface{}) {
	go _logger.printf(1, "hack", format, vs)
}

func ATrace(vs ...interface{}) {
	go _logger.printChain(1, "trace", vs)
}

func ATracef(format string, vs ...interface{}) {
	go _logger.printChainf(1, "trace", format, vs)
}
