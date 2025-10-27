package Comandos

import (
	utils "Backend/Structs"
	"fmt"
	"os"
	"strings"
)

func RemoveDisk(path string) error {
	path = strings.Trim(path, ` "'`)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		utils.TextoEnviar.WriteString("❌ Erro: No existe el disco que se desea borrar")
		fmt.Println(path)
		return fmt.Errorf("El disco no existe: %s", path)
	}

	err := os.Remove(path)
	if err != nil {
		utils.TextoEnviar.WriteString("Error eliminando el disco")
		return fmt.Errorf("Error eliminando el disco: %v", err)
	}
	log := fmt.Sprintf("🚨 Se elimino el disco en: %s", path)
	utils.TextoEnviar.WriteString(log)
	fmt.Printf("Disco Elimando con Éxito: %s\n", path)
	return nil
}
