package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"

	//"io/ioutil"
	"os"

	"github.com/disintegration/imaging"
)

func encryptImage(imagePath string) (string, error) {
	// This function takes a file path to an image as input, opens the file, and returns an error if there is any issue opening the file.
	// The defer statement ensures that the file is closed when the function returns or if an error occurs.
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Decode the image
	// image.Decode reads the file header to determine the format of the image.
	// Based on the format, image.Decode instantiates the appropriate image type
	// image.Decode will return an instance of *image.YCbCr since we have a  jpg/jpeg image
	//  --------------------------  HOW image.Decode WORKS --------------------
	// the metadata of the image is extracted and used to build a data structure that describes the image's format, size, color space, and other properties.
	// read the pixel data and create a new image object in memory
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// imaging.Rotate() function is called on the img object to rotate the image by 0 degrees.
	rotatedImg := imaging.Rotate(img, 0, color.Transparent)

	// Convert the rotated image to bytes

	// This creates a new buffer with an initial capacity of zero bytes.
	buf := new(bytes.Buffer)
	// rotatedImg as a JPEG image and write the encoded bytes to the buffer.
	// the 3rd arguement is for compression level and since it is nil this means it is set to default that is 75/100 in terms of
	//quality level
	if err := jpeg.Encode(buf, rotatedImg, nil); err != nil {
		return "", err
	}

	// to get the byte slice containing the encoded image data.
	imgBytes := buf.Bytes()

	// Define a fixed encryption key
	key := []byte("secret_key_1234")

	// Encrypt the image bytes using the key

	// This line initializes a new byte slice with the same length as the image bytes slice to store the encrypted bytes.
	encryptedBytes := make([]byte, len(imgBytes))

	//takes the image bytes and encrypts them using a key by adding the value of the key byte to the value of the corresponding image
	//byte, then taking the modulus of the result with 256. This produces a new byte value that is stored in the encrypted byte slice.
	//This process is repeated for each byte in the image data.
	for i := 0; i < len(imgBytes); i++ {
		// determining the current byte of the encryption key
		keyC := int16(key[i%len(key)])

		// to get a byte value
		encryptedC := byte((int16(imgBytes[i]) + keyC) % 256)

		// storing encrypted byte is then stored in the encryptedBytes slice at the corresponding index.
		encryptedBytes[i] = encryptedC
	}

	// Convert the encrypted bytes to base64-encoded text
	// takes a byte slice as input and returns a string
	encryptedText := base64.StdEncoding.EncodeToString(encryptedBytes)

	// 0644 means that the file is readable and writable

	if err := ioutil.WriteFile("encryptedimager.txt", []byte(encryptedText), 0644); err != nil {
		return "", err
	}
	// Return the encrypted text
	return encryptedText, nil
}

func main() {
	imagePath := "rotatedimage.jpg"
	encryptedText, err := encryptImage(imagePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Encrypted text:", encryptedText)
}
