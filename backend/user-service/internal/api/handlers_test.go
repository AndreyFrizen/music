// handlers_test.go
package handlers

import (
	"testing"
	"user-service/proto/user"
)

func TestServerAPI_ImplementsUserServiceServer(t *testing.T) {
	var _ user.UserServiceServer = (*serverAPI)(nil)
}
