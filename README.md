# Generic Kafka Event Producer

### Overview
This application will provide REST endpoints to publish kafka messages with following encodings for key and Value.
   - String
   - AVRO
   - JSON
   - MessagePack

Users can set following parameters in the kafka message.
   - Topic
   - Headers(any number of key-value pairs)
   - Key(kafka message key)
   - Value(kafka message payload)

### Running the application
#### Build the docker image
```sh 
   docker build -f Dockerfile -t generic-producer .
```

#### Pull the repository image
```sh
   docker pull wgarunap/generic-kafka-event-producer:v1.3.0
```

#### Run generic event producer
```sh
docker run -dit -e SCHEMAREG_URL=http://schemareg.event.com:8081 -e KAFKA_BROKERS=kafka-1:9092,kafka-2:9092 -p 8000:8000 wgarunap/generic-kafka-event-producer:v1.3.0
```

### Configs 
```shell
+---------------+------------+----------+---------+---------------------------------+-----------------------------------------------------+
|      ENV      | IsRequired | Env Type | Default |             Example             |                     Description                     |
+---------------+------------+----------+---------+---------------------------------+-----------------------------------------------------+
| SCHEMAREG_URL | Optional   | String   | -       | http://schemareg.event.com:8081 | Schema registry url to fetch the avro event schemas |
| KAFKA_BROKERS | Required   | []String | -       | kafka-1,kafka-2,kafka-2         | kafka brokers list                                  |
| PORT          | Optional   | Integer  | 8000    | 8000                            | http serving port                                   |
| LOG_LEVEL     | Optional   | String   | INFO    | FATAL,WARN,INFO,DEBUG,TRACE     | Application log level                               |
| LOG_COLOR     | Optional   | Boolean  | false   | false                           | Application log color Enable/Disable                |
+---------------+------------+----------+---------+---------------------------------+-----------------------------------------------------+
```

### Avro Producer 
#### Request URL
```sh
   <HOST>:<PORT>/publish
```
#### Example POST Request Payload
```json
{
   "topic":"test",
   "headers":{
      "key1":"value1",
      "key2":"value2",
   },
   "key":{
      "subject":"com.event.KeySchemaSubject", // Only Required if key.format=avro
      "version":1, // Only Required if key.format=avro
      "format":"avro", // default is string, key.format can be one of avro,json,string,bytes,messagepack
      "body":{
         "event_id":"123e4567-e89b-12d3-a456-426655440000"
      }
   },
   "value":{
      "subject":"com.event.EventName", // Only Required if value.format=avro
      "version":2, // Only Required if value.format=avro
      "format":"avro", // default is string, value.format can be one of avro,json,string,bytes,messagepack
      "body":{
         "meta":{
            "event_id":"123e4567-e89b-12d3-a456-426655440000",
            "trace_id":"123e4567-e89b-12d3-a456-426655440000",
            "account_id":"123e4567-e89b-12d3-a456-426655440000"
         }
      }
   }
}
```

### JSON Producer 
#### Request URL
```sh
   <HOST>:<PORT>/publish
```
#### Example POST Request Payload
```json
{
   "topic":"test",
   "headers":{
      "account_id":"123e4567-e89b-12d3-a456-426655440000"
   },
   "key":{
      "format":"string", //Optional, Default is string
      "body": "test"
   },
   "value":{
      "format":"json",
      "body":{
         "meta":{
            "event_id":"123e4567-e89b-12d3-a456-426655440000",
            "trace_id":"123e4567-e89b-12d3-a456-426655440000",
            "account_id":"123e4567-e89b-12d3-a456-426655440000"
         }
      }
   }
}
```

### Null Producer 
#### Request URL
```sh
   <HOST>:<PORT>/publish
```
#### Example POST Request Payload
```json
{
   "topic":"test",
   "headers":{
      "subject":"com.event.EventName",
      "account_id":"123e4567-e89b-12d3-a456-426655440000"
   },
   "key":{
      "format":"string", // Optional
      "body":"test"
   },
}
```


### Get Kafka Topic List
#### Request URL
```sh
   curl -XGET "<HOST>:<PORT>/topics"
```

### Get Schema Reg Subject List
#### Request URL
```sh
   curl -XGET "<HOST>:<PORT>/avrosubjects"
```

### TODO
   - Add following key/value encoding methods
      - Protobuf
   - Add authentication layer
   - Proper error handling and error responses
   - Imporve documentation about usage on each message encoding type
   - Github actions implementation to automate testing
   - write unit tests and integration tests 
