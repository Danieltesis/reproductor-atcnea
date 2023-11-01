package storage

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"sort"

	"github.com/Maiki/videolibreria/modelos"
)

func ExploreHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Directorio que deseas explorar
	directory := "./media"

	elements, err := Explore(directory)
	if err != nil {
		http.Error(w, "Error al explorar el directorio", http.StatusInternalServerError)
		return
	}

	// Generar la respuesta JSON
	jsonData, err := json.Marshal(elements)
	if err != nil {
		http.Error(w, "Error al generar la respuesta JSON", http.StatusInternalServerError)
		return
	}

	// Establecer las cabeceras de la respuesta
	w.Header().Set("Content-Type", "application/json")

	// Habilitar el acceso CORS
	w.Header().Set("Access-Control-Allow-Origin", "*") // Reemplaza "*" con el dominio permitido si es necesario

	w.WriteHeader(http.StatusOK)

	// Escribir la respuesta JSON en el cuerpo de la respuesta
	w.Write(jsonData)
}
func Explore(p string) ([]*modelos.Element, error) {
	elements := []*modelos.Element{}
	res, e := os.ReadDir(p)
	if e != nil {
		return nil, e
	}

	for _, r := range res {
		elements = append(elements, &modelos.Element{
			Name: r.Name(),
			Path: path.Join(p, r.Name()),
		})
	}
	sort.Slice(elements, func(i, j int) bool {
		b := elements[i].Name < elements[j].Name
		return b
	})
	return elements, nil
}
