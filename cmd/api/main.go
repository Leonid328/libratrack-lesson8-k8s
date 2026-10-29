package main

import (
	"log"
	"net/http"
	"os"

	"libratrack/internal/handler"
	webhandler "libratrack/internal/handler/web"
	"libratrack/internal/middleware"
	"libratrack/internal/repository"
	"libratrack/web"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return errEnv("DATABASE_URL is required")
	}
	migrPath := os.Getenv("MIGRATIONS_PATH")
	if migrPath == "" {
		migrPath = "migrations"
	}
	if err := repository.RunMigrations(dbURL, migrPath); err != nil {
		log.Printf("migrations: %v", err)
	}
	db, err := repository.Open(dbURL)
	if err != nil {
		return err
	}

	books := repository.NewBookRepository(db)
	authors := repository.NewAuthorRepository(db)
	cats := repository.NewCategoryRepository(db)
	copies := repository.NewCopyRepository(db)
	users := repository.NewUserRepository(db)

	bh := handler.NewBookHandler(books)
	ah := handler.NewAuthorHandler(authors)
	ch := handler.NewCategoryHandler(cats)
	cph := handler.NewCopyHandler(copies)
	authH := handler.NewAuthHandler(users)
	uh := handler.NewUserHandler(users)

	tmpl, err := web.ParseTemplates()
	if err != nil {
		return err
	}
	wh := webhandler.NewCatalogHandler(books, tmpl)

	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logging(), cors())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/swagger/index.html", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, swaggerHTML)
	})
	r.GET("/static/*filepath", gin.WrapH(http.StripPrefix("/static/", web.StaticHandler())))
	r.GET("/catalog", wh.Catalog)
	r.GET("/web/books", wh.BookRows)

	public := r.Group("/api/v1")
	{
		public.POST("/auth/register", authH.Register)
		public.POST("/auth/login", authH.Login)
		public.GET("/books", bh.List)
		public.GET("/books/:id", bh.Get)
		public.POST("/books", bh.Create)
		public.PUT("/books/:id", bh.Update)
		public.DELETE("/books/:id", bh.Delete)
		public.POST("/books/:id/copies", cph.Create)
		public.GET("/books/:id/copies", cph.List)
		public.GET("/books/:id/availability", cph.Availability)
		public.POST("/authors", ah.Create)
		public.GET("/authors", ah.List)
		public.GET("/authors/:id", ah.Get)
		public.POST("/categories", ch.Create)
		public.GET("/categories", ch.List)
		public.GET("/categories/:id", ch.Get)
	}

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me/loans", uh.MyLoans)
	}

	staff := r.Group("/web")
	staff.Use(middleware.AuthMiddleware(), middleware.RequireRole("librarian", "admin"))
	{
		staff.DELETE("/books/:id", wh.Delete)
	}

	admin := r.Group("/api/v1")
	admin.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		admin.GET("/users", uh.List)
		admin.PUT("/users/:id/role", uh.UpdateRole)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return r.Run(":" + port)
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

type envError string

func (e envError) Error() string { return string(e) }

func errEnv(s string) error { return envError(s) }

const swaggerHTML = `<!DOCTYPE html><html><head><title>LibraTrack API</title></head>
<body><h1>LibraTrack OpenAPI</h1><p>CRUD books, authors, categories, copies.</p></body></html>
`
