# wanutil
#### CLI utility tool for Wanchain

**WIP:** please note that development is in progress

## 1. Setup
Fetch and build wanutil
```
go get github.com/jsgoyette/wanutil
cd $GOPATH/src/!$
make

# and if you want to install the binary
go install
```

Set up the config file
```
mkdir ~/.wanutil
cp config.yml.example ~/.wanutil/config.yml

# edit config as needed
vi ~/.wanutil/config.yml
```

## 2. Usage

#### Show help
```
wanutil help
```

#### Get WAN balance
```
wanutil bal -a 0xecb4e4073a9bf5e024ee68d1f871635f1888030e
```

#### Get WAN balance at block number
```
wanutil bal -a 0xecb4e4073a9bf5e024ee68d1f871635f1888030e -b 1600000
```

#### Get token balance
```
wanutil bal -a 0xecb4e4073a9bf5e024ee68d1f871635f1888030e -t WETH
```

#### Get transaction
```
wanutil tx -hash ox48b53118a7ebaa8f1a587f12a1a1710dc38b578b6ef564b3b4caa2361551e368
```

#### Scan blockchain for transactions to an address, starting from block 1600000
```
wanutil scan -a 0xecb4e4073a9bf5e024ee68d1f871635f1888030e -b 1600000
```
