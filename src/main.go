package main


import(
	"fmt"
	"net/http"
	"io/ioutil"
	//"time"
	"os"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
)

type Data struct{
	Id string `json:"id"`
	Medical_plan string `json:"medical_plan"`
	Dental_plan string `json:"dental_plan"`
	Employee_name string `json:"employee_name"`
}

var tempFile *os.File

//criando http server
func uploadFile(w http.ResponseWriter, r *http.Request){
	fmt.Println("Uploading File")

	//parse input, type multipart/form-data
	r.ParseMultipartForm(10 << 20)// maximum upload 10 MB file size

	//retrieve file from posted form-data
	file, handler, err := r.FormFile("myFile")
	if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
	
	defer file.Close()

	fileName := handler.Filename
	fmt.Printf("Uploaded File: %+v\n", fileName)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("Mime Header: %+v\n", handler.Header)
	
	//write temporary file on the server
	tempFile, err := ioutil.TempFile("temp-files", "upload-*.csv")
	fmt.Println(tempFile)
	fmt.Printf("TempFile: %T\n", tempFile)
	checkError(err)
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	checkError(err)
	tempFile.Write(fileBytes)

	// return success or not
	fmt.Fprintf(w, "Successfully uploaded file\n")

	//testando a conversão aqui dentro
	convert(tempFile.Name())
}

//subindo na porta do localhost e acionando endpoint upload
func setupRoutes(){
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":4001", nil)
}

func main() {
	fmt.Println("Go File upload")
	setupRoutes()
}

//tratamento comum de erro
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

//parse de csv para json
func convert(path string){
	fmt.Printf("TempFile: %T\n", tempFile)
	pwd, _ := os.Getwd()
	fmt.Println(pwd)
	csvFile, err := os.Open(pwd + "/" + path)
	
	fmt.Printf("csvFile: %T\n", csvFile)
	if err != nil{
		fmt.Println("Deu ruim")
		fmt.Println(err.Error())
	}
	reader := csv.NewReader(csvFile)
	var dado []Data
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err !=nil{
			log.Fatal(err)
		}  
		dado = append(dado, Data{
			Id: line[0],
			Medical_plan: line[1],
			Dental_plan: line[2],
			Employee_name: line[3],
		})
	}

	dadoJson, _ := json.Marshal(dado)
	fmt.Println(string(dadoJson))
}

