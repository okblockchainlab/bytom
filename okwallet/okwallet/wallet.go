package okwallet

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bytom/account"
	"github.com/bytom/api"
	"github.com/bytom/blockchain/signers"
	"github.com/bytom/blockchain/txbuilder"
	"github.com/bytom/common"
	"github.com/bytom/consensus"
	"github.com/bytom/crypto"
	"github.com/bytom/crypto/ed25519/chainkd"
	btm_json "github.com/bytom/encoding/json"
	btm_errors "github.com/bytom/errors"
	"github.com/bytom/protocol/bc"
	"github.com/bytom/protocol/bc/types"
)

const (
	defaultTxTTL = 5 * time.Minute
)

type OKUTXO struct {
	OutputID            bc.Hash    `json:"id"`
	SourceID            bc.Hash    `json:"source_id"`
	AssetID             bc.AssetID `json:"asset_id"`
	Amount              uint64     `json:"amount"`
	SourcePos           uint64     `json:"source_pos"`
	ControlProgram      []byte     `json:"program"`
	AccountID           string     `json:"account_id"`
	Address             string     `json:"address"`
	ControlProgramIndex uint64     `json:"control_program_index"`
	ValidHeight         uint64     `json:"valid_height"`
	Change              bool       `json:"change"`
}

type OKAction struct {
	Type string
	Utxo OKUTXO
}

type OKBuildRequest struct {
	Tx *types.TxData `json:"base_transaction"`
	//Actions   []OKAction    `json:"actions"`
	Actions   []map[string]interface{} `json:"actions"`
	XPub      chainkd.XPub             `json:xpub`
	TTL       btm_json.Duration        `json:"ttl"`
	TimeRange uint64                   `json:"time_range"`
}

type GetAddressResp struct {
	XPub    string `json:xpub`
	Address string `json:address`
}

type signResp struct {
	Tx           *txbuilder.Template `json:"transaction"`
	SignComplete bool                `json:"sign_complete"`
}

type spendUTXOAction struct {
	UTXO *account.UTXO
	XPub *chainkd.XPub
}

func (a *spendUTXOAction) Build(_ context.Context, b *txbuilder.TemplateBuilder) error {
	accountSigner, err := createSigner([]chainkd.XPub{*a.XPub})
	if err != nil {
		return err
	}

	txInput, sigInst, err := account.UtxoToInputs(accountSigner, a.UTXO)
	if err != nil {
		return err
	}

	return b.AddInput(txInput, sigInst)
}

func createSigner(xpubs []chainkd.XPub) (*signers.Signer, error) {
	return signers.Create("account", xpubs, 1, 1)
}

func okUTXO2UTXO(oktx *OKUTXO) *account.UTXO {
	return &account.UTXO{
		OutputID:            oktx.OutputID,
		SourceID:            oktx.SourceID,
		AssetID:             oktx.AssetID,
		Amount:              oktx.Amount,
		SourcePos:           oktx.SourcePos,
		ControlProgram:      oktx.ControlProgram,
		AccountID:           oktx.AccountID,
		Address:             oktx.Address,
		ControlProgramIndex: oktx.ControlProgramIndex,
		ValidHeight:         oktx.ValidHeight,
		Change:              oktx.Change,
	}
}

func decodeSpendUTXOAction(data []byte, xpub *chainkd.XPub) (txbuilder.Action, error) {
	var okact OKAction
	if err := json.Unmarshal(data, &okact); err != nil {
		return nil, fmt.Errorf("unmarshal OKAction error. %v", err)
	}

	a := &spendUTXOAction{okUTXO2UTXO(&okact.Utxo), xpub}
	return a, nil
}

func decodeControlAddressAction(data []byte, _ *chainkd.XPub) (txbuilder.Action, error) {
	return txbuilder.DecodeControlAddressAction(data)
}

func decodeControlProgramAction(data []byte, _ *chainkd.XPub) (txbuilder.Action, error) {
	return txbuilder.DecodeControlProgramAction(data)
}

func decodeRetireAction(data []byte, _ *chainkd.XPub) (txbuilder.Action, error) {
	return txbuilder.DecodeRetireAction(data)
}

func actionDecoder(action string) (func([]byte, *chainkd.XPub) (txbuilder.Action, error), bool) {
	decoders := map[string]func([]byte, *chainkd.XPub) (txbuilder.Action, error){
		"control_address": decodeControlAddressAction,
		"control_program": decodeControlProgramAction,
		//"issue":                        a.wallet.AssetReg.DecodeIssueAction,
		"retire": decodeRetireAction,
		//"spend_account":                a.wallet.AccountMgr.DecodeSpendAction,
		"spend_account_unspent_output": decodeSpendUTXOAction,
	}
	decoder, ok := decoders[action]
	return decoder, ok
}

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
	signer, err := createSigner([]chainkd.XPub{xpub})
	if err != nil {
		return "", err
	}

	account := &account.Account{Signer: signer, ID: signers.IDGenerate(), Alias: "noused"}

	const idx = 1
	path := signers.Path(account.Signer, signers.AccountKeySpace, idx)
	derivedXPubs := chainkd.DeriveXPubs(account.XPubs, path)
	derivedPK := derivedXPubs[0].PublicKey()
	pubHash := crypto.Ripemd160(derivedPK)

	address, err := common.NewAddressWitnessPubKeyHash(pubHash, &consensus.ActiveNetParams)

	getAddrResp := GetAddressResp{XPub: hex.EncodeToString(xpub[:]), Address: address.EncodeAddress()}
	result, err := json.Marshal(getAddrResp)
	return string(result), nil
}

//reference: API.build
// action{type, output_id}, {amount, asset_id, address, type}, ttl, time_range
func CreateRawTransaction(reqStr string) (string, error) {
	var req OKBuildRequest
	if err := json.Unmarshal([]byte(reqStr), &req); err != nil {
		return "", err
	}

	actions := make([]txbuilder.Action, 0, len(req.Actions))
	for i, act := range req.Actions {
		typ, ok := act["type"].(string)
		if !ok {
			return "", btm_errors.WithDetailf(api.ErrBadActionType, "no action type provided on action %d", i)
		}

		decoder, ok := actionDecoder(typ)
		if !ok {
			return "", errors.New("unknown action type " + typ)
		}

		// Remarshal to JSON, the action may have been modified when we
		// filtered aliases.
		b, err := json.Marshal(act)
		if err != nil {
			return "", err
		}

		action, err := decoder(b, &req.XPub)
		if err != nil {
			return "", btm_errors.WithDetailf(api.ErrBadAction, "%s on action %d", err.Error(), i)
		}

		actions = append(actions, action)
	}
	actions = account.MergeSpendAction(actions)

	ttl := req.TTL.Duration
	if ttl == 0 {
		ttl = defaultTxTTL
	}
	maxTime := time.Now().Add(ttl)

	tpl, err := txbuilder.Build(nil, req.Tx, actions, maxTime, req.TimeRange)
	if btm_errors.Root(err) == txbuilder.ErrAction {
		// append each of the inner btm_errors contained in the data.
		var Errs string
		var rootErr error
		for i, innerErr := range btm_errors.Data(err)["actions"].([]error) {
			if i == 0 {
				rootErr = btm_errors.Root(innerErr)
			}
			Errs = Errs + innerErr.Error()
		}
		err = btm_errors.WithDetail(rootErr, Errs)
	}
	if err != nil {
		return "", err
	}

	// ensure null is never returned for signing instructions
	if tpl.SigningInstructions == nil {
		tpl.SigningInstructions = []*txbuilder.SigningInstruction{}
	}

	result, err := json.Marshal(tpl)
	return string(result), err
}

func pseudohsmSignTemplate(_ context.Context, _ chainkd.XPub, path [][]byte, data [32]byte, prvkeyStr string) ([]byte, error) {
	prvkeyBytes, err := hex.DecodeString(prvkeyStr)
	if err != nil {
		return nil, fmt.Errorf("decoding private key error: %v", err)
	}

	var xprv chainkd.XPrv
	if len(xprv) != len(prvkeyBytes) {
		return nil, errors.New("invalid private key")
	}
	copy(xprv[:], prvkeyBytes)

	if len(path) > 0 {
		xprv = xprv.Derive(path)
	}
	return xprv.Sign(data[:]), nil
}

//reference: API.pseudohsmSignTemplates
func SignRawTransaction(prvkeyStr, templateStr string) (string, error) {
	var tpl txbuilder.Template
	if err := json.Unmarshal([]byte(templateStr), &tpl); err != nil {
		return "", err
	}

	var unused context.Context
	if err := txbuilder.Sign(unused, &tpl, prvkeyStr, pseudohsmSignTemplate); err != nil {
		return "", fmt.Errorf("failed on sign transactio. %v", err)
	}

	resp := signResp{Tx: &tpl, SignComplete: txbuilder.SignProgress(&tpl)}
	result, err := json.Marshal(resp)
	return string(result), err
}
