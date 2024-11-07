package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/oktayudha05/backend-gin/models"
	"go.mongodb.org/mongo-driver/bson"
)

var validate *validator.Validate

func main(){
	initDB()
	validate = validator.New()
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

// get mahasiswa
func getMahasiswa(c *gin.Context){
	var mahasiswas []models.Mahasiswa
	cur, err := db.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal ambil data"})
		return
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var mahasiswa models.Mahasiswa
		cur.Decode(&mahasiswa)
		mahasiswas = append(mahasiswas, mahasiswa) 
	}
	
	if len(mahasiswas) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "belum ada mahasiswa"})
		return
	}
	c.IndentedJSON(http.StatusOK, mahasiswas)
}

// post mahasiswa
func postMahasiswa(c *gin.Context){
	var newMahasiswa models.Mahasiswa

	err := c.BindJSON(&newMahasiswa) 
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"message": "format data salah"})
		return
	}
	
	err = validate.Struct(newMahasiswa)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "data tidak lengkap"})
		return
	}

	count, err := db.CountDocuments(context.Background(), bson.M{"npm": newMahasiswa.NPM})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal cek npm"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "NPM sudah terdaftar"})
		return
	}

	_, err = db.InsertOne(context.Background(), newMahasiswa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal simpan data ke database"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newMahasiswa)
}

// get mahasiswa by npm
func getMahasiswaByNpm(c *gin.Context){
	npm := c.Param("npm")
	
	npmUint, err := strconv.ParseUint(npm, 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "NPM salah atau tidak ada"})
		return
	}

	var mahasiswa models.Mahasiswa
	err = db.FindOne(context.Background(), bson.M{"npm": uint(npmUint)}).Decode(&mahasiswa)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "mahasiswa dengan npm " + npm + " tidak ditemukan"})
		return
	}
	c.IndentedJSON(http.StatusOK, mahasiswa)
}

// delete mahasiswa by npm
func deleteMahasiswaByNpm(c *gin.Context){
	npm := c.Param("npm")

	npmUint, err := strconv.ParseUint(npm, 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "format npm salah"})
		return
	}

	result, err := db.DeleteOne(context.Background(), bson.M{"npm": uint(npmUint)})
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "gagal menghapus data mahasiswa"})
		return
	} else if result.DeletedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "npm yang diberikan tidak cocok"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "berhasil menghapus data dengan npm " + npm})
}