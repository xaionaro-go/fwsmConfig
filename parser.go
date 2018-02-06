package fwsmConfig

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func Parse(reader io.Reader) (cfg FwsmConfig, err error) {
	vlanIndexMap := map[int]*VLAN{}
	vlanNameMap := map[string]*VLAN{}
	aclMap := map[string]*ACL{}
	snatMap := map[string]*SNAT{}

	globalNatMap := map[int]net.IP{}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(strings.Trim(line, " "), " ")
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "interface":
			if len(words[1]) < 5 {
				panic(fmt.Errorf("invalid interface name: %v; should be vlanX[X[X[X]]]", words[1]))
			}
			vlan := VLAN{Interface: net.Interface{MTU: 1500, Flags: net.FlagUp | net.FlagMulticast}}
			vlan.Index, err = strconv.Atoi(words[1][4:]) // "Vlan10" -> vlan.Index: 10
			if err != nil {
				return
			}
			for scanner.Scan() {
				subLine := scanner.Text()
				if strings.Trim(subLine, " \t\r\n") == "!" {
					break
				}
				subWords := strings.Split(subLine, " ")[1:]
				switch subWords[0] {
				case "nameif":
					vlan.Name = subWords[1]
				case "security-level":
					vlan.SecurityLevel, err = strconv.Atoi(subWords[1])
				case "ip":
					switch subWords[1] {
					case "address":
						vlan.IPs = append(vlan.IPs, net.ParseIP(subWords[2]))
						var ipnet IPNet
						ipnet, err = parseIPNet(subWords[2], subWords[3])
						vlan.AttachedNetworks = append(vlan.AttachedNetworks, ipnet)
					default:
						warning("Cannot parse line: %v", subLine)
					}
				case "shutdown":
					vlan.Flags &= 0 ^ net.FlagUp
				case "no":
				default:
					warning("Cannot parse line: %v", subLine)
				}
				if err != nil {
					return
				}
			}

			cfg.VLANs = append(cfg.VLANs, &vlan)
			vlanIndexMap[vlan.Index] = &vlan
			vlanNameMap[vlan.Name] = &vlan

		/*case "dns":
		switch words[1] {
		case "name-server":
			cfg.DHCP.NSs = append(cfg.DHCP.NSs, parseNS(words[2]))
		default:
			warning("Cannot parse line: %v", line)
		}*/
		case "access-list":
			aclName := words[1]
			acl := aclMap[aclName]
			isToAppend := false
			if acl == nil {
				isToAppend = true
				acl = &ACL{Name: aclName}
			}

			err = acl.ParseAppendRule(words[2:])

			if isToAppend {
				cfg.ACLs = append(cfg.ACLs, acl)
				aclMap[acl.Name] = acl
			}
		case "mtu":
			ifaceName := words[1]
			var mtu int
			mtu, err = strconv.Atoi(words[2])
			if err != nil {
				return
			}

			vlanNameMap[ifaceName].MTU = mtu

		case "global":
			var globalNatId int
			globalNatId, err = strconv.Atoi(words[2])
			if err != nil {
				return
			}
			globalNatMap[globalNatId] = net.ParseIP(words[3])

		case "nat":
			var globalNatId int
			globalNatId, err = strconv.Atoi(words[2])
			if err != nil {
				return
			}
			natTo := globalNatMap[globalNatId]
			if natTo.String() == "" || natTo.String() == "<nil>" {
				continue
			}
			snat := snatMap[natTo.String()]
			isToAppend := false
			if snat == nil {
				isToAppend = true
				snat = &SNAT{}
			}

			var source IPNet
			source, err = parseIPNet(words[3], words[4])
			if err != nil {
				return
			}
			snat.Sources = append(snat.Sources, SNATSource{IPNet: source, IfName: strings.Trim(words[1], "()")})
			snat.NATTo = natTo

			if isToAppend {
				snat.FWSMGlobalId = globalNatId
				cfg.SNATs = append(cfg.SNATs, snat)
				snatMap[natTo.String()] = snat
			}

		case "static":
			unusedWords := words[2:]

			dnat := DNAT{}
			protocol := PROTO_IP
			if strings.Index(unusedWords[0], ".") == -1 {
				protocol = parseProtocol(unusedWords[0])
				unusedWords = unusedWords[1:]
			}
			dstHost, dstPort, unusedWords := parseHostPort(unusedWords)
			natToHost, natToPort, unusedWords := parseHostPort(unusedWords)
			if unusedWords[0] != "netmask" {
				panic("This shouldn't happened")
			}
			if unusedWords[1] != "255.255.255.255" {
				panic("This case is not implemented, yet")
			}
			dnat.Destinations = append(dnat.Destinations, IPPort{Protocol: &protocol, IP: dstHost, Port: dstPort})
			dnat.NATTo = IPPort{IP: natToHost, Port: natToPort}

			cfg.DNATs = append(cfg.DNATs, &dnat)

		case "access-group":
			if words[2] != "in" || words[3] != "interface" {
				panic("This case is not implemented")
			}

			aclName := words[1]
			ifName := words[4]

			acl := aclMap[aclName]

			acl.VLANNames = append(acl.VLANNames, ifName)

		case "route":
			dstNet, err := parseIPNet(words[2], words[3])
			if err != nil {
				panic(err)
			}
			gw := net.ParseIP(words[4])
			metric, err := strconv.Atoi(words[5])
			if err != nil {
				panic(err)
			}

			cfg.Routes = append(cfg.Routes,
				&Route{
					Sources:     IPNets{IPNet{IP: net.ParseIP("0.0.0.0"), Mask: net.IPv4Mask(0, 0, 0, 0)}},
					Destination: dstNet,
					Gateway:     gw,
					Metric:      metric,
				},
			)

		case "dhcpd":
			switch words[1] {
			case "address":
				dhcp := DHCP{}
				err := dhcp.ParseRange(words[2])
				if err != nil {
					panic(err)
				}
				cfg.DHCPs = append(cfg.DHCPs, dhcp)

			case "dns":
				for _, ns := range words[2:] {
					cfg.DHCP.NSs = append(cfg.DHCP.NSs, parseNS(ns))
				}

			case "domain":
				cfg.DHCP.Domain = Domain(words[2])

			case "option":
				dhcpOptionId, err := strconv.Atoi(words[2])
				if err != nil {
					panic(err)
				}
				dhcpOptionValueType := parseDHCPOptionValueType(words[3])
				cfg.DHCP.Options = append(cfg.DHCP.Options, DHCPOption{
					Id:        dhcpOptionId,
					ValueType: dhcpOptionValueType,
					Value:     []byte(words[4]),
				})

			case "enable", "lease", "ping_timeout": // is ignored, ATM
			default:
				panic(fmt.Errorf("Unknown DHCP command: %v", words))
			}

		default:
			warning("Cannot parse line: %v", line)
		}
		if err != nil {
			return
		}

		//fmt.Println(words)
	}
	err = scanner.Err()
	return
}
