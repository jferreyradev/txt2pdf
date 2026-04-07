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


**Prepara tu carpeta de trabajo**
1. Crea una carpeta para tus archivos `.txt` (por ejemplo, `documentos/`)
2. Coloca archivos `.txt` en esa carpeta

```bash
.\txt2pdf.exe -all -pdf -input ./documentos
```
Todos los archivos `.txt` en esa carpeta se procesan

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


### Procesar todos los archivos de una carpeta:
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

## Paso 4: Verificar integridad de archivos

### 📝 Verificar archivos TXT

¿Quieres asegurarte de que un archivo TXT no fue alterado?

**Ver hash de UN archivo TXT:**
```bash
.\txt2pdf.exe -file documento.txt -hash
```

**Resultado:**
```
SHA256: 2f04e232dedeca7d150c84e12d194b8500321c429a1f8d1a6c8b5e9f3a4c7d2e
Hash corto: 2f04e232dedeca7d
```
Compara con el que aparece en `hashes.txt` → ✅ Archivo TXT auténtico

**Ver hashes de TODOS los TXT:**
```bash
.\txt2pdf.exe -all -hash
```

### 📄 Verificar archivos PDF

¿Quieres asegurarte de que un PDF no fue alterado?

**Ver hash de UN PDF:**
```bash
.\txt2pdf.exe -file documento1.pdf -hash
```

**Resultado:**
```
SHA256: 4f790750acda0983c5313eded002b470212468bb56608e557fe3ac6af9c16369
Hash corto: 4f790750acda0983
```

Compara el **hash corto** con el que aparece en `hashes.txt` de la misma carpeta → ✅ PDF auténtico

### Ver hashes de TODOS los PDFs:
```bash
.\txt2pdf.exe -all -hash
```

---

## 📋 Archivo de Hashes (`hashes.txt`)

Se genera automáticamente en `hashes.txt` dentro de la carpeta de trabajo

```
=== REGISTRO DE AUTENTICIDAD DE DOCUMENTOS ===
Generado: 2026-03-30 11:03:12

Archivo: documento1.txt
Hash SHA256 TXT: 2f04e232dedeca7d150c84e12d194b8500321c429...
Hash SHA256 PDF: 4f790750acda0983c5313eded002b470212468bb...
Hash corto PDF: 4f790750acda0983
PDF: documento1.pdf
```

**¿Qué validan los hashes?**
- **Hash SHA256 TXT**: Verifica que el archivo de texto **original no fue alterado**
- **Hash SHA256 PDF**: Verifica que el archivo **PDF generado no fue modificado**

**⚠️ Guarda este archivo en lugar seguro** - Lo necesitarás para verificar autenticidad de AMBOS archivos (TXT y PDF)

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
✅ Sí, debes crear la carpeta de trabajo y especificarla con `-input ./tu_carpeta`

### ¿Dónde están los PDFs generados?
📁 En la **misma carpeta que los archivos `.txt`**:
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

1. ✅ Coloca archivos `.txt` en tu carpeta de trabajo
2. ✅ Ejecuta: `.\txt2pdf.exe -all -pdf -input ./tu_carpeta`
3. ✅ Tus PDFs están listos
4. ✅ Guarda `hashes.txt` de esa carpeta

**¡Ya tienes documentos auténticos!** 🎉

---

## 💖 Una Nota Personal

Esta herramienta existe porque alguien **creyó en hacerla bien**. No fue construida por automatización sino por **decisión deliberada**: cada parámetro, cada interacción, cada frase, pensada en que *te sea útil*.

Fue supervisada de cerca, iterada cuando no estaba lista, pulida cuando debía brillar. Con la ayuda de GitHub Copilot, esa visión se convirtió en código que realmente funciona.

No esperamos que pienses "qué buena herramienta". Esperamos que simplemente *funcione* para ti, sin fricciones, sin sorpresas desagradables. Eso es lo que significa hacerlo con 🫶.

*Úsalo, confía en él, y recuerda que fue hecho pensando en ti.* 💙

**Para más detalles técnicos, consulta [README.md](README.md)**
