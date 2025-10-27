package Structs

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func IniciarBusqueda(path []string, file *os.File, tempSuperblock Superblock, r bool) (int32, Superblock) {

	var Inode0 Inode

	if err := LeerEnDisco(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return -1, tempSuperblock
	}

	return EncontrarInodo(path, file, tempSuperblock, r, Inode0, 0)
}

func EncontrarInodo(path []string, file *os.File, tempSuperblock Superblock, r bool, Inodo Inode, id_inodo int32) (int32, Superblock) {

	carpeta := path[0]

	if len(path) > 1 {
		path = path[1:]
	} else {
		path = []string{}
	}

	index := int32(0)

	fmt.Println("========== Se esta buscando la carpeta:", carpeta)

	for _, block := range Inodo.I_block {
		if block != -1 {
			if index < 13 {

				var crrFolderBlock Folderblock

				if err := LeerEnDisco(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
					return -1, tempSuperblock
				}
				for _, folder := range crrFolderBlock.B_content {

					fmt.Println("Folder === Name:", string(folder.B_name[:]), "B_inodo", folder.B_inodo)
					fmt.Println("Carpeta buscada : ", carpeta)

					if len(carpeta) > 12 {
						carpeta = carpeta[:12]
					}

					if strings.Contains(string(folder.B_name[:]), carpeta) {

						fmt.Println("len(path)", len(path), "carpeta", carpeta)
						if len(path) == 0 {
							fmt.Println("Carpeta Encotrada======")
							return folder.B_inodo, tempSuperblock
						} else {
							fmt.Printf("Se encontro la carpeta, se procede a buscar la siguiente carpeta en la ruta======\n")
							var NextInode Inode
							if err := LeerEnDisco(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Inode{})))); err != nil {
								return -1, tempSuperblock
							}
							return EncontrarInodo(path, file, tempSuperblock, r, NextInode, folder.B_inodo)
						}
					}
				}

			} else {

			}
		}
		index++
	}

	if r {
		return CrearCarpeta(Inodo, file, tempSuperblock, carpeta, id_inodo, path)
	} else {
		return -1, tempSuperblock
	}

}
