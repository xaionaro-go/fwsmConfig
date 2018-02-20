package fwsmConfig

import (
	"bytes"
	"fmt"
	"github.com/xaionaro-go/networkControl"
	"io"
	"strings"
)

type ACLRule networkControl.ACLRule
type ACLRules networkControl.ACLRules
type ACL networkControl.ACL
type ACLs networkControl.ACLs

func (acl *ACL) ParseAppendRule(words []string) error {
	rule := ACLRule{}

	switch words[0] {
	case "extended":
		switch words[1] {
		case "permit":
			rule.Action = networkControl.ACL_ALLOW
		case "deny":
			rule.Action = networkControl.ACL_DENY
		default:
			return fmt.Errorf("Invalid ACL action: %v", words[1])
		}
		rule.Protocol = networkControl.Protocol(parseProtocol(words[2]))

		words = words[3:]

		rule.FromNet, rule.FromPortRanges, words = parseNetworkAndPortRanges(words)
		rule.ToNet, rule.ToPortRanges, words = parseNetworkAndPortRanges(words)

		if len(words) > 0 {
			switch words[0] {
			case "established":
				rule.Flags |= networkControl.ACLFL_ESTABLISHED
			}
		}
	default:
		warning("Cannot parse acl rule: %v", words)
		return nil
	}

	acl.Rules = append(acl.Rules, networkControl.ACLRule(rule))

	return nil
}

func (rule ACLRule) CiscoString() string {
	result := []string{}

	switch rule.Action {
	case networkControl.ACL_ALLOW:
		result = append(result, "permit")
	case networkControl.ACL_DENY:
		result = append(result, "deny")
	default:
		panic(fmt.Errorf("Invalid rule action: %v", rule.Action))
	}

	var s string

	s = Protocol(rule.Protocol).CiscoString()
	if s != "" {
		result = append(result, s)
	}

	s = IPNet(rule.FromNet).CiscoString()
	if s != "" {
		result = append(result, s)
	}

	s = PortRanges(rule.FromPortRanges).CiscoString()
	if s != "" {
		result = append(result, s)
	}

	s = IPNet(rule.ToNet).CiscoString()
	if s != "" {
		result = append(result, s)
	}

	s = PortRanges(rule.ToPortRanges).CiscoString()
	if s != "" {
		result = append(result, s)
	}

	return strings.Join(result, " ")
}

func (acl ACL) WriteTo(writer io.Writer) error {
	for _, rule := range acl.Rules {
		fmt.Fprintf(writer, "access-list %v extended %v\n", acl.Name, ACLRule(rule).CiscoString())
	}
	for _, ifName := range acl.VLANNames {
		fmt.Fprintf(writer, "access-group %v in interface %v\n", acl.Name, ifName)
	}
	return nil
}

func (acls ACLs) CiscoString() string {
	var buf bytes.Buffer
	for _, acl := range acls {
		ACL(*acl).WriteTo(&buf)
	}
	return buf.String()
}
