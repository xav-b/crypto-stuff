# [Analyzing cryptocurrency markets using Python](https://blog.patricktriest.com/analyzing-cryptocurrencies-python/?utm_source=pocket_mylist)

## Usage

```sh
# install dependencies
pyenv virtualenv 3.7.4 crypto
pyenv activate crypto
pip install -r requirements.txt

# setup kernel
mkdir $HOME/Library/Jupyter/kernels/crypto
# --> edit kernel.json with your own path
cp kernel.json $HOME/Library/Jupyter/kernels/crypto/

jupyter notebook
```
