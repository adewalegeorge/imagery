package main

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/h2non/bimg"
)

func main() {
	server := gin.Default()
	server.SetTrustedProxies(nil)

	server.GET("/api/opt", func(reqRes *gin.Context) {
		rel := reqRes.Query("rel") // relative path in bucket
		abs := reqRes.Query("abs") // absolute URL
		var imageURL string
		if rel != "" {
			bucketHost := os.Getenv("BUCKET_HOST")
			imageURL = strings.TrimRight(bucketHost, "/") + "/" + strings.TrimLeft(rel, "/")
		} else if abs != "" {
			imageURL = abs
		} else {
			reqRes.JSON(http.StatusBadRequest, gin.H{"error": "rel or abs parameter required"})
			return
		}

		widthStr := reqRes.Query("w")              // width:Integer
		heightStr := reqRes.Query("h")             // height:Integer (optional)
		crop := reqRes.DefaultQuery("c", "false")  // crop:Boolean (true, false)
		format := reqRes.DefaultQuery("f", "auto") // format:String (jpeg, png, webp, avif)
		blurStr := reqRes.DefaultQuery("b", "0")   // blur:Float (optional)
		gray := reqRes.DefaultQuery("g", "false")  // gray:Boolean (optional)

		if imageURL == "" || widthStr == "" {
			reqRes.JSON(http.StatusBadRequest, gin.H{"error": "url and w (width) are required"})
			return
		}

		width, err := strconv.Atoi(widthStr)
		if err != nil || width <= 0 {
			reqRes.JSON(http.StatusBadRequest, gin.H{"error": "invalid width"})
			return
		}

		height := 0
		if heightStr != "" {
			h, err := strconv.Atoi(heightStr)
			if err != nil || h < 0 {
				reqRes.JSON(http.StatusBadRequest, gin.H{"error": "invalid height"})
				return
			}
			height = h
		}

		blur := 0.0
		if blurStr != "" {
			b, err := strconv.ParseFloat(blurStr, 64)
			if err != nil || b < 0 {
				reqRes.JSON(http.StatusBadRequest, gin.H{"error": "invalid blur value"})
				return
			}
			blur = b
		}

		resp, err := http.Get(imageURL)
		if err != nil || resp.StatusCode != 200 {
			reqRes.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch image"})
			return
		}
		defer resp.Body.Close()
		imgBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			reqRes.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read image"})
			return
		}

		options := bimg.Options{
			Width:        width,
			Height:       height, // 0 means proportional scaling
			Crop:         crop == "true",
			Quality:      85,
			GaussianBlur: bimg.GaussianBlur{Sigma: blur},
			Interpretation: func() bimg.Interpretation {
				if gray == "true" {
					return bimg.InterpretationBW
				}
				return bimg.InterpretationSRGB
			}(),
		}

		// Auto format selection
		if format == "auto" {
			accept := reqRes.GetHeader("Accept")
			if accept != "" && (accept == "image/avif" || accept == "image/webp") {
				if accept == "image/avif" {
					options.Type = bimg.AVIF
				} else {
					options.Type = bimg.WEBP
				}
			} else {
				options.Type = bimg.JPEG
			}
		} else {
			switch format {
			case "jpeg":
				options.Type = bimg.JPEG
			case "png":
				options.Type = bimg.PNG
			case "webp":
				options.Type = bimg.WEBP
			case "avif":
				options.Type = bimg.AVIF
			default:
				options.Type = bimg.JPEG
			}
		}

		newImage, err := bimg.NewImage(imgBytes).Process(options)
		if err != nil {
			reqRes.JSON(http.StatusInternalServerError, gin.H{"error": "image processing failed"})
			return
		}

		var contentType string
		switch options.Type {
		case bimg.JPEG:
			contentType = "image/jpeg"
		case bimg.PNG:
			contentType = "image/png"
		case bimg.WEBP:
			contentType = "image/webp"
		case bimg.AVIF:
			contentType = "image/avif"
		default:
			contentType = "application/octet-stream"
		}

		reqRes.Data(http.StatusOK, contentType, newImage)
	})

	server.Run(":8080")
}
