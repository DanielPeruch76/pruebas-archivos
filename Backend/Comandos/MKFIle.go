package Comandos

import (
	"Backend/ParametrosStructs"
	"Backend/Structs"
	"fmt"
	"strconv"
	"strings"
)

func MkFile(parametros ParametrosStructs.ParametrosMkFile) {
	fmt.Println("======Start MKFILe======")

	path := parametros.Path
	r := parametros.R
	cont := parametros.Cont
	size := parametros.Size
	fmt.Println("Path: ", path)
	fmt.Println("R: ", r)
	fmt.Println("Cont: ", cont)
	fmt.Println("Size: ", size)

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
	pasosContenido := []string{}
	contenidoArchivoCopiar := ""
	exitoEncontrandoContenido := false

	if cont != "" {
		contenido := strings.Split(cont, "/")
		pasosContenido = contenido[1:]
	}

	if len(pasosContenido) > 0 {
		contenidoArchivoCopiar, exitoEncontrandoContenido = Structs.ObtenerContenido(pasosContenido, file, tempSuperblock)
	}

	if len(pasosContenido) > 0 && !exitoEncontrandoContenido {
		Structs.TextoEnviar.WriteString("❌ Error: No se encontró el archivo que se buscaba")
		return
	}

	if len(pasosCompletos) > 1 {
		fmt.Println("Se encontro una ruta completa ", pasosCompletos)
		nombreArchivo := pasosCompletos[len(pasosCompletos)-1]
		rutaArchivo := pasosCompletos[:len(pasosCompletos)-1]

		fmt.Printf("Nombre del archivo a crear %s\n", nombreArchivo)
		fmt.Println("Ruta del archivo ", rutaArchivo)
		fmt.Println("Cantidad de pasos", len(pasosCompletos))

		indiceInodo, tempSuperBloqueActualizado := Structs.IniciarBusqueda(rutaArchivo, file, tempSuperblock, r)

		if indiceInodo != -1 {
			fmt.Printf("El indice del inodo es: %d\n", indiceInodo)
			fmt.Println("El superBloque actualizado es:", tempSuperBloqueActualizado)

			contenidoCrear := ""

			if exitoEncontrandoContenido {
				contenidoCrear = contenidoArchivoCopiar
			} else {
				if size > 0 {
					contador := 0
					for i := 0; i < size; i++ {
						contenidoCrear += strconv.Itoa(contador)
						fmt.Println("El contador ", contador)
						contador++
						if contador > 9 {
							contador = 0
						}
					}
				} else {
					contenidoCrear = ""
					if size < 0 {
						Structs.TextoEnviar.WriteString("Error: EL tamaño no puede ser menor a 0\n")
						return
					}
				}
			}
			fmt.Printf("Se guardara el siguiente contenido en el archivo: %s\n", contenidoCrear)
			tempSuperBloqueActualizado = Structs.CrearNuevoArchivo(indiceInodo, contenidoCrear, file, tempSuperBloqueActualizado, nombreArchivo, int32(len(contenidoCrear)))

			errSuperBloque := Structs.EscribirEnDisco(file, tempSuperBloqueActualizado, int64(TempMBR.Partitions[index].Start))
			if errSuperBloque != nil {
				fmt.Println("Error: ", errSuperBloque)
				Structs.TextoEnviar.WriteString("❌ Error: No se pudo escribir el superbloque")
			}

			Structs.TextoEnviar.WriteString(fmt.Sprintf("✅ Se ha creado el archivo %s con exito \n", nombreArchivo))

		} else {
			Structs.TextoEnviar.WriteString("❌ Error: No se encontró la ruta donde se desea crear el archivo")
			return
		}
	} else {
		fmt.Println("Se creara un archivo en la raíz:", pasosCompletos)
		nombreArchivo := pasosCompletos[0]
		fmt.Printf("Nombre del archivo a crear %s\n", nombreArchivo)

		contenidoCrear := ""

		if exitoEncontrandoContenido {
			contenidoCrear = contenidoArchivoCopiar
		} else {
			if size > 0 {
				contador := 0
				for i := 0; i < size; i++ {
					contenidoCrear += string(contador)
					contador++
					if contador > 9 {
						contador = 0
					}
				}
			} else {
				contenidoCrear = ""
			}
		}
		fmt.Printf("Se guardara el siguiente contenido en el archivo: %s\n", contenidoCrear)

		tempSuperBloqueActualizado := Structs.CrearNuevoArchivo(0, contenidoCrear, file, tempSuperblock, nombreArchivo, int32(len(contenidoCrear)))

		errSuperBloque := Structs.EscribirEnDisco(file, tempSuperBloqueActualizado, int64(TempMBR.Partitions[index].Start))
		if errSuperBloque != nil {
			fmt.Println("Error: ", errSuperBloque)
			Structs.TextoEnviar.WriteString("❌ Error: No se pudo escribir el superbloque")
		}

		Structs.TextoEnviar.WriteString(fmt.Sprintf("✅ Se ha creado el archivo %s con exito \n", nombreArchivo))

	}
}
