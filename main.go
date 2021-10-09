package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
  "golang.org/x/crypto/bcrypt"

  "InstagramAPI/router"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {                                                                                 // {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`                          //   "name":"Riya",
	Name     string             `json:"name" bson:"name,omitempty"`                                  //   "email":"rkriya@gmail.com",
	Email    string             `json:"email" bson:"email,omitempty"`                                //   "password":"riya"
	Password string             `json:"password" bson:"password,omitempty"`                          // }
}

type Post struct {                                                                                 // {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`                   //     "userID":"id here",
	UserID          string             `json:"userId" bson:"userId,omitempt"`                        //     "Caption":"Hello World",
	Caption         string             `json:"caption" bson:"caption,omitempty"`                     //     "ImageURL":"https://cdn.pixabay.com/photo/2015/04/19/08/32/marguerite-729510__480.jpg"
	ImageURL        string             `json:"imageurl" bson:"imageurl,omitempty"`                   // }
	PostedTimeStamp string             `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

var users, posts = ConnecttoDB()

func main() {
	app := router.NewApp()

	app.Handle(`^/users$`, func(resp *router.Response, req *router.Request) {
		if req.Method == "POST" {
			resp.Header().Set("Content-Type", "application/json")

			var user User
			_ = json.NewDecoder(req.Body).Decode(&user)

      // Encrypting the password for secure storage
			hash, _ := HashPassword(user.Password)
      user.Password = hash
			result, err := users.InsertOne(context.TODO(), user)
			if err != nil {
				http.Error(resp, "Error", 404)
			}else {
				json.NewEncoder(resp).Encode(result)
			}
		} else {
			http.Error(resp, "Invalid request method.", 405)
		}
	})

	app.Handle(`/users/([\w\._-]+)$`, func(resp *router.Response, req *router.Request) {
		if req.Method == "GET" {
			resp.Header().Set("Content-Type", "application/json")
			var user User

			// string to primitive.ObjectID (typeCasting)
			id, _ := primitive.ObjectIDFromHex(req.Params[0])

			// creating filter of unordered map with ID as input
			filter := bson.M{"_id": id}

			//Searching in DB with given ID as keyword
			err := users.FindOne(context.TODO(), filter).Decode(&user)
			//Error Handling
			if err != nil {
				http.Error(resp, "Not Found", 404)
			} else {
				json.NewEncoder(resp).Encode(user)
			}
		} else {
			http.Error(resp, "Invalid request method.", 405)
		}
	})

	app.Handle(`/posts$`, func(resp *router.Response, req *router.Request) {
		if req.Method == "POST" {
			resp.Header().Set("Content-Type", "application/json")

			var post Post
			_ = json.NewDecoder(req.Body).Decode(&post)
			post.PostedTimeStamp = time.Now().Format(time.RFC850)
			result, err := posts.InsertOne(context.TODO(), post)
			if err != nil {
				http.Error(resp, "Error", 404)
			}else {
				json.NewEncoder(resp).Encode(result)
			}
		} else {
			http.Error(resp, "Invalid request method.", 405)
		}
	})

	app.Handle(`/posts/([\w\._-]+)$`, func(resp *router.Response, req *router.Request) {
		if req.Method == "GET" {
			resp.Header().Set("Content-Type", "application/json")
			var post Post

			// string to primitive.ObjectID (typeCasting)
			id, _ := primitive.ObjectIDFromHex(req.Params[0])

			// creating filter of unordered map with ID as input
			filter := bson.M{"_id": id}

			//Searching in DB with given ID as keyword
			err := posts.FindOne(context.TODO(), filter).Decode(&post)
			//Error Handling
			if err != nil {
				http.Error(resp, "Not Found", 404)
			} else {
				json.NewEncoder(resp).Encode(post)
			}
		} else {
			http.Error(resp, "Invalid request method.", 405)
		}
	})

	app.Handle(`/posts/user/([\w\._-]+)$`, func(resp *router.Response, req *router.Request) {
		if req.Method == "GET" {
			resp.Header().Set("Content-Type", "application/json")
			var postList []Post

			// creating filter of unordered map with ID as input
			filter := bson.M{"userId": req.Params[0]}

			//Searching in DB with given ID as keyword
			cur, err := posts.Find(context.TODO(), filter)
			//Error Handling
			if err != nil {
        http.Error(resp, "Not Found", 404)
      }
      defer cur.Close(context.TODO())
      for cur.Next(context.TODO()) {
          var post Post
          err := cur.Decode(&post)
          if err !=nil {
            http.Error(resp, "Not Found", 404)
          }
          postList = append(postList, post)
      }
			if err := cur.Err(); err != nil {
        http.Error(resp, "Error", 404)
			} else {
				json.NewEncoder(resp).Encode(postList)
			}
		} else {
			http.Error(resp, "Invalid request method.", 405)
		}
	})

	err := http.ListenAndServe(":9000", app)

	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

}

func ConnecttoDB() (*mongo.Collection, *mongo.Collection) {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	//Error Handling
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	//DB collection address which we are going to use
	users := client.Database("Appointy").Collection("users")
	posts := client.Database("Appointy").Collection("posts")

	return users, posts
}

//Function to encrypt password for secure storage using crupto library
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}
