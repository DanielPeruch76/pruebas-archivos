package Comandos

import (
	"Backend/ParametrosStructs"
	"Backend/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

func MkFs(parametros ParametrosStructs.ParametrosMkfs) {
	fmt.Println("======Start MKFS======")
	id := parametros.Id
	type_ := parametros.Type
	fmt.Println("Id:", id)
	fmt.Println("Type:", type_)

	path, encontrado := BuscarDisco(id)

	if !encontrado {
		Structs.TextoEnviar.WriteString("❌ Error: La particion no ha sido montada\n")
		return
	}

	file, err := Structs.AbrirArchivo(path)
	if err != nil {
		Structs.TextoEnviar.WriteString("❌ Error: No se pudo abrir el archivo del disco\n")
		return
	}

	var TempMBR Structs.MRB

	if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
		Structs.TextoEnviar.WriteString("❌ Error: No al leer el MBR del disco\n")
		return
	}

	var index int = -1

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if LimpiarString(strings.TrimRight(string(TempMBR.Partitions[i].Id[:]), "\x00")) == LimpiarString(id) {
				fmt.Println("Partición Encotrada")
				if LimpiarString(strings.TrimRight(string(TempMBR.Partitions[i].Status[:]), "\x00")) == "1" {
					fmt.Println("Partición es montada")
					index = i
				} else {
					fmt.Println("Partition no esta montada")
					Structs.TextoEnviar.WriteString("❌ Error: Partición no esta montanda\n")
					return
				}
				break
			}
		}
	}

	if index != -1 {
		Structs.PrintPartition(TempMBR.Partitions[index])
	} else {
		fmt.Println("Partition not found")
		Structs.TextoEnviar.WriteString("❌ Error: No se encontro la partición\n")
		return
	}

	numerador := int32(TempMBR.Partitions[index].Size - int32(binary.Size(Structs.Superblock{})))
	denrominador_base := int32(4 + int32(binary.Size(Structs.Inode{})) + 3*int32(binary.Size(Structs.Fileblock{})))
	var temp int32 = 0

	denrominador := denrominador_base + temp
	n := int32(numerador / denrominador)

	fmt.Println("N:", n)

	var newSuperblock Structs.Superblock
	newSuperblock.S_inodes_count = n
	newSuperblock.S_blocks_count = 3 * n

	newSuperblock.S_free_blocks_count = 3 * n
	newSuperblock.S_free_inodes_count = n

	copy(newSuperblock.S_mtime[:], time.Now().UTC().Format("2006-01-02"))
	copy(newSuperblock.S_umtime[:], time.Now().UTC().Format("2006-01-02"))
	newSuperblock.S_mnt_count = 0

	CreateExt2(n, TempMBR.Partitions[index], newSuperblock, time.Now().UTC().Format("2006-01-02"), file)

	defer file.Close()
	Structs.TextoEnviar.WriteString("✅ Se formateo el disco con el sistema EXT2\n")
	fmt.Println("======End MKFS======")
}

func CreateExt2(n int32, partition Structs.Partition, newSuperblock Structs.Superblock, date string, file *os.File) {
	fmt.Println("======Start CREATE EXT2======")
	fmt.Println("N:", n)
	fmt.Println("Superblock:", newSuperblock)
	fmt.Println("Date:", date)

	newSuperblock.S_filesystem_type = 2
	newSuperblock.S_bm_inode_start = partition.Start + int32(binary.Size(Structs.Superblock{}))
	newSuperblock.S_bm_block_start = newSuperblock.S_bm_inode_start + n
	newSuperblock.S_inode_start = newSuperblock.S_bm_block_start + 3*n
	newSuperblock.S_block_start = newSuperblock.S_inode_start + n*int32(binary.Size(Structs.Inode{}))

	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1
	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1

	for i := int32(0); i < n; i++ {
		err := Structs.EscribirEnDisco(file, byte(0), int64(newSuperblock.S_bm_inode_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
			Structs.TextoEnviar.WriteString("❌ Error: No se pudo escribir el bitmap de nodos\n")
		}
	}

	for i := int32(0); i < 3*n; i++ {
		err := Structs.EscribirEnDisco(file, byte(0), int64(newSuperblock.S_bm_block_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
			Structs.TextoEnviar.WriteString("❌ Error: No se pudo escribir el bitmap de bloques\n")
		}
	}

	var newInode Structs.Inode
	for i := int32(0); i < 15; i++ {
		newInode.I_block[i] = -1
	}

	for i := int32(0); i < n; i++ {
		err := Structs.EscribirEnDisco(file, newInode, int64(newSuperblock.S_inode_start+i*int32(binary.Size(Structs.Inode{}))))
		if err != nil {
			fmt.Println("Error: ", err)
			Structs.TextoEnviar.WriteString(fmt.Sprintf("❌ Error: No se pudo escribir el inodo no.%d\n", i))
		}
	}

	var newFileblock Structs.Fileblock
	for i := int32(0); i < 3*n; i++ {
		err := Structs.EscribirEnDisco(file, newFileblock, int64(newSuperblock.S_block_start+i*int32(binary.Size(Structs.Fileblock{}))))
		if err != nil {
			fmt.Println("Error: ", err)
			Structs.TextoEnviar.WriteString(fmt.Sprintf("❌ Error: No se pudo escribir el bloque no.%d\n", i))
		}
	}

	var Inode0 Structs.Inode
	Inode0.I_uid = 1
	Inode0.I_gid = 1
	Inode0.I_size = 0
	copy(Inode0.I_atime[:], date)
	copy(Inode0.I_ctime[:], date)
	copy(Inode0.I_mtime[:], date)
	copy(Inode0.I_type[:], "0")
	copy(Inode0.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode0.I_block[i] = -1
	}

	Inode0.I_block[0] = 0

	var Folderblock0 Structs.Folderblock
	Folderblock0.B_content[0].B_inodo = 0
	copy(Folderblock0.B_content[0].B_name[:], ".")
	Folderblock0.B_content[1].B_inodo = 0
	copy(Folderblock0.B_content[1].B_name[:], "..")
	Folderblock0.B_content[2].B_inodo = 1
	copy(Folderblock0.B_content[2].B_name[:], "users.txt")

	var Inode1 Structs.Inode
	Inode1.I_uid = 1
	Inode1.I_gid = 1
	Inode1.I_size = int32(binary.Size(Structs.Folderblock{}))
	copy(Inode1.I_atime[:], date)
	copy(Inode1.I_ctime[:], date)
	copy(Inode1.I_mtime[:], date)
	copy(Inode1.I_type[:], "1")
	copy(Inode1.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode1.I_block[i] = -1
	}

	Inode1.I_block[0] = 1

	data := "1,G,root\n1,U,root,root,123\n"
	var Fileblock1 Structs.Fileblock
	copy(Fileblock1.B_content[:], data)

	newSuperblock.S_mnt_count = 1
	newSuperblock.S_magic = 0xEF53
	newSuperblock.S_inode_size = int32(binary.Size(Structs.Inode{}))
	newSuperblock.S_block_size = int32(binary.Size(Structs.Fileblock{}))
	newSuperblock.S_fist_ino = int32(newSuperblock.S_inode_start + 2*int32(binary.Size(Structs.Inode{})))
	newSuperblock.S_first_blo = int32(newSuperblock.S_block_start + 2*int32(binary.Size(Structs.Fileblock{})))
	errSuperBloque := Structs.EscribirEnDisco(file, newSuperblock, int64(partition.Start))
	if errSuperBloque != nil {
		fmt.Println("Error: ", errSuperBloque)
		Structs.TextoEnviar.WriteString("❌ Error: No se pudo escribir el superbloque")
	}

	errBitMapInodo := Structs.EscribirEnDisco(file, byte(1), int64(newSuperblock.S_bm_inode_start))
	errBitMapInodo = Structs.EscribirEnDisco(file, byte(1), int64(newSuperblock.S_bm_inode_start+1))
	if errBitMapInodo != nil {
		fmt.Println("Error: ", errBitMapInodo)
		Structs.TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de inodos")
	}

	errBitMapBloques := Structs.EscribirEnDisco(file, byte(1), int64(newSuperblock.S_bm_block_start))
	errBitMapBloques = Structs.EscribirEnDisco(file, byte(1), int64(newSuperblock.S_bm_block_start+1))
	if errBitMapBloques != nil {
		fmt.Println("Error: ", errBitMapBloques)
		Structs.TextoEnviar.WriteString("❌ Error: No se pudo actualizar el bitmap de bloques")
	}

	errInodos := Structs.EscribirEnDisco(file, Inode0, int64(newSuperblock.S_inode_start))
	errInodos = Structs.EscribirEnDisco(file, Inode1, int64(newSuperblock.S_inode_start+int32(binary.Size(Structs.Inode{}))))
	if errInodos != nil {
		fmt.Println("Error: ", errInodos)
		Structs.TextoEnviar.WriteString("❌ Error: No se puedo escribir el inodo0 e inodo1")
	}

	errBloques := Structs.EscribirEnDisco(file, Folderblock0, int64(newSuperblock.S_block_start))
	errBloques = Structs.EscribirEnDisco(file, Fileblock1, int64(newSuperblock.S_block_start+int32(binary.Size(Structs.Fileblock{}))))

	if errBloques != nil {
		fmt.Println("Error: ", errBloques)
		Structs.TextoEnviar.WriteString("❌ Error: No se puedo escribir el bloque0 y el bloque1")
	}

	fmt.Println("======End CREATE EXT2======")
}

func BuscarDisco(id_particion string) (string, bool) {

	discos, err := Structs.LeerDiscos()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Array leído:", discos)
	}

	for i := 0; i < len(discos); i++ {

		if discos[i] == "" {
			continue
		}

		file, err := Structs.AbrirArchivo(discos[i])
		disco := discos[i]

		if err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No se pudo abrir el archivo del disco\n")
			continue
		}

		var TempMBR Structs.MRB

		if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No al leer el MBR del disco\n")
			continue
		}

		for i := 0; i < 4; i++ {
			if TempMBR.Partitions[i].Size != 0 {
				if LimpiarString(strings.TrimRight(string(TempMBR.Partitions[i].Id[:]), "\x00")) == LimpiarString(id_particion) {
					fmt.Println("Partición Encotrada")
					if LimpiarString(strings.TrimRight(string(TempMBR.Partitions[i].Status[:]), "\x00")) == "1" {
						fmt.Println("Partición es montada")
						return disco, true
					} else {
						fmt.Println("Partition no esta montada")
						Structs.TextoEnviar.WriteString("❌ Error: Partición no esta montanda\n")
					}
					break
				}
			}
		}
	}
	return "", false
}
