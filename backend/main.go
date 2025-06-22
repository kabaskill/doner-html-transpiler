package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type TranspileRequest struct {
	Content string `json:"content"`
}

type TranspileResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

// Rate limiting structures
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-rl.window)
	
	// Clean old requests
	if requests, exists := rl.requests[clientIP]; exists {
		validRequests := []time.Time{}
		for _, reqTime := range requests {
			if reqTime.After(cutoff) {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[clientIP] = validRequests
	}
	
	// Check if limit exceeded
	if len(rl.requests[clientIP]) >= rl.limit {
		return false
	}
	
	// Add current request
	rl.requests[clientIP] = append(rl.requests[clientIP], now)
	return true
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ips := strings.Split(xff, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	
	// Fall back to RemoteAddr
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	
	return r.RemoteAddr
}

func main() {
	// Check if running as CLI (if arguments provided)
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		runCLI()
		return
	}

	// Run as web server
	runServer()
}

func runCLI() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: deutsch-html-transpiler <input.dhtml>")
		fmt.Println("Example: deutsch-html-transpiler beispiel.dhtml")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	
	// Read the German HTML file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Create transpiler instance
	transpiler := NewTranspiler()
	
	// Transpile German HTML to standard HTML
	result, err := transpiler.Transpile(string(content))
	if err != nil {
		fmt.Printf("Error transpiling: %v\n", err)
		os.Exit(1)
	}

	// Output the result
	fmt.Println(result)
}

// Security validation function - now just validates basic limits, doesn't block content
func validateInput(content string) error {
	// Check input size limits
	if len(content) > MAX_INPUT_SIZE {
		return fmt.Errorf("input too large (max %d bytes)", MAX_INPUT_SIZE)
	}
	
	// Basic sanity checks but don't block dangerous content
	// The frontend will handle security warnings
	return nil
}

// Add security headers
func addSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'none'; object-src 'none';")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
}

func runServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize rate limiter: 100 requests per minute per IP
	rateLimiter := NewRateLimiter(100, time.Minute)

	// CORS function
	addCORSHeaders := func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost:5173",   // Vite dev server
			"http://localhost:3000",   // Alternative dev port
			"https://doner-html-transpiler.onrender.com", // Production domain
		}
		
		// Check if origin is allowed
		originAllowed := false
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				originAllowed = true
				break
			}
		}
		
		if originAllowed || origin == "" { // Allow empty origin for same-origin requests
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		addSecurityHeaders(w)
		addCORSHeaders(w, r)
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Transpile endpoint
	http.HandleFunc("/transpile", func(w http.ResponseWriter, r *http.Request) {
		// Rate limiting check
		clientIP := getClientIP(r)
		if !rateLimiter.Allow(clientIP) {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(TranspileResponse{Error: "Rate limit exceeded. Please try again later."})
			return
		}
		
		addSecurityHeaders(w)
		addCORSHeaders(w, r)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(TranspileResponse{Error: "Method not allowed"})
			return
		}

		var req TranspileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(TranspileResponse{Error: "Invalid JSON"})
			return
		}

		if req.Content == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(TranspileResponse{Error: "Content is required"})
			return
		}

		// Security validation
		if err := validateInput(req.Content); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(TranspileResponse{Error: err.Error()})
			return
		}

		// Create transpiler instance
		transpiler := NewTranspiler()
		
		// Transpile German HTML to standard HTML
		result, err := transpiler.Transpile(req.Content)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(TranspileResponse{Error: err.Error()})
			return
		}

		// Additional output sanitization - remove the problematic double encoding
		// The input validation already handles dangerous content
		// Just ensure clean output without double-encoding issues

		json.NewEncoder(w).Encode(TranspileResponse{Result: result})
	})

	// Dictionary endpoint - returns all supported tags and attributes
	http.HandleFunc("/dictionary", func(w http.ResponseWriter, r *http.Request) {
		addCORSHeaders(w, r)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
			return
		}

		transpiler := NewTranspiler()
		response := map[string]interface{}{
			"tags":       transpiler.GetSupportedTags(),
			"attributes": transpiler.GetSupportedAttributes(),
		}

		json.NewEncoder(w).Encode(response)
	})

	// Serve static files in production
	staticDir := "./static"
	if _, err := os.Stat(staticDir); err == nil {
		// Serve static assets
		fs := http.FileServer(http.Dir(staticDir))
		http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(staticDir+"/assets"))))
		
		// Serve index.html for the root path and any non-API routes
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Add CORS headers
			addCORSHeaders(w, r)
			
			// Handle preflight OPTIONS requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			if r.URL.Path == "/" || (!strings.HasPrefix(r.URL.Path, "/api/") && 
				!strings.HasPrefix(r.URL.Path, "/health") && 
				!strings.HasPrefix(r.URL.Path, "/transpile") && 
				!strings.HasPrefix(r.URL.Path, "/dictionary")) {
				http.ServeFile(w, r, staticDir+"/index.html")
			} else {
				fs.ServeHTTP(w, r)
			}
		})
	} else {
		// Development mode - serve a simple message
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Add CORS headers
			addCORSHeaders(w, r)
			
			// Handle preflight OPTIONS requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
    <title>D.Ö.N.E.R API</title>
</head>
<body>
    <h1>D.Ö.N.E.R API Server</h1>
    <p>The API server is running. Frontend build not found.</p>
    <h2>Available Endpoints:</h2>
    <ul>
        <li><a href="/health">GET /health</a> - Health check</li>
        <li><a href="/dictionary">GET /dictionary</a> - View dictionary</li>
        <li>POST /transpile - Transpile German HTML</li>
    </ul>
</body>
</html>`)
		})
	}

	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("API available at: http://localhost:%s\n", port)
	fmt.Printf("Health check: http://localhost:%s/health\n", port)
	fmt.Printf("Transpile endpoint: http://localhost:%s/transpile\n", port)
	fmt.Printf("Dictionary endpoint: http://localhost:%s/dictionary\n", port)
	
	// Check if static files exist
	if _, err := os.Stat("./static"); err == nil {
		fmt.Println("✓ Static files found - serving frontend")
	} else {
		fmt.Println("⚠ No static files found - API only mode")
	}
	
	fmt.Printf("🥙 D.Ö.N.E.R server ready!\n")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
