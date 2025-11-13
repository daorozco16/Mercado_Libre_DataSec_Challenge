// go run solution_summarizer.go --input articulo.txt --type bullet
// go run solution_summarizer.go -t short articulo.txt
//
// solution_summarizer.go
// Versión de Go: 1.21
//
// Este programa es una CLI (Command Line Interface) para resumir
// el contenido de un archivo de texto usando la API de Hugging Face.
// Usando el modelo facebook/bart-large-cnn
//
// Documentación de la API: https://huggingface.co/docs/api-inference/index
// Documentación de la API: https://huggingface.co/docs/inference-providers/tasks/summarization
// Documentación de Go https://pkg.go.dev/

package main

// Se importan librerías necesarias
import (
	"bytes"         // Para convertir datos a bytes y enviarlos por HTTP
	"encoding/json" // Para convertir estructuras Go a JSON y viceversa
	"flag"          // Para leer los argumentos desde la línea de comandos
	"fmt"           // Para imprimir en la pantalla
	"io"            // Para manejar el ReadAll
	"net/http"      // Para hacer peticiones HTTP
	"os"            // Para manejar errores, salir del programa y leer archivos
)

// Se declaran las variables que se van a usar
var prompt string

// Esto es una variable declrada que significa una lista [] de diccionarios {} con claves de texto {string} y valores de cualquier tipo interface{}
var apiResponse []map[string]interface{}

func main() {
	// 1. LEER ARGUMENTOS DEL 'CLI'
	// Se crean variables que van a guardar los argumentos del flag del CLI
	inputPath := flag.String("input", "", "Ruta del archivo de texto a resumir")
	summaryType := flag.String("type", "short", "Tipo de resumen: short, medium, bullet")

	// Atajo -t usando StringVar ya que se usa una variable ya existente
	flag.StringVar(summaryType, "t", "short", "Tipo de resumen: short, medium, bullet (atajo)")

	// Las descripciones se pueden visualizar corriendo el comando y --help

	// Se usa Parse para leer los argumentos que se pasaron en el flag del CLI antes de poder usarlas
	flag.Parse()

	// Si no hay archivo respondemos con un error (se debe usar * para extraer el valor real)
	if *inputPath == "" {
		fmt.Println("Error: se requiere la ruta del archivo de texto.")
		os.Exit(1) // Se usa Exit para salir inmediatamente del programa y se envia codigo de error 1
	}

	// 2. LEER EL ARCHIVO DE TEXTO
	content, err := os.ReadFile(*inputPath) // content almacena el contenido del archivo en bytes y err almacena errores siempre y cuando existan
	if err != nil {                         // nil es nada o vacio
		fmt.Printf("Error leyendo el archivo: %v\n", err) // Como si fuera un fstring en python, %v sirve para imprimer el valor del error, en un formato estandar
		os.Exit(1)
	}

	// 3. CREAR EL PROMPT SEGUN TIPO
	tipo := *summaryType // Se extrae el valor real del type que se envia en el CLI
	if tipo == "short" {
		prompt = "Por favor, resume el siguiente texto en 1 o 2 frases claras, informativas y precisas:\n\n" + string(content)
	} else if tipo == "medium" {
		prompt = "Genera un resumen en un solo parrafo bien estructurado, incluye las ideas principales y no seas redundante:\n\n" + string(content)
	} else if tipo == "bullet" {
		prompt = "Extrae y presenta las ideas clave del siguiente texto en formato de lista de viñetas, usa frases breves y directas:\n\n" + string(content)
	} else {
		fmt.Println("Tipo de resumen invalido. Disponibles: short, medium, bullet")
		os.Exit(1)
	}

	// 4. PREPARAR PETICION A LA API
	// Se crea una estructura 'HFRequest' simple que se usara para enviar datos a la API
	// Tiene un campo llamado Inputs que sera el texto que queremos resumir
	// Teniendo encuenta que es un POST debemos enviar el request en formato JSON
	// Se configura `json:"inputs"` para que indique como se llamara la clave en el JSON
	type HFRequest struct {
		Inputs string `json:"inputs"`
	}

	apiURL := "https://router.huggingface.co/hf-inference/models/facebook/bart-large-cnn"
	reqBody := HFRequest{Inputs: prompt}

	// Se usa 'Marshal' para convertir la estructura anterion a JSON, para que la API la pueda leer
	jsonData, _ := json.Marshal(reqBody)

	// Se crea una peticion HTTP tipo POST
	//(tipo de peticion, destinatario, contenido de la peticion)
	// Se usa NewBuffer para convertir los bytes del JSON en un formato especifico que lee la peticion HTTP
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creando request: %v\n", err)
		os.Exit(1)
	}

	// Se usa los headers y/o encabezados para agregar informacion
	// Se agrega el token de Hugging Face
	// Pongo de esta manera el Token por que es un entorno local y no supe como pasarla como variable de entorno o realizar un cifrado
	TKAuthor := "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"
	req.Header.Set("Authorization", "Bearer "+TKAuthor) // Tipo de autorizacion que se esta enviando o validando
	req.Header.Set("Content-Type", "application/json")  // Tipo de formato del contenido

	// Se crea un cliente HTTP y se envia la peticion
	client := &http.Client{}    // Esta linea segun documentacion se usa para crear un objeto y pasarle la direccion
	resp, err := client.Do(req) // Do literalmente significa hacer la solicitus pasandole el req que se genero con el NewRequest
	if err != nil {
		fmt.Printf("Error llamando a la API: %v\n", err)
		os.Exit(1)
	}

	// Se cierra la respuesta cuando se termine toda la ejecuacion
	// Se usa Body para extraer todo el cuerpo o contenido
	// Y tambien se usa Close para cerrar el resp y evitar perdida de memoria y no dejar conexion abiertas
	// Se usa 'defer' para agendar el cierre de la conexion una vez finalice la funcion
	// Segun buenas practicas de Go para API se debe escribir esta linea justo despues de enviar la solicitud para garantizar ese cierre y no se olvide escribirla al final del codigo
	defer resp.Body.Close()

	// Se usa StatusOK ya que es el equivalente a 200 y StatusCode es el estado del response
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // Aqui se lee el cuerpo de la respuesta para excepcionar errores
		fmt.Printf("Error de API: %s\n", string(bodyBytes))
		os.Exit(1)
	}

	// 5. LEER LA RESPUESTA DE LA API
	respBody, _ := io.ReadAll(resp.Body)

	// 6. CONVERTIR JSON A TEXTO
	// Unmarshal segun documentacion tomas los datos tipo JSON y los convierte en datos Go, como lo son listas, estructuras etc
	// Se utiliza & para que Unmarshal pueda modificar o escribir dentro de esa variable, con esto se le da la direccion en memoria
	err = json.Unmarshal(respBody, &apiResponse)
	if err != nil {
		fmt.Printf("Error parseando respuesta de la API: %v\n", err)
		os.Exit(1)
	}

	// 7. IMPRIMIR EL RESUMEN
	if len(apiResponse) > 0 {
		// El modelo Summarization de Hugging Face usa "summary_text" como clave dentro del Body
		text, ok := apiResponse[0]["summary_text"].(string)
		if ok {
			fmt.Println(text)
		} else { // Se usa si la clave "summary_text" no esta presente en el Body
			fmt.Println("No se encontró el resumen en la respuesta de la API")
		}
	} else {
		fmt.Println("Respuesta vacía de la API")
	}
}
