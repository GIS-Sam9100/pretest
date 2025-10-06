package controller

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateLokasi menyisipkan lokasi baru ke database
func CreateLokasi(w http.ResponseWriter, r *http.Request) {
	var lokasi model.Lokasi
	if err := at.ReadJSON(w, r, &lokasi); err != nil {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Validasi input GeoJSON
	if lokasi.Koordinat.Type != "Point" || len(lokasi.Koordinat.Coordinates) != 2 {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Koordinat tidak valid. Harus berupa GeoJSON Point dengan [longitude, latitude]."})
		return
	}

	insertedID, err := atdb.InsertOneDoc(config.Mongoconn, "lokasi", lokasi)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal membuat lokasi"})
		return
	}

	createdLokasi, _ := atdb.GetOneDoc[model.Lokasi](config.Mongoconn, "lokasi", bson.M{"_id": insertedID})
	at.WriteJSON(w, http.StatusCreated, createdLokasi)
}

// GetLokasiByID mengambil satu lokasi berdasarkan ID
func GetLokasiByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ID dibutuhkan"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Format ID tidak valid"})
		return
	}

	lokasi, err := atdb.GetOneDoc[model.Lokasi](config.Mongoconn, "lokasi", bson.M{"_id": objID})
	if err != nil {
		at.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Lokasi tidak ditemukan"})
		return
	}

	at.WriteJSON(w, http.StatusOK, lokasi)
}

// GetAllLokasi mengambil semua data lokasi
func GetAllLokasi(w http.ResponseWriter, r *http.Request) {
	lokasis, err := atdb.GetAllDoc[[]model.Lokasi](config.Mongoconn, "lokasi", bson.M{})
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil daftar lokasi"})
		return
	}

	if lokasis == nil {
		lokasis = []model.Lokasi{}
	}

	at.WriteJSON(w, http.StatusOK, lokasis)
}

// UpdateLokasi memperbarui data lokasi yang ada
func UpdateLokasi(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ID dibutuhkan"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Format ID tidak valid"})
		return
	}

	var lokasiData model.Lokasi
	if err := at.ReadJSON(w, r, &lokasiData); err != nil {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Validasi input GeoJSON
	if lokasiData.Koordinat.Type != "Point" || len(lokasiData.Koordinat.Coordinates) != 2 {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Koordinat tidak valid. Harus berupa GeoJSON Point dengan [longitude, latitude]."})
		return
	}

	updateFields := bson.M{
		"nama":      lokasiData.Nama,
		"kategori":  lokasiData.Kategori,
		"deskripsi": lokasiData.Deskripsi,
		"koordinat": lokasiData.Koordinat,
	}

	_, err = atdb.UpdateOneDoc(config.Mongoconn, "lokasi", bson.M{"_id": objID}, updateFields)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memperbarui lokasi"})
		return
	}

	updatedLokasi, _ := atdb.GetOneDoc[model.Lokasi](config.Mongoconn, "lokasi", bson.M{"_id": objID})
	at.WriteJSON(w, http.StatusOK, updatedLokasi)
}

// DeleteLokasi menghapus lokasi berdasarkan ID
func DeleteLokasi(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ID dibutuhkan"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		at.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Format ID tidak valid"})
		return
	}

	_, err = atdb.DeleteOneDoc(config.Mongoconn, "lokasi", bson.M{"_id": objID})
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal menghapus lokasi"})
		return
	}

	at.WriteJSON(w, http.StatusOK, map[string]string{"message": "Lokasi berhasil dihapus"})
}
