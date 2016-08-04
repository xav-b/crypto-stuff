#! /usr/bin/env python
# -*- coding: utf-8 -*-
# vim:fenc=utf-8

"""
Bitcoin historic rate analysis.
"""

import os

import quandl

SOURCE = 'GDAX/USD'
TOKEN = os.environ.get('QUANDL_API_TOKEN')

# TODO Add authentification
# quandl.ApiConfig.api_key = 'tEsTkEy123456789'
# quandl.ApiConfig.api_version = '2015-04-09'

data = quandl.get(source)
