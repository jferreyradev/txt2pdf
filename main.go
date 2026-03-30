package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

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

// calculateFileHash calcula el SHA256 de un archivo
func calculateFileHash(fileName string) (string, error) {
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

// generateAuditReport crea un archivo con registro de autenticidad
func generateAuditReport(inputDir string, auditFile string) error {
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
		hash, err := calculateFileHash(file)
		if err != nil {
			continue
		}

		auditContent.WriteString(fmt.Sprintf("Archivo: %s\n", baseName))
		auditContent.WriteString(fmt.Sprintf("Hash SHA256 TXT: %s\n", hash))

		// Calcular hash del PDF si existe
		pdfName := strings.TrimSuffix(baseName, ".txt") + ".pdf"
		pdfPath := filepath.Join(inputDir, pdfName)
		if pdfHash, err := calculateFileHash(pdfPath); err == nil {
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

	fmt.Printf("✓ Reporte de autenticidad generado: %s\n", auditFile)

	return nil
}

type LineEntry struct {
	Number    int
	BreakType string
	Content   string
	IsBlank   bool
}

func readFile(fileName string) ([]LineEntry, int, int, int, error) {
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

func generatePDF(fileName string, lines []LineEntry) error {
	pdfFileName := strings.TrimSuffix(fileName, ".txt") + ".pdf"

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)

	// Metadatos
	pdf.SetAuthor("Sistema de Procesamiento", true)
	pdf.SetCreator("TextIzer v1.0", true)
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

	fmt.Printf("✓ PDF generado: %s\n", pdfFileName)
	return nil
}

func main() {
	var fileName string
	var inputDir string
	var toPDF bool
	var processAll bool
	var genAudit bool
	var hashFile bool

	flag.StringVar(&fileName, "file", "", "Nombre del archivo a leer")
	flag.StringVar(&inputDir, "input", "./input", "Directorio con archivos .txt")
	flag.BoolVar(&toPDF, "pdf", false, "Generar salida en PDF")
	flag.BoolVar(&processAll, "all", false, "Procesar todos los archivos .txt del directorio")
	flag.BoolVar(&genAudit, "audit", false, "Generar reporte de autenticidad con hashes")
	flag.BoolVar(&hashFile, "hash", false, "Calcular hash SHA256 de archivo(s)")
	flag.Parse()

	// Crear directorio input si no existe
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		fmt.Printf("📁 Creando directorio: %s\n", inputDir)
		os.MkdirAll(inputDir, 0755)
	}

	// Calcular hash si se indica
	if hashFile {
		if fileName != "" {
			// Hash de un archivo específico
			hash, err := calculateFileHash(fileName)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("SHA256: %s\n", hash)
			fmt.Printf("Hash corto: %s\n", hash[:16])
		} else if processAll {
			// Hash de todos los PDFs
			files, err := filepath.Glob(filepath.Join(inputDir, "*.pdf"))
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			if len(files) == 0 {
				fmt.Println("No se encontraron archivos .pdf")
				return
			}

			fmt.Println("=== HASHES DE PDFS ===\n")
			for _, file := range files {
				baseName := filepath.Base(file)
				hash, err := calculateFileHash(file)
				if err != nil {
					fmt.Printf("  ✗ %s: Error\n", baseName)
					continue
				}
				fmt.Printf("Archivo: %s\n", baseName)
				fmt.Printf("SHA256: %s\n", hash)
				fmt.Printf("Hash corto: %s\n\n", hash[:16])
			}
		} else {
			fmt.Println("Use: -file archivo -hash  OR  -all -hash -input ./carpeta")
		}
		return
	}

	// Generar reporte de auditoría si se indica
	if genAudit {
		auditFile := filepath.Join(inputDir, "hashes.txt")
		err := generateAuditReport(inputDir, auditFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		return
	}

	// Procesar todos los archivos si se indica
	if processAll {
		files, err := filepath.Glob(filepath.Join(inputDir, "*.txt"))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(files) == 0 {
			fmt.Println("No se encontraron archivos .txt")
			return
		}

		fmt.Printf("Procesando %d archivos...\n\n", len(files))

		for _, file := range files {
			baseName := filepath.Base(file)

			// Ignorar archivo de hashes.txt
			if baseName == "hashes.txt" {
				continue
			}

			fmt.Printf("Procesando: %s\n", baseName)

			lines, _, _, _, err := readFile(file)
			if err != nil {
				fmt.Printf("  ✗ Error: %v\n\n", err)
				continue
			}

			if toPDF {
				err := generatePDF(file, lines)
				if err != nil {
					fmt.Printf("  ✗ Error al generar PDF: %v\n\n", err)
					continue
				}
			} else {
				fmt.Printf("  → %d líneas\n\n", len(lines))
			}
		}

		// Generar reporte de autenticidad automáticamente después de procesar PDFs
		if toPDF {
			fmt.Println("\nGenerando reporte de autenticidad...")
			auditFile := filepath.Join(inputDir, "hashes.txt")
			err := generateAuditReport(inputDir, auditFile)
			if err != nil {
				fmt.Printf("  ✗ Error al generar reporte: %v\n", err)
			}
		}
		return
	}

	// Procesar archivo único
	if fileName == "" {
		appName := filepath.Base(os.Args[0])
		fmt.Printf("\n🚀 %s - Convertidor de TXT a PDF con Autenticidad\n\n", strings.TrimSuffix(strings.ToUpper(appName), ".EXE"))
		fmt.Println("════════════════════════════════════════════════════════════════")
		fmt.Println("GUÍA RÁPIDA (3 comandos principales):")
		fmt.Println("════════════════════════════════════════════════════════════════\n")

		fmt.Println("1️⃣  Convertir UN archivo:")
		fmt.Printf("  %s -file input/documento.txt -pdf\n", appName)
		fmt.Println("  → Genera: documento.pdf + actualiza hashes.txt\n")

		fmt.Println("2️⃣  Convertir TODOS los archivos:")
		fmt.Printf("  %s -all -pdf\n", appName)
		fmt.Println("  → Genera: Todos los PDFs + actualiza hashes.txt\n")

		fmt.Println("3️⃣  Verificar integridad de PDF:")
		fmt.Printf("  %s -file documento.pdf -hash\n", appName)
		fmt.Println("  → Calcula: SHA256 + hash corto\n")

		fmt.Println("════════════════════════════════════════════════════════════════")
		fmt.Println("COMANDOS ADICIONALES:")
		fmt.Println("════════════════════════════════════════════════════════════════\n")

		fmt.Printf("  %s -all -hash\n", appName)
		fmt.Println("    Calcular hashes de TODOS los PDFs\n")

		fmt.Printf("  %s -all -pdf -input ./otra_carpeta\n", appName)
		fmt.Println("    Procesar archivos de otra carpeta\n")

		fmt.Printf("  %s -file documento.txt\n", appName)
		fmt.Println("    Solo leer y mostrar contenido (sin generar PDF)\n")

		fmt.Println("════════════════════════════════════════════════════════════════")
		fmt.Println()
		return
	}

	// Leer archivo
	lines, totalLines, blankLines, pageBreaks, err := readFile(fileName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Generar PDF si se solicita
	if toPDF {
		err := generatePDF(fileName, lines)
		if err != nil {
			fmt.Printf("Error al generar PDF: %v\n", err)
			return
		}

		// Generar reporte de autenticidad automáticamente
		fmt.Println("Generando reporte de autenticidad...")
		auditFile := filepath.Join(inputDir, "hashes.txt")
		err = generateAuditReport(inputDir, auditFile)
		if err != nil {
			fmt.Printf("  ✗ Error al generar reporte: %v\n", err)
		}
	} else {
		// Mostrar en consola
		fmt.Printf("Archivo: %s\n", fileName)
		fmt.Println("=" + strings.Repeat("=", 100))

		for _, entry := range lines {
			if entry.BreakType != "" {
				fmt.Printf("%5d: [%s] %s\n", entry.Number, entry.BreakType, entry.Content)
			} else if entry.IsBlank {
				fmt.Printf("%5d: (vacía)\n", entry.Number)
			} else {
				fmt.Printf("%5d: %s\n", entry.Number, entry.Content)
			}
		}

		fmt.Println("=" + strings.Repeat("=", 100))
		fmt.Printf("Total de líneas:    %d\n", totalLines)
		fmt.Printf("Líneas en blanco:   %d\n", blankLines)
		fmt.Printf("Saltos de página:   %d\n", pageBreaks)
	}
}
