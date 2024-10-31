package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()

	router.GET("/", func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})
	router.GET("/mahasiswa", getMahasiswa)
	router.GET("/mahasiswa/:npm", getMahasiswaByNpm)
	router.POST("/mahasiswa", postMahasiswa)
	router.DELETE("/mahasiswa/:npm", deleteMahasiswaByNpm)

	router.Run(":3000")
}

type Mahasiswa struct{
	Nama string `json:"nama"`
	NPM uint `json:"npm"`
	Prodi string `json:"prodi"`
}

var mahasiswas = []Mahasiswa{
	{Nama: "Oktario", NPM: 2320506044, Prodi: "Teknologi Informasi"},
	{Nama: "Bambang", NPM: 2320506046, Prodi: "Teknologi Informasi"},
	{Nama: "Istoro", NPM: 2340506088, Prodi: "Teknologi Informasi"},
}

// get mahasiswa
func getMahasiswa(c *gin.Context){
	c.IndentedJSON(http.StatusOK, mahasiswas)
}

// post mahasiswa
func postMahasiswa(c *gin.Context){
	var newMahasiswa Mahasiswa

	err := c.BindJSON(&newMahasiswa) 
	if err != nil{
		panic(err)
	}
	mahasiswas = append(mahasiswas, newMahasiswa)
	c.IndentedJSON(http.StatusCreated, newMahasiswa)
}

// get mahasiswa by npm
func getMahasiswaByNpm(c *gin.Context){
	npm := c.Param("npm")
	
	npmUint, err := strconv.ParseUint(npm, 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NPM salah atau tidak ada"})
		return
	}

	for _, a := range mahasiswas {
		if a.NPM == uint(npmUint) {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "mahasiswa not found"})
}

// delete mahasiswa by npm
func deleteMahasiswaByNpm(c *gin.Context){
	npm := c.Param("npm")

	npmUint, err := strconv.ParseUint(npm, 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NPM salah atau tidak ada"})
		return
	}

	index := -1
	for i, a := range mahasiswas {
		if a.NPM == uint(npmUint) {
			index = i
			break
		}
	}
	if index == -1 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NPM salah atau tidak ada"})
		return
	}

	mahasiswas = append(mahasiswas[:index], mahasiswas[index+1:]...)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "mahasiswa berhasil dihapus"})
}