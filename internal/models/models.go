package models

// type ConfigPath struct {
// 	NetworkConfigPath string            `json:"networkConfigPath"`
// 	NetworkDirectory  string            `json:"networkDirectory"`
// 	Orgs              map[string]string `json:"orgs,omitempty"`
// }

type Network struct {
	NetworkDirectory string            `json:"networkDirectory"`
	NetworkType      string            `json:"networkType"`
	ChannelName      string            `json:"channelName"`
	NetworkName      string            `json:"networkName"`
	ChaincodeName    string            `json:"chaincodeName"`
	ChaincodePath    string            `json:"chaincodePath"`
	ChaincodePkgPath string            `json:"chaincodePkgPath"`
	Orgs             map[string]string `json:"orgs,omitempty"`
}

type Ca struct {
	TlscaPort string `json:"tlscaPort"`
	OrgcaPort string `json:"orgcaPort"`
}

type Peer struct {
	PeerLis string `json:"peerLis"`
	PeerCc  string `json:"peerCc"`
	PeerOp  string `json:"peerOp"`
	Couchdb string `json:"couchdb"`
}

type Orderer struct {
	OrdererLis   string `json:"ordererLis"`
	OrdererAdm   string `json:"ordererAdm"`
	OrdererAdmOp string `json:"ordererAdmOp"`
}

type Peer_Orderer struct {
	OrdererName string `json:"ordererName"`
	OrdererLis  string `json:"ordererLis"`
}

type Org_Config struct {
	OrgName string       `json:"orgName"`
	OrgType string       `json:"orgType,omitempty"`
	Ca      Ca           `json:"ca"`
	Peers   []Peer       `json:"peers"`
	Orderer Peer_Orderer `json:"orderer"`
}

type Orderer_Config struct {
	OrgName  string    `json:"orgName"`
	OrgType  string    `json:"orgType,omitempty"`
	Ca       Ca        `json:"ca"`
	Orderers []Orderer `json:"orderers"`
}

type Network_Org_Config struct {
	Config  Org_Config
	Network Network
}

type Network_Orderer_Config struct {
	Config  Orderer_Config
	Network Network
}

type PdEndorsementPolicy struct {
	SignaturePolicy string `json:"signaturePolicy"`
}

type PdDetails struct {
	Name              string              `json:"name"`
	Policy            string              `json:"policy"`
	RequiredPeerCount int                 `json:"requiredPeerCount"`
	MaxPeerCount      int                 `json:"maxPeerCount"`
	BlockToLive       int                 `json:"blockToLive"`
	MemberOnlyRead    bool                `json:"memberOnlyRead"`
	MemberOnlyWrite   bool                `json:"memberOnlyWrite"`
	EndorsementPolicy PdEndorsementPolicy `json:"endorsementPolicy"`
}
