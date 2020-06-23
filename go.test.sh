#!/usr/bin/env bash

set -e
echo "" > coverage.txt

#go test ./... -mod=readonly -timeout 12m -race -coverprofile=coverage.txt -covermode=atomic \
go test -mod=readonly -timeout 12m -race -coverprofile=coverage.txt -covermode=atomic \
$(go list ./... | grep -v 'server' | grep -v '/simulation' | grep -v mock | grep -v 'netcloth-chain/tests' | grep -v crypto)

# filter out DONTCOVER
excludelist="$(find ./ -type f -name '*.go' | xargs grep -l 'DONTCOVER')"
excludelist+=" $(find ./ -type f -name '*.pb.go')"
excludelist+=" $(find ./ -type f -path './tests/mocks/*.go')"
for filename in ${excludelist}
do
  filename=$(echo $filename | sed 's/^./github.com\/cosmos\/cosmos-sdk/g')
  echo "Excluding ${filename} from coverage report..."
  sed -i.bak "/$(echo $filename | sed 's/\//\\\//g')/d" coverage.txt
done
