package policyrulemodeling

import (
	"context"
	"strings"
)

type Policy interface {
	Evaluate(ctx context.Context, subject Subject, resource Resource, action Action) Decision
	GetID() string
	GetName() string
}

type Rule interface {
	Matches(ctx context.Context, subject Subject, resource Resource, action Action) bool
	GetID() string
	Effect() Effect
	Priority() int
}

type Decision struct {
	Allow     bool
	Reason    string
	MatchedBy string // which rule/policy made the decision
}

type Effect int

const (
	EffectAllow Effect = iota
	EffectDeny
)

type Subject interface {
	GetID() string
	GetAttributes() map[string]interface{}
}

type Resource interface {
	GetType() string
	GetID() string
	GetAttributes() map[string]interface{}
}

type Action interface {
	GetName() string
}

type SimplePolicy struct {
	ID    string
	Name  string
	Rules []Rule
}

func (p *SimplePolicy) Evaluate(ctx context.Context, subject Subject, resource Resource, action Action) Decision {
	for _, rule := range p.Rules {
		if rule.Matches(ctx, subject, resource, action) {
			return Decision{
				Allow:     rule.Effect() == EffectAllow,
				Reason:    p.Name,
				MatchedBy: p.ID,
			}
		}
	}
	return Decision{Allow: false, Reason: "no matching rules"}
}

func (p *SimplePolicy) GetID() string {
	return p.ID
}

func (p *SimplePolicy) GetName() string {
	return p.Name
}

type AllMustAllowPolicy struct {
	ID    string
	Name  string
	Rules []Rule
}

func (p *AllMustAllowPolicy) Evaluate(ctx context.Context, subject Subject, resource Resource, action Action) Decision {
	matchedRules := []Rule{}

	for _, rule := range p.Rules {
		if rule.Matches(ctx, subject, resource, action) {
			matchedRules = append(matchedRules, rule)
		}
	}

	if len(matchedRules) == 0 {
		return Decision{Allow: false, Reason: "no matching rules"}
	}

	deniedBy := []string{}
	for _, rule := range matchedRules {
		if rule.Effect() != EffectAllow {
			deniedBy = append(deniedBy, rule.GetID())
		}
	}

	if len(deniedBy) > 0 {
		return Decision{
			Allow:     false,
			Reason:    "denied by rules: " + strings.Join(deniedBy, ", "),
			MatchedBy: p.ID,
		}
	}

	return Decision{
		Allow:     true,
		Reason:    "all rules allowed: " + strings.Join(getRuleIDs(matchedRules), ", "),
		MatchedBy: p.ID,
	}
}

func getRuleIDs(rules []Rule) []string {
	ids := make([]string, len(rules))
	for i, rule := range rules {
		ids[i] = rule.GetID()
	}
	return ids
}
