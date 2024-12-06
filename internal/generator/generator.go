package generator

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"

	"crypto/hmac"
	"crypto/sha512"

	"github.com/xssnick/tonutils-go/ton/wallet"
	"golang.org/x/crypto/pbkdf2"
)

type GeneratorConfig struct {
	Suffix    string
	CaseSense bool
	Workchain int
}

type WalletData struct {
	Address    string
	PrivateKey string
	PublicKey  string
	Mnemonic   string
}

type MnemonicInfo struct {
	Mnemonic   string
	PrivateKey string
	PublicKey  string
}

func GenerateWalletMnemonic() (MnemonicInfo, error) {
	mnemonic := wallet.NewSeed()

	mac := hmac.New(sha512.New, []byte(strings.Join(mnemonic, " ")))
	hash := mac.Sum(nil)
	k := pbkdf2.Key(hash, []byte("TON default seed"), 100000, 32, sha512.New)
	privateKey := ed25519.NewKeyFromSeed(k)
	publicKey := privateKey.Public().(ed25519.PublicKey)

	return MnemonicInfo{
		Mnemonic:   strings.Join(mnemonic, " "),
		PrivateKey: fmt.Sprintf("%x", k),
		PublicKey:  fmt.Sprintf("%x", publicKey),
	}, nil
}

func GenerateWalletData(workchain int) (*WalletData, error) {
	mnemonic_info, _ := GenerateWalletMnemonic()

	config := wallet.ConfigV5R1Final{
		NetworkGlobalID: -239,
		Workchain:       int8(workchain),
	}

	pub, err := hex.DecodeString(mnemonic_info.PublicKey)

	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	addr, err := wallet.AddressFromPubKey(pub, config, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}
	addr.SetBounce(false)

	return &WalletData{
		Address:   addr.String(),
		PublicKey: mnemonic_info.PublicKey,
		Mnemonic:  mnemonic_info.Mnemonic,
	}, nil
}

func CheckAddressSuffix(addr string, config GeneratorConfig) bool {
	suffix := config.Suffix

	if !config.CaseSense {
		addr = strings.ToLower(addr)
		suffix = strings.ToLower(suffix)
	}
	return strings.HasSuffix(addr, suffix)
}
