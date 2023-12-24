package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

func GenerateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	curve := elliptic.P256() // 使用 P-256 曲线
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

// 使用 ECDSA 签名消息
func ECDSASign(privateKey *ecdsa.PrivateKey, message string) (string, error) {
	hash := sha256.Sum256([]byte(message))
	// r 和 s 是椭圆曲线上的坐标
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", err
	}
	return r.String() + s.String(), nil
}

// 使用 ECDSA 验证签名
func ECDSAVerify(publicKey *ecdsa.PublicKey, message, signature string) bool {
	hash := sha256.Sum256([]byte(message))
	r, s := new(big.Int), new(big.Int)
	r.SetString(signature[:len(signature)/2], 10)
	s.SetString(signature[len(signature)/2:], 10)
	return ecdsa.Verify(publicKey, hash[:], r, s)
}
