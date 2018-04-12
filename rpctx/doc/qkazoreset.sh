#!/usr/bin/env bash

echo "start reset kafka"

mkdir /root/resource
cd /root/resource

# backup kafka config
cp /work/kafka/config/server.properties /root/resource
cp /work/kafka/bin/kafka-server-start.sh /root/resource

wget https://mirrors.tuna.tsinghua.edu.cn/apache/kafka/0.11.0.2/kafka_2.11-0.11.0.2.tgz
rm -rf /work/kafka
mkdir -p /work/kafka
cd /work/kafka
tar -zxf /root/source/kafka_2.11-0.11.0.2.tgz --strip 1

# restore kafka config
cp /root/resource/server.properties /work/kafka/config
cp /root/resource/kafka-server-start.sh /work/kafka/config

echo "reset kafka complete"

echo "reset zookeeper start"

service zookeeper stop

rm -rf /work/zookeeper
mkdir -p /work/zookeeper
mkdir /work/zookeeper/version-2
touch /work/zookeeper/myid
chown -R zookeeper:zookeeper /work/zookeeper
echo 1 > /work/zookeeper/myid

service zookeeper start

echo "reset zookeeper complete"

