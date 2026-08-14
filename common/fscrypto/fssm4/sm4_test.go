package fssm4

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"testing"

	"fsky.pro/fstest"
)

func TestSM4(t *testing.T) {
	fstest.PrintTestBegin("SM4")
	defer fstest.PrintTestEnd()

	// 转换成 hex 后的密钥
	keyHex := "0123456789abcdeffedcba9876543210"
	plainHex := "0123456789abcdeffedcba9876543210"

	// GB/T 32907 标准结果是下面这个
	expectHex := "681edf34d206965e86b3e94f536e4246"

	// 还原密钥
	key, _ := hex.DecodeString(keyHex)
	plain, _ := hex.DecodeString(plainHex)

	sm4c, err := NewSM4(key)
	if err != nil {
		fmt.Println("初始化失败:", err)
		return
	}

	cipher := make([]byte, 16)
	sm4c.EncryptBlock(cipher, plain)
	actualHex := hex.EncodeToString(cipher)

	fmt.Printf("Key:  %s\n", keyHex)
	fmt.Printf("明文: %s\n", plainHex)
	fmt.Printf("期望: %s\n", expectHex)
	fmt.Printf("实测: %s\n", actualHex)

	if actualHex == expectHex {
		fmt.Println("结果: ✅ 与 GB/T 32907 标准一致")
	} else {
		fmt.Println("结果: ❌ 不一致！")
	}
}

func TestEncryptKey(t *testing.T) {
	fstest.PrintTestBegin("EncryptKey")
	defer fstest.PrintTestEnd()

	// 用于加密数据的 sm4 key
	sm4Key := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, sm4Key)
	if err != nil {
		fmt.Println("create sm4 key fail,", err)
		return
	}
	hexSM4Key := hex.EncodeToString(sm4Key)
	fmt.Println("用于加密数据的 sm4-key =", hexSM4Key)

	// 用于加密 sm4 key 的 key
	keyOfSM4Key := []byte("*YXHT->20260322*")

	sm4, err := NewSM4(keyOfSM4Key)
	if err != nil {
		fmt.Println("create sm4 fail,", err)
		return
	}
	fmt.Println("用于加密 key 的 key =", hex.EncodeToString(keyOfSM4Key))

	// 加密 sm4 key
	encKeyOfKey := make([]byte, 16)
	sm4.EncryptBlock(encKeyOfKey, sm4Key)
	fmt.Println("加密后的 sm4 key =", hex.EncodeToString(encKeyOfKey))

	// ---------------------------------
	// 还原 sm4 key
	sm4, _ = NewSM4(keyOfSM4Key)
	sm4.DecryptBlock(sm4Key, encKeyOfKey)
	fmt.Println("还原的 sm4 key =", hex.EncodeToString(sm4Key))
}

func TestEncryptBytes(t *testing.T) {
	fstest.PrintTestBegin("EncryptBytes")
	defer fstest.PrintTestEnd()

	origin := "abc/123/efgh/456.yxwav"
	fmt.Println("解密前：", origin)

	sm4, _ := NewSM4([]byte("*YXHT->20260322*"))
	data := sm4.EncryptBytes([]byte(origin))
	fmt.Println("加密后：", data)

	borigin, _ := sm4.DecryptBytes(data)
	fmt.Println("解密后：", string(borigin))
}
