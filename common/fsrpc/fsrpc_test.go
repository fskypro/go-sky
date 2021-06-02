package fsrpc

import "testing"

func TestIsEmptyArg(t *testing.T) {
    a := EArg{}
    b := &a
    t.Log("isEmptyArg(nil):", IsEmptyArg(nil))
    t.Log("isEmptyArg(EArg{}): ", IsEmptyArg(a))
    t.Log("isEmotyArg(&EArg{}): ", IsEmptyArg(b))
}

func TestIsEmptyReply(t *testing.T) {
    a := EReply{}
    b := &a
    t.Log("isEmptyReply(nil)", IsEmptyReply(nil))
    t.Log("isEmptyReply(EReply{})", IsEmptyReply(a))
    t.Log("isEmptyReply(EReply{})", IsEmptyReply(b))
}
