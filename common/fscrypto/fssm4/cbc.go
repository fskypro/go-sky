/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: cbc 模式
@author: fanky
@version: 1.0
@date: 2026-03-20
**/

package fssm4

import "errors"

/*
SM4 每次只能处理 固定 16 字节，但音频数据是任意长度的，所以需要解决两个问题：

	数据不是 16 的倍数怎么办？ → PKCS7 填充
	每块独立加密安全吗？ → CBC 模式

数据不够 16 字节的整数倍，就在末尾补字节，补几个就填数字几：

	原始数据: [AA BB CC DD EE]  ← 5字节，差11字节
	填充后:   [AA BB CC DD EE 0B 0B 0B 0B 0B 0B 0B 0B 0B 0B 0B]
	                              ↑ 补了11个 0x0B（11的十六进制）
	解密后看最后一个字节的值，就知道要去掉几个填充字节。

CBC 模式

	如果每块独立加密（ECB模式），相同的数据块会产生相同的密文：
		ECB 模式（不安全）:
		明文块A → 加密 → 密文X
		明文块A → 加密 → 密文X   ← 一样！攻击者能发现规律
		明文块B → 加密 → 密文Y

CBC 解决方法：

	每块加密前先和前一块密文 XOR：
	CBC 模式:
		IV(随机)  ──XOR──  明文块1  →  加密  →  密文块1
		密文块1   ──XOR──  明文块2  →  加密  →  密文块2
		密文块2   ──XOR──  明文块3  →  加密  →  密文块3
*/

// 将不够 16 的整数倍数据块，凑够整数倍数据块
func pkcs7Pad(data []byte) []byte {
	pad := 16 - len(data)%16
	padded := make([]byte, len(data)+pad)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = byte(pad)
	}
	return padded
}

// 将凑够 16 整数倍的数据块去除
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 || len(data)%16 != 0 {
		return nil, errors.New("invalid padded data length")
	}
	pad := int(data[len(data)-1])
	if pad == 0 || pad > 16 {
		return nil, errors.New("invalid PKCS7 padding value")
	}
	for i := len(data) - pad; i < len(data); i++ {
		if data[i] != byte(pad) {
			return nil, errors.New("PKCS7 padding mismatch")
		}
	}
	return data[:len(data)-pad], nil
}

// -------------------------------------------------------------------
// 加密解密字节块
// -------------------------------------------------------------------
func CBCEncrypt(sm4 *SM4, iv, plaintext []byte) []byte {
	if len(plaintext) == 0 {
		return []byte{}
	}
	padded := pkcs7Pad(plaintext)
	ciphertext := make([]byte, len(padded))
	prev := iv
	buf := make([]byte, 16, 16)
	for i := 0; i < len(padded); i += 16 {
		for j := range 16 {
			buf[j] = padded[i+j] ^ prev[j]
		}
		sm4.EncryptBlock(ciphertext[i:], buf)
		prev = ciphertext[i : i+16]
	}
	return ciphertext
}

func CBCDecrypt(sm4 *SM4, iv, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return []byte{}, nil
	}
	if len(ciphertext)%16 != 0 {
		return nil, errors.New("ciphertext length is not a multiple of 16")
	}
	plaintext := make([]byte, len(ciphertext))
	prev := iv
	buf := make([]byte, 16, 16)
	for i := 0; i < len(ciphertext); i += 16 {
		sm4.DecryptBlock(buf, ciphertext[i:i+16])
		for j := range 16 {
			plaintext[i+j] = buf[j] ^ prev[j]
		}
		prev = ciphertext[i : i+16]
	}
	return pkcs7Unpad(plaintext)
}
