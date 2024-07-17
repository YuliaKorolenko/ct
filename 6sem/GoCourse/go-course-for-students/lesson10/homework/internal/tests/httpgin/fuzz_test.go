package httpgin

import (
	"testing"
)

func FuzzTest(f *testing.F) {
	client := getTestClient()

	f.Fuzz(func(t *testing.T, name string) {
		email := name + "gmail.com"
		user, err := client.createUser(name, email)
		if err == nil {
			if len(user.Data.Nickname) == 0 {
				t.Errorf("Created user with epmty name. Expect: %s, but got: %s", "error", "nil")
			}
			if len(user.Data.Email) == 0 {
				t.Errorf("Created user with epmty name. Expect: %s, but got: %s", "error", "nil")
			}
			if len(user.Data.Nickname) > 100 {
				t.Errorf("Created user with to long name. Expect: %s, but got: %s", "error", "nil")
			}
			if len(user.Data.Email) > 100 {
				t.Errorf("Created user with to long email. Expect: %s, but got: %s", "error", "nil")
			}
		}
	})
}
