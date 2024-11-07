package models

type Mahasiswa struct{
	Nama string `json:"nama" validate:"required"`
	NPM uint `json:"npm" validate:"required"`
	Prodi string `json:"prodi"`
	Angkatan uint `json:"angkatan"`
	Asal string `json:"asal"`
	Instagram string `json:"instagram"`
}
