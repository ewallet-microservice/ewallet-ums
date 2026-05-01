package external

import (
	"context"
	"errors"

	"github.com/mhasnanr/ewallet-ums/bootstrap"
	"github.com/mhasnanr/ewallet-ums/cmd/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WalletGRPC struct {
	client wallet.WalletClient
}

func NewWalletGRPC() (*WalletGRPC, *grpc.ClientConn, error) {
	serverAddr := bootstrap.GetEnv("WALLET_GRPC_URL", "")

	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, errors.New("failed to dial wallet grpc")
	}

	client := wallet.NewWalletClient(conn)

	return &WalletGRPC{
		client: client,
	}, conn, err
}

func (w *WalletGRPC) CreateWallet(ctx context.Context, userID int) (*wallet.CreateWalletResponse, error) {
	req := &wallet.CreateWalletRequest{
		UserId: int64(userID),
	}

	response, err := w.client.CreateWallet(ctx, req)
	if err != nil {
		return nil, errors.New("failed to create wallet")
	}

	return response, nil
}
