package Comandos

import (
	"Backend/ParametrosStructs"
	"Backend/Structs"
	"fmt"
	"strconv"
	"strings"
)

var siguienteLetra int = 0

func DeterminarLetra(index int) [1]byte {
	if index < 0 || index > 25 {
		return [1]byte{'?'}
	}
	return [1]byte{byte('A' + index)}
}

func LimpiarString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\x00")
	s = strings.Map(func(r rune) rune {
		if r < 32 || r > 126 {
			return -1
		}
		return r
	}, s)
	return s
}

func Mount(parametros ParametrosStructs.ParametrosMount) {
	fmt.Println("======Start MOUNT======")
	path := parametros.Path
	name := parametros.Name
	name = LimpiarString(name)
	fmt.Println("Path:", path)
	fmt.Println("Name:", name)

	file, err := Structs.AbrirArchivo(path)
	if err != nil {
		Structs.TextoEnviar.WriteString(fmt.Sprintln("Error no se encontró el disco:", err.Error()))
		return
	}

	var TempMBR Structs.MRB

	if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
		Structs.TextoEnviar.WriteString(fmt.Sprintln("Error en leer el archivo:", err.Error()))
		return
	}

	fmt.Printf("MBR actual\n")

	//Structs.PrintMBR(TempMBR)

	fmt.Println("-------------")

	var index int = -1

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			fmt.Println(string(TempMBR.Partitions[i].Name[:]))
			fmt.Println(name)
			if LimpiarString(strings.TrimRight(string(TempMBR.Partitions[i].Name[:]), "\x00")) == name {
				fmt.Println("Partition found")
				index = i
				break
			}
		}
	}

	if index != -1 {
		fmt.Println("Particion encontrada")
		Structs.PrintPartition(TempMBR.Partitions[index])

		if LimpiarString(strings.TrimRight(string(TempMBR.Partitions[index].Status[:]), "\x00")) == "1" {
			Structs.TextoEnviar.WriteString("Error: Ya esta montada la particion")
			return
		}

	} else {
		fmt.Println("Partition not found")
		fmt.Println("Se buscara si es una particion logica\n")
		index2 := -1
		for i := 0; i < 4; i++ {
			if TempMBR.Partitions[i].Size != 0 {
				if TempMBR.Partitions[i].Type == [1]byte{'e'} {
					fmt.Println("Partition extendida encontrada\n")
					index2 = i
					break
				}
			}
		}

		if index2 != -1 {
			fmt.Println("Particion extendida")
			Structs.PrintPartition(TempMBR.Partitions[index2])

			if TempMBR.Letter[0] == 0 {
				TempMBR.Letter = DeterminarLetra(siguienteLetra)
				siguienteLetra++
			}

			TempMBR.NumPartitions++
			id := strconv.Itoa(69) + strconv.Itoa(int(TempMBR.NumPartitions)) + string(TempMBR.Letter[0])
			Structs.ParticionesMontadas = append(Structs.ParticionesMontadas, id)

			indexEBR := TempMBR.Partitions[index2].Start

			paticionMontada := false

			for !paticionMontada {
				var tempEBR Structs.EBR
				if errNuevoEBR := Structs.LeerEnDisco(file, &tempEBR, int64(indexEBR)); errNuevoEBR != nil {
					return
				}

				if LimpiarString(strings.TrimRight(string(tempEBR.Part_name[:]), "\x00")) == name {

					copy(tempEBR.Part_mount[:], "1")
					copy(tempEBR.Part_name[:], id)

					if errActualizarEBR := Structs.EscribirEnDisco(file, tempEBR, int64(indexEBR)); errActualizarEBR != nil {
						return
					}

					var tempEBRImprimir Structs.EBR
					if errImprimirEBR := Structs.LeerEnDisco(file, &tempEBRImprimir, int64(indexEBR)); errImprimirEBR != nil {
						Structs.TextoEnviar.WriteString("Error al leer el ebr guardado ----------____>")
						return
					}

					ImprimirEBR(tempEBRImprimir)

					paticionMontada = true

				} else {
					if tempEBR.Part_next != -1 {
						indexEBR = tempEBR.Part_next
					} else {
						break
					}
				}
			}

			if paticionMontada {
				Structs.TextoEnviar.WriteString(fmt.Sprintln("✅Se monto la particion %s\n", name))
				if err := Structs.EscribirEnDisco(file, TempMBR, 0); err != nil {
					Structs.TextoEnviar.WriteString(fmt.Sprintln("Error al montar la particion:", err.Error()))
					return
				}

				defer file.Close()
				Structs.TextoEnviar.WriteString(fmt.Sprintf("🔨 Se monto la particion %s\n", id))
				fmt.Println("======End MOUNT======")
				return
			} else {
				Structs.TextoEnviar.WriteString(fmt.Sprintln("Erro: No se encontró la particion lógica %s\n", name))
				return
			}

		} else {
			Structs.TextoEnviar.WriteString("Error: No exite parcion extendida para buscar una logica ni parcion primaria de este nombre\n")
			return
		}

	}

	if TempMBR.Letter[0] == 0 {
		TempMBR.Letter = DeterminarLetra(siguienteLetra)
		siguienteLetra++
	}

	TempMBR.NumPartitions++
	id := strconv.Itoa(69) + strconv.Itoa(int(TempMBR.NumPartitions)) + string(TempMBR.Letter[0])
	Structs.ParticionesMontadas = append(Structs.ParticionesMontadas, id)

	copy(TempMBR.Partitions[index].Id[:], id)
	copy(TempMBR.Partitions[index].Status[:], "1")

	if err := Structs.EscribirEnDisco(file, TempMBR, 0); err != nil {
		Structs.TextoEnviar.WriteString(fmt.Sprintln("Error al montar la particion:", err.Error()))
		return
	}

	var TempMBR2 Structs.MRB
	if err := Structs.LeerEnDisco(file, &TempMBR2, 0); err != nil {
		print("Error al leer el MBR actualizado despues del mount:", err.Error())
		return
	}

	//Structs.PrintMBR(TempMBR2)

	defer file.Close()
	Structs.TextoEnviar.WriteString(fmt.Sprintf("🔨 Se monto la particion %s\n", id))
	fmt.Println("======End MOUNT======")
}

func ImprimirEBR(ebr Structs.EBR) {
	fmt.Println(fmt.Sprintf("Part_Mount %s\n Part_Fit %s\n Part_Start %d\n Part_Size %d\n Part_Next. %d\n Part_Name %s",
		string(ebr.Part_mount[:]),
		string(ebr.Part_fit[:]),
		ebr.Part_start,
		ebr.Part_size,
		ebr.Part_next,
		string(ebr.Part_name[:])))
}
