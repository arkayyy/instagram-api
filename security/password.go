package pass

import (
	//PRE-DEFINED PACKAGE IMPORTS:-
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	// "strings"

	//CUSTOM PACKAGE IMPORTS:-
	"main/config"
)



func Encrypt(pass string) string{
    //fmt.Println("Encryption Program v0.01")

		text := []byte(pass)
		key := []byte(config.Pkey)

		// generate a new aes cipher using our 32 byte long key
		c, err := aes.NewCipher(key)
		
		if err != nil {
			fmt.Println(err)
		}

		gcm, err := cipher.NewGCM(c)

		if err != nil {
			fmt.Println(err)
		}

		// creates a new byte array the size of the nonce
		// which must be passed to Seal
		nonce := make([]byte, gcm.NonceSize())
		// populates our nonce with a cryptographically secure
		// random sequence
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			fmt.Println(err)
		}

		// here we encrypt our text using the Seal function
		
		// var str string = (string(gcm.Seal(nonce, nonce, text, nil)))
		// var b []byte = []byte(str)
		// fmt.Println("original: ", gcm.Seal(nonce, nonce, text, nil))
		// fmt.Println("converted: ",b)
		return string(gcm.Seal(nonce, nonce, text, nil))
}





func Decrypt(txt []byte) string{


		fmt.Println("Decryption Program v0.01")

		key := []byte(config.Pkey)
		ciphertext :=[]byte{246,132,85,118,135,40,197,30,13,75,10,4,219,57,128,150,22,160,86,83,174,210,239,151,197,221,114,217,131,227,186,61,203}
		//taking dummy value as of now

		c, err := aes.NewCipher(key)
		if err != nil {
			fmt.Println(err)
		}

		gcm, err := cipher.NewGCM(c)
		if err != nil {
			fmt.Println(err)
		}

		nonceSize := gcm.NonceSize()
		if len(ciphertext) < nonceSize {
			fmt.Println(err)
		}

		nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(plaintext))
		return string(plaintext)
}