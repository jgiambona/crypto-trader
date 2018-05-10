# CryptoTrader

CrayptoTrader automates building of volume for a certain token. Currently hardcoded
in the system is NOX to ETH trade.

## Production

#### Prerequisites:

+ Docker 18.1
+ Grafana
+ InfluxDB

The first thing to do in order to run system is to download docker and pull golang image into the host system.

## Development

#### Prerequisites:

+ Go 1.10

Assuming, Go is already in the system execute this in the workspace root.

~~~
$ make
~~~

Set the configuration `config.yaml` and restart.

~~~
$ docker run --rm -it -v $PWD:/go/src/github.com/ffimnsr/trader -p 8000:8000 golang:1.10.1 bash
> go get -u github.com/golang/dep/cmd/dep
> dep ensure
> make
~~~

## Frontend Settings

![alt text](https://raw.githubusercontent.com/jgiambona/crypto-trader/master/images/bot-settings.png)

#### Account Configuration:

Need two accounts to bounce buy and sell of trades. If in case one account is missing it will miss the buy/sell trade.

~~~
API Key     - Livecoin generated API key
API Secret  - Livecoin secret key
~~~

After filling up the details, save the new account configuration. The popup notification should state that the bot account is saved.
Otherwise there is an error.

#### Rule Configuration:

If the rule configuration is not filled up. It will run the default bot configuration.

~~~
Interval                        - May come in seconds, minutes, hours, and days (e.g. 1s means 1 second, 1m for 1 minute). Default is 7 seconds.
Maximum Volume                  - Is the maximum volume the bot will reach before put to halt.
Transaction Volume              - Is the minimum quantity of NOX token needed to deal in the exchange.
Variance of Transaction Volume  - The random percentage on which it will add to transaction volume.
Bid Price Step Down             - Is the deduction it will do the achieve the lowest price ask in exchange.
Minimum Bid Price               - The minimum bid price on which the bot will play if the price is lower then bot will halt.
~~~

Change all the details, both rules have separate interval rules.

#### What is currently on the road map and WIP?

[ ] Proper interaction to Grafana UI
[ ] Proper handling of other crypto currency pairs

