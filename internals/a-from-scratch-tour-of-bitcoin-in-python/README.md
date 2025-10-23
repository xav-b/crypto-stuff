# [A from-scratch tour of Bitcoin in Python](https://karpathy.github.io/2021/06/21/blockchain/)

> We are going to create, digitally sign, and broadcast a Bitcoin
> transaction in pure Python, from scratch, and with zero dependencies.

**Disclaimer: this is the work of Andrej, linked above. The code here is purely
*based on his  blog post**

## Requirements

- Python 3.7+

```sh
pyenv virtualenv 3.7.2 crypto
pyenv activate crypto

pip install -r dev-requirements.txt
```

## Utils

[Blog post github repository](https://github.com/karpathy/cryptos)
[Bitcoin testnet faucet](https://bitcoinfaucet.uo1.net/)

## Todo

- [ ] Understand hash algos
- [ ] Go further: [Bitcoins the hard way: Using the raw Bitcoin protocol](http://www.righto.com/2014/02/bitcoins-hard-way-using-raw-bitcoin.html)
