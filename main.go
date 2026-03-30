package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var fileName string
	var inputDir string
	var toPDF bool
	var processAll bool
	var genAudit bool
	var hashFile bool
	var forcePortrait bool
	var forceLandscape bool
	var apiMode bool
	var apiPort string

	flag.StringVar(&fileName, "file", "", "Nombre del archivo a leer")
	flag.StringVar(&inputDir, "input", "./input", "Directorio con archivos .txt")
	flag.BoolVar(&toPDF, "pdf", false, "Generar salida en PDF")
	flag.BoolVar(&processAll, "all", false, "Procesar todos los archivos .txt del directorio")
	flag.BoolVar(&genAudit, "audit", false, "Generar reporte de autenticidad con hashes")
	flag.BoolVar(&hashFile, "hash", false, "Calcular hash SHA256 de archivo(s)")
	flag.BoolVar(&forcePortrait, "portrait", false, "Forzar orientación vertical (Portrait)")
	flag.BoolVar(&forceLandscape, "landscape", false, "Forzar orientación horizontal (Landscape)")
	flag.BoolVar(&apiMode, "api", false, "Iniciar en modo servidor API REST")
	flag.StringVar(&apiPort, "port", "8080", "Puerto para el servidor API")
	flag.Parse()

	// Modo API
	if apiMode {
		fmt.Println("🚀 txt2pdf API REST iniciado")
		fmt.Printf("🌐 Servidor escuchando en http://localhost:%s\n", apiPort)
		fmt.Printf("📚 Documentación: http://localhost:%s/help\n", apiPort)
		fmt.Printf("🔗 Status: http://localhost:%s/status\n", apiPort)
		fmt.Println("⏹️  Presiona Ctrl+C para detener\n")

		server := NewAPIServer()
		if err := server.Start(apiPort); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Modo CLI
	// Crear directorio input si no existe
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		fmt.Printf("📁 Creando directorio: %s\n", inputDir)
		os.MkdirAll(inputDir, 0755)
	}

	// Calcular hash si se indica
	if hashFile {
		if fileName != "" {
			// Hash de un archivo específico
			hash, err := CalculateFileHash(fileName)
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
				hash, err := CalculateFileHash(file)
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
		err := GenerateAuditReport(inputDir, auditFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("✓ Reporte de autenticidad generado: %s\n", auditFile)
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

			lines, _, _, _, err := ReadFile(file)
			if err != nil {
				fmt.Printf("  ✗ Error: %v\n\n", err)
				continue
			}

			if toPDF {
				// Determinar orientación
				var orientation string
				if forcePortrait {
					orientation = "P"
				} else if forceLandscape {
					orientation = "L"
				} else {
					// Por defecto: auto-detecta
					orientation = DetectOrientation(lines)
				}

				err := GeneratePDFToFile(file, lines, orientation)
				if err != nil {
					fmt.Printf("  ✗ Error al generar PDF: %v\n\n", err)
					continue
				}
				fmt.Printf("✓ PDF generado: %s\n\n", strings.TrimSuffix(baseName, ".txt")+".pdf")
			} else {
				fmt.Printf("  → %d líneas\n\n", len(lines))
			}
		}

		// Generar reporte de autenticidad automáticamente después de procesar PDFs
		if toPDF {
			fmt.Println("\nGenerando reporte de autenticidad...")
			auditFile := filepath.Join(inputDir, "hashes.txt")
			err := GenerateAuditReport(inputDir, auditFile)
			if err != nil {
				fmt.Printf("  ✗ Error al generar reporte: %v\n", err)
			} else {
				fmt.Printf("✓ Reporte generado: %s\n", auditFile)
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
		fmt.Println("ORIENTACIÓN DEL PDF (opcional):")
		fmt.Println("════════════════════════════════════════════════════════════════\n")

		fmt.Printf("  %s -file documento.txt -pdf -portrait\n", appName)
		fmt.Println("    Fuerza orientación vertical (Portrait)\n")

		fmt.Printf("  %s -file documento.txt -pdf -landscape\n", appName)
		fmt.Println("    Fuerza orientación horizontal (Landscape)\n")

		fmt.Println("════════════════════════════════════════════════════════════════")
		fmt.Println("MODO API REST:")
		fmt.Println("════════════════════════════════════════════════════════════════\n")

		fmt.Printf("  %s -api -port 8080\n", appName)
		fmt.Println("    Iniciar servidor API REST en puerto 8080\n")

		fmt.Printf("  curl -F 'file=@documento.txt' http://localhost:8080/convert\n")
		fmt.Println("    Convertir TXT a PDF vía API\n")

		fmt.Println("════════════════════════════════════════════════════════════════")
		fmt.Println()
		return
	}

	// Leer archivo
	lines, totalLines, blankLines, pageBreaks, err := ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Generar PDF si se solicita
	if toPDF {
		// Determinar orientación
		var orientation string
		if forcePortrait {
			orientation = "P"
		} else if forceLandscape {
			orientation = "L"
		} else {
			// Por defecto: auto-detecta
			orientation = DetectOrientation(lines)
		}

		err := GeneratePDFToFile(fileName, lines, orientation)
		if err != nil {
			fmt.Printf("Error al generar PDF: %v\n", err)
			return
		}
		fmt.Printf("✓ PDF generado: %s\n", strings.TrimSuffix(fileName, ".txt")+".pdf")

		// Generar reporte de autenticidad automáticamente
		fmt.Println("Generando reporte de autenticidad...")
		// Usar el directorio del archivo si se especificó -file, sino usar inputDir
		reportDir := inputDir
		if fileName != "" {
			reportDir = filepath.Dir(fileName)
			if reportDir == "." {
				reportDir = "."
			}
		}
		auditFile := filepath.Join(reportDir, "hashes.txt")
		err = GenerateAuditReport(reportDir, auditFile)
		if err != nil {
			fmt.Printf("  ✗ Error al generar reporte: %v\n", err)
		} else {
			fmt.Printf("✓ Reporte generado: %s\n", auditFile)
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
