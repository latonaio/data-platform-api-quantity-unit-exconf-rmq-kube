package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-quantity-unit-exconf-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-quantity-unit-exconf-rmq-kube/DPFM_API_Output_Formatter"
	"encoding/json"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	database "github.com/latonaio/golang-mysql-network-connector"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
	"golang.org/x/xerrors"
)

type ExistenceConf struct {
	ctx context.Context
	db  *database.Mysql
	l   *logger.Logger
}

func NewExistenceConf(ctx context.Context, db *database.Mysql, l *logger.Logger) *ExistenceConf {
	return &ExistenceConf{
		ctx: ctx,
		db:  db,
		l:   l,
	}
}

func (e *ExistenceConf) Conf(msg rabbitmq.RabbitmqMessage) interface{} {
	var ret interface{}
	ret = map[string]interface{}{
		"ExistenceConf": false,
	}
	input := make(map[string]interface{})
	err := json.Unmarshal(msg.Raw(), &input)
	if err != nil {
		return ret
	}

	_, ok := input["QuantityUnit"]
	if ok {
		input := &dpfm_api_input_reader.SDC{}
		err = json.Unmarshal(msg.Raw(), input)
		ret = e.confQuantityUnit(input)
		goto endProcess
	}

	err = xerrors.Errorf("can not get exconf check target")
endProcess:
	if err != nil {
		e.l.Error(err)
	}
	return ret
}

func (e *ExistenceConf) confQuantityUnit(input *dpfm_api_input_reader.SDC) *dpfm_api_output_formatter.QuantityUnit {
	exconf := dpfm_api_output_formatter.QuantityUnit{
		ExistenceConf: false,
	}
	if input.QuantityUnit.QuantityUnit == nil {
		return &exconf
	}
	exconf = dpfm_api_output_formatter.QuantityUnit{
		QuantityUnit:  *input.QuantityUnit.QuantityUnit,
		ExistenceConf: false,
	}

	rows, err := e.db.Query(
		`SELECT QuantityUnit 
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_quantity_unit_quantity_unit_data 
		WHERE QuantityUnit = ?;`, exconf.QuantityUnit,
	)
	if err != nil {
		e.l.Error(err)
		return &exconf
	}
	defer rows.Close()

	exconf.ExistenceConf = rows.Next()
	return &exconf
}
