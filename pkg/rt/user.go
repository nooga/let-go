//go:build !js

/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import "os/user"

func currentUser() (home, name string) {
	if u, err := user.Current(); err == nil {
		return u.HomeDir, u.Username
	}
	return "", ""
}
