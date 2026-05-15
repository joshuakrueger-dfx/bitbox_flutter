package api

import (
	"encoding/hex"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/BitBoxSwiss/bitbox02-api-go/api/firmware"
	"github.com/BitBoxSwiss/bitbox02-api-go/api/firmware/messages"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type fakeBitboxDevice struct {
	initErr error

	channelHash         string
	channelHashOk       bool
	channelHashVerified *bool

	deviceInfo         *firmware.DeviceInfo
	deviceInfoErr      error
	rootFingerprint    []byte
	rootFingerprintErr error

	supportsETH       bool
	supportsLTC       bool
	supportsBluetooth bool
	supportedERC20    map[string]bool

	ethPubResult              string
	ethPubErr                 error
	ethSignResult             []byte
	ethSignErr                error
	ethSignEIP1559Result      []byte
	ethSignEIP1559Err         error
	ethSignMessageResult      []byte
	ethSignMessageErr         error
	ethSignTypedMessageResult []byte
	ethSignTypedMessageErr    error

	btcXPubResult       string
	btcXPubErr          error
	btcSignPSBTErr      error
	btcSignMessageSig   []byte
	btcSignMessageErr   error
	panicOnETHSignTyped bool

	calls []string
}

func (f *fakeBitboxDevice) Init() error {
	f.calls = append(f.calls, "Init")
	return f.initErr
}

func (f *fakeBitboxDevice) ChannelHash() (string, bool) {
	f.calls = append(f.calls, "ChannelHash")
	return f.channelHash, f.channelHashOk
}

func (f *fakeBitboxDevice) ChannelHashVerify(ok bool) {
	f.calls = append(f.calls, "ChannelHashVerify")
	f.channelHashVerified = &ok
}

func (f *fakeBitboxDevice) DeviceInfo() (*firmware.DeviceInfo, error) {
	f.calls = append(f.calls, "DeviceInfo")
	return f.deviceInfo, f.deviceInfoErr
}

func (f *fakeBitboxDevice) RootFingerprint() ([]byte, error) {
	f.calls = append(f.calls, "RootFingerprint")
	return f.rootFingerprint, f.rootFingerprintErr
}

func (f *fakeBitboxDevice) SupportsETH(chainID uint64) bool {
	f.calls = append(f.calls, "SupportsETH")
	return f.supportsETH && chainID != 0
}

func (f *fakeBitboxDevice) SupportsLTC() bool {
	f.calls = append(f.calls, "SupportsLTC")
	return f.supportsLTC
}

func (f *fakeBitboxDevice) SupportsBluetooth() bool {
	f.calls = append(f.calls, "SupportsBluetooth")
	return f.supportsBluetooth
}

func (f *fakeBitboxDevice) SupportsERC20(contractAddress string) bool {
	f.calls = append(f.calls, "SupportsERC20")
	return f.supportedERC20[contractAddress]
}

func (f *fakeBitboxDevice) ETHPub(uint64, []uint32, messages.ETHPubRequest_OutputType, bool, []byte) (string, error) {
	f.calls = append(f.calls, "ETHPub")
	return f.ethPubResult, f.ethPubErr
}

func (f *fakeBitboxDevice) ETHSign(uint64, []uint32, uint64, *big.Int, uint64, [20]byte, *big.Int, []byte, messages.ETHAddressCase) ([]byte, error) {
	f.calls = append(f.calls, "ETHSign")
	return f.ethSignResult, f.ethSignErr
}

func (f *fakeBitboxDevice) ETHSignEIP1559(uint64, []uint32, uint64, *big.Int, *big.Int, uint64, [20]byte, *big.Int, []byte, messages.ETHAddressCase) ([]byte, error) {
	f.calls = append(f.calls, "ETHSignEIP1559")
	return f.ethSignEIP1559Result, f.ethSignEIP1559Err
}

func (f *fakeBitboxDevice) ETHSignMessage(uint64, []uint32, []byte) ([]byte, error) {
	f.calls = append(f.calls, "ETHSignMessage")
	return f.ethSignMessageResult, f.ethSignMessageErr
}

func (f *fakeBitboxDevice) ETHSignTypedMessage(uint64, []uint32, []byte) ([]byte, error) {
	f.calls = append(f.calls, "ETHSignTypedMessage")
	if f.panicOnETHSignTyped {
		panic("simulated typed-data panic")
	}
	return f.ethSignTypedMessageResult, f.ethSignTypedMessageErr
}

func (f *fakeBitboxDevice) BTCXPub(messages.BTCCoin, []uint32, messages.BTCPubRequest_XPubType, bool) (string, error) {
	f.calls = append(f.calls, "BTCXPub")
	return f.btcXPubResult, f.btcXPubErr
}

func (f *fakeBitboxDevice) BTCSignPSBT(messages.BTCCoin, *psbt.Packet, *firmware.PSBTSignOptions) error {
	f.calls = append(f.calls, "BTCSignPSBT")
	return f.btcSignPSBTErr
}

func (f *fakeBitboxDevice) BTCSignMessage(messages.BTCCoin, *messages.BTCScriptConfigWithKeypath, []byte) (*firmware.BTCSignMessageResult, error) {
	f.calls = append(f.calls, "BTCSignMessage")
	return &firmware.BTCSignMessageResult{Signature: f.btcSignMessageSig}, f.btcSignMessageErr
}

func withFakeBitbox(t *testing.T, fake *fakeBitboxDevice) {
	t.Helper()
	previous := bitbox
	bitbox = fake
	t.Cleanup(func() {
		bitbox = previous
	})
}

func TestFakeBitboxHarnessSimulatesPairingAndCapabilities(t *testing.T) {
	verified := false
	fake := &fakeBitboxDevice{
		channelHash:         "PAIR-CODE",
		channelHashOk:       true,
		channelHashVerified: &verified,
		deviceInfo:          &firmware.DeviceInfo{Name: "Simulated BitBox"},
		rootFingerprint:     []byte{0x01, 0x02, 0x03, 0x04},
		supportsETH:         true,
		supportsLTC:         true,
		supportsBluetooth:   true,
		supportedERC20:      map[string]bool{"0xToken": true},
	}
	withFakeBitbox(t, fake)

	if !InitDevice() {
		t.Fatal("expected simulated init to succeed")
	}
	if got := GetChannelHash(); got != "PAIR-CODE" {
		t.Fatalf("expected channel hash, got %q", got)
	}
	ChannelHashVerify(true)
	if fake.channelHashVerified == nil || !*fake.channelHashVerified {
		t.Fatal("expected simulated channel hash confirmation")
	}
	if !SupportsETH(1) || !SupportsLTC() || !SupportsBluetooth() || !SupportsERC20("0xToken") {
		t.Fatal("expected simulated capabilities to be exposed")
	}
	if got := DeviceInfo().Name; got != "Simulated BitBox" {
		t.Fatalf("expected simulated device info, got %q", got)
	}
	if got := GetMasterFingerprint(); !reflect.DeepEqual(got, []byte{0x01, 0x02, 0x03, 0x04}) {
		t.Fatalf("expected simulated root fingerprint, got %x", got)
	}
}

func TestFakeBitboxHarnessSimulatesEthereumAndBitcoinOperations(t *testing.T) {
	fake := &fakeBitboxDevice{
		ethPubResult:              "0x0000000000000000000000000000000000000001",
		ethSignResult:             []byte{0x11},
		ethSignEIP1559Result:      []byte{0x22},
		ethSignMessageResult:      []byte{0x33},
		ethSignTypedMessageResult: []byte{0x44},
		btcXPubResult:             "xpub-simulated",
		btcSignMessageSig:         []byte{0x55},
	}
	withFakeBitbox(t, fake)

	keypath := "0000002c0000003c000000000000000000000000"
	if got := ETHGetAddress(1, keypath, int(messages.ETHPubRequest_ADDRESS), false, nil); got != fake.ethPubResult {
		t.Fatalf("expected fake ETH address, got %q", got)
	}
	if got := ETHSignTransaction(1, keypath, 1, "01", 21000, make([]byte, 20), "01", nil, int(messages.ETHAddressCase_ETH_ADDRESS_CASE_LOWER)); !reflect.DeepEqual(got, []byte{0x11}) {
		t.Fatalf("expected fake ETH legacy signature, got %x", got)
	}
	encodedLegacyTx, err := rlp.EncodeToBytes(legacyTxPayload{
		Nonce:    1,
		GasPrice: big.NewInt(1),
		Gas:      21000,
		To:       ptr(common.HexToAddress("0x0000000000000000000000000000000000000001")),
		Value:    big.NewInt(1),
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := ETHSignRPLTx(1, keypath, hex.EncodeToString(encodedLegacyTx), false); !reflect.DeepEqual(got, []byte{0x11}) {
		t.Fatalf("expected fake ETH RLP signature, got %x", got)
	}
	if got := ETHSignEIP1559(1, keypath, 1, "01", "02", 21000, make([]byte, 20), "01", nil, int(messages.ETHAddressCase_ETH_ADDRESS_CASE_LOWER)); !reflect.DeepEqual(got, []byte{0x22}) {
		t.Fatalf("expected fake ETH EIP1559 signature, got %x", got)
	}
	if got := ETHSignMessage(1, keypath, []byte("hello")); !reflect.DeepEqual(got, []byte{0x33}) {
		t.Fatalf("expected fake ETH message signature, got %x", got)
	}
	if got := ETHSignTypedMessage(1, keypath, []byte(`{"types":{},"primaryType":"Mail"}`)); !reflect.DeepEqual(got, []byte{0x44}) {
		t.Fatalf("expected fake ETH typed signature, got %x", got)
	}
	if got := BTCXPub(int(messages.BTCCoin_BTC), keypath, int(messages.BTCPubRequest_XPUB), false); got != "xpub-simulated" {
		t.Fatalf("expected fake xpub, got %q", got)
	}
	const validPSBT = "cHNidP8BAHECAAAAAfbXTun4YYxDroWyzRq3jDsWFVlsZ7HUzxiORY/iR4goAAAAAAD9////AuLCAAAAAAAAFgAUg3w5W0zt3AmxRmgA5Q6wZJUDRhUowwAAAAAAABYAFJjQqUoXDcwUEqfExu9pnaSn5XBct0ElAAABAR+ghgEAAAAAABYAFHn03igII+hp819N2Zlb5LnN8atRAQDfAQAAAAABAZ9EJlMJnXF5bFVrb1eFBYrEev3pg35WpvS3RlELsMMrAQAAAAD9////AqCGAQAAAAAAFgAUefTeKAgj6GnzX03ZmVvkuc3xq1EoRs4JAAAAABYAFKG2PzjYjknaA6lmXFqPaSgHwXX9AkgwRQIhAL0v0r3LisQ9KOlGzMhM/xYqUmrv2a5sORRlkX1fqDC8AiB9XqxSNEdb4mPnp7ylF1cAlbAZ7jMhgIxHUXylTww3bwEhA0AEOM0yYEpexPoKE3vT51uxZ+8hk9sOEfBFKOeo6oDDAAAAACIGAyNQfmAT/YLmZaxxfDwClmVNt2BkFnfQu/i8Uc/hHDUiGBKiwYlUAACAAQAAgAAAAIAAAAAAAAAAAAAAIgIDnxFM7Qr9LvJwQDB9GozdTRIe3MYVuHOqT7dU2EuvHrIYEqLBiVQAAIABAACAAAAAgAEAAAAAAAAAAA=="
	if got := BTCSignPSBT(int(messages.BTCCoin_BTC), validPSBT); got == "" {
		t.Fatal("expected fake PSBT signing to return an encoded PSBT")
	}
	if got := BTCSignMessage(int(messages.BTCCoin_BTC), keypath, []byte("hello")); !reflect.DeepEqual(got, []byte{0x55}) {
		t.Fatalf("expected fake BTC message signature, got %x", got)
	}
}

func ptr[T any](value T) *T {
	return &value
}

func TestFakeBitboxHarnessSimulatesErrorsAndPanicsWithoutCrashing(t *testing.T) {
	fake := &fakeBitboxDevice{
		initErr:                   errors.New("init failed"),
		ethPubErr:                 errors.New("address rejected"),
		ethSignMessageErr:         errors.New("message rejected"),
		btcXPubErr:                errors.New("xpub rejected"),
		rootFingerprintErr:        errors.New("fingerprint rejected"),
		panicOnETHSignTyped:       true,
		ethSignTypedMessageErr:    errors.New("unused"),
		ethSignTypedMessageResult: []byte{0xaa},
	}
	withFakeBitbox(t, fake)

	if InitDevice() {
		t.Fatal("expected simulated init error")
	}
	keypath := "0000002c0000003c000000000000000000000000"
	if got := ETHGetAddress(1, keypath, int(messages.ETHPubRequest_ADDRESS), false, nil); got != "" {
		t.Fatalf("expected empty address on simulated error, got %q", got)
	}
	if got := ETHSignMessage(1, keypath, []byte("hello")); got != nil {
		t.Fatalf("expected nil ETH signature on simulated error, got %x", got)
	}
	if got := BTCXPub(int(messages.BTCCoin_BTC), keypath, int(messages.BTCPubRequest_XPUB), false); got != "" {
		t.Fatalf("expected empty xpub on simulated error, got %q", got)
	}
	if got := GetMasterFingerprint(); len(got) != 0 {
		t.Fatalf("expected empty fingerprint on simulated error, got %x", got)
	}
	if got := ETHSignTypedMessage(1, keypath, []byte(`{}`)); got != nil {
		t.Fatalf("expected panic recovery to return nil signature, got %x", got)
	}
}

func TestExportedAPIsReturnZeroValuesWithoutDeviceInsteadOfCrashing(t *testing.T) {
	previous := bitbox
	bitbox = nil
	t.Cleanup(func() {
		bitbox = previous
	})

	keypath := "0000002c0000003c000000000000000000000000"
	if InitDevice() {
		t.Fatal("expected InitDevice to fail without a device")
	}
	if got := GetChannelHash(); got != "" {
		t.Fatalf("expected empty channel hash without device, got %q", got)
	}
	ChannelHashVerify(false)
	if SupportsETH(1) || SupportsLTC() || SupportsBluetooth() || SupportsERC20("0xToken") {
		t.Fatal("expected no capabilities without device")
	}
	if got := DeviceInfo(); got.Name != "" {
		t.Fatalf("expected zero device info without device, got %+v", got)
	}
	if got := ETHGetAddress(1, keypath, int(messages.ETHPubRequest_ADDRESS), false, nil); got != "" {
		t.Fatalf("expected empty ETH address without device, got %q", got)
	}
	if got := ETHSignMessage(1, keypath, []byte("hello")); got != nil {
		t.Fatalf("expected nil ETH message signature without device, got %x", got)
	}
	if got := ETHSignTypedMessage(1, keypath, []byte(`{}`)); got != nil {
		t.Fatalf("expected nil ETH typed signature without device, got %x", got)
	}
	if got := BTCXPub(int(messages.BTCCoin_BTC), keypath, int(messages.BTCPubRequest_XPUB), false); got != "" {
		t.Fatalf("expected empty xpub without device, got %q", got)
	}
	if got := BTCSignMessage(int(messages.BTCCoin_BTC), keypath, []byte("hello")); got != nil {
		t.Fatalf("expected nil BTC message signature without device, got %x", got)
	}
	if got := GetMasterFingerprint(); len(got) != 0 {
		t.Fatalf("expected empty fingerprint without device, got %x", got)
	}
}
