"""Testing Bitcoin curve and transaction security."""

import random

from bitcoin.bitcoin import BITCOIN
from bitcoin.hashes import sha256, ripemd160


random.seed(1337)

G = BITCOIN.gen.G
p = G.curve.p


def test_g_on_curve():
    assert (G.y**2 - G.x**3 - 7) % p == 0


def test_random_point_not_on_the_curve():
    """Some other totally random point will of course not be on the curve, _MOST_ likely."""
    x = random.randrange(0, p)
    y = random.randrange(0, p)

    assert (y**2 - x**3 - 7) % p != 0


def test_public_keys_on_curve():
    sk, pk = 1, G
    assert (pk.y**2 - pk.x**3 - 7) % p == 0

    sk, pk = 2, G + G
    assert (pk.y**2 - pk.x**3 - 7) % p == 0

    sk, pk = 3, G + G + G
    assert (pk.y**2 - pk.x**3 - 7) % p == 0


def test_g_multiplication():
    assert G == 1*G
    assert G + G == 2*G
    assert G + G + G == 3*G


def test_custom_sha256_empty():
    assert sha256(b'').hex() == 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855'


def test_custom_ripemd160():
    assert ripemd160(b'').hex() == '9c1185a5c5e9fc54612808977ee8f548b2258d31'
