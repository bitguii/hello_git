package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

//这里钱包是一个结构，每一个钱包里保存了公钥，私钥对

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	//这里的PubKey不存储原始公钥，而是存储x和y的拼接字符串，在校验段重新拆分（参考r,s传递）
	PubKey []byte
}

// 创建钱包
func NewWallet() *Wallet {
	//创建曲线
	curve := elliptic.P256()

	//生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	//生成公钥
	pubKeyOrig := privateKey.PublicKey

	//拼接x，y
	pubKey := append(pubKeyOrig.X.Bytes(), pubKeyOrig.Y.Bytes()...)

	return &Wallet{
		privateKey,
		pubKey,
	}
}

// 生成地址
func (w *Wallet) NewAddress() string {
	pubKey := w.PubKey

	rip160HashValue := HashPubKey(pubKey)
	version := byte(00)
	payload := append([]byte{version}, rip160HashValue...)

	//checkSum
	checkCode := CheckSum(payload)

	payload = append(payload, checkCode...)

	//go语言有一个库，是go语言实现btc全节点源码
	address := base58.Encode(payload)

	return address
}

func HashPubKey(data []byte) []byte {
	hash := sha256.Sum256(data)

	//相当于生成一个编码器
	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}

	//返回rip160哈希结果
	rip160HashValue := rip160hasher.Sum(nil)

	return rip160HashValue
}

func CheckSum(data []byte) []byte {
	//checkSum
	//两次sha256
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])

	//前4字节校验码
	checkCode := hash2[:4]

	return checkCode
}

func IsValidAddress(address string) bool {
	//解码
	addressByte := base58.Decode(address)
	if len(addressByte) < 4 {
		return false
	}

	//取数据
	payload := addressByte[:len(addressByte)-4]
	CheckSum1 := addressByte[len(addressByte)-4:]

	//做checkSum函数
	CheckSum2 := CheckSum(payload)

	//比较
	return bytes.Equal(CheckSum1, CheckSum2)
}
