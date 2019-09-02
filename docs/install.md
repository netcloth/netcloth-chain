
### Install go

Go 1.12.5+ is required.

Install ```go``` from [here](https://golang.org/doc/install)

set env
```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
echo "export GOBIN=\$GOPATH/bin" >> ~/.bash_profile
echo "export PATH=\$PATH:\$GOBIN" >> ~/.bash_profile
echo "export GO111MODULE=on" >> ~/.bash_profile
source ~/.bash_profile
```

### Build and install nch

```bash
# get source code
git clone https://github.com/NetCloth/netcloth-chain.git


# Install the app into your $GOBIN
make install

# check version
nchd version
nchcli version

# Now you should be able to run the following commands:
nchd help
nchcli help
```