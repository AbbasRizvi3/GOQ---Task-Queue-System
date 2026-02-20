package recovery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware(c *gin.Context) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
			jsonBody, _ := json.Marshal(map[string]string{
				"error": "There was an internal server error",
			})

			c.Header("Content-Type", "application/json")
			c.AbortWithStatus(http.StatusInternalServerError)
			_, err := c.Writer.Write(jsonBody)
			if err != nil {
				fmt.Printf("Error writing JSON response: %v\n", err)
			}
		}
	}()

	c.Next()
}
