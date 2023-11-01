
console.log("asdasd")
const videoList = document.getElementById('videoList');

fetch('http://localhost:8031/api/files')
  .then(response => response.json())
  .then(data => {
    console.log(data.Name)
    data.forEach(video => {
      const listItem = document.createElement('div');
      const videoLink = document.createElement('label');
      videoLink.textContent = video.Name;
      
      // Agregar evento click al elemento <li>
      listItem.addEventListener('click', () => {
        cargar(video.Name);
      });
      
      listItem.appendChild(videoLink);
      videoList.appendChild(listItem);
    });
  })
  .catch(error => console.log('Error:', error));

  function cargar(name) {
    console.log(name)
    const videoSource = document.getElementById('repro');
    videoSource.src = `http://localhost:8031/video?path=${name}`;
    
    // Resto del código de la función "cargar"
  }
  document.addEventListener("DOMContentLoaded", () => {
    const $inputArchivo = document.querySelector("#archivo"),
        $btnEnviar = document.querySelector("#btnEnviar");

        $inputArchivo.addEventListener("change", async () => {
          const archivos = $inputArchivo.files;
          if (archivos.length <= 0) {
            return alert("No hay archivos seleccionados");
          }
        
          const primerArchivo = archivos[0];
          const formdata = new FormData();
          formdata.append("archivo", primerArchivo);
          const nombre = "Parzibyte"; // Dato de tipo cadena para ejemplificar
          formdata.append("nombre", nombre);
          const URL_SERVIDOR = "http://localhost:8031/subida"; // Servidor de Go
        
          try {
            const response = await fetch(URL_SERVIDOR, {
              method: "POST",
              body: formdata,
            });
            const respuesta = await response.text();
            alert("El servidor dijo: " + respuesta);
          } catch (e) {
            alert("Error en el servidor: " + e.message);
          }
        });

});