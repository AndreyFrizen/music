package auth

import (
	"context"
	"music/user-service/proto/user"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserServiceClient - клиент для общения с User Service
type UserServiceClient struct {
	client user.UserServiceClient
	conn   *grpc.ClientConn
}

var (
	instance *UserServiceClient
	once     sync.Once
	initErr  error
)

// InitUserServiceClient - инициализирует клиент (вызывается в main каждого сервиса)
func InitUserServiceClient(userServiceAddr string) error {
	once.Do(func() {
		conn, err := grpc.NewClient(userServiceAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			initErr = err
			return
		}

		instance = &UserServiceClient{
			client: user.NewUserServiceClient(conn),
			conn:   conn,
		}
	})
	return initErr
}

func GetUserServiceClient() *UserServiceClient {
	return instance
}

// ValidateToken - проверяет токен через User Service
func (c *UserServiceClient) ValidateToken(ctx context.Context, token string) (bool, int64, string, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.client.ValidateToken(ctx, &user.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return false, 0, "", err
	}

	return resp.Valid, resp.UserId, resp.Role, nil
}

// Close - закрывает соединение
func (c *UserServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
