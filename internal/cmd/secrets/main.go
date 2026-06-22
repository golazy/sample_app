package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"sample_app/internal/secretkeys"
)

const (
	recipientsPath = ".secrets/recipients.txt"
	sopsConfigPath = ".sops.yaml"
	keyDir         = ".secrets/keys"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "secrets:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		usage(stderr)
		return errors.New("missing command")
	}

	switch args[0] {
	case "new-key":
		if len(args) != 2 {
			usage(stderr)
			return errors.New("usage: new-key <user>")
		}
		return createKey(args[1], stdout, stderr)
	case "add-key":
		if len(args) != 3 {
			usage(stderr)
			return errors.New("usage: add-key <user> <age-recipient>")
		}
		return addKey(args[1], args[2], stdout)
	case "remove-user":
		if len(args) != 2 {
			usage(stderr)
			return errors.New("usage: remove-user <user>")
		}
		return removeUser(args[1], stdout)
	case "users":
		if len(args) != 1 {
			usage(stderr)
			return errors.New("usage: users")
		}
		return listUsers(stdout)
	default:
		usage(stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage:")
	fmt.Fprintln(w, "  secrets new-key <user>")
	fmt.Fprintln(w, "  secrets add-key <user> <age-recipient>")
	fmt.Fprintln(w, "  secrets remove-user <user>")
	fmt.Fprintln(w, "  secrets users")
}

func createKey(user string, stdout, stderr io.Writer) error {
	if err := secretkeys.ValidateUser(user); err != nil {
		return err
	}
	recipients, err := loadRecipients()
	if err != nil {
		return err
	}
	for _, recipient := range recipients {
		if recipient.User == user {
			return fmt.Errorf("user %q already has a registered recipient; use secrets:add-key to rotate it", user)
		}
	}

	if err := os.MkdirAll(keyDir, 0o700); err != nil {
		return err
	}
	keyPath := filepath.Join(keyDir, user+".txt")
	if _, err := os.Stat(keyPath); err == nil {
		return fmt.Errorf("%s already exists", keyPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	cmd := exec.Command("age-keygen", "-o", keyPath)
	cmd.Stdout = stderr
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("age-keygen: %w", err)
	}
	if err := os.Chmod(keyPath, 0o600); err != nil {
		return err
	}

	key, err := publicRecipient(keyPath)
	if err != nil {
		return err
	}
	recipients, _, err = secretkeys.AddRecipient(recipients, user, key)
	if err != nil {
		return err
	}
	if err := saveRecipients(recipients); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "created private age identity: %s\n", keyPath)
	fmt.Fprintf(stdout, "registered recipient for %s: %s\n", user, key)
	printUpdateKeysHint(stdout)
	return nil
}

func publicRecipient(keyPath string) (string, error) {
	cmd := exec.Command("age-keygen", "-y", keyPath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("age-keygen -y: %w", err)
	}
	key := strings.TrimSpace(string(output))
	if err := secretkeys.ValidateAgeRecipient(key); err != nil {
		return "", fmt.Errorf("age-keygen -y returned invalid recipient: %w", err)
	}
	return key, nil
}

func addKey(user, key string, stdout io.Writer) error {
	recipients, err := loadRecipients()
	if err != nil {
		return err
	}
	recipients, updated, err := secretkeys.AddRecipient(recipients, user, key)
	if err != nil {
		return err
	}
	if err := saveRecipients(recipients); err != nil {
		return err
	}
	if updated {
		fmt.Fprintf(stdout, "updated recipient for %s\n", user)
	} else {
		fmt.Fprintf(stdout, "added recipient for %s\n", user)
	}
	printUpdateKeysHint(stdout)
	return nil
}

func removeUser(user string, stdout io.Writer) error {
	recipients, err := loadRecipients()
	if err != nil {
		return err
	}
	recipients, err = secretkeys.RemoveRecipient(recipients, user)
	if err != nil {
		return err
	}
	if err := saveRecipients(recipients); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "removed recipient for %s\n", user)
	printUpdateKeysHint(stdout)
	return nil
}

func listUsers(stdout io.Writer) error {
	recipients, err := loadRecipients()
	if err != nil {
		return err
	}
	if len(recipients) == 0 {
		fmt.Fprintln(stdout, "No users have access yet.")
		return nil
	}
	for _, recipient := range recipients {
		fmt.Fprintf(stdout, "%-20s %s\n", recipient.User, recipient.Key)
	}
	return nil
}

func loadRecipients() ([]secretkeys.Recipient, error) {
	data, err := os.ReadFile(recipientsPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return secretkeys.ParseRecipients(data)
}

func saveRecipients(recipients []secretkeys.Recipient) error {
	if err := os.MkdirAll(filepath.Dir(recipientsPath), 0o755); err != nil {
		return err
	}
	if err := writeFile(recipientsPath, secretkeys.FormatRecipients(recipients), 0o644); err != nil {
		return err
	}
	return writeFile(sopsConfigPath, secretkeys.FormatSOPSConfig(recipients), 0o644)
}

func writeFile(path string, data []byte, perm os.FileMode) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, perm); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func printUpdateKeysHint(stdout io.Writer) {
	matches, err := filepath.Glob(".secrets/*.sops.*")
	if err != nil || len(matches) == 0 {
		return
	}
	sort.Strings(matches)
	fmt.Fprintln(stdout, "Apply recipient changes to existing encrypted files:")
	for _, match := range matches {
		fmt.Fprintf(stdout, "  sops updatekeys -y %s\n", match)
	}
}
