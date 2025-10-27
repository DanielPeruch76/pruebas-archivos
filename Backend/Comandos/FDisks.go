package Comandos

import (
	"Backend/ParametrosStructs"
	"Backend/Structs"
	"encoding/binary"
	"fmt"
	"strings"
)

func Fdisk(parametros ParametrosStructs.ParametrosFDisk) {

	size := parametros.Size
	name := parametros.Name
	unit := parametros.Unit
	type_ := parametros.Type
	fit := parametros.Fit
	path := parametros.Path

	fmt.Println("======Start FDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Name:", name)
	fmt.Println("Unit:", unit)
	fmt.Println("Type:", type_)
	fmt.Println("Fit:", fit)
	fmt.Println("Path", path)

	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Ajuste no adecuado")
		Structs.TextoEnviar.WriteString("Error: Ajuste no adecuado\n")
		return
	}

	if size <= 0 {
		fmt.Println("Error: El tamaño debe ser mayor que cero")
		Structs.TextoEnviar.WriteString("Error: El tamaño debe ser mayor que cero\n")
		return
	}

	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: Las unidades deben ser k o m")
		Structs.TextoEnviar.WriteString("Error: Las unidades deben ser k,m o b\n")
		return
	}

	if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Println("Error: El tipo indicado es incorrecto")
		Structs.TextoEnviar.WriteString("Error: Las unidades deben ser k,m o b\n")
		return
	}

	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	} else {
		fmt.Println("El size es en bytes")
	}

	file, err := Structs.AbrirArchivo(path)
	if err != nil {
		Structs.TextoEnviar.WriteString("Error: No se encotro el disco " + path + " " + err.Error() + "\n")
		return
	}

	var TempMBR Structs.MRB

	if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
		Structs.TextoEnviar.WriteString("Error al leer el archivo\n")
		return
	}

	fmt.Println("El MBR actual es este:")
	//Structs.PrintMBR(TempMBR)

	fmt.Println("-------------")

	if int(TempMBR.MbrSize) < size {
		Structs.TextoEnviar.WriteString("Error: El tamaño de la partición es mayor a el disco\n")
		defer file.Close()
		return
	}

	if type_ != "l" {

		var count = 0
		var gap = int32(0)

		for i := 0; i < 4; i++ {
			if TempMBR.Partitions[i].Size != 0 {
				count++
				gap = TempMBR.Partitions[i].Start + TempMBR.Partitions[i].Size
				if TempMBR.Partitions[i].Type == [1]byte{'e'} && strings.EqualFold(type_, "e") {
					Structs.TextoEnviar.WriteString("Error: Ya fue creada una partición extendida previamente\n")
					return
				}
			}
		}

		if count == 4 {
			Structs.TextoEnviar.WriteString("Error: El disco ya tiene 4 particiones")
			return
		}

		for i := 0; i < 4; i++ {
			if TempMBR.Partitions[i].Size == 0 {
				TempMBR.Partitions[i].Size = int32(size)

				if count == 0 {
					TempMBR.Partitions[i].Start = int32(binary.Size(TempMBR))
				} else {
					TempMBR.Partitions[i].Start = gap
				}

				copy(TempMBR.Partitions[i].Name[:], name)
				copy(TempMBR.Partitions[i].Fit[:], fit)
				copy(TempMBR.Partitions[i].Status[:], "0")
				copy(TempMBR.Partitions[i].Type[:], type_)
				TempMBR.Partitions[i].Correlative = -1
				break
			}
		}

		if type_ == "e" {
			fmt.Println("Se creó una partición extendida\n")
			var nuevoEBR Structs.EBR
			copy(nuevoEBR.Part_mount[:], "0")
			copy(nuevoEBR.Part_fit[:], "w")
			nuevoEBR.Part_start = gap + int32(binary.Size(Structs.EBR{}))
			nuevoEBR.Part_size = 0
			nuevoEBR.Part_next = -1
			copy(nuevoEBR.Part_name[:], "")
			if err := Structs.EscribirEnDisco(file, nuevoEBR, int64(gap)); err != nil {
				fmt.Println("Ocurrio un error al escribir el ebt por defecto\n")
				return
			}

			fmt.Println("Se escribio el ebr")

			var tempEBRImprimir Structs.EBR
			if errImprimirEBR := Structs.LeerEnDisco(file, &tempEBRImprimir, int64(gap)); errImprimirEBR != nil {
				fmt.Sprintf("Error al leer el ebr que se crea por defecto,%s", errImprimirEBR.Error)
				Structs.TextoEnviar.WriteString("Error al leer el ebr que se crear por defecto ----------____>")
				return
			}
			fmt.Println("-------->Este es el EBR  que se crea por defecto en la partición extendida")
			ImprimirEBR(tempEBRImprimir)
		}

		if err := Structs.EscribirEnDisco(file, TempMBR, 0); err != nil {
			return
		}

		var TempMBR2 Structs.MRB

		if err := Structs.LeerEnDisco(file, &TempMBR2, 0); err != nil {
			return
		}

		fmt.Println("🛠️ El mbr actualizado,despues de crear la particion:")
		Structs.TextoEnviar.WriteString(fmt.Sprintf("🛠️ Se creó la particion %s\n", name))
		Structs.TextoEnviar.WriteString("Se ha creado la particion\n")
		//Structs.PrintMBR(TempMBR2)

		defer file.Close()
		fmt.Println("======End FDISK======")

	} else {

		count := 0
		existeExtendida := false

		fmt.Println("Se entro a buscar si hay parcion extendida")

		for i := 0; i < 4; i++ {
			if TempMBR.Partitions[i].Size != 0 {
				if TempMBR.Partitions[i].Type == [1]byte{'e'} {
					fmt.Printf("Se encontro una partición extendida\n")
					existeExtendida = true
					break
				}
			}
			fmt.Println(count)
			count++
		}

		fmt.Printf("__________________________--->Indice %d,\n", count)

		if !existeExtendida {
			Structs.TextoEnviar.WriteString("Error: NO existe una partición extendida\n")
			return
		}

		var tempEBR Structs.EBR
		if errNuevoEBR := Structs.LeerEnDisco(file, &tempEBR, int64(TempMBR.Partitions[count].Start)); errNuevoEBR != nil {
			return
		}

		if tempEBR.Part_next == -1 {
			copy(tempEBR.Part_mount[:], "0")
			copy(tempEBR.Part_fit[:], fit)
			tempEBR.Part_size = int32(size)
			tempEBR.Part_next = tempEBR.Part_start + int32(size)
			copy(tempEBR.Part_name[:], name)
			if errActualizarEBR := Structs.EscribirEnDisco(file, tempEBR, int64(TempMBR.Partitions[count].Start)); errActualizarEBR != nil {
				return
			}

			var tempEBRImprimir Structs.EBR
			if errImprimirEBR := Structs.LeerEnDisco(file, &tempEBRImprimir, int64(TempMBR.Partitions[count].Start)); errImprimirEBR != nil {
				Structs.TextoEnviar.WriteString("Error al leer al actualizar el ebr que se creo por defecto\n")
				return
			}
			fmt.Println("------>Se imprime la actulización del ebr que se creo por defecto cuando se guarda la primera partición logica")
			ImprimirEBR(tempEBRImprimir)

			var nuevoEBR Structs.EBR
			copy(nuevoEBR.Part_mount[:], "0")
			copy(nuevoEBR.Part_fit[:], "w")
			nuevoEBR.Part_start = tempEBR.Part_next + int32(binary.Size(Structs.EBR{}))
			nuevoEBR.Part_size = 0
			nuevoEBR.Part_next = -1
			copy(nuevoEBR.Part_name[:], "")
			if errNUEVOEBR := Structs.EscribirEnDisco(file, nuevoEBR, int64(tempEBR.Part_next)); errNUEVOEBR != nil {
				return
			}

			var tempEBRImprimir2 Structs.EBR
			if errImprimirEBR := Structs.LeerEnDisco(file, &tempEBRImprimir2, int64(tempEBR.Part_next)); errImprimirEBR != nil {
				Structs.TextoEnviar.WriteString("Error al leer al actualizar el segundo ebr\n")
				return
			}
			fmt.Println("------>Se imprime el ebr que se creo por segunda vez")
			ImprimirEBR(tempEBRImprimir2)

			Structs.TextoEnviar.WriteString("Se creo con exito la parció logica\n")

		} else {
			terminado := false
			indiceNextEBR := tempEBR.Part_next
			for !terminado {
				var temporalEBR Structs.EBR
				if err1 := Structs.LeerEnDisco(file, &temporalEBR, int64(indiceNextEBR)); err1 != nil {
					return
				}

				if temporalEBR.Part_next == -1 {
					terminado = true

					copy(temporalEBR.Part_mount[:], "0")
					copy(temporalEBR.Part_fit[:], fit)
					temporalEBR.Part_size = int32(size)
					temporalEBR.Part_next = temporalEBR.Part_start + int32(size)
					copy(temporalEBR.Part_name[:], name)
					if errActualizarEBR := Structs.EscribirEnDisco(file, temporalEBR, int64(indiceNextEBR)); errActualizarEBR != nil {
						return
					}

					var tempEBRImprimir Structs.EBR
					if errImprimirEBR := Structs.LeerEnDisco(file, &tempEBRImprimir, int64(indiceNextEBR)); errImprimirEBR != nil {
						Structs.TextoEnviar.WriteString("Error al leer el ebr nuevo que se crea en el for\n")
						return
					}
					fmt.Println("------>Se imprime la actualizacion del ebr que se encontro en el for\n")
					ImprimirEBR(tempEBRImprimir)

					var nuevoEBR Structs.EBR
					copy(nuevoEBR.Part_mount[:], "0")
					copy(nuevoEBR.Part_fit[:], "w")
					nuevoEBR.Part_start = temporalEBR.Part_next + int32(binary.Size(Structs.EBR{}))
					nuevoEBR.Part_size = 0
					nuevoEBR.Part_next = -1
					copy(nuevoEBR.Part_name[:], "")
					if errNUEVOEBR := Structs.EscribirEnDisco(file, nuevoEBR, int64(temporalEBR.Part_next)); errNUEVOEBR != nil {
						return
					}

					var tempEBRImprimir2 Structs.EBR
					if errImprimirEBR := Structs.LeerEnDisco(file, &tempEBRImprimir2, int64(temporalEBR.Part_next)); errImprimirEBR != nil {
						Structs.TextoEnviar.WriteString("Error al leer el ebr nuevo que se crea  por defecto en el for\n")
						return
					}
					fmt.Println("------>Se imprime el ebr que se crea por defecto en el for\n")
					ImprimirEBR(tempEBRImprimir2)

					Structs.TextoEnviar.WriteString("✅Se creo con exito la parció logica\n")

				} else {
					fmt.Println("Se buscara el siguiente EBR\n")
					indiceNextEBR = temporalEBR.Part_next
				}
			}

		}
	}
}
