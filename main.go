package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	godotenv.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("❌ Error al conectar a Mongo:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ Mongo no responde:", err)
	}
	log.Println("🍃 ¡Conexión exitosa a MongoDB en el puerto 27017!")

	r := gin.Default()
	r.Use(CORSMiddleware())

	// Ruta 1: Test de conexión
	r.GET("/api/test-mongo", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"exito": true, "mensaje": "¡API de Go conectada a MongoDB!"})
	})

	// Ruta 2: POBLAR LA BASE DE DATOS (Solo la corremos una vez)
	r.GET("/api/seed", func(c *gin.Context) {
		coleccion := client.Database(os.Getenv("DB_NAME")).Collection("kardex")

		// Borramos lo anterior para no duplicar si le das F5
		coleccion.Drop(context.TODO())

		// Tu Kardex real sacado de la imagen
		miKardex := bson.M{
			"boleta":          "2022350438",
			"promedio_global": 8.81,
			"plan":            "04",
			"carrera":         "INGENIERIA EN COMPUTACION",
			"historial": []bson.M{
				{
					"semestre": 1,
					"materias": []bson.M{
						bson.M{"clave": "C101", "materia": "CALCULO DIFERENCIAL E INTEGRAL", "calificacion": 7},
						bson.M{"clave": "C102", "materia": "FISICA CLASICA", "calificacion": 10},
						bson.M{"clave": "C103", "materia": "FUNDAMENTOS DE PROGRAMACION", "calificacion": 8},
					},
				},
				{
					"semestre": 2,
					"materias": []bson.M{
						bson.M{"clave": "C207", "materia": "CALCULO VECTORIAL", "calificacion": 8},
						bson.M{"clave": "C208", "materia": "ELECTRICIDAD Y MAGNETISMO", "calificacion": 9},
					},
				},
			},
		}

		_, err := coleccion.InsertOne(context.TODO(), miKardex)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error guardando en Mongo"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"exito": true, "mensaje": "Kardex de 2022350438 insertado exitosamente"})
	})

	// Ruta 3: OBTENER EL KARDEX (Esta es la que usará el Frontend)
	r.GET("/api/kardex/:boleta", func(c *gin.Context) {
		boletaParam := c.Param("boleta")
		coleccion := client.Database(os.Getenv("DB_NAME")).Collection("kardex")

		var resultado bson.M
		err := coleccion.FindOne(context.TODO(), bson.M{"boleta": boletaParam}).Decode(&resultado)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"exito": false, "mensaje": "Kardex no encontrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"exito": false, "mensaje": "Error buscando en base de datos"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"exito": true, "data": resultado})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("🚀 Servidor Go escuchando en http://localhost:%s", port)
	r.Run(":" + port)
}
