#!/usr/bin/env python

from dataclasses import dataclass

import asyncio
from binance import AsyncClient, BinanceSocketManager


@dataclass
class TickerStream:
    """Individual Symbol Ticker Streams.

    24hr rollwing window ticker statistics for a single symbol. These are NOT
    the statistics of the UTC day, but a 24hr rolling window from requestTime to
    24hrs before.
    
    https://binance-docs.github.io/apidocs/futures/en/#individual-symbol-ticker-streams
    
    """
    event_time: int
    symbol: str
    price_change: float
    price_change_pct: float
    weighted_avg_price: float
    ytc_close: float
    today_close: float
    close_trade_qty: int
    best_bid_price: float
    best_bid_qty: int
    best_ask_price: float
    best_ask_qty: int
    open_price: float
    high_price: float
    low_price: float
    tot_traded_base_volume: int
    tot_traded_quote_volume: int
    statistics_open_time: int
    statistics_close_time: int
    first_trade_id: int
    last_trade_id: int
    tot_trades: int
    event_type: str = '24hrTicker'


async def kline_listener(client, symbol):
    bm = BinanceSocketManager(client, user_timeout=60)

    async with bm.symbol_ticker_socket(symbol=symbol) as stream:
        while True:
            msg = await stream.recv()
            ticker = TickerStream(
                event_time=msg['E'],
                symbol=msg['s'],
                price_change=msg['p'],
                price_change_pct=msg['P'],
                weighted_avg_price=msg['w'],
                ytc_close=msg['x'],
                today_close=msg['c'],
                close_trade_qty=msg['Q'],
                best_bid_price=msg['b'],
                best_bid_qty=msg['B'],
                best_ask_price=msg['a'],
                best_ask_qty=msg['A'],
                open_price=msg['o'],
                high_price=msg['h'],
                low_price=msg['l'],
                tot_traded_base_volume=msg['v'],
                tot_traded_quote_volume=msg['q'],
                statistics_open_time=msg['O'],
                statistics_close_time=msg['C'],
                first_trade_id=msg['F'],
                last_trade_id=msg['L'],
                tot_trades=msg['n'],
            )
            print(ticker.event_time, ticker.symbol, ticker.price_change, ticker.price_change_pct)


async def main():
    client = await AsyncClient.create()

    status = await client.get_system_status()
    print(status)

    symbol = 'BNBBTC'
    await kline_listener(client, symbol)

if __name__ == "__main__":

    loop = asyncio.get_event_loop()
    loop.run_until_complete(main())