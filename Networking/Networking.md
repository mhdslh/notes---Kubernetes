1- [Container Networking From Scratch - Kristen Jacobs, Oracle](https://www.youtube.com/watch?v=6v_BDHIgOY8)

2- [Tutorial: Communication Is Key - Understanding Kubernetes Networking - Jeff Poole, Vivint Smart Home](https://www.youtube.com/watch?v=InZVNuKY5GY)

3a- Operating system: [Lecturs](https://www.youtube.com/watch?v=QEZKgWqwMeA&list=PLzBbfbHQmjyuqBFJ8KpDdcvnLNkvPXbS-) + [Homepage](https://pages.cs.wisc.edu/~remzi/Classes/537/Fall2013/)
3b- [Kubernetes Networking Series](https://www.youtube.com/playlist?list=PLSAko72nKb8QWsfPpBlsw-kOdMBD7sra-)

4- [Cloud Networking](https://www.coursera.org/learn/cloud-networking)

To Do:\
vlan => vxlan\
iptables (nat)\
bgp, netfilter subsystem, socket programming, loopback vs physical interface

**Appendix A**:
![image](https://github.com/mhdslh/notes---Kubernetes/assets/61638154/5859cb96-9af6-4294-9f0a-06738b0e5b8e)


**Appendix B: Life of a Packet**
![image](https://github.com/mhdslh/notes---Kubernetes/assets/61638154/6f4e3323-67e9-4157-8a11-d610d0fb75cc)
- Remark1: Source and Destination MAC addresses of a packet are updated hop by hop. The source and destination IP addresses always remain unchanged end to end.

- Remark2: DNS resolves FQDN (fully qualified domain names) to IP addresses. DNS cache is a local storage of DNS records (stores IP addresses of previously visited FQDNs). ARP cache entries are created when an IP address is resolved to a MAC address.

- Remark 3: If the destination IP address is not within the local subnet, the packet will be forwarded to its default gateway by setting the destination MAC address to the default gateway's MAC address. 

**Appendix C: Routing**

- Remark1: Each entry in the routing table typically includes information about the destination network, the next hop or outgoing interface to use for reaching that destination, and sometimes additional parameters such as the cost or metric associated with each route. Next hop is always directly reachable by a connected interface. The routing table is updated dynamically through routing protocols (dynamic routes) or manually configured by network administrators to ensure efficient and reliable routing of network traffic (static routes).

- Remark2: The gateway of last resort is the network point that a router uses to forward packets when it doesn't have specific instructions about a destination address in its routing table. Essentially, it's the "fallback" route or last-resort route. (Do not confuse it with default gateway. Default Gateway is commonly used in the context of host and client configurations. It's the "way out" of the local network.)

- Remark3: If the routing table contains multiple entries that match the destination IP, the router uses the longest prefix match rule. This means that the router selects the entry that has the most number of leading bits in common with the destination IP address (i.e., most specific route).  When a router has multiple routes, it can distribute the data packets among these routes to balance the load and maximize throughput, minimize response time and avoid overload of any single path.
