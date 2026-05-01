package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mhasnanr/ewallet-ums/bootstrap"
	"github.com/mhasnanr/ewallet-ums/constants"
)

type ExternalWallet struct{}

type WalletRequest struct {
	UserID int `json:"user_id"`
}

func (w *ExternalWallet) CreateWallet(userID int) error {
	var err error

	postBody, _ := json.Marshal(WalletRequest{UserID: userID})
	requestBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(fmt.Sprintf("%s/wallets/v1", bootstrap.GetEnv("WALLET_API_BASE_URL", "")), "application/json", requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		errorStr := fmt.Errorf("wallet service error: status %d: %s", resp.StatusCode, string(body))
		fmt.Println(errorStr)
		return constants.ErrorFailedToCreateWallet
	}

	return nil
}
