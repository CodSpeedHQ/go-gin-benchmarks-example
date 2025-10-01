package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func BenchmarkGetAlbums(b *testing.B) {
	req, _ := http.NewRequest("GET", "/albums", nil)
	benchmarkRequest(b, req)
}

func BenchmarkGetAlbumByIDExists(b *testing.B) {
	req, _ := http.NewRequest("GET", "/albums/1", nil)
	benchmarkRequest(b, req)
}

func BenchmarkGetAlbumByIDNotFound(b *testing.B) {
	req, _ := http.NewRequest("GET", "/albums/999", nil)
	benchmarkRequest(b, req)
}

func BenchmarkPostAlbumsValid(b *testing.B) {
	newAlbum := album{
		Title:  "Kind of Blue",
		Artist: "Miles Davis",
		Price:  29.99,
	}
	albumJSON, _ := json.Marshal(newAlbum)
	req, _ := http.NewRequest("POST", "/albums", bytes.NewBuffer(albumJSON))
	req.Header.Set("Content-Type", "application/json")
	benchmarkRequest(b, req)
}

func BenchmarkPostAlbumsInvalidJSON(b *testing.B) {
	invalidJSON := `{"title": "Invalid Album", "artist": "Test Artist", "price": "invalid_price"}`
	req, _ := http.NewRequest("POST", "/albums", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	benchmarkRequest(b, req)
}

func BenchmarkPostAlbumsEmptyBody(b *testing.B) {
	req, _ := http.NewRequest("POST", "/albums", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	benchmarkRequest(b, req)
}
