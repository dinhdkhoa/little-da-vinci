package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	const (
		port    = "8080"
		timeout = 30 * time.Second
	)

	logger := log.New(os.Stdout, "[perma-app] ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("Starting Permacomputing Web Application")

	db, err := initDb(logger)
	if err != nil {
		logger.Fatalf("cant connect to db: %v", err)
	}

	migrationsql := getMigrationSql()
	_, err = db.Exec(migrationsql)
	if err != nil {
		logger.Printf("Migration warning/error: %v", err)
	}
	defer db.Close()

	limiter := newRateLimiter()

	mux := http.NewServeMux()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		logger.Println("WARNING: API_KEY environment variable not set. Endpoint /sync-db/{id} will require API_KEY.")
	}

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	mux.HandleFunc("POST /sign-up", rateLimitMiddleware(limiter, signUp(db)))
	mux.HandleFunc("GET /sync-db/{id}", syncDb(db, apiKey))

	handler := loggingMiddleware(logger)(panicRecoveryMiddleware(logger)(mux))

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  timeout,
	}

	logger.Printf("Server starting on port :%s", port)
	logger.Println("Press Ctrl+C to stop")

	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}


func signUp(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.Header().Set("HX-Retarget", "#form-error")
			w.Header().Set("HX-Reswap", "innerHTML")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`<div class="bg-red-500/20 border border-red-500/50 text-red-200 p-4 rounded-xl text-center text-sm font-semibold">
				Dữ liệu biểu mẫu không hợp lệ.
			</div>`))
			return
		}

		title := strings.TrimSpace(r.FormValue("title"))
		parentName := strings.TrimSpace(r.FormValue("parentName"))
		phone := strings.TrimSpace(r.FormValue("phone"))
		studentName := strings.TrimSpace(r.FormValue("studentName"))
		birthYearStr := strings.TrimSpace(r.FormValue("birthYear"))
		gender := strings.TrimSpace(r.FormValue("gender"))
		source := strings.TrimSpace(r.FormValue("source"))

		var errors []string
		if parentName == "" {
			errors = append(errors, "Họ và tên phụ huynh không được để trống.")
		}
		if phone == "" {
			errors = append(errors, "Số điện thoại không được để trống.")
		}
		if studentName == "" {
			errors = append(errors, "Họ và tên học sinh không được để trống.")
		}
		birthYear, err := strconv.Atoi(birthYearStr)
		if err != nil || birthYear < 2005 || birthYear > 2025 {
			errors = append(errors, "Năm sinh học sinh phải từ 2005 đến 2025.")
		}
		if gender == "" {
			errors = append(errors, "Vui lòng chọn giới tính học sinh.")
		}

		if len(errors) > 0 {
			w.Header().Set("HX-Retarget", "#form-error")
			w.Header().Set("HX-Reswap", "innerHTML")
			w.WriteHeader(http.StatusUnprocessableEntity)
			errList := ""
			for _, e := range errors {
				errList += fmt.Sprintf("<li>%s</li>", html.EscapeString(e))
			}
			w.Write([]byte(fmt.Sprintf(`<div class="bg-red-500/20 border border-red-500/50 text-red-200 p-4 rounded-xl text-sm font-semibold">
				<ul class="list-disc list-inside space-y-1">%s</ul>
			</div>`, errList)))
			return
		}

		query := `INSERT INTO registrations (title, parent_name, phone, student_name, birth_year, gender, source)
		          VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err = db.Exec(query, title, parentName, phone, studentName, birthYear, gender, source)
		if err != nil {
			log.Printf("DB error on sign-up: %v", err)
			w.Header().Set("HX-Retarget", "#form-error")
			w.Header().Set("HX-Reswap", "innerHTML")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`<div class="bg-red-500/20 border border-red-500/50 text-red-200 p-4 rounded-xl text-center text-sm font-semibold">
				Có lỗi hệ thống xảy ra. Vui lòng thử lại sau.
			</div>`))
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(fmt.Sprintf(`<div id="registration-form" class="bg-white/10 backdrop-blur-md p-8 md:p-12 rounded-3xl border border-white/20 shadow-2xl space-y-6 text-center text-white">
			<div class="w-16 h-16 bg-primary text-secondary rounded-full flex items-center justify-center mx-auto text-3xl font-black">
				✓
			</div>
			<h3 class="text-2xl md:text-3xl font-black font-display text-primary">
				Đăng Ký Thành Công!
			</h3>
			<p class="text-cyan-100 text-base max-w-md mx-auto leading-relaxed">
				Cảm ơn <strong class="text-white">%s %s</strong> đã đăng ký cho bé <strong class="text-white">%s</strong>! Little Da Vinci sẽ liên hệ qua Zalo/SĐT <strong class="text-white">%s</strong> trong thời gian sớm nhất để xác nhận lớp học.
			</p>
			<div class="pt-4">
				<button type="button"
					hx-get="/"
					hx-select="#registration-form"
					hx-target="#registration-form"
					hx-swap="outerHTML"
					class="bg-primary hover:bg-primary-dark text-secondary font-bold py-3 px-8 rounded-2xl shadow-lg hover:shadow-primary/30 transition-all duration-300 text-sm inline-flex items-center gap-2 cursor-pointer">
					<span>Đăng Ký Cho Học Sinh Khác</span>
					<span class="material-symbols-outlined text-base">add</span>
				</button>
			</div>
		</div>`, html.EscapeString(title), html.EscapeString(parentName), html.EscapeString(studentName), html.EscapeString(phone))))
	}
}

type Registration struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ParentName  string    `json:"parent_name"`
	Phone       string    `json:"phone"`
	StudentName string    `json:"student_name"`
	BirthYear   int       `json:"birth_year"`
	Gender      string    `json:"gender"`
	Source      string    `json:"source"`
	CreatedAt   time.Time `json:"created_at"`
}

func syncDb(db *sql.DB, apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if apiKey == "" || r.Header.Get("X-API-Key") != apiKey {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}

		idStr := r.PathValue("id")
		lastID, err := strconv.Atoi(idStr)
		if err != nil || lastID < 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID parameter"})
			return
		}

		query := `SELECT id, title, parent_name, phone, student_name, birth_year, gender, COALESCE(source, ''), created_at
		          FROM registrations
		          WHERE id > $1
		          ORDER BY id ASC`

		rows, err := db.Query(query, lastID)
		if err != nil {
			log.Printf("DB error on sync-db: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
			return
		}
		defer rows.Close()

		registrations := make([]Registration, 0)
		for rows.Next() {
			var reg Registration
			if err := rows.Scan(
				&reg.ID,
				&reg.Title,
				&reg.ParentName,
				&reg.Phone,
				&reg.StudentName,
				&reg.BirthYear,
				&reg.Gender,
				&reg.Source,
				&reg.CreatedAt,
			); err != nil {
				log.Printf("Row scan error on sync-db: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Error processing data"})
				return
			}
			registrations = append(registrations, reg)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Rows error on sync-db: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Error processing data"})
			return
		}

		json.NewEncoder(w).Encode(registrations)
	}
}

