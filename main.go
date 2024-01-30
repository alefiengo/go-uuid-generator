package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

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

	err := r.Run(":8081")
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func generateUUID(version string) (string, error) {
	var uuidString string

	switch version {
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
