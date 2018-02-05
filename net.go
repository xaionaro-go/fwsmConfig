package fwsmConfig

import (
	"net"
	"strconv"
)

type IPPort struct {
	IP   net.IP
	Port *uint16
}

type IPs []net.IP
type IPNets []net.IPNet
type IPPorts []IPPort
type NSs []net.NS

func (ipport IPPort) String() string {
	if ipport.Port == nil {
		return ipport.IP.String()
	}

	return ipport.IP.String()+":"+strconv.Itoa(int(*ipport.Port))
}

func parseIPNet(ipStr string, maskStr string) (ipnet net.IPNet, err error) {
	ip := net.ParseIP(ipStr)
	ipnet.Mask = net.IPMask(net.ParseIP(maskStr))
	ipnet.IP = ip.Mask(ipnet.Mask)

	return
}

func parseNS(nsStr string) net.NS {
	return net.NS{Host: nsStr}
}

type Protocol int

// awk 'BEGIN{prevId=-1} {if($1 == "#" || $1 == "" || $2 <= prevId){next} gsub("[-.]", "", $1) ;printf "%s", "PROTO_"toupper($1)" = Protocol("$2") // "; $1=""; prevId=$2; $2=""; print $0}' < /etc/protocols

const (
	PROTO_IP             = Protocol(0)   // IP # internet protocol, pseudo protocol number
	PROTO_ICMP           = Protocol(1)   // ICMP # internet control message protocol
	PROTO_IGMP           = Protocol(2)   // IGMP # Internet Group Management
	PROTO_GGP            = Protocol(3)   // GGP # gateway-gateway protocol
	PROTO_IPENCAP        = Protocol(4)   // IP-ENCAP # IP encapsulated in IP (officially ``IP'')
	PROTO_ST             = Protocol(5)   // ST # ST datagram mode
	PROTO_TCP            = Protocol(6)   // TCP # transmission control protocol
	PROTO_EGP            = Protocol(8)   // EGP # exterior gateway protocol
	PROTO_IGP            = Protocol(9)   // IGP # any private interior gateway (Cisco)
	PROTO_PUP            = Protocol(12)  // PUP # PARC universal packet protocol
	PROTO_UDP            = Protocol(17)  // UDP # user datagram protocol
	PROTO_HMP            = Protocol(20)  // HMP # host monitoring protocol
	PROTO_XNSIDP         = Protocol(22)  // XNS-IDP # Xerox NS IDP
	PROTO_RDP            = Protocol(27)  // RDP # "reliable datagram" protocol
	PROTO_ISOTP4         = Protocol(29)  // ISO-TP4 # ISO Transport Protocol class 4 [RFC905]
	PROTO_DCCP           = Protocol(33)  // DCCP # Datagram Congestion Control Prot. [RFC4340]
	PROTO_XTP            = Protocol(36)  // XTP # Xpress Transfer Protocol
	PROTO_DDP            = Protocol(37)  // DDP # Datagram Delivery Protocol
	PROTO_IDPRCMTP       = Protocol(38)  // IDPR-CMTP # IDPR Control Message Transport
	PROTO_IPV6           = Protocol(41)  // IPv6 # Internet Protocol, version 6
	PROTO_IPV6ROUTE      = Protocol(43)  // IPv6-Route # Routing Header for IPv6
	PROTO_IPV6FRAG       = Protocol(44)  // IPv6-Frag # Fragment Header for IPv6
	PROTO_IDRP           = Protocol(45)  // IDRP # Inter-Domain Routing Protocol
	PROTO_RSVP           = Protocol(46)  // RSVP # Reservation Protocol
	PROTO_GRE            = Protocol(47)  // GRE # General Routing Encapsulation
	PROTO_ESP            = Protocol(50)  // IPSEC-ESP # Encap Security Payload [RFC2406]
	PROTO_AH             = Protocol(51)  // IPSEC-AH # Authentication Header [RFC2402]
	PROTO_SKIP           = Protocol(57)  // SKIP # SKIP
	PROTO_IPV6ICMP       = Protocol(58)  // IPv6-ICMP # ICMP for IPv6
	PROTO_IPV6NONXT      = Protocol(59)  // IPv6-NoNxt # No Next Header for IPv6
	PROTO_IPV6OPTS       = Protocol(60)  // IPv6-Opts # Destination Options for IPv6
	PROTO_RSPF           = Protocol(73)  // RSPF CPHB # Radio Shortest Path First (officially CPHB)
	PROTO_VMTP           = Protocol(81)  // VMTP # Versatile Message Transport
	PROTO_EIGRP          = Protocol(88)  // EIGRP # Enhanced Interior Routing Protocol (Cisco)
	PROTO_OSPF           = Protocol(89)  // OSPFIGP # Open Shortest Path First IGP
	PROTO_AX25           = Protocol(93)  // AX.25 # AX.25 frames
	PROTO_IPIP           = Protocol(94)  // IPIP # IP-within-IP Encapsulation Protocol
	PROTO_ETHERIP        = Protocol(97)  // ETHERIP # Ethernet-within-IP Encapsulation [RFC3378]
	PROTO_ENCAP          = Protocol(98)  // ENCAP # Yet Another IP encapsulation [RFC1241]
	PROTO_PIM            = Protocol(103) // PIM # Protocol Independent Multicast
	PROTO_IPCOMP         = Protocol(108) // IPCOMP # IP Payload Compression Protocol
	PROTO_VRRP           = Protocol(112) // VRRP # Virtual Router Redundancy Protocol [RFC5798]
	PROTO_L2TP           = Protocol(115) // L2TP # Layer Two Tunneling Protocol [RFC2661]
	PROTO_ISIS           = Protocol(124) // ISIS # IS-IS over IPv4
	PROTO_SCTP           = Protocol(132) // SCTP # Stream Control Transmission Protocol
	PROTO_FC             = Protocol(133) // FC # Fibre Channel
	PROTO_MOBILITYHEADER = Protocol(135) // Mobility-Header # Mobility Support for IPv6 [RFC3775]
	PROTO_UDPLITE        = Protocol(136) // UDPLite # UDP-Lite [RFC3828]
	PROTO_MPLSINIP       = Protocol(137) // MPLS-in-IP # MPLS-in-IP [RFC4023]
	PROTO_MANET          = Protocol(138) // # MANET Protocols [RFC5498]
	PROTO_HIP            = Protocol(139) // HIP # Host Identity Protocol
	PROTO_SHIM6          = Protocol(140) // Shim6 # Shim6 Protocol [RFC5533]
	PROTO_WESP           = Protocol(141) // WESP # Wrapped Encapsulating Security Payload
	PROTO_ROHC           = Protocol(142) // ROHC # Robust Header Compression
)
