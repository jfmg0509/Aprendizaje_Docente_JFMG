package http // Paquete HTTP: middlewares aplicados a todas las requests

import (
	"context"   // Contexto para pasar valores entre middlewares y handlers
	"log"       // Logging del sistema
	"math/rand" // Generación de valores aleatorios
	"net/http"  // Tipos HTTP
	"time"      // Manejo de tiempo (duración y seed)
)

// ctxKey es un tipo propio para evitar colisiones de claves en context.Context.
// Buena práctica recomendada por Go.
type ctxKey string

// requestIDKey es la clave usada para guardar el request ID en el contexto.
const requestIDKey ctxKey = "request_id"

// requestIDMiddleware genera un ID único por request.
// - Lo guarda en el contexto
// - Lo agrega como header HTTP (X-Request-Id)
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Genera un ID aleatorio
		rid := randomID()

		// Crea un nuevo contexto con el request ID
		ctx := context.WithValue(r.Context(), requestIDKey, rid)

		// Agrega el ID como header de respuesta
		w.Header().Set("X-Request-Id", rid)

		// Continúa la cadena de middlewares/handlers con el nuevo contexto
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// loggingMiddleware registra información de cada request.
// Se ejecuta antes y después del handler.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Marca el inicio del procesamiento
		start := time.Now()

		// Ejecuta el siguiente middleware o handler
		next.ServeHTTP(w, r)

		// Obtiene el request ID desde el contexto
		rid, _ := r.Context().Value(requestIDKey).(string)

		// Registra método, ruta, request ID y duración
		log.Printf(
			"[%s] %s %s (%s)",
			rid,
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

// methodOverrideMiddleware permite usar PUT y DELETE desde formularios HTML.
// HTML solo soporta GET y POST, así que se usa un campo hidden "_method".
//
// Ejemplo en HTML:
// <input type="hidden" name="_method" value="PUT">
func methodOverrideMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Solo aplica si el método original es POST
		if r.Method == http.MethodPost {

			// ParseForm lee los valores del formulario
			if err := r.ParseForm(); err == nil {

				// Si existe el campo "_method", se sobreescribe el método HTTP
				if m := r.FormValue("_method"); m != "" {
					r.Method = m
				}
			}
		}

		// Continúa con la request (ya modificada si aplicaba)
		next.ServeHTTP(w, r)
	})
}

// randomID genera un string aleatorio para identificar requests.
// Usa letras minúsculas y números.
func randomID() string {

	// Inicializa la semilla del generador aleatorio con el tiempo actual
	rand.Seed(time.Now().UnixNano())

	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

	// Crea un slice de bytes de longitud fija
	b := make([]byte, 10)

	// Llena el slice con caracteres aleatorios
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	// Convierte []byte a string y lo retorna
	return string(b)
}
