package Comandos

import (
	"Backend/ParametrosStructs"
	"Backend/Structs"
	"fmt"
	"strings"
)

func MkDir(parametros ParametrosStructs.ParametrosMkDir) {

	fmt.Println("======Start MKFILe======")

	path := parametros.Path
	p := parametros.P
	fmt.Println("Path: ", path)
	fmt.Println("P: ", p)

	if !Structs.Usuario.Status {
		fmt.Println("No hay una sesión activa")
		Structs.TextoEnviar.WriteString("⚠️ Error: No hay una sesión activa\n")
		return
	}

	ruta := Structs.Usuario.Path
	id := Structs.Usuario.ID
	file, err := Structs.AbrirArchivo(ruta)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB

	if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
		return
	}

	var index int = -1

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if LimpiarString(strings.TrimRight(string(TempMBR.Partitions[i].Id[:]), "\x00")) == LimpiarString(id) {
				fmt.Println("Partición encontrada")
				if LimpiarString(strings.TrimRight(string(TempMBR.Partitions[i].Status[:]), "\x00")) == "1" {
					fmt.Println("La partición esta montada")
					index = i
				} else {
					fmt.Println("La partición no esta montada")
					return
				}
				break
			}
		}
	}

	if index != -1 {
		Structs.PrintPartition(TempMBR.Partitions[index])
	} else {
		fmt.Println("Partition no fue encontrada")
		return
	}

	fmt.Println("ID:", string(Structs.Usuario.ID[:]))
	fmt.Println("index:", index)

	var tempSuperblock Structs.Superblock

	if err := Structs.LeerEnDisco(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		return
	}

	pasos := strings.Split(path, "/")
	pasosCompletos := pasos[1:]

	if len(pasosCompletos) > 1 {
		fmt.Println("Se encontro una ruta completa ", pasosCompletos)
		nombreCarpeta := pasosCompletos[len(pasosCompletos)-1]
		rutaCarpeta := pasosCompletos[:len(pasosCompletos)-1]

		fmt.Printf("Nombre de la carpeta a crear %s\n", nombreCarpeta)
		fmt.Println("Ruta del archivo ", rutaCarpeta)
		fmt.Println("Cantidad de pasos", len(rutaCarpeta))

		indiceInodo, tempSuperBloqueActualizado := Structs.IniciarBusqueda(rutaCarpeta, file, tempSuperblock, false)

		if indiceInodo != -1 {
			fmt.Printf("Se encontro las carpetas padres de : %s\n", nombreCarpeta)
			indiceNuevoInodoCarpeta, superBloqueActualizado := Structs.IniciarBusqueda(pasosCompletos, file, tempSuperBloqueActualizado, true)

			if indiceNuevoInodoCarpeta != -1 {
				Structs.TextoEnviar.WriteString(fmt.Sprintf("✅ Se ha creado la carpeta %s con éxito \n", nombreCarpeta))
			} else {
				Structs.TextoEnviar.WriteString(fmt.Sprintf("❌ Error: No se pudo crear la carpeta\n"))
			}

			errSuperBloque := Structs.EscribirEnDisco(file, superBloqueActualizado, int64(TempMBR.Partitions[index].Start))
			if errSuperBloque != nil {
				fmt.Println("Error: ", errSuperBloque)
				Structs.TextoEnviar.WriteString("❌ Error: No se pudo escribir el superbloque")
			}

		} else {
			fmt.Printf("No se encontro encontro las carpetas padres de : %s\n", nombreCarpeta)
			if p {
				fmt.Printf("Se proporciono el permiso de crear carpetas padres\n")
				indiceNuevoInodoCarpeta, superBloqueActualizado := Structs.IniciarBusqueda(pasosCompletos, file, tempSuperBloqueActualizado, true)

				if indiceNuevoInodoCarpeta != -1 {
					Structs.TextoEnviar.WriteString(fmt.Sprintf("✅ Se ha creado la carpeta %s con éxito \n", nombreCarpeta))
				} else {
					Structs.TextoEnviar.WriteString(fmt.Sprintf("❌ Error: No se pudo crear la carpeta\n"))
				}

				errSuperBloque := Structs.EscribirEnDisco(file, superBloqueActualizado, int64(TempMBR.Partitions[index].Start))
				if errSuperBloque != nil {
					fmt.Println("Error: ", errSuperBloque)
					Structs.TextoEnviar.WriteString("❌ Error: No se pudo escribir el superbloque")
				}
			} else {
				Structs.TextoEnviar.WriteString(fmt.Sprintf("❌ Error: No existen las carpetas padre\n"))
			}
		}
	} else {
		fmt.Println("Se creara la carpeta en la raiz\n")

		indiceNuevoInodoCarpeta, superBloqueActualizado := Structs.IniciarBusqueda(pasosCompletos, file, tempSuperblock, true)

		if indiceNuevoInodoCarpeta != -1 {
			Structs.TextoEnviar.WriteString(fmt.Sprintf("✅ Se ha creado la carpeta %s con éxito \n", pasosCompletos[0]))
		} else {
			Structs.TextoEnviar.WriteString(fmt.Sprintf("❌ Error: No se pudo crear la carpeta\n"))
		}

		errSuperBloque := Structs.EscribirEnDisco(file, superBloqueActualizado, int64(TempMBR.Partitions[index].Start))
		if errSuperBloque != nil {
			fmt.Println("Error: ", errSuperBloque)
			Structs.TextoEnviar.WriteString("❌ Error: No se pudo escribir el superbloque")
		}
	}
}
