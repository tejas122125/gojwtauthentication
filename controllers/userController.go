package controllers

import (
	"context"
	"fmt"
	"jwt/database"
	helpers "jwt/helpers"
	models "jwt/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)
 var userCollection *mongo.Collection = database.OpenCollection(database.Client,"monujindabad")
func HashPassword(password string) string {
	passwordGenerated, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(passwordGenerated)
}
func VerifyPassword(userPassword string, foundPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(foundPassword), []byte(userPassword))
	if err != nil {
		return false, "enter the correct password"
	} else {
		return true, "password matched"
	}

}
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
		var user models.User
		err := c.BindJSON(&user)
		
		if err != nil {
			
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot validate"})
		}

		validate := validator.New()
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		countEmail, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "server error in teechnical default",
			})
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		countPhone, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			fmt.Println("mmmmmmmmmmm")

			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "server error in technical default",
			})
		}

		if countEmail > 0 || countPhone > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": " this email or phone already exists in database enter a unique one"})
		}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		name:= user.ID.String()
		user.User_id = &name
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "there is some error while insertinf=g user informationin data base"})
		}
		c.JSON(http.StatusOK, resultInsertionNumber)
	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, Cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		var foundUser models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		err = userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer Cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}
		passwordValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer Cancel()
		if passwordValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		signedToken, signedRefreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *foundUser.User_id)

		helpers.UpdateAllTokens(signedToken, signedRefreshToken, *foundUser.User_id)
		userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
	}

}

//nhi kar sakta
// func GetUsers() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		err := helpers.ChechUserType(c, "ADMIN")
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "you idiot u are not the admin",
// 			})
// 		}
// 		var ctx,cancel = context.WithTimeout(context.Background(),time.Second*100)

// 	}

// }
func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userId := ctx.Param("user_id")

		err := helpers.MatchTypeToUid(ctx, userId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		var c, cancel = context.WithTimeout(context.Background(), time.Second*100)
		var user models.User

		err = userCollection.FindOne(c, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusAccepted, user)

	}
}
