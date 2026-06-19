package appinit

import "testing"

func TestSecureCookieKeyUsesEnvironment(t *testing.T) {
	t.Setenv("SECURE_COOKIE_KEY", "env-cookie-key")

	if got := secureCookieKey(); got != "env-cookie-key" {
		t.Fatalf("secureCookieKey() = %q, want %q", got, "env-cookie-key")
	}
}

func TestSecureCookieKeyFallsBackToDevelopmentValue(t *testing.T) {
	t.Setenv("SECURE_COOKIE_KEY", "")

	if got := secureCookieKey(); got != developmentSecureCookieKey {
		t.Fatalf("secureCookieKey() = %q, want %q", got, developmentSecureCookieKey)
	}
}
