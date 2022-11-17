package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-quantity-unit-exconf-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-quantity-unit-exconf-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-quantity-unit-exconf-rmq-kube/database"
	"sync"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
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

func (e *ExistenceConf) Conf(input *dpfm_api_input_reader.SDC) *dpfm_api_output_formatter.QuantityUnit {
	quantityUnit := *input.QuantityUnit.QuantityUnit
	notKeyExistence := make([]string, 0, 1)
	KeyExistence := make([]string, 0, 1)

	existData := &dpfm_api_output_formatter.QuantityUnit{
		QuantityUnit:       quantityUnit,
		ExistenceConf: false,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if !e.confQuantityUnit(quantityUnit) {
			notKeyExistence = append(notKeyExistence, quantityUnit)
			return
		}
		KeyExistence = append(KeyExistence, quantityUnit)
	}()

	wg.Wait()

	if len(KeyExistence) == 0 {
		return existData
	}
	if len(notKeyExistence) > 0 {
		return existData
	}

	existData.ExistenceConf = true
	return existData
}

func (e *ExistenceConf) confQuantityUnit(val string) bool {
	rows, err := e.db.Query(
		`SELECT QuantityUnit 
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_quantity_unit_quantity_unit_data 
		WHERE QuantityUnit = ?;`, val,
	)
	if err != nil {
		e.l.Error(err)
		return false
	}

	for rows.Next() {
		var quantityUnit string
		err := rows.Scan(&quantityUnit)
		if err != nil {
			e.l.Error(err)
			continue
		}
		if quantityUnit == val {
			return true
		}
	}
	return false
}
