package fwsmConfig

import (
	"bytes"
	"io"
	"net"
)

type Route struct {
	Sources     IPNets
	Destination IPNet
	Gateway     net.IP
	Metric      int
}

type Routes []*Route

func (a Routes) Len() int           { return len(a) }
func (a Routes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Routes) Less(i, j int) bool { return a[i].GetPos() < a[j].GetPos() }

func (route Route) GetPos() string {
	return route.Gateway.String()
}

func (route Route) WriteTo(writer io.Writer) error {
	return nil
}

func (routes Routes) CiscoString() string {
	var buf bytes.Buffer
	for _, route := range routes {
		route.WriteTo(&buf)
	}
	return buf.String()
}

