package fwsmConfig

import (
	"bytes"
	"fmt"
	"github.com/xaionaro-go/networkControl"
	"io"
	"sort"
	"net"
)

/*type VLAN struct {
	net.Interface
	VlanId        int
	SecurityLevel int
	IPs           IPNets
	//IPs              IPs
	//AttachedNetworks IPNets
}*/

type VLAN networkControl.VLAN

type VLANs []*VLAN

func (a VLANs) Len() int           { return len(a) }
func (a VLANs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VLANs) Less(i, j int) bool { return a[i].VlanId < a[j].VlanId }
func (a VLANs) Sort() VLANs { sort.Sort(a); return a }

func (vlan VLAN) WriteTo(writer io.Writer) error {
	fmt.Fprintf(writer, "interface Vlan%v\n", vlan.VlanId)
	fmt.Fprintf(writer, " nameif %v\n", vlan.Name)
	fmt.Fprintf(writer, " security-level %v\n", vlan.SecurityLevel)
	if len(vlan.IPs) > 0 {
		if len(vlan.IPs) > 1 {
			panic("Not implemented, yet")
		}
		ip := vlan.IPs[0]
		mask := net.IP(ip.Mask)
		fmt.Fprintf(writer, " ip address %v %v\n", ip.IP.String(), mask.String())
	} else {
		fmt.Fprintf(writer, " no ip address\n")
	}
	fmt.Fprintf(writer, "!\n")
	return nil
}

func (vlan VLAN) CiscoString() string {
	var buf bytes.Buffer
	vlan.WriteTo(&buf)
	return buf.String()
}
func (vlans VLANs) CiscoString() string {
	var buf bytes.Buffer
	for _, vlan := range vlans {
		vlan.WriteTo(&buf)
	}
	return buf.String()
}

func (vlans VLANs) ToNetworkControlVLANs() (result networkControl.VLANs) {
	result = networkControl.VLANs{}
	for _, vlan := range vlans {
		result[vlan.VlanId] = (*networkControl.VLAN)(vlan)
	}
	return
}

func (vlans VLANs) Find(vlanId int) (vlan VLAN, found bool) {
	for _, vlan := range vlans {
		if vlan.VlanId == vlanId {
			return *vlan, true
		}
	}

	return VLAN{}, false
}

func (vlans VLANs) Remove(netHost networkControl.HostI, vlanIds ...int) (err error) {
	todeleteMap := map[int]bool{}
	for _, vlanId := range vlanIds {
		todeleteMap[vlanId] = true
	}
	toRemoveIndexes := []int{}
	for idx, vlan := range vlans {
		if !todeleteMap[vlan.VlanId] {
			continue
		}

		if netHost != nil {
			err = netHost.RemoveBridgedVLAN(vlan.VlanId)
		}
		if err != nil {
			break
		}
		toRemoveIndexes = append(toRemoveIndexes, idx)
	}

	vlans = removeIndexes(vlans, toRemoveIndexes...).(VLANs)

	return err
}

func (vlans VLANs) Add(netHost networkControl.HostI, newVLANs ...VLAN) (err error) {
	alreadyExistsMap := map[int]bool{}
	for _, vlan := range vlans {
		alreadyExistsMap[vlan.VlanId] = true
	}

	for _, vlan := range newVLANs {
		if alreadyExistsMap[vlan.VlanId] {
			continue
		}

		if netHost != nil {
			err = netHost.AddBridgedVLAN(*networkControl.NewVLAN(vlan.Interface))
		}
		if err != nil {
			break
		}
		vlans = append(vlans, &vlan)
	}

	return err
}

/*func (vlans VLANs) ReconsiderSecurityLevels(netHost networkControl.HostI) (err error) {
	if netHost == nil {
		return nil
	}

	netHost.Firewall().ForwardIfaces_Flush()
	for _, vlanA := range vlans {
		for _, vlanB := range vlans {
			if vlanA.SecurityLevel >= vlanB.SecurityLevel {
				continue
			}
			err = netHost.Firewall().ForwardIfaces_AddReturn(vlanB.Interface, vlanA.Interface)
			if err != nil {
				return err
			}
		}
		err = netHost.Firewall().ForwardIfaces_AddReject(nil, vlanA.Interface)
		if err != nil {
			return err
		}
	}

	return nil
}
*/
