package filesystem

import "bytes"

func heal(content []byte) []byte {
	return bytes.ReplaceAll(content, []byte("}}"), []byte("}"))
}
