package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// LineEntry representa una línea del archivo
type LineEntry struct {
	Number    int
	BreakType string
	Content   string
	IsBlank   bool
}

// ConversionOptions contiene las opciones para la conversión
type ConversionOptions struct {
	ForcePortrait  bool
	ForceLandscape bool
}

// CalculateFileHash calcula el SHA256 de un archivo
func CalculateFileHash(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// CalculateBufferHash calcula el SHA256 de un buffer
func CalculateBufferHash(buffer []byte) string {
	hash := sha256.Sum256(buffer)
	return hex.EncodeToString(hash[:])
}

// isPageBreak detecta si una línea es un salto de página
func isPageBreak(line string) string {
	trimmed := strings.TrimSpace(line)

	// Form Feed
	if strings.ContainsAny(line, "\f\x0c") {
		return "FF"
	}

	// Marcador PAGE BREAK
	if strings.Contains(trimmed, "PAGE BREAK") {
		return "MARKER"
	}

	// Número de página (1-999)
	if len(trimmed) > 0 && len(trimmed) <= 3 {
		_, err := strconv.Atoi(trimmed)
		if err == nil {
			return "PAGE_NUM"
		}
	}

	return ""
}

// DetectOrientation analiza las líneas y detecta la mejor orientación
// Retorna "P" para Portrait (vertical) o "L" para Landscape (horizontal)
func DetectOrientation(lines []LineEntry) string {
	if len(lines) == 0 {
		return "L" // Por defecto Landscape
	}

	// Analizar primeras 100 líneas no vacías
	totalLen := 0
	count := 0
	maxCount := 100

	for _, entry := range lines {
		if count >= maxCount {
			break
		}
		if !entry.IsBlank && entry.BreakType == "" {
			totalLen += len(entry.Content)
			count++
		}
	}

	if count == 0 {
		return "L" // Si no hay líneas, usar Landscape
	}

	avgLen := totalLen / count

	// Si línea promedio > 80 caracteres → Landscape
	// Si línea promedio <= 80 caracteres → Portrait
	if avgLen > 80 {
		return "L" // Landscape para líneas largas
	}
	return "P" // Portrait para líneas cortas
}

// ReadFile lee un archivo y lo convierte en líneas
func ReadFile(fileName string) ([]LineEntry, int, int, int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []LineEntry
	lineNum := 0
	blankLines := 0
	pageBreaks := 0
	lastBreakType := ""

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		breakType := isPageBreak(line)

		// Si hay dos saltos de página seguidos, omitir el segundo
		if breakType != "" && lastBreakType != "" {
			breakType = "" // Omitir este salto
		}

		entry := LineEntry{
			Number:    lineNum,
			BreakType: breakType,
			Content:   line,
			IsBlank:   len(line) == 0,
		}

		if breakType != "" {
			pageBreaks++
			lastBreakType = breakType
		} else {
			lastBreakType = ""
		}

		if len(line) == 0 {
			blankLines++
		}

		lines = append(lines, entry)
	}

	return lines, lineNum, blankLines, pageBreaks, scanner.Err()
}

// ReadBuffer lee un buffer y lo convierte en líneas
func ReadBuffer(buffer []byte) ([]LineEntry, error) {
	reader := bytes.NewReader(buffer)
	scanner := bufio.NewScanner(reader)
	var lines []LineEntry
	lineNum := 0
	lastBreakType := ""

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		breakType := isPageBreak(line)

		// Si hay dos saltos de página seguidos, omitir el segundo
		if breakType != "" && lastBreakType != "" {
			breakType = "" // Omitir este salto
		}

		entry := LineEntry{
			Number:    lineNum,
			BreakType: breakType,
			Content:   line,
			IsBlank:   len(line) == 0,
		}

		if breakType != "" {
			lastBreakType = breakType
		} else {
			lastBreakType = ""
		}

		lines = append(lines, entry)
	}

	return lines, scanner.Err()
}

// GeneratePDFToBuffer genera un PDF en un buffer (sin guardar a disco)
func GeneratePDFToBuffer(lines []LineEntry, fileName string, orientation string) ([]byte, error) {
	// Si no se especifica orientación, auto-detectar
	if orientation == "" {
		orientation = DetectOrientation(lines)
	}

	pdf := gofpdf.New(orientation, "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)

	// Metadatos
	pdf.SetAuthor("Sistema de Procesamiento", true)
	pdf.SetCreator("txt2pdf v1.0", true)
	pdf.SetTitle(fileName, true)
	pdf.SetSubject("Conversión de TXT a PDF", true)
	pdf.SetKeywords("procesamiento, documentos", true)

	// Configurar encabezado con logo
	pdf.SetHeaderFunc(func() {
		// Logo como marca de agua - semi-transparente
		logoPath := "logo/logo_dgs.png"
		if _, err := os.Stat(logoPath); err == nil {
			// Guardar posición actual
			x, y := pdf.GetXY()

			// Hacer la imagen semi-transparente (50% opacidad)
			pdf.SetAlpha(0.5, "Normal")

			// Posición esquina superior derecha: x=240, y=3
			// Tamaño: 40x10mm
			pdf.Image(logoPath, 240, 3, 40, 10, false, "", 0, "")

			// Restaurar opacidad normal
			pdf.SetAlpha(1.0, "Normal")

			// Restaurar posición original
			pdf.SetXY(x, y)
		}
	})

	// Configurar pie de página
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Courier", "", 6)
		// Fecha y hora
		dateStr := time.Now().Format("02/01/2006 15:04:05")
		pdf.Cell(0, 10, dateStr)
		// Número de página
		pageStr := fmt.Sprintf("Página %d", pdf.PageNo())
		pdf.SetX(250)
		pdf.Cell(0, 10, pageStr)
	})

	pdf.AddPage()
	pdf.SetFont("Courier", "", 7)

	// Contenido
	needsNewPage := false

	for _, entry := range lines {
		// Si hay salto de página (solo respetar FF, ignorar PAGE_NUM)
		if entry.BreakType == "FF" {
			needsNewPage = true
			continue
		}

		// Ignorar otros tipos de saltos de página
		if entry.BreakType != "" {
			continue
		}

		// Saltar líneas en blanco
		if entry.IsBlank {
			continue
		}

		// Si necesitamos nueva página, crearla antes del contenido
		if needsNewPage {
			pdf.AddPage()
			pdf.SetFont("Courier", "", 7)
			needsNewPage = false
		}

		y := pdf.GetY()

		// Verificar si se necesita nueva página por altura (margen de seguridad)
		if y > 270 {
			pdf.AddPage()
		}

		// Usar MultiCell para que envuelva automáticamente líneas largas
		pdf.MultiCell(0, 2, entry.Content, "", "L", false)
	}

	// Generar en buffer
	buffer := new(bytes.Buffer)
	err := pdf.Output(buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// GeneratePDFToFile genera un PDF y lo guarda en disco
func GeneratePDFToFile(fileName string, lines []LineEntry, orientation string) error {
	pdfFileName := strings.TrimSuffix(fileName, ".txt") + ".pdf"

	// Si no se especifica orientación, auto-detectar
	if orientation == "" {
		orientation = DetectOrientation(lines)
	}

	pdf := gofpdf.New(orientation, "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)

	// Metadatos
	pdf.SetAuthor("Sistema de Procesamiento", true)
	pdf.SetCreator("txt2pdf v1.0", true)
	pdf.SetTitle(filepath.Base(fileName), true)
	pdf.SetSubject("Conversión de TXT a PDF", true)
	pdf.SetKeywords("procesamiento, documentos", true)

	// Configurar encabezado con logo
	pdf.SetHeaderFunc(func() {
		// Logo como marca de agua - semi-transparente
		logoPath := "logo/logo_dgs.png"
		if _, err := os.Stat(logoPath); err == nil {
			// Guardar posición actual
			x, y := pdf.GetXY()

			// Hacer la imagen semi-transparente (50% opacidad)
			pdf.SetAlpha(0.5, "Normal")

			// Posición esquina superior derecha: x=240, y=3
			// Tamaño: 40x10mm
			pdf.Image(logoPath, 240, 3, 40, 10, false, "", 0, "")

			// Restaurar opacidad normal
			pdf.SetAlpha(1.0, "Normal")

			// Restaurar posición original
			pdf.SetXY(x, y)
		}
	})

	// Configurar pie de página
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Courier", "", 6)
		// Fecha y hora
		dateStr := time.Now().Format("02/01/2006 15:04:05")
		pdf.Cell(0, 10, dateStr)
		// Número de página
		pageStr := fmt.Sprintf("Página %d", pdf.PageNo())
		pdf.SetX(250)
		pdf.Cell(0, 10, pageStr)
	})

	pdf.AddPage()
	pdf.SetFont("Courier", "", 7)

	// Contenido
	needsNewPage := false

	for _, entry := range lines {
		// Si hay salto de página (solo respetar FF, ignorar PAGE_NUM)
		if entry.BreakType == "FF" {
			needsNewPage = true
			continue
		}

		// Ignorar otros tipos de saltos de página
		if entry.BreakType != "" {
			continue
		}

		// Saltar líneas en blanco
		if entry.IsBlank {
			continue
		}

		// Si necesitamos nueva página, crearla antes del contenido
		if needsNewPage {
			pdf.AddPage()
			pdf.SetFont("Courier", "", 7)
			needsNewPage = false
		}

		y := pdf.GetY()

		// Verificar si se necesita nueva página por altura (margen de seguridad)
		if y > 270 {
			pdf.AddPage()
		}

		// Usar MultiCell para que envuelva automáticamente líneas largas
		pdf.MultiCell(0, 2, entry.Content, "", "L", false)
	}

	// Guardar PDF
	err := pdf.OutputFileAndClose(pdfFileName)
	if err != nil {
		return err
	}

	return nil
}

// GenerateAuditReport crea un archivo con registro de autenticidad
func GenerateAuditReport(inputDir string, auditFile string) error {
	files, err := filepath.Glob(filepath.Join(inputDir, "*.txt"))
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no se encontraron archivos .txt")
	}

	var auditContent strings.Builder
	auditContent.WriteString("=== REGISTRO DE AUTENTICIDAD DE DOCUMENTOS ===\n")
	auditContent.WriteString(fmt.Sprintf("Generado: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	auditContent.WriteString(strings.Repeat("=", 50) + "\n\n")

	for _, file := range files {
		baseName := filepath.Base(file)
		hash, err := CalculateFileHash(file)
		if err != nil {
			continue
		}

		auditContent.WriteString(fmt.Sprintf("Archivo: %s\n", baseName))
		auditContent.WriteString(fmt.Sprintf("Hash SHA256 TXT: %s\n", hash))

		// Calcular hash del PDF si existe
		pdfName := strings.TrimSuffix(baseName, ".txt") + ".pdf"
		pdfPath := filepath.Join(inputDir, pdfName)
		if pdfHash, err := CalculateFileHash(pdfPath); err == nil {
			auditContent.WriteString(fmt.Sprintf("Hash SHA256 PDF: %s\n", pdfHash))
			auditContent.WriteString(fmt.Sprintf("Hash corto PDF: %s\n", pdfHash[:16]))
		}

		auditContent.WriteString(fmt.Sprintf("PDF: %s\n", pdfName))
		auditContent.WriteString(strings.Repeat("-", 50) + "\n\n")
	}

	err = os.WriteFile(auditFile, []byte(auditContent.String()), 0644)
	if err != nil {
		return err
	}

	return nil
}
