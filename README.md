# Text-to-PDF Converter con Autenticidad

Herramienta Go para convertir archivos TXT a PDF con **autenticidad verificable** mediante QR code y hash SHA256.

## Características

✅ Convierte TXT a PDF (línea por línea)  
✅ Detecta **saltos de página** (Form Feed)  
✅ **Código QR de autenticidad** en última página  
✅ **Hash SHA256** embebido en QR  
✅ **Logo watermark** semi-transparente  
✅ **Metadata** en PDF (Autor, Título, Creador, etc.)  
✅ **Reporte de autenticidad** con todos los hashes  
✅ Procesamiento **batch** para múltiples archivos  
✅ **Verificación simple** sin PowerShell required  

## Instalación

```bash
go build -o txt2pdf.exe
```

## Uso

### 1. Convertir un archivo TXT a PDF
```bash
.\txt2pdf.exe -file input/LIBRAMIENTOS.txt -pdf
```

### 2. Convertir todos los archivos a PDF
```bash
.\txt2pdf.exe -all -pdf
```

### 3. Generar reporte de autenticidad
```bash
.\txt2pdf.exe -audit
```
Genera:
- `input/autenticidad.txt` - Tabla con todos los hashes
- `input/autenticidad.pdf` - Reporte en formato PDF

### 4. Especificar carpeta de entrada personalizada
```bash
.\txt2pdf.exe -all -pdf -input ./mi_carpeta
.\txt2pdf.exe -audit -input ./mi_carpeta
```

## Verificación de Autenticidad

### Paso 1: Escanear el QR en la última página del PDF
El QR contiene información en formato: `Archivo | Hash | Timestamp`

Ejemplo: `LIBRAMIENTOS.txt | fa7db23065f80f21 | 2026-03-27 14:23:15`

### Paso 2: Consultar el archivo `autenticidad.txt`
Abre `input/autenticidad.txt` generado con `-audit`

```
=== REGISTRO DE AUTENTICIDAD DE DOCUMENTOS ===
Generado: 2026-03-27 14:23:15

Archivo: LIBRAMIENTOS.txt
Hash SHA256: fa7db23065f80f212769a7bb18f8d21854ea2d2216d8e321af727e6feee0b39b
Hash QR: fa7db23065f80f21
PDF: LIBRAMIENTOS.pdf
```

### Paso 3: Comparar hashes
- **Hash del QR**: `fa7db23065f80f21` (primeros 16 caracteres)
- **Hash en autenticidad.txt**: `fa7db23065f80f21`

✅ **Si coinciden** → Documento es auténtico, sin modificaciones  
❌ **Si NO coinciden** → Archivo fue alterado después de generar el PDF

## Estructura del Proyecto

```
text-analyzer/
├── main.go                 (código principal)
├── go.mod                  (módulo Go)
├── go.sum                  (checksums de dependencias)
├── README.md              (este archivo)
│
├── logo/
│   └── logo_dgs.png       (logo watermark)
│
├── input/                 (archivos de entrada)
│   ├── *.txt             (archivos TXT source)
│   ├── *.pdf             (PDFs generados)
│   ├── autenticidad.txt  (reporte de autenticidad)
│   └── autenticidad.pdf  (reporte en PDF)
│
└── txt2pdf.exe           (ejecutable compilado)
```

## Dependencias

- `github.com/jung-kurt/gofpdf` - Generación de PDF
- `github.com/skip2/go-qrcode` - Códigos QR

## Características Técnicas

- **PDF Orientation**: Landscape A4
- **Font**: Courier 7pt
- **QR Position**: Esquina inferior derecha de última página
- **QR Size**: 20x20mm
- **Logo**: Semi-transparente (50% opacidad)
- **Hash Algorithm**: SHA256
- **QR Error Correction**: High level

## Ejemplo Completo

```bash
# 1. Generar todos los PDFs
.\txt2pdf.exe -all -pdf

# 2. Generar reporte de autenticidad
.\txt2pdf.exe -audit

# 3. Verificar un documento
# - Escanea QR en LIBRAMIENTOS.pdf última página
# - Lee autenticidad.txt
# - Compara hashes

# Resultado: Documento auténtico ✓
```
