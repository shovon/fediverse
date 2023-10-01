package accountaddress

import "testing"

func TestParsing(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		a, err := ParseAccountAddress("@user@host")
		if err != nil {
			t.Fatal(err)
		}
		if a.User != "user" {
			t.Error("user not parsed correctly")
			t.Fail()
		}
		if a.Host != "host" {
			t.Error("host not parsed correctly")
			t.Fail()
		}
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := ParseAccountAddress("user@host")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
