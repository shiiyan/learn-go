package policyrulemodeling

import (
	"context"
	"testing"
)

// Mock implementations for testing

type MockSubject struct {
	ID         string
	Attributes map[string]interface{}
}

func (s *MockSubject) GetID() string {
	return s.ID
}

func (s *MockSubject) GetAttributes() map[string]interface{} {
	return s.Attributes
}

type MockResource struct {
	Type       string
	ID         string
	Attributes map[string]interface{}
}

func (r *MockResource) GetType() string {
	return r.Type
}

func (r *MockResource) GetID() string {
	return r.ID
}

func (r *MockResource) GetAttributes() map[string]interface{} {
	return r.Attributes
}

type MockAction struct {
	Name string
}

func (a *MockAction) GetName() string {
	return a.Name
}

type MockRule struct {
	ID       string
	RuleEffect Effect
	RulePriority int
	MatchFunc  func(ctx context.Context, subject Subject, resource Resource, action Action) bool
}

func (r *MockRule) Matches(ctx context.Context, subject Subject, resource Resource, action Action) bool {
	if r.MatchFunc != nil {
		return r.MatchFunc(ctx, subject, resource, action)
	}
	return false
}

func (r *MockRule) GetID() string {
	return r.ID
}

func (r *MockRule) Effect() Effect {
	return r.RuleEffect
}

func (r *MockRule) Priority() int {
	return r.RulePriority
}

// Test SimplePolicy

func TestSimplePolicy_Evaluate(t *testing.T) {
	ctx := context.Background()
	subject := &MockSubject{ID: "user123", Attributes: map[string]interface{}{"role": "admin"}}
	resource := &MockResource{Type: "document", ID: "doc456", Attributes: map[string]interface{}{"owner": "user123"}}
	action := &MockAction{Name: "read"}

	t.Run("allow when rule matches and allows", func(t *testing.T) {
		allowRule := &MockRule{
			ID:         "allow-rule-1",
			RuleEffect: EffectAllow,
			MatchFunc: func(ctx context.Context, subject Subject, resource Resource, action Action) bool {
				return subject.GetID() == "user123" && action.GetName() == "read"
			},
		}

		policy := &SimplePolicy{
			ID:    "simple-policy-1",
			Name:  "Simple Test Policy",
			Rules: []Rule{allowRule},
		}

		decision := policy.Evaluate(ctx, subject, resource, action)

		if !decision.Allow {
			t.Errorf("Expected Allow=true, got %v", decision.Allow)
		}
		if decision.MatchedBy != "simple-policy-1" {
			t.Errorf("Expected MatchedBy='simple-policy-1', got %v", decision.MatchedBy)
		}
		if decision.Reason != "Simple Test Policy" {
			t.Errorf("Expected Reason='Simple Test Policy', got %v", decision.Reason)
		}
	})

	t.Run("deny when rule matches but denies", func(t *testing.T) {
		denyRule := &MockRule{
			ID:         "deny-rule-1",
			RuleEffect: EffectDeny,
			MatchFunc: func(ctx context.Context, subject Subject, resource Resource, action Action) bool {
				return subject.GetID() == "user123" && action.GetName() == "read"
			},
		}

		policy := &SimplePolicy{
			ID:    "simple-policy-2",
			Name:  "Simple Deny Policy",
			Rules: []Rule{denyRule},
		}

		decision := policy.Evaluate(ctx, subject, resource, action)

		if decision.Allow {
			t.Errorf("Expected Allow=false, got %v", decision.Allow)
		}
		if decision.MatchedBy != "simple-policy-2" {
			t.Errorf("Expected MatchedBy='simple-policy-2', got %v", decision.MatchedBy)
		}
	})

	t.Run("deny when no rules match", func(t *testing.T) {
		noMatchRule := &MockRule{
			ID:         "no-match-rule",
			RuleEffect: EffectAllow,
			MatchFunc: func(ctx context.Context, subject Subject, resource Resource, action Action) bool {
				return false // never matches
			},
		}

		policy := &SimplePolicy{
			ID:    "simple-policy-3",
			Name:  "No Match Policy",
			Rules: []Rule{noMatchRule},
		}

		decision := policy.Evaluate(ctx, subject, resource, action)

		if decision.Allow {
			t.Errorf("Expected Allow=false, got %v", decision.Allow)
		}
		if decision.Reason != "no matching rules" {
			t.Errorf("Expected Reason='no matching rules', got %v", decision.Reason)
		}
	})
}

// Test AllMustAllowPolicy

func TestAllMustAllowPolicy_Evaluate(t *testing.T) {
	ctx := context.Background()
	subject := &MockSubject{ID: "user123", Attributes: map[string]interface{}{"role": "admin"}}
	resource := &MockResource{Type: "document", ID: "doc456", Attributes: map[string]interface{}{"owner": "user123"}}
	action := &MockAction{Name: "read"}

	t.Run("allow when all matching rules allow", func(t *testing.T) {
		allowRule1 := &MockRule{
			ID:         "allow-rule-1",
			RuleEffect: EffectAllow,
			MatchFunc: func(ctx context.Context, subject Subject, resource Resource, action Action) bool {
				return subject.GetID() == "user123"
			},
		}
		allowRule2 := &MockRule{
			ID:         "allow-rule-2",
			RuleEffect: EffectAllow,
			MatchFunc: func(ctx context.Context, subject Subject, resource Resource, action Action) bool {
				return action.GetName() == "read"
			},
		}

		policy := &AllMustAllowPolicy{
			ID:    "all-must-allow-1",
			Name:  "All Must Allow Policy",
			Rules: []Rule{allowRule1, allowRule2},
		}

		decision := policy.Evaluate(ctx, subject, resource, action)

		if !decision.Allow {
			t.Errorf("Expected Allow=true, got %v", decision.Allow)
		}
		if decision.MatchedBy != "all-must-allow-1" {
			t.Errorf("Expected MatchedBy='all-must-allow-1', got %v", decision.MatchedBy)
		}
		if decision.Reason != "all rules allowed: allow-rule-1, allow-rule-2" {
			t.Errorf("Expected specific reason, got %v", decision.Reason)
		}
	})

	t.Run("deny when some matching rules deny", func(t *testing.T) {
		allowRule := &MockRule{
			ID:         "allow-rule-1",
			RuleEffect: EffectAllow,
			MatchFunc: func(ctx context.Context, subject Subject, resource Resource, action Action) bool {
				return subject.GetID() == "user123"
			},
		}
		denyRule := &MockRule{
			ID:         "deny-rule-1",
			RuleEffect: EffectDeny,
			MatchFunc: func(ctx context.Context, subject Subject, resource Resource, action Action) bool {
				return action.GetName() == "read"
			},
		}

		policy := &AllMustAllowPolicy{
			ID:    "all-must-allow-2",
			Name:  "Mixed Rules Policy",
			Rules: []Rule{allowRule, denyRule},
		}

		decision := policy.Evaluate(ctx, subject, resource, action)

		if decision.Allow {
			t.Errorf("Expected Allow=false, got %v", decision.Allow)
		}
		if decision.Reason != "denied by rules: deny-rule-1" {
			t.Errorf("Expected specific deny reason, got %v", decision.Reason)
		}
	})

	t.Run("deny when no rules match", func(t *testing.T) {
		noMatchRule := &MockRule{
			ID:         "no-match-rule",
			RuleEffect: EffectAllow,
			MatchFunc: func(ctx context.Context, subject Subject, resource Resource, action Action) bool {
				return false // never matches
			},
		}

		policy := &AllMustAllowPolicy{
			ID:    "all-must-allow-3",
			Name:  "No Match Policy",
			Rules: []Rule{noMatchRule},
		}

		decision := policy.Evaluate(ctx, subject, resource, action)

		if decision.Allow {
			t.Errorf("Expected Allow=false, got %v", decision.Allow)
		}
		if decision.Reason != "no matching rules" {
			t.Errorf("Expected Reason='no matching rules', got %v", decision.Reason)
		}
	})

	t.Run("deny when multiple rules deny", func(t *testing.T) {
		denyRule1 := &MockRule{
			ID:         "deny-rule-1",
			RuleEffect: EffectDeny,
			MatchFunc: func(ctx context.Context, subject Subject, resource Resource, action Action) bool {
				return subject.GetID() == "user123"
			},
		}
		denyRule2 := &MockRule{
			ID:         "deny-rule-2",
			RuleEffect: EffectDeny,
			MatchFunc: func(ctx context.Context, subject Subject, resource Resource, action Action) bool {
				return action.GetName() == "read"
			},
		}

		policy := &AllMustAllowPolicy{
			ID:    "all-must-allow-4",
			Name:  "Multiple Deny Policy",
			Rules: []Rule{denyRule1, denyRule2},
		}

		decision := policy.Evaluate(ctx, subject, resource, action)

		if decision.Allow {
			t.Errorf("Expected Allow=false, got %v", decision.Allow)
		}
		if decision.Reason != "denied by rules: deny-rule-1, deny-rule-2" {
			t.Errorf("Expected multiple deny reason, got %v", decision.Reason)
		}
	})
}

// Test helper function
func TestGetRuleIDs(t *testing.T) {
	rule1 := &MockRule{ID: "rule-1", RuleEffect: EffectAllow}
	rule2 := &MockRule{ID: "rule-2", RuleEffect: EffectDeny}
	rule3 := &MockRule{ID: "rule-3", RuleEffect: EffectAllow}

	rules := []Rule{rule1, rule2, rule3}
	ids := getRuleIDs(rules)

	expected := []string{"rule-1", "rule-2", "rule-3"}
	if len(ids) != len(expected) {
		t.Errorf("Expected %d IDs, got %d", len(expected), len(ids))
	}

	for i, expectedID := range expected {
		if ids[i] != expectedID {
			t.Errorf("Expected ID at index %d to be %s, got %s", i, expectedID, ids[i])
		}
	}
}

// Integration test with realistic scenario
func TestRealisticPolicyScenario(t *testing.T) {
	ctx := context.Background()

	// Admin user
	adminUser := &MockSubject{
		ID:         "admin123",
		Attributes: map[string]interface{}{"role": "admin", "department": "IT"},
	}

	// Regular user
	regularUser := &MockSubject{
		ID:         "user456",
		Attributes: map[string]interface{}{"role": "user", "department": "Sales"},
	}

	// Sensitive document
	sensitiveDoc := &MockResource{
		Type:       "document",
		ID:         "sensitive-doc-789",
		Attributes: map[string]interface{}{"classification": "confidential", "owner": "admin123"},
	}

	// Read action
	readAction := &MockAction{Name: "read"}

	// Rules
	adminRule := &MockRule{
		ID:         "admin-access-rule",
		RuleEffect: EffectAllow,
		MatchFunc: func(ctx context.Context, subject Subject, resource Resource, action Action) bool {
			attrs := subject.GetAttributes()
			role, ok := attrs["role"].(string)
			return ok && role == "admin"
		},
	}

	ownerRule := &MockRule{
		ID:         "owner-access-rule",
		RuleEffect: EffectAllow,
		MatchFunc: func(ctx context.Context, subject Subject, resource Resource, action Action) bool {
			resourceAttrs := resource.GetAttributes()
			owner, ok := resourceAttrs["owner"].(string)
			return ok && owner == subject.GetID()
		},
	}

	confidentialRule := &MockRule{
		ID:         "confidential-restriction-rule",
		RuleEffect: EffectDeny,
		MatchFunc: func(ctx context.Context, subject Subject, resource Resource, action Action) bool {
			resourceAttrs := resource.GetAttributes()
			classification, ok := resourceAttrs["classification"].(string)
			if !ok || classification != "confidential" {
				return false
			}
			
			subjectAttrs := subject.GetAttributes()
			role, ok := subjectAttrs["role"].(string)
			return !ok || role != "admin"
		},
	}

	// Test with AllMustAllowPolicy
	policy := &AllMustAllowPolicy{
		ID:    "secure-document-policy",
		Name:  "Secure Document Access Policy",
		Rules: []Rule{adminRule, ownerRule, confidentialRule},
	}

	t.Run("admin can access confidential document they own", func(t *testing.T) {
		decision := policy.Evaluate(ctx, adminUser, sensitiveDoc, readAction)
		
		if !decision.Allow {
			t.Errorf("Expected admin to access their confidential document, got: %v", decision.Reason)
		}
	})

	t.Run("regular user cannot access confidential document", func(t *testing.T) {
		decision := policy.Evaluate(ctx, regularUser, sensitiveDoc, readAction)
		
		if decision.Allow {
			t.Errorf("Expected regular user to be denied access to confidential document, got: %v", decision.Reason)
		}
	})
}
