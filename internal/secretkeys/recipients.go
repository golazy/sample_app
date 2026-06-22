package secretkeys

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Recipient struct {
	User string
	Key  string
}

var (
	userPattern         = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)
	ageRecipientPattern = regexp.MustCompile(`^age1[0-9a-z]+$`)
)

func ValidateUser(user string) error {
	if !userPattern.MatchString(user) {
		return fmt.Errorf("user must start with a letter or number and contain only letters, numbers, dots, underscores, or dashes")
	}
	return nil
}

func ValidateAgeRecipient(key string) error {
	if !ageRecipientPattern.MatchString(key) {
		return fmt.Errorf("age recipient must look like age1...")
	}
	return nil
}

func ParseRecipients(data []byte) ([]Recipient, error) {
	var recipients []Recipient
	seen := map[string]struct{}{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	line := 0
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		fields := strings.Fields(text)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected <user> <age-recipient>", line)
		}
		user, key := fields[0], fields[1]
		if err := ValidateUser(user); err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		if err := ValidateAgeRecipient(key); err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		if _, ok := seen[user]; ok {
			return nil, fmt.Errorf("line %d: duplicate user %q", line, user)
		}
		seen[user] = struct{}{}
		recipients = append(recipients, Recipient{User: user, Key: key})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	sortRecipients(recipients)
	return recipients, nil
}

func FormatRecipients(recipients []Recipient) []byte {
	ordered := append([]Recipient(nil), recipients...)
	sortRecipients(ordered)

	var b strings.Builder
	b.WriteString("# Public age recipients allowed to decrypt development SOPS secrets.\n")
	b.WriteString("# Format: <user> <age-recipient>\n")
	b.WriteString("#\n")
	b.WriteString("# Create and register a local identity:\n")
	b.WriteString("#   mise run secrets:new-key -- alice\n")
	b.WriteString("#\n")
	b.WriteString("# Register a public recipient shared by another user:\n")
	b.WriteString("#   mise run secrets:add-key -- alice age1...\n")
	for _, recipient := range ordered {
		fmt.Fprintf(&b, "%s %s\n", recipient.User, recipient.Key)
	}
	return []byte(b.String())
}

func AddRecipient(recipients []Recipient, user, key string) ([]Recipient, bool, error) {
	if err := ValidateUser(user); err != nil {
		return nil, false, err
	}
	if err := ValidateAgeRecipient(key); err != nil {
		return nil, false, err
	}

	updated := append([]Recipient(nil), recipients...)
	for i := range updated {
		if updated[i].User == user {
			updated[i].Key = key
			sortRecipients(updated)
			return updated, true, nil
		}
	}
	updated = append(updated, Recipient{User: user, Key: key})
	sortRecipients(updated)
	return updated, false, nil
}

func RemoveRecipient(recipients []Recipient, user string) ([]Recipient, error) {
	if err := ValidateUser(user); err != nil {
		return nil, err
	}
	updated := make([]Recipient, 0, len(recipients))
	removed := false
	for _, recipient := range recipients {
		if recipient.User == user {
			removed = true
			continue
		}
		updated = append(updated, recipient)
	}
	if !removed {
		return nil, fmt.Errorf("user %q is not in .secrets/recipients.txt", user)
	}
	sortRecipients(updated)
	return updated, nil
}

func FormatSOPSConfig(recipients []Recipient) []byte {
	ordered := append([]Recipient(nil), recipients...)
	sortRecipients(ordered)

	var b strings.Builder
	b.WriteString("# Managed by mise secrets:* tasks. Public user mappings live in .secrets/recipients.txt.\n")
	if len(ordered) == 0 {
		b.WriteString("creation_rules: []\n")
		return []byte(b.String())
	}

	b.WriteString("creation_rules:\n")
	b.WriteString("  - path_regex: '^\\.secrets/.*\\.sops\\.(env|json|yaml|yml)$'\n")
	b.WriteString("    age: ")
	for i, recipient := range ordered {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(recipient.Key)
	}
	b.WriteString("\n")
	return []byte(b.String())
}

func sortRecipients(recipients []Recipient) {
	sort.Slice(recipients, func(i, j int) bool {
		return recipients[i].User < recipients[j].User
	})
}
