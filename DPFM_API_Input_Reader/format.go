package dpfm_api_input_reader

import (
	"data-platform-api-quantity-unit-exconf-rmq-kube/DPFM_API_Caller/requests"
)

func (sdc *SDC) ConvertToQuantityUnit() *requests.QuantityUnit {
	data := sdc.QuantityUnit
	return &requests.QuantityUnit{
		QuantityUnit: data.QuantityUnit,
	}
}
