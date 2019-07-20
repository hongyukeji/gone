package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"github.com/wx11055/gone/conv"
)

type aes_mode int

func (m aes_mode) int() int {
	return int(m)
}

const (
	aes_128 aes_mode = 128 / 8
	aes_192          = 192 / 8
	aes_256          = 256 / 8
)

var byte2str = _conv.BytesToStr

//https://zh.wikipedia.org/wiki/%E5%88%86%E7%BB%84%E5%AF%86%E7%A0%81%E5%B7%A5%E4%BD%9C%E6%A8%A1%E5%BC%8F#%E5%B8%B8%E7%94%A8%E6%A8%A1%E5%BC%8F
// 分组工作密码模式
// CBC模式加密  key长度32以内,如果超过32截断,不足32个字节,用0补齐
func AESEnCrypt(key, str []byte) (string, error) {
	k := keyPadding(key, aes_256)
	return AESCBCEnCrypt(k, str)
}

func AESDeCrypt(key []byte, str string) (string, error) {
	k := keyPadding(key, aes_256)
	return AESCBCDeCrypt(k, str)
}

// AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
func AESCBCEnCrypt(key, str []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	str = pkcs5Padding(str, blockSize)
	// str = ZeroPadding(str, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(str))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := str
	blockMode.CryptBlocks(crypted, str)
	return hex.EncodeToString(crypted), nil
}

// CBC模式解密
func AESCBCDeCrypt(key []byte, str string) (string, error) {
	strBytes, err := hex.DecodeString(str)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(strBytes))
	// origData := crypted
	blockMode.CryptBlocks(origData, strBytes)
	origData = pkcs5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return byte2str(origData), nil
}

func aesEnCrypt_v1(key []byte, in []byte) (string, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	out := make([]byte, len(in))
	b.Encrypt(out, []byte(in))
	return _conv.BytesToStr(out), nil
}

func aesDeCrypt_v1(key []byte, in string) (string, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	out := make([]byte, len(in))
	b.Decrypt(out, []byte(in))
	return _conv.BytesToStr(out), nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	//只要少于256就能放到一个byte中，默认的blockSize=16(即采用16*8=128, AES-128长的密钥)
	padding := blockSize - len(ciphertext)%blockSize //需要padding的数目
	//最少填充1个byte，如果原文刚好是blocksize的整数倍，则再填充一个blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding) //生成填充的文本
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// 填充0
func zeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

// 去掉尾部填充的0
func zeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}

func keyPadding(k []byte, l aes_mode) []byte {
	if len(k) > l.int() {
		k = k[:l]
	}
	return zeroPadding(k, l.int())
}
