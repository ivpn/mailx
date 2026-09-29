package repository

import (
	"strings"
	"testing"
)

func TestSanitizeAliasSort(t *testing.T) {
	tests := []struct {
		name          string
		columnPrefix  string
		sortBy        string
		sortOrder     string
		expectedBy    string
		expectedOrder string
	}{
		{name: "valid name/ASC", columnPrefix: "a.", sortBy: "name", sortOrder: "ASC", expectedBy: "a.name", expectedOrder: "ASC"},
		{name: "valid created_at/DESC", columnPrefix: "a.", sortBy: "created_at", sortOrder: "DESC", expectedBy: "a.created_at", expectedOrder: "DESC"},
		{name: "valid updated_at/ASC", columnPrefix: "a.", sortBy: "updated_at", sortOrder: "ASC", expectedBy: "a.updated_at", expectedOrder: "ASC"},
		{name: "empty sortBy defaults", columnPrefix: "a.", sortBy: "", sortOrder: "ASC", expectedBy: "a.created_at", expectedOrder: "ASC"},
		{name: "empty sortOrder defaults", columnPrefix: "a.", sortBy: "name", sortOrder: "", expectedBy: "a.name", expectedOrder: "DESC"},
		{name: "both empty default", columnPrefix: "a.", sortBy: "", sortOrder: "", expectedBy: "a.created_at", expectedOrder: "DESC"},
		{name: "lowercase sortOrder rejected", columnPrefix: "a.", sortBy: "name", sortOrder: "asc", expectedBy: "a.name", expectedOrder: "DESC"},
		{name: "unknown column falls back", columnPrefix: "a.", sortBy: "password", sortOrder: "ASC", expectedBy: "a.created_at", expectedOrder: "ASC"},
		{name: "SQL injection via sortBy stacked query", columnPrefix: "a.", sortBy: "created_at; DROP TABLE aliases;--", sortOrder: "DESC", expectedBy: "a.created_at", expectedOrder: "DESC"},
		{name: "SQL injection via sortBy subquery", columnPrefix: "a.", sortBy: "(SELECT password FROM users)", sortOrder: "DESC", expectedBy: "a.created_at", expectedOrder: "DESC"},
		{name: "SQL injection via sortBy comma expression", columnPrefix: "a.", sortBy: "name, (SELECT password FROM users)", sortOrder: "DESC", expectedBy: "a.created_at", expectedOrder: "DESC"},
		{name: "SQL injection via sortOrder stacked query", columnPrefix: "a.", sortBy: "name", sortOrder: "ASC; DROP TABLE aliases;--", expectedBy: "a.name", expectedOrder: "DESC"},
		{name: "case-mismatched column rejected", columnPrefix: "a.", sortBy: "Name", sortOrder: "DESC", expectedBy: "a.created_at", expectedOrder: "DESC"},
		{name: "no prefix valid column", columnPrefix: "", sortBy: "name", sortOrder: "ASC", expectedBy: "name", expectedOrder: "ASC"},
		{name: "no prefix falls back", columnPrefix: "", sortBy: "password", sortOrder: "DESC", expectedBy: "created_at", expectedOrder: "DESC"},
		{name: "no prefix SQL injection via sortBy", columnPrefix: "", sortBy: "created_at; DROP TABLE aliases;--", sortOrder: "DESC", expectedBy: "created_at", expectedOrder: "DESC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBy, gotOrder := sanitizeAliasSort(tt.columnPrefix, tt.sortBy, tt.sortOrder)
			if gotBy != tt.expectedBy {
				t.Errorf("sanitizeAliasSort(%q, %q, %q) sortBy = %q, want %q", tt.columnPrefix, tt.sortBy, tt.sortOrder, gotBy, tt.expectedBy)
			}
			if gotOrder != tt.expectedOrder {
				t.Errorf("sanitizeAliasSort(%q, %q, %q) sortOrder = %q, want %q", tt.columnPrefix, tt.sortBy, tt.sortOrder, gotOrder, tt.expectedOrder)
			}
		})
	}
}

func TestAliasSearchFilter(t *testing.T) {
	t.Run("empty inputs produce no filter", func(t *testing.T) {
		filter, args := aliasSearchFilter("a.", "", "")
		if filter != "" {
			t.Errorf("expected empty filter, got %q", filter)
		}
		if len(args) != 0 {
			t.Errorf("expected no args, got %v", args)
		}
	})

	t.Run("wildcard true binds a parameter", func(t *testing.T) {
		filter, args := aliasSearchFilter("a.", "true", "")
		if !strings.Contains(filter, "a.catch_all = ?") {
			t.Errorf("expected bound catch_all clause, got %q", filter)
		}
		if len(args) != 1 || args[0] != true {
			t.Errorf("expected args [true], got %v", args)
		}
	})

	t.Run("wildcard false binds a parameter", func(t *testing.T) {
		filter, args := aliasSearchFilter("", "false", "")
		if !strings.Contains(filter, "catch_all = ?") {
			t.Errorf("expected bound catch_all clause, got %q", filter)
		}
		if len(args) != 1 || args[0] != false {
			t.Errorf("expected args [false], got %v", args)
		}
	})

	invalidWildcard := []string{"TRUE", "1", "true OR 1=1", "false; DROP TABLE aliases;--", " true"}
	for _, v := range invalidWildcard {
		t.Run("invalid wildcard is ignored: "+v, func(t *testing.T) {
			filter, args := aliasSearchFilter("a.", v, "")
			if filter != "" || len(args) != 0 {
				t.Errorf("aliasSearchFilter(%q) = (%q, %v), want empty filter/args", v, filter, args)
			}
		})
	}

	t.Run("search binds two LIKE parameters", func(t *testing.T) {
		filter, args := aliasSearchFilter("a.", "", "invoice")
		if strings.Count(filter, "LIKE ?") != 2 {
			t.Errorf("expected two LIKE ? placeholders, got %q", filter)
		}
		if len(args) != 2 || args[0] != "%invoice%" || args[1] != "%invoice%" {
			t.Errorf("expected args [%%invoice%% %%invoice%%], got %v", args)
		}
	})

	sqliPayloads := []string{
		"x' OR '1'='1",
		"'; DROP TABLE aliases;--",
		"' UNION SELECT password FROM users--",
		"%' OR 1=1-- -",
		"a\") OR (\"1\"=\"1",
	}
	for _, payload := range sqliPayloads {
		t.Run("SQLi payload never appears in filter text: "+payload, func(t *testing.T) {
			filter, args := aliasSearchFilter("a.", "", payload)
			if strings.Contains(filter, payload) {
				t.Errorf("filter %q must not contain raw payload %q", filter, payload)
			}
			if !strings.Contains(filter, "LIKE ?") {
				t.Errorf("expected LIKE ? placeholders, got %q", filter)
			}
			found := false
			for _, a := range args {
				if s, ok := a.(string); ok && strings.Contains(s, payload) {
					found = true
				}
			}
			if !found {
				t.Errorf("expected payload to appear safely bound in args %v", args)
			}
		})
	}

	t.Run("columnPrefix a. qualifies columns", func(t *testing.T) {
		filter, _ := aliasSearchFilter("a.", "true", "x")
		for _, col := range []string{"a.catch_all", "a.name", "a.description"} {
			if !strings.Contains(filter, col) {
				t.Errorf("expected filter to reference %q, got %q", col, filter)
			}
		}
	})

	t.Run("columnPrefix empty leaves columns unqualified", func(t *testing.T) {
		filter, _ := aliasSearchFilter("", "true", "x")
		for _, col := range []string{"catch_all", "name", "description"} {
			if !strings.Contains(filter, col) {
				t.Errorf("expected filter to reference %q, got %q", col, filter)
			}
		}
		if strings.Contains(filter, "a.catch_all") {
			t.Errorf("expected unqualified columns, got %q", filter)
		}
	})
}

func TestAliasStatusFilter(t *testing.T) {
	tests := []struct {
		name           string
		status         string
		expectedFilter string
		expectedArgs   []any
		expectedUnscop bool
	}{
		{name: "active_inactive", status: "active_inactive", expectedFilter: "AND a.deleted_at IS NULL", expectedArgs: nil, expectedUnscop: false},
		{name: "active", status: "active", expectedFilter: "AND a.deleted_at IS NULL AND a.enabled = ?", expectedArgs: []any{true}, expectedUnscop: false},
		{name: "inactive", status: "inactive", expectedFilter: "AND a.deleted_at IS NULL AND a.enabled = ?", expectedArgs: []any{false}, expectedUnscop: false},
		{name: "deleted", status: "deleted", expectedFilter: "AND a.deleted_at IS NOT NULL", expectedArgs: nil, expectedUnscop: true},
		{name: "all", status: "all", expectedFilter: "", expectedArgs: nil, expectedUnscop: true},
		{name: "empty falls back to active_inactive shape", status: "", expectedFilter: "AND a.deleted_at IS NULL", expectedArgs: nil, expectedUnscop: false},
		{name: "unrecognized falls back to active_inactive shape", status: "bogus; DROP TABLE aliases;--", expectedFilter: "AND a.deleted_at IS NULL", expectedArgs: nil, expectedUnscop: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, args, unscoped := aliasStatusFilter("a.", tt.status)
			if filter != tt.expectedFilter {
				t.Errorf("aliasStatusFilter(%q) filter = %q, want %q", tt.status, filter, tt.expectedFilter)
			}
			if len(args) != len(tt.expectedArgs) {
				t.Errorf("aliasStatusFilter(%q) args = %v, want %v", tt.status, args, tt.expectedArgs)
			}
			for i := range tt.expectedArgs {
				if args[i] != tt.expectedArgs[i] {
					t.Errorf("aliasStatusFilter(%q) args[%d] = %v, want %v", tt.status, i, args[i], tt.expectedArgs[i])
				}
			}
			if unscoped != tt.expectedUnscop {
				t.Errorf("aliasStatusFilter(%q) unscoped = %v, want %v", tt.status, unscoped, tt.expectedUnscop)
			}
		})
	}

	t.Run("columnPrefix empty leaves columns unqualified", func(t *testing.T) {
		filter, _, _ := aliasStatusFilter("", "active")
		if !strings.Contains(filter, "enabled = ?") || strings.Contains(filter, "a.enabled") {
			t.Errorf("expected unqualified enabled column, got %q", filter)
		}
	})
}
