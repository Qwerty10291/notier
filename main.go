package main

import (
	"encoding/json"
	"fmt"
	"gin_backend/storage"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Status struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

var DB *sqlx.DB

func main() {
	config := LoadConfig()
	router := gin.Default()
	conn, err := storage.InitDb("db.sqlite")
	if err != nil {
		panic(err)
	}
	DB = conn
	router.Static("/static", "./static")
	router.LoadHTMLFiles("templates/index.html")
	router.GET("/", Index)
	router.GET("/notes", GetNotesHandler)
	router.POST("/notes", NewNoteHandler)
	router.DELETE("/notes/:id", DeleteNotesHandler)
	router.Run(fmt.Sprintf("%s:%s", config.Host, config.Port))
}

func Index(c *gin.Context){
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func GetNotesHandler(c *gin.Context) {
	notes, err := storage.GetAllNotes(DB)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, Status{Error: "unknown server error"})
		return
	}
	c.JSON(http.StatusOK, notes)
}

func NewNoteHandler(c *gin.Context) {
	text := c.PostForm("text")
	if len(text) == 0 {
		c.JSON(http.StatusBadRequest, Status{Error: "note text must be specified"})
		return
	}
	note, err := storage.NewNote(DB, text)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, Status{Error: "unknown server error"})
		return
	}
	c.JSON(http.StatusCreated, note)
}

func DeleteNotesHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Status{Error: fmt.Sprintf("cannot parse %s as int", c.Param("id"))})
		return
	}
	err = storage.DeleteNote(DB, int(id))
	if err != nil {
		if err == storage.ErrorNoteNotFound {
			c.JSON(http.StatusNotFound, Status{Error: fmt.Sprintf("error with id %d not found", id)})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, Status{Error: "unknown server error"})
	}
	c.JSON(http.StatusOK, Status{Status: true})
}

func LoadConfig() ServerConfig {
	if file, err := os.Stat("config.json"); err != nil || file == nil {
		log.Println("config file open failed. Use default")
		return ServerConfig{
			Host: "127.0.0.1",
			Port: "8000",
		}
	}
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Println("config file open failed. Use default")
		return ServerConfig{
			Host: "127.0.0.1",
			Port: "8000",
		}
	}
	decoder := json.NewDecoder(configFile)
	var config ServerConfig
	err = decoder.Decode(&config)
	if err != nil {
		log.Println("config file corrupted. Use default")
		return ServerConfig{
			Host: "127.0.0.1",
			Port: "8000",
		}
	}
	log.Println("config file loaded")
	return config
}
