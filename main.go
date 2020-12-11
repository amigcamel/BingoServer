package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/juliangruber/go-intersect"
)

type bingoResponse struct {
	Token     string   `json:"token" binding:"required"`
	Sids      []string `json:"sids" binding:"required"`
	ClientSid string   `json:"clientSid" binding:"required"`
}

func testEq(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func main() {
	fmt.Println()
	router := gin.Default()

	router.Use(cors.Default()) // TODO: limit origin

	router.GET("/winners", func(c *gin.Context) {
		winners := getWinners()
		fmt.Println("===============================")
		fmt.Println(winners)
		c.JSON(200, gin.H{
			"winners": winners,
		})
	})

	router.POST("/api", func(c *gin.Context) {
		var json bingoResponse
		var status bool
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if json.Token != "fake-token" {
			c.JSON(400, gin.H{"message": "Bad token"})
		}

		allSids := getTargetSids()
		fmt.Println(json.Sids)
		res := intersect.Simple(json.Sids, allSids)
		var intersec []string
		for _, value := range res.([]interface{}) {
			intersec = append(intersec, value.(string))
		}

		ok := testEq(intersec, json.Sids)
		if ok {
			status = true
			insertWinner(json.ClientSid)
		} else {
			status = false
		}

		c.JSON(200, gin.H{
			"status": status,
		})
	})
	router.Run(":8000")
}
