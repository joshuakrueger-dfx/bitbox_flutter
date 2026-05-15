import 'dart:async';
import 'dart:typed_data';

import 'package:bitbox_flutter/bitbox_manager.dart';
import 'package:bitbox_flutter/testing.dart';
import 'package:bitbox_flutter/usb/bitbox_usb_platform_interface.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  late BitboxUsbPlatform previousPlatform;

  setUp(() {
    previousPlatform = BitboxUsbPlatform.instance;
  });

  tearDown(() {
    BitboxUsbPlatform.instance = previousPlatform;
  });

  test('simulates pairing and records the connection flow', () async {
    final platform = installSimulatedBitboxPlatform(
      channelHash: 'hash-visible-to-user',
    );
    final manager = BitboxManager();

    final devices = await manager.devices;
    await manager.connect(devices.single);

    expect(await manager.initBitBox(), isTrue);
    expect(await manager.getChannelHash(), 'hash-visible-to-user');
    expect(platform.channelHashVerified, isFalse);
    expect(await manager.channelHashVerify(), isTrue);
    expect(platform.channelHashVerified, isTrue);

    expect(
        platform.calls.map((call) => call.method),
        containsAllInOrder([
          SimulatedBitboxMethod.getDevices,
          SimulatedBitboxMethod.requestPermission,
          SimulatedBitboxMethod.open,
          SimulatedBitboxMethod.initBitBox,
          SimulatedBitboxMethod.getChannelHash,
          SimulatedBitboxMethod.channelHashVerify,
        ]));
  });

  test('covers generic BTC and ETH operations without hardware', () async {
    final platform = installSimulatedBitboxPlatform(
      masterFingerprint: Uint8List.fromList(<int>[1, 2, 3, 4]),
      btcXPub: 'xpub-test',
      btcPsbt: 'signed-psbt',
      btcMessageSignature: Uint8List.fromList(<int>[5, 6]),
      ethAddress: '0x2222222222222222222222222222222222222222',
      ethTransactionSignature: Uint8List.fromList(<int>[7, 8]),
      ethEip1559Signature: Uint8List.fromList(<int>[9, 10]),
      ethRlpSignature: Uint8List.fromList(<int>[11, 12]),
      ethMessageSignature: Uint8List.fromList(<int>[13, 14]),
      ethTypedMessageSignature: Uint8List.fromList(<int>[15, 16]),
    );
    final manager = BitboxManager();
    await manager.connect((await manager.devices).single);

    expect(await manager.supportsETH(1), isTrue);
    expect(await manager.supportsERC20('0xToken'), isTrue);
    expect(await manager.supportsLTC(), isTrue);
    expect(await manager.getMasterFingerprint(), <int>[1, 2, 3, 4]);
    expect(await manager.getBTCXPub(0, "m/84'/0'/0'"), 'xpub-test');
    expect(await manager.signBTCPsbt(0, 'psbt'), 'signed-psbt');
    expect(
      await manager.signBTCMessage(
        0,
        "m/84'/0'/0'/0/0",
        Uint8List.fromList(<int>[1]),
      ),
      <int>[5, 6],
    );
    expect(
      await manager.getETHAddress(1, "m/44'/60'/0'/0/0"),
      '0x2222222222222222222222222222222222222222',
    );
    expect(
      await manager.signETHTransaction(
        1,
        "m/44'/60'/0'/0/0",
        1,
        BigInt.from(2),
        21000,
        Uint8List(20),
        BigInt.from(3),
        Uint8List(0),
        0,
      ),
      <int>[7, 8],
    );
    expect(
      await manager.signETHTransactionEIP1559(
        1,
        "m/44'/60'/0'/0/0",
        1,
        BigInt.from(2),
        BigInt.from(3),
        21000,
        Uint8List(20),
        BigInt.from(4),
        Uint8List(0),
        0,
      ),
      <int>[9, 10],
    );
    expect(
      await manager.signETHRLPTransaction(1, "m/44'/60'/0'/0/0", '0x01', true),
      <int>[11, 12],
    );
    expect(
      await manager.signETHMessage(
        1,
        "m/44'/60'/0'/0/0",
        Uint8List.fromList(<int>[1]),
      ),
      <int>[13, 14],
    );
    expect(
      await manager.signETHTypedMessage(
        1,
        "m/44'/60'/0'/0/0",
        Uint8List.fromList(<int>[123, 125]),
      ),
      <int>[15, 16],
    );

    expect(platform.count(SimulatedBitboxMethod.getETHAddress), 1);
    final call =
        platform.callsFor(SimulatedBitboxMethod.signETHTransaction).single;
    expect(call.argument<String>('gasPrice'), '2');
    expect(call.argument<String>('value'), '3');
  });

  test('simulates delays and hardware errors deterministically', () async {
    final releaseSigning = Completer<void>();
    final platform = installSimulatedBitboxPlatform(
      behaviors: <String, SimulatedBitboxBehavior>{
        SimulatedBitboxMethod.signETHMessage: (_) async {
          await releaseSigning.future;
          return Uint8List.fromList(<int>[42]);
        },
      },
    );
    final manager = BitboxManager();
    await manager.connect((await manager.devices).single);

    final pending = manager.signETHMessage(
      1,
      "m/44'/60'/0'/0/0",
      Uint8List.fromList(<int>[1]),
    );
    await Future<void>.delayed(Duration.zero);
    expect(platform.count(SimulatedBitboxMethod.signETHMessage), 1);

    releaseSigning.complete();
    expect(await pending, <int>[42]);

    platform.throwOn(
      SimulatedBitboxMethod.signETHMessage,
      StateError('signature aborted'),
    );
    expect(
      manager.signETHMessage(
        1,
        "m/44'/60'/0'/0/0",
        Uint8List.fromList(<int>[1]),
      ),
      throwsA(isA<StateError>()),
    );
  });

  test('simulates rejected pairing, unsupported capabilities, and close',
      () async {
    final platform = installSimulatedBitboxPlatform(
      channelHashVerifyResult: false,
      supportsETHResult: false,
      supportsERC20Result: false,
      supportsLTCResult: false,
    );
    final manager = BitboxManager();
    await manager.connect((await manager.devices).single);

    expect(await manager.channelHashVerify(), isFalse);
    expect(platform.channelHashVerified, isFalse);
    expect(await manager.supportsETH(1), isFalse);
    expect(await manager.supportsERC20('0xToken'), isFalse);
    expect(await manager.supportsLTC(), isFalse);

    await manager.disconnect();
    expect(platform.isOpen, isFalse);
    expect(
      manager.getMasterFingerprint(),
      throwsA(isA<SimulatedBitboxStateException>()),
    );
  });

  test('guards against signing before the simulated device is opened',
      () async {
    installSimulatedBitboxPlatform();
    final manager = BitboxManager();

    expect(
      manager.signETHMessage(
        1,
        "m/44'/60'/0'/0/0",
        Uint8List.fromList(<int>[1]),
      ),
      throwsA(isA<SimulatedBitboxStateException>()),
    );
  });
}
