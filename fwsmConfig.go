package fwsmConfig

import (
	"encoding/json"
	"fmt"
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

func (cfg FwsmConfig) CiscoString() (result string) {
	result += cfg.DHCP.CiscoString()
	result += cfg.VLANs.CiscoString()
	result += cfg.ACLs.CiscoString()
	result += cfg.SNATs.CiscoString()
	result += cfg.DNATs.CiscoString()
	result += cfg.Routes.CiscoString()
	result += cfg.DHCPs.CiscoString()
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
