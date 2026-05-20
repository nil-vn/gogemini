package http

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gogemini/internal/domain"
	"gogemini/internal/repo"
	"gogemini/internal/service"
)

func registerAdminRoutes(r *gin.Engine, db *sql.DB) {
	repository := repo.AdminRepo{DB: db}

	r.POST("/api/auth/login", func(c *gin.Context) {
		var req struct{ Login, Password string }
		if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return }
		u, hash, err := repository.FindUserByUsernameOrEmail(req.Login)
		if err != nil || !service.CheckWerkzeugPasswordHash(hash, req.Password) { c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"}); return }
		c.SetCookie("session", service.BuildSessionToken(u.ID, os.Getenv("AUTH_SECRET")), 86400, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"user": u})
	})

	admin := r.Group("/api/admin")
	admin.Use(func(c *gin.Context) { if _, err := c.Cookie("session"); err != nil { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}); return }; c.Next() })

	admin.GET("/users", func(c *gin.Context) { data, err := repository.ListUsers(); respondList(c, data, err) })
	admin.GET("/cars", func(c *gin.Context) { data, err := repository.ListCars(); respondList(c, data, err) })
	admin.GET("/customers", func(c *gin.Context) { data, err := repository.ListCustomers(); respondList(c, data, err) })
	admin.GET("/transactions", func(c *gin.Context) { data, err := repository.ListTransactions(); respondList(c, data, err) })
	admin.GET("/dashboard", func(c *gin.Context) { data, err := repository.Dashboard(); respondMap(c, data, err) })

	admin.POST("/users", func(c *gin.Context) {
		var req struct{ Username, Role, Email, Status, PasswordHash string }
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.PasswordHash) == "" { c.JSON(400, gin.H{"error": "username and password_hash required"}); return }
		id, err := repository.CreateUser(req.Username, req.Role, req.Email, req.PasswordHash, req.Status)
		if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
		c.JSON(201, gin.H{"id": id})
	})
	admin.PUT("/users/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64); if err != nil { c.JSON(400, gin.H{"error": "invalid id"}); return }
		var req struct{ Username, Role, Email, Status string }
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Username) == "" { c.JSON(400, gin.H{"error": "username required"}); return }
		if err := repository.UpdateUser(id, req.Username, req.Role, req.Email, req.Status); err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
		c.Status(204)
	})
	admin.DELETE("/users/:id", func(c *gin.Context) { deleteByID(c, repository, "users") })

	admin.POST("/cars", func(c *gin.Context) {
		var req domain.Car
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" { c.JSON(400, gin.H{"error": "name required"}); return }
		id, err := repository.CreateCar(req); if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
		c.JSON(201, gin.H{"id": id})
	})
	admin.DELETE("/cars/:id", func(c *gin.Context) { deleteByID(c, repository, "car") })

	admin.POST("/customers", func(c *gin.Context) {
		var req domain.Customer
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" { c.JSON(400, gin.H{"error": "name required"}); return }
		id, err := repository.CreateCustomer(req); if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
		c.JSON(201, gin.H{"id": id})
	})
	admin.DELETE("/customers/:id", func(c *gin.Context) { deleteByID(c, repository, "customer") })

	admin.POST("/transactions", func(c *gin.Context) {
		var req domain.Transaction
		if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"error": "invalid body"}); return }
		id, err := repository.CreateTransaction(req); if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
		c.JSON(201, gin.H{"id": id})
	})
	admin.DELETE("/transactions/:id", func(c *gin.Context) { deleteByID(c, repository, "transaction") })

	admin.GET("/search", func(c *gin.Context) { data, err := repository.SearchAll(c.Query("q")); respondMap(c, data, err) })

	admin.GET("/system", func(c *gin.Context) {
		m, err := repository.GetSettings()
		respondMap(c, m, err)
	})
	admin.PUT("/system", func(c *gin.Context) {
		var req map[string]string
		if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"error": "invalid body"}); return }
		if err := repository.UpsertSettings(req); err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
		c.Status(204)
	})

	admin.POST("/upload/:module", func(c *gin.Context) {
		module := c.Param("module")
		if module != "cars" && module != "customers" { c.JSON(400, gin.H{"error": "module must be cars or customers"}); return }
		file, header, err := c.Request.FormFile("file")
		if err != nil { c.JSON(400, gin.H{"error": "file is required"}); return }
		defer file.Close()
		if header.Size > 5*1024*1024 { c.JSON(400, gin.H{"error": "file too large (max 5MB)"}); return }
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		if ct := http.DetectContentType(buf[:n]); !strings.HasPrefix(ct, "image/") { c.JSON(400, gin.H{"error": "only image files are allowed"}); return }
		if _, err := file.Seek(0, io.SeekStart); err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
		rel := filepath.Join("uploads", module, fmt.Sprintf("%s_%s", uuid.NewString(), filepath.Base(header.Filename)))
		root := os.Getenv("UPLOAD_ROOT")
		if root == "" { root = "static" }
		abs := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
		out, err := os.Create(abs); if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
		defer out.Close()
		if _, err := io.Copy(out, file); err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
		c.JSON(201, gin.H{"path": filepath.ToSlash(rel), "uploaded_at": time.Now().UTC().Format(time.RFC3339)})
	})
}

func deleteByID(c *gin.Context, repository repo.AdminRepo, table string) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil { c.JSON(400, gin.H{"error": "invalid id"}); return }
	if err := repository.DeleteByID(table, id); err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	c.Status(204)
}

func respondList(c *gin.Context, data any, err error) { if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }; c.JSON(200, gin.H{"items": data}) }
func respondMap(c *gin.Context, data any, err error) { if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }; c.JSON(200, data) }
