#! /bin/bash

IP1=172.16.2.11
IP2=172.16.2.12
GATEWAY=172.16.2.1
TUN=172.16.2.100
TO_NETWORK=172.16.1.0
TO_NODE=10.0.0.11

sudo ip netns add NS1
sudo ip link add veth10 type veth peer name veth11
sudo ip link set veth11 netns NS1
sudo ip netns exec NS1 ip addr add $IP1/24 dev veth11
sudo ip netns exec NS1 ip link set veth11 up
sudo ip netns exec NS1 ip link set lo up

sudo ip netns add NS2
sudo ip link add veth20 type veth peer name veth21
sudo ip link set veth21 netns NS2
sudo ip netns exec NS2 ip addr add $IP2/24 dev veth21
sudo ip netns exec NS2 ip link set veth21 up
sudo ip netns exec NS2 ip link set lo up

sudo ip link add br0 type bridge
sudo ip addr add $GATEWAY/24 dev br0
sudo ip link set br0 up
sudo ip link set veth10 master br0
sudo ip link set veth10 up
sudo ip link set veth20 master br0
sudo ip link set veth20 up

sudo ip netns exec NS1 ip route add default via $GATEWAY dev veth11
sudo ip netns exec NS2 ip route add default via $GATEWAY dev veth21

sudo sysctl -w net.ipv4.ip_forward=1

#----
# sudo ip route add $TO_NETWORK/24 via $TO_NODE dev enp0s8
# sudo socat tun:$TUN/16,iff-up udp:$TO_NODE:9000,bind=$NODE:9000

