package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
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
				"description": "Convertir TXT a PDF (soporta uno o múltiples archivos)",
				"parameters": gin.H{
					"file":        "Archivo(s) TXT (multipart form - repetir para múltiples)",
					"orientation": "auto | portrait | landscape (opcional, por defecto: auto)",
				},
				"examples": []string{
					"Archivo único: curl -F 'file=@documento.txt' http://localhost:8080/convert",
					"Múltiples: curl -F 'file=@doc1.txt' -F 'file=@doc2.txt' http://localhost:8080/convert",
				},
				"returns": "PDF binario (1 archivo) o ZIP (múltiples archivos)",
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

// handleConvert convierte TXT a PDF (soporta uno o múltiples archivos)
func (s *APIServer) handleConvert(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Error al procesar formulario multipart",
		})
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		c.JSON(400, gin.H{
			"error": "No se proporcionó archivo. Use: -F 'file=@documento.txt'",
		})
		return
	}

	// Obtener orientación (por defecto: auto)
	orientation := c.DefaultPostForm("orientation", "auto")
	if orientation == "auto" {
		orientation = "" // Se detectará para cada archivo
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

	// Procesar múltiples archivos
	var pdfResults []gin.H

	for _, file := range files {
		// Validar que sea TXT
		ext := filepath.Ext(file.Filename)
		if ext != ".txt" {
			pdfResults = append(pdfResults, gin.H{
				"file":  file.Filename,
				"error": "Solo se aceptan archivos .txt",
			})
			continue
		}

		// Abrir el archivo
		openFile, err := file.Open()
		if err != nil {
			pdfResults = append(pdfResults, gin.H{
				"file":  file.Filename,
				"error": "Error al leer el archivo",
			})
			continue
		}

		// Leer contenido en buffer
		content, err := io.ReadAll(openFile)
		openFile.Close()
		if err != nil {
			pdfResults = append(pdfResults, gin.H{
				"file":  file.Filename,
				"error": "Error al procesar el archivo",
			})
			continue
		}

		// Parsear líneas
		lines, err := ReadBuffer(content)
		if err != nil {
			pdfResults = append(pdfResults, gin.H{
				"file":  file.Filename,
				"error": "Error al parsear archivo",
			})
			continue
		}

		// Determinar orientación
		finalOrientation := orientation
		if finalOrientation == "" {
			finalOrientation = DetectOrientation(lines)
		}

		// Generar PDF en buffer
		pdfBuffer, err := GeneratePDFToBuffer(lines, file.Filename, finalOrientation)
		if err != nil {
			pdfResults = append(pdfResults, gin.H{
				"file":  file.Filename,
				"error": "Error al generar PDF",
			})
			continue
		}

		// Calcular hash del PDF
		pdfHash := CalculateBufferHash(pdfBuffer)

		// Si es un solo archivo, devolverlo directamente
		if len(files) == 1 {
			pdfFileName := file.Filename[:len(file.Filename)-4] + ".pdf"
			c.Header("Content-Type", "application/pdf")
			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", pdfFileName))
			c.Header("X-PDF-Hash", pdfHash)
			c.Header("X-PDF-Hash-Short", pdfHash[:16])
			c.Data(http.StatusOK, "application/pdf", pdfBuffer)
			return
		}

		// Para múltiples archivos, guardar en memoria
		pdfResults = append(pdfResults, gin.H{
			"file":       filepath.Base(file.Filename[:len(file.Filename)-4] + ".pdf"),
			"size_bytes": len(pdfBuffer),
			"sha256":     pdfHash,
			"short_hash": pdfHash[:16],
			"buffer":     pdfBuffer,
		})
	}

	// Si hay múltiples archivos, empaquetarlos en ZIP
	if len(files) > 1 {
		// Crear ZIP in-memory
		zipBuffer := new(bytes.Buffer)
		zipWriter := zip.NewWriter(zipBuffer)
		defer zipWriter.Close()

		successCount := 0
		for _, result := range pdfResults {
			// Saltar errores
			if _, ok := result["error"]; ok {
				continue
			}

			pdfFileName := result["file"].(string)
			pdfBuffer := result["buffer"].([]byte)

			// Agregar PDF al ZIP
			w, err := zipWriter.Create(pdfFileName)
			if err != nil {
				continue
			}
			_, err = w.Write(pdfBuffer)
			if err != nil {
				continue
			}
			successCount++
		}

		zipWriter.Close()

		// Enviar ZIP
		// Construir mapa de hashes de PDFs
		pdfHashes := make(map[string]string)
		for _, result := range pdfResults {
			if _, ok := result["error"]; ok {
				continue
			}
			pdfFileName := result["file"].(string)
			pdfHash := result["sha256"].(string)
			pdfHashes[pdfFileName] = pdfHash
		}
		// Serializar a JSON
		hashesJSON, err := json.Marshal(pdfHashes)
		if err == nil {
			c.Header("X-Pdf-Hashes", string(hashesJSON))
			c.Header("Access-Control-Expose-Headers", "X-Pdf-Hashes")
		}
		c.Header("Content-Type", "application/zip")
		c.Header("Content-Disposition", "attachment; filename=documentos.zip")
		c.Data(http.StatusOK, "application/zip", zipBuffer.Bytes())
		return
	}

	// Si algo salió mal con un solo archivo
	c.JSON(400, pdfResults[0])
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
