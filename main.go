package main

import (
	"bufio"
	"flag"
	"fmt"
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
		logoPath := "logo_dgs.png"
		if _, err := os.Stat(logoPath); err == nil {
			// Guardar posición actual
			x, y := pdf.GetXY()

			// Hacer la imagen semi-transparente (30% opacidad)
			pdf.SetAlpha(0.5, "Normal")

			// Posición esquina superior derecha: x=240, y=3
			// Tamaño: 40x10mm
			pdf.Image(logoPath, 240, 3, 40, 10, false, "", 0, "")

			// Restaurar opacidad normal
			pdf.SetAlpha(1.0, "Normal")

			// Restaurar posición original sin agregar líneas
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

		// Mostrar solo el contenido sin números ni marcas
		lineStr := entry.Content

		// Usar MultiCell para que envuelva automáticamente líneas largas
		pdf.MultiCell(0, 2, lineStr, "", "L", false)
	}

	// Guardar PDF
	pdfFileName := strings.TrimSuffix(fileName, ".txt") + ".pdf"
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

	flag.StringVar(&fileName, "file", "", "Nombre del archivo a leer")
	flag.StringVar(&inputDir, "input", "./input", "Directorio con archivos .txt")
	flag.BoolVar(&toPDF, "pdf", false, "Generar salida en PDF")
	flag.BoolVar(&processAll, "all", false, "Procesar todos los archivos .txt del directorio")
	flag.Parse()

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
		return
	}

	// Procesar archivo único
	if fileName == "" {
		appName := filepath.Base(os.Args[0])
		fmt.Printf("\n%s - Convertidor de TXT a PDF\n\n", strings.TrimSuffix(strings.ToUpper(appName), ".EXE"))
		fmt.Println("Uso:")
		fmt.Printf("  %s -file archivo.txt          (leer un archivo)\n", appName)
		fmt.Printf("  %s -file archivo.txt -pdf    (generar PDF)\n", appName)
		fmt.Printf("  %s -all -pdf                 (generar PDFs de todos)\n", appName)
		fmt.Printf("  %s -all -input ./mi_carpeta  (procesar todos de otra carpeta)\n\n", appName)
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
