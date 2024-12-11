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