package handler

import (
	"github.com/CycleZero/mc-msmp-go/dto"
	"log"
)

func DefaultHandler(request *dto.MsmpRequest, response dto.MsmpResponse) {
	log.Println("=============================")
	log.Println("DefaultHandler")
	log.Println("Request:", request)
	log.Println("Response:", response)
	log.Println("=============================")
}
