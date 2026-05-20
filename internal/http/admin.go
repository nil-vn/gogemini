package http

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gogemini/internal/config"
	"gogemini/internal/domain"
	"gogemini/internal/repo"
	"gogemini/internal/service"
)

type loginAttemptStore struct {
	mu       sync.Mutex
	fails    map[string]int
	lockedTo map[string]time.Time
}

func newLoginAttemptStore() *loginAttemptStore {
	return &loginAttemptStore{fails: map[string]int{}, lockedTo: map[string]time.Time{}}
}

func (s *loginAttemptStore) allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if until, ok := s.lockedTo[key]; ok {
		if time.Now().Before(until) {
			return false
		}
		delete(s.lockedTo, key)
		delete(s.fails, key)
	}
	return true
}

func (s *loginAttemptStore) fail(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fails[key]++
	if s.fails[key] >= 5 {
		s.lockedTo[key] = time.Now().Add(15 * time.Minute)
	}
}

func (s *loginAttemptStore) success(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.fails, key)
	delete(s.lockedTo, key)
}

func registerAdminRoutes(r *gin.Engine, db *sql.DB, cfg config.Config) {
	repository := repo.AdminRepo{DB: db}
	attempts := newLoginAttemptStore()
	secureCookie := cfg.Environment != "development"
	authSecret := cfg.AuthSecret

	r.POST("/api/auth/login", func(c *gin.Context) {
		var req struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			writeError(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
			return
		}
		loginKey := strings.ToLower(strings.TrimSpace(req.Login))
		if !attempts.allow(loginKey) {
			writeError(c, http.StatusTooManyRequests, "AUTH_LOCKED", "too many failed login attempts")
			return
		}
		u, hash, err := repository.FindUserByUsernameOrEmail(req.Login)
		if err != nil || !service.CheckWerkzeugPasswordHash(hash, req.Password) {
			attempts.fail(loginKey)
			writeError(c, http.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS", "invalid credentials")
			return
		}
		attempts.success(loginKey)
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("session", service.BuildSessionToken(u.ID, authSecret), 3600, "/", "", secureCookie, true)
		c.JSON(http.StatusOK, gin.H{"user": u})
	})
	r.POST("/api/auth/logout", func(c *gin.Context) {
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("session", "", 0, "/", "", secureCookie, true)
		c.Status(http.StatusNoContent)
	})

	admin := r.Group("/api/admin")
	admin.Use(func(c *gin.Context) {
		token, err := c.Cookie("session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorEnvelope{Error: apiError{Code: "AUTH_UNAUTHORIZED", Message: "unauthorized"}})
			return
		}
		if _, ok := service.ValidateSessionToken(token, authSecret); !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorEnvelope{Error: apiError{Code: "AUTH_UNAUTHORIZED", Message: "unauthorized"}})
			return
		}
		c.Next()
	})

	admin.GET("/users", func(c *gin.Context) {
		q, err := parseListQuery(c, []string{"id", "username", "email", "status"})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data, err := repository.ListUsers(q)
		respondList(c, data, err, q)
	})
	admin.GET("/cars", func(c *gin.Context) {
		q, err := parseListQuery(c, []string{"id", "name", "status", "vin"})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data, err := repository.ListCars(q)
		respondList(c, data, err, q)
	})
	admin.GET("/customers", func(c *gin.Context) {
		q, err := parseListQuery(c, []string{"id", "name", "phone", "status"})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data, err := repository.ListCustomers(q)
		respondList(c, data, err, q)
	})
	admin.GET("/transactions", func(c *gin.Context) {
		q, err := parseListQuery(c, []string{"id", "status", "purchase_date", "selling_price"})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data, err := repository.ListTransactions(q)
		respondList(c, data, err, q)
	})
	admin.GET("/dashboard", func(c *gin.Context) { data, err := repository.Dashboard(); respondMap(c, data, err) })

	admin.POST("/users", func(c *gin.Context) {
		var req struct {
			Username     string `json:"username"`
			Role         string `json:"role"`
			Email        string `json:"email"`
			Status       string `json:"status"`
			PasswordHash string `json:"password_hash"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.PasswordHash) == "" {
			c.JSON(400, gin.H{"error": "username and password_hash required"})
			return
		}
		id, err := repository.CreateUser(req.Username, req.Role, req.Email, req.PasswordHash, req.Status)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"id": id})
	})
	admin.PUT("/users/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid id"})
			return
		}
		var req struct {
			Username string `json:"username"`
			Role     string `json:"role"`
			Email    string `json:"email"`
			Status   string `json:"status"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Username) == "" {
			c.JSON(400, gin.H{"error": "username required"})
			return
		}
		if err := repository.UpdateUser(id, req.Username, req.Role, req.Email, req.Status); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(204)
	})
	admin.GET("/users/:id", func(c *gin.Context) { getByID(c, func(id int64) (any, error) { return repository.GetUser(id) }) })
	admin.DELETE("/users/:id", func(c *gin.Context) { deleteByID(c, repository, "users") })

	admin.POST("/cars", func(c *gin.Context) {
		var req domain.Car
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
			c.JSON(400, gin.H{"error": "name required"})
			return
		}
		id, err := repository.CreateCar(req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"id": id})
	})
	admin.GET("/cars/:id", func(c *gin.Context) { getByID(c, func(id int64) (any, error) { return repository.GetCar(id) }) })
	admin.PUT("/cars/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid id"})
			return
		}
		var req domain.Car
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
			c.JSON(400, gin.H{"error": "name required"})
			return
		}
		if err := repository.UpdateCar(id, req); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(204)
	})
	admin.DELETE("/cars/:id", func(c *gin.Context) { deleteByID(c, repository, "car") })

	admin.POST("/customers", func(c *gin.Context) {
		var req domain.Customer
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
			c.JSON(400, gin.H{"error": "name required"})
			return
		}
		id, err := repository.CreateCustomer(req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"id": id})
	})
	admin.GET("/customers/:id", func(c *gin.Context) { getByID(c, func(id int64) (any, error) { return repository.GetCustomer(id) }) })
	admin.PUT("/customers/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid id"})
			return
		}
		var req domain.Customer
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
			c.JSON(400, gin.H{"error": "name required"})
			return
		}
		if err := repository.UpdateCustomer(id, req); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(204)
	})
	admin.DELETE("/customers/:id", func(c *gin.Context) { deleteByID(c, repository, "customer") })

	admin.POST("/transactions", func(c *gin.Context) {
		var req domain.Transaction
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid body"})
			return
		}
		id, err := repository.CreateTransaction(req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"id": id})
	})
	admin.GET("/transactions/:id", func(c *gin.Context) { getByID(c, func(id int64) (any, error) { return repository.GetTransaction(id) }) })
	admin.PUT("/transactions/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid id"})
			return
		}
		var req domain.Transaction
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid body"})
			return
		}
		if err := repository.UpdateTransaction(id, req); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(204)
	})
	admin.DELETE("/transactions/:id", func(c *gin.Context) { deleteByID(c, repository, "transaction") })

	admin.GET("/search", func(c *gin.Context) { data, err := repository.SearchAll(c.Query("q")); respondMap(c, data, err) })

	admin.GET("/system", func(c *gin.Context) {
		m, err := repository.GetSettings()
		respondMap(c, m, err)
	})
	admin.PUT("/system", func(c *gin.Context) {
		token, _ := c.Cookie("session")
		uid, ok := service.ValidateSessionToken(token, authSecret)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		u, err := repository.GetUser(uid)
		if err != nil || !strings.EqualFold(strings.TrimSpace(u.Role), "admin") {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		var req map[string]string
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid body"})
			return
		}
		if err := validateSystemSettingsInput(req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := repository.UpsertSettings(req); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(204)
	})

	admin.POST("/upload/:module", func(c *gin.Context) {
		module := c.Param("module")
		if module != "cars" && module != "customers" {
			c.JSON(400, gin.H{"error": "module must be cars or customers"})
			return
		}
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "file is required"})
			return
		}
		defer file.Close()
		if err := validateUploadFile(file, header); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		fileName := sanitizeUploadFilename(header.Filename)
		rel := filepath.Join(module, fmt.Sprintf("%s_%s", uuid.NewString(), fileName))
		abs := filepath.Join(cfg.UploadDir, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		out, err := os.Create(abs)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer out.Close()
		if _, err := io.Copy(out, file); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"path": filepath.ToSlash(filepath.Join("uploads", rel)), "uploaded_at": time.Now().UTC().Format(time.RFC3339)})
	})
}

var uploadNameSanitizer = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func sanitizeUploadFilename(name string) string {
	base := filepath.Base(strings.TrimSpace(name))
	if base == "." || base == string(filepath.Separator) || base == "" {
		base = "upload"
	}
	base = strings.ReplaceAll(base, "..", "")
	base = uploadNameSanitizer.ReplaceAllString(base, "_")
	base = strings.Trim(base, "._-")
	if base == "" {
		return "upload"
	}
	return base
}

func validateUploadFile(file multipart.File, header *multipart.FileHeader) error {
	if header.Size <= 0 {
		return errors.New("file is empty")
	}
	if header.Size > 5*1024*1024 {
		return errors.New("file too large (max 5MB)")
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExt := map[string]struct{}{".jpg": {}, ".jpeg": {}, ".png": {}, ".webp": {}}
	if _, ok := allowedExt[ext]; !ok {
		return errors.New("unsupported file extension")
	}
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	allowedMIME := map[string]struct{}{"image/jpeg": {}, "image/png": {}, "image/webp": {}}
	ct := strings.ToLower(http.DetectContentType(buf[:n]))
	if _, ok := allowedMIME[ct]; !ok {
		return errors.New("only jpeg/png/webp images are allowed")
	}
	_, err := file.Seek(0, io.SeekStart)
	return err
}

func deleteByID(c *gin.Context, repository repo.AdminRepo, table string) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	if err := repository.DeleteByID(table, id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}

func parseListQuery(c *gin.Context, allowedSort []string) (repo.ListQuery, error) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		return repo.ListQuery{}, errors.New("page must be >= 1")
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		return repo.ListQuery{}, errors.New("page_size must be in range 1..100")
	}
	sort := strings.TrimSpace(c.DefaultQuery("sort", "id"))
	validSort := map[string]struct{}{}
	for _, s := range allowedSort {
		validSort[s] = struct{}{}
	}
	if _, ok := validSort[sort]; !ok {
		return repo.ListQuery{}, fmt.Errorf("sort must be one of: %s", strings.Join(allowedSort, ","))
	}
	order := strings.ToLower(strings.TrimSpace(c.DefaultQuery("order", "desc")))
	if order != "asc" && order != "desc" {
		return repo.ListQuery{}, errors.New("order must be asc or desc")
	}
	return repo.ListQuery{Page: page, PageSize: pageSize, Sort: sort, Order: order, Search: c.Query("q"), Status: strings.TrimSpace(c.Query("status"))}, nil
}
func getByID(c *gin.Context, fn func(id int64) (any, error)) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	data, err := fn(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}
func respondList[T any](c *gin.Context, res repo.ListResult[T], err error, q repo.ListQuery) {
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	sort := q.Sort
	if sort == "" {
		sort = "id"
	}
	order := strings.ToLower(q.Order)
	if order != "asc" {
		order = "desc"
	}
	c.JSON(200, gin.H{"items": res.Items, "total": res.Total, "page": page, "page_size": pageSize, "sort": sort, "order": order})
}
func respondMap(c *gin.Context, data any, err error) {
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}

func validateSystemSettingsInput(input map[string]string) error {
	if len(input) == 0 {
		return errors.New("at least one setting is required")
	}
	allowed := map[string]map[string]struct{}{
		"currency": {"JPY": {}, "USD": {}, "VND": {}},
		"theme":    {"dark": {}, "light": {}},
		"language": {"vi": {}, "ja": {}, "en": {}},
	}
	for k, v := range input {
		values, ok := allowed[k]
		if !ok {
			return fmt.Errorf("unsupported setting key: %s", k)
		}
		normalized := strings.TrimSpace(v)
		if _, ok := values[normalized]; !ok {
			return fmt.Errorf("invalid value for %s", k)
		}
		input[k] = normalized
	}
	return nil
}
