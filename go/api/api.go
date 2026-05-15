package api

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	_ "log"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/BitBoxSwiss/bitbox02-api-go/api/common"
	"github.com/BitBoxSwiss/bitbox02-api-go/api/firmware"
	"github.com/BitBoxSwiss/bitbox02-api-go/api/firmware/mocks"
	"github.com/BitBoxSwiss/bitbox02-api-go/util/semver"
	"github.com/konstantinullrich/bitbox_flutter/u2fhid"
)

// fixTimezone sets the local timezone on Android. This is a workaround to the bug that on Android,
// time.Local is hard-coded to UTC. See https://github.com/golang/go/issues/20455.
//
// We need the correct timezone to be able to send the `time.Now().Zone()` offset to the BitBox02.
// Without it, the BitBox02 will always display UTC time instead of local time.
//
// This fix is copied from https://github.com/golang/go/issues/20455#issuecomment-342287698.
func fixTimezone() {
	if runtime.GOOS != "android" {
		// Only run the fix on Android.
		return
	}
	out, err := exec.Command("/system/bin/getprop", "persist.sys.timezone").Output()
	if err != nil {
		return
	}
	z, err := time.LoadLocation(strings.TrimSpace(string(out)))
	if err != nil {
		return
	}
	time.Local = z
}

func init() {
	fixTimezone()
}

// GoReadWriteCloserInterface adapts io.ReadWriteCloser's Read method to return the byte read byte slice
// instead of the .Read([]byte) pattern. This is as gomobile bind seems to make a copy of passed
// slices instead of writing directly to it, so the byte slice never makes it back to Go-land.
type GoReadWriteCloserInterface interface {
	Read(n int) ([]byte, error)
	io.Writer
	io.Closer
}

// GoDeviceInfoInterface adapts usb.DeviceInfo's Open method to return the adapted ReadWriteCloser.
type GoDeviceInfoInterface interface {
	IsBluetooth() bool
	VendorID() int
	ProductID() int
	UsagePage() int
	Interface() int
	Serial() string
	Product() string
	Identifier() string
	Open() (GoReadWriteCloserInterface, error)
}

// readWriteCloser implements io.ReadWriteCloser, translating from GoReadWriteCloserInterface. All methods
// are as-is except for Read().
type readWriteCloser struct {
	GoReadWriteCloserInterface
}

// Read implements io.ReadWriteCloser, translating GoReadWriteCloserInterface.Read, which returns a slice
//
//	instead of receiving it as an argument.
func (r readWriteCloser) Read(readBytesOut []byte) (int, error) {
	readBytes, err := r.GoReadWriteCloserInterface.Read(len(readBytesOut))
	if err != nil {
		return 0, err
	}
	copy(readBytesOut, readBytes)
	return len(readBytes), nil
}

// deviceInfo implements usb.DeviceInfo, translating from GoDeviceInfoInterface. All methods are as-is except
// for the Open method.
type deviceInfo struct {
	GoDeviceInfoInterface
}

// Open implements usb.DeviceInfo.
func (d deviceInfo) Open() (io.ReadWriteCloser, error) {
	device, err := d.GoDeviceInfoInterface.Open()
	if err != nil {
		return nil, err
	}
	return readWriteCloser{device}, nil
}

//export GetDevice
func GetDevice(device GoReadWriteCloserInterface) {
	defer recoverPanic("GetDevice")

	const bitboxCMD = 0x80 + 0x40 + 0x01
	comm := u2fhid.NewCommunication(readWriteCloser{device}, bitboxCMD)
	bitbox = firmware.NewDevice(nil, nil, &mocks.Config{}, comm, &mocks.Logger{})
}

// GetDeviceWithInfo is like GetDevice but accepts version and product info for Bluetooth connections.
// version should be like "v9.25.0", product should be like "bb02p-multi" or "bb02p-btconly".
//
//export GetDeviceWithInfo
func GetDeviceWithInfo(device GoReadWriteCloserInterface, versionStr string, productStr string) {
	defer recoverPanic("GetDeviceWithInfo")

	const bitboxCMD = 0x80 + 0x40 + 0x01
	comm := u2fhid.NewCommunication(readWriteCloser{device}, bitboxCMD)

	// Parse version string (e.g., "v9.25.0" -> semver). Fall back to a known
	// default rather than panicking, so a malformed version from the device
	// does not crash the host engine.
	version, err := semver.NewSemVerFromString(strings.TrimPrefix(versionStr, "v"))
	if err != nil {
		fmt.Printf("[GetDeviceWithInfo] invalid version %q, falling back: %v\n", versionStr, err)
		version = semver.NewSemVer(9, 25, 0)
	}

	// Map product string to common.Product
	// bb02p-* are BitBox02 Plus (Nova) products — must use Plus variants
	var product common.Product
	switch productStr {
	case "bb02p-multi", "BitBox02 Nova Multi":
		product = common.ProductBitBox02PlusMulti
	case "bb02p-btconly", "BitBox02 Nova BTC-only":
		product = common.ProductBitBox02PlusBTCOnly
	case "bb02p-bl-multi", "BitBox02 Nova Multi bl":
		product = common.ProductBitBox02PlusMulti
	case "bb02p-bl-btconly", "BitBox02 Nova BTC-only bl":
		product = common.ProductBitBox02PlusBTCOnly
	// Original BitBox02 (non-Nova) products
	case "BitBox02Multi":
		product = common.ProductBitBox02Multi
	case "BitBox02BTCOnly":
		product = common.ProductBitBox02BTCOnly
	default:
		product = common.ProductBitBox02PlusMulti
	}

	bitbox = firmware.NewDevice(version, &product, &mocks.Config{}, comm, &mocks.Logger{})
}

//export GetChannelHash
func GetChannelHash() (hash string) {
	defer recoverPanic("GetChannelHash")

	if bitbox == nil {
		return ""
	}
	hash, _ = bitbox.ChannelHash()
	return hash
}

//export ChannelHashVerify
func ChannelHashVerify(ok bool) {
	defer recoverPanic("ChannelHashVerify")

	if bitbox == nil {
		return
	}
	bitbox.ChannelHashVerify(ok)
}

//export InitDevice
func InitDevice() (success bool) {
	defer recoverPanic("InitDevice")

	if bitbox == nil {
		fmt.Println("[InitDevice] device pointer is nil")
		return false
	}
	err := bitbox.Init()
	if err != nil {
		fmt.Println("[InitDevice] error:", err)
		return false
	}
	return true
}

//export SupportsETH
func SupportsETH(chainId int) (supported bool) {
	defer recoverPanic("SupportsETH")

	if bitbox == nil {
		return false
	}
	return bitbox.SupportsETH(uint64(chainId))
}

//export SupportsLTC
func SupportsLTC() (supported bool) {
	defer recoverPanic("SupportsLTC")

	if bitbox == nil {
		return false
	}
	return bitbox.SupportsLTC()
}

//export SupportsBluetooth
func SupportsBluetooth() (supported bool) {
	defer recoverPanic("SupportsBluetooth")

	if bitbox == nil {
		return false
	}
	return bitbox.SupportsBluetooth()
}

//export SupportsERC20
func SupportsERC20(contractAddress string) (supported bool) {
	defer recoverPanic("SupportsERC20")

	if bitbox == nil {
		return false
	}
	return bitbox.SupportsERC20(contractAddress)
}

//export DeviceInfo
func DeviceInfo() (out firmware.DeviceInfo) {
	defer recoverPanic("DeviceInfo")

	if bitbox == nil {
		return firmware.DeviceInfo{}
	}
	info, err := bitbox.DeviceInfo()
	if err != nil || info == nil {
		fmt.Printf("[DeviceInfo] error: %v\n", err)
		return firmware.DeviceInfo{}
	}
	return *info
}

func hexToUint32Slice(hexStr string) ([]uint32, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("error decoding the hex string: %v", err)
	}

	if len(bytes)%4 != 0 {
		return nil, fmt.Errorf("bytelength has to be divisable by 4 for uint32")
	}

	result := make([]uint32, len(bytes)/4)
	for i := 0; i < len(bytes); i += 4 {
		result[i/4] = binary.BigEndian.Uint32(bytes[i : i+4])
	}

	return result, nil
}
