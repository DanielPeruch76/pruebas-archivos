package Structs

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var TextoEnviar strings.Builder
var ParticionesMontadas []string

func CrearArchivo(name string) error {

	name = strings.Trim(name, ` "'`)
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Err CreateFile dir==", err)
		return err
	}

	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Err CreateFile create==", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

func AbrirArchivo(name string) (*os.File, error) {
	name = strings.Trim(name, ` "'`)
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile==", err)
		return nil, err
	}
	return file, nil
}

func EscribirEnDisco(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WriteObject==", err)
		return err
	}
	return nil
}

func LeerEnDisco(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	return nil
}

func InitSearch(path string, file *os.File, tempSuperblock Superblock) int32 {
	fmt.Println("======Start INITSEARCH======")
	fmt.Println("path:", path)

	TempStepsPath := strings.Split(path, "/")
	StepsPath := TempStepsPath[1:]

	fmt.Println("StepsPath:", StepsPath, "len(StepsPath):", len(StepsPath))
	for _, step := range StepsPath {
		fmt.Println("step:", step)
	}

	var Inode0 Inode

	if err := LeerEnDisco(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return -1
	}

	fmt.Println("======End INITSEARCH======")

	return SarchInodeByPath(StepsPath, Inode0, file, tempSuperblock)
}

func pop(s *[]string) string {
	lastIndex := len(*s) - 1
	last := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return last
}

func SarchInodeByPath(StepsPath []string, Inodo Inode, file *os.File, tempSuperblock Superblock) int32 {
	fmt.Println("======Start SARCHINODEBYPATH======")
	index := int32(0)
	SearchedName := strings.Replace(pop(&StepsPath), " ", "", -1)

	fmt.Println("========== SearchedName:", SearchedName)

	for _, block := range Inodo.I_block {
		if block != -1 {
			if index < 13 {

				var crrFolderBlock Folderblock

				if err := LeerEnDisco(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Folderblock{})))); err != nil {
					return -1
				}

				for _, folder := range crrFolderBlock.B_content {

					fmt.Println("Folder === Name:", string(folder.B_name[:]), "B_inodo", folder.B_inodo)

					if strings.Contains(string(folder.B_name[:]), SearchedName) {

						fmt.Println("len(StepsPath)", len(StepsPath), "StepsPath", StepsPath)
						if len(StepsPath) == 0 {
							fmt.Println("Folder found======")
							return folder.B_inodo
						} else {
							fmt.Println("NextInode======")
							var NextInode Inode

							if err := LeerEnDisco(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Inode{})))); err != nil {
								return -1
							}
							return SarchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
						}
					}
				}

			} else {

			}
		}
		index++
	}

	fmt.Println("======End SARCHINODEBYPATH======")

	return 0
}

func GetInodeFileData(Inode Inode, file *os.File, tempSuperblock Superblock) string {
	fmt.Println("======Start GETINODEFILEDATA======")
	index := int32(0)

	var content string

	for _, block := range Inode.I_block {
		if block != -1 {
			if index < 13 {

				var crrFileBlock Fileblock

				if err := LeerEnDisco(file, &crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Fileblock{})))); err != nil {
					return ""
				}

				content += string(crrFileBlock.B_content[:encontrarFinValido(crrFileBlock.B_content)])

			} else {

			}
		}
		index++
	}

	fmt.Println("======End GETINODEFILEDATA======")
	return content
}

func GuardarDiscos(stringsArray []string) error {

	file, err := AbrirArchivo("./Manejo Discos/Manager.mia")
	if err != nil {
		TextoEnviar.WriteString("❌ Error al abrir archivo manager\n")
		return fmt.Errorf("Error abriendo archivo Manager: %v\n", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(stringsArray); err != nil {
		return fmt.Errorf("error codificando array: %v", err)
	}

	return nil
}

func LeerDiscos() ([]string, error) {

	if _, err := os.Stat("./Manejo Discos/Manager.mia"); os.IsNotExist(err) {
		return []string{}, nil
	}

	file, err := AbrirArchivo("./Manejo Discos/Manager.mia")
	if err != nil {
		return nil, fmt.Errorf("error abriendo archivo: %v", err)
	}
	defer file.Close()

	var stringsArray []string
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&stringsArray); err != nil {
		return nil, fmt.Errorf("error decodificando array: %v", err)
	}

	return stringsArray, nil
}

func AgregarDisco(ruta string) error {
	arrayDiscos, err := LeerDiscos()
	if err != nil {
		return err
	}

	arrayDiscos = append(arrayDiscos, ruta)
	return GuardarDiscos(arrayDiscos)
}

func EliminarString(indice int) error {
	discos, err := LeerDiscos()
	if err != nil {
		return err
	}

	if indice < 0 || indice >= len(discos) {
		return fmt.Errorf("índice fuera de rango")
	}

	discos = append(discos[:indice], discos[indice+1:]...)
	return GuardarDiscos(discos)
}

func encontrarFinValido(content [64]byte) int {
	for i := 0; i < 64; i++ {
		if content[i] == 0 {
			return i
		}
	}
	return 64
}

func ObtenerDiscos(ruta string) string {
	ruta = strings.Trim(ruta, `"'`)
	ruta = strings.TrimSpace(ruta)

	rutaDisco := strings.Split(ruta, `/`)
	nombreDisco := rutaDisco[len(rutaDisco)-1]
	return nombreDisco
}

func ObtenerMBR(ruta string) MRB {
	file, err := AbrirArchivo(ruta)
	if err != nil {
		TextoEnviar.WriteString("Error: No se encotro el disco " + ruta + " " + err.Error() + "\n")
		return MRB{}
	}

	var TempMBR MRB

	if err := LeerEnDisco(file, &TempMBR, 0); err != nil {
		TextoEnviar.WriteString("Error al leer el archivo\n")
		return MRB{}
	}
	return TempMBR
}

func BytesToString(b []byte) string {
	// Encuentra el primer byte nulo y corta hasta ahí
	for i, v := range b {
		if v == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}
