package main

// Config represents the configuration for a role.
type Config struct {
	CertPath     string `json:"certPath"`
	KeyDirectory string `json:"keyPath"`
	TLSCertPath  string `json:"tlsCertPath"`
	PeerEndpoint string `json:"peerEndpoint"`
	GatewayPeer  string `json:"gatewayPeer"`
	MSPID        string `json:"mspID"`
}

// Create a Profile map
var profile = map[string]Config{

	"manufacturer": {
		CertPath:     "../Fabric-network/organizations/peerOrganizations/manufacturer.auto.com/users/User1@manufacturer.auto.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Fabric-network/organizations/peerOrganizations/manufacturer.auto.com/users/User1@manufacturer.auto.com/msp/keystore/",
		TLSCertPath:  "../Fabric-network/organizations/peerOrganizations/manufacturer.auto.com/peers/peer0.manufacturer.auto.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.manufacturer.auto.com",
		MSPID:        "ManufacturerMSP",
	},

	"dealer": {
		CertPath:     "../Fabric-network/organizations/peerOrganizations/dealer.auto.com/users/User1@dealer.auto.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Fabric-network/organizations/peerOrganizations/dealer.auto.com/users/User1@dealer.auto.com/msp/keystore/",
		TLSCertPath:  "../Fabric-network/organizations/peerOrganizations/dealer.auto.com/peers/peer0.dealer.auto.com/tls/ca.crt",
		PeerEndpoint: "localhost:9051",
		GatewayPeer:  "peer0.dealer.auto.com",
		MSPID:        "DealerMSP",
	},

	"pharmacies": {
		CertPath:     "../Fabric-network/organizations/peerOrganizations/pharmacies.auto.com/users/User1@pharmacies.auto.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Fabric-network/organizations/peerOrganizations/pharmacies.auto.com/users/User1@pharmacies.auto.com/msp/keystore/",
		TLSCertPath:  "../Fabric-network/organizations/peerOrganizations/pharmacies.auto.com/peers/peer0.pharmacies.auto.com/tls/ca.crt",
		PeerEndpoint: "localhost:11051",
		GatewayPeer:  "peer0.pharmacies.auto.com",
		MSPID:        "PharmaciesMSP",
	},
}
