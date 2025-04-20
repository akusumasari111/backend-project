package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"backend-project/config"

	"github.com/gin-gonic/gin"
)

const MaxUploadSize = 5 << 20 // 5 MB
const UploadPath = "./uploads"

func UploadProductImage(c *gin.Context) {
	id := c.Param("id")

	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check product existence"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded file"})
		return
	}
	defer file.Close()

	if !isValidImageFormat(header) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format. Only PNG, JPG, JPEG allowed."})
		return
	}

	if err := os.MkdirAll(UploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload folder"})
		return
	}

	filename := fmt.Sprintf("%s%s", id, filepath.Ext(header.Filename))
	filepath := filepath.Join(UploadPath, filename)

	out, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}

	_, err = config.DB.Exec("UPDATE products SET image_url = ? WHERE id = ?", filename, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded and database updated", "filename": filename})
}

func DownloadProductImage(c *gin.Context) {
	id := c.Param("id")

	var imageUrl sql.NullString
	err := config.DB.QueryRow("SELECT image_url FROM products WHERE id = ?", id).Scan(&imageUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		}
		return
	}

	if !imageUrl.Valid || imageUrl.String == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not uploaded for this product"})
		return
	}

	imagePath := filepath.Join(UploadPath, imageUrl.String)
	c.File(imagePath)
}

func DeleteProductImage(c *gin.Context) {
	id := c.Param("id")

	var imageUrl sql.NullString
	err := config.DB.QueryRow("SELECT image_url FROM products WHERE id = ?", id).Scan(&imageUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		}
		return
	}

	if !imageUrl.Valid || imageUrl.String == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found for this product"})
		return
	}

	imagePath := filepath.Join(UploadPath, imageUrl.String)
	if err := os.Remove(imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image file"})
		return
	}

	_, err = config.DB.Exec("UPDATE products SET image_url = NULL WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}

func isValidImageFormat(header *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	switch ext {
	case ".png", ".jpg", ".jpeg":
		return true
	default:
		return false
	}
}
