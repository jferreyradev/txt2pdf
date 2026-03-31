# 🌐 API REST - QUICKSTART

Guía rápida para integrar txt2pdf en tus aplicaciones via API REST.

---

## Paso 1: Iniciar el servidor

```bash
txt2pdf.exe -api -port 8080
```

**Salida esperada:**
```
🚀 txt2pdf API REST iniciado
🌐 Servidor escuchando en http://localhost:8080
📚 Documentación: http://localhost:8080/help
🔗 Status: http://localhost:8080/status
```

✅ El servidor está listo. Accede a http://localhost:8080/help en tu navegador.

---

## Paso 2: Probar con cURL

### Convertir un archivo TXT a PDF

```bash
curl -F "file=@BOLETAS.txt" http://localhost:8080/convert \
  --output BOLETAS.pdf
```

**Respuesta:** Se descarga `BOLETAS.pdf` con headers:
- `X-PDF-Hash`: Hash SHA256 completo
- `X-PDF-Hash-Short`: Primeros 16 caracteres

### Convertir múltiples archivos a PDF (ZIP)

```bash
curl -F "file=@BOLETAS.txt" -F "file=@LIBRAMIENTOS.txt" -F "file=@PLANILLAS.txt" \
  http://localhost:8080/convert --output reportes.zip
```

**Respuesta:** Se descarga `reportes.zip` con todos los PDFs generados

### Con orientación personalizada

```bash
curl -F "file=@PLANILLAS.txt" -F "file=@RESUMEN.txt" \
  -F "orientation=landscape" http://localhost:8080/convert --output apaisado.zip
```

**Opciones disponibles:** `auto`, `portrait`, `landscape`

### Calcular hash de archivo

```bash
curl -F "file=@documento.pdf" http://localhost:8080/hash
```

**Respuesta JSON:**
```json
{
  "filename": "documento.pdf",
  "sha256": "abc123def456...",
  "short_hash": "abc123def456abcd",
  "size_bytes": 17537
}
```

### Ver estado del servidor

```bash
curl http://localhost:8080/status | jq
```

---

## Paso 3: Integraciónes

### JavaScript + HTML (Frontend)

```html
<!DOCTYPE html>
<html>
<head>
    <title>txt2pdf API</title>
</head>
<body>
    <h1>txt2pdf Batch Converter</h1>
    
    <input type="file" id="fileInput" accept=".txt" multiple />
    
    <select id="orientation">
        <option value="auto">Auto-detect</option>
        <option value="portrait">Portrait</option>
        <option value="landscape">Landscape</option>
    </select>
    
    <button onclick="convertToPDF()">Convertir</button>
    <div id="status"></div>

    <script>
        async function convertToPDF() {
            const files = document.getElementById('fileInput').files;
            const orientation = document.getElementById('orientation').value;
            const status = document.getElementById('status');
            
            if (files.length === 0) {
                status.innerHTML = '❌ Selecciona uno o más archivos';
                return;
            }

            const formData = new FormData();
            for (let file of files) {
                formData.append('file', file);
            }
            formData.append('orientation', orientation);

            try {
                status.innerHTML = `⏳ Procesando ${files.length} archivo(s)...`;
                
                const response = await fetch('http://localhost:8080/convert', {
                    method: 'POST',
                    body: formData
                });

                if (!response.ok) {
                    throw new Error(`Error: ${response.statusText}`);
                }

                // Descargar resultado
                const blob = await response.blob();
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = files.length === 1 
                    ? files[0].name.replace('.txt', '.pdf')
                    : 'documentos.zip';
                a.click();
                window.URL.revokeObjectURL(url);

                status.innerHTML = `
                    ✅ ${files.length === 1 ? 'PDF generado' : 'ZIP generado'} correctamente<br>
                    📦 Descargando...
                `;
            } catch (error) {
                status.innerHTML = `❌ ${error.message}`;
            }
        }
    </script>
</body>
</html>
```

---

## Paso 4: Documentación de endpoints
}

# Uso
convert_txt_to_pdf "documento.txt" "documento.pdf" "portrait"
get_hash "documento.pdf"
```

---

## Casos de Uso

### 1️⃣ Conversión Automática en Lotes

```bash
#!/bin/bash

# Procesar todos los TXT en una carpeta
for file in *.txt; do
    curl -F "file=@$file" http://localhost:8080/convert \
         --output "${file%.txt}.pdf" \
         --progress-bar
    echo "✓ Procesado: $file"
done
```

### 2️⃣ Sistema de Auditoría

```python
import requests
import json
from datetime import datetime

audit_log = []

for txt_file in ['reporte1.txt', 'reporte2.txt']:
    # Convertir
    response = requests.post(
        'http://localhost:8080/convert',
        files={'file': open(txt_file, 'rb')}
    )
    
    pdf_hash = response.headers.get('X-PDF-Hash-Short')
    
    # Registrar en auditoría
    audit_log.append({
        'timestamp': datetime.now().isoformat(),
        'file': txt_file,
        'pdf_hash': pdf_hash,
        'status': 'OK'
    })

# Guardar auditoría
with open('audit.json', 'w') as f:
    json.dump(audit_log, f, indent=2)
```

### 3️⃣ Verificación de Integridad

```javascript
async function verificarIntegridad(pdfFile, expectedHash) {
    const formData = new FormData();
    formData.append('file', pdfFile);
    
    const response = await fetch('http://localhost:8080/hash', {
        method: 'POST',
        body: formData
    });
    
    const data = await response.json();
    const coincide = data.sha256 === expectedHash;
    
    return {
        verificado: coincide,
        hash_esperado: expectedHash,
        hash_actual: data.sha256,
        estado: coincide ? '✅ Intacto' : '❌ Alterado'
    };
}
```

---

## ⚙️ Configuración Avanzada

### Cambiar Puerto

```bash
txt2pdf.exe -api -port 3000
```

### Detrás de Proxy/Load Balancer

```bash
# Habilitar X-Forwarded-* headers
export GIN_MODE=release
txt2pdf.exe -api -port 8080
```

### HTTPS (con reverse proxy)

```nginx
server {
    listen 443 ssl;
    server_name api.txt2pdf.local;
    
    ssl_certificate /etc/ssl/certs/cert.pem;
    ssl_certificate_key /etc/ssl/private/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header X-Forwarded-For $remote_addr;
    }
}
```

---

## 🆘 Troubleshooting

### "Connection refused"
```bash
# Verificar que el servidor está corriendo
curl http://localhost:8080/status

# Si no funciona, inicia el servidor:
txt2pdf.exe -api -port 8080
```

### "File not found"
```bash
# Asegúrate de que el archivo existe y está accesible
# Los archivos deben ser TXT o cualquier formato

# Para calcular hash, puede ser cualquier archivo
curl -F "file=@archivo.pdf" http://localhost:8080/hash
```

### "Invalid file type"
```bash
# El endpoint /convert solo acepta contenido de texto
# Para otros formatos, usa /hash
```

### "Timeout"
```bash
# Aumentar timeout en cliente
curl --max-time 30 -F "file=@archivo_grande.txt" \
     http://localhost:8080/convert --output salida.pdf
```

---

## 📚 Referencias Rápidas

| Endpoint | Método | Parámetros | Respuesta |
|----------|--------|-----------|----------|
| `/convert` | POST | file, orientation | PDF binary |
| `/hash` | POST | file | JSON con hashes |
| `/status` | GET | — | JSON con estado |
| `/help` | GET | — | JSON con documentación |

**Headers útiles en respuesta de /convert:**
- `X-PDF-Hash` - SHA256 completo
- `X-PDF-Hash-Short` - Primeros 16 caracteres
- `Content-Type: application/pdf`
- `Content-Disposition: attachment`

---

## 💙 Listo para usar

Ya tienes todo lo necesario para integrar txt2pdf en tus aplicaciones.

**Siguiente paso:** Consulta [README.md](README.md) para detalles técnicos avanzados.
