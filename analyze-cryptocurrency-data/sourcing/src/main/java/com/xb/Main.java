package com.xb;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.PrintWriter;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Stream;

import okhttp3.*;
import com.google.gson.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

class CoinPlatform {
    int id;
    String name;
    String symbol;
    String slug;
    String token_address;
}

class CMCCoin {
    int id;
    int rank;
    String name;
    String symbol;
    String slug;
    int is_active;
    // NOTE can we serialize them to dates?
    // NOTE what's about camelcase/snake case
    String first_historical_data;
    String last_historical_data;
    CoinPlatform platform;

    public String toCSVString() {
        return String.join(",", String.valueOf(id), String.valueOf(rank), name, symbol, slug);
    }

    public String toString() {
        return "Coin mapping #" + rank + ": " + name + " " + symbol + " " + slug;
    }
}

class CMCResponseStatus {
    protected String timestamp;
    protected int error_code;
    protected String error_message;
    protected int elapsed;
    protected int credit_count;

    public String getTimestamp() { return timestamp; }
}

class CoinsCMCMapResponse {
    protected List<CMCCoin> data;
    protected CMCResponseStatus status;

    public CMCResponseStatus getStatus() { return status; }
    public List<CMCCoin> getData() { return data; }
}

class CoinMarketCap {
    private String apikey;
    private OkHttpClient client;
    private final String BASE_URL = "https://pro-api.coinmarketcap.com";

    private static final Logger log = LoggerFactory.getLogger(CoinMarketCap.class);

    public CoinMarketCap() {
        apikey = System.getenv("CMC_API_KEY");
        if (apikey != null)
            log.info("configuration successfully found: ********-****-****-****-{}", apikey.substring(23));
        else {
            log.error("CMC API Key not found, aborting");
            System.exit(1);
        }

        // this is a paid API so let's leverage cache
        final int cacheSize = 50 * 1024 * 1024;  // 50 MiB
        final File cacheDirectory = new File("src/main/resources/cache");
        final Cache cache = new Cache(cacheDirectory, cacheSize);

        client = new OkHttpClient.Builder()
            .cache(cache)
            .build();
    }

    public boolean dump(String csvFilename, List<CMCCoin> dataLines) throws FileNotFoundException {
        log.info("writing data to csv: {}", csvFilename);
        try (PrintWriter writer = new PrintWriter(new File(csvFilename))) {
            // write CSV header
            writer.println("id,rank,name,symbol,slug");
            // write rows
            dataLines.stream()
                     .map(coin -> coin.toCSVString())
                     .forEach(writer::println);
        } catch (FileNotFoundException e) {
            log.error(e.getMessage());
            return false;
        }

        return true;
    }

    public List<CMCCoin> getCoins() throws IOException {
        final Gson gson = new Gson();
        final String endpoint = "/v1/cryptocurrency/map";

        final Request request = new Request.Builder()
                .url(BASE_URL + endpoint)
                .addHeader("Content-Type", "application/json")
                .addHeader("Accept", "application/json")
                .addHeader("X-CMC_PRO_API_KEY", apikey)
                .build();

        log.info("querying " + endpoint);
        Response response = client.newCall(request).execute();

        if (!response.isSuccessful()) throw new IOException("Unexpected code " + response);
        else log.debug("response code: " + response.code());

        CoinsCMCMapResponse payload = gson.fromJson(response.body().charStream(), CoinsCMCMapResponse.class);

        // TODO check payload.getStatus();

        return payload.getData();
    }
}

public class Main {
    private static final Logger log = LoggerFactory.getLogger(Main.class);

    public static void main(String args[]) throws IOException {
        log.info("crypto data sourcing starting...");
        CoinMarketCap client = new CoinMarketCap();

        List<CMCCoin> coins = client.getCoins();

        log.info("successfully fetched {} coin mappings", coins.size());
        for(CMCCoin coin : coins) {
            log.info(coin.toString());
        }

        log.info("storing coins mapping");
        boolean success = client.dump("./coins.csv", coins);
        log.info("coins dump " + (success ? "successful" : "failed"));
    }
}