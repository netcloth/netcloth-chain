### Requirement
 * macOS / Ubuntu 14.04 LTS +
 * RAM: 8 GB+
 * Disk: 100 GB+

### Install go

Go 1.12.5+ is required.

Install ```go``` from [here](https://golang.org/doc/install)

```bash
# download go package
wget https://dl.google.com/go/go1.12.2.linux-amd64.tar.gz
tar -xvf go1.12.2.linux-amd64.tar.gz
mv go $HOME/
```

set env
```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
echo "export GOBIN=\$GOPATH/bin" >> ~/.bash_profile
echo "export PATH=\$PATH:\$GOBIN" >> ~/.bash_profile
echo "export GO111MODULE=on" >> ~/.bash_profile
source ~/.bash_profile
```

check go version
```cassandraql
go version

```

### Build and install nch

```bash
# get source code
git clone https://github.com/NetCloth/netcloth-chain.git
cd netcloth-chain && git checkout v1.0.0

# Install the app into your $GOBIN
make install

# check version
nchd version
nchcli version

# Now you should be able to run the following commands:
nchd help
nchcli help
```
