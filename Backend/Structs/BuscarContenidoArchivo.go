package Structs

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func ObtenerContenido(path []string, file *os.File, tempSuperblock Superblock) (string, bool) {

	var Inode0 Inode

	if err := LeerEnDisco(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return "", false
	}

	return EncontrarInodoTexto(path, file, tempSuperblock, Inode0)
}

func EncontrarInodoTexto(path []string, file *os.File, tempSuperblock Superblock, Inodo Inode) (string, bool) {
	fmt.Println("===== Se inicia busqueda de texto ======\n")
	index := int32(0)

	nombre := path[0]

	if len(path) > 1 {
		path = path[1:]
	} else {
		path = []string{}
	}

	fmt.Println("========== Nombre Buscado:", nombre)

	for _, block := range Inodo.I_block {
		if block != -1 {
			if index < 13 {

				var crrFolderBlock Folderblock

				if err := LeerEnDisco(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
					return "", false
				}

				for _, folder := range crrFolderBlock.B_content {

					fmt.Println("Folder === Name:", string(folder.B_name[:]), "B_inodo", folder.B_inodo)

					if len(nombre) > 12 {
						nombre = nombre[:12]
					}

					if strings.Contains(string(folder.B_name[:]), nombre) {

						fmt.Println("len(path)", len(path), "Nombre", nombre)
						if len(path) == 0 {
							fmt.Println("Folder found======")
							var inodoArchivo Inode
							if err := LeerEnDisco(file, &inodoArchivo, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Inode{})))); err != nil {
								return "", false
							}

							return ObtnerTexto(inodoArchivo, file, tempSuperblock), true
						} else {
							fmt.Println("NextInode======")
							var NextInode Inode

							if err := LeerEnDisco(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Inode{})))); err != nil {
								return "", false
							}
							return EncontrarInodoTexto(path, file, tempSuperblock, NextInode)
						}
					}
				}

			} else {

			}
		}
		index++
	}

	return "", false
}

func ObtnerTexto(Inode Inode, file *os.File, tempSuperblock Superblock) string {

	index := int32(0)

	var content string = ""

	for _, block := range Inode.I_block {
		if block != -1 {
			if index < 13 {

				var crrFileBlock Fileblock

				if err := LeerEnDisco(file, &crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Fileblock{})))); err != nil {
					return ""
				}

				content += string(crrFileBlock.B_content[:encontrarFinValidoCadena(crrFileBlock.B_content)])
				fmt.Printf("Se encontro el siguiente contenido en el block de archivo %s", string(crrFileBlock.B_content[:encontrarFinValidoCadena(crrFileBlock.B_content)]))

			} else {

			}
		}
		index++
	}

	return content
}

func encontrarFinValidoCadena(content [64]byte) int {
	for i := 0; i < 64; i++ {
		if content[i] == 0 {
			return i
		}

	}
	return 64
}
