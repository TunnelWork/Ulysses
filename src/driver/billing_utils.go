package driver

import (
	"fmt"
	"strconv"
	"strings"

	harpocrates "github.com/TunnelWork/Harpocrates"
	"github.com/TunnelWork/Ulysses.Lib/billing"
	"github.com/TunnelWork/Ulysses.Lib/logging"
	"github.com/TunnelWork/Ulysses.Lib/payment"
)

func (Wallet) DepositPaymentOnStatusChange(referenceID string, newResult payment.PaymentResult) {
	// When Status shows PAID, deposit the amount specified into the wallet
	if newResult.Status == payment.PAID {
		// Parse the ReferenceID
		walletID, err := referenceID2WalletID(referenceID)
		if err != nil {
			logging.Error("driver: failed to parse referenceID: %s", referenceID)
			return
		}

		// get the wallet
		wallet, err := billing.GetWalletByID(walletID)
		if err != nil {
			logging.Error("driver: failed to get wallet: %s", err.Error())
			return
		}

		// Deposit the payment
		if newResult.Unit.Currency == "USD" {
			err = wallet.Deposit(newResult.Unit.Price)
			if err != nil {
				logging.Error("driver: failed to deposit payment: %s", err.Error())
				return
			}
		} else {
			logging.Error("driver: failed to deposit payment: unsupported currency: %s", newResult.Unit.Currency)
			return
		}
	} else {
		// In current version, just ignore
		logging.Debug("driver: ignoring payment status change to: %s", newResult.Status)
		return
	}
}

// Generate a random reference ID for a wallet ID
func walletID2ReferenceID(walletID uint64) string {
	// represent walletID in hex
	walletIDHex := fmt.Sprintf("%x", walletID)

	// get a random hex suffix
	suffix, err := harpocrates.GetRandomHex(8)
	if err != nil {
		suffix = "DEADBEEF"
	}

	return fmt.Sprintf("%s-%s", walletIDHex, suffix)
}

// Parse the walletID from the referenceID
func referenceID2WalletID(referenceID string) (uint64, error) {
	// Split the referenceID into parts
	parts := strings.Split(referenceID, "-")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid referenceID: %s", referenceID)
	}

	// Convert the walletID from hex
	walletID, err := strconv.ParseUint(parts[0], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid referenceID: %s", referenceID)
	}

	return walletID, nil
}
