# Test Pulpoline API 🚀

Un servicio web en Go que procesa texto usando inteligencia artificial. Permite enviar preguntas o textos y recibir respuestas generadas por IA, todo manejado de forma concurrente y eficiente.

## ¿Qué hace este proyecto?

Imagínate que tienes un asistente virtual que responde preguntas. Este servicio hace exactamente eso: recibes texto (como "¿Qué es Java?"), lo envías a una API de IA, y obtienes una respuesta inteligente de vuelta. Todo esto funcionando de forma concurrente, es decir, puede manejar muchas peticiones al mismo tiempo sin problemas.

## Características principales

- ✅ **API REST simple**: Solo envías un POST con tu texto y recibes la respuesta
- ✅ **Soporte múltiple de IA**: Funciona con OpenAI o Groq API
- ✅ **Súper rápido**: Usa goroutines de Go para manejar múltiples peticiones en paralelo
- ✅ **Confiable**: Manejo de errores, timeouts, y cierre seguro del servidor
- ✅ **Fácil de usar**: Health check para verificar que todo funciona
- ✅ **Opción gratuita**: Groq API es completamente gratuita para pruebas

## ¿Qué necesito para empezar?

**Requisitos básicos:**
- Go 1.21 o superior instalado
- Una API key de Groq (gratuita) o OpenAI

**Opciones de API:**
- **Groq API**: Gratuita, solo necesitas crear una cuenta en https://console.groq.com/
- **OpenAI**: Requiere créditos (aunque suelen dar algunos gratis al registrarte)

## Primeros pasos

### 1. Clonar el proyecto

```bash
git clone <url-del-repositorio>
cd test-pulpoline-api
```

### 2. Instalar dependencias

```bash
go mod download
```

Esto descargará todas las librerías que necesita el proyecto. Solo toma unos segundos.

### 3. Configurar el entorno

Tienes dos opciones, elige la que más te convenga:

#### Opción A: Groq API (Gratuita - Recomendada) ⭐

Perfecta para empezar, es completamente gratuita:

1. Ve a https://console.groq.com/ y crea una cuenta (es gratis, no requiere tarjeta)
2. Obtén tu API key desde el dashboard
3. Configura tu `.env`:

```env
AI_PROVIDER=groq
GROQ_API_KEY=tu_groq_api_key_aqui
SERVER_ADDR=:8080
```

#### Opción B: OpenAI (Si tienes créditos)

Si prefieres usar OpenAI (recuerda que tiene costo, aunque dan créditos iniciales):

```env
AI_PROVIDER=openai
OPENAI_API_KEY=tu_openai_key_aqui
SERVER_ADDR=:8080
```

**Tip:** Para usar variables de entorno directamente sin archivo `.env`:

```bash
# Windows PowerShell
$env:AI_PROVIDER="groq"
$env:GROQ_API_KEY="tu_key_aqui"
$env:SERVER_ADDR=":8080"

# Linux/Mac
export AI_PROVIDER=groq
export GROQ_API_KEY=tu_key_aqui
export SERVER_ADDR=:8080
```

## Ejecutar el servidor

Tienes dos formas de hacerlo:

### Opción 1: Ejecutar directamente (Recomendado para desarrollo)

```bash
go run ./cmd/api
```

### Opción 2: Compilar primero

```bash
go build -o test-pulpoline-api ./cmd/api
./test-pulpoline-api
```

Una vez ejecutado, verás algo como:
```
Archivo .env cargado correctamente desde: .env
Configuración cargada - Provider: groq, ServerAddr: :8080
Usando Groq API (gratuita)
Servidor iniciado en :8080
```

¡Listo! Tu servidor está corriendo en `http://localhost:8080` 🎉

## Cómo usar la API

### Verificar que funciona (Health Check)

Primero, asegúrate de que el servidor está vivo:

```bash
curl http://localhost:8080/health
```

Deberías ver:
```json
{
  "status": "healthy",
  "service": "test-pulpoline-api"
}
```

### Procesar un texto

Ahora sí, envía tu pregunta o texto a procesar:

```bash
curl -X POST http://localhost:8080/api/process \
  -H "Content-Type: application/json" \
  -d '{"text": "Explica qué es la programación concurrente en Go"}'
```

**Respuesta exitosa:**
```json
{
  "id": "78ce2d49-3c73-4123-b65d-00a503f8113a",
  "text": "Explica qué es la programación concurrente en Go",
  "response": "La programación concurrente en Go...",
  "status": "success"
}
```

**Si algo sale mal:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "error": "error al procesar con IA: API key no está configurada",
  "status": "error"
}
```

## Cómo está organizado el código

El proyecto sigue una estructura profesional tipo microservicio. Cada cosa tiene su lugar:

```
test-pulpoline-api/
├── cmd/api/              # 🚪 La puerta de entrada - aquí empieza todo
├── internal/
│   ├── config/          # ⚙️ Configuración y variables de entorno
│   ├── handler/         # 📨 Maneja las peticiones HTTP que llegan
│   ├── service/         # 🧠 La lógica de negocio - aquí se procesa todo
│   ├── client/ai/       # 🤖 Cliente para hablar con las APIs de IA
│   └── queue/           # 📬 Sistema de colas para manejar concurrencia
└── pkg/errors/          # ❌ Errores personalizados
```

**En palabras simples:**
- `cmd/api` = El main, donde todo empieza
- `handler` = Recibe las peticiones del usuario
- `service` = Decide qué hacer con esas peticiones
- `client/ai` = Habla con OpenAI, Groq, etc.
- `queue` = Organiza las peticiones para no sobrecargarse

## Cómo funciona internamente

### El flujo completo

1. **Llega una petición** → El handler la recibe y valida
2. **Se genera un ID único** → Para rastrear cada petición
3. **Se encola o procesa directamente** → Depende de cuántas peticiones hay
4. **Se envía a la IA** → Tu texto va a OpenAI o Groq
5. **Se recibe la respuesta** → La IA genera una respuesta
6. **Se devuelve al cliente** → Tu recibe la respuesta en JSON

Todo esto pasa de forma **concurrente**, así que si llegan 10 peticiones al mismo tiempo, todas se procesan en paralelo sin bloquearse.

### Los "ingredientes" de la concurrencia

El sistema usa las herramientas de Go:
- **Goroutines**: Como trabajadores que procesan tareas en paralelo
- **Canales**: Como tubos de comunicación entre los trabajadores
- **Context**: Para controlar timeouts y cancelaciones
- **WaitGroups**: Para coordinar que todos terminen correctamente

## Probar la aplicación

### Prueba básica

```bash
# Verificar que está vivo
curl http://localhost:8080/health

# Hacer una pregunta
curl -X POST http://localhost:8080/api/process \
  -H "Content-Type: application/json" \
  -d '{"text": "Hola, ¿cómo estás?"}'
```

### Probar concurrencia

¿Quieres ver cómo maneja múltiples peticiones a la vez? Prueba esto:

```bash
# En Linux/Mac
for i in {1..5}; do
  curl -X POST http://localhost:8080/api/process \
    -H "Content-Type: application/json" \
    -d "{\"text\": \"Petición número $i\"}" &
done
wait
```

Esto enviará 5 peticiones al mismo tiempo. ¡Verás cómo todas se procesan en paralelo!

## Configuración avanzada

Si quieres ajustar cosas, puedes modificar estos valores en el código:

- **Tamaño de la cola**: En `main.go`, cambia `NewRequestQueue(10)` (10 = 10 peticiones en cola)
- **Número de workers**: En `queue.go`, cambia `workers: 5` (5 = 5 procesadores paralelos)
- **Timeout de peticiones**: En `handler.go`, cambia `30*time.Second` (30 segundos máximo)
- **Modelo de IA**: En los clientes de IA, cambia el modelo según lo que necesites

## Solución de problemas comunes

### "No me toma las variables de entorno"

- Verifica que el archivo `.env` esté en la raíz del proyecto
- Asegúrate de que no tenga espacios extra o caracteres raros
- Si usas Cursor/VSCode, reinicia el servidor después de cambiar el `.env`

### "La cola está llena"

El servidor está recibiendo muchas peticiones muy rápido. Puedes:
- Aumentar el tamaño de la cola en `main.go`
- Esperar un momento y volver a intentar
- Revisar si hay algún proceso que esté enviando demasiadas peticiones

### "El servidor no responde"

Verifica:
- ¿Está corriendo? Busca el mensaje "Servidor iniciado en :8080"
- ¿El puerto está libre? Asegúrate de que no haya otro proceso usando el puerto 8080
- ¿Hay errores en la consola? Revisa los logs que aparecen al iniciar

### "Error: API key no está configurada"

- Asegúrate de tener la API key correcta en tu `.env`:
  - Para Groq: `GROQ_API_KEY=tu_key_aqui`
  - Para OpenAI: `OPENAI_API_KEY=tu_key_aqui`
- Verifica que `AI_PROVIDER` esté configurado correctamente (`groq` o `openai`)
- Verifica que el `.env` esté siendo cargado (deberías ver un log al iniciar)

## Lo que hace especial este código

- ✅ **Bien organizado**: Cada cosa en su lugar, fácil de encontrar y modificar
- ✅ **Manejo de errores**: Si algo falla, sabrás exactamente qué pasó
- ✅ **Preparado para producción**: Cierre seguro, timeouts, validaciones
- ✅ **Fácil de extender**: Agregar nuevas funcionalidades es sencillo
- ✅ **Documentado**: El código explica qué hace cada cosa
- ✅ **Seguro**: Usa canales y mutexes correctamente para evitar condiciones de carrera

## Ideas para mejorar (si quieres seguir trabajando en esto)

- [ ] Agregar autenticación (API keys para proteger los endpoints)
- [ ] Implementar rate limiting (evitar abusos)
- [ ] Agregar métricas (ver cuántas peticiones se procesan, tiempos, etc.)
- [ ] Tests automatizados (para asegurar que todo funciona)
- [ ] Cache de respuestas (para no repetir las mismas peticiones)
- [ ] Documentación OpenAPI/Swagger (para que otros sepan cómo usar tu API)
- [ ] Soporte para más modelos de IA

## Licencia

Este proyecto es parte de una prueba técnica para Pulpoline.

## Autor

**Gabriel Alejandro Pina**  
Desarrollador FullStack
