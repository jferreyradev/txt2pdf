package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func main() {
	var fileName string

	flag.StringVar(&fileName, "file", "", "Nombre del archivo a leer")
	flag.Parse()

	if fileName == "" {
		fmt.Println("Error: especifica el archivo con -file NOMBRE.txt")
		return
	}

	// Abrir archivo (busca en el directorio actual)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer file.Close()

	// Leer línea por línea
	scanner := bufio.NewScanner(file)
	lineNum := 0
	blankLines := 0
	pageBreaks := 0

	fmt.Printf("Archivo: %s\n", fileName)
	fmt.Println("=" + strings.Repeat("=", 100))

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Detectar saltos de página
		breakType := isPageBreak(line)

		if breakType != "" {
			fmt.Printf("%5d: [%s] %s\n", lineNum, breakType, line)
			pageBreaks++
		} else if len(line) == 0 {
			fmt.Printf("%5d: (vacía)\n", lineNum)
			blankLines++
		} else {
			fmt.Printf("%5d: %s\n", lineNum, line)
		}
	}

	// Resumen
	fmt.Println("=" + strings.Repeat("=", 100))
	fmt.Printf("Total de líneas:    %d\n", lineNum)
	fmt.Printf("Líneas en blanco:   %d\n", blankLines)
	fmt.Printf("Saltos de página:   %d\n", pageBreaks)
}
