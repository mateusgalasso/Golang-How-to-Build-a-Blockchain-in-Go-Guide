package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"goblockchain/utils"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

func NewWallet() *Wallet {
	wallet := new(Wallet)
	//1 -Criate private and publickey
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	wallet.privateKey = privateKey
	wallet.publicKey = &privateKey.PublicKey
	//2 - perform sha-256 hashing on the public key (32 bytes)
	h2 := sha256.New()
	h2.Write(wallet.publicKey.X.Bytes())
	h2.Write(wallet.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	//3 - Perform RIPEMD-160 hashing on result of SHA-256 (20 bytes)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	//4 - Add verison byte in front RIPEMD-160 hash (0x00 for Mainnet)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])
	//5 -  Perform SHA-256 on the extended result
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	//6 - Perform SHA-256 on the previous resulto of SHA-256
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	//7 - Take 4 first bytes of the second SHA-256 hash of checksum
	chsum := digest6[:4]
	//8 - Add the 4 checksum bytes from 7 at the end of extended RIPEMD-160 hash from 4
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum)
	//9 - convert the result from a byte string into base58
	address := base58.Encode(dc8)

	wallet.blockchainAddress = address
	return wallet
}
func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}
func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}
func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}
func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}
func (w *Wallet) BlockChainAddress() string {
	return w.blockchainAddress
}

type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBlocchainAddress     string
	recipientBlockchainAddress string
	value                      float32
}

func (w *Wallet) NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string, recipient string, value float32) *Transaction {
	return &Transaction{
		senderPrivateKey:           privateKey,
		senderPublicKey:            publicKey,
		senderBlocchainAddress:     sender,
		recipientBlockchainAddress: recipient,
		value:                      value,
	}
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}

func (t Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlocchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}
