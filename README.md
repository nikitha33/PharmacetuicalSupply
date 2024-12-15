# Steps to build Fabric network using  and list fabric ca and list all the containers

```
* Create folder fabric network
```
```
* Register the ca admin for each organization
```
Build the docker-compose-ca.yaml in the docker folder
```
docker compose -f docker/docker-compose-ca.yaml up -d
```
Giving permission for all the folders to use:
```
sudo chmod -R 777 organizations/
```
```
******Register and enroll the users for each organization******
```
Build the registerEnroll.sh script file
```
chmod +x registerEnroll.sh

./registerEnroll.sh
```
```
******Build the infrastructure******
```
Build the docker-compose-3org.yaml in the docker folder

```
docker compose -f docker/docker-compose-3org.yaml up -d
```
```
*******Generate the genesis block*******
```
Build the configtx.yaml file in the config folder

```
export FABRIC_CFG_PATH=./config

export CHANNEL_NAME=autochannel
```
```
configtxgen -profile ThreeOrgsChannel -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
```
```
******Create the application channel*******
```
```
export ORDERER_CA=./organizations/ordererOrganizations/auto.com/orderers/orderer.auto.com/msp/tlscacerts/tlsca.auto.com-cert.pem

export ORDERER_ADMIN_TLS_SIGN_CERT=./organizations/ordererOrganizations/auto.com/orderers/orderer.auto.com/tls/server.crt

export ORDERER_ADMIN_TLS_PRIVATE_KEY=./organizations/ordererOrganizations/auto.com/orderers/orderer.auto.com/tls/server.key
Command to give permission to script file to execute
```
```
osnadmin channel join --channelID $CHANNEL_NAME --config-block ./channel-artifacts/$CHANNEL_NAME.block -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
```
```
osnadmin channel list -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
```
![alt text](<Screenshot from 2024-12-09 18-05-32.png>)

# To submit transaction as ManufacturerMSP

```
export CHANNEL_NAME=autochannel
export FABRIC_CFG_PATH=./peercfg
export CORE_PEER_LOCALMSPID=ManufacturerMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/manufacturer.auto.com/peers/peer0.manufacturer.auto.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/manufacturer.auto.com/users/Admin@manufacturer.auto.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/auto.com/orderers/orderer.auto.com/msp/tlscacerts/tlsca.auto.com-cert.pem
export MANUFACTURER_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/manufacturer.auto.com/peers/peer0.manufacturer.auto.com/tls/ca.crt
export DEALER_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/dealer.auto.com/peers/peer0.dealer.auto.com/tls/ca.crt
export PHARMACIES_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/pharmacies.auto.com/peers/peer0.pharmacies.auto.com/tls/ca.crt

```

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.auto.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Project-Pharma --peerAddresses localhost:7051 --tlsRootCertFiles $MANUFACTURER_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $DEALER_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $PHARMACIES_PEER_TLSROOTCERT -c '{"function":"CreateDrug","Args":["01", "Aspirin", "BharatPharma", "Manufactured"]}'
```
# To read a drug and to add it sucessfully
```
 peer chaincode query -C $CHANNEL_NAME -n Project-Pharma -c '{"Args":["CreateDrug","serialNumber","drugName","manufacturer"]}'
 ```

 # setting up the environment for DEALERMSP
 ```
export CHANNEL_NAME=autochannel
export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=autochannel 
export CORE_PEER_LOCALMSPID=DealerMSP 
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=localhost:9051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/dealer.auto.com/peers/peer0.dealer.auto.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/dealer.auto.com/users/Admin@dealer.auto.com/msp
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/auto.com/orderers/orderer.auto.com/msp/tlscacerts/tlsca.auto.com-cert.pem
export MANUFACTURER_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/manufacturer.auto.com/peers/peer0.manufacturer.auto.com/tls/ca.crt
export DEALER_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/dealer.auto.com/peers/peer0.dealer.auto.com/tls/ca.crt
export PHARMACIES_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/pharmacies.auto.com/peers/peer0.pharmacies.auto.com/tls/ca.crt

 ```
 # To invoke the private transcation setting up environment 
```
export SERIALNUMBER=$(echo -n "SN12345" | base64 | tr -d \\n)

export DRUGNAME=$(echo -n "PainKiller" | base64 | tr -d \\n)

export MANUFACTURER=$(echo -n "PharmaInc" | base64 | tr -d \\n)

export BATCHNUMBER=$(echo -n "Batch001" | base64 | tr -d \\n)

export EXPIRYDATE=$(echo -n "2025-12-31" | base64 | tr -d \\n)

export QUANTITY=$(echo -n "100" | base64 | tr -d \\n)

export STATUS=$(echo -n "Manufactured" | base64 | tr -d \\n)
```
# To Invoke the chaincode in dealermsp
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.auto.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Project-Pharma --peerAddresses localhost:7051 --tlsRootCertFiles $MANUFACTURER_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $DEALER_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $PHARMACIES_PEER_TLSROOTCERT -c '{"Args":["SmartContract:CreateDrug","SN12345"]}' --transient "{\"drugName\":\"$DRUGNAME\",\"manufacturer\":\"$MANUFACTURER\",\"batchNumber\":\"$BATCHNUMBER\",\"expiryDate\":\"$EXPIRYDATE\",\"quantity\":\"$QUANTITY\",\"status\":\"$STATUS\"}"

```
# To read the query
```
peer chaincode query -C $CHANNEL_NAME -n Project-Pharma -c '{"Args":["SmartContract:ReadDrug","SN12345"]}'
```

 ########Creating the client folder############

 ```
 cd ..

 mkdir client
 ```

 Building the client application

 ```
 cd Client

go mod init client

Create and build profile.go, connect.go, client.go, main.go

go get google.golang.org/grpc@v1.67.1

go mod tidy

go run .
```
# To register and enroll a new user 
```
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/manufacturer.auto.com/

fabric-ca-client register --caname ca-manufacturer --id.name user2 --id.secret user2pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/manufacturer/ca-cert.pem"

fabric-ca-client enroll -u https://user2:user2pw@localhost:7054 --caname ca-manufacturer -M "${PWD}/organizations/peerOrganizations/manufacturer.auto.com/users/User2@manufacturer.auto.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/manufacturer/ca-cert.pem"  
```



