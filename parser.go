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
	vlanMap := map[int]*VLAN{}
	aclMap := map[string]*ACL{}

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
						var ipnet net.IPNet
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

			vlanMap[vlan.Index] = cfg.VLANs.Append(vlan)
		case "dns":
			switch words[1] {
			case "name-server":
				cfg.DHCP.NSs = append(cfg.DHCP.NSs, parseNS(words[2]))
			default:
				warning("Cannot parse line: %v", line)
			}
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
				aclMap[acl.Name] = cfg.ACLs.Append(*acl)
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
