package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

/*
In Bitcoin, “target bits” is the block header storing the difficulty at which the block was mined.
We won’t implement a target adjusting algorithm,
for now, so we can just define the difficulty as a global constant.
*/
const (
	targetBits = 15 // 0~255, diffculty of calculate nonce
	maxNonce   = math.MaxInt64
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			// nonce here is the counter from the Hashcash description above
			/*
				Nonce是一個在加密通訊只能使用一次的數字。 在認證協定中，它往往是一個隨機或偽隨機數，以避免重送攻擊。
				Nonce也用於串流加密法以確保安全。 如果需要使用相同的金鑰加密一個以上的訊息，
				就需要Nonce來確保不同的訊息與該金鑰加密的金鑰流不同。
			*/
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0 // counter

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	/*
		1. Prepare data.
		2. Hash it with SHA-256.
		3. Convert the hash to a big integer.
		4. Compare the integer with the target.

	*/

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
