package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := NewTenantDB(
		TenantDBInput{
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
		ctx.JSON(200, gin.H{"user": user})
	})
	r.POST("/users/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		tenantID := ctx.Request.Header.Get("tenant-id")
		var reqParams UpdateUserRequestParams
		if err := ctx.BindJSON(&reqParams); err != nil {
			ctx.JSON(
				400,
				gin.H{"error_reason": err.Error()},
			)
			return
		}
		if err := usecase.Update(ctx, tenantID, *ToUser(id, reqParams)); err != nil {
			ctx.JSON(
				500,
				gin.H{"error_reason": err.Error()},
			)
			return
		}
		ctx.Writer.WriteHeader(200)
	})
	r.Run("0.0.0.0:8080") // 0.0.0.0:8080 でサーバーを立てます。
}

func getEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		v = defaultValue
	}
	return v
}

type UpdateUserRequestParams struct {
	Name   string `json:"name,omitempty"`
	Gender string `json:"gender,omitempty"`
	Age    int    `json:"age,omitempty"`
}

func ToUser(userID string, p UpdateUserRequestParams) *User {
	return &User{
		ID:     userID,
		Name:   p.Name,
		Gender: p.Gender,
		Age:    p.Age,
	}
}
