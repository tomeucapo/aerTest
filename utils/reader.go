package utils

import (
	"aerTest/domain"
	"log"
	"strconv"

	r "gopkg.in/gorethink/gorethink.v3"
)

func LogsReader(requestType string, connOpts r.ConnectOpts, c chan domain.Request) {
	session, err := r.Connect(connOpts)

	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	log.Println("Connected to rethinkdb, and get all requests: ", requestType)
	res, err := r.Table("connectorlog").GetAllByIndex("action", requestType).Run(session)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Close()

	if res.IsNil() {
		log.Println("Rows not found for ", requestType)
		return
	}

	count := 0
	var response interface{}
	for res.Next(&response) {
		if res.Err() != nil {
			log.Fatalln(res.Err())
			continue
		}

		rqData := response.(map[string]interface{})

		chain, _ := strconv.Atoi(rqData["chain"].(string))
		parameters := rqData["parameters"].([]interface{})

		//log.Println("Parameters = ", parameters)

		hotelsInterface := parameters[2].([]interface{})
		hotels := make([]string, len(hotelsInterface))
		for i := 0; i < len(hotelsInterface); i++ {
			hotels[i] = hotelsInterface[0].(string)
		}

		rq := domain.Request{Usuario: domain.Usuario{Codigo: rqData["user"].(string), CadenaId: chain},
			CodIdioma:         parameters[1].(string),
			Hoteles:           hotels,
			Entrada:           parameters[3].(string),
			Salida:            parameters[4].(string),
			RespuestaGenerica: true}

		paramOpcionales := domain.Opcionales{}
		posOptionals := 0
		if requestType == "motor.dispo" {
			posOptionals = 6
		} else if requestType == "motor.dispo_limit_ocupa" {
			posOptionals = 7
		}

		if len(parameters) > posOptionals && posOptionals != 0 {
			optionals := parameters[posOptionals].(map[string]interface{})
			if _, ok := optionals["aplicar_dto_general"]; ok {
				paramOpcionales.AplicarDtoGral = optionals["aplicar_dto_general"].(bool)
			}
			if _, ok := optionals["canal"]; ok {
				paramOpcionales.Canal = optionals["canal"].(string)
			}
			if _, ok := optionals["cod_pais_mercado"]; ok {
				paramOpcionales.CodPaisMercado = optionals["cod_pais_mercado"].(string)
			}

			// Ocupacion opcional
			if _, ok := optionals["num_adultos"]; ok {
				paramOpcionales.Adultos = int(optionals["num_adultos"].(float64))
			}
			if _, ok := optionals["num_ninos"]; ok {
				paramOpcionales.Ninos = int(optionals["num_ninos"].(float64))
			}
			if _, ok := optionals["num_bebes"]; ok {
				paramOpcionales.Bebes = int(optionals["num_bebes"].(float64))
			}
		}

		if parameters[5] != nil {
			if requestType == "motor.dispo" {
				rq.CodigoPromocion = parameters[5].(string)
			} else if requestType == "motor.dispo_limit_ocupa" {
				paramOpcionales.Adultos = int(parameters[5].(float64))
				paramOpcionales.Ninos = int(parameters[6].(float64))
			}
		}

		rq.Opcionales = paramOpcionales
		
		c <- rq
		count++
	}

	log.Printf("*** Loaded %d records of %s from logs!\n", count, requestType)
}
