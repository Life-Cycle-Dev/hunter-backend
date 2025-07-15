package util

import (
	"errors"
	"fmt"
)

func CheckAllow(userPermission []string, requirePermission string) error {
	for _, perm := range userPermission {
		if perm == requirePermission {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("permission denied (%s)", requirePermission))
}
