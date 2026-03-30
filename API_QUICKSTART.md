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

### Convertir TXT a PDF

```bash
curl -F "file=@documento.txt" http://localhost:8080/convert \
  --output documento.pdf
```

**Resultado:** `documento.pdf` se descarga automáticamente con headers:
- `X-PDF-Hash`: Hash SHA256 completo
- `X-PDF-Hash-Short`: Primeros 16 caracteres

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
    <h1>Convertir TXT a PDF</h1>
    
    <input type="file" id="fileInput" accept=".txt" />
    
    <select id="orientation">
        <option value="auto">Auto-detect</option>
        <option value="portrait">Portrait</option>
        <option value="landscape">Landscape</option>
    </select>
    
    <button onclick="convertToPDF()">Convertir</button>
    <div id="status"></div>

    <script>
        async function convertToPDF() {
            const file = document.getElementById('fileInput').files[0];
            const orientation = document.getElementById('orientation').value;
            const status = document.getElementById('status');
            
            if (!file) {
                status.innerHTML = '❌ Selecciona un archivo';
                return;
            }

            const formData = new FormData();
            formData.append('file', file);
            formData.append('orientation', orientation);

            try {
                status.innerHTML = '⏳ Procesando...';
                
                const response = await fetch('http://localhost:8080/convert', {
                    method: 'POST',
                    body: formData
                });

                if (!response.ok) {
                    throw new Error(`Error: ${response.statusText}`);
                }

                // Obtener headers con los hashes
                const fullHash = response.headers.get('X-PDF-Hash');
                const shortHash = response.headers.get('X-PDF-Hash-Short');

                // Descargar PDF
                const blob = await response.blob();
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = file.name.replace('.txt', '.pdf');
                a.click();
                window.URL.revokeObjectURL(url);

                status.innerHTML = `
                    ✅ PDF generado correctamente<br>
                    📄 Hash: <code>${shortHash}</code><br>
                    🔒 Verifica con el hash completo si lo necesitas
                `;
            } catch (error) {
                status.innerHTML = `❌ ${error.message}`;
            }
        }
    </script>
</body>
</html>
```

### Python

```python
import requests

# Convertir TXT a PDF
files = {'file': open('documento.txt', 'rb')}
data = {'orientation': 'auto'}

response = requests.post(
    'http://localhost:8080/convert',
    files=files,
    data=data
)

if response.status_code == 200:
    # Guardar PDF
    with open('documento.pdf', 'wb') as f:
        f.write(response.content)
    
    # Obtener hashes
    full_hash = response.headers.get('X-PDF-Hash')
    short_hash = response.headers.get('X-PDF-Hash-Short')
    
    print(f"✅ PDF generado")
    print(f"Hash corto: {short_hash}")
else:
    print(f"❌ Error: {response.status_code}")

# Calcular hash
files = {'file': open('documento.pdf', 'rb')}
response = requests.post(
    'http://localhost:8080/hash',
    files=files
)

result = response.json()
print(f"SHA256: {result['sha256']}")
print(f"Tamaño: {result['size_bytes']} bytes")
```

### Node.js + Express

```javascript
const express = require('express');
const formData = require('form-data');
const axios = require('axios');
const fs = require('fs');

const app = express();

app.post('/convertir', async (req, res) => {
    try {
        // Preparar archivo
        const form = new formData();
        form.append('file', fs.createReadStream('documento.txt'));
        form.append('orientation', 'auto');

        // Enviar a txt2pdf API
        const response = await axios.post(
            'http://localhost:8080/convert',
            form,
            { headers: form.getHeaders() }
        );

        // Guardar PDF
        fs.writeFileSync('documento.pdf', response.data);

        // Obtener hashes
        const fullHash = response.headers['x-pdf-hash'];
        const shortHash = response.headers['x-pdf-hash-short'];

        res.json({
            success: true,
            file: 'documento.pdf',
            hash_short: shortHash,
            hash_full: fullHash
        });
    } catch (error) {
        res.json({ success: false, error: error.message });
    }
});

app.listen(3000);
```

### cURL + Bash

```bash
#!/bin/bash

# Función para convertir TXT a PDF
convert_txt_to_pdf() {
    local input_file=$1
    local output_file=$2
    local orientation=${3:-auto}

    echo "⏳ Convirtiendo $input_file..."

    curl -F "file=@$input_file" \
         -F "orientation=$orientation" \
         http://localhost:8080/convert \
         --output "$output_file" \
         -H "User-Agent: bash-script"

    if [ $? -eq 0 ]; then
        echo "✅ $output_file generado"
    else
        echo "❌ Error al procesar"
        return 1
    fi
}

# Función para obtener hash
get_hash() {
    local file=$1
    curl -F "file=@$file" http://localhost:8080/hash | jq .
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
