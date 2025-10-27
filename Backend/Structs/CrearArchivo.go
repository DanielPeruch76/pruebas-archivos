package Structs

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

func CrearNuevoArchivo(indiceInodo int32, contenido string, file *os.File, tempSuperblock Superblock, nombreArchivo string, tamañoArchivo int32) Superblock {

	var inodoCarpeta Inode

	if err := LeerEnDisco(file, &inodoCarpeta, int64(tempSuperblock.S_inode_start+indiceInodo*int32(binary.Size(Inode{})))); err != nil {
		return tempSuperblock
	}

	indexBloque := int32(0)

	for _, block := range inodoCarpeta.I_block {
		if block != -1 {
			if indexBloque < 13 {

				var blockCarpetas Folderblock

				if err := LeerEnDisco(file, &blockCarpetas, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
					return tempSuperblock
				}

				indexFolder := 0

				for _, folder := range blockCarpetas.B_content {
					fmt.Println("Folder === Name:", string(folder.B_name[:]), "B_inodo", folder.B_inodo)

					if estaVacio(folder.B_name) {

						numeradorInodo := tempSuperblock.S_fist_ino - tempSuperblock.S_inode_start
						denominadorInodo := int32(binary.Size(Inode{}))
						idNuevoInodo := numeradorInodo / denominadorInodo
						tempSuperblock.S_fist_ino = tempSuperblock.S_fist_ino + int32(binary.Size(Inode{}))

						copy(folder.B_name[:], nombreArchivo)
						folder.B_inodo = idNuevoInodo

						fmt.Printf("Este es el nombre que se guarda la carpeta\" %s\"/n", nombreArchivo)

						copy(blockCarpetas.B_content[indexFolder].B_name[:], nombreArchivo)
						blockCarpetas.B_content[indexFolder].B_inodo = idNuevoInodo

						var nuevoInodo Inode
						nuevoInodo.I_uid = 1
						nuevoInodo.I_gid = 1
						nuevoInodo.I_size = tamañoArchivo
						copy(nuevoInodo.I_atime[:], time.Now().UTC().Format("2006-01-02"))
						copy(nuevoInodo.I_ctime[:], time.Now().UTC().Format("2006-01-02"))
						copy(nuevoInodo.I_mtime[:], time.Now().UTC().Format("2006-01-02"))
						copy(nuevoInodo.I_type[:], "1")
						copy(nuevoInodo.I_perm[:], "664")

						for i := int32(0); i < 15; i++ {
							nuevoInodo.I_block[i] = -1
						}

						numeradorBloqueArchivo := tempSuperblock.S_first_blo - tempSuperblock.S_block_start
						denominadorBloqueArchivo := int32(binary.Size(Folderblock{}))
						idSegundoBloque := numeradorBloqueArchivo / denominadorBloqueArchivo
						nuevoInodo.I_block[0] = idSegundoBloque
						tempSuperblock.S_first_blo = tempSuperblock.S_first_blo + int32(binary.Size(Folderblock{}))

						errBloques := EscribirEnDisco(file, blockCarpetas, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{}))))

						errInodos := EscribirEnDisco(file, nuevoInodo, int64(tempSuperblock.S_inode_start+idNuevoInodo*int32(binary.Size(Inode{}))))

						errBitMapInodo := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_inode_start+idNuevoInodo))

						if errBitMapInodo != nil {
							fmt.Println("Error: ", errBitMapInodo)
							TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de inodos")
							return tempSuperblock
						}

						if errBloques != nil {
							fmt.Println("Error: ", errBloques)
							TextoEnviar.WriteString("❌ Error: No se puedo escribir el bloque actualizado")
							return tempSuperblock
						}

						if errInodos != nil {
							fmt.Println("Error: ", errInodos)
							TextoEnviar.WriteString("❌ Error: No se puedo escribir el inodo de archivos")
							return tempSuperblock
						}

						if len(contenido) < 65 {
							var Fileblock1 Fileblock
							copy(Fileblock1.B_content[:], contenido)
							errArchivo := EscribirEnDisco(file, Fileblock1, int64(tempSuperblock.S_block_start+idSegundoBloque*int32(binary.Size(Fileblock{}))))

							errBitMapSegundoBloque := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_block_start+idSegundoBloque))

							if errArchivo != nil {
								fmt.Println("Error: ", errArchivo)
								TextoEnviar.WriteString("❌ Error: No se puedo escribir el bloque de archivos\n")
								return tempSuperblock
							}

							if errBitMapSegundoBloque != nil {
								fmt.Println("Error: ", errBitMapSegundoBloque)
								TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de bloques para el block de archivos")
								return tempSuperblock
							}

							return tempSuperblock

						} else {
							var Fileblock1 Fileblock
							copy(Fileblock1.B_content[:], contenido[:64])
							errArchivo := EscribirEnDisco(file, Fileblock1, int64(tempSuperblock.S_block_start+idSegundoBloque*int32(binary.Size(Fileblock{}))))

							errBitMapSegundoBloque := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_block_start+idSegundoBloque))

							if errArchivo != nil {
								fmt.Println("Error: ", errArchivo)
								TextoEnviar.WriteString("❌ Error: No se puedo actualizar el bloque de archivo")
								return tempSuperblock
							}

							if errBitMapSegundoBloque != nil {
								fmt.Println("Error: ", errBitMapSegundoBloque)
								TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de bloques para el block de archivos")
								return tempSuperblock
							}

							contenido = contenido[64:]
							return GuardarContenidoFaltante(idNuevoInodo, file, tempSuperblock, contenido)
						}
					} else {
						fmt.Println("Carpeta ocupada\n")
					}
					indexFolder++
				}
			} else {
				fmt.Println("Apuntador indirecto -----------------------------")
				return tempSuperblock
			}
		}
		indexBloque++
	}

	indexBloque = int32(0)

	inodoPadre := int32(0)
	inodoActual := int32(0)
	buscarInfo := true

	for _, block := range inodoCarpeta.I_block {

		if block != -1 && buscarInfo {
			if indexBloque < 13 {
				var crrFolderBlock Folderblock
				if err := LeerEnDisco(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
					return tempSuperblock
				}
				inodoPadre = crrFolderBlock.B_content[0].B_inodo
				inodoActual = crrFolderBlock.B_content[1].B_inodo
				buscarInfo = false
				indexBloque++
				continue
			}
		}

		if block == -1 {
			if indexBloque < 13 {

				numerador := tempSuperblock.S_first_blo - tempSuperblock.S_block_start
				denominador := int32(binary.Size(Folderblock{}))
				idNuevoBloque := numerador / denominador
				tempSuperblock.S_first_blo = tempSuperblock.S_first_blo + int32(binary.Size(Folderblock{}))

				inodoCarpeta.I_block[indexBloque] = idNuevoBloque

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
				copy(nuevoFolderBlock.B_content[2].B_name[:], nombreArchivo)

				var nuevoInodo Inode
				nuevoInodo.I_uid = 1
				nuevoInodo.I_gid = 1
				nuevoInodo.I_size = tamañoArchivo
				copy(nuevoInodo.I_atime[:], time.Now().UTC().Format("2006-01-02"))
				copy(nuevoInodo.I_ctime[:], time.Now().UTC().Format("2006-01-02"))
				copy(nuevoInodo.I_mtime[:], time.Now().UTC().Format("2006-01-02"))
				copy(nuevoInodo.I_type[:], "1")
				copy(nuevoInodo.I_perm[:], "664")

				for i := int32(0); i < 15; i++ {
					nuevoInodo.I_block[i] = -1
				}

				numeradorBloqueArchivo := tempSuperblock.S_first_blo - tempSuperblock.S_block_start
				denominadorBloqueArchivo := int32(binary.Size(Folderblock{}))
				idSegundoBloque := numeradorBloqueArchivo / denominadorBloqueArchivo
				nuevoInodo.I_block[0] = idSegundoBloque
				tempSuperblock.S_first_blo = tempSuperblock.S_first_blo + int32(binary.Size(Folderblock{}))

				errInodoActual := EscribirEnDisco(file, inodoCarpeta, int64(tempSuperblock.S_inode_start+indiceInodo*int32(binary.Size(Inode{}))))

				errBloques := EscribirEnDisco(file, nuevoFolderBlock, int64(tempSuperblock.S_block_start+idNuevoBloque*int32(binary.Size(Folderblock{}))))

				errInodos := EscribirEnDisco(file, nuevoInodo, int64(tempSuperblock.S_inode_start+idNuevoInodo*int32(binary.Size(Inode{}))))

				errBitMapInodo := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_inode_start+idNuevoInodo))

				errBitMapPrimerBloque := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_block_start+idNuevoBloque))

				if errBitMapInodo != nil {
					fmt.Println("Error: ", errBitMapInodo)
					TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de inodos")
					return tempSuperblock
				}

				if errBitMapPrimerBloque != nil {
					fmt.Println("Error: ", errBitMapPrimerBloque)
					TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de bloques para el nuevo bloque de archivos")
					return tempSuperblock
				}

				if errInodoActual != nil {
					fmt.Println("Error: ", errInodoActual)
					TextoEnviar.WriteString("❌ Error: No se puedo escribir el inodo de archivos")
					return tempSuperblock
				}

				if errBloques != nil {
					fmt.Println("Error: ", errBloques)
					TextoEnviar.WriteString("❌ Error: No se puedo escribir el bloque actualizado")
					return tempSuperblock
				}

				if errInodos != nil {
					fmt.Println("Error: ", errInodos)
					TextoEnviar.WriteString("❌ Error: No se puedo escribir el inodo de archivos")
					return tempSuperblock
				}

				if len(contenido) < 65 {
					var Fileblock1 Fileblock
					copy(Fileblock1.B_content[:], contenido)
					errArchivo := EscribirEnDisco(file, Fileblock1, int64(tempSuperblock.S_block_start+idSegundoBloque*int32(binary.Size(Fileblock{}))))

					errBitMapSegundoBloque := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_block_start+idSegundoBloque))

					if errArchivo != nil {
						fmt.Println("Error: ", errArchivo)
						TextoEnviar.WriteString("❌ Error: No se puedo escribir el bloque de archivos\n")
						return tempSuperblock
					}

					if errBitMapSegundoBloque != nil {
						fmt.Println("Error: ", errBitMapSegundoBloque)
						TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de bloques para el block de archivos")
					}

					return tempSuperblock
				} else {
					var Fileblock1 Fileblock
					copy(Fileblock1.B_content[:], contenido[:64])
					errArchivo := EscribirEnDisco(file, Fileblock1, int64(tempSuperblock.S_block_start+idSegundoBloque*int32(binary.Size(Fileblock{}))))

					errBitMapSegundoBloque := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_block_start+idSegundoBloque))

					if errArchivo != nil {
						fmt.Println("Error: ", errArchivo)
						TextoEnviar.WriteString("❌ Error: No se puedo actualizar el bloque de archivo")
						return tempSuperblock
					}

					if errBitMapSegundoBloque != nil {
						fmt.Println("Error: ", errBitMapSegundoBloque)
						TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de bloques para el block de archivos")
					}

					contenido = contenido[64:]
					return GuardarContenidoFaltante(idNuevoInodo, file, tempSuperblock, contenido)
				}
			} else {
				fmt.Println("Apuntador indirecto -----------------------------")
				return tempSuperblock
			}
		}
		indexBloque++
	}

	return tempSuperblock
}

func GuardarContenidoFaltante(indexInodo int32, file *os.File, tempSuperblock Superblock, contenido string) Superblock {

	var inodoArchivo Inode

	if err := LeerEnDisco(file, &inodoArchivo, int64(tempSuperblock.S_inode_start+indexInodo*int32(binary.Size(Inode{})))); err != nil {
		return tempSuperblock
	}

	indexBloque := int32(0)

	for _, block := range inodoArchivo.I_block {
		if block == -1 {
			if indexBloque < 13 {

				numeradorBloqueArchivo := tempSuperblock.S_first_blo - tempSuperblock.S_block_start
				denominadorBloqueArchivo := int32(binary.Size(Folderblock{}))
				idnuevoBloque := numeradorBloqueArchivo / denominadorBloqueArchivo
				inodoArchivo.I_block[indexBloque] = idnuevoBloque
				tempSuperblock.S_first_blo = tempSuperblock.S_first_blo + int32(binary.Size(Folderblock{}))

				errInodos := EscribirEnDisco(file, inodoArchivo, int64(tempSuperblock.S_inode_start+indexInodo*int32(binary.Size(Inode{}))))

				if errInodos != nil {
					fmt.Println("Error: ", errInodos)
					TextoEnviar.WriteString("❌ Error: No se puedo escribir el inodo de archivos")
					return tempSuperblock
				}

				if len(contenido) < 65 {
					var Fileblock1 Fileblock
					copy(Fileblock1.B_content[:], contenido)

					errArchivo := EscribirEnDisco(file, Fileblock1, int64(tempSuperblock.S_block_start+idnuevoBloque*int32(binary.Size(Fileblock{}))))

					errBitMapSegundoBloque := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_block_start+idnuevoBloque))

					if errBitMapSegundoBloque != nil {
						fmt.Println("Error: ", errBitMapSegundoBloque)
						TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de bloques para el block de archivos")
					}

					if errArchivo != nil {
						fmt.Println("Error: ", errArchivo)
						TextoEnviar.WriteString("❌ Error: No se puedo escribir el bloque de archivos\n")
					}
					return tempSuperblock

				} else {
					var Fileblock1 Fileblock
					copy(Fileblock1.B_content[:], contenido[:64])

					errArchivo := EscribirEnDisco(file, Fileblock1, int64(tempSuperblock.S_block_start+idnuevoBloque*int32(binary.Size(Fileblock{}))))

					errBitMapSegundoBloque := EscribirEnDisco(file, byte(1), int64(tempSuperblock.S_bm_block_start+idnuevoBloque))

					if errBitMapSegundoBloque != nil {
						fmt.Println("Error: ", errBitMapSegundoBloque)
						TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de bloques para el block de archivos")
					}

					if errArchivo != nil {
						fmt.Println("Error: ", errArchivo)
						TextoEnviar.WriteString("❌ Error: No se puedo actualizar el bloque de archivo")
					}
					contenido = contenido[64:]

					return GuardarContenidoFaltante(indexInodo, file, tempSuperblock, contenido)
				}

			} else {
				fmt.Println("Apuntandor indirecto------------------------------------------")
				return tempSuperblock
			}
		}
		indexBloque++
	}

	return tempSuperblock

}
