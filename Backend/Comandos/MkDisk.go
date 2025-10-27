package Comandos

import (
	"Backend/ParametrosStructs"
	"Backend/Structs"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func Mkdisk(parametros ParametrosStructs.ParametrosMKDisk) {

	size := parametros.Size
	fit := strings.ToLower(parametros.Fit)
	unit := strings.ToLower(parametros.Unit)
	path := parametros.Path

	fmt.Println("======Start MKDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Fit:", fit)
	fmt.Println("Unit:", unit)
	fmt.Println("Path", path)

	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Ajuste no adecuado")
		Structs.TextoEnviar.WriteString("Error: Ajuste no adecuado")
		return
	}

	if size <= 0 {
		fmt.Println("Error: El tamaño debe ser mayor que cero")
		Structs.TextoEnviar.WriteString("Error: El tamaño debe ser mayor que cero")
		return
	}

	if unit != "k" && unit != "m" {
		fmt.Println("Error: Las unidades deben ser k o m")
		Structs.TextoEnviar.WriteString("Error: Las unidades deben ser k o m")
		return
	}
	err := Structs.CrearArchivo(path)
	if err != nil {
		fmt.Println("Error: ", err)
		Structs.TextoEnviar.WriteString("Error abriendo el archivo")
		Structs.TextoEnviar.WriteString(err.Error())
		return
	}

	if unit == "k" {
		size = size * 1024
	} else {
		size = size * 1024 * 1024
	}

	file, err := Structs.AbrirArchivo(path)
	if err != nil {
		Structs.TextoEnviar.WriteString("Error abriendo el archivo")
		return
	}

	zeroBuffer := make([]byte, 1024)

	for i := 0; i < size/1024; i++ {
		err := Structs.EscribirEnDisco(file, zeroBuffer, int64(i*1024))
		if err != nil {
			Structs.TextoEnviar.WriteString("Error creando el archivo")
			return
		}
	}

	rand.Seed(time.Now().UnixNano())
	random_signature := rand.Int31n(100)

	var newMRB Structs.MRB
	newMRB.MbrSize = int32(size)
	newMRB.NumPartitions = 0
	newMRB.Signature = random_signature
	copy(newMRB.Fit[:], fit)
	copy(newMRB.CreationDate[:], time.Now().UTC().Format("2006-01-02"))

	if err := Structs.EscribirEnDisco(file, newMRB, 0); err != nil {
		Structs.TextoEnviar.WriteString("Error escribiendo el MBR")
		return
	}

	var TempMBR Structs.MRB

	if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
		Structs.TextoEnviar.WriteString("Error al leer el MBR")
		return
	}

	//Structs.PrintMBR(TempMBR)

	defer file.Close()

	arrayLeido, err := Structs.LeerDiscos()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Array leído:", arrayLeido)
	}

	if err := Structs.AgregarDisco(path); err != nil {
		fmt.Println("Error:", err)
	}

	arrayLeido2, err := Structs.LeerDiscos()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Array leído:", arrayLeido2)
	}

	fmt.Println("======End MKDISK======")

	Structs.TextoEnviar.WriteString(fmt.Sprintf("✅ Se ha creado el disco con éxito\n"))

}
