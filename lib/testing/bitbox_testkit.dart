import 'dart:async';
import 'dart:typed_data';

import 'package:bitbox_flutter/usb/bitbox_device.dart';
import 'package:bitbox_flutter/usb/bitbox_usb_platform_interface.dart';

typedef SimulatedBitboxBehavior = FutureOr<Object?> Function(
  SimulatedBitboxCall call,
);

final class SimulatedBitboxMethod {
  static const getDevices = 'getDevices';
  static const startScan = 'startScan';
  static const requestPermission = 'requestPermission';
  static const open = 'open';
  static const initBitBox = 'initBitBox';
  static const getMasterFingerprint = 'getMasterFingerprint';
  static const getChannelHash = 'getChannelHash';
  static const channelHashVerify = 'channelHashVerify';
  static const supportsETH = 'supportsETH';
  static const supportsERC20 = 'supportsERC20';
  static const supportsLTC = 'supportsLTC';
  static const getBTCXPub = 'getBTCXPub';
  static const signBTCPsbt = 'signBTCPsbt';
  static const signBTCMessage = 'signBTCMessage';
  static const getETHAddress = 'getETHAddress';
  static const signETHRPLTransaction = 'signETHRPLTransaction';
  static const signETHTransaction = 'signETHTransaction';
  static const signETHTransactionEIP1559 = 'signETHTransactionEIP1559';
  static const signETHMessage = 'signETHMessage';
  static const signETHTypedMessage = 'signETHTypedMessage';
  static const close = 'close';

  const SimulatedBitboxMethod._();
}

final class SimulatedBitboxCall {
  const SimulatedBitboxCall(this.method, this.arguments);

  final String method;
  final Map<String, Object?> arguments;

  T argument<T>(String name) => arguments[name] as T;

  @override
  String toString() => '$method($arguments)';
}

final class SimulatedBitboxStateException implements Exception {
  const SimulatedBitboxStateException(this.message);

  final String message;

  @override
  String toString() => 'SimulatedBitboxStateException: $message';
}

class SimulatedBitboxPlatform extends BitboxUsbPlatform {
  SimulatedBitboxPlatform({
    List<BitboxDevice>? devices,
    this.defaultDelay = Duration.zero,
    this.requireOpen = true,
    bool channelHashVerified = false,
    this.startScanResult = true,
    this.permissionResult = true,
    this.openResult = true,
    this.initResult = true,
    this.channelHashVerifyResult = true,
    this.supportsETHResult = true,
    this.supportsERC20Result = true,
    this.supportsLTCResult = true,
    String? channelHash,
    Uint8List? masterFingerprint,
    String? btcXPub,
    String? btcPsbt,
    Uint8List? btcMessageSignature,
    String? ethAddress,
    Uint8List? ethTransactionSignature,
    Uint8List? ethEip1559Signature,
    Uint8List? ethRlpSignature,
    Uint8List? ethMessageSignature,
    Uint8List? ethTypedMessageSignature,
    Map<String, Duration>? delays,
    Map<String, Object>? errors,
    Map<String, SimulatedBitboxBehavior>? behaviors,
  })  : _devices = List<BitboxDevice>.of(
          devices ?? <BitboxDevice>[_defaultDevice],
        ),
        channelHash = channelHash ?? _defaultChannelHash,
        masterFingerprint = _copy(masterFingerprint ?? _defaultFingerprint),
        btcXPub = btcXPub ?? _defaultBtcXPub,
        btcPsbt = btcPsbt ?? _defaultSignedPsbt,
        btcMessageSignature = _copy(
          btcMessageSignature ?? _defaultBtcSignature,
        ),
        ethAddress = ethAddress ?? _defaultEthAddress,
        ethTransactionSignature = _copy(
          ethTransactionSignature ?? _defaultEthSignature,
        ),
        ethEip1559Signature = _copy(
          ethEip1559Signature ?? _defaultEthEip1559Signature,
        ),
        ethRlpSignature = _copy(ethRlpSignature ?? _defaultEthRlpSignature),
        ethMessageSignature = _copy(
          ethMessageSignature ?? _defaultEthMessageSignature,
        ),
        ethTypedMessageSignature = _copy(
          ethTypedMessageSignature ?? _defaultEthTypedMessageSignature,
        ),
        _delays = Map<String, Duration>.of(delays ?? const {}),
        _errors = Map<String, Object>.of(errors ?? const {}),
        _behaviors = Map<String, SimulatedBitboxBehavior>.of(
          behaviors ?? const {},
        ),
        _channelHashVerified = channelHashVerified;

  static final BitboxDevice _defaultDevice = BitboxDevice(
    identifier: 'simulated-bitbox-02',
    vendorId: 0x03eb,
    productId: 0x2403,
    productName: 'BitBox02 Simulator',
    deviceId: 1,
    deviceName: 'BitBox02 Simulator',
    manufacturerName: 'Shift Crypto',
    configurationCount: 1,
  );

  static const _defaultChannelHash = 'simulated-channel-hash';
  static const _defaultBtcXPub = 'xpub-simulated';
  static const _defaultSignedPsbt = 'psbt-signed-by-simulator';
  static const _defaultEthAddress =
      '0x1111111111111111111111111111111111111111';

  static final _defaultFingerprint = Uint8List.fromList(
    <int>[0x12, 0x34, 0x56, 0x78],
  );
  static final _defaultBtcSignature = Uint8List.fromList(
    <int>[0x30, 0x44, 0x02, 0x20],
  );
  static final _defaultEthSignature = Uint8List.fromList(<int>[0x01, 0x02]);
  static final _defaultEthEip1559Signature = Uint8List.fromList(
    <int>[0x03, 0x04],
  );
  static final _defaultEthRlpSignature = Uint8List.fromList(<int>[0x05, 0x06]);
  static final _defaultEthMessageSignature = Uint8List.fromList(
    <int>[0x07, 0x08],
  );
  static final _defaultEthTypedMessageSignature = Uint8List.fromList(
    <int>[0x09, 0x0a],
  );

  final Duration defaultDelay;
  final bool requireOpen;
  final bool startScanResult;
  final bool permissionResult;
  final bool openResult;
  final bool initResult;
  final bool channelHashVerifyResult;
  final bool supportsETHResult;
  final bool supportsERC20Result;
  final bool supportsLTCResult;
  final String channelHash;
  final Uint8List masterFingerprint;
  final String btcXPub;
  final String btcPsbt;
  final Uint8List btcMessageSignature;
  final String ethAddress;
  final Uint8List ethTransactionSignature;
  final Uint8List ethEip1559Signature;
  final Uint8List ethRlpSignature;
  final Uint8List ethMessageSignature;
  final Uint8List ethTypedMessageSignature;

  final List<SimulatedBitboxCall> calls = <SimulatedBitboxCall>[];
  final Map<String, Duration> _delays;
  final Map<String, Object> _errors;
  final Map<String, SimulatedBitboxBehavior> _behaviors;
  final List<BitboxDevice> _devices;

  bool _isOpen = false;
  bool _channelHashVerified;

  static SimulatedBitboxPlatform install({
    List<BitboxDevice>? devices,
    Duration defaultDelay = Duration.zero,
    bool requireOpen = true,
    bool channelHashVerified = false,
    bool startScanResult = true,
    bool permissionResult = true,
    bool openResult = true,
    bool initResult = true,
    bool channelHashVerifyResult = true,
    bool supportsETHResult = true,
    bool supportsERC20Result = true,
    bool supportsLTCResult = true,
    String? channelHash,
    Uint8List? masterFingerprint,
    String? btcXPub,
    String? btcPsbt,
    Uint8List? btcMessageSignature,
    String? ethAddress,
    Uint8List? ethTransactionSignature,
    Uint8List? ethEip1559Signature,
    Uint8List? ethRlpSignature,
    Uint8List? ethMessageSignature,
    Uint8List? ethTypedMessageSignature,
    Map<String, Duration>? delays,
    Map<String, Object>? errors,
    Map<String, SimulatedBitboxBehavior>? behaviors,
  }) {
    final platform = SimulatedBitboxPlatform(
      devices: devices,
      defaultDelay: defaultDelay,
      requireOpen: requireOpen,
      channelHashVerified: channelHashVerified,
      startScanResult: startScanResult,
      permissionResult: permissionResult,
      openResult: openResult,
      initResult: initResult,
      channelHashVerifyResult: channelHashVerifyResult,
      supportsETHResult: supportsETHResult,
      supportsERC20Result: supportsERC20Result,
      supportsLTCResult: supportsLTCResult,
      channelHash: channelHash,
      masterFingerprint: masterFingerprint,
      btcXPub: btcXPub,
      btcPsbt: btcPsbt,
      btcMessageSignature: btcMessageSignature,
      ethAddress: ethAddress,
      ethTransactionSignature: ethTransactionSignature,
      ethEip1559Signature: ethEip1559Signature,
      ethRlpSignature: ethRlpSignature,
      ethMessageSignature: ethMessageSignature,
      ethTypedMessageSignature: ethTypedMessageSignature,
      delays: delays,
      errors: errors,
      behaviors: behaviors,
    );
    BitboxUsbPlatform.instance = platform;
    return platform;
  }

  bool get isOpen => _isOpen;

  bool get channelHashVerified => _channelHashVerified;

  List<BitboxDevice> get devices => List<BitboxDevice>.unmodifiable(_devices);

  void setDelay(String method, Duration delay) {
    _delays[method] = delay;
  }

  void throwOn(String method, Object error) {
    _errors[method] = error;
  }

  void clearError(String method) {
    _errors.remove(method);
  }

  void when(String method, SimulatedBitboxBehavior behavior) {
    _behaviors[method] = behavior;
  }

  int count(String method) =>
      calls.where((call) => call.method == method).length;

  List<SimulatedBitboxCall> callsFor(String method) =>
      calls.where((call) => call.method == method).toList(growable: false);

  @override
  Future<List<BitboxDevice>> getDevices() => _run(
        SimulatedBitboxMethod.getDevices,
        const <String, Object?>{},
        List<BitboxDevice>.of(_devices),
        needsOpen: false,
      );

  @override
  Future<bool> startScan() => _run(
        SimulatedBitboxMethod.startScan,
        const <String, Object?>{},
        startScanResult,
        needsOpen: false,
      );

  @override
  Future<bool> requestPermission(BitboxDevice usbDevice) => _run(
        SimulatedBitboxMethod.requestPermission,
        <String, Object?>{'device': usbDevice},
        permissionResult,
        needsOpen: false,
      );

  @override
  Future<bool> open(BitboxDevice usbDevice) async {
    final result = await _run(
      SimulatedBitboxMethod.open,
      <String, Object?>{'device': usbDevice},
      openResult,
      needsOpen: false,
    );
    _isOpen = result;
    return result;
  }

  @override
  Future<bool> initBitBox() => _run(
        SimulatedBitboxMethod.initBitBox,
        const <String, Object?>{},
        initResult,
      );

  @override
  Future<Uint8List> getMasterFingerprint() => _run(
        SimulatedBitboxMethod.getMasterFingerprint,
        const <String, Object?>{},
        _copy(masterFingerprint),
      );

  @override
  Future<String> getChannelHash() => _run(
        SimulatedBitboxMethod.getChannelHash,
        const <String, Object?>{},
        channelHash,
      );

  @override
  Future<bool> channelHashVerify() async {
    final result = await _run(
      SimulatedBitboxMethod.channelHashVerify,
      const <String, Object?>{},
      channelHashVerifyResult,
    );
    _channelHashVerified = result;
    return result;
  }

  @override
  Future<bool> supportsETH(int chainId) => _run(
        SimulatedBitboxMethod.supportsETH,
        <String, Object?>{'chainId': chainId},
        supportsETHResult,
      );

  @override
  Future<bool> supportsERC20(String contractAddress) => _run(
        SimulatedBitboxMethod.supportsERC20,
        <String, Object?>{'contractAddress': contractAddress},
        supportsERC20Result,
      );

  @override
  Future<bool> supportsLTC() => _run(
        SimulatedBitboxMethod.supportsLTC,
        const <String, Object?>{},
        supportsLTCResult,
      );

  @override
  Future<String> getBTCXPub(
    int coinType,
    Uint8List keypath,
    int addressType,
    bool display,
  ) =>
      _run(
        SimulatedBitboxMethod.getBTCXPub,
        <String, Object?>{
          'coinType': coinType,
          'keypath': _copy(keypath),
          'addressType': addressType,
          'display': display,
        },
        btcXPub,
      );

  @override
  Future<String> signBTCPsbt(int coinType, String psbt) => _run(
        SimulatedBitboxMethod.signBTCPsbt,
        <String, Object?>{'coinType': coinType, 'psbt': psbt},
        btcPsbt,
      );

  @override
  Future<Uint8List> signBTCMessage(
    int chainId,
    Uint8List keypath,
    Uint8List message,
  ) =>
      _run(
        SimulatedBitboxMethod.signBTCMessage,
        <String, Object?>{
          'chainId': chainId,
          'keypath': _copy(keypath),
          'message': _copy(message),
        },
        _copy(btcMessageSignature),
      );

  @override
  Future<String> getETHAddress(
    int chainId,
    Uint8List keypath,
    int outputType,
    bool display,
  ) =>
      _run(
        SimulatedBitboxMethod.getETHAddress,
        <String, Object?>{
          'chainId': chainId,
          'keypath': _copy(keypath),
          'outputType': outputType,
          'display': display,
        },
        ethAddress,
      );

  @override
  Future<Uint8List> signETHRPLTransaction(
    int chainId,
    Uint8List keypath,
    String transactionData,
    bool isEIP1559,
  ) =>
      _run(
        SimulatedBitboxMethod.signETHRPLTransaction,
        <String, Object?>{
          'chainId': chainId,
          'keypath': _copy(keypath),
          'transactionData': transactionData,
          'isEIP1559': isEIP1559,
        },
        _copy(ethRlpSignature),
      );

  @override
  Future<Uint8List> signETHTransaction(
    int chainId,
    Uint8List keypath,
    int nonce,
    String gasPrice,
    int gasLimit,
    Uint8List recipient,
    String value,
    Uint8List data,
    int recipientAddressCase,
  ) =>
      _run(
        SimulatedBitboxMethod.signETHTransaction,
        <String, Object?>{
          'chainId': chainId,
          'keypath': _copy(keypath),
          'nonce': nonce,
          'gasPrice': gasPrice,
          'gasLimit': gasLimit,
          'recipient': _copy(recipient),
          'value': value,
          'data': _copy(data),
          'recipientAddressCase': recipientAddressCase,
        },
        _copy(ethTransactionSignature),
      );

  @override
  Future<Uint8List> signETHTransactionEIP1559(
    int chainId,
    Uint8List keypath,
    int nonce,
    String maxPriorityFeePerGas,
    String maxFeePerGas,
    int gasLimit,
    Uint8List recipient,
    String value,
    Uint8List data,
    int recipientAddressCase,
  ) =>
      _run(
        SimulatedBitboxMethod.signETHTransactionEIP1559,
        <String, Object?>{
          'chainId': chainId,
          'keypath': _copy(keypath),
          'nonce': nonce,
          'maxPriorityFeePerGas': maxPriorityFeePerGas,
          'maxFeePerGas': maxFeePerGas,
          'gasLimit': gasLimit,
          'recipient': _copy(recipient),
          'value': value,
          'data': _copy(data),
          'recipientAddressCase': recipientAddressCase,
        },
        _copy(ethEip1559Signature),
      );

  @override
  Future<Uint8List> signETHMessage(
    int chainId,
    Uint8List keypath,
    Uint8List message,
  ) =>
      _run(
        SimulatedBitboxMethod.signETHMessage,
        <String, Object?>{
          'chainId': chainId,
          'keypath': _copy(keypath),
          'message': _copy(message),
        },
        _copy(ethMessageSignature),
      );

  @override
  Future<Uint8List> signETHTypedMessage(
    int chainId,
    Uint8List keypath,
    Uint8List jsonMessage,
  ) =>
      _run(
        SimulatedBitboxMethod.signETHTypedMessage,
        <String, Object?>{
          'chainId': chainId,
          'keypath': _copy(keypath),
          'jsonMessage': _copy(jsonMessage),
        },
        _copy(ethTypedMessageSignature),
      );

  @override
  Future<bool> close() async {
    final result = await _run(
      SimulatedBitboxMethod.close,
      const <String, Object?>{},
      true,
      needsOpen: false,
    );
    if (result) _isOpen = false;
    return result;
  }

  Future<T> _run<T>(
    String method,
    Map<String, Object?> arguments,
    T fallback, {
    bool needsOpen = true,
  }) async {
    final call =
        SimulatedBitboxCall(method, Map<String, Object?>.of(arguments));
    calls.add(call);

    final delay = _delays[method] ?? defaultDelay;
    if (delay > Duration.zero) await Future<void>.delayed(delay);

    final error = _errors[method];
    if (error != null) throw error;

    if (requireOpen && needsOpen && !_isOpen) {
      throw const SimulatedBitboxStateException('BitBox is not open');
    }

    final behavior = _behaviors[method];
    if (behavior == null) return fallback;

    final result = await behavior(call);
    return result as T;
  }

  static Uint8List _copy(Uint8List bytes) => Uint8List.fromList(bytes);
}

SimulatedBitboxPlatform installSimulatedBitboxPlatform({
  List<BitboxDevice>? devices,
  Duration defaultDelay = Duration.zero,
  bool requireOpen = true,
  bool channelHashVerified = false,
  bool startScanResult = true,
  bool permissionResult = true,
  bool openResult = true,
  bool initResult = true,
  bool channelHashVerifyResult = true,
  bool supportsETHResult = true,
  bool supportsERC20Result = true,
  bool supportsLTCResult = true,
  String? channelHash,
  Uint8List? masterFingerprint,
  String? btcXPub,
  String? btcPsbt,
  Uint8List? btcMessageSignature,
  String? ethAddress,
  Uint8List? ethTransactionSignature,
  Uint8List? ethEip1559Signature,
  Uint8List? ethRlpSignature,
  Uint8List? ethMessageSignature,
  Uint8List? ethTypedMessageSignature,
  Map<String, Duration>? delays,
  Map<String, Object>? errors,
  Map<String, SimulatedBitboxBehavior>? behaviors,
}) =>
    SimulatedBitboxPlatform.install(
      devices: devices,
      defaultDelay: defaultDelay,
      requireOpen: requireOpen,
      channelHashVerified: channelHashVerified,
      startScanResult: startScanResult,
      permissionResult: permissionResult,
      openResult: openResult,
      initResult: initResult,
      channelHashVerifyResult: channelHashVerifyResult,
      supportsETHResult: supportsETHResult,
      supportsERC20Result: supportsERC20Result,
      supportsLTCResult: supportsLTCResult,
      channelHash: channelHash,
      masterFingerprint: masterFingerprint,
      btcXPub: btcXPub,
      btcPsbt: btcPsbt,
      btcMessageSignature: btcMessageSignature,
      ethAddress: ethAddress,
      ethTransactionSignature: ethTransactionSignature,
      ethEip1559Signature: ethEip1559Signature,
      ethRlpSignature: ethRlpSignature,
      ethMessageSignature: ethMessageSignature,
      ethTypedMessageSignature: ethTypedMessageSignature,
      delays: delays,
      errors: errors,
      behaviors: behaviors,
    );
