package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
)

const reward = 50

// 1.定义交易结构
type Transaction struct {
	TXID     []byte     //交易ID
	TXINPUT  []TXInput  //交易输入
	TXOUTPUT []TXOutput //交易输出
}

// 定义交易输入
type TXInput struct {
	TXid  []byte //交易ID
	Index int64  //索引
	//Sig   string //解锁脚本，先用地址来模拟
	Signature []byte //真正的数字签名，由r,s组成的[]byte
	PubKey    []byte //这里的PubKey不存储原始公钥，而是存储x和y的拼接字符串，在校验段重新拆分（参考r,s传递），是公钥不是哈希，也不是地址
}

// 定义交易输出
type TXOutput struct {
	Value float64 //转账金额
	//PutKeyHash string  //锁定脚本，用地址模拟
	PubKeyHash []byte //收款方的公钥的哈希，注意是哈希不是公钥，也不是地址
}

// 由于存储的字段是地址的公钥哈希，所以无法直接创建TXOutput
// 为了能够得到公钥哈希，需要一个Lock函数来处理
func (output *TXOutput) Lock(address string) {
	//锁定动作
	output.PubKeyHash = GetPubKeyFromAddress(address)
}

// 给TXOutput提供一个创建的方法，否则无法调用Lock
func NewTXOutput(value float64, address string) *TXOutput {
	output := TXOutput{
		Value: value,
	}
	output.Lock(address)
	return &output
}

// 设置交易ID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

// 实现一个函数判断当前函数是否为挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	if len(tx.TXINPUT) == 1 {
		//1.交易input只有一个
		//2.交易input的id为空
		//3.交易的index为-1
		if len(tx.TXINPUT) == 1 && len(tx.TXINPUT[0].TXid) == 0 && tx.TXINPUT[0].Index == -1 {
			return true
		}
	}
	return false
}

// 2.提供创建交易方法(挖矿交易)
func NewCoinbaseTX(address string, data string) *Transaction {
	//挖矿交易只有一个input，无需引用交易id，无需引用index
	//矿工由于挖矿时无需指定签名，所以PubKey字段可以由矿工自由填写，一般填写矿池的名字
	//签名先填写为空，后面创建完整交易后，最后做一次签名即可
	input := TXInput{
		[]byte{},
		-1,
		nil,
		[]byte(data),
	}

	//新的创建方法
	output := NewTXOutput(reward, address)

	tx := &Transaction{
		[]byte{},
		[]TXInput{input},
		[]TXOutput{*output},
	}
	tx.SetHash()

	return tx
}

// 创建普通转账交易
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	//创建交易之后钥进行数字签名->所以需要私钥->打开钱包（newwallets()）
	ws := NewWallets()

	//找到自己的钱包，根据地址找到自己的wallet
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		fmt.Println("没有找到该地址的钱包，交易创建失败！")
		return nil
	}
	//得到对应的私钥，公钥
	pubKey := wallet.PubKey
	privateKey := wallet.PrivateKey

	// a.找到最合理的utxo集合，map[string][]uint64
	pubKeyHash := HashPubKey(pubKey)
	utxos, resValue := bc.FindNeedUTXOs(pubKeyHash, amount)

	if resValue < amount {
		fmt.Println("余额不足，交易失败！")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	// b.创建交易输入，将这些utxo逐一转成inputs
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TXInput{
				[]byte(id),
				int64(i),
				nil,
				pubKey,
			}
			inputs = append(inputs, input)
		}
	}

	// c.创建outputs
	output := NewTXOutput(amount, to)
	outputs = append(outputs, *output)

	// d.如果有零钱，要找零
	if resValue > amount {
		//找零
		output := NewTXOutput(resValue-amount, to)
		outputs = append(outputs, *output)
	}
	tx := Transaction{
		[]byte{},
		inputs,
		outputs,
	}
	tx.SetHash()

	bc.SignTransaction(&tx, privateKey)

	return &tx
}

// 签名的具体实现,参数为：私钥，inputs里面所有引用的交易的结构map[string]Transaction
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}
	//创建一个当前交易的副本，txCopy，使用函数TrimmedCopy：要把Signature和PubKey字段设置为nil
	txCopy := tx.TrimmedCopy()

	//遍历循环txCopy的inputs，得到这个input索引的output的公钥哈希
	for i, intput := range txCopy.TXINPUT {
		prevTX := prevTXs[string(intput.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("引用无效的交易")
		}
		//不要对intput进行赋值，这是一个副本，要对txCopy.TXINPUT[xx]进行操作,否则无法吧pubKeyHash传进来
		txCopy.TXINPUT[i].PubKey = prevTX.TXOUTPUT[intput.Index].PubKeyHash

		//所需的三个数据都具备了，开始做哈希处理
		//生成要签名的数据，要签名的数据一定是哈希值
		//我们对每一个input都要签名一次，签名的数据是由当前input引用的output的哈希+当前的outputs（都承载在当前这个txCopy里面）
		//要对这个拼好的txCopy进行哈希处理，SetHash得到TXID，这个TXID就是我们要签名的最终数据
		txCopy.SetHash()
		//还原，以免影响后面input的签名
		txCopy.TXINPUT[i].PubKey = nil
		signDataHash := txCopy.TXID
		//执行签名动作得到r,s字节流
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, signDataHash)
		if err != nil {
			log.Panic(err)
		}
		//放到我们所签名的input的Signture中
		signature := append(r.Bytes(), s.Bytes()...)
		tx.TXINPUT[i].Signature = signature

	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var intputs []TXInput
	var outputs []TXOutput

	for _, input := range tx.TXINPUT {
		intputs = append(intputs, TXInput{input.TXid, input.Index, nil, nil})
	}

	for _, output := range tx.TXOUTPUT {
		outputs = append(outputs, output)
	}
	return Transaction{tx.TXID, intputs, outputs}
}

// 校验
// 所需要的数据：公钥，数据（txCopy，生成哈希），签名
// 要对每一个签名过的过的input进行校验
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	//得到签名的数据
	txCopy := tx.TrimmedCopy()
	for i, input := range tx.TXINPUT {
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("引用的交易无效")
		}
		txCopy.TXINPUT[i].PubKey = prevTX.TXOUTPUT[input.Index].PubKeyHash
		txCopy.SetHash()
		dataHash := txCopy.TXID
		//得到Signature，反推回r,s
		signature := input.Signature
		//拆解PubKey,X,Y得到原生公钥
		pubKey := input.PubKey

		//定义两个辅助的big.Int{}
		r := big.Int{}
		s := big.Int{}

		//拆分signature,平均分，前半部分给r，后半部分给s
		r.SetBytes(signature[0 : len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])

		//定义两个辅助的big.Int{}
		X := big.Int{}
		Y := big.Int{}

		//拆分signature,平均分，前半部分给X，后半部分给Y
		X.SetBytes(pubKey[0 : len(signature)/2])
		Y.SetBytes(pubKey[len(signature)/2:])

		pubKeyOrigin := ecdsa.PublicKey{elliptic.P256(), &X, &Y}

		//Verify
		if !ecdsa.Verify(&pubKeyOrigin, dataHash, &r, &s) {
			return false
		}
	}
	return true
}
