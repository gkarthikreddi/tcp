package stack

const (
	ETH_IP        = 0x0800
	ICMP_PRO      = 1
	ICMP_ECHO_REQ = 8
	ICMP_ECHO_REP = 0
)

type ipHeader struct {
	Version     uint8 // IPv4 or IPv6, always 4 for IPv4 (Generally 4 bits)
	IHL         uint8 // Lenght of IP header (Generally 4 bits)
	TOS         uint8
	TotalLength uint16

	// Fragmentation related members, we don't use these
	Identification uint16
	UnusedFlag     bool
	DfFlag         bool
	MoreFlag       bool
	FragOffset     uint16

	TTL       uint8
	Protocol  uint8
	CheckSum  uint16
	SrcIpAddr [4]byte
	DstIpAddr [4]byte
}

func newIpHeader() ipHeader {
	return ipHeader{
		Version: 4,
		IHL:     5, // We will not be using option field, hence hdr size shall always be 5*4 = 20B
		DfFlag:  true,
		TTL:     64,
	}
}
