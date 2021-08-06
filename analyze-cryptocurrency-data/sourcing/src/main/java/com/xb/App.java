package com.xb;

import java.io.IOException;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class App {
    private static final Logger log = LoggerFactory.getLogger(App.class);

    public static void main(String args[]) throws IOException {
        log.info("crypto data sourcing starting...");
        CoinMarketCap client = new CoinMarketCap();

        List<CMCCoin> coins = client.getCoins();

        log.info("successfully fetched {} coin mappings", coins.size());
        for(CMCCoin coin : coins) {
            log.info(coin.toString());
        }

        // TODO: skip if already exists?
        log.info("storing coins mapping");
        String header = "id,rank,name,symbol,slug";
        boolean success = client.writeCSV("./coins.csv", header, coins);
        log.info("coins dump " + (success ? "successful" : "failed"));
    }
}