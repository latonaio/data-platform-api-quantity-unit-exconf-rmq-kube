package dpfm_api_output_formatter

type MetaData struct {
	ConnectionKey     string        `json:"connection_key"`
	Result            bool          `json:"result"`
	RedisKey          string        `json:"redis_key"`
	Filepath          string        `json:"filepath"`
	APIStatusCode     int           `json:"api_status_code"`
	RuntimeSessionID  string        `json:"runtime_session_id"`
	BusinessPartnerID *int          `json:"business_partner"`
	ServiceLabel      string        `json:"service_label"`
	QuantityUnit      *QuantityUnit `json:"QuantityUnit,omitempty"`
	APISchema         string        `json:"api_schema"`
	Accepter          []string      `json:"accepter"`
	Deleted           bool          `json:"deleted"`
}

type QuantityUnit struct {
	QuantityUnit  string `json:"QuantityUnit"`
	ExistenceConf bool   `json:"ExistenceConf"`
}
