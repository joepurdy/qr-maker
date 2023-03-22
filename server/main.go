package main

import (
	"bytes"
	"image"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func main() {
	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/generate-qr-code", func(c *gin.Context) {
		logoFile, _ := c.FormFile("logo")
		qrCodeSize, _ := strconv.Atoi(c.PostForm("size"))
		qrCodeData := c.PostForm("data")

		// Create a temporary file to store the logo
		dst := filepath.Join(os.TempDir(), logoFile.Filename)
		defer os.Remove(dst)

		c.SaveUploadedFile(logoFile, dst)

		logoImage, err := imaging.Open(dst)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid logo file"})
			return
		}

		qrCodeWithLogo, err := generateQRCodeWithLogo(qrCodeData, qrCodeSize, logoImage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "image/png")
		c.Status(http.StatusOK)
		c.Writer.Write(qrCodeWithLogo)
	})

	r.Run()
}

func resizeAndCropLogo(logo image.Image, size int) image.Image {
	logo = imaging.Resize(logo, size, size, imaging.Lanczos)
	logo = imaging.CropCenter(logo, size, size)
	return logo
}

func generateQRCodeWithLogo(data string, size int, logo image.Image) ([]byte, error) {
	qrCode, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	qrCode.DisableBorder = true
	qrCodeImage := qrCode.Image(size)

	logoSize := size / 4
	resizedLogo := resizeAndCropLogo(logo, logoSize)

	offset := (size - logoSize) / 2
	qrCodeWithLogo := imaging.Overlay(qrCodeImage, resizedLogo, image.Pt(offset, offset), 1.0)

	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, qrCodeWithLogo, imaging.PNG)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
