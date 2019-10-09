package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sinistra/go-gin-api/auth"
	"sinistra/go-gin-api/models"
)

type UserController struct{}

func (u UserController) Login(c *gin.Context) {
	var user models.User
	var jwt models.JWT

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "login failed binding.", "error": err.Error()})
	}

	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username is missing.", "data": user})
		return
	}
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "password is missing.", "data": user})
		return
	}

	ok := auth.LdapValidate(user.Username, user.Password)

	if ok == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "authentication failed", "data": user})
		return
	}

	token, err := auth.GenerateToken(user)

	if err != nil {
		log.Fatal(err)
	}

	jwt.Token = token
	c.JSON(http.StatusOK, gin.H{"message": "ok", "data": jwt})

}

func (u UserController) CheckForToken(c *gin.Context) {

	status, msg := auth.TokenVerify(c)

	if status == http.StatusOK {
		output := auth.DecodeToken(c)
		c.JSON(http.StatusOK, gin.H{"message": "ok", "data": output})
	} else {
		c.JSON(status, gin.H{"message": msg})
	}
}

func (u UserController) TestAuth(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username not in context"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok", "data": username})
}
