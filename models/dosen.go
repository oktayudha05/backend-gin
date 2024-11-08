package models

type Dosen struct {
	Nama string `json:"nama" validate:"required"`
	NIP uint `json:"nip" validate:"required"`
	Jabatan string `json:"jabatan"`
	Asal string `json:"asal"`
}