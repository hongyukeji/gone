package _crypto_test

import (
	"crypto/aes"
	"gone/crypto"
	"testing"
)

type CryptEntity struct {
	key string
	val string
}

var encryptTests = []CryptEntity{
	{"1234567890123456", "qwerasdfzxcvtyui"},
	{"asdf", "qwerasdfzxcvtyui"},
	{"asdf", "asdf"},
	{"aa", "zxcvxv"},
	{"1234", "asdf"},
	{"", "asdfaa"},
}

// 测试原始加密解密方式
func TestCipherEncryptAndDecrypt(t *testing.T) {
	for i, tt := range encryptTests {
		block, err := aes.NewCipher([]byte(tt.key))
		if err != nil {
			t.Errorf("NewCipher(%d bytes) = %s", len(tt.key), err)
			continue
		}
		out := make([]byte, len(tt.val))
		block.Encrypt(out, []byte(tt.val))
		b2, err := aes.NewCipher([]byte(tt.key))
		if err != nil {
			t.Errorf("NewCipher(%d bytes) = %s", len(tt.key), err)
			continue
		}
		o := make([]byte, len(tt.val))
		b2.Decrypt(o, out)
		for j, v := range o {
			if v != tt.val[j] {
				t.Errorf("Cipher.Encrypt %d: out[%d] = %#x, want %#x", i, j, v, tt.val[j])
				break
			}
		}
	}
}

// 测试 v1版本加密解密
func TestCipherEncryptAndDecrypt_V1(t *testing.T) {
	for i, tt := range encryptTests {
		s, err := _crypto.AESEnCrypt([]byte(tt.key), []byte(tt.val))
		if err != nil {
			t.Errorf("NewCipher(%d bytes) = %s", len(tt.key), err)
			continue
		}
		o, err := _crypto.AESDeCrypt([]byte(tt.key), s)
		if err != nil {
			t.Errorf("NewCipher(%d bytes) = %s", len(tt.key), err)
			continue
		}
		if o != tt.val {
			t.Errorf("Cipher.Encrypt %d ==> %s, want %s", i, o, tt.val)
		}
	}
}
