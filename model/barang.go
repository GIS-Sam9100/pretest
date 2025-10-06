package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Barang struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Nama        string             `json:"nama" bson:"nama"`
	Harga       float64            `json:"harga" bson:"harga"`
	Stok        int                `json:"stok" bson:"stok"`
	Deskripsi   string             `json:"deskripsi,omitempty" bson:"deskripsi,omitempty"`
	DibuatPada  time.Time          `json:"dibuat_pada,omitempty" bson:"dibuat_pada,omitempty"`
	DiupdatePada time.Time          `json:"diupdate_pada,omitempty" bson:"diupdate_pada,omitempty"`
}

type BarangRequest struct {
	Nama      string  `json:"nama" validate:"required"`
	Harga     float64 `json:"harga" validate:"required,gt=0"`
	Stok      int     `json:"stok" validate:"gte=0"`
	Deskripsi string  `json:"deskripsi,omitempty"`
}
