package fwsmConfig

import (
	"encoding/json"
	"fmt"
	"github.com/xaionaro-go/networkControl"
	"io"
	"os"
	"sort"
)

type FwsmConfig struct {
	DHCP   DHCP
	VLANs  VLANs
	ACLs   ACLs
	SNATs  SNATs
	DNATs  DNATs
	Routes Routes
}

func NewFwsmConfig() *FwsmConfig {
	return &FwsmConfig{DHCP: *NewDHCP()}
}

func (cfg *FwsmConfig) prepareToWrite() {
	sort.Sort(cfg.VLANs)
	sort.Sort(networkControl.ACLs(cfg.ACLs))
	sort.Sort(networkControl.SNATs(cfg.SNATs))
	sort.Sort(networkControl.DNATs(cfg.DNATs))
	sort.Sort(networkControl.Routes(cfg.Routes))
}

func (cfg FwsmConfig) CiscoString() (result string) {
	result += cfg.VLANs.CiscoString()
	result += cfg.DHCP.CiscoString(cfg.VLANs)
	result += cfg.ACLs.CiscoString()
	result += cfg.SNATs.CiscoString()
	result += cfg.DNATs.CiscoString()
	result += cfg.Routes.CiscoString()
	return
}

func (cfg FwsmConfig) WriteTo(writer io.Writer) error {
	cfg.prepareToWrite()

	fmt.Fprintf(writer, cfg.CiscoString())
	return nil
}

func (cfg FwsmConfig) WriteJsonTo(writer io.Writer) (err error) {
	cfg.prepareToWrite()
	jsonEncoder := json.NewEncoder(writer)
	jsonEncoder.SetIndent("", "  ")
	return jsonEncoder.Encode(cfg)
}


func (cfg FwsmConfig) ToNetworkControlState() networkControl.State {
	return networkControl.State{
		DHCP:         networkControl.DHCP(cfg.DHCP),
		BridgedVLANs: cfg.VLANs.ToNetworkControlVLANs(),
		ACLs:         networkControl.ACLs(cfg.ACLs),
		SNATs:        networkControl.SNATs(cfg.SNATs),
		DNATs:        networkControl.DNATs(cfg.DNATs),
		Routes:       networkControl.Routes(cfg.Routes),
	}
}

func (cfg FwsmConfig) Apply(netHost networkControl.HostI) error {
	err := netHost.SetNewState(cfg.ToNetworkControlState())
	if err != nil {
		return err
	}
	return netHost.Apply()
}

func (cfg FwsmConfig) Save(netHost networkControl.HostI, cfgPath string) error {
	f, err := os.Create(cfgPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = cfg.WriteTo(f)
	if err != nil {
		return err
	}

	return netHost.Save()
}

func (cfg FwsmConfig) Revert(netHost networkControl.HostI) error {
	panic(errNotImplemented)
	return errNotImplemented
}
