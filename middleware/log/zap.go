package log

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ZapUnaryServerInterceptor ...
func ZapUnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var (
			startTime = time.Now()
		)
		defer func() {
			var (
				code  = grpc.Code(err)
				level = levelFromCode(code)
			)
			logger.Check(level, fmt.Sprintf("%s: %s", code.String(), info.FullMethod)).
				Write(
					zap.Error(err),
					zap.String("grpc.code", code.String()),
					zap.Float32("grpc.duration_ms", float32(time.Since(startTime).Nanoseconds())/1000000.0),
					zap.String("grpc.method", info.FullMethod),
					zap.Reflect("grpc.request", req),
					zap.Reflect("grpc.response", resp),
				)
		}()
		resp, err = handler(ctx, req)
		return
	}
}

// ZapStreamServerInterceptor ...
func ZapStreamServerInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		var (
			startTime = time.Now()
		)
		defer func() {
			var (
				code  = grpc.Code(err)
				level = levelFromCode(code)
			)
			logger.Check(level, fmt.Sprintf("%s: %s", code.String(), info.FullMethod)).
				Write(
					zap.Error(err),
					zap.String("grpc.code", code.String()),
					zap.Float32("grpc.duration_ms", float32(time.Since(startTime).Nanoseconds())/1000000.0),
					zap.String("grpc.method", info.FullMethod),
				)
		}()
		err = handler(srv, stream)
		return
	}
}

// wrappedServerStream is a thin wrapper around grpc.ServerStream that allows modifying context.
type wrappedServerStream struct {
	grpc.ServerStream
	// WrappedContext is the wrapper's own Context. You can assign it.
	WrappedContext context.Context
}

// Context returns the wrapper's WrappedContext, overwriting the nested grpc.ServerStream.Context()
func (w *wrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

// TODO: unused function
// WrapServerStream returns a ServerStream that has the ability to overwrite context.
func wrapServerStream(stream grpc.ServerStream) *wrappedServerStream {
	if existing, ok := stream.(*wrappedServerStream); ok {
		return existing
	}
	return &wrappedServerStream{ServerStream: stream, WrappedContext: stream.Context()}
}

func levelFromCode(code codes.Code) zapcore.Level {
	switch code {
	case codes.OK:
		return zap.InfoLevel
	case codes.Canceled:
		return zap.InfoLevel
	case codes.Unknown:
		return zap.ErrorLevel
	case codes.InvalidArgument:
		return zap.InfoLevel
	case codes.DeadlineExceeded:
		return zap.WarnLevel
	case codes.NotFound:
		return zap.InfoLevel
	case codes.AlreadyExists:
		return zap.InfoLevel
	case codes.PermissionDenied:
		return zap.WarnLevel
	case codes.Unauthenticated:
		return zap.InfoLevel // unauthenticated requests can happen
	case codes.ResourceExhausted:
		return zap.WarnLevel
	case codes.FailedPrecondition:
		return zap.WarnLevel
	case codes.Aborted:
		return zap.WarnLevel
	case codes.OutOfRange:
		return zap.WarnLevel
	case codes.Unimplemented:
		return zap.ErrorLevel
	case codes.Internal:
		return zap.ErrorLevel
	case codes.Unavailable:
		return zap.WarnLevel
	case codes.DataLoss:
		return zap.ErrorLevel
	default:
		return zap.ErrorLevel
	}
}
