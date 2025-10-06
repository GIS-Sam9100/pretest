package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GeoJSON represents a GeoJSON structure for MongoDB.
// MongoDB expects coordinates in [longitude, latitude] order.
type GeoJSON struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

// Lokasi represents a geographic location entity.
type Lokasi struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Nama        string             `json:"nama" bson:"nama"`
	Kategori    string             `json:"kategori" bson:"kategori"`
	Deskripsi   string             `json:"deskripsi" bson:"deskripsi"`
	Koordinat   GeoJSON            `json:"koordinat" bson:"koordinat"`
}
