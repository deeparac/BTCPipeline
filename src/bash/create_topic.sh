# define kafka path
KAFKA="/arac/tools/kafka"

# check topic name
set -e
if [ "$#" -ne 1 ]; then
	echo "Illegal number of arguments! Please enter only one string as a valid topic name."
	exit 1
fi

# create
$KAFKA/bin/kafka-topics.sh --create \
	--zookeeper localhost:2181 \
	--replication-factor 1 \
	--partitions 1 \
	--topic $1

# done
echo "Successfully created topic '$1'."
