# update apt-get
apt update

# install wget & curl
apt install wget
apt install curl

# install vim
apt install vim

# install git
apt install git

# download kafka 2.11-1.0.1
curl -O http://apache.claz.org/kafka/1.0.1/kafka_2.11-1.0.1.tgz

# install kafka
tar -xzf kafka_2.11-1.0.1.tgz && mv kafka_2.11-1.0.1/ kafka/

# download Golang
curl -O https://dl.google.com/go/go1.10.linux-amd64.tar.gz

# install Golang
tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz

# add Golang PATH Variable
export PATH=$PATH:/usr/local/go/bin

# install Golang package dependencies
go get github.com/Shopify/sarama
go get github.com/jasonlvhit/gocron

# install java 1.8
apt install openjdk1.8-jre

# Done
echo "All Done!" 
