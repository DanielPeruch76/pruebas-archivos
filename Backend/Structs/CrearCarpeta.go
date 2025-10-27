package Structs

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

func CrearCarpeta(Inodo Inode, file *os.File, tempSuperblock Superblock, carpeta string, direccionInodo int32, path []string) (int32, Superblock) {

	index := int32(0)

	for _, block := range Inodo.I_block {
		if block != -1 {
			if index < 13 {

				var crrFolderBlock Folderblock

				if err := LeerEnDisco(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
					return -1, tempSuperblock
				}
				indexFolder := 0
				for _, folder := range crrFolderBlock.B_content {

					fmt.Println("Folder === Name:", string(folder.B_name[:]), "B_inodo", folder.B_inodo)

					if estaVacio(folder.B_name) {

						numeradorInodo := tempSuperblock.S_fist_ino - tempSuperblock.S_inode_start
						denominadorInodo := int32(binary.Size(Inode{}))
						idNuevoInodo := numeradorInodo / denominadorInodo
						tempSuperblock.S_fist_ino = tempSuperblock.S_fist_ino + int32(binary.Size(Inode{}))

						fmt.Printf("Este es el nombre que se guarda la carpeta \"%s\"\nn", carpeta)

						copy(folder.B_name[:], carpeta)
						folder.B_inodo = idNuevoInodo

						copy(crrFolderBlock.B_content[indexFolder].B_name[:], carpeta)
						crrFolderBlock.B_content[indexFolder].B_inodo = idNuevoInodo

						fmt.Printf(" Este el el bloque actualizado Nombre: %s\n Inodo: %d \n", crrFolderBlock.B_content[3].B_name, crrFolderBlock.B_content[3].B_inodo)

						var nuevoInodo Inode
						nuevoInodo.I_uid = 1
						nuevoInodo.I_gid = 1
						nuevoInodo.I_size = 0
						copy(nuevoInodo.I_atime[:], time.Now().UTC().Format("2006-01-02"))
						copy(nuevoInodo.I_ctime[:], time.Now().UTC().Format("2006-01-02"))
						copy(nuevoInodo.I_mtime[:], time.Now().UTC().Format("2006-01-02"))
						copy(nuevoInodo.I_type[:], "0")
						copy(nuevoInodo.I_perm[:], "664")

						for i := int32(0); i < 15; i++ {
							nuevoInodo.I_block[i] = -1
						}

						numerador := tempSuperblock.S_first_blo - tempSuperblock.S_block_start
						denominador := int32(binary.Size(Folderblock{}))
						idNuevoBloque := numerador / denominador
						nuevoInodo.I_block[0] = idNuevoBloque
						tempSuperblock.S_first_blo = tempSuperblock.S_first_blo + int32(binary.Size(Folderblock{}))

						var nuevoFolder Folderblock
						nuevoFolder.B_content[0].B_inodo = direccionInodo
						copy(nuevoFolder.B_content[0].B_name[:], "..")
						nuevoFolder.B_content[1].B_inodo = idNuevoInodo
						copy(nuevoFolder.B_content[1].B_name[:], ".")

						errInodos := EscribirEnDisco(file, nuevoInodo, int64(tempSuperblock.S_inode_start+idNuevoInodo*int32(binary.Size(Inode{}))))

						errBloqueNuevo := EscribirEnDisco(file, nuevoFolder, int64(tempSuperblock.S_block_start+idNuevoBloque*int32(binary.Size(Folderblock{}))))

						errBloqueActual := EscribirEnDisco(file, crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{}))))

						errBitMapInodo := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_inode_start+idNuevoInodo))

						errBitMapBloques := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_block_start+idNuevoBloque))

						if errBloqueNuevo != nil {
							fmt.Println("Error: ", errBloqueNuevo)
							TextoEnviar.WriteString("❌ Error: No se puedo escribir el nuevo bloque de carpetas")
						}

						if errBloqueActual != nil {
							fmt.Println("Error: ", errBloqueActual)
							TextoEnviar.WriteString("❌ Error: No se puedo escribir la actualizacion del bloque actual")
						}

						if errInodos != nil {
							fmt.Println("Error: ", errInodos)
							TextoEnviar.WriteString("❌ Error: No se puedo crear el nuevo Inodo")
						}

						if errBitMapInodo != nil {
							fmt.Println("Error: ", errBitMapInodo)
							TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de inodos")
						}

						if errBitMapBloques != nil {
							fmt.Println("Error: ", errBitMapBloques)
							TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de bloques")
						}

						var crrInode Inode

						if err := LeerEnDisco(file, &crrInode, int64(tempSuperblock.S_inode_start+direccionInodo*int32(binary.Size(Inode{})))); err != nil {
						}

						fmt.Printf("El inodo apunta a: ", crrInode.I_block)

						var FolderPrueba Folderblock

						if err := LeerEnDisco(file, &FolderPrueba, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
							return -1, tempSuperblock
						}

						fmt.Printf("Nombre: %s\n Inodo: %d \n", FolderPrueba.B_content[0].B_name, FolderPrueba.B_content[0].B_inodo)
						fmt.Printf("Nombre: %s\n Inodo: %d \n", FolderPrueba.B_content[1].B_name, FolderPrueba.B_content[1].B_inodo)
						fmt.Printf("Nombre: %s\n Inodo: %d \n", FolderPrueba.B_content[2].B_name, FolderPrueba.B_content[2].B_inodo)
						fmt.Printf("Nombre: %s\n Inodo: %d \n", FolderPrueba.B_content[3].B_name, FolderPrueba.B_content[3].B_inodo)

						if len(path) == 0 {
							fmt.Println("Se creo un nuevo inodo y con esto es suficiente, exito al crear la ruta del archivo")
							return idNuevoInodo, tempSuperblock
						} else {
							carpetaBuscar := path[0]

							if len(path) > 1 {
								path = path[1:]
							} else {
								path = []string{}
							}
							fmt.Println("Se creo un nuevo inodo para crear más subcarpetas")
							return CrearCarpeta(nuevoInodo, file, tempSuperblock, carpetaBuscar, idNuevoInodo, path)
						}
					} else {
						fmt.Printf("Espacio ocupado en carpeta por %s", folder.B_name[:])
					}

					indexFolder++

				}

			} else {

			}
		}
		index++
	}

	fmt.Println("No se encontro espacio en algun bloque de carpetas\n")

	inodoPadre := int32(0)
	inodoActual := int32(0)
	buscarInfo := true

	index2 := int32(0)

	for _, block := range Inodo.I_block {

		if block != -1 && buscarInfo {
			if index2 < 13 {
				var crrFolderBlock Folderblock
				if err := LeerEnDisco(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
					return -1, tempSuperblock
				}
				inodoPadre = crrFolderBlock.B_content[0].B_inodo
				inodoActual = crrFolderBlock.B_content[1].B_inodo
				buscarInfo = false
				index2++
				continue
			} else {

			}
		}

		if block == -1 {
			if index2 < 13 {

				numerador := tempSuperblock.S_first_blo - tempSuperblock.S_block_start
				denominador := int32(binary.Size(Folderblock{}))
				idNuevoBloque := numerador / denominador
				tempSuperblock.S_first_blo = tempSuperblock.S_first_blo + int32(binary.Size(Folderblock{}))

				Inodo.I_block[index2] = idNuevoBloque
				var nuevoFolderBlock Folderblock
				nuevoFolderBlock.B_content[0].B_inodo = inodoPadre
				copy(nuevoFolderBlock.B_content[0].B_name[:], "..")
				nuevoFolderBlock.B_content[1].B_inodo = inodoActual
				copy(nuevoFolderBlock.B_content[1].B_name[:], ".")

				numeradorInodo := tempSuperblock.S_fist_ino - tempSuperblock.S_inode_start
				denominadorInodo := int32(binary.Size(Inode{}))
				idNuevoInodo := numeradorInodo / denominadorInodo
				tempSuperblock.S_fist_ino = tempSuperblock.S_fist_ino + int32(binary.Size(Inode{}))

				nuevoFolderBlock.B_content[2].B_inodo = idNuevoInodo
				copy(nuevoFolderBlock.B_content[2].B_name[:], carpeta)

				var nuevoInodo Inode
				nuevoInodo.I_uid = 1
				nuevoInodo.I_gid = 1
				nuevoInodo.I_size = 0
				copy(nuevoInodo.I_atime[:], time.Now().UTC().Format("2006-01-02"))
				copy(nuevoInodo.I_ctime[:], time.Now().UTC().Format("2006-01-02"))
				copy(nuevoInodo.I_mtime[:], time.Now().UTC().Format("2006-01-02"))
				copy(nuevoInodo.I_type[:], "0")
				copy(nuevoInodo.I_perm[:], "664")

				for i := int32(0); i < 15; i++ {
					nuevoInodo.I_block[i] = -1
				}

				numeradorSegundoBloque := tempSuperblock.S_first_blo - tempSuperblock.S_block_start
				denominadorSegundoBloque := int32(binary.Size(Folderblock{}))
				idSegundoBloque := numeradorSegundoBloque / denominadorSegundoBloque
				nuevoInodo.I_block[0] = idSegundoBloque
				tempSuperblock.S_first_blo = tempSuperblock.S_first_blo + int32(binary.Size(Folderblock{}))

				var nuevoSegundoFolder Folderblock
				nuevoSegundoFolder.B_content[0].B_inodo = direccionInodo
				copy(nuevoSegundoFolder.B_content[0].B_name[:], "..")
				nuevoSegundoFolder.B_content[1].B_inodo = idNuevoInodo
				copy(nuevoSegundoFolder.B_content[1].B_name[:], ".")

				errInodoActual := EscribirEnDisco(file, Inodo, int64(tempSuperblock.S_inode_start+direccionInodo*int32(binary.Size(Inode{}))))

				errPrimerBloque := EscribirEnDisco(file, nuevoFolderBlock, int64(tempSuperblock.S_block_start+idNuevoBloque*int32(binary.Size(Folderblock{}))))

				errInodoNuevo := EscribirEnDisco(file, nuevoInodo, int64(tempSuperblock.S_inode_start+idNuevoInodo*int32(binary.Size(Inode{}))))

				errSegundoBloque := EscribirEnDisco(file, nuevoSegundoFolder, int64(tempSuperblock.S_block_start+idSegundoBloque*int32(binary.Size(Folderblock{}))))

				errBitMapInodo := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_inode_start+idNuevoInodo))

				errBitMapPrimerBloque := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_block_start+idNuevoBloque))

				errBitMapSegundoBloque := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_block_start+idSegundoBloque))

				if errInodoActual != nil {
					fmt.Println("Error: ", errInodoActual)
					TextoEnviar.WriteString("❌ Error: No se puedo actualizar el inodo actual")
				}
				if errPrimerBloque != nil {
					fmt.Println("Error: ", errPrimerBloque)
					TextoEnviar.WriteString("❌ Error: No se puedo escribir el primer bloque a crear")
				}
				if errInodoNuevo != nil {
					fmt.Println("Error: ", errInodoNuevo)
					TextoEnviar.WriteString("❌ Error: No se puedo escribir el nuevo inodo")
				}

				if errSegundoBloque != nil {
					fmt.Println("Error: ", errSegundoBloque)
					TextoEnviar.WriteString("❌ Error: No se puedo escribir el segundo bloque")
				}

				if errBitMapInodo != nil {
					fmt.Println("Error: ", errBitMapInodo)
					TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de inodos")
				}

				if errBitMapPrimerBloque != nil {
					fmt.Println("Error: ", errBitMapPrimerBloque)
					TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de bloques para el primer bloque")
				}

				if errBitMapSegundoBloque != nil {
					fmt.Println("Error: ", errBitMapSegundoBloque)
					TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de bloques para el segundo bloque")
				}

				var crrInode Inode

				if err := LeerEnDisco(file, &crrInode, int64(tempSuperblock.S_inode_start+direccionInodo*int32(binary.Size(Inode{})))); err != nil {
				}

				fmt.Printf("El inodo apunta a: ", crrInode.I_block)

				var FolderPrueba Folderblock

				if err := LeerEnDisco(file, &FolderPrueba, int64(tempSuperblock.S_block_start+idNuevoBloque*int32(binary.Size(Folderblock{})))); err != nil {
					return -1, tempSuperblock
				}

				fmt.Printf("Nombre: %s\n Inodo: %d \n", FolderPrueba.B_content[0].B_name, FolderPrueba.B_content[0].B_inodo)
				fmt.Printf("Nombre: %s\n Inodo: %d \n", FolderPrueba.B_content[1].B_name, FolderPrueba.B_content[1].B_inodo)
				fmt.Printf("Nombre: %s\n Inodo: %d \n", FolderPrueba.B_content[2].B_name, FolderPrueba.B_content[2].B_inodo)
				fmt.Printf("Nombre: %s\n Inodo: %d \n", FolderPrueba.B_content[3].B_name, FolderPrueba.B_content[3].B_inodo)

				if len(path) == 0 {
					fmt.Println("Se creo un nuevo inodo y con est es suficiente, exito al crear la ruta del archivo")
					return idNuevoInodo, tempSuperblock
				} else {
					fmt.Println("Se creo un nuevo inodo para crear más subcarpetas")
					carpetaBuscar := path[0]

					if len(path) > 1 {
						path = path[1:]
					} else {
						path = []string{}
					}
					fmt.Println("Se creo un nuevo inodo para crear más subcarpetas")
					return CrearCarpeta(nuevoInodo, file, tempSuperblock, carpetaBuscar, idNuevoInodo, path)
				}

			} else {

			}
		}
		index2++
	}
	return -1, tempSuperblock
}

func estaVacio(nombre [12]byte) bool {
	fmt.Printf("¿Esta vacio nombre: %s?\n", nombre)

	zeroBytes := [12]byte{}

	if nombre != zeroBytes {
		fmt.Printf(" No ...\n")
		return false
	} else {
		fmt.Printf("Si ... \n")
		return true
	}
}
