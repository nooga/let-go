/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import "encoding/base64"

func decodeBase64URL(s string) ([]byte, error) {
	data, err := base64.RawURLEncoding.DecodeString(s)
	if err == nil {
		return data, nil
	}
	data, paddedErr := base64.URLEncoding.DecodeString(s)
	if paddedErr == nil {
		return data, nil
	}
	return nil, err
}
