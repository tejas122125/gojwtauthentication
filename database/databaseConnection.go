package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client{
err := godotenv.Load(".env")

if err != nil {
	log.Fatal("error while loading .env file in go")

}
mongoDb := os.Getenv("MONGODB_URL")
ctx,cancel := context.WithTimeout(context.Background(),time.Second*10)
defer cancel()
client , err := mongo.Connect(ctx,options.Client().ApplyURI(mongoDb))
if err != nil{
	log.Fatal("mongo error")
}
fmt.Println("mongo connection successful")

return client


}
var Client *mongo.Client = DBinstance()
func OpenCollection (client *mongo.Client,collectionName string )*mongo.Collection{


	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}