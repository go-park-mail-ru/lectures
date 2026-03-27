#!/bin/sh

/opt/bitnami/kafka/bin/kafka-topics.sh --create --topic cool_topic --bootstrap-server localhost:9092
/opt/bitnami/kafka/bin/kafka-topics.sh --describe --topic cool_topic --bootstrap-server localhost:9092
