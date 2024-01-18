package helpers

import (
	"context"
	"fmt"
	"jwt/database"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type signedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.StandardClaims
}

var collections *mongo.Collection = database.OpenCollection(database.Client, "user")
var SECRET_KEY string = "tejaswee"

func GenerateAllTokens(email string, firstname string, lastname string, usertype string, uid string) (signedToken string, signedRefreshToken string, err error) {

	claims := &signedDetails{
		Email:      email,
		First_name: firstname,
		Last_name:  lastname,
		Uid:        uid,
		User_type:  usertype,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaim := &signedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString([]byte(SECRET_KEY))

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaim).SignedString([]byte(SECRET_KEY))

	if err != nil {
		fmt.Println("mmmmmmmmmmm")

		log.Panic(err)
		return
	}
	return token, refreshToken, err

}
func ValidateToken(signedToken string) (claims *signedDetails, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&signedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*signedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token has expired")
		msg = err.Error()
		return
	}

	return claims, msg

}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	var userCollection *mongo.Collection = database.OpenCollection(database.Client,"monujindabad")
	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}

	return
}
