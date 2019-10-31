#!/bin/sh

./create.sh sw1 host11 host12
./add-intf.sh A host11 sw1
./add-intf.sh B sw1 host12

./create.sh sw2 host21 host22
./add-intf.sh A host21 sw2
./add-intf.sh B sw2 host22

./add-intf.sh C sw1 sw2
