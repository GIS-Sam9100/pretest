package controller

import (
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


// CreateBarang creates a new Barang
func CreateBarang(w http.ResponseWriter, r *http.Request) {
	var barang model.BarangRequest
	if err := at.ReadJSON(w, r, &barang); err != nil {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Validate input
	if barang.Nama == "" || barang.Harga <= 0 {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Nama and Harga are required and Harga must be greater than 0"})
		return
	}

	// Create new Barang document
	newBarang := model.Barang{
		Nama:         barang.Nama,
		Harga:        barang.Harga,
		Stok:         barang.Stok,
		Deskripsi:    barang.Deskripsi,
		DibuatPada:   time.Now(),
		DiupdatePada: time.Now(),
	}

	// Insert to database
	insertedID, err := atdb.InsertOneDoc(config.Mongoconn, "barang", newBarang)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create barang"})
		return
	}

	// Get the created barang
	var createdBarang model.Barang
	createdBarang, err = atdb.GetOneDoc[model.Barang](config.Mongoconn, "barang", bson.M{"_id": insertedID})
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve created barang"})
		return
	}

	at.WriteJSON(w, http.StatusCreated, createdBarang)
}

// GetBarangByID gets a single Barang by ID
func GetBarangByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ID is required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
		return
	}

	var barang model.Barang
		barang, err = atdb.GetOneDoc[model.Barang](config.Mongoconn, "barang", bson.M{"_id": objID})
	if err != nil {
		at.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Barang not found"})
		return
	}

	at.WriteJSON(w, http.StatusOK, barang)
}

// GetAllBarang gets all Barang
func GetAllBarang(w http.ResponseWriter, r *http.Request) {
	var barangs []model.Barang
		barangs, err := atdb.GetAllDoc[[]model.Barang](config.Mongoconn, "barang", bson.M{})
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get barang list"})
		return
	}

	if barangs == nil {
		barangs = []model.Barang{}
	}

	at.WriteJSON(w, http.StatusOK, barangs)
}

// UpdateBarang updates an existing Barang
func UpdateBarang(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ID is required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
		return
	}

	var barang model.BarangRequest
	if err := at.ReadJSON(w, r, &barang); err != nil {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Check if barang exists
	count, err := atdb.GetCountDoc(config.Mongoconn, "barang", bson.M{"_id": objID})
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to check document existence"})
		return
	}
	if count == 0 {
		at.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Barang not found"})
		return
	}

	// Prepare update data
	updatefields := bson.M{
		"nama":          barang.Nama,
		"harga":         barang.Harga,
		"stok":          barang.Stok,
		"deskripsi":     barang.Deskripsi,
		"diupdate_pada": time.Now(),
	}

	_, err = atdb.UpdateOneDoc(config.Mongoconn, "barang", bson.M{"_id": objID}, updatefields)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update barang"})
		return
	}

	// Get the updated barang
	var updatedBarang model.Barang
	updatedBarang, err = atdb.GetOneDoc[model.Barang](config.Mongoconn, "barang", bson.M{"_id": objID})
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve updated barang"})
		return
	}

	at.WriteJSON(w, http.StatusOK, updatedBarang)
}

// DeleteBarang deletes a Barang by ID
func DeleteBarang(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ID is required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
		return
	}

	// Check if barang exists
	count, err := atdb.GetCountDoc(config.Mongoconn, "barang", bson.M{"_id": objID})
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to check document existence"})
		return
	}
	if count == 0 {
		at.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Barang not found"})
		return
	}

	_, err = atdb.DeleteOneDoc(config.Mongoconn, "barang", bson.M{"_id": objID})
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete barang"})
		return
	}

	at.WriteJSON(w, http.StatusOK, map[string]string{"message": "Barang deleted successfully"})
}
