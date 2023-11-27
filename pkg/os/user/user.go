package user

import "os/user"

func HomeDir() (string, error) {
	currUser, err := user.Current()
	if err != nil {
		return "", err
	}

	return currUser.HomeDir, nil
}
