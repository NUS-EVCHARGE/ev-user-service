package handler

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)
//func publicKeyHandler(w http.ResponseWriter, r *http.Request) {
//	publicKey := privateKey.PublicKey
//	pubASN1, _ := x509.MarshalPKIXPublicKey(&publicKey)
//	pemKey := pem.EncodeToMemory(&pem.Block{
//		Type:  "RSA PUBLIC KEY",
//		Bytes: pubASN1,
//	})
//
//	w.Header().Set("Content-Type", "text/plain")
//	w.Write(pemKey) // Send the public key to the client
//}

func GetPublicKey(c *gin.Context) {
	publicKey, err := loadPublicKey("public.pem")

	if err != nil {
		c.JSON(http.StatusInternalServerError, CreateResponse(http.StatusInternalServerError, "Error", "Failed to load public key"))
		return
	}

	pemkey, err := publicKeyToPEM(publicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CreateResponse(http.StatusInternalServerError, "Error", "Failed to convert public key to PEM"))
		return
	}

	//c.Header("Content-Type", "text/plain")
	//c.String(http.StatusOK, pemkey)
	c.JSON(http.StatusOK, CreateResponse(http.StatusOK, pemkey, "success"))

	return

}

func publicKeyToPEM(publicKey *rsa.PublicKey) (string, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return string(pemKey), nil
}


func loadPublicKey(fileName string) (*rsa.PublicKey, error) {
	pubKeyFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer pubKeyFile.Close()

	pubKeyPEM, err := ioutil.ReadAll(pubKeyFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pubKeyPEM)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("invalid public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}

	return rsaPubKey, nil
}




