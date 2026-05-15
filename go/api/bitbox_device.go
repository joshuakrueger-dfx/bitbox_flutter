package api

import (
	"math/big"

	"github.com/BitBoxSwiss/bitbox02-api-go/api/firmware"
	"github.com/BitBoxSwiss/bitbox02-api-go/api/firmware/messages"
	"github.com/btcsuite/btcd/btcutil/psbt"
)

type bitboxDevice interface {
	Init() error
	ChannelHash() (string, bool)
	ChannelHashVerify(ok bool)
	DeviceInfo() (*firmware.DeviceInfo, error)
	RootFingerprint() ([]byte, error)
	SupportsETH(chainID uint64) bool
	SupportsLTC() bool
	SupportsBluetooth() bool
	SupportsERC20(contractAddress string) bool

	ETHPub(
		chainID uint64,
		keypath []uint32,
		outputType messages.ETHPubRequest_OutputType,
		display bool,
		contractAddress []byte,
	) (string, error)
	ETHSign(
		chainID uint64,
		keypath []uint32,
		nonce uint64,
		gasPrice *big.Int,
		gasLimit uint64,
		recipient [20]byte,
		value *big.Int,
		data []byte,
		recipientAddressCase messages.ETHAddressCase,
	) ([]byte, error)
	ETHSignEIP1559(
		chainID uint64,
		keypath []uint32,
		nonce uint64,
		maxPriorityFeePerGas *big.Int,
		maxFeePerGas *big.Int,
		gasLimit uint64,
		recipient [20]byte,
		value *big.Int,
		data []byte,
		recipientAddressCase messages.ETHAddressCase,
	) ([]byte, error)
	ETHSignMessage(chainID uint64, keypath []uint32, msg []byte) ([]byte, error)
	ETHSignTypedMessage(chainID uint64, keypath []uint32, jsonMsg []byte) ([]byte, error)

	BTCXPub(
		coin messages.BTCCoin,
		keypath []uint32,
		xpubType messages.BTCPubRequest_XPubType,
		display bool,
	) (string, error)
	BTCSignPSBT(coin messages.BTCCoin, psbt *psbt.Packet, options *firmware.PSBTSignOptions) error
	BTCSignMessage(
		coin messages.BTCCoin,
		scriptConfig *messages.BTCScriptConfigWithKeypath,
		message []byte,
	) (*firmware.BTCSignMessageResult, error)
}

var bitbox bitboxDevice
