package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"github.com/btcsuite/btcutil/base58"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

const walletFile = "wallet.dat"

// 定义一个Wallets结构，来保存所有wallet以及他的地址
type Wallets struct {
	//map[地址]钱包
	WalletsMap map[string]*Wallet
}

// SerializableWallet Solve gob: type elliptic.p256Curve has no exported fields
// 因为之前gob.Register(elliptic.P256())的用法在现在的go版本无法使用了
type SerializableWallet struct {
	D         *big.Int
	X, Y      *big.Int
	PublicKey []byte
}

// 创建方法
func NewWallets() *Wallets {
	var ws Wallets
	ws.WalletsMap = make(map[string]*Wallet)

	ws.LoadFile()

	return &ws
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.NewAddress()

	ws.WalletsMap[address] = wallet

	ws.saveToFile()
	return address
}

// 保存方法，把新建的wallet添加进去
func (ws *Wallets) saveToFile() {
	var buffer bytes.Buffer

	//fmt.Println(walletFile)

	gob.Register(SerializableWallet{})

	wallets := make(map[string]SerializableWallet)
	for k, v := range ws.WalletsMap {
		wallets[k] = SerializableWallet{
			D:         v.PrivateKey.D,
			X:         v.PrivateKey.PublicKey.X,
			Y:         v.PrivateKey.PublicKey.Y,
			PublicKey: v.PubKey,
		}
	}

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(wallets)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, buffer.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

// 读取文件方法，把所有的wallet读出来
func (ws *Wallets) LoadFile() {
	//fmt.Println(walletFile)
	_, err := os.Stat(walletFile)
	if os.IsNotExist(err) {
		ws.WalletsMap = make(map[string]*Wallet)
		return
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets map[string]SerializableWallet
	//gob.Register(elliptic.P256())
	gob.Register(SerializableWallet{})
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.WalletsMap = make(map[string]*Wallet)
	//ws.Wallets = wallets.Wallets
	for k, v := range wallets {
		ws.WalletsMap[k] = &Wallet{
			PrivateKey: &ecdsa.PrivateKey{
				PublicKey: ecdsa.PublicKey{
					Curve: elliptic.P256(),
					X:     v.X,
					Y:     v.Y,
				},
				D: v.D,
			},
			PubKey: v.PublicKey,
		}
	}
}

func (ws *Wallets) ListAllAddress() []string {
	var addresses []string
	//遍历钱包，将所有的key取出来返回
	for address := range ws.WalletsMap {
		addresses = append(addresses, address)
	}

	return addresses
}

// 通过地址返回公钥哈希
func GetPubKeyFromAddress(address string) []byte {
	//解码
	addressByte := base58.Decode(address) //25字节

	//截取出公钥哈希，去除version和校验码(25-1-4)
	len := len(addressByte)
	pubKeyHash := addressByte[1 : len-4]

	return pubKeyHash
}
