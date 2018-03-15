# specify kafka path
KAFKA_PATH="/arac/tools/kafka"

# start zookeeper server
echo "zookeeper server started, running in background"
$KAFKA_PATH/bin/zookeeper-server-start.sh $KAFKA_PATH/config/zookeeper.properties &

# start kafka server
echo "kafka server started"
$KAFKA_PATH/bin/kafka-server-start.sh $KAFKA_PATH/config/server.properties
