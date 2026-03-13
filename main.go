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
	// Ruta 2: POBLAR LA BASE DE DATOS (Catálogo y Kardex)
	r.GET("/api/seed", func(c *gin.Context) {
		// 1. POBLAR EL PLAN DE ESTUDIOS
		coleccionPlan := client.Database(os.Getenv("DB_NAME")).Collection("plan_estudios")
		coleccionPlan.Drop(context.TODO()) // Limpiamos para no duplicar

		// Arreglo con las materias, semestres y créditos (Sacados de tu PDF)
		// Arreglo con las materias, semestres y créditos exactos
		materias := []interface{}{
			// SEMESTRE 1
			bson.M{"clave": "C101", "materia": "CÁLCULO DIFERENCIAL E INTEGRAL", "semestre": 1, "creditos": 12.0},
			bson.M{"clave": "C102", "materia": "FÍSICA CLÁSICA", "semestre": 1, "creditos": 10.5},
			bson.M{"clave": "C103", "materia": "FUNDAMENTOS DE ÁLGEBRA", "semestre": 1, "creditos": 6.0},
			bson.M{"clave": "C104", "materia": "FUNDAMENTOS DE PROGRAMACIÓN", "semestre": 1, "creditos": 9.0},
			bson.M{"clave": "C105", "materia": "HUMANIDADES I: INGENIERÍA, CIENCIA Y SOCIEDAD", "semestre": 1, "creditos": 9.0},
			bson.M{"clave": "C106", "materia": "QUÍMICA BÁSICA", "semestre": 1, "creditos": 7.5},

			// SEMESTRE 2
			bson.M{"clave": "C207", "materia": "ÁLGEBRA LINEAL", "semestre": 2, "creditos": 6.0},
			bson.M{"clave": "C208", "materia": "CÁLCULO VECTORIAL", "semestre": 2, "creditos": 12.0},
			bson.M{"clave": "C209", "materia": "ELECTRICIDAD Y MAGNETISMO", "semestre": 2, "creditos": 10.5},
			bson.M{"clave": "C210", "materia": "HUMANIDADES II: LA COMUNICACIÓN Y LA INGENIERÍA", "semestre": 2, "creditos": 6.0},
			bson.M{"clave": "C211", "materia": "MATEMÁTICAS DISCRETAS", "semestre": 2, "creditos": 12.0},
			bson.M{"clave": "C212", "materia": "PROGRAMACIÓN ORIENTADA A OBJETOS", "semestre": 2, "creditos": 9.0},

			// SEMESTRE 3
			bson.M{"clave": "C313", "materia": "CIRCUITOS DE CA Y CD", "semestre": 3, "creditos": 7.5},
			bson.M{"clave": "C314", "materia": "CIRCUITOS LÓGICOS I", "semestre": 3, "creditos": 10.5},
			bson.M{"clave": "C315", "materia": "ECUACIONES DIFERENCIALES", "semestre": 3, "creditos": 9.0},
			bson.M{"clave": "C316", "materia": "ESTRUCTURA DE DATOS", "semestre": 3, "creditos": 9.0},
			bson.M{"clave": "C317", "materia": "HUMANIDADES III: DESARROLLO HUMANO", "semestre": 3, "creditos": 6.0},
			bson.M{"clave": "C318", "materia": "LENGUAJES DE BAJO NIVEL", "semestre": 3, "creditos": 9.0},

			// SEMESTRE 4
			bson.M{"clave": "C419", "materia": "ANÁLISIS NUMÉRICO", "semestre": 4, "creditos": 7.5},
			bson.M{"clave": "C420", "materia": "CIRCUITOS LÓGICOS II", "semestre": 4, "creditos": 7.5},
			bson.M{"clave": "C421", "materia": "ELECTRÓNICA ANALÓGICA", "semestre": 4, "creditos": 10.5},
			bson.M{"clave": "C422", "materia": "HUMANIDADES IV: DESARROLLO PERSONAL Y PROFESIONAL", "semestre": 4, "creditos": 6.0},
			bson.M{"clave": "C423", "materia": "TEORÍA DE AUTÓMATAS", "semestre": 4, "creditos": 12.0},
			bson.M{"clave": "C424", "materia": "VARIABLE COMPLEJA Y ANÁLISIS DE FOURIER", "semestre": 4, "creditos": 12.0},

			// SEMESTRE 5
			bson.M{"clave": "C525", "materia": "ANÁLISIS DE SEÑALES ANALÓGICAS", "semestre": 5, "creditos": 10.5},
			bson.M{"clave": "C526", "materia": "ANÁLISIS DE ALGORITMOS", "semestre": 5, "creditos": 7.5},
			bson.M{"clave": "C527", "materia": "COMPILADORES", "semestre": 5, "creditos": 7.5},
			bson.M{"clave": "C528", "materia": "HUMANIDADES V: EL HUMANISMO FRENTE A LA GLOBALIZACIÓN", "semestre": 5, "creditos": 9.0},
			bson.M{"clave": "C529", "materia": "ORGANIZACIÓN DE COMPUTADORAS", "semestre": 5, "creditos": 10.5},
			bson.M{"clave": "C530", "materia": "PROBABILIDAD Y ESTADÍSTICA", "semestre": 5, "creditos": 9.0},

			// SEMESTRE 6
			bson.M{"clave": "C631", "materia": "ARQUITECTURA DE COMPUTADORAS", "semestre": 6, "creditos": 7.5},
			bson.M{"clave": "C632", "materia": "INGENIERÍA DE SOFTWARE", "semestre": 6, "creditos": 7.5},
			bson.M{"clave": "C633", "materia": "METODOLOGÍA DE LA INVESTIGACIÓN Ó TÓPICOS SELECTOS DE INGENIERÍA I", "semestre": 6, "creditos": 6.0},
			bson.M{"clave": "C634", "materia": "MODULACIÓN DIGITAL", "semestre": 6, "creditos": 10.5},
			bson.M{"clave": "C635", "materia": "SISTEMAS OPERATIVOS", "semestre": 6, "creditos": 9.0},
			bson.M{"clave": "C636", "materia": "TEORÍA DE CONTROL ANALÓGICO", "semestre": 6, "creditos": 10.5},

			// SEMESTRE 7
			bson.M{"clave": "C737", "materia": "ADMINISTRACIÓN EN LA INGENIERÍA", "semestre": 7, "creditos": 9.0},
			bson.M{"clave": "C738", "materia": "BASES DE DATOS", "semestre": 7, "creditos": 7.5},
			bson.M{"clave": "C739", "materia": "NUEVAS TECNOLOGÍAS EN LA TRANSFERENCIA DE LA INFORMACIÓN", "semestre": 7, "creditos": 7.5},
			bson.M{"clave": "C740", "materia": "OPTATIVA I", "semestre": 7, "creditos": 6.0},
			bson.M{"clave": "C741", "materia": "TEORÍA DE CONTROL DIGITAL", "semestre": 7, "creditos": 10.5},
			bson.M{"clave": "C742", "materia": "TEORÍA DE LA INFORMACIÓN Y CODIFICACIÓN", "semestre": 7, "creditos": 7.5},

			// SEMESTRE 8
			bson.M{"clave": "C843", "materia": "FORMULACIÓN Y EVALUACIÓN DE PROYECTOS", "semestre": 8, "creditos": 7.5},
			bson.M{"clave": "C844", "materia": "PROYECTO DE INGENIERÍA Ó TÓPICOS SELECTOS DE INGENIERÍA II", "semestre": 8, "creditos": 6.0},
			bson.M{"clave": "C845", "materia": "OPTATIVA II", "semestre": 8, "creditos": 4.5},
			bson.M{"clave": "C846", "materia": "OPTATIVA III", "semestre": 8, "creditos": 7.5},
			bson.M{"clave": "C847", "materia": "REDES DE COMPUTADORAS", "semestre": 8, "creditos": 7.5},
			bson.M{"clave": "C848", "materia": "SISTEMAS DISTRIBUIDOS", "semestre": 8, "creditos": 10.5},
		}

		_, errPlan := coleccionPlan.InsertMany(context.TODO(), materias)
		if errPlan != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error guardando el plan de estudios"})
			return
		}

		// 2. POBLAR EL KARDEX DE PRUEBA
		coleccionKardex := client.Database(os.Getenv("DB_NAME")).Collection("kardex")
		coleccionKardex.Drop(context.TODO())

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

		_, errKardex := coleccionKardex.InsertOne(context.TODO(), miKardex)
		if errKardex != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error guardando el kardex en Mongo"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"exito":                true,
			"mensaje":              "¡Plan de estudios y Kardex insertados exitosamente!",
			"materias_registradas": len(materias),
		})
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

	// Ruta 4: APERTURA DE GRUPOS (Admin) - Con Horario Detallado por Día
	r.POST("/api/grupos", func(c *gin.Context) {
		// Definimos la sub-estructura para los bloques de cada día
		type DetalleHorario struct {
			Dia    string `json:"dia"`
			Bloque string `json:"bloque"`
		}

		// Estructura principal que recibe el JSON del Frontend
		var input struct {
			ClaveMateria string           `json:"clave_materia"`
			NombreGrupo  string           `json:"nombre_grupo"`
			RFCProfesor  string           `json:"rfc_profesor"`
			CupoMaximo   int              `json:"cupo_maximo"`
			Horario      []DetalleHorario `json:"horario"` // Arreglo de objetos {dia, bloque}
		}

		// Validamos que el JSON venga correcto
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"exito": false, "mensaje": "Datos inválidos o incompletos"})
			return
		}

		coleccion := client.Database(os.Getenv("DB_NAME")).Collection("grupos")

		// Armamos el documento final para MongoDB
		nuevoGrupo := bson.M{
			"clave_materia":     input.ClaveMateria,
			"grupo":             input.NombreGrupo,
			"rfc_profesor":      input.RFCProfesor,
			"cupo_maximo":       input.CupoMaximo,
			"cupo_disponible":   input.CupoMaximo, // Inicialmente todo el cupo está libre
			"horario_detallado": input.Horario,    // Guardamos la lista de días y sus bloques
			"alumnos_inscritos": []string{},       // Lista de boletas vacía para empezar
			"fecha_creacion":    time.Now(),
		}

		_, err := coleccion.InsertOne(context.TODO(), nuevoGrupo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"exito": false, "mensaje": "Error al guardar el grupo en Mongo"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"exito":   true,
			"mensaje": "Grupo " + input.NombreGrupo + " aperturado con horario flexible exitosamente",
		})
	})
	// Ruta 5: OBTENER ALUMNOS DE UN GRUPO (Para el Profesor)
	r.GET("/api/grupos/:id_grupo/alumnos", func(c *gin.Context) {
		idGrupo := c.Param("id_grupo")
		coleccion := client.Database(os.Getenv("DB_NAME")).Collection("grupos")

		var grupo bson.M
		err := coleccion.FindOne(context.TODO(), bson.M{"grupo": idGrupo}).Decode(&grupo)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"exito": false, "mensaje": "Grupo no encontrado"})
			return
		}

		// Aquí Go buscaría en la lista de 'alumnos_inscritos' y traería sus nombres
		// Por ahora, devolvemos un éxito para que el Front pueda pintar la tabla
		c.JSON(http.StatusOK, gin.H{
			"exito":   true,
			"grupo":   idGrupo,
			"materia": grupo["clave_materia"],
			"alumnos": []bson.M{
				{"boleta": "2022350438", "nombre": "José Antonio Godoy López", "p1": 0, "p2": 0, "p3": 0},
				{"boleta": "2024601111", "nombre": "Compañero Prueba", "p1": 0, "p2": 0, "p3": 0},
			},
		})
	})

	// Ruta 6: BUSCAR GRUPOS (Para Inscripción)
	r.GET("/api/buscar-grupos", func(c *gin.Context) {
		query := c.Query("q") // Obtenemos el parámetro de búsqueda
		coleccion := client.Database(os.Getenv("DB_NAME")).Collection("grupos")

		// Filtro de "O" (OR): busca si 'grupo', 'clave_materia' o 'rfc_profesor' contienen el texto
		filtro := bson.M{
			"$or": []bson.M{
				{"grupo": bson.M{"$regex": query, "$options": "i"}},
				{"clave_materia": bson.M{"$regex": query, "$options": "i"}},
				{"rfc_profesor": bson.M{"$regex": query, "$options": "i"}},
			},
		}

		cursor, _ := coleccion.Find(context.TODO(), filtro)
		var resultados []bson.M
		cursor.All(context.TODO(), &resultados)

		c.JSON(http.StatusOK, gin.H{"exito": true, "data": resultados})
	})

	// Ruta 7: Calificaciones Actuales (Parciales del semestre en curso)
	r.GET("/api/alumno/:boleta/calificaciones", func(c *gin.Context) {
		boleta := c.Param("boleta")
		coleccion := client.Database(os.Getenv("DB_NAME")).Collection("grupos")

		// Buscamos los grupos donde el alumno está inscrito
		filtro := bson.M{"alumnos_inscritos": boleta}
		cursor, _ := coleccion.Find(context.TODO(), filtro)

		var resultados []bson.M
		cursor.All(context.TODO(), &resultados)

		c.JSON(http.StatusOK, gin.H{"exito": true, "data": resultados})
	})

	// Ruta 8: Horario del Alumno (Cruza grupos inscritos con sus bloques)
	r.GET("/api/alumno/:boleta/horario", func(c *gin.Context) {
		boleta := c.Param("boleta")
		coleccion := client.Database(os.Getenv("DB_NAME")).Collection("grupos")

		filtro := bson.M{"alumnos_inscritos": boleta}
		cursor, _ := coleccion.Find(context.TODO(), filtro)

		var grupos []bson.M
		cursor.All(context.TODO(), &grupos)

		c.JSON(http.StatusOK, gin.H{"exito": true, "materias": grupos})
	})

	// Ruta 9: Horario del Profesor (Filtra por RFC)
	r.GET("/api/profesor/:rfc/horario", func(c *gin.Context) {
		rfc := c.Param("rfc")
		coleccion := client.Database(os.Getenv("DB_NAME")).Collection("grupos")
		filtro := bson.M{"rfc_profesor": rfc}
		cursor, _ := coleccion.Find(context.TODO(), filtro)
		var grupos []bson.M
		cursor.All(context.TODO(), &grupos)
		c.JSON(http.StatusOK, gin.H{"exito": true, "materias": grupos})
	})

	// Ruta 10: Obtener grupos de un profesor (Para llenar los selects)
	r.GET("/api/profesor/:rfc/grupos", func(c *gin.Context) {
		rfc := c.Param("rfc")
		coleccion := client.Database(os.Getenv("DB_NAME")).Collection("grupos")
		cursor, _ := coleccion.Find(context.TODO(), bson.M{"rfc_profesor": rfc})
		var resultados []bson.M
		cursor.All(context.TODO(), &resultados)
		c.JSON(http.StatusOK, gin.H{"exito": true, "data": resultados})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("🚀 Servidor Go escuchando en http://localhost:%s", port)
	r.Run(":" + port)
}
