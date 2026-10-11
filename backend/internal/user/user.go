package user

import "sync"

var (
	usernameMu sync.Mutex
	usernameDB = make(map[string]struct{})
)

func ClaimUsername(name string) bool {
	usernameMu.Lock()
	defer usernameMu.Unlock()

	if _, exists := usernameDB[name]; exists {
		return false
	}
	usernameDB[name] = struct{}{}
	return true
}

func ReleaseUsername(name string) {
	usernameMu.Lock()
	defer usernameMu.Unlock()

	delete(usernameDB, name)
}
