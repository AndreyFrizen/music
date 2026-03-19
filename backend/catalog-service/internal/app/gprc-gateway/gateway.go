package gateway

import (
	configate "catalog-service/config/gateway"
	"catalog-service/proto/catalog"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson" // ADD THIS IMPORT
)

type GatewayApp struct {
	router *gin.Engine
	log    *slog.Logger
	config *configate.Config
	conn   *grpc.ClientConn
	mux    *runtime.ServeMux
	server *http.Server
}

func NewGatewayApp(log *slog.Logger, config *configate.Config, opts ...grpc.DialOption) (*GatewayApp, error) {
	const op = "GatewayApp.New"

	defaultOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10 * 1024 * 1024)),
	}
	allOpts := append(defaultOpts, opts...)

	conn, err := grpc.NewClient(config.GRPCPort, allOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %s, %w", op, err)
	}

	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		// FIX: Add custom error handler to properly translate gRPC errors to HTTP
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
			st := status.Convert(err)
			httpStatus := runtime.HTTPStatusFromCode(st.Code())

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(httpStatus)

			// Return clean JSON error instead of gRPC status
			json.NewEncoder(w).Encode(map[string]any{
				"error":  st.Message(),
				"code":   int(st.Code()),
				"status": httpStatus,
			})
		}),
		runtime.WithOutgoingHeaderMatcher(func(key string) (string, bool) {
			allowed := map[string]bool{
				"x-request-id": true,
				"x-trace-id":   true,
			}
			if allowed[key] {
				return key, true
			}
			return runtime.DefaultHeaderMatcher(key)
		}),
	)

	ctx := context.Background()
	if err := catalog.RegisterCatalogServiceHandler(ctx, gwmux, conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to register catalog service handler: %w", err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Any("/api/v1/*path", gin.WrapH(gwmux))

	router.GET("/health", func(c *gin.Context) {
		state := conn.GetState()
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"grpc":   state.String(),
		})
	})

	server := &http.Server{
		Addr:         config.HTTPPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &GatewayApp{
		router: router,
		log:    log,
		config: config,
		conn:   conn,
		mux:    gwmux,
		server: server,
	}, nil
}

func (g *GatewayApp) Run() error {
	const op = "GatewayApp.Run"

	g.log.Info("Starting HTTP gateway",
		"port", g.config.HTTPPort,
		"grpc_endpoint", g.config.GRPCPort,
	)

	if err := g.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (g *GatewayApp) Stop(ctx context.Context) error {
	const op = "GatewayApp.Stop"

	g.log.Info("Shutting down gateway...")

	if err := g.server.Shutdown(ctx); err != nil {
		g.log.Error("HTTP server shutdown error", "error", err)
	}

	if err := g.conn.Close(); err != nil {
		g.log.Error("gRPC connection close error", "error", err)
	}

	g.log.Info("Gateway stopped")
	return nil
}

// wrapGateway adapts runtime.ServeMux to Gin handler
// func wrapGateway(gwmux *runtime.ServeMux) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Let gRPC-Gateway handle the request
// 		gwmux.ServeHTTP(c.Writer, c.Request)
// 	}
// }

// grpcClientInterceptor adds tracing to gRPC calls
// func grpcClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
// 	// Add request ID to outgoing context
// 	if md, ok := metadata.FromOutgoingContext(ctx); ok {
// 		if vals := md.Get("x-request-id"); len(vals) > 0 {
// 			ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", vals[0])
// 		}
// 	}

// 	start := time.Now()
// 	err := invoker(ctx, method, req, reply, cc, opts...)
// 	log.Printf("gRPC call: %s, duration: %v, error: %v", method, time.Since(start), err)
// 	return err
// }
