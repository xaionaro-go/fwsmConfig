package fwsmConfig

import (
	"bufio"
	"fmt"
	"github.com/xaionaro-go/networkControl"
	"io"
	"net"
	"strconv"
	"strings"
)

func Parse(reader io.Reader) (cfg FwsmConfig, err error) {
	cfg = *NewFwsmConfig()

	vlanVlanIdMap := map[int]*VLAN{}
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
			vlan.VlanId, err = strconv.Atoi(words[1][4:]) // "Vlan10" -> vlan.VlanId: 10
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
						/*vlan.IPs = append(vlan.IPs, net.ParseIP(subWords[2]))
						var ipnet IPNet
						ipnet, err = parseIPNet(subWords[2], subWords[3])
						vlan.AttachedNetworks = append(vlan.AttachedNetworks, ipnet)*/
						var ipnet networkControl.IPNet
						ipnet, err = parseIPNetUnmasked(subWords[2], subWords[3])
						vlan.IPs = append(vlan.IPs, ipnet)
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

			if vlan.Name == "" {
				warning("Empty interface name of vlan: %v", vlan)
				continue
			}

			cfg.VLANs = append(cfg.VLANs, &vlan)
			vlanVlanIdMap[vlan.VlanId] = &vlan
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
				cfg.ACLs = append(cfg.ACLs, (*networkControl.ACL)(acl))
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

			var source networkControl.IPNet
			source, err = parseIPNet(words[3], words[4])
			if err != nil {
				return
			}
			snat.Sources = append(snat.Sources, networkControl.SNATSource{IPNet: source, IfName: strings.Trim(words[1], "()")})
			snat.NATTo = natTo

			if isToAppend {
				snat.FWSMGlobalId = globalNatId
				cfg.SNATs = append(cfg.SNATs, (*networkControl.SNAT)(snat))
				snatMap[natTo.String()] = snat
			}

		case "static":
			unusedWords := words[2:]

			dnat := DNAT{}
			protocol := networkControl.PROTO_IP
			if strings.Index(unusedWords[0], ".") == -1 {
				protocol = parseProtocol(unusedWords[0])
				unusedWords = unusedWords[1:]
			}
			dstHost, dstPort, unusedWords := parseHostPort(unusedWords)
			natToHost, natToPort, unusedWords := parseHostPort(unusedWords)
			if unusedWords[0] != "netmask" {
				panic(fmt.Errorf("This shouldn't happened: %v", words))
			}
			if unusedWords[1] != "255.255.255.255" {
				panic(fmt.Errorf("This case is not implemented, yet: %v", words))
			}
			dnat.Destinations = append(dnat.Destinations, networkControl.IPPort{Protocol: &protocol, IP: dstHost, Port: dstPort})
			dnat.NATTo = networkControl.IPPort{Protocol: &protocol, IP: natToHost, Port: natToPort}
			dnat.IfName = strings.Split(strings.Trim(words[1], "()"), ",")[0]

			cfg.DNATs = append(cfg.DNATs, (*networkControl.DNAT)(&dnat))

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
				&networkControl.Route{
					Sources:     networkControl.IPNets{networkControl.IPNet{IP: net.ParseIP("0.0.0.0"), Mask: net.IPv4Mask(0, 0, 0, 0)}},
					Destination: dstNet,
					Gateway:     gw,
					Metric:      metric,
					IfName:      words[1],
				},
			)

		case "dhcpd":
			switch words[1] {
			case "address":
				subnet := DHCPSubnet{}
				err := (*DHCPRange)(&subnet.Options.Range).Parse(words[2])
				if err != nil {
					panic(err)
				}
				iface := vlanNameMap[words[3]]
				if iface == nil {
					panic(fmt.Errorf("Cannot find interface %v", words[3]))
				}
				for _, ip := range iface.IPs {
					if !ip.Contains(subnet.Options.Range.Start) {
						continue
					}
					subnet.Network = net.IPNet(ip)
					subnet.Network.IP = subnet.Network.IP.Mask(subnet.Network.Mask)
				}
				err = cfg.DHCP.Subnets.ISet(subnet)
				if err != nil {
					panic(err)
				}

			case "dns":
				cfg.DHCP.Options.DomainNameServers.Set(words[2:])

			case "domain":
				cfg.DHCP.Options.DomainName = words[2]

			case "option":
				dhcpOptionId, err := strconv.Atoi(words[2])
				if err != nil {
					panic(err)
				}
				value, valueType := parseDHCPOptionValue(words[3], words[4])
				cfg.DHCP.Options.Custom[dhcpOptionId] = []byte(value)

				cfg.DHCP.UserDefinedOptionFields.Set(fmt.Sprintf("option%v", dhcpOptionId), dhcpOptionId, valueType)

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
