package secretkeys

import (
	"strings"
	"testing"
)

const (
	aliceKey = "age1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"
	bobKey   = "age1ppppppppppppppppppppppppppppppppppppppppppppppppppppppppppp"
)

func TestParseRecipientsSkipsCommentsAndSorts(t *testing.T) {
	recipients, err := ParseRecipients([]byte(`
# comment
bob ` + bobKey + `
alice ` + aliceKey + `
`))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(recipients), 2; got != want {
		t.Fatalf("len = %d, want %d", got, want)
	}
	if got, want := recipients[0].User, "alice"; got != want {
		t.Fatalf("first user = %q, want %q", got, want)
	}
}

func TestAddRecipientAddsAndUpdates(t *testing.T) {
	recipients, updated, err := AddRecipient(nil, "alice", aliceKey)
	if err != nil {
		t.Fatal(err)
	}
	if updated {
		t.Fatal("updated = true, want false for new user")
	}

	recipients, updated, err = AddRecipient(recipients, "alice", bobKey)
	if err != nil {
		t.Fatal(err)
	}
	if !updated {
		t.Fatal("updated = false, want true for existing user")
	}
	if got, want := recipients[0].Key, bobKey; got != want {
		t.Fatalf("key = %q, want %q", got, want)
	}
}

func TestRemoveRecipient(t *testing.T) {
	recipients := []Recipient{{User: "alice", Key: aliceKey}, {User: "bob", Key: bobKey}}
	recipients, err := RemoveRecipient(recipients, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(recipients), 1; got != want {
		t.Fatalf("len = %d, want %d", got, want)
	}
	if got, want := recipients[0].User, "bob"; got != want {
		t.Fatalf("remaining user = %q, want %q", got, want)
	}
	if _, err := RemoveRecipient(recipients, "alice"); err == nil {
		t.Fatal("removing a missing user succeeded")
	}
}

func TestFormatSOPSConfigIncludesRecipients(t *testing.T) {
	config := string(FormatSOPSConfig([]Recipient{
		{User: "bob", Key: bobKey},
		{User: "alice", Key: aliceKey},
	}))
	if !strings.Contains(config, "path_regex: '^\\.secrets/.*\\.sops\\.(env|json|yaml|yml)$'") {
		t.Fatalf("config missing path regex:\n%s", config)
	}
	if !strings.Contains(config, aliceKey+","+bobKey) {
		t.Fatalf("config did not include sorted recipients:\n%s", config)
	}
}

func TestValidationRejectsUnsafeInput(t *testing.T) {
	if err := ValidateUser("../alice"); err == nil {
		t.Fatal("unsafe user accepted")
	}
	if err := ValidateAgeRecipient("ssh-ed25519 AAAA"); err == nil {
		t.Fatal("non-age recipient accepted")
	}
}
