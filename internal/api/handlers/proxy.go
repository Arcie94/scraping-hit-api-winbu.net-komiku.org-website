package handlers

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"

	"komiku-scraper/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
)

// ImageProxyHandler handles image proxy requests
func ImageProxyHandler(c *fiber.Ctx) error {
	// Get URL from query parameter (base64 encoded for safety)
	encodedURL := c.Query("url")
	if encodedURL == "" {
		return c.JSON(models.ErrorResponse("INVALID_URL", "URL parameter is required"))
	}

	// Decode URL
	urlBytes, err := base64.URLEncoding.DecodeString(encodedURL)
	if err != nil {
		return c.JSON(models.ErrorResponse("INVALID_URL", "Invalid base64 URL"))
	}
	imageURL := string(urlBytes)

	// Get size parameter (optional)
	size := c.Query("size", "") // small, medium, large

	// Fetch image
	resp, err := http.Get(imageURL)
	if err != nil {
		return c.JSON(models.ErrorResponse("FETCH_FAILED", "Failed to fetch image"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return c.JSON(models.ErrorResponse("FETCH_FAILED", "Image not found"))
	}

	// Read image data
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(models.ErrorResponse("READ_FAILED", "Failed to read image"))
	}

	// If no resize requested, return original
	if size == "" {
		contentType := resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "image/jpeg" // Default
		}

		c.Set("Content-Type", contentType)
		c.Set("Cache-Control", "public, max-age=86400") // 24 hours
		return c.Send(imgData)
	}

	// Resize image
	resizedData, err := resizeImage(imgData, size)
	if err != nil {
		// If resize fails, return original
		c.Set("Content-Type", "image/jpeg")
		c.Set("Cache-Control", "public, max-age=86400")
		return c.Send(imgData)
	}

	c.Set("Content-Type", "image/jpeg")
	c.Set("Cache-Control", "public, max-age=86400")
	return c.Send(resizedData)
}

// resizeImage resizes an image based on size parameter
func resizeImage(data []byte, size string) ([]byte, error) {
	// Decode image
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// Determine target width based on size
	var width uint
	switch size {
	case "small":
		width = 150
	case "medium":
		width = 300
	case "large":
		width = 600
	default:
		return data, nil // Unknown size, return original
	}

	// Resize maintaining aspect ratio
	resized := resize.Resize(width, 0, img, resize.Lanczos3)

	// Encode to JPEG
	buf := new(bytes.Buffer)

	if format == "png" {
		err = png.Encode(buf, resized)
	} else {
		err = jpeg.Encode(buf, resized, &jpeg.Options{Quality: 85})
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
