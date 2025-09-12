package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// album represents data about a record album.
type album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./albums.db")
	if err != nil {
		return err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS albums (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		artist TEXT NOT NULL,
		price REAL NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	// Check if table is empty and seed with initial data
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM albums").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		seedData := []album{
			{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
			{Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
			{Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
		}

		for _, a := range seedData {
			_, err = db.Exec("INSERT INTO albums (title, artist, price) VALUES (?, ?, ?)",
				a.Title, a.Artist, a.Price)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// setupRouter configures and returns the Gin router with all routes
func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	return router
}

func main() {
	if err := initDB(); err != nil {
		panic(err)
	}
	defer db.Close()

	router := setupRouter()
	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, artist, price FROM albums")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var albums []album
	for rows.Next() {
		var a album
		err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		albums = append(albums, a)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Insert the new album into the database and get the generated ID.
	result, err := db.Exec("INSERT INTO albums (title, artist, price) VALUES (?, ?, ?)",
		newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the auto-generated ID
	id, err := result.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newAlbum.ID = id

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	var a album
	err := db.QueryRow("SELECT id, title, artist, price FROM albums WHERE id = ?", id).
		Scan(&a.ID, &a.Title, &a.Artist, &a.Price)

	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, a)
}
