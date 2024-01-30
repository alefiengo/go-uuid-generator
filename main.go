package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	r.GET("/uuid/:version", func(c *gin.Context) {
		version := c.Param("version")
		uuidString, err := generateUUID(version)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to generate UUID: %v", err),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"uuid": uuidString,
		})
	})

	err := r.Run(":80")
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func generateUUID(version string) (string, error) {
	var uuidString string

	switch version {
	case "0":
		uuidInstance := uuid.Nil
		uuidString = uuidInstance.String()
	case "1":
		uuidInstance, err := uuid.NewV1()
		if err != nil {
			return "", err
		}
		uuidString = uuidInstance.String()
	case "4":
		uuidInstance, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		uuidString = uuidInstance.String()
	case "6":
		uuidInstance, err := uuid.NewV6()
		if err != nil {
			return "", err
		}
		uuidString = uuidInstance.String()
	case "7":
		uuidInstance, err := uuid.NewV7()
		if err != nil {
			return "", err
		}
		uuidString = uuidInstance.String()
	default:
		return "", fmt.Errorf("unsupported UUID version: %s", version)
	}

	return uuidString, nil
}
