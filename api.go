package main

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// APIServer encapsula el servidor Gin
type APIServer struct {
	router *gin.Engine
}

// NewAPIServer crea un nuevo servidor API
func NewAPIServer() *APIServer {
	router := gin.Default()
	server := &APIServer{router: router}

	// Habilitar CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Rutas
	router.POST("/convert", server.handleConvert)
	router.POST("/hash", server.handleHash)
	router.GET("/status", server.handleStatus)
	router.GET("/help", server.handleHelp)

	return server
}

// handleStatus devuelve el estado de la API
func (s *APIServer) handleStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "🟢 txt2pdf API activa",
		"version": "1.0",
		"endpoints": gin.H{
			"POST /convert": "Convertir TXT a PDF",
			"POST /hash":    "Calcular hash SHA256",
			"GET  /status":  "Estado de la API",
			"GET  /help":    "Documentación",
		},
	})
}

// handleHelp devuelve documentación
func (s *APIServer) handleHelp(c *gin.Context) {
	c.JSON(200, gin.H{
		"info": "txt2pdf API REST 🫶",
		"endpoints": []gin.H{
			{
				"endpoint":    "POST /convert",
				"description": "Convertir archivo TXT a PDF",
				"parameters": gin.H{
					"file":        "Archivo TXT (multipart form)",
					"orientation": "auto | portrait | landscape (opcional, por defecto: auto)",
				},
				"example": "curl -F 'file=@documento.txt' http://localhost:8080/convert",
				"returns": "PDF binario",
			},
			{
				"endpoint":    "POST /hash",
				"description": "Calcular hash SHA256 de archivo",
				"parameters": gin.H{
					"file": "Archivo a hashear (multipart form)",
				},
				"example": "curl -F 'file=@documento.pdf' http://localhost:8080/hash",
				"returns": gin.H{
					"filename":   "nombre del archivo",
					"sha256":     "hash completo",
					"short_hash": "primeros 16 caracteres",
				},
			},
		},
	})
}

// handleConvert convierte TXT a PDF
func (s *APIServer) handleConvert(c *gin.Context) {
	// Obtener el archivo del formulario
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "No se proporcionó archivo. Use: -F 'file=@documento.txt'",
		})
		return
	}

	// Validar que sea TXT
	ext := filepath.Ext(file.Filename)
	if ext != ".txt" {
		c.JSON(400, gin.H{
			"error": "Solo se aceptan archivos .txt",
			"got":   ext,
		})
		return
	}

	// Abrir el archivo
	openFile, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Error al leer el archivo",
		})
		return
	}
	defer openFile.Close()

	// Leer contenido en buffer
	content, err := io.ReadAll(openFile)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Error al procesar el archivo",
		})
		return
	}

	// Parsear líneas
	lines, err := ReadBuffer(content)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error al parsear archivo",
		})
		return
	}

	// Obtener orientación (por defecto: auto)
	orientation := c.DefaultPostForm("orientation", "")
	if orientation == "" {
		orientation = DetectOrientation(lines)
	} else if orientation == "auto" {
		orientation = DetectOrientation(lines)
	} else if orientation == "portrait" {
		orientation = "P"
	} else if orientation == "landscape" {
		orientation = "L"
	} else {
		c.JSON(400, gin.H{
			"error": "Orientación inválida. Use: auto, portrait, landscape",
			"got":   orientation,
		})
		return
	}

	// Generar PDF en buffer
	pdfBuffer, err := GeneratePDFToBuffer(lines, file.Filename, orientation)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Error al generar PDF",
		})
		return
	}

	// Calcular hash del PDF
	pdfHash := CalculateBufferHash(pdfBuffer)

	// Enviar PDF como descarga
	pdfFileName := filepath.Base(file.Filename[:len(file.Filename)-4]) + ".pdf"
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", pdfFileName))
	c.Header("X-PDF-Hash", pdfHash)
	c.Header("X-PDF-Hash-Short", pdfHash[:16])
	c.Data(http.StatusOK, "application/pdf", pdfBuffer)
}

// handleHash calcula el hash SHA256 de un archivo
func (s *APIServer) handleHash(c *gin.Context) {
	// Obtener el archivo del formulario
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "No se proporcionó archivo. Use: -F 'file=@documento.pdf'",
		})
		return
	}

	// Abrir el archivo
	openFile, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Error al leer el archivo",
		})
		return
	}
	defer openFile.Close()

	// Leer contenido en buffer
	content, err := io.ReadAll(openFile)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Error al procesar el archivo",
		})
		return
	}

	// Calcular hash
	hash := CalculateBufferHash(content)

	c.JSON(200, gin.H{
		"filename":   file.Filename,
		"sha256":     hash,
		"short_hash": hash[:16],
		"size_bytes": len(content),
	})
}

// Start inicia el servidor
func (s *APIServer) Start(port string) error {
	return s.router.Run(":" + port)
}
