# Text Analyzer

Herramienta para analizar archivos de texto y visualizar saltos de línea y página.

## Características

✅ Detecta **números de página** (1, 2, 3...)
✅ Detecta **Form Feed** (`\f`)
✅ Detecta **marcadores PAGE BREAK**
✅ Reporta **líneas en blanco**
✅ Muestra contexto alrededor de cada salto
✅ Reporte formateado en terminal

## Instalación

```bash
go build -o analyzer.exe
```

## Uso

```bash
# Analizar archivos en input/
.\analyzer.exe

# O especificar directorio personalizado
.\analyzer.exe -input="C:\ruta\a\archivos"
```

## Salida

Genera un reporte detallado con:
- Número de línea del salto
- Tipo de salto (PAGE_NUM, FF, MARKER, BLANK)
- Posición en bytes
- Contexto (línea anterior/siguiente)

## Estructura

```
text-analyzer/
├── main.go
├── go.mod
├── README.md
└── input/          (tus archivos .txt)
```
