package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const (
	testAliceKey = "age1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"
	testBobKey   = "age1ppppppppppppppppppppppppppppppppppppppppppppppppppppppppppp"
)

func TestRunAddListAndRemoveUser(t *testing.T) {
	dir := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	}()

	var stdout, stderr bytes.Buffer
	if err := run([]string{"add-key", "alice", testAliceKey}, &stdout, &stderr); err != nil {
		t.Fatalf("add-key: %v\nstderr:\n%s", err, stderr.String())
	}
	if err := run([]string{"add-key", "bob", testBobKey}, &stdout, &stderr); err != nil {
		t.Fatalf("add-key bob: %v\nstderr:\n%s", err, stderr.String())
	}

	recipients, err := os.ReadFile(recipientsPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(recipients), "alice "+testAliceKey) {
		t.Fatalf("recipients missing alice:\n%s", recipients)
	}
	config, err := os.ReadFile(sopsConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(config), testAliceKey+","+testBobKey) {
		t.Fatalf("sops config missing sorted keys:\n%s", config)
	}

	stdout.Reset()
	if err := run([]string{"users"}, &stdout, &stderr); err != nil {
		t.Fatalf("users: %v", err)
	}
	if !strings.Contains(stdout.String(), "alice") || !strings.Contains(stdout.String(), "bob") {
		t.Fatalf("users output missing registered users:\n%s", stdout.String())
	}

	if err := run([]string{"remove-user", "alice"}, &stdout, &stderr); err != nil {
		t.Fatalf("remove-user: %v", err)
	}
	recipients, err = os.ReadFile(recipientsPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(recipients), "\nalice ") {
		t.Fatalf("recipients still contains alice:\n%s", recipients)
	}
}
