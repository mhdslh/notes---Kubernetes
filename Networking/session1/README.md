![image](https://github.com/mhdslh/notes---Kubernetes/assets/61638154/bfdf82f0-52b5-4d2c-8334-412ce685870f)

Two ways to enable container networking:

1- Advertising containers' routes to nodes and gateways sitting in between, in nodes are in different L3 networks. In our setup, we can enable this by running `sudo ip route add $TO_NETWORK/24 via $TO_NODE dev enp0s8` on each node.

2- The problem with the previous approach is that we need to advertise the routes for all container networks to all the nodes and gateways. We can overcome this by overlay networks and tunneling. In our setup, we can enable this by running `sudo socat tun:$TUN/16,iff-up udp:$TO_NODE:9000,bind=$NODE:9000` on both nodes. First, note that when a packet is sent through tunnel it first will be decapsulate on the other node by socat, when arrives at TO_NODE interface. Then it will be sent to tun0 since it is handleing (listening for) all the traffics arriving at TO_NODE:9000. **/16 subnet in tunnel address provides connectivity to both container networks on the local node and remote.** All container networks are subnets of 172.16.0.0/16. On node1 we are using 172.16.1.0/24 range for containers and on node2 we are using 172.16.2.0/24. If a packet is not destined for the local container network (due to longer prefix match), it will be forwarded to the other node.

Appendix:
socat establishes bidirectional data transfer between two data channels. The first argument configures the listener and the second argument configures the other endpoint. For example, `socat tcp-listen:8080 tcp:localhost:8080' forwards any data received on port 8080 to port 8080 on the same machine, i.e., loopback interface. We can use bind option for listener to bind the listening port to a specific network interface. Similarly, bind option can be used on the remote endpoint to set the local endpoint that is used for outbound traffic going to the remote endpoint. When traffic arrives at the remote endpoint it needs to be handled by a process.

References:
  1- [kristenjacobs-container-networking](https://github.com/kristenjacobs/container-networking)

  2- [Understanding Kubernetes Networking. Part 1: Container Networking](https://www.youtube.com/watch?v=B6FsWNUnRo0&list=PLSAko72nKb8QWsfPpBlsw-kOdMBD7sra-&index=1)
