package handler

import (
	"encoding/json"
	"github.com/CycleZero/mc-msmp-go/dto"
	"log"
)

func DefaultHandler(request *dto.MsmpRequest, response dto.MsmpResponse) {
	reqStr, err := json.Marshal(request)
	if err != nil {
		log.Println("DefaultHandler json.Marshal error:", err)
	}
	resStr, err := json.Marshal(response)
	if err != nil {
		log.Println("DefaultHandler json.Marshal error:", err)
	}

	log.Println("=============================")
	log.Println("DefaultHandler")
	log.Println("Request:", string(reqStr))
	log.Println("Response:", string(resStr))
	log.Println("=============================")
}
