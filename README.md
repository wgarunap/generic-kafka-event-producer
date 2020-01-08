# generic-avro-producer

    <HOST>:<PORT>/publish

```json
{
	"topic":"test",
	"subject":"com.company.events.budget.BudgetChanged",
	"version": 2,
	"headers":{
		"subject":"com.company.events.budget.BudgetChanged",
		"account_id":"123e4567-e89b-12d3-a456-426655440000"
	},
	"key": "test",
	"value": {
	  "meta": {
	    "event_id": "123e4567-e89b-12d3-a456-426655440000",
	    "trace_id": "123e4567-e89b-12d3-a456-426655440000",
	    "account_id": "123e4567-e89b-12d3-a456-426655440000",

	  }
	}
}
```