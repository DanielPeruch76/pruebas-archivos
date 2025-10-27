package main

import (
	utils "Backend/Analizador"
	funcionesComandos "Backend/Comandos"
	"Backend/ParametrosStructs"
	cadena "Backend/Structs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var TextoEnviar string = ""

type Comandos struct {
	Comandos string `json:"comandos,omitempty"`
}

type Login struct {
	Usuario  string `json:"usuario,omitempty"`
	Password string `json:"password,omitempty"`
	Id       string `json:"id,omitempty"`
}

type RespuestaLogin struct {
	Discos []string `json:"discos"`
	MRBs   []MRB    `json:"mrbs"`
}

// Struct específico para la respuesta JSON
type MRB struct {
	MbrSize       int32       `json:"mbrSize"`
	CreationDate  string      `json:"creationDate"`
	Signature     int32       `json:"signature"`
	Fit           string      `json:"fit"`
	Partitions    []Partition `json:"partitions"`
	Letter        string      `json:"letter"`
	NumPartitions int32       `json:"numPartitions"`
}

type Partition struct {
	Status      string `json:"status"`
	Type        string `json:"type"`
	Fit         string `json:"fit"`
	Start       int32  `json:"start"`
	Size        int32  `json:"size"`
	Name        string `json:"name"`
	Correlative int32  `json:"correlative"`
	Id          string `json:"id"`
}

func RecibirComandos(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var comandos Comandos
	err := json.NewDecoder(req.Body).Decode(&comandos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Se recibió la cadena de texto:", comandos.Comandos)

	stringsArray := []string{""}
	if err := cadena.GuardarDiscos(stringsArray); err != nil {
		fmt.Println("Error:", err)
	}

	lineas := strings.Split(comandos.Comandos, "\n")
	for i, linea := range lineas {
		fmt.Println("----------------------------------------------------")
		fmt.Println("Comando #", i, " ", linea)
		fmt.Printf("Este el comando completo: %d: %s\n", i, linea)
		fmt.Printf("Su tamaño es: %d", len(linea))
		comando := strings.TrimSpace(linea)
		comandoParametros := strings.SplitN(comando, " ", 2)
		if len(comandoParametros) > 1 {
			utils.Analizar(comandoParametros[0], comandoParametros[1])
		} else if strings.EqualFold(comandoParametros[0], "mounted") {
			cadena.TextoEnviar.WriteString("-----------------Inicio Mounted---------------\n")
			fmt.Println("Se encontro un comando \"mounted\"")
			if len(cadena.ParticionesMontadas) < 1 {
				cadena.TextoEnviar.WriteString("⚠️ No hay particiones montadas\n")
			} else {
				cadena.TextoEnviar.WriteString("🗃️ Particiones Mondatas\n")
				for i, valor := range cadena.ParticionesMontadas {
					cadena.TextoEnviar.WriteString(fmt.Sprintf("%d. %s\n", i, valor))
				}
			}
			cadena.TextoEnviar.WriteString("-----------------Fin Mounted---------------\n")

		} else if strings.EqualFold(comandoParametros[0], "logout") {
			fmt.Println("Se encontro un comando \"logout\"")
			cadena.TextoEnviar.WriteString("-----------------Inicio LogOut-----------------\n")
			funcionesComandos.LogOut()
			cadena.TextoEnviar.WriteString("-----------------Fin LogOut--------------------\n")
		} else {
			fmt.Println("Posible espacio en blanco o error")
		}
		fmt.Printf("----------------------------------------------------\n\n")
	}
	comandos.Comandos = cadena.TextoEnviar.String()
	cadena.TextoEnviar.Reset()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comandos)
}

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var login Login
	err := json.NewDecoder(req.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Datos de login recibidos - ID: %s, Usuario: %s, Password: %s\n",
		login.Id, login.Usuario, login.Password)

	var parametros = ParametrosStructs.ParametrosLogin{
		Id:   login.Id,
		User: login.Usuario,
		Pass: login.Password,
	}

	funcionesComandos.Login(parametros)

	var discos []string
	var mbrs []MRB

	if cadena.Usuario.Status {
		rutaDiscos, err := cadena.LeerDiscos()
		rutaDiscos = rutaDiscos[1:]
		if err != nil {
			fmt.Println("Error:", err)
			discos = []string{}
		} else {
			for _, disco := range rutaDiscos {
				discos = append(discos, cadena.ObtenerDiscos(disco))
				mbrTemp := cadena.ObtenerMBR(disco)
				var mbr MRB
				mbr.MbrSize = mbrTemp.MbrSize
				mbr.CreationDate = cadena.BytesToString(mbrTemp.CreationDate[:])
				mbr.Signature = mbrTemp.Signature
				mbr.Fit = cadena.BytesToString(mbrTemp.Fit[:])
				mbr.Letter = cadena.BytesToString(mbrTemp.Letter[:])
				mbr.NumPartitions = mbrTemp.NumPartitions
				mbr.Partitions = ObtenerParticiones(disco)
				mbrs = append(mbrs, mbr)
			}
		}

	} else {
		discos = []string{}
	}

	respuestaLogin := RespuestaLogin{
		Discos: discos,
		MRBs:   mbrs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuestaLogin)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/enviarComando/", RecibirComandos).Methods("POST", "OPTIONS")

	router.HandleFunc("/login", LoginHandler).Methods("POST", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	fmt.Println("Servidor escuchando en el puerto 3000")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", handler))
}

func ObtenerParticiones(ruta string) []Partition {
	tempMbr := cadena.ObtenerMBR(ruta)
	var partitions []Partition
	var existeExtendida bool = false
	var indiceParticionExtendida int = 0
	for i := 0; i < 4; i++ {
		if tempMbr.Partitions[i].Size != 0 {
			var partition Partition
			partition.Status = cadena.BytesToString(tempMbr.Partitions[i].Status[:])
			partition.Type = cadena.BytesToString(tempMbr.Partitions[i].Type[:])
			partition.Fit = cadena.BytesToString(tempMbr.Partitions[i].Fit[:])
			partition.Start = tempMbr.Partitions[i].Start
			partition.Size = tempMbr.Partitions[i].Size
			partition.Name = cadena.BytesToString(tempMbr.Partitions[i].Name[:])
			partition.Correlative = tempMbr.Partitions[i].Correlative
			partition.Id = cadena.BytesToString(tempMbr.Partitions[i].Id[:])
			if tempMbr.Partitions[i].Type == [1]byte{'e'} {
				existeExtendida = true
				indiceParticionExtendida = i
			}
			partitions = append(partitions, partition)
		}
	}

	if existeExtendida {
		file, err := cadena.AbrirArchivo(ruta)
		if err != nil {
			fmt.Println("Error en abrir el archivo para EBR------------------------> ", err)
			return partitions
		}

		particionExtendida := tempMbr.Partitions[indiceParticionExtendida]
		var tempEBR cadena.EBR
		if errNuevoEBR := cadena.LeerEnDisco(file, &tempEBR, int64(particionExtendida.Start)); errNuevoEBR != nil {
			return partitions
		}
		fmt.Println("Se imprime el primer ebr")
		funcionesComandos.ImprimirEBR(tempEBR)

		if tempEBR.Part_size == 0 {
			return partitions
		} else {
			var partition Partition
			mountStatus := cadena.BytesToString(tempEBR.Part_mount[:])
			if mountStatus == "1" {
				partition.Status = "1"
				partition.Name = cadena.BytesToString(tempEBR.Part_name[0:4])
				partition.Id = cadena.BytesToString(tempEBR.Part_name[0:4])
			} else {
				partition.Status = "0"
				partition.Name = cadena.BytesToString(tempEBR.Part_name[:])
				partition.Id = cadena.BytesToString(tempEBR.Part_name[:])
			}

			//partition.Status = cadena.BytesToString(tempEBR.Part_mount[:])
			partition.Type = "L"
			partition.Fit = cadena.BytesToString(tempEBR.Part_fit[:])
			partition.Start = tempEBR.Part_start
			partition.Size = tempEBR.Part_size
			partition.Correlative = -1
			partitions = append(partitions, partition)
			terminado := false
			indiceNextEBR := tempEBR.Part_next
			if tempEBR.Part_next == -1 {
				return partitions
			} else {
				for !terminado {
					var temporalEBR cadena.EBR
					if err1 := cadena.LeerEnDisco(file, &temporalEBR, int64(indiceNextEBR)); err1 != nil {
						return partitions
					}
					fmt.Println("Se imprime el EBR encontrado en el for")
					funcionesComandos.ImprimirEBR(temporalEBR)
					if temporalEBR.Part_size == 0 {
						terminado = true
						return partitions
					} else {
						var partitionEBR Partition

						mountStatus2 := cadena.BytesToString(temporalEBR.Part_mount[:])
						if mountStatus2 == "1" {
							partitionEBR.Status = "1"
							partitionEBR.Name = cadena.BytesToString(temporalEBR.Part_name[0:4])
							partitionEBR.Id = cadena.BytesToString(temporalEBR.Part_name[0:4])
						} else {
							partitionEBR.Status = "0"
							partitionEBR.Name = cadena.BytesToString(temporalEBR.Part_name[:])
							partitionEBR.Id = cadena.BytesToString(temporalEBR.Part_name[:])
						}

						//partitionEBR.Status = cadena.BytesToString(temporalEBR.Part_mount[:])
						partitionEBR.Type = "l"
						partitionEBR.Fit = cadena.BytesToString(temporalEBR.Part_fit[:])
						partitionEBR.Start = temporalEBR.Part_start
						partitionEBR.Size = temporalEBR.Part_size
						partitionEBR.Correlative = -1
						partitions = append(partitions, partitionEBR)
						if temporalEBR.Part_next == -1 {
							terminado = true
							return partitions
						} else {
							terminado = false
							indiceNextEBR = temporalEBR.Part_next
						}
					}
				}
			}
		}
	}
	return partitions
}
