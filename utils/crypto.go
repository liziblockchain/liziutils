package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

func Base64Encode(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

func Base64Decode(in string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil
	}

	return decoded
}

func HexEncode(in []byte) string {
	return hex.EncodeToString(in)
}

func HexDecode(in string) []byte {
	decoded, err := hex.DecodeString(in)
	if err != nil {
		return nil
	}

	return decoded
}

func Sha256(indata []byte) string {
	encode := sha256.Sum256((indata))
	return hex.EncodeToString(encode[:])
}

// AesEncrypt
// [in] src : original data
// [out] encryped data with base64 encode
func AesEncrypt(src, key, iv string) (string, error) {
	if (len(src) < 1) || (len(key) < 1) || (len(iv) < 1) {
		return "", errors.New("invalid param")
	}

	encrypt, err := aesEncrypt([]byte(src), []byte(key), []byte(iv))
	if err != nil {
		return "", err
	}
	return Base64Encode(encrypt), nil
}

func aesEncrypt(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if (err != nil) || (block == nil) {
		return nil, errors.New("invalid decrypt key")
	}
	blockSize := block.BlockSize()
	plaintextNew := PKCS5Padding(plaintext, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintextNew))
	blockMode.CryptBlocks(ciphertext, plaintextNew)

	return ciphertext, nil
}

// AesDecrypt
// [in] src : encryped data with base64 encode
// [out] original data
func AesDecrypt(src, key, iv string) (string, error) {
	if (len(src) < 1) || (len(key) < 1) || (len(iv) < 1) {
		return "", errors.New("invalid param")
	}

	srcDecode := Base64Decode(src)
	orig, err := aesDecrypt(srcDecode, []byte(key), []byte(iv))
	if err != nil {
		return "", err
	}
	return string(orig), nil
}

func aesDecrypt(src []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if (err != nil) || (block == nil) {
		return nil, errors.New("invalid decrypt key")
	}

	blockSize := block.BlockSize()
	if len(src) < blockSize {
		return nil, errors.New("ciphertext too short")
	}

	if len(src)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	blockModel := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(src))
	blockModel.CryptBlocks(plaintext, src)
	plaintext = PKCS5UnPadding(plaintext)

	return plaintext, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	// if padding < 1 {
	// 	return src
	// }
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	if length < 1 {
		return src
	}
	unpadding := int(src[length-1])
	lenNew := length - unpadding
	if (lenNew < 0) || (lenNew > length-1) {
		return src
	}
	return src[:lenNew]
}

//convert base64 encoded string to binary as input
/*base64 encoded string*/
func RsaDecrypt(src string, privateKey []byte) (string, error) {

	block, _ := pem.Decode(privateKey)
	if block == nil {
		// fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAA no pem data")
		return "", errors.New("no pem data")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// base64Origin := Base64Decode(src)
	// fmt.Println("bbbbbbbbbbbbbbbb, base64Origin:", src, string(base64Origin))
	// decoded, err := rsa.DecryptPKCS1v15(rand.Reader, priv, base64Origin)
	decoded, err := rsa.DecryptPKCS1v15(rand.Reader, priv, Base64Decode(src))
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
