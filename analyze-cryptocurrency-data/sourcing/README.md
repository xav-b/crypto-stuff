# Crypto data sourcing - Java

!! Written in java for educational purpose, should be in Python !!

## Guide

**Runtime**

```sh
$ java -version
java version "1.8.0_221"
Java(TM) SE Runtime Environment (build 1.8.0_221-b11)
Java HotSpot(TM) 64-Bit Server VM (build 25.221-b11, mixed mode)

$ javac -version
javac 1.8.0_221

$ mvn -version
Apache Maven 3.8.1 (05c21c65bdfed0f71a2f2ada8b84da59348c4c5d)
Maven home: /usr/local/Cellar/maven/3.8.1/libexec
Java version: 15.0.2, vendor: N/A, runtime: /usr/local/Cellar/openjdk/15.0.2/libexec/openjdk.jdk/Contents/Home
Default locale: en_SG, platform encoding: UTF-8
OS name: "mac os x", version: "11.4", arch: "x86_64", family: "mac"
```

```sh
# build the jar in `target`
mvn clean package

# run it
java -cp target/crypto-1.jar com.xb.Main
```

FIXME: doesn't include dependencies in the jar. This does but I don't
think this is the way to package an app.

```sh
# build
mvn clean compile
# run
mvn exec:java -Dexec.mainClass="com.xb.Main"
```

---

## Todo

- [ ] Properly structure the project files
- [ ] Have 2 modules: one CMC to CSV / Another CSV to Database
- [ ] make it a CLI
