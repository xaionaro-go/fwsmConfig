package fwsmConfig

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type ACLAction int

const (
	ACL_ALLOW = ACLAction(1)
	ACL_DENY  = ACLAction(2)
)

type ACLFlags uint16

const (
	ACLFL_ESTABLISHED = ACLFlags(0x01)
)

type ACLRule struct {
	Action         ACLAction
	Protocol       Protocol
	FromNet        IPNet
	FromPortRanges PortRanges
	ToNet          IPNet
	ToPortRanges   PortRanges
	Flags          ACLFlags
}

type ACLRules []ACLRule

type ACL struct {
	Name      string
	Rules     ACLRules
	VLANNames []string
}

type ACLs []*ACL

func (a ACLs) Len() int           { return len(a) }
func (a ACLs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ACLs) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (rule ACLRule) CiscoString() string {
	result := []string{}

	switch rule.Action {
	case ACL_ALLOW:
		result = append(result, "permit")
	case ACL_DENY:
		result = append(result, "deny")
	default:
		panic(fmt.Errorf("Invalid rule action: %v", rule.Action))
	}

	var s string

	s = rule.Protocol.CiscoString()
	if s != "" {
		result = append(result, s)
	}

	s = rule.FromNet.CiscoString()
	if s != "" {
		result = append(result, s)
	}

	s = rule.FromPortRanges.CiscoString()
	if s != "" {
		result = append(result, s)
	}

	s = rule.ToNet.CiscoString()
	if s != "" {
		result = append(result, s)
	}

	s = rule.ToPortRanges.CiscoString()
	if s != "" {
		result = append(result, s)
	}

	return strings.Join(result, " ")
}

func (acl ACL) WriteTo(writer io.Writer) error {
	for _, rule := range acl.Rules {
		fmt.Fprintf(writer, "access-list %v extended %v\n", acl.Name, rule.CiscoString())
	}
	return nil
}

func (acls ACLs) CiscoString() string {
	var buf bytes.Buffer
	for _, acl := range acls {
		acl.WriteTo(&buf)
	}
	return buf.String()
}

func (acl *ACL) ParseAppendRule(words []string) error {
	rule := ACLRule{}

	switch words[0] {
	case "extended":
		switch words[1] {
		case "permit":
			rule.Action = ACL_ALLOW
		case "deny":
			rule.Action = ACL_DENY
		default:
			return fmt.Errorf("Invalid ACL action: %v", words[1])
		}
		rule.Protocol = parseProtocol(words[2])

		words = words[3:]

		rule.FromNet, rule.FromPortRanges, words = parseNetworkAndPortRanges(words)
		rule.ToNet, rule.ToPortRanges, words = parseNetworkAndPortRanges(words)

		if len(words) > 0 {
			switch words[0] {
			case "established":
				rule.Flags |= ACLFL_ESTABLISHED
			}
		}
	default:
		warning("Cannot parse acl rule: %v", words)
		return nil
	}

	acl.Rules = append(acl.Rules, rule)

	return nil
}
