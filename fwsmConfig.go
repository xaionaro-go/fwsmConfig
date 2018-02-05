package fwsmConfig

import (
	"io"
	"sort"
)

type FwsmConfig struct {
	DHCP   DHCPCommon
	VLANs  VLANs
	ACLs   ACLs
	SNATs  SNATs
	DNATs  DNATs
	Routes Routes
}

func (cfg FwsmConfig) WriteTo(writer io.Writer) error {
	sort.Sort(cfg.VLANs)
	sort.Sort(cfg.ACLs)
	sort.Sort(cfg.SNATs)
	sort.Sort(cfg.DNATs)
	sort.Sort(cfg.Routes)

	err := cfg.DHCP.WriteTo(writer)
	if err != nil {
		return err
	}

	for _, vlan := range cfg.VLANs {
		err := vlan.WriteTo(writer)
		if err != nil {
			return err
		}
	}
	for _, acl := range cfg.ACLs {
		err := acl.WriteTo(writer)
		if err != nil {
			return err
		}
	}
	for _, snat := range cfg.SNATs {
		err := snat.WriteTo(writer)
		if err != nil {
			return err
		}
	}
	for _, dnat := range cfg.DNATs {
		err := dnat.WriteTo(writer)
		if err != nil {
			return err
		}
	}
	for _, route := range cfg.Routes {
		err := route.WriteTo(writer)
		if err != nil {
			return err
		}
	}

	return nil
}
