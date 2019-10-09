package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, data interface{}) {
	_ = json.NewEncoder(w).Encode(data)
}

func DumpStructAsJson(object interface{}) {
	b, err := json.Marshal(object)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
