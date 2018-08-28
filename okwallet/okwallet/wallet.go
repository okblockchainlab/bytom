package okwallet

import (
  "encoding/hex"
  "fmt"
  "errors"

  "github.com/bytom/crypto/ed25519/chainkd"
  "github.com/bytom/account"
	"github.com/bytom/blockchain/signers"
	"github.com/bytom/common"
	"github.com/bytom/crypto"
	"github.com/bytom/consensus"
)

//reference: API.createAccount, API.createAccountReceiver
func GetAddressByPrivateKey(prvkeyStr string) (string, error) {
  prvkeyBytes, err := hex.DecodeString(prvkeyStr)
  if err != nil {
    return "", fmt.Errorf("decoding private key error: %v", err)
  }

  var xprv chainkd.XPrv
  if len(xprv) != len(prvkeyBytes) {
    return "", errors.New("invalid private key")
  }
  copy(xprv[:], prvkeyBytes)

  xpub := xprv.XPub()

  const quorum = 1
	signer, err := signers.Create("account", []chainkd.XPub{xpub}, quorum, 1)
	id := signers.IDGenerate()
	if err != nil {
		return "", fmt.Errorf("generate id for account failed: %v", err)
	}

	account := &account.Account{Signer: signer, ID: id, Alias: "noused"}

  const idx = 1
	path := signers.Path(account.Signer, signers.AccountKeySpace, idx)
	derivedXPubs := chainkd.DeriveXPubs(account.XPubs, path)
	derivedPK := derivedXPubs[0].PublicKey()
	pubHash := crypto.Ripemd160(derivedPK)

	address, err := common.NewAddressWitnessPubKeyHash(pubHash, &consensus.ActiveNetParams)
  return address.EncodeAddress(), nil
}
