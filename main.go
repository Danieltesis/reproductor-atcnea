package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Maiki/videolibreria/storage"
)

const (
	DirectorioArchivos        = "media"
	MaximoTamanioFotosEnBytes = 1024 << 20 // 5 megabytes, recuerda que debe haber espacio para tamaño foto + datos adicionales (o sea, formulario)
	HostPermitidoParaCORS     = "http://127.0.0.1:5500/template/bibloteca.html"
)

func main() {

	http.HandleFunc("/subida", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			io.WriteString(w, "Solo se permiten peticiones POST")
			return
		}

		// CORS

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Prevenir que envíen peticiones muy grandes. Recuerda dejar espacio para máximo tamaño de foto + datos adicionales
		r.Body = http.MaxBytesReader(w, r.Body, MaximoTamanioFotosEnBytes)
		// Parsea y aloja en RAM o disco duro, dependiendo del límite que le indiquemos
		err := r.ParseMultipartForm(MaximoTamanioFotosEnBytes)
		if err != nil {
			log.Printf("Error al parsear: %v", err)
			return
		}
		encabezadosDeArchivos := r.MultipartForm.File["archivo"]

		nombre := r.Form.Get("nombre")
		log.Printf("Nombre: %v", nombre)

		encabezadoPrimerArchivo := encabezadosDeArchivos[0]
		primerArchivo, err := encabezadoPrimerArchivo.Open()
		if err != nil {
			log.Printf("Error al abrir archivo: %v", err)
			return
		}
		defer primerArchivo.Close()

		archivoParaGuardar, err := os.Create(filepath.Join(DirectorioArchivos, encabezadoPrimerArchivo.Filename))
		if err != nil {
			log.Printf("Error al crear archivo: %v", err)
			return
		}
		defer archivoParaGuardar.Close()
		_, err = io.Copy(archivoParaGuardar, primerArchivo)
		if err != nil {
			log.Printf("Error al guardar archivo: %v", err)
			return
		}
		io.WriteString(w, "Subido correctamente")
	})

	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir("./media"))))

	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		videoPath := r.URL.Query().Get("path")
		if videoPath == "" {
			http.Error(w, "Se requiere 'path' como parámetro de consulta", http.StatusBadRequest)
			return
		}

		videoFile, err := os.Open(filepath.Join(DirectorioArchivos, videoPath))
		if err != nil {
			http.Error(w, "El archivo de video no se pudo abrir", http.StatusInternalServerError)
			return
		}
		defer videoFile.Close()

		w.Header().Set("Content-Type", "video/mp4")
		_, err = io.Copy(w, videoFile)
		if err != nil {
			http.Error(w, "Error al enviar el video", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/api/files", storage.ExploreHandler)

	http.ListenAndServe(":8031", nil)
}
