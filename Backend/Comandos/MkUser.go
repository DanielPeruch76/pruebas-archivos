package Comandos

import (
	"Backend/ParametrosStructs"
	"Backend/Structs"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func MkUser(parametros ParametrosStructs.ParametrosMkUser) {
	fmt.Println("======Start MKUSR======")
	user := parametros.User
	pass := parametros.Pass
	grp := parametros.Grp
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Grp:", grp)

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
	var id_usuario int32 = 0
	existe_grupo := false
	existe_usuario := false

	for _, line := range lines {

		words := strings.Split(line, ",")

		if len(words) == 5 && words[0] != "0" {
			id_usuario += int32(1)
			fmt.Println("Si se encotnro un usuario previo")
			if words[3] == user {
				existe_usuario = true
			}
		}

		if len(words) == 3 {
			if words[2] == grp {
				existe_grupo = true
			}
		}

	}

	if !existe_grupo {
		Structs.TextoEnviar.WriteString("❌ Error: No existe el grupo")
		return
	}

	if existe_usuario {
		Structs.TextoEnviar.WriteString("❌ Error: Ya existe el usuario")
		return
	}

	id_str := strconv.Itoa(int(id_usuario + 1))
	fmt.Println("Id_nuevo_usuario:", id_str)

	data += id_str
	data += ",U,"
	data += grp
	data += ","
	data += user
	data += ","
	data += pass
	data += "\n"

	fmt.Println("Esto es el string completo: ", data)
	fmt.Println("-----------------------------------\n")

	indexBloque := int32(0)

	for _, block := range crrInode.I_block {
		if data == "" {
			continue
		}
		if block != -1 {
			if indexBloque < 13 {
				var crrFileBlock Structs.Fileblock
				if err := Structs.LeerEnDisco(file, &crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Fileblock{})))); err != nil {
					return
				}

				if len(data) < 65 {
					copy(crrFileBlock.B_content[:], data)
					errBloques := Structs.EscribirEnDisco(file, crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Fileblock{}))))
					if errBloques != nil {
						fmt.Println("Error: ", errBloques)
						Structs.TextoEnviar.WriteString("❌ Error: No se puedo actualizar el bloque de archivo")
					}
					data = ""
				} else {
					copy(crrFileBlock.B_content[:], data[:64])
					errBloques := Structs.EscribirEnDisco(file, crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Fileblock{}))))
					if errBloques != nil {
						fmt.Println("Error: ", errBloques)
						Structs.TextoEnviar.WriteString("❌ Error: No se puedo actualizar el bloque de archivo")
					}
					data = data[64:]
				}
			} else {

			}
		}
		indexBloque++
	}

	if data == "" {
		fmt.Println("No hay necesidad de crear un nuevo bloque")
	} else {

		for data != "" {
			fmt.Println(fmt.Sprintf("Esto falta por guardar: %s", data))

			numerador := tempSuperblock.S_first_blo - tempSuperblock.S_block_start
			denominador := int32(binary.Size(Structs.Fileblock{}))
			idNuevoBloque := numerador / denominador
			crrInode.I_block[ObtenerIndiceNuevoBloque(crrInode)] = idNuevoBloque
			tempSuperblock.S_first_blo = tempSuperblock.S_first_blo + int32(binary.Size(Structs.Fileblock{}))

			var nuevoBloque Structs.Fileblock

			if len(data) < 65 {
				copy(nuevoBloque.B_content[:], data)
				errBloques := Structs.EscribirEnDisco(file, nuevoBloque, int64(tempSuperblock.S_block_start+idNuevoBloque*int32(binary.Size(Structs.Fileblock{}))))

				errBitMapSegundoBloque := Structs.EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_block_start+idNuevoBloque))

				if errBitMapSegundoBloque != nil {
					fmt.Println("Error: ", errBitMapSegundoBloque)
					Structs.TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de bloques para el block de archivos")
				}

				if errBloques != nil {
					fmt.Println("Error: ", errBloques)
					Structs.TextoEnviar.WriteString("❌ Error: No se puedo actualizar el bloque de archivo")
				}
				data = ""
			} else {
				copy(nuevoBloque.B_content[:], data[:64])
				errBloques := Structs.EscribirEnDisco(file, nuevoBloque, int64(tempSuperblock.S_block_start+idNuevoBloque*int32(binary.Size(Structs.Fileblock{}))))

				errBitMapSegundoBloque := Structs.EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_block_start+idNuevoBloque))

				if errBitMapSegundoBloque != nil {
					fmt.Println("Error: ", errBitMapSegundoBloque)
					Structs.TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de bloques para el block de archivos")
				}

				if errBloques != nil {
					fmt.Println("Error: ", errBloques)
					Structs.TextoEnviar.WriteString("❌ Error: No se puedo actualizar el bloque de archivo")
				}
				data = data[64:]
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

	Structs.TextoEnviar.WriteString("✅ Se guardo el nuevo usuario con exito\n")

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

	fmt.Println("======End MKUSR======")
}

func ObtenerIndiceNuevoBloque(crrInode Structs.Inode) int32 {
	indexNuevoBloque := int32(0)

	for _, block := range crrInode.I_block {
		if block == -1 {
			if indexNuevoBloque < 13 {
				return indexNuevoBloque
			} else {

			}
		}
		indexNuevoBloque++
	}

	return -1
}
