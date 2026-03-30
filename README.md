# Text-to-PDF Converter con Autenticidad

Herramienta Go para convertir archivos TXT a PDF con **autenticidad verificable** mediante hash SHA256 embebido en cada página.

## Características

✅ Convierte TXT a PDF (línea por línea)  
✅ Detecta **saltos de página** (Form Feed)  
✅ **Hash SHA256** en archivo separado `hashes.txt`  
✅ **Logo watermark** semi-transparente  
✅ **Metadata** en PDF (Autor, Título, Creador, etc.)  
✅ **Reporte de autenticidad** con todos los hashes  
✅ Procesamiento **batch** para múltiples archivos  
✅ **Generación automática** de hashes al crear PDFs  
✅ **Calculadora de hashes** para verificación offline  
✅ Verificación simple sin herramientas adicionales

## Instalación

**Requisitos previos:**
- Go 1.21 o superior instalado

**Compilar el ejecutable:**
```bash
go build -o txt2pdf.exe
```

Listo. El ejecutable se crear en el mismo directorio.

## ✨ Características

✅ Convierte TXT a PDF (línea por línea)  
✅ Detecta **saltos de página** (Form Feed)  
✅ **Hash SHA256** en archivo separado `hashes.txt`  
✅ **Logo watermark** semi-transparente  
✅ **Metadata** en PDF (Autor, Título, Creador, etc.)  
✅ **Reporte de autenticidad** con todos los hashes  
✅ Procesamiento **batch** para múltiples archivos  
✅ **Generación automática** de hashes al crear PDFs  
✅ **Calculadora de hashes** para verificación offline  
✅ Verificación simple sin herramientas adicionales  
✅ **Auto-creación** de carpetas necesarias

## 🚀 Guía Rápida (2 pasos)

### Paso 1: Convertir archivos a PDF
```bash
.\txt2pdf.exe -all -pdf
```
Convierte todos los `.txt` en `input/` a PDF **y automáticamente genera `hashes.txt`**.

### Paso 2: Verificar documentos
1. Abre cualquier PDF generado
2. En el **footer (pie de página)** verás: fecha y número de página
3. Ejecuta `.\txt2pdf.exe -file documento.pdf -hash` para calcular su hash
4. Compara con el `hash corto` en `hashes.txt` → ✅ Auténtico

### Casos de uso comunes

| Necesidad | Comando |
|-----------|---------|
| Convertir 1 archivo (+ hashes) | `.\txt2pdf.exe -file archivo.txt -pdf` |
| Convertir todos (+ hashes) | `.\txt2pdf.exe -all -pdf` |
| Otra carpeta | `.\txt2pdf.exe -all -pdf -input ./carpeta` |
| Calcular hash de PDF | `.\txt2pdf.exe -file documento.pdf -hash` |
| Hashes de todos los PDFs | `.\txt2pdf.exe -all -hash` |
| Solo leer archivo | `.\txt2pdf.exe -file archivo.txt` |
| Ver opciones | `.\txt2pdf.exe` |

**Nota:** El archivo `hashes.txt` se genera **automáticamente** después de crear PDFs. Ya no necesitas ejecutar `-audit` por separado.

## Uso Detallado

### 1. Convertir un archivo TXT a PDF
```bash
.\txt2pdf.exe -file input/ARCHIVO.txt -pdf
```
✅ Genera `ARCHIVO.pdf`  
✅ Automáticamente actualiza `hashes.txt` con los hash SHA256 del documento

### 2. Convertir todos los archivos
```bash
.\txt2pdf.exe -all -pdf
```
✅ Procesa todos los `.txt` de la carpeta `input/`  
✅ Automáticamente actualiza `hashes.txt` con todos los hash SHA256

### 3. Especificar carpeta personalizada
```bash
.\txt2pdf.exe -all -pdf -input ./mi_carpeta
```
✅ Genera PDFs en carpeta personalizada  
✅ Crea `hashes.txt` en la misma carpeta

### 5. Calcular hash de un PDF
```bash
# Calcular hash de un PDF específico
.\txt2pdf.exe -file documento.pdf -hash

# Calcular hash de todos los PDFs
.\txt2pdf.exe -all -hash

# Otra carpeta
.\txt2pdf.exe -all -hash -input ./carpeta
```

## Verificación de Autenticidad

### 🔐 Sistema Simples y Efectivo

El archivo **`hashes.txt`** contiene:
- **Hash SHA256 del TXT original**: Valida que el documento fuente no fue alterado
- **Hash SHA256 del PDF**: Valida que el PDF generado no fue modificado
- **Hash corto**: Primeros 16 caracteres para referencia visual rápida

```
=== REGISTRO DE AUTENTICIDAD DE DOCUMENTOS ===
Generado: 2026-03-30 10:56:44

Archivo: DOCUMENTO.txt
Hash SHA256 TXT: fa7db23065f80f212769a7bb18f8d21854ea2d2216d8e321af727e6feee0b39b
Hash SHA256 PDF: 8b4c73e5a2f9d1c6e4b7f3a9c2e8d5a1f9b3c6e2a5d7f1e4c8b2a6d9f3c5e8
Hash corto PDF: 8b4c73e5a2f9d1c6
PDF: DOCUMENTO.pdf
```

### 🔍 ¿Cómo verificar que un PDF es auténtico?

**Verificación automática:**
```bash
# Calcular hash del PDF
.\txt2pdf.exe -file documento.pdf -hash

# Resultado:
# SHA256: 8b4c73e5a2f9d1c6e4b7f3a9c2e8d5a1f9b3c6e2a5d7f1e4c8b2a6d9f3c5e8
# Hash corto: 8b4c73e5a2f9d1c6

# Compara con hashes.txt
# Si coincide → ✅ PDF auténtico
# Si difiere → ❌ Alteración detectada
```

### ¿Por qué es efectivo?

- Si alguien modifica el PDF → el hash cambiaría inmediatamente
- El `hashes.txt` almacenado separado revela cualquier alteración
- Detecta cambios accidentales y modificaciones con herramientas
- Se puede verificar **sin dependencias** (solo necesitas calcular SHA256)

## 🔒 Modelo de Seguridad

| Componente | Función | Uso |
|-----------|----------|-----|
| **PDF con watermark** | Documento procesado | Lectura y distribución |
| **hashes.txt** | Registro centralizado de integridad | Auditoría y verificación |
| **-hash flag** | Calculadora independiente | Validación posterior |

**Limitaciones por diseño:**
- ⚠️ Valida integridad del documento (contra alteraciones)
- ⚠️ Para máxima seguridad, mantén `hashes.txt` en lugar protegido
- ⚠️ No autentica la identidad del autor (se requeriría certificado digital)

## Estructura del Proyecto

```
txt2pdf/
├── main.go                 (código principal)
├── go.mod                  (módulo Go)
├── go.sum                  (checksums de dependencias)
├── README.md              (este archivo)
├── logo/
│   └── logo_dgs.png       (logo watermark)
├── input/                 (archivos de entrada)
│   ├── *.txt             (archivos TXT fuente)
│   ├── *.pdf             (PDFs generados)
│   └── hashes.txt        (reporte de hashes)
└── txt2pdf.exe           (ejecutable compilado)
```

## Dependencias

- `github.com/jung-kurt/gofpdf` - Generación de PDF

## Características Técnicas

- **PDF Orientation**: Landscape A4
- **Font**: Courier 7pt
- **Footer**: Fecha | Hash (primeros 16 caracteres) | Página N
- **Logo**: Semi-transparente (50% opacidad)
- **Hash Algorithm**: SHA256
- **Page Break Detection**: Form Feed (FF) character
- **Encoding**: UTF-8

## Ejemplo Completo

```bash
# 1. Compilar
go build -o txt2pdf.exe

# 2. Convertir todos los TXT a PDF
.\txt2pdf.exe -all -pdf
# Resultado: input/*.pdf (con hash en footer)

# 3. Generar reporte de autenticidad
.\txt2pdf.exe -audit
# Resultado: input/hashes.txt (con todos los hashes)

# 4. Verificar un documento
.\txt2pdf.exe -file input/documento.pdf -hash
# Resultado: SHA256 completo + hash corto

# 5. Verificar todos
.\txt2pdf.exe -all -hash
# Resultado: Tabla con todos los hashes

# VERIFICACIÓN: Compara hashes con hashes.txt
# Si coinciden → ✅ Documentos auténticos
```

## Flujo de Trabajo Recomendado

1. **Generar PDFs**
   ```bash
   .\txt2pdf.exe -all -pdf
   ```

2. **Crear reporte**
   ```bash
   .\txt2pdf.exe -audit
   # Guardar hashes.txt en lugar seguro (USB, email, nube)
   ```

3. **Verificar integridad posterior**
   ```bash
   # Si sospechas alteración
   .\txt2pdf.exe -file documento_sospechoso.pdf -hash
   # Compara con hashes.txt original
   ```

4. **Compartir documentos**
   - Enviar PDF + referencia a hashes.txt
   - Destinatario puede verificar con: `.\txt2pdf.exe -file documento.pdf -hash`

## Preguntas Frecuentes

**P: ¿Puedo editar un PDF después de generarlo?**
R: Sí, pero el hash en `hashes.txt` no coincidirá con el nuevo PDF. La alteración será detectada.

**P: ¿Qué pasa si pierdo hashes.txt?**
R: Sin el archivo de referencia, no podrás verificar si el PDF fue alterado. Guárdalo en lugar seguro.

**P: ¿Es seguro contra expertos?**
R: Valida contra cambios accidentales y herramientas básicas. Para máxima seguridad se requiere firma digital certificada.

**P: ¿Funciona en Windows/Linux/Mac?**
R: Sí, es código Go puro. Solo necesitas compilar: `go build -o txt2pdf`
