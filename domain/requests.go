package domain

type Usuario struct {
	AfiliadoId int    `json:"afiliado_id,omitempty"`
	CadenaId   int    `json:"cadena_id"`
	Codigo     string `json:"codigo"`
}

type Opcionales struct {
	AplicarDtoGral bool   `json:"aplicar_dto_general"`
	Canal          string `json:"canal"`
	Adultos        int    `json:"num_adultos,omitempty"`
	Ninos          int    `json:"num_ninos,omitempty"`
	Bebes          int    `json:"num_bebes,omitempty"`
	Edades         []int  `json:"edades_ninos,omitempty"`
	CodPaisMercado string `json:"cod_pais_mercado,omitempty"`
}

type Request struct {
	Usuario           Usuario    `json:"usuario"`
	CodIdioma         string     `json:"cod_idioma"`
	Hoteles           []string   `json:"hoteles"`
	Entrada           string     `json:"entrada"`
	Salida            string     `json:"salida"`
	Opcionales        Opcionales `json:"opcionales"`
	RespuestaGenerica bool       `json:"respuesta_generica"`
	CodigoPromocion   string     `json:"cod_promocion,omitempty"`
}
