package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oktayudha05/backend-gin/database"
	"github.com/oktayudha05/backend-gin/models"
	"go.mongodb.org/mongo-driver/bson"
)

var dbDosen = database.GetDbDosen()

//get dosen
func GetDosen(c *gin.Context){
	var dosens []models.Dosen
	cur, err := dbDosen.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mendapatkan data"})
		return
	}

	for cur.Next(context.Background()) {
		var dosen models.Dosen
		cur.Decode(&dosen)
		dosens = append(dosens, dosen)
	}

	if len(dosens) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "data dosen belum ada"})
		return
	}
	c.IndentedJSON(http.StatusOK, dosens)
}

//post dosen
func PostDosen(c *gin.Context){
	var newDosen models.Dosen
	err := c.BindJSON(&newDosen)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "gagal bind data"})
		return
	}

	err = validate.Struct(newDosen)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "data yang dikirim tidak sesuai"})
		return
	}

	count, err := dbDosen.CountDocuments(context.Background(), bson.M{"nip": newDosen.NIP})
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "gagal mencari nip"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "nip sudah ada"})
		return
	}

	_, err = dbDosen.InsertOne(context.Background(), newDosen)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menambahkan data ke database"})
		return
	}
	c.IndentedJSON(http.StatusOK, newDosen)
}

// get dosen by nip
func GetDosenByNip(c *gin.Context){
	nip := c.Param("nip")
	nipUint, err := strconv.ParseUint(nip, 10, 32)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "gagal format nip"})
		return
	}

	var dosen models.Dosen
	err = dbDosen.FindOne(context.Background(), bson.M{"nip": uint(nipUint)}).Decode(&dosen)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "data dengan nip " + nip + " tidak ditemukan"})
		return
	}
	c.IndentedJSON(http.StatusOK, dosen)
}

// delete dosen
func DeleteDosenByNip(c *gin.Context){
	nip := c.Param("nip")
	nipUint, err := strconv.ParseUint(nip, 10, 32)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "gagal convert nip"})
		return
	}
	
	result, err := dbDosen.DeleteOne(context.Background(), bson.M{"nip": uint(nipUint)})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "gagal menghapus data"})
		return
	} else if result.DeletedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "nip yang diberikan tidak cocok"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "berhasil menghapus data dengan nip " + nip})
}