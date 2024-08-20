package apartment_sale_api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var someKey = []byte("someKey")

func generateUserToken(user_type string) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Subject:   user_type,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(someKey)
}

func authorization(c *gin.Context, user_type string) {
	strToken, err := generateUserToken(user_type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", strToken, 3600*24*30, "/", "", false, true)

	c.Status(http.StatusOK)
}

func (this *SaleServer) GetDummyLogin(c *gin.Context) {
	user_type := c.Params.ByName("user_type")

	authorization(c, user_type)
}

func (this *SaleServer) PostLogin(c *gin.Context) {
	type requestInput struct {
		Id       string `json:"id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var input requestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := this.db.GetUser(input.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	authorization(c, user.UserType)
}

func (this *SaleServer) PostRegister(c *gin.Context) {
	type requestInput struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		UserType string `json:"user_type" binding:"required"`
	}
	var input requestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	uid, err := this.db.AddUser(input.Email, string(hash), input.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": uid})
}

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return someKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiresAt, err := claims.GetExpirationTime()
		if err != nil || time.Now().Unix() > expiresAt.Unix() {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		user_type, err := claims.GetSubject()
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user_type", user_type)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
