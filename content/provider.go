package content

import (
	"TgBot/config"
	"fmt"

	providerClient "github.com/Hogants/tg-bot-proto/gen/go/content"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Provider struct {
	Client providerClient.ProviderClient
}

func New(cfg config.ProviderConfig) (*Provider, error) {
	const op = "grpc.New"

	// Опции для интерсептора grpcretry
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(2)),
		grpcretry.WithPerRetryTimeout(cfg.Timeout),
	}
	cc, err := grpc.NewClient(cfg.Target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Provider{
		providerClient.NewProviderClient(cc),
	}, nil
}
