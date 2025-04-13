<h1>
<p align="center">
  <br>transaction-outbox-pattern
</h1>
</p>

![transaction_outbox_pattern](https://github.com/user-attachments/assets/f7da1c25-f260-4549-9539-b785835d42cb)

### Developing

* Gin
* Kafka
* Kafka-connect
* Debezium Mysql Connector
* Mysql
* Postgres
* Docker

### How to perform a test

#### You need to have [Docker](https://www.docker.com/) installed.

```bash

# Clone this repository and access project folder
$ git clone https://github.com/Daiven75/transaction-outbox-pattern.git && cd transaction-outbox-pattern

# Run the command in the terminal to go up all stack
$ docker-compose up -d

# The services are running on port:
# - flight-service - http:localhost:8888
# - passenger-service - http:localhost:8887
# - kafka-connect - http:localhost:8083

```

#### Creating a flight
```bash
curl -iX POST "http://localhost:8888/flights" -d "{\"company\":\"Company-1\",\"origin\":\"ABC\",\"destination\":\"DEF\"}"
```

#### Creating and linking a passenger to the flight
```bash
curl -iX POST "http://localhost:8887/passengers" -d "{\"first_name\":\"passenger-1\",\"plan_type\":\"DEFAULT\",\"dispatch\":\"false\",\"flight_id\":\1\}"
```

#### Searching all flights
```bash
curl -i "http://localhost:8888/flights"
```

#### [PLUS] Kafka connect provides an API, and in it, we can query the previously installed connector
```bash
# Search for all installed connectors
curl -i "http://localhost:8083/connectors"

# Displays existing connector configuration details
curl -i "http://localhost:8083/connectors/mysql-source-connector"
```

## How about we chat?

Made by Lucas Silva 👋🏽 [Contact](https://www.linkedin.com/in/lucas-silva-959102169)

[![Gmail Badge](https://img.shields.io/badge/-75.lucas.slima@gmail.com-c14438?style=flat-square&logo=Gmail&logoColor=white&link=mailto:75.lucas.slima@gmail.com)](mailto:75.lucas.slima@gmail.com)
