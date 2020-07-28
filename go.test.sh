#!/usr/bin/env bash

set -e
echo "" > coverage.txt

go test -mod=readonly -timeout 20m -race -coverprofile=coverage.txt -covermode=atomic \
$(go list ./... | grep -v 'netcloth-chain/server' | grep -v '/simulation' | grep -v mock | grep -v 'netcloth-chain/tests' | grep -v 'crypto' | grep -v '/simapp' | grep -v '/app/genesis')

# filter out DONTCOVER
excludelist="$(find . -type f -name '*.go' | xargs grep -l 'DONTCOVER')"
excludelist+=" $(find . -type f -name '*.pb.go')"
excludelist+=" $(find . -type f -name 'test_common.go')"
excludelist+=" $(find . -type f -name 'common_test.go')"
excludelist+=" $(find . -type f -path './tests/mocks/*.go')"
for filename in ${excludelist}
do
  filename=$(echo $filename | sed 's/^./github.com\/netcloth\/netcloth-chain/g')
  echo "Excluding ${filename} from coverage report..."
  sed -i.bak "/$(echo $filename | sed 's/\//\\\//g')/d" coverage.txt
done
