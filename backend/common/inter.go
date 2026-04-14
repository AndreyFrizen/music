package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// contextKey - тип для ключей контекста (чтобы избежать коллизий)
type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

// AuthUnaryInterceptor - главный интерсептор для всех сервисов (кроме User Service)
func AuthUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	// 1. Пропускаем публичные методы (если нужно)
	if isPublicMethod(info.FullMethod) {
		return handler(ctx, req)
	}

	// 2. Извлекаем токен из запроса
	token, err := extractToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication required: %v", err)
	}

	// 3. Проверяем токен через User Service
	client := GetUserServiceClient()
	if client == nil {
		return nil, status.Errorf(codes.Internal, "user service client not initialized")
	}

	isValid, userID, role, err := client.ValidateToken(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate token: %v", err)
	}

	if !isValid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid or expired token")
	}

	ctx = context.WithValue(ctx, UserIDKey, userID)
	ctx = context.WithValue(ctx, UserRoleKey, role)

	return handler(ctx, req)
}

func extractToken(ctx context.Context) (string, error) {
	// Получаем metadata из контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no metadata in request")
	}

	// Ищем заголовок authorization (ВНИМАНИЕ: в нижнем регистре!)
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return "", status.Error(codes.Unauthenticated, "authorization header not found")
	}

	// Формат: "Bearer <token>"
	authHeader := authHeaders[0]
	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", status.Error(codes.Unauthenticated, "invalid authorization header format. Expected 'Bearer <token>'")
	}

	return parts[1], nil
}

// isPublicMethod - определяет, какие методы не требуют авторизации
func isPublicMethod(method string) bool {
	// Эти методы доступны без токена
	publicMethods := map[string]bool{
		"/health.Health/Check":    true,
		"/user.UserService/Login": true,
		// Добавьте другие публичные методы
	}

	return publicMethods[method]
}
