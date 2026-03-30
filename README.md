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

## 🚀 Guía Rápida (3 pasos)

### Paso 1: Convertir archivos
```bash
.\txt2pdf.exe -all -pdf
```
Convierte todos los `.txt` en `input/` a PDF con QR de autenticidad.

### Paso 2: Generar reporte
```bash
.\txt2pdf.exe -audit
```
Crea `autenticidad.txt` con todos los hashes para verificar.

### Paso 3: Verificar documentos
1. Abre cualquier PDF generado
2. Escanea el QR en la última página (abajo a la derecha)
3. Busca ese hash en `autenticidad.txt` → ✅ Auténtico

**¡Listo!** Tus documentos están certificados y verificables.

### Casos de uso comunes

| Necesidad | Comando |
|-----------|---------|
| Convertir 1 archivo | `.\txt2pdf.exe -file archivo.txt -pdf` |
| Convertir todos | `.\txt2pdf.exe -all -pdf` |
| Otra carpeta | `.\txt2pdf.exe -all -pdf -input ./carpeta` |
| Verificar autenticidad | `.\txt2pdf.exe -audit` |
| Ver opciones | `.\txt2pdf.exe` |

## Uso

### 1. Convertir un archivo TXT a PDF
```bash
.\txt2pdf.exe -file input/ARCHIVO.txt -pdf
```

### 2. Convertir todos los archivos de la carpeta \input a PDF
```bash
.\txt2pdf.exe -all -pdf
```

### 3. Generar reporte de autenticidad
```bash
.\txt2pdf.exe -audit
```
Genera:
- `input/autenticidad.txt` - Tabla con todos los hashes para verificación

### 4. Especificar carpeta de entrada personalizada
```bash
.\txt2pdf.exe -all -pdf -input ./mi_carpeta
.\txt2pdf.exe -audit -input ./mi_carpeta
```

## Verificación de Autenticidad

### Paso 1: Escanear el QR en la última página del PDF
El QR contiene información en formato: `Archivo | Hash | Timestamp`

Ejemplo: `ARCHIVO.txt | fa7db23065f80f21 | 2026-03-27 14:23:15`

### Paso 2: Consultar el archivo `autenticidad.txt`
Abre `input/autenticidad.txt` generado con `-audit`

```
=== REGISTRO DE AUTENTICIDAD DE DOCUMENTOS ===
Generado: 2026-03-27 14:23:15

Archivo: ARCHIVO.txt
Hash SHA256 TXT: fa7db23065f80f212769a7bb18f8d21854ea2d2216d8e321af727e6feee0b39b
Hash QR: fa7db23065f80f21
Hash SHA256 PDF: 8b4c73e5a2f9d1c6e4b7f3a9c2e8d5a1f9b3c6e2a5d7f1e4c8b2a6d9f3c5e8
PDF: ARCHIVO.pdf
```

### Paso 3: Verificación en 3 capas

**Capa 1: Verificación Rápida**
- Escanea el QR en la última página del PDF
- Compara los 16 primeros caracteres con el reporte

**Capa 2: Verificación Completa TXT**
- Compara `Hash SHA256 TXT` del reporte con el archivo .txt original
- Si coincide → archivo origen no fue alterado

**Capa 3: Verificación del PDF**
- Compara `Hash SHA256 PDF` del reporte con el PDF
- Si coincide → PDF embebido es idéntico al generado
- Si NO coincide → PDF fue modificado después de su generación

### Comparación de hashes
- **QR**: `fa7db23065f80f21` (primeros 16 caracteres del TXT)
- **Reporte TXT**: `fa7db23065f80f212769a7bb18f8d21854ea2d2216d8e321af727e6feee0b39b` ✅ Coincide
- **Reporte PDF**: `8b4c73e5a2f9d1c6e4b7f3a9c2e8d5a1f9b3c6e2a5d7f1e4c8b2a6d9f3c5e8` (verificar PDF no fue modificado)

**Resultado final:**
- ✅ **TXT auténtico** si hashes TXT coinciden
- ✅ **PDF íntegro** si hashes PDF coinciden
- ❌ **Alteración detectada** si algún hash NO coincide

## 🔒 Modelo de Seguridad

El sistema implementa una **cadena de verificación triple**:

| Nivel | Propósito | Verifica |
|-------|----------|----------|
| **QR (Rápida)** | Verificación visual en el PDF | Primeros 16 caracteres del TXT |
| **Reporte (Estándar)** | Verificación offline separada | Hash completo TXT + PDF |
| **Integridad PDF** | Detección de manipulación | Cambios en el documento generado |

**Ventajas de este diseño:**
- ✅ PDF con QR para verificación rápida (práctico)
- ✅ Reporte separado e inmutable (protección adicional)
- ✅ Hash del PDF detecta cualquier modificación posterior
- ✅ No requiere firma digital ni infraestructura adicional
- ✅ Compatible con cualquier lector PDF estándar

**Limitaciones (por diseño):**
- ⚠️ Valida contra alteraciones accidentales o herramientas básicas
- ⚠️ Requiere que el archivo `autenticidad.txt` se conserve
- ⚠️ No prueba identidad del autor (sin certificados digitales)

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
