package root

import (
	"os"
)

func logout() error {
	return os.Remove(filename)
}
