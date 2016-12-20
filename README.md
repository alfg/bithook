bithook
=======

Bithook is a Bitcoin Webhook CLI Utility. You can use this program to listen for new blocks and
transactions for addresses. When a transaction is received, Bithook will fire off a web POST 
request any URL you specify along with the transaction JSON data.

Bithook uses the [Blockchain.info](http://blockchain.info) websocket API to listen for transactions.

## Features
* Subscribe to new blocks
* Subscribe to transactions for a specified address
* Subscribe to unconfirmed addresses
* Sends JSON data as a POST request to any URL

## Install From Source

```
go get github.com/alfg/bithook
./bin/bithook help
```

## Install from Homebrew
```
brew cask alfg/tap
brew install alfg/tap/bithook
bithook help
```

## Usage

`bithook <command> -webhook=http://webhook/path`

```
Commands:
bithook blocks -- Subscribe to new blocks.
bithook unconfirmed -- Subscribe to new unconfirmed transactions.
bithook address <address> -- Subscribe to address.
bithook test -- Receives latest transaction. Use for testing.
bithook help -- This help menu.
bithook version -- This version.

Flags:
-webhook=<webhook path> -- This is optional. The results will just echo to output if flag not set.
```

#### Example

The following example will listen for transactions of 1dice8EMZmqKvrGE4Qc9bUFf9PX3xaYDp and POST the json results to http://requestb.in/nt5bcnnt. You can view the results at http://requestb.in/nt5bcnnt?inspect.

`bithook address 1dice8EMZmqKvrGE4Qc9bUFf9PX3xaYDp -webhook=http://requestb.in/nt5bcnnt`

## Develop

```
git clone git@github.com:alfg/bithook.git
export GOPATH=$HOME/path/to/project
cd /to/project
go run main.go
```

## License
MIT
