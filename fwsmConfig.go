package fwsmConfig

import (
	"encoding/json"
	"io"
	"sort"
)

type FwsmConfig struct {
	DHCP   DHCPCommon
	VLANs  VLANs
	ACLs   ACLs
	SNATs  SNATs
	DNATs  DNATs
	DHCPs  DHCPs
	Routes Routes
}

func (cfg *FwsmConfig) prepareToWrite() {
	sort.Sort(cfg.VLANs)
	sort.Sort(cfg.ACLs)
	sort.Sort(cfg.SNATs)
	sort.Sort(cfg.DNATs)
	sort.Sort(cfg.Routes)
}

func (cfg FwsmConfig) WriteTo(writer io.Writer) error {
	cfg.prepareToWrite()

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
	for _, dhcp := range cfg.DHCPs {
		err := dhcp.WriteTo(writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cfg FwsmConfig) WriteJsonTo(writer io.Writer) (err error) {
	cfg.prepareToWrite()
	jsonEncoder := json.NewEncoder(writer)
	jsonEncoder.SetIndent("", "  ")
	return jsonEncoder.Encode(cfg)
}
