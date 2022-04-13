// Package Classification for Prime API
//
// Documentation for Prime API
//
// 	Schemes: http
//  Host: localhost
//	BasePath: /
// 	Version: 1.0.0
//
//	Consumes:
// 	- application/json
//
//	Produces:
// 	- application/json
//
// Swagger:meta

package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-mongo/model"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://sombir:root@cluster0.udy9f.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
const dbName, colName = "prime", "watchlist"

var collection *mongo.Collection

// Connect with mongoDB

func init() {
	//client option
	clientOption := options.Client().ApplyURI(connectionString)

	//connect to mongoDB
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDb connection established successfully!")
	collection = client.Database(dbName).Collection(colName)

	//collection Instance -->
	fmt.Println("Collection Instance is ready...")
}

func SourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	path := r.FormValue("path")
	lineStr := r.FormValue("line")
	line, err := strconv.Atoi(lineStr)
	if err != nil {
		line = -1
	}
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b := bytes.NewBuffer(nil)

	_, err = io.Copy(b, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var lines [][2]int
	if line > 0 {
		lines = append(lines, [2]int{line, line})
	}
	lexer := lexers.Get("go")
	iterator, _ := lexer.Tokenise(nil, b.String())
	style := styles.Get("github")

	if style == nil { 
		style = styles.Fallback
	}

	formatter := html.New(html.TabWidth(2), html.HighlightLines(lines), html.WithLineNumbers(true))
	formatter.Format(w, style, iterator)

	// _ = quick.Highlight(w, b.String(), "go", "html", "monokai")
}

func MakeLinks(stack string) string {
	lines := strings.Split(stack, "\n")
	for li, line := range lines {
		if len(line) == 0 || line[0] != '\t' {
			continue
		}
		file := ""
		for i, ch := range line {
			if ch == ':' {
				file = line[1:i]
				break
			}
		}

		var linestr strings.Builder
		for i := len(file) + 2; i < len(line); i++ {
			if line[i] < '0' || line[i] > '9' {
				break
			}
			linestr.WriteByte(line[i])
		}

		//dynamic path
		v := url.Values{}
		v.Set("path", file)
		v.Set("line", linestr.String())
		lines[li] = "\t<a href= \"/debug/?" + v.Encode() + "\">" + file + ":" + linestr.String() + "</a>" + line[len(file)+2+len(linestr.String()):]
	}
	return strings.Join(lines, "\n")
}

// GET by ID
func GetAMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	_id, _ := primitive.ObjectIDFromHex(params["id"])
	var movieResult model.Prime
	filter := bson.M{"_id": _id}

	err := collection.FindOne(context.Background(), filter).Decode(&movieResult)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(movieResult)
}

// GET Method
func GetMyAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencode")

	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	var movies []primitive.M

	for cur.Next(context.Background()) {
		var movie bson.M
		err := cur.Decode(&movie)
		if err != nil {
			fmt.Println(err)
		}
		movies = append(movies, movie)
	}

	// allmovies := getAllMovies()
	json.NewEncoder(w).Encode(movies)
}

// POST Method
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Prime
	_ = json.NewDecoder(r.Body).Decode(&movie)

	res, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(res)
}

// UPDATE Method
func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)

	_id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(result)
}

// Delete ONE
func DeleteAMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	_id, _ := primitive.ObjectIDFromHex(params["id"])
	dCount, err := collection.DeleteOne(context.Background(), bson.M{"_id": _id})
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(dCount)
}

// Delete Method
func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	count, err := collection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(count)
}
