# 🚀 QUICKSTART - txt2pdf

Guía rápida para empezar a usar la herramienta en 5 minutos.

---

## Paso 1: Preparar el ejecutable

### Opción A: Usar el ejecutable ya compilado
Si tienes `txt2pdf.exe` en la carpeta, **salta al Paso 2**.

### Opción B: Compilar (solo una vez)
Requiere: [Go instalado](https://go.dev/dl)

```bash
go build -o txt2pdf.exe
```

---

## Paso 2: Preparar tus documentos

**Opción A: Usar la carpeta por defecto `input/`**
1. Crear carpeta `input/` (se crea automáticamente si la necesitas)
2. Coloca archivos `.txt` en `input/`

**Opción B: Usar carpeta personalizada**
```bash
.\txt2pdf.exe -all -pdf -input ./mi_carpeta_personalizada
```
- Se crea automáticamente si no existe
- Todos los archivos `.txt` en esa carpeta se procesan

```
Ejemplo con carpeta personalizada:
├── txt2pdf.exe
└── documentos/           ← Tu carpeta personalizada
    ├── reporte1.txt
    ├── reporte2.txt
    ├── reporte1.pdf      ← Se genera aquí
    ├── reporte2.pdf      ← Se genera aquí
    └── hashes.txt        ← Se genera aquí
```

---

## Paso 3: Generar PDFs

### Con carpeta por defecto:
```bash
.\txt2pdf.exe -all -pdf
```
Procesa todos los `.txt` en `input/`

### Con carpeta personalizada:
```bash
.\txt2pdf.exe -all -pdf -input ./documentos
```
Procesa todos los `.txt` en `./documentos/`

### Convertir UN archivo específico:
```bash
.\txt2pdf.exe -file documento.txt -pdf
# O de otra carpeta:
.\txt2pdf.exe -file ./documentos/reporte.txt -pdf
```

### ⭐ Orientación del PDF (Nuevo)

**Por defecto: Auto-detecta automáticamente**
```bash
# Simplemente genera PDF - analiza y elige orientación automáticamente
.\txt2pdf.exe -file documento.txt -pdf
.\txt2pdf.exe -all -pdf
```

**Fuerza orientación específica (opcional):**
```bash
# Vertical (Portrait) - fuerza líneas cortas
.\txt2pdf.exe -file documento.txt -pdf -portrait

# Horizontal (Landscape) - fuerza líneas largas
.\txt2pdf.exe -file documento.txt -pdf -landscape
```

**Cómo funciona la auto-detección (por defecto):**
- Analiza primeras 100 líneas del documento
- Si línea promedio ≤ 80 caracteres → Portrait
- Si línea promedio > 80 caracteres → Landscape

---

## Paso 4: Verificar integridad

¿Quieres asegurarte de que un PDF no fue alterado?

### Ver hash de UN PDF:
```bash
.\txt2pdf.exe -file documento1.pdf -hash
```

**Resultado:**
```
SHA256: 4f790750acda0983c5313eded002b470212468bb56608e557fe3ac6af9c16369
Hash corto: 4f790750acda0983
```

Compara el **hash corto** con el que aparece en `input/hashes.txt` → ✅ Es auténtico

### Ver hashes de TODOS los PDFs:
```bash
.\txt2pdf.exe -all -hash
```

---

## 📋 Archivo de Hashes (`hashes.txt`)

Se genera automáticamente en `input/hashes.txt`

```
=== REGISTRO DE AUTENTICIDAD DE DOCUMENTOS ===
Generado: 2026-03-30 11:03:12

Archivo: documento1.txt
Hash SHA256 TXT: 2f04e232dedeca7d150c84e12d194b8500321c429...
Hash SHA256 PDF: 4f790750acda0983c5313eded002b470212468bb...
Hash corto PDF: 4f790750acda0983
PDF: documento1.pdf
```

**⚠️ Guarda este archivo en lugar seguro** - Lo necesitarás para verificar autenticidad

---

## 🎯 Caso de Uso Real

### Escenario: Auditor que verifica documentos de diferentes carpetas

**Día 1 - Procesar auditoría 2024:**
```bash
.\txt2pdf.exe -all -pdf -input ./auditoria_2024

# Se generan PDFs + hashes.txt en ./auditoria_2024/
# Guardo hashes.txt en lugar seguro
```

**Caso 2 - Procesar auditoría 2025 (otra carpeta):**
```bash
.\txt2pdf.exe -all -pdf -input ./auditoria_2025

# Se generan PDFs + hashes.txt en ./auditoria_2025/
# Guardo hashes.txt en lugar seguro
```

**Años después - Verificar cualquiera:**
```bash
.\txt2pdf.exe -file ./auditoria_2024/reporte.pdf -hash
# SHA256: 4f790750acda0983...
# Comparo con hashes.txt guardado
```

---

## 🆘 Ayuda Rápida

### Ver todos los comandos disponibles:
```bash
.\txt2pdf.exe
```

### ¿Necesito crear las carpetas?
✅ **No**, se crean automáticamente. 
- Si usas la carpeta por defecto, se crea `input/`
- Si especificas `-input ./personal`, se crea `personal/`

### ¿Dónde están los PDFs generados?
📁 En la **misma carpeta que los archivos `.txt`**:
- Con `-all -pdf` → en `input/`
- Con `-all -pdf -input ./documentos` → en `./documentos/`

### ¿Puedo editar los PDFs después?
✅ Sí, pero el hash será diferente y se detectará la alteración

### ¿Puedo trabajar con múltiples carpetas?
✅ Sí, especifica `-input` para cada carpeta:
```bash
.\txt2pdf.exe -all -pdf -input ./carpeta1
.\txt2pdf.exe -all -pdf -input ./carpeta2
```

---

## ⚙️ Parámetros Avanzados

| Parámetro | Uso |
|-----------|-----|
| `-file archivo.txt` | Procesar archivo específico |
| `-pdf` | Generar PDF (con auto-detección de orientación) |
| `-all` | Procesar todos los archivos de la carpeta |
| `-hash` | Calcular SHA256 |
| `-input ./carpeta` | Especificar carpeta diferente a `input/` |
| `-portrait` | Fuerza orientación vertical (opcional) |
| `-landscape` | Fuerza orientación horizontal (opcional) |

**Ejemplos:**
```bash
# Procesar todos con auto-detección (defecto)
.\txt2pdf.exe -all -pdf

# Procesar todos de otra carpeta con auto-detección
.\txt2pdf.exe -all -pdf -input ./documentos

# Forzar portrait en todos
.\txt2pdf.exe -all -pdf -portrait -input ./cartas

# Forzar landscape específicamente
.\txt2pdf.exe -all -pdf -landscape -input ./reportes

# Solo leer (sin generar PDF)
.\txt2pdf.exe -file documento.txt

# Solo calcular hash
.\txt2pdf.exe -file documento.pdf -hash
```

---

## ✨ ¿Listo?

1. ✅ Coloca archivos `.txt` en `input/`
2. ✅ Ejecuta: `.\txt2pdf.exe -all -pdf`
3. ✅ Tus PDFs están listos
4. ✅ Guarda `input/hashes.txt`

**¡Ya tienes documentos auténticos!** 🎉
