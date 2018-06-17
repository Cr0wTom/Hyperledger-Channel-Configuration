# Hyperledger-Voting

## Changing channel configuration with docker CLI (First network example case)
	#Generate the crypto material for the new org join
	cd org3-artifacts
	../../bin/cryptogen generate --config=./org3-crypto.yaml
	export FABRIC_CFG_PATH=$PWD && ../../bin/configtxgen -printOrg Org3MSP > ../channel-artifacts/org3.json
	cd ../ && cp -r crypto-config/ordererOrganizations org3-artifacts/crypto-config/

	#Execute the Org1 bash CLI
	docker exec -it cli bash

	#Install the jq tool
	apt update && apt install -y jq

	#Export the ORDERER_CA and CHANNEL_NAME variables
	export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/ptwist.eu/orderers/orderer.ptwist.eu/msp/tlscacerts/tlsca.ptwist.eu-cert.pem  && export CHANNEL_NAME=mychannel

	#Fetch the configuration
	peer channel fetch config config_block.pb -o orderer.ptwist.eu:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA

	#Convert the Configuration to JSON and Trim It Down
	configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config > config.json

		#Copy JSON file for edit
		docker cp <containerId>:/file/path/within/container /host/path/target

	#Add the new org crypto material
	jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"Org3MSP":.[1]}}}}}' config.json ./channel-artifacts/org3.json > modified_config.json

	configtxlator proto_encode --input config.json --type common.Config --output config.pb

	configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb

	configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output org3_update.pb

	configtxlator proto_decode --input org3_update.pb --type common.ConfigUpdate | jq . > org3_update.json

	echo '{"payload":{"header":{"channel_header":{"channel_id":"mychannel", "type":2}},"data":{"config_update":'$(cat org3_update.json)'}}}' | jq . > org3_update_in_envelope.json

	configtxlator proto_encode --input org3_update_in_envelope.json --type common.Envelope --output org3_update_in_envelope.pb

	#Sign the configuration update by channel admins
	peer channel signconfigtx -f org3_update_in_envelope.pb

	export CORE_PEER_LOCALMSPID="Org2MSP" &&	export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt && export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp && export CORE_PEER_ADDRESS=peer0.org2.example.com:7051

	peer channel update -f org3_update_in_envelope.pb -c $CHANNEL_NAME -o orderer.example.com:7050 --tls --cafile $ORDERER_CA
