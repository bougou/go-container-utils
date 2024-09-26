package container

import "os"

func FileExists(item string) (bool, error) {
	info, err := os.Stat(item)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// item exists
	return !info.IsDir(), nil
}
