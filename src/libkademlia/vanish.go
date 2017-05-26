package libkademlia

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	mathrand "math/rand"
	"sss"
	"time"
)

type VanashingDataObject struct {
	AccessKey  int64
	Ciphertext []byte
	NumberKeys byte
	Threshold  byte
}

func GenerateRandomCryptoKey() (ret []byte) {
	for i := 0; i < 32; i++ {
		ret = append(ret, uint8(mathrand.Intn(256)))
	}
	return
}

func GenerateRandomAccessKey() (accessKey int64) {
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	accessKey = r.Int63()
	return
}

func CalculateSharedKeyLocations(accessKey int64, count int64) (ids []ID) {
	r := mathrand.New(mathrand.NewSource(accessKey))
	ids = make([]ID, count)
	for i := int64(0); i < count; i++ {
		for j := 0; j < IDBytes; j++ {
			ids[i][j] = uint8(r.Intn(256))
		}
	}
	return
}

func encrypt(key []byte, text []byte) (ciphertext []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ciphertext = make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], text)
	return
}

func decrypt(key []byte, ciphertext []byte) (text []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext is not long enough")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}

func (k *Kademlia) VanishData(data []byte, numberKeys byte, threshold byte, timeoutSeconds int) (V VanashingDataObject) {
	key := GenerateRandomCryptoKey()
	V.Ciphertext = encrypt(key, data)
	V.AccessKey = GenerateRandomAccessKey()
	V.NumberKeys = numberKeys
	V.Threshold = threshold
	skey, err := sss.Split(numberKeys, threshold, key)
	if err != nil {
		V.NumberKeys = 0 // NumberKeys = 0 means error
		return V
	}
	addrs := CalculateSharedKeyLocations(V.AccessKey, int64(numberKeys))
	i := 0
	for kid, kv := range skey {
		packed := append([]byte{kid}, kv...)
		_, err := k.DoIterativeStore(addrs[i], packed)
		i++
		if err != nil {
			V.NumberKeys = 0 // NumberKeys = 0 means error
			return V
		}
	}
	return V
}

func (k *Kademlia) UnvanishData(vdo VanashingDataObject) (data []byte) {
	keys := make(map[byte][]byte)
	addrs := CalculateSharedKeyLocations(vdo.AccessKey, int64(vdo.NumberKeys))
	for i := 0; i < len(addrs); i++ {
		packed, err := k.DoIterativeFindValue(addrs[i])
		if err != nil {
			kid := packed[0]
			kv := packed[1:]
			keys[kid] = kv
		}
	}
	if len(keys) < int(vdo.Threshold) {
		return nil
	}
	key := sss.Combine(keys)
	data = decrypt(key, vdo.Ciphertext)
	return data
}
