package main

import (
	"go-config/core"
	"log"
)

func init(){
	log.SetPrefix("[Go-heracles-client]")
	log.SetFlags(log.LstdFlags | log.Lshortfile |log.LUTC)
}

func main(){
	log.Println("***********************Go Heracles Client Start*********************")
	core.InitConfig()
}
