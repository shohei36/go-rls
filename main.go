package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := NewTenantDB(
		TenantDBConfig{
			Host:       getEnv("DB_HOST", "localhost"),
			Port:       getEnv("DB_PORT", "5432"),
			User:       getEnv("DB_USER", "tenant"),
			Password:   getEnv("DB_PASSWORD", "password"),
			DBName:     getEnv("DB_DB_NAME", "mydb"),
			SchemaName: getEnv("DB_SCHEMA_NAME", "myschema"),
		},
	)
	if err != nil {
		panic(fmt.Sprintf("failed to init db: %s", err.Error()))
	}

	repo := NewUserRepository(db)
	tx := NewTransaction(db)
	usecase := NewUserUsecase(tx, repo)

	r := gin.Default()
	r.GET("/users/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		tenantID := ctx.Request.Header.Get("tenant-id")
		user, err := usecase.FetchByID(ctx, tenantID, id)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error_reason": err.Error(),
			})
			return
		}
		type DTO struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Gender string `json:"gender"`
			Age    int    `json:"age"`
		}
		dto := DTO{
			ID:     id,
			Name:   user.Name,
			Gender: user.Gender,
			Age:    user.Age,
		}
		ctx.JSON(200, gin.H{"user": dto})
	})
	r.POST("/users/:id", func(ctx *gin.Context) {
		tenantID := ctx.Request.Header.Get("tenant-id")
		type DTO struct {
			Name   string `json:"name" binding:"required"`
			Gender string `json:"gender" binding:"required"`
			Age    int    `json:"age" binding:"required"`
		}
		var dto DTO
		if err := ctx.BindJSON(&dto); err != nil {
			ctx.JSON(
				400,
				gin.H{"error_reason": err.Error()},
			)
			return
		}
		m := User{
			ID:     ctx.Param("id"),
			Name:   dto.Name,
			Gender: dto.Gender,
			Age:    dto.Age,
		}
		if err := usecase.Update(ctx, tenantID, m); err != nil {
			ctx.JSON(
				500,
				gin.H{"error_reason": err.Error()},
			)
			return
		}
		ctx.Writer.WriteHeader(200)
	})
	r.Run("0.0.0.0:8080")
}

func getEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		v = defaultValue
	}
	return v
}
