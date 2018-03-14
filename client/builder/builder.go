package builder

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/wire"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	cmn "github.com/tendermint/tmlibs/common"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Broadcast the transaction bytes to Tendermint
func BroadcastTx(tx []byte) (*ctypes.ResultBroadcastTxCommit, error) {

	node, err := client.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxCommit(tx)
	if err != nil {
		return res, err
	}

	if res.CheckTx.Code != uint32(0) {
		return res, errors.Errorf("CheckTx failed: (%d) %s",
			res.CheckTx.Code,
			res.CheckTx.Log)
	}
	if res.DeliverTx.Code != uint32(0) {
		return res, errors.Errorf("DeliverTx failed: (%d) %s",
			res.DeliverTx.Code,
			res.DeliverTx.Log)
	}
	return res, err
}

// Query from Tendermint with the provided key and storename
func Query(key cmn.HexBytes, storeName string) (res []byte, err error) {

	path := fmt.Sprintf("/%s/key", storeName)
	node, err := client.GetNode()
	if err != nil {
		return res, err
	}

	opts := rpcclient.ABCIQueryOptions{
		Height:  viper.GetInt64(client.FlagHeight),
		Trusted: viper.GetBool(client.FlagTrustNode),
	}
	result, err := node.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}
	resp := result.Response
	if resp.Code != uint32(0) {
		return res, errors.Errorf("Query failed: (%d) %s", resp.Code, resp.Log)
	}
	return resp.Value, nil
}

// Get the from address from the name flag
func GetFromAddress(accountName string) (from sdk.Address, err error) {

	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	info, err := keybase.Get(accountName)
	if err != nil {
		return nil, errors.Errorf("No key for: %s", accountName)
	}

	return info.PubKey.Address(), nil
}

// sign and build the transaction from the msg
func SignAndBuild(name string, passphrase string, msg sdk.Msg, cdc *wire.Codec) ([]byte, error) {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	// sign and build
	bz := msg.GetSignBytes()
	sig, pubkey, err := keybase.Sign(name, passphrase, bz)
	if err != nil {
		return nil, err
	}
	sigs := []sdk.StdSignature{{
		PubKey:    pubkey,
		Signature: sig,
		Sequence:  viper.GetInt64(client.FlagSequence),
	}}

	// marshal bytes
	tx := sdk.NewStdTx(msg, sigs)

	return cdc.MarshalBinary(tx)
}

// sign and build the transaction from the msg
func SignBuildBroadcast(name string, passphrase string, msg sdk.Msg, cdc *wire.Codec) (*ctypes.ResultBroadcastTxCommit, error) {
	txBytes, err := SignAndBuild(name, passphrase, msg, cdc)
	if err != nil {
		return nil, err
	}

	return BroadcastTx(txBytes)
}
