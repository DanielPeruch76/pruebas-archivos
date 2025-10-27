package Comandos

import (
	"Backend/ParametrosStructs"
	"Backend/Structs"
	"encoding/binary"
	"fmt"
	"strings"
)

func Login(parametros ParametrosStructs.ParametrosLogin) {
	fmt.Println("======Comando Login======")
	user := parametros.User
	pass := parametros.Pass
	id := parametros.Id
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Id:", id)

	if Structs.Usuario.Status {
		fmt.Println("Error: Usuario logueado")
		Structs.TextoEnviar.WriteString("❌ Error: Ya existe una sesión activa\n")
		return
	}

	path, encontrado := BuscarDisco(id)
	var login bool = false

	if !encontrado {
		Structs.TextoEnviar.WriteString("❌ Error: No existe la partición indicada\n")
		return
	}

	file, err := Structs.AbrirArchivo(path)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB

	if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
		Structs.TextoEnviar.WriteString("❌ Error: No se pudo leer el MBR\n")
		return
	}

	//Structs.PrintMBR(TempMBR)

	fmt.Println("-------------")

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
		fmt.Println("Partition not found")
		return
	}

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

	fmt.Println("Fileblock------------")

	lines := strings.Split(data, "\n")

	for _, line := range lines {

		words := strings.Split(line, ",")

		if len(words) == 5 {
			if (strings.Contains(words[3], user)) && (strings.Contains(words[4], pass)) {
				login = true
				break
			}
		}
	}

	fmt.Println("Inode", crrInode.I_block)

	defer file.Close()

	if login {
		fmt.Println("Sesión iniciada :)")
		Structs.TextoEnviar.WriteString(fmt.Sprintf("👨🏻‍💻 El usuario %s ha iniciado sesión\n", user))
		Structs.Usuario.ID = id
		Structs.Usuario.Status = true
		Structs.Usuario.Path = path
		Structs.Usuario.User = user
	} else {
		fmt.Println("Credenciales erróneas")
		Structs.TextoEnviar.WriteString("🚨 Credenciales incorrectas\n")
		Structs.Usuario.ID = ""
		Structs.Usuario.Status = false
		Structs.Usuario.Path = ""
		Structs.Usuario.User = ""
	}

	fmt.Println("======End LOGIN======")
}
