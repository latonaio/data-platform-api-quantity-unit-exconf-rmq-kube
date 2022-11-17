# data-platform-api-quantity-unit-exconf-rmq-kube
data-platform-api-quantity-unit-exconf-rmq-kube は、データ連携基盤において、API で数量単位の存在性チェックを行うためのマイクロサービスです。

## 動作環境
・ OS: LinuxOS  
・ CPU: ARM/AMD/Intel  

## 存在確認先テーブル名
以下のsqlファイルに対して、数量単位の存在確認が行われます。

* data-platform-quantity-unit-sql-quantity-unit-data.sql（データ連携基盤 数量単位 - 数量単位データ）

## caller.go による存在性確認
Input で取得されたファイルに基づいて、caller.go で、 API がコールされます。
caller.go の 以下の箇所が、指定された API をコールするソースコードです。

```
func (e *ExistenceConf) Conf(input *dpfm_api_input_reader.SDC) *dpfm_api_output_formatter.QuantityUnit {
	quantityUnit := *input.QuantityUnit.QuantityUnit
	notKeyExistence := make([]int, 0, 1)
	KeyExistence := make([]int, 0, 1)

	existData := &dpfm_api_output_formatter.QuantityUnit{
		QuantityUnit: quantityUnit,
		ExistenceConf:   false,
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
```

## Input
data-platform-api-quantity-unit-exconf-rmq-kube では、以下のInputファイルをRabbitMQからJSON形式で受け取ります。  

```
{
	"connection_key": "request",
	"result": true,
	"redis_key": "abcdefg",
	"api_status_code": 200,
	"runtime_session_id": "boi9ar543dg91ipdnspi099u231280ab0v8af0ew",
	"business_partner": 201,
	"filepath": "/var/lib/aion/Data/rededge_sdc/abcdef.json",
	"service_label": "ORDERS",
	"QuantityUnit": {
		"QuantityUnit": "KG"
	},
	"api_schema": "DPFMOrdersCreates",
	"accepter": ["Header"],
	"order_id": 1,
	"deleted": false
}
```

## Output
data-platform-api-quantity-unit-exconf-rmq-kube では、[golang-logging-library-for-data-platform](https://github.com/latonaio/golang-logging-library-for-data-platform) により、Output として、RabbitMQ へのメッセージを JSON 形式で出力します。数量単位の対象値が存在する場合 true、存在しない場合 false、を返します。"cursor" ～ "time"は、golang-logging-library-for-data-platform による 定型フォーマットの出力結果です。

```
{
	"cursor": "/Users/latona2/bitbucket/data-platform-api-quantity-unit-exconf-rmq-kube/main.go#L69",
	"function": "main.dataCallProcess",
	"level": "INFO",
	"message": {
		"QuantityUnit": {
			"QuantityUnit": "KG",
			"ExistenceConf": true
		}
	},
	"runtime_session_id": "boi9ar543dg91ipdnspi099u231280ab0v8af0ew",
	"time": "2022-11-17T15:14:23+09:00"
}
```

