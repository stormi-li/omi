package main

import (
	"fmt"

	"github.com/stormi-li/omi/app-login/database"
)

func main(){
	fmt.Println(database.DB.AllowGlobalUpdate)
}