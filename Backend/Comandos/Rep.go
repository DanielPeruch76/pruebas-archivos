package Comandos

import (
	"Backend/ParametrosStructs"
	"Backend/Structs"
	"encoding/binary"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Rep(parametros ParametrosStructs.ParametrosRep) {
	id := parametros.Id
	name := parametros.Name
	path := parametros.Path
	path_file_ls := parametros.PathFileLs

	fmt.Println("======Start Rep======")
	fmt.Println("Name:", name)
	fmt.Println("Path:", path)
	fmt.Println("Path file ls:", path_file_ls)
	fmt.Println("Id", id)

	if name == "mbr" {
		diskPath, encontrado := BuscarDisco(id)
		if !encontrado {
			Structs.TextoEnviar.WriteString("❌ Error: La particion no ha sido montada\n")
			return
		}

		file, err := Structs.AbrirArchivo(diskPath)
		if err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No se pudo abrir el archivo del disco\n")
			return
		}
		defer file.Close()

		var TempMBR Structs.MRB
		if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No al leer el MBR del disco\n")
			return
		}

		var builder strings.Builder
		builder.WriteString(`digraph G {
    rankdir=TB;
    node [shape=none];
    
    // Tabla MBR
    mbr_table [label=<
        <TABLE BORDER="1" CELLBORDER="1" CELLSPACING="0" CELLPADDING="4" STYLE="ROUNDED">
            <TR>
                <TD COLSPAN="2" BGCOLOR="darkblue" ALIGN="CENTER">
                    <FONT COLOR="white"><B>REPORTE MBR</B></FONT>
                </TD>
            </TR>
`)

		colorFondo1 := "lightgray"
		colorFondo2 := "white"

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "tamaño", colorFondo2, TempMBR.MbrSize))

		var fit string
		fitStr := LimpiarString(strings.TrimRight(string(TempMBR.Fit[:]), "\x00"))
		if fitStr == "bf" {
			fit = "Best Fit"
		} else if fitStr == "ff" {
			fit = "First Fit"
		} else {
			fit = "Worst Fit"
		}

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Ajuste", colorFondo2, html.EscapeString(fit)))

		creationDate := LimpiarString(string(TempMBR.CreationDate[:]))
		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Fecha Creacion", colorFondo2, html.EscapeString(creationDate)))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "DSK Signature", colorFondo2, TempMBR.Signature))

		builder.WriteString(`        </TABLE>
    >];
`)

		particionCount := 0
		for i := 0; i < 4; i++ {
			if TempMBR.Partitions[i].Size != 0 {
				builder.WriteString(fmt.Sprintf(`
    // Tabla Partición %d
    part_table_%d [label=<
        <TABLE BORDER="1" CELLBORDER="1" CELLSPACING="0" CELLPADDING="4" STYLE="ROUNDED">
            <TR>
                <TD COLSPAN="2" BGCOLOR="darkgreen" ALIGN="CENTER">
                    <FONT COLOR="white"><B>REPORTE PARTICION %d</B></FONT>
                </TD>
            </TR>
`, i, i, i))

				var status string
				statusStr := LimpiarString(strings.TrimRight(string(TempMBR.Partitions[i].Status[:]), "\x00"))
				if statusStr == "0" {
					status = "No montado"
				} else {
					status = "Montado"
				}

				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "ESTADO", colorFondo2, html.EscapeString(status)))

				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "tamaño", colorFondo2, TempMBR.Partitions[i].Size))

				var partFit string
				partFitStr := LimpiarString(strings.TrimRight(string(TempMBR.Partitions[i].Fit[:]), "\x00"))
				if partFitStr == "b" {
					partFit = "Best Fit"
				} else if partFitStr == "f" {
					partFit = "First Fit"
				} else {
					partFit = "Worst Fit"
				}

				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Ajuste", colorFondo2, html.EscapeString(partFit)))

				nombre := LimpiarString(string(TempMBR.Partitions[i].Name[:]))
				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Nombre", colorFondo2, html.EscapeString(nombre)))

				tipo := LimpiarString(string(TempMBR.Partitions[i].Type[:]))
				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Tipo", colorFondo2, html.EscapeString(tipo)))

				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "Correlativo", colorFondo2, TempMBR.Partitions[i].Correlative))

				partitionId := LimpiarString(string(TempMBR.Partitions[i].Id[:]))
				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Id", colorFondo2, html.EscapeString(partitionId)))

				builder.WriteString(`        </TABLE>
    >];
`)
				particionCount++
			}
		}

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
			fmt.Println("Particion extendida encontrada")
			Structs.PrintPartition(TempMBR.Partitions[index2])

			indexEBR := TempMBR.Partitions[index2].Start
			i := particionCount

			for {
				var tempEBR Structs.EBR
				if errNuevoEBR := Structs.LeerEnDisco(file, &tempEBR, int64(indexEBR)); errNuevoEBR != nil {
					fmt.Printf("Error leyendo EBR en offset %d: %v\n", indexEBR, errNuevoEBR)
					break
				}

				if tempEBR.Part_size == 0 {
					fmt.Printf("EBR vacío encontrado en offset %d, terminando lectura\n", indexEBR)
					break
				}

				fmt.Printf("EBR encontrado: %s, tamaño: %d, next: %d\n",
					string(tempEBR.Part_name[:]), tempEBR.Part_size, tempEBR.Part_next)

				builder.WriteString(fmt.Sprintf(`
    // Tabla EBR %d
    ebr_table_%d [label=<
        <TABLE BORDER="1" CELLBORDER="1" CELLSPACING="0" CELLPADDING="4" STYLE="ROUNDED">
            <TR>
                <TD COLSPAN="2" BGCOLOR="darkred" ALIGN="CENTER">
                    <FONT COLOR="white"><B>REPORTE PARTICION %d</B></FONT>
                </TD>
            </TR>
`, i, i, i))

				var status string
				statusStr := LimpiarString(strings.TrimRight(string(tempEBR.Part_mount[:]), "\x00"))
				if statusStr == "0" {
					status = "No montado"
				} else {
					status = "Montado"
				}

				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "ESTADO", colorFondo2, html.EscapeString(status)))

				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "tamaño", colorFondo2, tempEBR.Part_size))

				var partFit string
				partFitStr := LimpiarString(strings.TrimRight(string(tempEBR.Part_fit[:]), "\x00"))
				if partFitStr == "b" {
					partFit = "Best Fit"
				} else if partFitStr == "f" {
					partFit = "First Fit"
				} else {
					partFit = "Worst Fit"
				}

				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Ajuste", colorFondo2, html.EscapeString(partFit)))

				nombre := LimpiarString(string(tempEBR.Part_name[:]))
				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Nombre", colorFondo2, html.EscapeString(nombre)))

				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Tipo", colorFondo2, "Logica"))

				builder.WriteString(`        </TABLE>
    >];
`)

				i++
				particionCount++

				if tempEBR.Part_next == -1 {
					fmt.Println("No hay más EBRs, terminando lectura")
					break
				}

				indexEBR = tempEBR.Part_next
				fmt.Printf("Moviento al siguiente EBR en offset: %d\n", indexEBR)
			}
		}

		if particionCount > 0 {
			builder.WriteString(`
    // Organizar tablas verticalmente
    mbr_table -> part_table_0 [style=invis];
`)

			firstPartIndex := -1
			prevPartIndex := -1
			for i := 0; i < 4; i++ {
				if TempMBR.Partitions[i].Size != 0 {
					if firstPartIndex == -1 {
						firstPartIndex = i
					}
					if prevPartIndex != -1 && prevPartIndex != i {
						builder.WriteString(fmt.Sprintf("    part_table_%d -> part_table_%d [style=invis];\n", prevPartIndex, i))
					}
					prevPartIndex = i
				}
			}
		}

		builder.WriteString("}")

		fmt.Println("Código Graphviz generado:")
		fmt.Println(builder.String())

		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creando directorio: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error creando directorio\n")
			return
		}

		tempDot := filepath.Join(dir, "temp.dot")
		err = os.WriteFile(tempDot, []byte(builder.String()), 0644)
		if err != nil {
			fmt.Printf("Error creando archivo .dot: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error creando archivo temporal\n")
			return
		}
		defer os.Remove(tempDot)

		cmd := exec.Command("dot", "-Tjpg", tempDot, "-o", path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error ejecutando Graphviz: %v\n", err)
			fmt.Printf("Output: %s\n", output)
			Structs.TextoEnviar.WriteString("❌ Error generando imagen Graphviz\n")
			debugDot := filepath.Join(dir, "debug.dot")
			os.WriteFile(debugDot, []byte(builder.String()), 0644)
			fmt.Printf("Código DOT guardado en: %s para debug\n", debugDot)

			return
		}

		fmt.Printf("Reporte generado exitosamente en: %s\n", path)
		Structs.TextoEnviar.WriteString("✅ Reporte MBR generado exitosamente\n")
	} else if name == "disk" {
		diskPath, encontrado := BuscarDisco(id)
		if !encontrado {
			Structs.TextoEnviar.WriteString("❌ Error: La particion no ha sido montada\n")
			return
		}

		file, err := Structs.AbrirArchivo(diskPath)
		if err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No se pudo abrir el archivo del disco\n")
			return
		}
		defer file.Close()

		var TempMBR Structs.MRB
		if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No al leer el MBR del disco\n")
			return
		}

		var builder strings.Builder

		builder.WriteString(`digraph D {
    rankdir=LR;
    subgraph cluster_0 {
        bgcolor="#68d9e2"
        node [style="rounded" style=filled];
        label="Estructura del Disco";
        labelloc="t";
        
        disk_node [shape=record label="`)

		partitions := make([]string, 0)

		for i := 0; i < 4; i++ {
			if TempMBR.Partitions[i].Size != 0 {
				nombre := LimpiarString(string(TempMBR.Partitions[i].Name[:]))
				tipo := LimpiarString(string(TempMBR.Partitions[i].Type[:]))

				nombre = strings.ReplaceAll(nombre, "|", "\\|")
				nombre = strings.ReplaceAll(nombre, "{", "\\{")
				nombre = strings.ReplaceAll(nombre, "}", "\\}")
				nombre = strings.ReplaceAll(nombre, "\"", "\\\"")

				if tipo == "e" {
					ebrParts := make([]string, 0)
					indexEBR := TempMBR.Partitions[i].Start

					for {
						var tempEBR Structs.EBR
						if errNuevoEBR := Structs.LeerEnDisco(file, &tempEBR, int64(indexEBR)); errNuevoEBR != nil {
							break
						}

						if tempEBR.Part_size == 0 {
							break
						}

						nombreEBR := LimpiarString(string(tempEBR.Part_name[:]))
						nombreEBR = strings.ReplaceAll(nombreEBR, "|", "\\|")
						nombreEBR = strings.ReplaceAll(nombreEBR, "{", "\\{")
						nombreEBR = strings.ReplaceAll(nombreEBR, "}", "\\}")
						nombreEBR = strings.ReplaceAll(nombreEBR, "\"", "\\\"")

						if nombreEBR != "" {
							ebrParts = append(ebrParts, "EBR|"+nombreEBR)
						} else {
							ebrParts = append(ebrParts, "EBR|Libre")
						}

						if tempEBR.Part_next == -1 {
							break
						}
						indexEBR = tempEBR.Part_next
					}

					if len(ebrParts) > 0 {
						partitions = append(partitions, "{Extendida|{"+strings.Join(ebrParts, "|")+"}}")
					} else {
						partitions = append(partitions, "{Extendida|Libre}")
					}
				} else {
					partitions = append(partitions, nombre)
				}
			} else {
				partitions = append(partitions, "Libre")
			}
		}

		builder.WriteString("MBR|")
		builder.WriteString(strings.Join(partitions, "|"))
		builder.WriteString(`"];
    }
}`)

		fmt.Println("Código Graphviz generado:")
		fmt.Println(builder.String())

		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creando directorio: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error creando directorio\n")
			return
		}

		tempDot := filepath.Join(dir, "temp.dot")
		err = os.WriteFile(tempDot, []byte(builder.String()), 0644)
		if err != nil {
			fmt.Printf("Error creando archivo .dot: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error creando archivo temporal\n")
			return
		}
		defer os.Remove(tempDot)

		cmd := exec.Command("dot", "-Tjpg", tempDot, "-o", path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error ejecutando Graphviz: %v\n", err)
			fmt.Printf("Output: %s\n", output)
			Structs.TextoEnviar.WriteString("❌ Error generando imagen Graphviz\n")
			debugDot := filepath.Join(dir, "debug.dot")
			os.WriteFile(debugDot, []byte(builder.String()), 0644)
			fmt.Printf("Código DOT guardado en: %s para debug\n", debugDot)

			return
		}

		fmt.Printf("Reporte generado exitosamente en: %s\n", path)
		Structs.TextoEnviar.WriteString("✅ Reporte de disco generado exitosamente\n")
	} else if name == "inode" {

		diskPath, encontrado := BuscarDisco(id)
		if !encontrado {
			Structs.TextoEnviar.WriteString("❌ Error: La particion no ha sido montada\n")
			return
		}

		file, err := Structs.AbrirArchivo(diskPath)
		if err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No se pudo abrir el archivo del disco\n")
			return
		}
		defer file.Close()

		var TempMBR Structs.MRB
		if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No al leer el MBR del disco\n")
			return
		}
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
			fmt.Println("Partition no fue encontrada")
			return
		}

		fmt.Println("ID:", string(Structs.Usuario.ID[:]))
		fmt.Println("index:", index)

		var tempSuperblock Structs.Superblock

		if err := Structs.LeerEnDisco(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
			return
		}

		numeradorInodo := tempSuperblock.S_fist_ino - tempSuperblock.S_inode_start
		denominadorInodo := int32(binary.Size(Structs.Inode{}))
		cantidadRecorrido := numeradorInodo / denominadorInodo

		limpiarParaHTML := func(s string) string {
			cleaned := strings.TrimRight(s, "\x00")
			cleaned = strings.ReplaceAll(cleaned, "&", "&amp;")
			cleaned = strings.ReplaceAll(cleaned, "<", "&lt;")
			cleaned = strings.ReplaceAll(cleaned, ">", "&gt;")
			cleaned = strings.ReplaceAll(cleaned, "\"", "&quot;")
			cleaned = strings.ReplaceAll(cleaned, "'", "&#39;")
			return cleaned
		}

		var builder strings.Builder

		builder.WriteString(`digraph G {
    rankdir=LR;  // Dirección horizontal
    node [shape=none];
    nodesep=0.5;  // Separación entre nodos

`)

		var inodoTables []string
		var indiceInodo = tempSuperblock.S_inode_start

		for i := 0; i < int(cantidadRecorrido); i++ {
			var Inode0 Structs.Inode
			if err1 := Structs.LeerEnDisco(file, &Inode0, int64(indiceInodo)); err1 != nil {
				return
			}

			builder.WriteString(fmt.Sprintf(`
    inodo_%d [label=<
        <TABLE BORDER="1" CELLBORDER="1" CELLSPACING="0" CELLPADDING="4" STYLE="ROUNDED">
            <TR>
                <TD COLSPAN="2" BGCOLOR="darkblue" ALIGN="CENTER">
                    <FONT COLOR="white"><B>INODO %d</B></FONT>
                </TD>
            </TR>
`, i, i))

			colorFondo1 := "lightgray"
			colorFondo2 := "white"

			builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "UID", colorFondo2, Inode0.I_uid))

			builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "GID", colorFondo2, Inode0.I_gid))

			builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "Size", colorFondo2, Inode0.I_size))

			fechaAcceso := limpiarParaHTML(string(Inode0.I_atime[:]))
			fechaModificacion := limpiarParaHTML(string(Inode0.I_mtime[:]))
			fechaCreacion := limpiarParaHTML(string(Inode0.I_ctime[:]))

			builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Fecha acceso", colorFondo2, fechaAcceso))

			builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Fecha modificación", colorFondo2, fechaModificacion))

			builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Fecha creación", colorFondo2, fechaCreacion))

			for j := 0; j < 12; j++ {
				if Inode0.I_block[j] != -1 {
					builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>Bloque %d</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, j, colorFondo2, Inode0.I_block[j]))
				}
			}

			if Inode0.I_block[12] != -1 {
				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>Bloque indirecto</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, colorFondo2, Inode0.I_block[12]))
			}

			if Inode0.I_block[13] != -1 {
				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>Bloque doble indirecto</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, colorFondo2, Inode0.I_block[13]))
			}

			if Inode0.I_block[14] != -1 {
				builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>Bloque triple indirecto</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, colorFondo2, Inode0.I_block[14]))
			}

			builder.WriteString(`        </TABLE>
    >];
`)

			inodoTables = append(inodoTables, fmt.Sprintf("inodo_%d", i))

			indiceInodo += int32(binary.Size(Structs.Inode{}))
		}

		if len(inodoTables) > 1 {
			for i := 0; i < len(inodoTables)-1; i++ {
				builder.WriteString(fmt.Sprintf("    %s -> %s [weight=2];\n", inodoTables[i], inodoTables[i+1]))
			}
		}

		builder.WriteString("}")

		fmt.Println("Código Graphviz generado:")
		fmt.Println(builder.String())

		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creando directorio: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error creando directorio\n")
			return
		}

		tempDot := filepath.Join(dir, "temp.dot")
		err = os.WriteFile(tempDot, []byte(builder.String()), 0644)
		if err != nil {
			fmt.Printf("Error creando archivo .dot: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error creando archivo temporal\n")
			return
		}
		defer os.Remove(tempDot)

		cmd := exec.Command("dot", "-Tjpg", tempDot, "-o", path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error ejecutando Graphviz: %v\n", err)
			fmt.Printf("Output: %s\n", output)
			Structs.TextoEnviar.WriteString("❌ Error generando imagen Graphviz\n")

			debugDot := filepath.Join(dir, "debug.dot")
			os.WriteFile(debugDot, []byte(builder.String()), 0644)
			fmt.Printf("Código DOT guardado en: %s para debug\n", debugDot)
			return
		}

		fmt.Printf("Reporte de inodos generado exitosamente en: %s\n", path)
		Structs.TextoEnviar.WriteString("✅ Reporte de inodos generado exitosamente\n")
	} else if name == "block" {
		diskPath, encontrado := BuscarDisco(id)
		if !encontrado {
			Structs.TextoEnviar.WriteString("❌ Error: La particion no ha sido montada\n")
			return
		}

		file, err := Structs.AbrirArchivo(diskPath)
		if err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No se pudo abrir el archivo del disco\n")
			return
		}
		defer file.Close()

		var TempMBR Structs.MRB
		if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No al leer el MBR del disco\n")
			return
		}
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
			fmt.Println("Partition no fue encontrada")
			return
		}

		fmt.Println("ID:", string(Structs.Usuario.ID[:]))
		fmt.Println("index:", index)

		var tempSuperblock Structs.Superblock

		if err := Structs.LeerEnDisco(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
			return
		}

		numeradorBloque := tempSuperblock.S_first_blo - tempSuperblock.S_block_start
		denominadorBloque := int32(binary.Size(Structs.Folderblock{}))
		cantidadRecorrido := numeradorBloque / denominadorBloque

		limpiarParaHTML := func(s string) string {
			cleaned := strings.TrimRight(s, "\x00")
			cleaned = strings.ReplaceAll(cleaned, "&", "&amp;")
			cleaned = strings.ReplaceAll(cleaned, "<", "&lt;")
			cleaned = strings.ReplaceAll(cleaned, ">", "&gt;")
			cleaned = strings.ReplaceAll(cleaned, "\"", "&quot;")
			cleaned = strings.ReplaceAll(cleaned, "'", "&#39;")
			return cleaned
		}

		var builder strings.Builder

		builder.WriteString(`digraph G {
    rankdir=LR;  // Dirección horizontal
    node [shape=none];
    nodesep=0.5;  // Separación entre nodos

`)

		var bloqueTables []string
		var indiceBloque = tempSuperblock.S_block_start

		for i := 0; i < int(cantidadRecorrido); i++ {
			var folderBlock Structs.Folderblock
			if err := Structs.LeerEnDisco(file, &folderBlock, int64(indiceBloque)); err != nil {
				return
			}

			var fileBlock Structs.Fileblock
			if err := Structs.LeerEnDisco(file, &fileBlock, int64(indiceBloque)); err != nil {
				return
			}

			esFolderBlock := false
			esFileBlock := false

			for _, content := range folderBlock.B_content {
				if content.B_inodo != -1 && content.B_inodo != 0 {
					esFolderBlock = true
					break
				}
			}

			contenidoArchivo := strings.TrimRight(string(fileBlock.B_content[:]), "\x00")
			if contenidoArchivo != "" && !esFolderBlock {
				esFileBlock = true
			}

			builder.WriteString(fmt.Sprintf(`
    bloque_%d [label=<
        <TABLE BORDER="1" CELLBORDER="1" CELLSPACING="0" CELLPADDING="4" STYLE="ROUNDED">
            <TR>
                <TD COLSPAN="2" BGCOLOR="darkgreen" ALIGN="CENTER">
                    <FONT COLOR="white"><B>BLOQUE %d</B></FONT>
                </TD>
            </TR>
`, i, i))

			colorFondo1 := "lightgray"
			colorFondo2 := "white"

			if esFolderBlock {
				builder.WriteString(fmt.Sprintf(`            <TR>
            <TD BGCOLOR="%s"><B>%s</B></TD>
            <TD BGCOLOR="%s">%s</TD>
        </TR>
`, colorFondo1, "Tipo", colorFondo2, "Folder Block"))

				for j, content := range folderBlock.B_content {
					if content.B_inodo != -1 && content.B_inodo != 0 {
						nombre := limpiarParaHTML(strings.TrimRight(string(content.B_name[:]), "\x00"))
						builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>Contenido %d</B></TD>
                <TD BGCOLOR="%s">Nombre: %s, Inodo: %d</TD>
            </TR>
`, colorFondo1, j, colorFondo2, nombre, content.B_inodo))
					}
				}

			} else if esFileBlock {
				builder.WriteString(fmt.Sprintf(`            <TR>
            <TD BGCOLOR="%s"><B>%s</B></TD>
            <TD BGCOLOR="%s">%s</TD>
        </TR>
`, colorFondo1, "Tipo", colorFondo2, "File Block"))

				contenido := limpiarParaHTML(strings.TrimRight(string(fileBlock.B_content[:]), "\x00"))
				builder.WriteString(fmt.Sprintf(`            <TR>
            <TD BGCOLOR="%s"><B>Contenido</B></TD>
            <TD BGCOLOR="%s">%s</TD>
        </TR>
`, colorFondo1, colorFondo2, contenido))

			} else {
				builder.WriteString(fmt.Sprintf(`            <TR>
            <TD BGCOLOR="%s"><B>%s</B></TD>
            <TD BGCOLOR="%s">%s</TD>
        </TR>
`, colorFondo1, "Tipo", colorFondo2, "Bloque Libre/Vacío"))
			}

			builder.WriteString(`        </TABLE>
    >];
`)

			bloqueTables = append(bloqueTables, fmt.Sprintf("bloque_%d", i))

			indiceBloque += int32(binary.Size(Structs.Folderblock{}))
		}

		if len(bloqueTables) > 1 {
			for i := 0; i < len(bloqueTables)-1; i++ {
				builder.WriteString(fmt.Sprintf("    %s -> %s [weight=2];\n", bloqueTables[i], bloqueTables[i+1]))
			}
		}

		builder.WriteString("}")

		fmt.Println("Código Graphviz generado:")
		fmt.Println(builder.String())

		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creando directorio: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error creando directorio\n")
			return
		}

		tempDot := filepath.Join(dir, "temp.dot")
		err = os.WriteFile(tempDot, []byte(builder.String()), 0644)
		if err != nil {
			fmt.Printf("Error creando archivo .dot: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error creando archivo temporal\n")
			return
		}
		defer os.Remove(tempDot)

		cmd := exec.Command("dot", "-Tjpg", tempDot, "-o", path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error ejecutando Graphviz: %v\n", err)
			fmt.Printf("Output: %s\n", output)
			Structs.TextoEnviar.WriteString("❌ Error generando imagen Graphviz\n")

			debugDot := filepath.Join(dir, "debug.dot")
			os.WriteFile(debugDot, []byte(builder.String()), 0644)
			fmt.Printf("Código DOT guardado en: %s para debug\n", debugDot)
			return
		}

		fmt.Printf("Reporte de bloques generado exitosamente en: %s\n", path)
		Structs.TextoEnviar.WriteString("✅ Reporte de bloques generado exitosamente\n")
	} else if name == "bm_inode" {
		diskPath, encontrado := BuscarDisco(id)
		if !encontrado {
			Structs.TextoEnviar.WriteString("❌ Error: La particion no ha sido montada\n")
			return
		}

		file, err := Structs.AbrirArchivo(diskPath)
		if err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No se pudo abrir el archivo del disco\n")
			return
		}
		defer file.Close()

		var TempMBR Structs.MRB
		if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No al leer el MBR del disco\n")
			return
		}
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
			fmt.Println("Partition no fue encontrada")
			return
		}

		fmt.Println("ID:", string(Structs.Usuario.ID[:]))
		fmt.Println("index:", index)

		var tempSuperblock Structs.Superblock

		if err := Structs.LeerEnDisco(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
			return
		}

		numeradorInodo := tempSuperblock.S_fist_ino - tempSuperblock.S_inode_start
		denominadorInodo := int32(binary.Size(Structs.Inode{}))
		cantidadRecorrido := numeradorInodo / denominadorInodo
		cantidadRecorrido = tempSuperblock.S_bm_block_start - tempSuperblock.S_bm_inode_start

		indice_bm_inodo := tempSuperblock.S_bm_inode_start

		var valores []byte
		for i := 0; i < int(cantidadRecorrido); i++ {
			var valor [1]byte
			if err := Structs.LeerEnDisco(file, &valor, int64(indice_bm_inodo)); err != nil {
				return
			}
			valores = append(valores, valor[0])
			indice_bm_inodo++
		}

		var contenido strings.Builder
		contador := 0

		for i, valor := range valores {
			if valor == 0 {
				contenido.WriteString("0")
			} else if valor == 1 {
				contenido.WriteString("1")
			} else {
				contenido.WriteString(fmt.Sprintf("%d", valor))
			}

			contador++

			if contador < 20 && i < len(valores)-1 {
				contenido.WriteString(" ")
			}

			if contador == 20 && i < len(valores)-1 {
				contenido.WriteString("\n")
				contador = 0
			}
		}

		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creando directorio: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error creando directorio\n")
			return
		}

		err = ioutil.WriteFile(path, []byte(contenido.String()), 0644)
		if err != nil {
			fmt.Printf("Error escribiendo archivo: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error escribiendo archivo\n")
			return
		}

		fmt.Printf("Archivo creado exitosamente: %s\n", path)
		fmt.Printf("Total de valores escritos: %d\n", len(valores))
		Structs.TextoEnviar.WriteString("✅ Archivo de bitmap creado exitosamente\n")

	} else if name == "bm_bloc" {
		diskPath, encontrado := BuscarDisco(id)
		if !encontrado {
			Structs.TextoEnviar.WriteString("❌ Error: La particion no ha sido montada\n")
			return
		}

		file, err := Structs.AbrirArchivo(diskPath)
		if err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No se pudo abrir el archivo del disco\n")
			return
		}
		defer file.Close()

		var TempMBR Structs.MRB
		if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No al leer el MBR del disco\n")
			return
		}
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
			fmt.Println("Partition no fue encontrada")
			return
		}

		fmt.Println("ID:", string(Structs.Usuario.ID[:]))
		fmt.Println("index:", index)

		var tempSuperblock Structs.Superblock

		if err := Structs.LeerEnDisco(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
			return
		}

		numeradorInodo := tempSuperblock.S_fist_ino - tempSuperblock.S_inode_start
		denominadorInodo := int32(binary.Size(Structs.Inode{}))
		cantidadRecorrido := numeradorInodo / denominadorInodo
		cantidadRecorrido = tempSuperblock.S_inode_start - tempSuperblock.S_bm_block_start

		indice_bm_inodo := tempSuperblock.S_bm_block_start

		// Leer todos los valores del bitmap
		var valores []byte
		for i := 0; i < int(cantidadRecorrido); i++ {
			var valor [1]byte
			if err := Structs.LeerEnDisco(file, &valor, int64(indice_bm_inodo)); err != nil {
				return
			}
			valores = append(valores, valor[0])
			indice_bm_inodo++
		}

		var contenido strings.Builder
		contador := 0

		for i, valor := range valores {
			if valor == 0 {
				contenido.WriteString("0")
			} else if valor == 1 {
				contenido.WriteString("1")
			} else {
				contenido.WriteString(fmt.Sprintf("%d", valor))
			}

			contador++

			if contador < 20 && i < len(valores)-1 {
				contenido.WriteString(" ")
			}

			if contador == 20 && i < len(valores)-1 {
				contenido.WriteString("\n")
				contador = 0
			}
		}

		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creando directorio: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error creando directorio\n")
			return
		}

		err = ioutil.WriteFile(path, []byte(contenido.String()), 0644)
		if err != nil {
			fmt.Printf("Error escribiendo archivo: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error escribiendo archivo\n")
			return
		}

		fmt.Printf("Archivo creado exitosamente: %s\n", path)
		fmt.Printf("Total de valores escritos: %d\n", len(valores))
		Structs.TextoEnviar.WriteString("✅ Archivo de bitmap creado exitosamente\n")
	} else if name == "sb" {
		diskPath, encontrado := BuscarDisco(id)
		if !encontrado {
			Structs.TextoEnviar.WriteString("❌ Error: La particion no ha sido montada\n")
			return
		}

		file, err := Structs.AbrirArchivo(diskPath)
		if err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No se pudo abrir el archivo del disco\n")
			return
		}
		defer file.Close()

		var TempMBR Structs.MRB
		if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
			Structs.TextoEnviar.WriteString("❌ Error: No al leer el MBR del disco\n")
			return
		}
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
			fmt.Println("Partition no fue encontrada")
			return
		}

		fmt.Println("ID:", string(Structs.Usuario.ID[:]))
		fmt.Println("index:", index)

		var tempSuperblock Structs.Superblock

		if err := Structs.LeerEnDisco(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
			return
		}

		limpiarParaHTML := func(s string) string {
			cleaned := strings.TrimRight(s, "\x00")
			cleaned = strings.ReplaceAll(cleaned, "&", "&amp;")
			cleaned = strings.ReplaceAll(cleaned, "<", "&lt;")
			cleaned = strings.ReplaceAll(cleaned, ">", "&gt;")
			cleaned = strings.ReplaceAll(cleaned, "\"", "&quot;")
			cleaned = strings.ReplaceAll(cleaned, "'", "&#39;")
			return cleaned
		}

		var builder strings.Builder

		builder.WriteString(`digraph G {
    rankdir=TB;
    node [shape=none];
    
    // Tabla Superblock
    superblock_table [label=<
        <TABLE BORDER="1" CELLBORDER="1" CELLSPACING="0" CELLPADDING="4" STYLE="ROUNDED">
            <TR>
                <TD COLSPAN="2" BGCOLOR="darkblue" ALIGN="CENTER">
                    <FONT COLOR="white"><B>REPORTE SUPERBLOCK</B></FONT>
                </TD>
            </TR>
`)

		colorFondo1 := "lightgray"
		colorFondo2 := "white"

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "Filesystem Type", colorFondo2, tempSuperblock.S_filesystem_type))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "Inodes Count", colorFondo2, tempSuperblock.S_inodes_count))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "Blocks Count", colorFondo2, tempSuperblock.S_blocks_count))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "Free Blocks Count", colorFondo2, tempSuperblock.S_free_blocks_count))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "Free Inodes Count", colorFondo2, tempSuperblock.S_free_inodes_count))

		mtime := limpiarParaHTML(string(tempSuperblock.S_mtime[:]))
		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Modification Time", colorFondo2, mtime))

		umtime := limpiarParaHTML(string(tempSuperblock.S_umtime[:]))
		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%s</TD>
            </TR>
`, colorFondo1, "Unmount Time", colorFondo2, umtime))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "Mount Count", colorFondo2, tempSuperblock.S_mnt_count))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "Magic", colorFondo2, tempSuperblock.S_magic))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "Inode Size", colorFondo2, tempSuperblock.S_inode_size))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "Block Size", colorFondo2, tempSuperblock.S_block_size))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "First Inode", colorFondo2, tempSuperblock.S_fist_ino))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "First Block", colorFondo2, tempSuperblock.S_first_blo))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "BM Inode Start", colorFondo2, tempSuperblock.S_bm_inode_start))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "BM Block Start", colorFondo2, tempSuperblock.S_bm_block_start))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "Inode Start", colorFondo2, tempSuperblock.S_inode_start))

		builder.WriteString(fmt.Sprintf(`            <TR>
                <TD BGCOLOR="%s"><B>%s</B></TD>
                <TD BGCOLOR="%s">%d</TD>
            </TR>
`, colorFondo1, "Block Start", colorFondo2, tempSuperblock.S_block_start))

		builder.WriteString(`        </TABLE>
    >];
}`)

		fmt.Println("Código Graphviz generado:")
		fmt.Println(builder.String())

		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creando directorio: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error creando directorio\n")
			return
		}

		tempDot := filepath.Join(dir, "temp.dot")
		err = os.WriteFile(tempDot, []byte(builder.String()), 0644)
		if err != nil {
			fmt.Printf("Error creando archivo .dot: %v\n", err)
			Structs.TextoEnviar.WriteString("❌ Error creando archivo temporal\n")
			return
		}
		defer os.Remove(tempDot)

		cmd := exec.Command("dot", "-Tjpg", tempDot, "-o", path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error ejecutando Graphviz: %v\n", err)
			fmt.Printf("Output: %s\n", output)
			Structs.TextoEnviar.WriteString("❌ Error generando imagen Graphviz\n")
			debugDot := filepath.Join(dir, "debug.dot")
			os.WriteFile(debugDot, []byte(builder.String()), 0644)
			fmt.Printf("Código DOT guardado en: %s para debug\n", debugDot)
			return
		}

		fmt.Printf("Reporte de Superblock generado exitosamente en: %s\n", path)
		Structs.TextoEnviar.WriteString("✅ Reporte de Superblock generado exitosamente\n")

		fmt.Printf("\n=== INFORMACIÓN DEL SUPERBLOCK ===\n")
		fmt.Printf("Filesystem Type: %d\n", tempSuperblock.S_filesystem_type)
		fmt.Printf("Total Inodes: %d, Libres: %d\n", tempSuperblock.S_inodes_count, tempSuperblock.S_free_inodes_count)
		fmt.Printf("Total Blocks: %d, Libres: %d\n", tempSuperblock.S_blocks_count, tempSuperblock.S_free_blocks_count)
		fmt.Printf("Mount Count: %d\n", tempSuperblock.S_mnt_count)
		fmt.Printf("Magic Number: %d\n", tempSuperblock.S_magic)
		fmt.Printf("Inode Size: %d bytes, Block Size: %d bytes\n", tempSuperblock.S_inode_size, tempSuperblock.S_block_size)
		fmt.Printf("Modification Time: %s\n", mtime)
		fmt.Printf("Unmount Time: %s\n", umtime)
	} else if name == "file" {

		if path_file_ls == "" {
			Structs.TextoEnviar.WriteString("Error: No se mando el parametro para buscar el archivo\n")
		} else {
			diskPath, encontrado := BuscarDisco(id)
			if !encontrado {
				Structs.TextoEnviar.WriteString("❌ Error: La particion no ha sido montada\n")
				return
			}

			file, err := Structs.AbrirArchivo(diskPath)
			if err != nil {
				Structs.TextoEnviar.WriteString("❌ Error: No se pudo abrir el archivo del disco\n")
				return
			}
			defer file.Close()

			var TempMBR Structs.MRB
			if err := Structs.LeerEnDisco(file, &TempMBR, 0); err != nil {
				Structs.TextoEnviar.WriteString("❌ Error: No al leer el MBR del disco\n")
				return
			}
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
				fmt.Println("Partition no fue encontrada")
				return
			}

			fmt.Println("ID:", string(Structs.Usuario.ID[:]))
			fmt.Println("index:", index)

			var tempSuperblock Structs.Superblock

			if err := Structs.LeerEnDisco(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
				return
			}

			contenido := strings.Split(path_file_ls, "/")
			pathContenido := contenido[1:]
			contenidoArchivoCopiar, exitoEncontrandoContenido := "", false

			if len(pathContenido) > 0 {
				contenidoArchivoCopiar, exitoEncontrandoContenido = Structs.ObtenerContenido(pathContenido, file, tempSuperblock)
				if exitoEncontrandoContenido {
					Structs.TextoEnviar.WriteString(fmt.Sprintf("📜 Este es el contendio de %s: \n%s\n", pathContenido[len(pathContenido)-1], contenidoArchivoCopiar))
				} else {
					Structs.TextoEnviar.WriteString(fmt.Sprintf("❌ Error: No se encontro el archivo: %s \n", pathContenido[len(pathContenido)-1]))
				}
			} else {
				Structs.TextoEnviar.WriteString("❌ Error: Parametro invalido")
			}

			if exitoEncontrandoContenido {
				dir := filepath.Dir(path)
				if err := os.MkdirAll(dir, 0755); err != nil {
					return
				}

				err := ioutil.WriteFile(path, []byte(contenidoArchivoCopiar), 0644)
				if err != nil {
					return
				}
				fmt.Printf("Archivo creado exitosamente: %s\n", path)
			} else {

			}

		}

	}

}

func contarValores(valores []byte, buscado byte) int {
	contador := 0
	for _, v := range valores {
		if v == buscado {
			contador++
		}
	}
	return contador
}
