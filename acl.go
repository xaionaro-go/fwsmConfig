package fwsmConfig

import (
	"net"
	"io"
)

type ACLAction int

const (
	ACL_ALLOW = ACLAction(0)
	ACL_DENY  = ACLAction(1)
)

type ACLRule struct {
	Action   ACLAction
	Protocol Protocol
	From     net.IPNet
	To       net.IPNet
}

type ACLRules []ACLRule

type ACL struct {
	Name      string
	Rules     ACLRules
	VLANNames []string
}

type ACLs []ACL

func (a ACLs) Len() int           { return len(a) }
func (a ACLs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ACLs) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (acl ACL) WriteTo(writer io.Writer) error {
	return nil
}
