# Hyperledger-Voting

## Changing channel configuration with docker CLI

	docker exec -it cli bash

	apt update && apt install -y jq

	export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/ptwist.eu/orderers/orderer.ptwist.eu/msp/tlscacerts/tlsca.ptwist.eu-cert.pem  && export CHANNEL_NAME=mychannel

	peer channel fetch config config_block.pb -o orderer.ptwist.eu:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA

	configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config > config.json
