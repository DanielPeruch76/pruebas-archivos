package Analizador

import (
	comandos "Backend/Comandos"
	"Backend/ParametrosStructs"
	utils "Backend/Structs"
	"fmt"
	"regexp"
	"strings"
)

func Analizar(comando string, parametros string) {
	fmt.Println("Se comienza a analizar el comando")
	fmt.Printf("El comando es: %s\n", comando)
	fmt.Printf("Los parametros son: %s\n", parametros)

	if comando[0] == '#' {
		fmt.Println("************* Es un comentario *************")
		return
	}

	if strings.EqualFold(comando, "mkdisk") {
		fmt.Println("Se encontró un comando \"mkdisk\"")
		utils.TextoEnviar.WriteString("-------------Inicio Comando MKDisk-----------\n")

		params, sizeEncontrado, pathEncontrado, parametroDesconocido := procesarParametrosMKDisk(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Se encontro un parametro desconocido en mkdisk\n")
			utils.TextoEnviar.WriteString("------------Fin Comando MKDisk-------------\n")
			return
		} else if sizeEncontrado && pathEncontrado {
			comandos.Mkdisk(params)
		} else if !sizeEncontrado && pathEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono el parámetro size\n")
		} else if !pathEncontrado && sizeEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono el parámetro path\n")
		} else {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono ni path ni size\n")
		}
		utils.TextoEnviar.WriteString("------------Fin Comando MKDisk--------------\n")
	} else if strings.EqualFold(comando, "rmdisk") {
		fmt.Println("Se encontro un comando \"rmdisk\"")
		utils.TextoEnviar.WriteString("---------------Inicio Comando RMDisk----------------\n")
		params, pathEncontrado, parametroDesconocido := ProcesarInputRmDisk(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Se encontró un parametro desconocido en rmdisk\n")
			utils.TextoEnviar.WriteString("-----------------Fin Comando RMDiks---------------------\n")
			return
		} else if pathEncontrado {
			comandos.RemoveDisk(params.Path)
		}
		utils.TextoEnviar.WriteString("-----------------Fin Comando RMDiks---------------------\n")

	} else if strings.EqualFold(comando, "fdisk") {
		fmt.Println("Se encontró un comando \"fdisk\"")
		utils.TextoEnviar.WriteString("---------- Inicio Comando FDisk ----------\n")
		params, sizeEncontrado, pathEncontrado, nameEncontrado, parametroDesconocido := procesarParametrosFDisk(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Se encontro un parametro desconocido en fdisk\n")
			utils.TextoEnviar.WriteString("--------- Fin Comando FDiks -------------\n")
			return
		} else if sizeEncontrado && pathEncontrado && nameEncontrado {
			comandos.Fdisk(params)
		} else if !sizeEncontrado && pathEncontrado && nameEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono el parámetro size")
		} else if !pathEncontrado && sizeEncontrado && nameEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono el parámetro path")
		} else if pathEncontrado && sizeEncontrado && !nameEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono el parámetro name")
		} else {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono todos los parametros obligatorios")
		}
		utils.TextoEnviar.WriteString("------------Fin Comando FDiks-----------------\n")
	} else if strings.EqualFold(comando, "mount") {
		fmt.Println("Se encontro un comando \"mount\"")
		utils.TextoEnviar.WriteString("-----------Inicio Comando Mount-----------\n ")
		params, pathEncontrado, nombreEncontrado, parametroDesconocido := ProcesarInputMount(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Se encontro un parametro desconocido en mount\n")
			utils.TextoEnviar.WriteString("-------- Fin Comando Mount -----------\n")
			return
		} else if pathEncontrado && nombreEncontrado {
			comandos.Mount(params)
		} else if !pathEncontrado && nombreEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: path no proporcionado\n")
		} else if !nombreEncontrado && pathEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: name no proporcionado\n")
		} else {
			utils.TextoEnviar.WriteString("❌ Eror: No se proporciono ningun parametro obligatorio\n")
		}
		utils.TextoEnviar.WriteString("--------- Fin Comando Mount --------------\n")
	} else if strings.EqualFold(comando, "mkfs") {
		fmt.Println("Se encontro un comando \"mkfs\"")
		utils.TextoEnviar.WriteString("---------- Inicio Comando MkFs -------------")
		params, idEncontrado, parametroDesconocido := ProcesarInputMkfs(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Se encontro un parametro desconocido en mkfs\n")
			utils.TextoEnviar.WriteString("----------- Fin Comando MkFs -----------")
			return
		} else if idEncontrado {
			comandos.MkFs(params)
		} else {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono el id de la partición\n")
		}
		utils.TextoEnviar.WriteString("----------- Fin Comando MkFs -----------")
	} else if strings.EqualFold(comando, "login") {
		fmt.Println("Se encontro un comando \"login\"")
		utils.TextoEnviar.WriteString("--------- Inicio Comando Login -------\n")
		params, userEncontrado, passEncontrado, idEncontrado, parametroDesconocido := ProcesarInputLogin(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Parametro desconocido en login\n")
			utils.TextoEnviar.WriteString("---------- Fin Comando Login ----------\n")
			return
		} else if userEncontrado && passEncontrado && idEncontrado {
			comandos.Login(params)
		} else if !userEncontrado && passEncontrado && idEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono el parámetro \"user\"\n")
		} else if !passEncontrado && idEncontrado && userEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono el parámetro \"pass\"\n")
		} else if !idEncontrado && passEncontrado && userEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono el parámetro \"id\"\n")
		}
		utils.TextoEnviar.WriteString("------------ Fin Comando Login ----------\n")
	} else if strings.EqualFold(comando, "mkusr") {
		fmt.Println("Se encontro un comando \"mkuser\"")
		utils.TextoEnviar.WriteString("---------- Inicio Comando MkUser -------\n")
		params, userEncontrado, passEncontrado, grpEncontrado, parametroDesconocido := procesarInputMkUser(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Parametro desconocido en comando MKUser\n")
			utils.TextoEnviar.WriteString("------------- Fin Comando MKUser -----------\n")
			return
		} else if userEncontrado && passEncontrado && grpEncontrado {
			comandos.MkUser(params)
		} else if !userEncontrado && passEncontrado && grpEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono el parámetro \"user\"\n")
		} else if !passEncontrado && grpEncontrado && userEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono el parámetro \"pass\"\n")
		} else if !grpEncontrado && passEncontrado && userEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se proporciono el parámetro \"grp\"\n")
		}
		utils.TextoEnviar.WriteString("------------ Fin Comando MKUser -----------\n")
	} else if strings.EqualFold(comando, "mkgrp") {
		fmt.Println("Se encontro un comando \"mkgrp\"")
		utils.TextoEnviar.WriteString("----------- Inicio Comando MkGRP -------\n")
		params, nameEncontrado, parametroDesconocido := procesarInputMkGrp(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Parametro desconocido en comando MKGRP\n")
			utils.TextoEnviar.WriteString("---------- Fin Comando MKGRP ----------\n")
			return
		} else if nameEncontrado {
			comandos.MkGrp(params)
		} else if !nameEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se encontro name para el grupo\n")
			utils.TextoEnviar.WriteString("---------- Fin Comando MKGRP ----------\n")
		}
		utils.TextoEnviar.WriteString("---------- Fin Comando MKGRP ----------\n")
	} else if strings.EqualFold(comando, "rmgrp") {
		fmt.Println("Se encontro un comando \"rmgrp\"")
		utils.TextoEnviar.WriteString("--------- Inicio Comando RMGRP -------\n")
		params, nameEncontrado, parametroDesconocido := procesarInputRMGrp(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Parametro desconocido en comando RMGRP\n")
			utils.TextoEnviar.WriteString("---------- Fin Comando RMGRP ----------\n")
			return
		} else if nameEncontrado {
			comandos.RMGrp(params)
		} else if !nameEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se encontro name para el grupo que se desea eliminar\n")
			utils.TextoEnviar.WriteString("---------- Fin Comando RMGRP -----------\n")
		}
		utils.TextoEnviar.WriteString("---------- Fin Comando RMGRP ---------\n")
	} else if strings.EqualFold(comando, "rmusr") {
		fmt.Println("Se encontro un comando \"rmusr\"")
		utils.TextoEnviar.WriteString("----------- Inicio Comando RMUSR --------\n")
		params, userEncontrado, parametroDesconocido := procesarInputRMUser(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Parametro desconocido en comando RMUSR\n")
			utils.TextoEnviar.WriteString("-------- Fin Comando RMUSR ----------\n")
			return
		} else if userEncontrado {
			comandos.RMUsr(params)
		} else if !userEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se encontro el parametro usuario que se desea eliminar\n")
			utils.TextoEnviar.WriteString("--------- Fin Comando RMUSR ----------\n")
		}
		utils.TextoEnviar.WriteString("---------- Fin Comando RMUSR ----------\n")
	} else if strings.EqualFold(comando, "mkfile") {
		fmt.Println("Se encontro un comando \"mkfile\"")
		utils.TextoEnviar.WriteString("------- Inicio Comando MKFILE --------\n")
		params, pathEncontrado, parametroDesconocido := procesarInputMkFile(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Parametro desconocido en comando MKFILE\n")
			utils.TextoEnviar.WriteString("---------- Fin Comando MKFILE ----------\n")
			return
		} else if pathEncontrado {
			comandos.MkFile(params)
		} else if !pathEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se encontro el parametro path\n")
			utils.TextoEnviar.WriteString("------- Fin Comando MKFILE ---------\n")
		}
		utils.TextoEnviar.WriteString("-------- Fin Comando MKFILE ---------\n")
	} else if strings.EqualFold(comando, "cat") {
		fmt.Println("Se encontro un comando \"cat\"")
		utils.TextoEnviar.WriteString("---------- Inicio Comando CAT -------\n")
		params, pathEncontrado, parametroDesconocido := procesarInputCat(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Parametro desconocido en comando CAT\n")
			utils.TextoEnviar.WriteString("---------- Fin Comando CAT ----------\n")
			return
		} else if pathEncontrado {
			comandos.Cat(params)
		} else if !pathEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se encontro el parametro path\n")
			utils.TextoEnviar.WriteString("---------- Fin Comando CAT ---------\n")
		}
		utils.TextoEnviar.WriteString("-------- Fin Comando CAT --------\n")
	} else if strings.EqualFold(comando, "mkdir") {
		fmt.Println("Se encontro un comando \"mkdir\"")
		utils.TextoEnviar.WriteString("-------- Inicio Comando MKDIR ------\n")
		params, pathEncontrado, parametroDesconocido := procesarInputMkDir(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Parametro desconocido en comando MKDIR\n")
			utils.TextoEnviar.WriteString("--------- Fin Comando MKDIR ----------\n")
			return
		} else if pathEncontrado {
			comandos.MkDir(params)
		} else if !pathEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se encontro el parametro path\n")
			utils.TextoEnviar.WriteString("---------- Fin Comando MKDIR ---------\n")
		}
		utils.TextoEnviar.WriteString("---------- Fin Comando MKDIR ---------\n")
	} else if strings.EqualFold(comando, "chgrp") {
		fmt.Println("Se encontro un comando \"chgrp\"")
		utils.TextoEnviar.WriteString("-------- Inicio Comando CHGRP ------\n")
		params, userEncontrado, grpEncotrado, parametroDesconocido := procesarInputChGrp(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Parametro desconocido en comando CHGRP\n")
			utils.TextoEnviar.WriteString("--------- Fin Comando CHGRP ----------\n")
			return
		} else if userEncontrado && grpEncotrado {
			comandos.ChangeGrp(params)
		} else if !grpEncotrado {
			utils.TextoEnviar.WriteString("❌ Error: No se encontro el parametro grp\n")
			utils.TextoEnviar.WriteString("---------- Fin Comando CHGRP ---------\n")
		} else if !userEncontrado {
			utils.TextoEnviar.WriteString("❌ Error: No se encontro el parametro user\n")
			utils.TextoEnviar.WriteString("---------- Fin Comando CHGRP ---------\n")
		}
		utils.TextoEnviar.WriteString("---------- Fin Comando CHGRP ---------\n")
	} else if strings.EqualFold(comando, "rep") {
		fmt.Println("Se encontro un comando \"rep\"")
		utils.TextoEnviar.WriteString("-------- Inicio Comando rep ------\n")
		params, nameEncontrado, pathEncotrado, idEncotrado, parametroDesconocido := procesarParametrosRep(parametros)
		if parametroDesconocido {
			utils.TextoEnviar.WriteString("❌ Error: Parametro desconocido en comando rep\n")
			utils.TextoEnviar.WriteString("--------- Fin Comando rep ----------\n")
			return
		} else if nameEncontrado && pathEncotrado && idEncotrado {
			comandos.Rep(params)
		} else if !nameEncontrado && pathEncotrado && idEncotrado {
			utils.TextoEnviar.WriteString("❌ Error: No se encontro el parametro name\n")
			utils.TextoEnviar.WriteString("---------- Fin Comando rep ---------\n")
		} else if !pathEncotrado && nameEncontrado && idEncotrado {
			utils.TextoEnviar.WriteString("❌ Error: No se encontro el parametro path\n")
			utils.TextoEnviar.WriteString("---------- Fin Comando rep ---------\n")
		}
		utils.TextoEnviar.WriteString("---------- Fin Comando rep ---------\n")
	}
}

func procesarParametrosMKDisk(parametrosAnalizar string) (ParametrosStructs.ParametrosMKDisk, bool, bool, bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)

	parametros := ParametrosStructs.ParametrosMKDisk{
		Fit:  "ff",
		Unit: "m",
	}

	sizeEncontrado := false
	pathEncontrado := false
	desconocidoEncontrado := false

	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := match[2]
		if nombreParametro != "path" {
			valorParametro = strings.Trim(match[2], "\"")
		}

		switch nombreParametro {
		case "size":
			fmt.Sscanf(valorParametro, "%d", &parametros.Size)
			sizeEncontrado = true
			break
		case "fit":
			parametros.Fit = strings.ToLower(valorParametro)
			break
		case "unit":
			parametros.Unit = strings.ToLower(valorParametro)
			break
		case "path":
			parametros.Path = valorParametro
			pathEncontrado = true
			break
		default:
			desconocidoEncontrado = true
			fmt.Printf("⚠️ Error: Parametro no reconocida '%s'\n", valorParametro)
			return parametros, pathEncontrado, pathEncontrado, desconocidoEncontrado
		}
	}
	return parametros, sizeEncontrado, pathEncontrado, desconocidoEncontrado
}

func ProcesarInputRmDisk(parametrosAnalizar string) (ParametrosStructs.ParametrosRMDisk, bool, bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)
	parametros := ParametrosStructs.ParametrosRMDisk{}
	pathEncontrado := false
	parametroDesconocido := false
	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := match[2]
		if nombreParametro != "path" {
			valorParametro = strings.Trim(match[2], "\"")
		}
		if nombreParametro == "path" {
			parametros.Path = valorParametro
			pathEncontrado = true
		} else {
			parametroDesconocido = true
			fmt.Printf("⚠️ Error: Parametro no reconocida '%s'\n", valorParametro)
			return parametros, pathEncontrado, parametroDesconocido
		}
	}
	return parametros, pathEncontrado, parametroDesconocido
}

func procesarParametrosFDisk(parametrosAnalizar string) (ParametrosStructs.ParametrosFDisk, bool, bool, bool, bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)

	parametros := ParametrosStructs.ParametrosFDisk{
		Unit: "k",
		Type: "p",
		Fit:  "wf",
	}

	sizeEncontrado := false
	pathEncontrado := false
	nombreEncontrado := false
	desconocidoEncontrado := false

	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := match[2]
		if nombreParametro != "path" {
			valorParametro = strings.Trim(match[2], "\"")
		}

		switch nombreParametro {
		case "size":
			fmt.Sscanf(valorParametro, "%d", &parametros.Size)
			sizeEncontrado = true
			break
		case "unit":
			parametros.Unit = strings.ToLower(valorParametro)
			break
		case "path":
			parametros.Path = valorParametro
			pathEncontrado = true
			break
		case "type":
			parametros.Type = strings.ToLower(valorParametro)
			break
		case "fit":
			parametros.Fit = strings.ToLower(valorParametro)
			break
		case "name":
			parametros.Name = valorParametro
			nombreEncontrado = true
			break
		default:
			desconocidoEncontrado = true
			fmt.Printf("⚠️ Error: Parametro no reconocida '%s'\n", valorParametro)
			return parametros, sizeEncontrado, pathEncontrado, nombreEncontrado, desconocidoEncontrado
		}
	}
	return parametros, sizeEncontrado, pathEncontrado, nombreEncontrado, desconocidoEncontrado
}

func ProcesarInputMount(parametrosAnalizar string) (ParametrosStructs.ParametrosMount, bool, bool, bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)
	parametros := ParametrosStructs.ParametrosMount{}
	pathEncontrado := false
	nombreEncontrado := false
	desconocidoEncontrado := false
	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := match[2]
		if nombreParametro != "path" {
			valorParametro = strings.Trim(match[2], "\"")
		}
		if nombreParametro == "path" {
			parametros.Path = valorParametro
			pathEncontrado = true
		} else if nombreParametro == "name" {
			parametros.Name = valorParametro
			nombreEncontrado = true
		} else {
			desconocidoEncontrado = true
			return parametros, pathEncontrado, pathEncontrado, nombreEncontrado
			fmt.Println("Parametro no reconocido en mount")
		}
	}
	return parametros, pathEncontrado, nombreEncontrado, desconocidoEncontrado
}

func ProcesarInputMkfs(parametrosAnalizar string) (ParametrosStructs.ParametrosMkfs, bool, bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)
	parametros := ParametrosStructs.ParametrosMkfs{Type: "full"}
	idEncontrado := false
	desconocidoEncontrado := false

	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := strings.Trim(match[2], "\"")

		if nombreParametro == "id" {
			parametros.Id = valorParametro
			idEncontrado = true
		} else if nombreParametro == "type" {
			parametros.Type = strings.ToLower(valorParametro)
		} else {
			desconocidoEncontrado = true
			fmt.Println("Parametro no reconocido en mkfs")
		}
	}

	return parametros, idEncontrado, desconocidoEncontrado
}

func ProcesarInputLogin(parametrosAnalizar string) (ParametrosStructs.ParametrosLogin, bool, bool, bool, bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)
	parametros := ParametrosStructs.ParametrosLogin{}
	userEncontrado := false
	passEncontrado := false
	idEncontrado := false
	parametroDesconocido := false

	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := strings.Trim(match[2], "\"")

		if nombreParametro == "id" {
			parametros.Id = valorParametro
			idEncontrado = true
		} else if nombreParametro == "user" {
			parametros.User = valorParametro
			userEncontrado = true
		} else if nombreParametro == "pass" {
			parametros.Pass = valorParametro
			passEncontrado = true
		} else {
			parametroDesconocido = true
			fmt.Println("Parametro no reconocido en login")
			return parametros, userEncontrado, passEncontrado, idEncontrado, parametroDesconocido
		}
	}

	return parametros, userEncontrado, passEncontrado, idEncontrado, parametroDesconocido
}

func procesarInputMkUser(parametrosAnalizar string) (ParametrosStructs.ParametrosMkUser, bool, bool, bool, bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)
	parametros := ParametrosStructs.ParametrosMkUser{}
	userEncontrado := false
	passEncontrado := false
	grpEncontrado := false
	parametroDesconocido := false

	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := strings.Trim(match[2], "\"")

		if nombreParametro == "grp" {
			parametros.Grp = valorParametro
			grpEncontrado = true
		} else if nombreParametro == "user" {
			parametros.User = valorParametro
			userEncontrado = true
		} else if nombreParametro == "pass" {
			parametros.Pass = valorParametro
			passEncontrado = true
		} else {
			parametroDesconocido = true
			fmt.Println("Parametro no reconocido en mkuser")
			return parametros, userEncontrado, passEncontrado, grpEncontrado, parametroDesconocido
		}
	}

	return parametros, userEncontrado, passEncontrado, grpEncontrado, parametroDesconocido
}

func procesarInputMkGrp(parametrosAnalizar string) (ParametrosStructs.ParametrosMkGrp, bool, bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)
	parametros := ParametrosStructs.ParametrosMkGrp{}
	nameEncontrado := false
	parametroDesconocido := false

	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := strings.Trim(match[2], "\"")

		if nombreParametro == "name" {
			parametros.Name = valorParametro
			nameEncontrado = true
		} else {
			parametroDesconocido = true
			fmt.Println("Parametro no reconocido en mkgrp")
			return parametros, nameEncontrado, parametroDesconocido
		}
	}

	return parametros, nameEncontrado, parametroDesconocido
}

func procesarInputRMGrp(parametrosAnalizar string) (ParametrosStructs.ParametrosRmGrp, bool, bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)
	parametros := ParametrosStructs.ParametrosRmGrp{}
	nameEncontrado := false
	parametroDesconocido := false

	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := strings.Trim(match[2], "\"")

		if nombreParametro == "name" {
			parametros.Name = valorParametro
			nameEncontrado = true
		} else {
			parametroDesconocido = true
			fmt.Println("Parametro no reconocido en mkgrp")
			return parametros, nameEncontrado, parametroDesconocido
		}
	}

	return parametros, nameEncontrado, parametroDesconocido
}

func procesarInputRMUser(parametrosAnalizar string) (ParametrosStructs.ParametrosRmUser, bool, bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)
	parametros := ParametrosStructs.ParametrosRmUser{}
	userEncontrado := false
	parametroDesconocido := false

	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := strings.Trim(match[2], "\"")

		if nombreParametro == "user" {
			parametros.User = valorParametro
			userEncontrado = true
		} else {
			parametroDesconocido = true
			fmt.Println("Parametro no reconocido en rmuser")
			return parametros, userEncontrado, parametroDesconocido
		}
	}

	return parametros, userEncontrado, parametroDesconocido
}

func procesarInputMkFile(parametrosAnalizar string) (ParametrosStructs.ParametrosMkFile, bool, bool) {
	re := regexp.MustCompile(`-(\w+)(?:=("[^"]+"|\S+))?`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)
	parametros := ParametrosStructs.ParametrosMkFile{
		Cont: "",
		Size: 0,
		R:    false,
	}
	pathEncontrado := false
	parametroDesconocido := false

	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := ""

		if len(match) > 2 && match[2] != "" {
			valorParametro = strings.Trim(match[2], "\"")
		}

		switch nombreParametro {
		case "path":
			if valorParametro == "" {
				fmt.Println("Error: path debe tener un valor")
				parametroDesconocido = true
				return parametros, pathEncontrado, parametroDesconocido
			}
			parametros.Path = valorParametro
			pathEncontrado = true
		case "r":
			fmt.Printf("Flag simple encontrado: -%s\n", nombreParametro)
			parametros.R = true
		case "size":
			fmt.Sscanf(valorParametro, "%d", &parametros.Size)
		case "cont":
			parametros.Cont = valorParametro
		default:
			parametroDesconocido = true
			fmt.Printf("Parámetro no reconocido: -%s\n", nombreParametro)
			return parametros, pathEncontrado, parametroDesconocido
		}
	}
	return parametros, pathEncontrado, parametroDesconocido
}

func procesarInputCat(parametrosAnalizar string) (ParametrosStructs.ParametrosCat, bool, bool) {
	re := regexp.MustCompile(`-\w+\s*=\s*("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)

	parametros := ParametrosStructs.ParametrosCat{
		ListaPath: []string{},
	}
	pathEncontrado := false

	for _, match := range matches {
		if len(match) > 1 && match[1] != "" {
			valorParametro := strings.Trim(match[1], "\"")
			parametros.ListaPath = append(parametros.ListaPath, valorParametro)
			pathEncontrado = true
			fmt.Printf("Ruta Encontrada: %s\n", valorParametro)
		}
	}

	return parametros, pathEncontrado, false
}

func procesarInputMkDir(parametrosAnalizar string) (ParametrosStructs.ParametrosMkDir, bool, bool) {
	re := regexp.MustCompile(`-(\w+)(?:=("[^"]+"|\S+))?`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)
	parametros := ParametrosStructs.ParametrosMkDir{
		P: false,
	}
	pathEncontrado := false
	parametroDesconocido := false

	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := ""

		if len(match) > 2 && match[2] != "" {
			valorParametro = strings.Trim(match[2], "\"")
		}

		switch nombreParametro {
		case "path":
			if valorParametro == "" {
				fmt.Println("Error: path debe tener un valor")
				parametroDesconocido = true
				return parametros, pathEncontrado, parametroDesconocido
			}
			parametros.Path = valorParametro
			pathEncontrado = true
		case "p":
			fmt.Printf("Flag simple encontrado: -%s\n", nombreParametro)
			parametros.P = true
		default:
			parametroDesconocido = true
			fmt.Printf("Parámetro no reconocido: -%s\n", nombreParametro)
			return parametros, pathEncontrado, parametroDesconocido
		}
	}
	return parametros, pathEncontrado, parametroDesconocido
}

func procesarInputChGrp(parametrosAnalizar string) (ParametrosStructs.ParametrosChGrp, bool, bool, bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)
	parametros := ParametrosStructs.ParametrosChGrp{}
	userEncontrado := false
	grpEncontrado := false
	parametroDesconocido := false

	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := strings.Trim(match[2], "\"")

		if nombreParametro == "user" {
			parametros.User = valorParametro
			userEncontrado = true
		} else if nombreParametro == "grp" {
			parametros.Grp = valorParametro
			grpEncontrado = true
		} else {
			parametroDesconocido = true
			fmt.Println("Parametro no reconocido en rmuser")
			return parametros, userEncontrado, grpEncontrado, parametroDesconocido
		}
	}

	return parametros, userEncontrado, grpEncontrado, parametroDesconocido
}

func procesarParametrosRep(parametrosAnalizar string) (ParametrosStructs.ParametrosRep, bool, bool, bool, bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(parametrosAnalizar, -1)

	parametros := ParametrosStructs.ParametrosRep{
		PathFileLs: "",
	}

	idEncontrado := false
	pathEncontrado := false
	nameEncontrado := false
	desconocidoEncontrado := false

	for _, match := range matches {
		nombreParametro := strings.ToLower(match[1])
		valorParametro := match[2]
		if nombreParametro != "path" {
			valorParametro = strings.Trim(match[2], "\"")
		}

		switch nombreParametro {
		case "path_file_ls":
			parametros.PathFileLs = valorParametro
			break
		case "path":
			parametros.Path = valorParametro
			pathEncontrado = true
			break
		case "id":
			parametros.Id = valorParametro
			idEncontrado = true
			break
		case "name":
			parametros.Name = strings.ToLower(valorParametro)
			nameEncontrado = true
			break
		default:
			desconocidoEncontrado = true
			fmt.Printf("⚠️ Error: Parametro no reconocida '%s'\n", valorParametro)
			return parametros, nameEncontrado, pathEncontrado, idEncontrado, desconocidoEncontrado
		}
	}
	return parametros, nameEncontrado, pathEncontrado, idEncontrado, desconocidoEncontrado
}
