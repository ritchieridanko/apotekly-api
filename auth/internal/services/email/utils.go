package email

import "encoding/base64"

func encodeSubject(subject string) (encodedSubject string) {
	return "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?="
}
