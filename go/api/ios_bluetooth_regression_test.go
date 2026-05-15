package api

import (
	"os"
	"strings"
	"testing"
)

func TestIOSBluetoothDoesNotDeduplicateReadPackets(t *testing.T) {
	contentBytes, err := os.ReadFile("../../ios/Classes/Bluetooth.swift")
	if err != nil {
		t.Fatal(err)
	}
	content := string(contentBytes)

	forbiddenPatterns := []string{
		"BLEPacketDeduper",
		"packetDeduper",
		"seenPackets",
		"skipping duplicate packet",
		"shouldAccept(",
	}
	for _, pattern := range forbiddenPatterns {
		if strings.Contains(content, pattern) {
			t.Fatalf("Bluetooth.swift must not contain BLE packet dedup pattern %q", pattern)
		}
	}

	if !strings.Contains(content, "let waitResult = ctx.semaphore.wait(timeout: .now() + 60)") {
		t.Fatal("Bluetooth.swift must keep the 60s BLE read timeout for long BitBox confirmations")
	}
	if strings.Contains(content, "read timed out after 10s") || strings.Contains(content, ".now() + 10") {
		t.Fatal("Bluetooth.swift must not regress to the old 10s BLE read timeout")
	}
}
