#! /usr/bin/env python

"""A secured Bitcoin Identity.

Credits: https://karpathy.github.io/2021/06/21/blockchain/

"""

from __future__ import annotations  # PEP 563: Postponed Evaluation of Annotations
from dataclasses import dataclass   # https://docs.python.org/3/library/dataclasses.html
import os

from bitcoin.bitcoin import BITCOIN
from bitcoin.curves import Point
from bitcoin.hashes import ripemd160, sha256


def gen_secret_key(n: int, force: bytes = None) -> int:
    """A random integer that satisfies 1 <= key < n.

    n is the upper bound on the key, typically the order of the elliptic curve
    we are using. The function will return a valid key, i.e. 1 <= key < n.

    For the sake of reproducibility one can force the input string and always
    get the same result.

    """
    while True:
        key = int.from_bytes(force or os.urandom(32), 'big')
        if 1 <= key < n:
            break  # the key is valid, break out
        elif force:
            # the key generated from `force` is not valid, give up
            raise ValueError(f"'{force}' doesnt produce a valid key")

    return key


class PublicKey(Point):
    """The public key is just a Point on a Curve, but has some additional specific
    encoding / decoding functionality that this class implements.

    """

    @classmethod
    def from_point(cls, pt: Point):
        """Promote a Point to be a PublicKey."""
        return cls(pt.curve, pt.x, pt.y)

    @classmethod
    def from_private_key(cls, sk):
        """ sk can be an int or a hex string """
        assert isinstance(sk, (int, str))
        sk = int(sk, 16) if isinstance(sk, str) else sk
        pk = sk * BITCOIN.gen.G

        return cls.from_point(pk)

    def encode(self, compressed: bool, hash160: bool = False) -> str:
        """Return the SEC bytes encoding of the public key Point."""
        # calculate the bytes
        if compressed:
            # (x,y) is very redundant. Because y^2 = x^3 + 7,
            # we can just encode x, and then y = +/- sqrt(x^3 + 7),
            # so we need one more bit to encode whether it was the + or the -
            # but because this is modular arithmetic there is no +/-, instead
            # it can be shown that one y will always be even and the other odd.
            prefix = b'\x02' if self.y % 2 == 0 else b'\x03'
            pkb = prefix + self.x.to_bytes(32, 'big')
        else:
            pkb = b'\x04' + self.x.to_bytes(32, 'big') + self.y.to_bytes(32, 'big')

        # hash if desired
        return ripemd160(sha256(pkb)) if hash160 else pkb

    def address(self, net: str, compressed: bool) -> str:
        """Return the associated bitcoin address for this public key as string."""
        # encode the public key into bytes and hash to get the payload
        pkb_hash = self.encode(compressed=compressed, hash160=True)
        # add version byte (0x00 for Main Network, or 0x6f for Test Network)
        version = {'main': b'\x00', 'test': b'\x6f'}
        ver_pkb_hash = version[net] + pkb_hash
        # calculate the checksum
        checksum = sha256(sha256(ver_pkb_hash))[:4]
        # append to form the full 25-byte binary Bitcoin Address
        byte_address = ver_pkb_hash + checksum
        # finally b58 encode the result
        b58check_address = b58encode(byte_address)

        return b58check_address


# -----------------------------------------------------------------------------
# base58 encoding / decoding utilities
# reference: https://en.bitcoin.it/wiki/Base58Check_encoding


# characters of the alphabet that are very unambiguous. For example it does
# not use ‘O’ and ‘0’, because they are very easy to mess up on paper.
alphabet = '123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz'


def b58encode(b: bytes) -> str:
    """Base58 encoding / decoding utility.

    Reference: https://en.bitcoin.it/wiki/Base58Check_encoding

    """
    assert len(b) == 25  # version is 1 byte, pkb_hash 20 bytes, checksum 4 bytes
    n = int.from_bytes(b, 'big')

    chars = []
    while n:
        n, i = divmod(n, 58)
        chars.append(alphabet[i])

    # special case handle the leading 0 bytes... ¯\_(ツ)_/¯
    num_leading_zeros = len(b) - len(b.lstrip(b'\x00'))
    res = num_leading_zeros * alphabet[0] + ''.join(reversed(chars))

    return res


# -----------------------------------------------------------------------------
# convenience functions


def gen_key_pair(force: str = None):
    """ generate a (secret, public) key pair in one shot """
    sk = gen_secret_key(BITCOIN.gen.n, force)
    pk = PublicKey.from_private_key(sk)
    return sk, pk