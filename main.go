package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/juliangruber/go-intersect"
)

type bingoResponse struct {
	Token string   `json:"token" binding:"required"`
	Sids  []string `json:"sids" binding:"required"`
}

type targetsidResponse struct {
	Sid string `json:"sid" binding:"required"`
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
	router := gin.Default()

	router.Use(cors.Default()) // TODO: limit origin

	// TODO
	tokenMapping := make(map[string]string)
	tokenMapping["0416"] = "0416"
	tokenMapping["0432"] = "0432"

	router.GET("/token/:token", func(c *gin.Context) {
		var res bool
		token := c.Param("token")
		if tokenMapping[token] == "" {
			res = false
		} else {
			res = true
		}
		c.JSON(200, gin.H{
			"status": res,
		})

	})

	router.GET("/targetsids", func(c *gin.Context) {
		targetSids := getTargetSids()
		c.JSON(200, gin.H{
			"targetSids": targetSids,
		})
	})

	router.PUT("/targetsids", func(c *gin.Context) {
		var json targetsidResponse
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		insertTargetSid(json.Sid)
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.DELETE("/targetsids", func(c *gin.Context) {
		clearTargetSids()
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/winners", func(c *gin.Context) {
		winners := getWinners()
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

		clientSid := tokenMapping[json.Token]
		if clientSid == "" {
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
			insertWinner(clientSid)
		} else {
			status = false
		}

		c.JSON(200, gin.H{
			"status": status,
		})
	})
	router.Run(":8000")
}
