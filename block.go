package main

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/boltdb/bolt"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int // required to verify a proof.
}

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func DeserializeBlock(d []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) error {

	lastHash := bc.tip

	newBlock := NewBlock(data, lastHash)

	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		serialData, err := newBlock.Serialize()
		if err != nil {
			return err
		}
		if err := b.Put(newBlock.Hash, serialData); err != nil {
			return err
		}

		if err := b.Put([]byte("l"), newBlock.Hash); err != nil {
			return err
		}
		bc.tip = newBlock.Hash

		return nil
	})
	return err
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// genesisblock is block#0, here is not previous block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// func NewBlockchain() *Blockchain {
// 	return &Blockchain{[]*Block{NewGenesisBlock()}}
// }

const (
	dbFile       = "blockchain.db"
	blocksBucket = "block"
)

/*
32-byte block-hash -> Block structure (serialized)
'l' -> the hash of the last block in a chain
*/
func NewBlockchain() (*Blockchain, error) {
	var tip []byte

	// Open a DB file.
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {

		// open bucket
		b := tx.Bucket([]byte(blocksBucket))
		// Check if there’s a blockchain stored in it.
		if b == nil {
			// it is empty, create new block
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				return err
			}

			data, err := genesis.Serialize()
			if err != nil {
				return err
			}

			// add genesis block
			if err = b.Put(genesis.Hash, data); err != nil {
				return err
			}
			// update last block for chain
			if err = b.Put([]byte("l"), genesis.Hash); err != nil {
				return err
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	bc := Blockchain{tip, db}

	return &bc, nil
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	_ = i.db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block, err = DeserializeBlock(encodedBlock)
		return err
	})

	i.currentHash = block.PrevBlockHash

	return block
}
