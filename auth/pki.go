package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Komplementariteten/lutra"

	"github.com/Komplementariteten/lutra/util"
)

const privateKeyFile = "./res/key.pem"
const publicKeyFile = "./res/pub.pem"

func loadEcdsaKey(config *lutra.LutraConfig) *ecdsa.PrivateKey {
	f, err := os.OpenFile(privateKeyFile, os.O_RDONLY, 0700)
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	pem, _ := pem.Decode(bytes)
	decryptedBlock, err := x509.DecryptPEMBlock(pem, []byte(config.PrivateKeyPw))
	if err != nil {
		panic(err)
	}
	if len(decryptedBlock) == 0 {
		panic(fmt.Errorf("No Data received from Encrypted PEM Block"))
	}
	key, err := x509.ParseECPrivateKey(decryptedBlock)
	if err != nil {
		panic(err)
	}
	return key
}

func createNewEcdsaKey(config *lutra.LutraConfig) *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		panic(err)
	}
	err = util.CreateResourceFolderIfNotExists()
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(privateKeyFile, os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	marshal, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	pemBlock, err := x509.EncryptPEMBlock(rand.Reader, "EC PRIVATE KEY", marshal, []byte(config.PrivateKeyPw), x509.PEMCipherAES256)
	if err != nil {
		panic(err)
	}
	err = pem.Encode(f, pemBlock)
	if err != nil {
		panic(err)
	}

	return privateKey
}

func LoadEscdaKeyAsBin(config *lutra.LutraConfig) []byte {
	f, err := os.OpenFile(privateKeyFile, os.O_RDONLY, 0700)
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	pem, _ := pem.Decode(bytes)
	decryptedBlock, err := x509.DecryptPEMBlock(pem, []byte(config.PrivateKeyPw))
	if err != nil {
		panic(err)
	}
	if len(decryptedBlock) == 0 {
		panic(fmt.Errorf("No Data received from Encrypted PEM Block"))
	}
	return decryptedBlock
}

func GetEcdsaKey(config *lutra.LutraConfig) *ecdsa.PrivateKey {
	if ok, err := util.IfFileExists(privateKeyFile); !ok {
		if err != nil {
			panic(err)
		}
		k := createNewEcdsaKey(config)
		return k
	} else {
		k := loadEcdsaKey(config)
		return k
	}

}
