# Blockchain

someone say "learning by doing"

OK! I want to learn blockchain, that doing it!


## basic


blockchain: chain of block, back  linked list, ordered


```
basic block
    previous block hash
    hash
    timestamp
    data
```


```
chain

block#0             <- block#1                 <- block#2

prevHash: none         previousHash: hash0        previousHash: hash1
hash: hash0            hash: hash1                hash: hash2
timestamp: 00          timestamp: 01              timestamp:02
data: "I'm genesis"    data: "I'm block #1"       data: "I'm block #2"
```

hash using `PrevBlockHash + Data + timestamp`

## Hashcash

Bitcoin uses Hashcash, a Proof-of-Work algorithm that was initially developed to prevent email spam. It can be split into the following steps:

1. Take some publicly known data (in case of email, it’s receiver’s email address; in case of Bitcoin, it’s block headers).
2. Add a counter to it. The counter starts at 0.
3. Get a hash of the data + counter combination.
4. Check that the hash meets certain requirements.
    a. If it does, you’re done.
    b. If it doesn’t, increase the counter and repeat the steps 3 and 4.


## CLI


```
./blockchain printchain

./blockchain -data "Send 1 BTC to kmollee"
```