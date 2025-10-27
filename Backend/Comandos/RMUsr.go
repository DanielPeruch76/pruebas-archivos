package Comandos

import (
	"Backend/ParametrosStructs"
	"Backend/Structs"
	"encoding/binary"
	"fmt"
	"strings"
)

func RMUsr(parametros ParametrosStructs.ParametrosRmUser) {
	fmt.Println("======Start RMGRP======")
	user := parametros.User
	fmt.Println("Name:", user)

	if !Structs.Usuario.Status {
		fmt.Println("No hay una sesión activa")
		Structs.TextoEnviar.WriteString("⚠️ Error: No hay una sesión activa\n")
		return
	} else if !(Structs.Usuario.User == "root") {
		Structs.TextoEnviar.WriteString("⚠️ Error: No esta logueado con el usuario \"root\"")
		return
	}

	path := Structs.Usuario.Path
	id := Structs.Usuario.ID
	file, err := Structs.AbrirArchivo(path)
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

	indexInode := Structs.InitSearch("/users.txt", file, tempSuperblock)

	var crrInode Structs.Inode

	if err := Structs.LeerEnDisco(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))); err != nil {
		return
	}

	data := Structs.GetInodeFileData(crrInode, file, tempSuperblock)
	fmt.Println(fmt.Sprintf("Este el contenido del archivo users.txt %s", data))

	lines := strings.Split(data, "\n")
	existe_usuario := false
	usuario_eliminado := false

	var nuevoData string
	var usuarioModificado string

	for _, line := range lines {

		words := strings.Split(line, ",")

		if len(words) == 5 {

			if words[3] == user {

				if words[0] != "0" {
					existe_usuario = true
					usuarioModificado = "0," + words[1] + "," + words[2] + "," + words[3] + "," + words[4]
					nuevoData += usuarioModificado + "\n"
					continue
				} else {
					usuario_eliminado = true
				}

			}
		}

		nuevoData += line
		nuevoData += "\n"
	}

	if usuario_eliminado {
		Structs.TextoEnviar.WriteString("❌ Error: El usuario ya fue eliminado previamente\n")
		return
	}

	if !existe_usuario {
		Structs.TextoEnviar.WriteString("❌ Error: El usuario que desea eliminar no existe\n")
		return
	}

	fmt.Println("Esto es el nuevo string completo: ", nuevoData)
	fmt.Println("-----------------------------------\n")

	indexBloque := int32(0)

	for _, block := range crrInode.I_block {
		if nuevoData == "" {
			continue
		}
		if block != -1 {
			if indexBloque < 13 {
				var crrFileBlock Structs.Fileblock
				if err := Structs.LeerEnDisco(file, &crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Fileblock{})))); err != nil {
					return
				}

				if len(nuevoData) < 65 {
					copy(crrFileBlock.B_content[:], nuevoData)
					errBloques := Structs.EscribirEnDisco(file, crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Fileblock{}))))
					if errBloques != nil {
						fmt.Println("Error: ", errBloques)
						Structs.TextoEnviar.WriteString("❌ Error: No se puedo actualizar el bloque de archivo")
					}
					nuevoData = ""
				} else {
					copy(crrFileBlock.B_content[:], nuevoData[:64])
					errBloques := Structs.EscribirEnDisco(file, crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Fileblock{}))))
					if errBloques != nil {
						fmt.Println("Error: ", errBloques)
						Structs.TextoEnviar.WriteString("❌ Error: No se puedo actualizar el bloque de archivo")
					}
					nuevoData = nuevoData[64:]
				}
			} else {

			}
		}
		indexBloque++
	}

	if nuevoData == "" {
		fmt.Println("No hay necesidad de crear un nuevo bloque")
	} else {

		for nuevoData != "" {
			fmt.Println(fmt.Sprintf("Esto falta por guardar: %s", nuevoData))

			numerador := tempSuperblock.S_first_blo - tempSuperblock.S_block_start
			denominador := int32(binary.Size(Structs.Fileblock{}))
			idNuevoBloque := numerador / denominador
			crrInode.I_block[ObtenerIndiceNuevoBloque(crrInode)] = idNuevoBloque
			tempSuperblock.S_first_blo = tempSuperblock.S_first_blo + int32(binary.Size(Structs.Fileblock{}))

			var nuevoBloque Structs.Fileblock

			if len(nuevoData) < 65 {
				copy(nuevoBloque.B_content[:], nuevoData)
				errBloques := Structs.EscribirEnDisco(file, nuevoBloque, int64(tempSuperblock.S_block_start+idNuevoBloque*int32(binary.Size(Structs.Fileblock{}))))
				if errBloques != nil {
					fmt.Println("Error: ", errBloques)
					Structs.TextoEnviar.WriteString("❌ Error: No se puedo actualizar el bloque de archivo")
				}
				nuevoData = ""
			} else {
				copy(nuevoBloque.B_content[:], nuevoData[:64])
				errBloques := Structs.EscribirEnDisco(file, nuevoBloque, int64(tempSuperblock.S_block_start+idNuevoBloque*int32(binary.Size(Structs.Fileblock{}))))
				if errBloques != nil {
					fmt.Println("Error: ", errBloques)
					Structs.TextoEnviar.WriteString("❌ Error: No se puedo actualizar el bloque de archivo")
				}
				nuevoData = nuevoData[64:]
			}

		}
	}

	errSuperBloque := Structs.EscribirEnDisco(file, tempSuperblock, int64(TempMBR.Partitions[index].Start))
	if errSuperBloque != nil {
		fmt.Println("Error: ", errSuperBloque)
		Structs.TextoEnviar.WriteString("❌ Error: No se pudo escribir el superbloque")
	}

	errInodos := Structs.EscribirEnDisco(file, crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{}))))
	if errInodos != nil {
		fmt.Println("Error: ", errInodos)
		Structs.TextoEnviar.WriteString("❌ Error: No se puedo escribir el inodo0 e inodo1")
	}

	Structs.TextoEnviar.WriteString(" ✅ Se eliminó usuario con exito\n")

	var index1 int = -1

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if LimpiarString(strings.TrimRight(string(TempMBR.Partitions[i].Id[:]), "\x00")) == LimpiarString(id) {
				fmt.Println("Partición encontrada")
				if LimpiarString(strings.TrimRight(string(TempMBR.Partitions[i].Status[:]), "\x00")) == "1" {
					fmt.Println("La partición esta montada")
					index1 = i
				} else {
					fmt.Println("La partición no esta montada")
					return
				}
				break
			}
		}
	}

	if index1 != -1 {
		Structs.PrintPartition(TempMBR.Partitions[index])
	} else {
		fmt.Println("Partition no fue encontrada")
		return
	}

	fmt.Println("ID:", string(Structs.Usuario.ID[:]))
	fmt.Println("index:", index1)

	var tempSuperblock1 Structs.Superblock

	if err := Structs.LeerEnDisco(file, &tempSuperblock1, int64(TempMBR.Partitions[index1].Start)); err != nil {
		return
	}

	indexInode1 := Structs.InitSearch("/users.txt", file, tempSuperblock1)

	var crrInode1 Structs.Inode

	if err := Structs.LeerEnDisco(file, &crrInode1, int64(tempSuperblock1.S_inode_start+indexInode1*int32(binary.Size(Structs.Inode{})))); err != nil {
		return
	}

	data1 := Structs.GetInodeFileData(crrInode1, file, tempSuperblock1)
	fmt.Println(fmt.Sprintf("Este el contenido del archivo users.txt %s", data1))

	fmt.Println("======End RMGRP======")

}
