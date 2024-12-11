FROM confluentinc/cp-kafka-connect:7.7.1

# Installing Connector
RUN confluent-hub install debezium/debezium-connector-mysql:latest --no-prompt

# Connector Properties
COPY debezium-mysql-config.json /etc/kafka/debezium-mysql-config.json

# Launching Kafka Connect Worker and Registering Mysql Connector
CMD ["/bin/bash", "-c", "/etc/confluent/docker/run & sleep 40 && curl -X POST -H \"Content-Type: application/json\" --data @/etc/kafka/debezium-mysql-config.json http://localhost:8083/connectors && wait"]
