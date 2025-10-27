package Comandos

import (
	"Backend/Structs"
	"fmt"
)

func LogOut() {
	fmt.Println("======Start LOGOUT======")
	if Structs.Usuario.Status {
		usuario := Structs.Usuario.User
		Structs.Usuario.ID = ""
		Structs.Usuario.Status = false
		Structs.Usuario.Path = ""
		Structs.Usuario.User = ""
		fmt.Println("User logged out")
		Structs.TextoEnviar.WriteString(fmt.Sprintf("👋 El usuario \"%s\" ha finalizado sesión\n", usuario))
	} else {
		fmt.Println("No user logged in")
		Structs.TextoEnviar.WriteString("🚩 Error: Comando inválido, no hay sesión activa\n")
	}
	fmt.Println("======End LOGOUT======")
}
