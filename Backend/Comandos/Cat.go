package Comandos

import (
	"Backend/ParametrosStructs"
	"Backend/Structs"
	"fmt"
	"strings"
)

func Cat(parametros ParametrosStructs.ParametrosCat) {

	fmt.Println("======Start CAT======")
	listaPath := parametros.ListaPath
	fmt.Println("Path: ", listaPath)

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

	for _, path := range listaPath {
		fmt.Println("Path a analizar: ", path)

		contenido := strings.Split(path, "/")
		pathContenido := contenido[1:]

		if len(pathContenido) > 0 {
			contenidoArchivoCopiar, exitoEncontrandoContenido := Structs.ObtenerContenido(pathContenido, file, tempSuperblock)
			if exitoEncontrandoContenido {
				Structs.TextoEnviar.WriteString(fmt.Sprintf("📜 Este es el contendio de %s: \n%s\n", pathContenido[len(pathContenido)-1], contenidoArchivoCopiar))
			} else {
				Structs.TextoEnviar.WriteString(fmt.Sprintf("❌ Error: No se encontro el archivo: %s \n", pathContenido[len(pathContenido)-1]))
			}
		} else {
			Structs.TextoEnviar.WriteString("❌ Error: Parametro invalido")
		}

	}

}
