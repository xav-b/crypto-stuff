"""Bitcoin transaction from scratch.

Reference: https://en.bitcoin.it/wiki/Transaction

"""


from __future__ import annotations # PEP 563: Postponed Evaluation of Annotations
from dataclasses import dataclass
import random
from typing import List, Union

from bitcoin.hashes import sha256


class Opcode:
    OP_DUP = 118
    OP_HASH160 = 169
    OP_EQUALVERIFY = 136
    OP_CHECKSIG = 172


def encode_varint(i):
    """Encode a (possibly but rarely large) integer into bytes with a super simple compression scheme."""
    if i < 0xfd:
        return bytes([i])
    elif i < 0x10000:
        return b'\xfd' + i.to_bytes(2, 'little')
    elif i < 0x100000000:
        return b'\xfe' + i.to_bytes(4, 'little')
    elif i < 0x10000000000000000:
        return b'\xff' + i.to_bytes(8, 'little')
    else:
        raise ValueError("integer too large: %d" % (i, ))


def txn_size(tx_ins, tx_outs) -> int:
    """Compute an estimated transaction byte.

    Reference: 
        - https://en.bitcoin.it/wiki/Protocol_documentation#tx
        - https://bitcoin.stackexchange.com/questions/1195/how-to-calculate-transaction-size-before-sending-legacy-non-segwit-p2pkh-p2sh

    """
    return 180 * tx_ins + 140 * tx_outs + 10 + random.randrange(-tx_ins, tx_ins)


@dataclass
class Tx:
    version: int
    tx_ins: List[TxIn]
    tx_outs: List[TxOut]
    locktime: int = 0

    def id(self) -> str:
        # little/big endian conventions require byte order swap
        return sha256(sha256(self.encode()))[::-1].hex()
    
    def fee(self) -> int:
        input_total = sum(tx_in.value() for tx_in in self.tx_ins)
        output_total = sum(tx_out.amount for tx_out in self.tx_outs)
        return input_total - output_total

    def size(self) -> int:
        return txn_size(len(self.tx_ins), len(self.tx_outs))

    def encode(self, sig_index=-1) -> bytes:
        """Encode this transaction as bytes.

        If sig_index is given then return the modified transaction
        encoding of this tx with respect to the single input index.
        This result then constitutes the "message" that gets signed
        by the aspiring transactor of this input.

        """
        out = []
        # encode metadata
        out += [self.version.to_bytes(4, 'little')]
        # encode inputs
        out += [encode_varint(len(self.tx_ins))]

        if sig_index == -1:
            # we are just serializing a fully formed transaction
            out += [tx_in.encode() for tx_in in self.tx_ins]
        else:
            # used when crafting digital signature for a specific input index
            out += [tx_in.encode(script_override=(sig_index == i))
                    for i, tx_in in enumerate(self.tx_ins)]

        # encode outputs
        out += [encode_varint(len(self.tx_outs))]
        out += [tx_out.encode() for tx_out in self.tx_outs]
        # encode... other metadata
        out += [self.locktime.to_bytes(4, 'little')]
        out += [int(1).to_bytes(4, 'little') if sig_index != -1 else b''] # 1 = SIGHASH_ALL

        return b''.join(out)


@dataclass
class TxIn:
    prev_tx: bytes # prev transaction ID: hash256 of prev tx contents
    prev_index: int # UTXO output index in the transaction
    script_sig: Script = None # unlocking script, Script class coming a bit later below
    sequence: int = 0xffffffff # originally intended for "high frequency trades", with locktime
    prev_tx_script_pubkey: Script = None

    def encode(self, script_override: bool = None):
        out = []
        out += [self.prev_tx[::-1]] # little endian vs big endian encodings... sigh
        out += [self.prev_index.to_bytes(4, 'little')]

        if script_override is None:
            # None = just use the actual script
            out += [self.script_sig.encode()]
        elif script_override is True:
            # True = override the script with the script_pubkey of the associated input
            out += [self.prev_tx_script_pubkey.encode()]
        elif script_override is False:
            # False = override with an empty script
            out += [Script([]).encode()]
        else:
            raise ValueError("script_override must be one of None|True|False")

        out += [self.sequence.to_bytes(4, 'little')]

        return b''.join(out)


@dataclass
class TxOut:
    amount: int # in units of satoshi (1e-8 of a bitcoin)
    script_pubkey: Script = None # locking script

    def encode(self):
        out = []
        out += [int(self.amount).to_bytes(8, 'little')]
        out += [self.script_pubkey.encode()]

        return b''.join(out)


@dataclass
class Script:
    cmds: List[Union[int, bytes]]

    def encode(self):
        """Encode to bytes OP instructions and public hash."""
        out = []
        for cmd in self.cmds:
            if isinstance(cmd, int):
                # an int is just an opcode, encode as a single byte
                out += [cmd.to_bytes(1, 'little')]
            elif isinstance(cmd, bytes):
                # bytes represent an element, encode its length and then content
                length = len(cmd)
                # any longer than this requires a bit of tedious handling that we'll skip here
                assert length < 75
                out += [length.to_bytes(1, 'little'), cmd]

        ret = b''.join(out)
        return encode_varint(len(ret)) + ret
    
    def __add__(self, other):
        return Script(self.cmds + other.cmds)

    def __repr__(self):
        repr_int = lambda cmd: OP_CODE_NAMES.get(cmd, 'OP_[{}]'.format(cmd))
        repr_bytes = lambda cmd: cmd.hex()
        repr_cmd = lambda cmd: repr_int(cmd) if isinstance(cmd, int) else repr_bytes(cmd)
        return ' '.join(map(repr_cmd, self.cmds))


def to_satoshi(bitcoin: float) -> float:
    return bitcoin / 1e-8


OP_CODE_NAMES = {
    0: 'OP_0',
    # values 1..75 are not opcodes but indicate elements
    76: 'OP_PUSHDATA1',
    77: 'OP_PUSHDATA2',
    78: 'OP_PUSHDATA4',
    79: 'OP_1NEGATE',
    81: 'OP_1',
    82: 'OP_2',
    83: 'OP_3',
    84: 'OP_4',
    85: 'OP_5',
    86: 'OP_6',
    87: 'OP_7',
    88: 'OP_8',
    89: 'OP_9',
    90: 'OP_10',
    91: 'OP_11',
    92: 'OP_12',
    93: 'OP_13',
    94: 'OP_14',
    95: 'OP_15',
    96: 'OP_16',
    97: 'OP_NOP',
    99: 'OP_IF',
    100: 'OP_NOTIF',
    103: 'OP_ELSE',
    104: 'OP_ENDIF',
    105: 'OP_VERIFY',
    106: 'OP_RETURN',
    107: 'OP_TOALTSTACK',
    108: 'OP_FROMALTSTACK',
    109: 'OP_2DROP',
    110: 'OP_2DUP',
    111: 'OP_3DUP',
    112: 'OP_2OVER',
    113: 'OP_2ROT',
    114: 'OP_2SWAP',
    115: 'OP_IFDUP',
    116: 'OP_DEPTH',
    117: 'OP_DROP',
    118: 'OP_DUP',
    119: 'OP_NIP',
    120: 'OP_OVER',
    121: 'OP_PICK',
    122: 'OP_ROLL',
    123: 'OP_ROT',
    124: 'OP_SWAP',
    125: 'OP_TUCK',
    130: 'OP_SIZE',
    135: 'OP_EQUAL',
    136: 'OP_EQUALVERIFY',
    139: 'OP_1ADD',
    140: 'OP_1SUB',
    143: 'OP_NEGATE',
    144: 'OP_ABS',
    145: 'OP_NOT',
    146: 'OP_0NOTEQUAL',
    147: 'OP_ADD',
    148: 'OP_SUB',
    154: 'OP_BOOLAND',
    155: 'OP_BOOLOR',
    156: 'OP_NUMEQUAL',
    157: 'OP_NUMEQUALVERIFY',
    158: 'OP_NUMNOTEQUAL',
    159: 'OP_LESSTHAN',
    160: 'OP_GREATERTHAN',
    161: 'OP_LESSTHANOREQUAL',
    162: 'OP_GREATERTHANOREQUAL',
    163: 'OP_MIN',
    164: 'OP_MAX',
    165: 'OP_WITHIN',
    166: 'OP_RIPEMD160',
    167: 'OP_SHA1',
    168: 'OP_SHA256',
    169: 'OP_HASH160',
    170: 'OP_HASH256',
    171: 'OP_CODESEPARATOR',
    172: 'OP_CHECKSIG',
    173: 'OP_CHECKSIGVERIFY',
    174: 'OP_CHECKMULTISIG',
    175: 'OP_CHECKMULTISIGVERIFY',
    176: 'OP_NOP1',
    177: 'OP_CHECKLOCKTIMEVERIFY',
    178: 'OP_CHECKSEQUENCEVERIFY',
    179: 'OP_NOP4',
    180: 'OP_NOP5',
    181: 'OP_NOP6',
    182: 'OP_NOP7',
    183: 'OP_NOP8',
    184: 'OP_NOP9',
    185: 'OP_NOP10',
}