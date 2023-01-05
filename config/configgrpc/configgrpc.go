// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configgrpc // import "go.opentelemetry.io/collector/config/configgrpc"

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/mostynb/go-grpc-compression/snappy"
	"github.com/mostynb/go-grpc-compression/zstd"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configauth"
	"go.opentelemetry.io/collector/config/configcompression"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/config/internal"
	"go.opentelemetry.io/collector/extension/auth"
)

var errMetadataNotFound = errors.New("no request metadata found")

// Allowed balancer names to be set in grpclb_policy to discover the servers.
var allowedBalancerNames = []string{roundrobin.Name, grpc.PickFirstBalancerName}

// KeepaliveClientConfig exposes the keepalive.ClientParameters to be used by the exporter.
// Refer to the original data-structure for the meaning of each parameter:
// https://godoc.org/google.golang.org/grpc/keepalive#ClientParameters
type KeepaliveClientConfig struct {
	Time                time.Duration `mapstructure:"time"`
	Timeout             time.Duration `mapstructure:"timeout"`
	PermitWithoutStream bool          `mapstructure:"permit_without_stream"`
}

// GRPCClientSettings defines common settings for a gRPC client configuration.
type GRPCClientSettings struct {
	Keepalive       *KeepaliveClientConfig            `mapstructure:"keepalive"`
	Headers         map[string]string                 `mapstructure:"headers"`
	Auth            *configauth.Authentication        `mapstructure:"auth"`
	Endpoint        string                            `mapstructure:"endpoint"`
	Compression     configcompression.CompressionType `mapstructure:"compression"`
	BalancerName    string                            `mapstructure:"balancer_name"`
	TLSSetting      configtls.TLSClientSetting        `mapstructure:"tls"`
	ReadBufferSize  int                               `mapstructure:"read_buffer_size"`
	WriteBufferSize int                               `mapstructure:"write_buffer_size"`
	WaitForReady    bool                              `mapstructure:"wait_for_ready"`
}

// KeepaliveServerConfig is the configuration for keepalive.
type KeepaliveServerConfig struct {
	ServerParameters  *KeepaliveServerParameters  `mapstructure:"server_parameters"`
	EnforcementPolicy *KeepaliveEnforcementPolicy `mapstructure:"enforcement_policy"`
}

// KeepaliveServerParameters allow configuration of the keepalive.ServerParameters.
// The same default values as keepalive.ServerParameters are applicable and get applied by the server.
// See https://godoc.org/google.golang.org/grpc/keepalive#ServerParameters for details.
type KeepaliveServerParameters struct {
	MaxConnectionIdle     time.Duration `mapstructure:"max_connection_idle"`
	MaxConnectionAge      time.Duration `mapstructure:"max_connection_age"`
	MaxConnectionAgeGrace time.Duration `mapstructure:"max_connection_age_grace"`
	Time                  time.Duration `mapstructure:"time"`
	Timeout               time.Duration `mapstructure:"timeout"`
}

// KeepaliveEnforcementPolicy allow configuration of the keepalive.EnforcementPolicy.
// The same default values as keepalive.EnforcementPolicy are applicable and get applied by the server.
// See https://godoc.org/google.golang.org/grpc/keepalive#EnforcementPolicy for details.
type KeepaliveEnforcementPolicy struct {
	MinTime             time.Duration `mapstructure:"min_time"`
	PermitWithoutStream bool          `mapstructure:"permit_without_stream"`
}

// GRPCServerSettings defines common settings for a gRPC server configuration.
type GRPCServerSettings struct {
	TLSSetting           *configtls.TLSServerSetting `mapstructure:"tls"`
	Keepalive            *KeepaliveServerConfig      `mapstructure:"keepalive"`
	Auth                 *configauth.Authentication  `mapstructure:"auth"`
	NetAddr              confignet.NetAddr           `mapstructure:",squash"`
	MaxRecvMsgSizeMiB    uint64                      `mapstructure:"max_recv_msg_size_mib"`
	ReadBufferSize       int                         `mapstructure:"read_buffer_size"`
	WriteBufferSize      int                         `mapstructure:"write_buffer_size"`
	MaxConcurrentStreams uint32                      `mapstructure:"max_concurrent_streams"`
	IncludeMetadata      bool                        `mapstructure:"include_metadata"`
}

// SanitizedEndpoint strips the prefix of either http:// or https:// from configgrpc.GRPCClientSettings.Endpoint.
func (gcs *GRPCClientSettings) SanitizedEndpoint() string {
	switch {
	case gcs.isSchemeHTTP():
		return strings.TrimPrefix(gcs.Endpoint, "http://")
	case gcs.isSchemeHTTPS():
		return strings.TrimPrefix(gcs.Endpoint, "https://")
	default:
		return gcs.Endpoint
	}
}

func (gcs *GRPCClientSettings) isSchemeHTTP() bool {
	return strings.HasPrefix(gcs.Endpoint, "http://")
}

func (gcs *GRPCClientSettings) isSchemeHTTPS() bool {
	return strings.HasPrefix(gcs.Endpoint, "https://")
}

// ToClientConn creates a client connection to the given target. By default, it's
// a non-blocking dial (the function won't wait for connections to be
// established, and connecting happens in the background). To make it a blocking
// dial, use grpc.WithBlock() dial option.
func (gcs *GRPCClientSettings) ToClientConn(ctx context.Context, host component.Host, settings component.TelemetrySettings, extraOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts, err := gcs.toDialOptions(host, settings)
	if err != nil {
		return nil, err
	}
	opts = append(opts, extraOpts...)
	return grpc.DialContext(ctx, gcs.SanitizedEndpoint(), opts...)
}

func (gcs *GRPCClientSettings) toDialOptions(host component.Host, settings component.TelemetrySettings) ([]grpc.DialOption, error) {
	var opts []grpc.DialOption
	if configcompression.IsCompressed(gcs.Compression) {
		cp, err := getGRPCCompressionName(gcs.Compression)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.UseCompressor(cp)))
	}

	tlsCfg, err := gcs.TLSSetting.LoadTLSConfig()
	if err != nil {
		return nil, err
	}
	cred := insecure.NewCredentials()
	if tlsCfg != nil {
		cred = credentials.NewTLS(tlsCfg)
	} else if gcs.isSchemeHTTPS() {
		cred = credentials.NewTLS(&tls.Config{})
	}
	opts = append(opts, grpc.WithTransportCredentials(cred))

	if gcs.ReadBufferSize > 0 {
		opts = append(opts, grpc.WithReadBufferSize(gcs.ReadBufferSize))
	}

	if gcs.WriteBufferSize > 0 {
		opts = append(opts, grpc.WithWriteBufferSize(gcs.WriteBufferSize))
	}

	if gcs.Keepalive != nil {
		keepAliveOption := grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                gcs.Keepalive.Time,
			Timeout:             gcs.Keepalive.Timeout,
			PermitWithoutStream: gcs.Keepalive.PermitWithoutStream,
		})
		opts = append(opts, keepAliveOption)
	}

	if gcs.Auth != nil {
		if host.GetExtensions() == nil {
			return nil, errors.New("no extensions configuration available")
		}

		grpcAuthenticator, cerr := gcs.Auth.GetClientAuthenticator(host.GetExtensions())
		if cerr != nil {
			return nil, cerr
		}

		perRPCCredentials, perr := grpcAuthenticator.PerRPCCredentials()
		if perr != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithPerRPCCredentials(perRPCCredentials))
	}

	if gcs.BalancerName != "" {
		valid := validateBalancerName(gcs.BalancerName)
		if !valid {
			return nil, fmt.Errorf("invalid balancer_name: %s", gcs.BalancerName)
		}
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, gcs.BalancerName)))
	}

	otelOpts := []otelgrpc.Option{
		otelgrpc.WithTracerProvider(settings.TracerProvider),
		otelgrpc.WithMeterProvider(settings.MeterProvider),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	}

	// Enable OpenTelemetry observability plugin.
	opts = append(opts, grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor(otelOpts...)))
	opts = append(opts, grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor(otelOpts...)))

	return opts, nil
}

func validateBalancerName(balancerName string) bool {
	for _, item := range allowedBalancerNames {
		if item == balancerName {
			return true
		}
	}
	return false
}

// ToListener returns the net.Listener constructed from the settings.
func (gss *GRPCServerSettings) ToListener() (net.Listener, error) {
	return gss.NetAddr.Listen()
}

func (gss *GRPCServerSettings) ToServer(host component.Host, settings component.TelemetrySettings, extraOpts ...grpc.ServerOption) (*grpc.Server, error) {
	opts, err := gss.toServerOption(host, settings)
	if err != nil {
		return nil, err
	}
	opts = append(opts, extraOpts...)
	return grpc.NewServer(opts...), nil
}

func (gss *GRPCServerSettings) toServerOption(host component.Host, settings component.TelemetrySettings) ([]grpc.ServerOption, error) {
	switch gss.NetAddr.Transport {
	case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6":
		internal.WarnOnUnspecifiedHost(settings.Logger, gss.NetAddr.Endpoint)
	}

	var opts []grpc.ServerOption

	if gss.TLSSetting != nil {
		tlsCfg, err := gss.TLSSetting.LoadTLSConfig()
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsCfg)))
	}

	if gss.MaxRecvMsgSizeMiB > 0 {
		opts = append(opts, grpc.MaxRecvMsgSize(int(gss.MaxRecvMsgSizeMiB*1024*1024)))
	}

	if gss.MaxConcurrentStreams > 0 {
		opts = append(opts, grpc.MaxConcurrentStreams(gss.MaxConcurrentStreams))
	}

	if gss.ReadBufferSize > 0 {
		opts = append(opts, grpc.ReadBufferSize(gss.ReadBufferSize))
	}

	if gss.WriteBufferSize > 0 {
		opts = append(opts, grpc.WriteBufferSize(gss.WriteBufferSize))
	}

	// The default values referenced in the GRPC docs are set within the server, so this code doesn't need
	// to apply them over zero/nil values before passing these as grpc.ServerOptions.
	// The following shows the server code for applying default grpc.ServerOptions.
	// https://github.com/grpc/grpc-go/blob/120728e1f775e40a2a764341939b78d666b08260/internal/transport/http2_server.go#L184-L200
	if gss.Keepalive != nil {
		if gss.Keepalive.ServerParameters != nil {
			svrParams := gss.Keepalive.ServerParameters
			opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle:     svrParams.MaxConnectionIdle,
				MaxConnectionAge:      svrParams.MaxConnectionAge,
				MaxConnectionAgeGrace: svrParams.MaxConnectionAgeGrace,
				Time:                  svrParams.Time,
				Timeout:               svrParams.Timeout,
			}))
		}
		// The default values referenced in the GRPC are set within the server, so this code doesn't need
		// to apply them over zero/nil values before passing these as grpc.ServerOptions.
		// The following shows the server code for applying default grpc.ServerOptions.
		// https://github.com/grpc/grpc-go/blob/120728e1f775e40a2a764341939b78d666b08260/internal/transport/http2_server.go#L202-L205
		if gss.Keepalive.EnforcementPolicy != nil {
			enfPol := gss.Keepalive.EnforcementPolicy
			opts = append(opts, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				MinTime:             enfPol.MinTime,
				PermitWithoutStream: enfPol.PermitWithoutStream,
			}))
		}
	}

	var uInterceptors []grpc.UnaryServerInterceptor
	var sInterceptors []grpc.StreamServerInterceptor

	if gss.Auth != nil {
		authenticator, err := gss.Auth.GetServerAuthenticator(host.GetExtensions())
		if err != nil {
			return nil, err
		}

		uInterceptors = append(uInterceptors, func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return authUnaryServerInterceptor(ctx, req, info, handler, authenticator)
		})
		sInterceptors = append(sInterceptors, func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return authStreamServerInterceptor(srv, ss, info, handler, authenticator)
		})
	}

	otelOpts := []otelgrpc.Option{
		otelgrpc.WithTracerProvider(settings.TracerProvider),
		otelgrpc.WithMeterProvider(settings.MeterProvider),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	}

	// Enable OpenTelemetry observability plugin.
	// TODO: Pass construct settings to have access to Tracer.
	uInterceptors = append(uInterceptors, otelgrpc.UnaryServerInterceptor(otelOpts...))
	sInterceptors = append(sInterceptors, otelgrpc.StreamServerInterceptor(otelOpts...))

	uInterceptors = append(uInterceptors, enhanceWithClientInformation(gss.IncludeMetadata))
	sInterceptors = append(sInterceptors, enhanceStreamWithClientInformation(gss.IncludeMetadata))

	opts = append(opts, grpc.ChainUnaryInterceptor(uInterceptors...), grpc.ChainStreamInterceptor(sInterceptors...))

	return opts, nil
}

// getGRPCCompressionName returns compression name registered in grpc.
func getGRPCCompressionName(compressionType configcompression.CompressionType) (string, error) {
	switch compressionType {
	case configcompression.Gzip:
		return gzip.Name, nil
	case configcompression.Snappy:
		return snappy.Name, nil
	case configcompression.Zstd:
		return zstd.Name, nil
	default:
		return "", fmt.Errorf("unsupported compression type %q", compressionType)
	}
}

// enhanceWithClientInformation intercepts the incoming RPC, replacing the incoming context with one that includes
// a client.Info, potentially with the peer's address.
func enhanceWithClientInformation(includeMetadata bool) func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(contextWithClient(ctx, includeMetadata), req)
	}
}

func enhanceStreamWithClientInformation(includeMetadata bool) func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, wrapServerStream(contextWithClient(ss.Context(), includeMetadata), ss))
	}
}

// contextWithClient attempts to add the peer address to the client.Info from the context. When no
// client.Info exists in the context, one is created.
func contextWithClient(ctx context.Context, includeMetadata bool) context.Context {
	cl := client.FromContext(ctx)
	if p, ok := peer.FromContext(ctx); ok {
		cl.Addr = p.Addr
	}
	if includeMetadata {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			copiedMD := md.Copy()
			if len(md[client.MetadataHostName]) == 0 && len(md[":authority"]) > 0 {
				copiedMD[client.MetadataHostName] = md[":authority"]
			}
			cl.Metadata = client.NewMetadata(copiedMD)
		}
	}
	return client.NewContext(ctx, cl)
}

func authUnaryServerInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler, server auth.Server) (interface{}, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMetadataNotFound
	}

	ctx, err := server.Authenticate(ctx, headers)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func authStreamServerInterceptor(srv interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler, server auth.Server) error {
	ctx := stream.Context()
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errMetadataNotFound
	}

	ctx, err := server.Authenticate(ctx, headers)
	if err != nil {
		return err
	}

	return handler(srv, wrapServerStream(ctx, stream))
}
