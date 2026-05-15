# Testing

## Fast PR gate

Run the hardware-wallet tests that do not need a physical BitBox:

```sh
flutter analyze --no-fatal-infos
flutter test
cd go
go test ./api
```

The Go API tests include a generic fake BitBox device. It is not RealUnit-specific:
it can simulate pairing, channel hashes, confirmations, capabilities, ETH address
lookup, ETH signing, BTC xpubs, BTC message signing, device errors, and recovered
panics without connecting hardware.

## Official BitBox simulator

`github.com/BitBoxSwiss/bitbox02-api-go` also ships official `TestSimulator*`
integration tests. Their README documents:

```sh
go test -v -run TestSimulator ./...
SIMULATOR=/path/to/simulator go test -v -run TestSimulator ./...
```

The published simulator binaries referenced by the dependency are Linux amd64
binaries, so they are best suited for Linux CI or a Linux development machine.
The fast fake-device tests in this plugin remain the default local and PR gate.

## Regression coverage

The tests explicitly guard against these hardware-wallet regressions:

- gomobile-exported API functions without `recoverPanic`
- iOS BLE packet deduplication being reintroduced
- iOS BLE read timeout regressing from 60 seconds to 10 seconds
- Pairing/channel-hash behavior not being simulatable without hardware
- ETH/BTC success, error, and panic flows not being simulatable without hardware
