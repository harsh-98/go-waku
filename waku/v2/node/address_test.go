package node

import (
	"net"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

func TestExternalAddressSelection(t *testing.T) {
	a1, _ := ma.NewMultiaddr("/ip4/192.168.0.106/tcp/60000/p2p/16Uiu2HAmUVVrJo1KMw4QwUANYF7Ws4mfcRqf9xHaaGP87GbMuY2f")                                                                                                                // Valid
	a2, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/60000/p2p/16Uiu2HAmUVVrJo1KMw4QwUANYF7Ws4mfcRqf9xHaaGP87GbMuY2f")                                                                                                                    // Valid but should not be prefered
	a3, _ := ma.NewMultiaddr("/ip4/192.168.1.20/tcp/19710/p2p/16Uiu2HAmUVVrJo1KMw4QwUANYF7Ws4mfcRqf9xHaaGP87GbMuY2f")                                                                                                                 // Valid
	a4, _ := ma.NewMultiaddr("/dns4/www.status.im/tcp/2012/ws/p2p/16Uiu2HAmUVVrJo1KMw4QwUANYF7Ws4mfcRqf9xHaaGP87GbMuY2f")                                                                                                             // Invalid (it's useless)
	a5, _ := ma.NewMultiaddr("/dns4/www.status.im/tcp/443/wss/p2p/16Uiu2HAmUVVrJo1KMw4QwUANYF7Ws4mfcRqf9xHaaGP87GbMuY2f")                                                                                                             // Valid
	a6, _ := ma.NewMultiaddr("/ip4/192.168.1.20/tcp/19710/wss/p2p/16Uiu2HAmUVVrJo1KMw4QwUANYF7Ws4mfcRqf9xHaaGP87GbMuY2f")                                                                                                             // Invalid (local + wss)
	a7, _ := ma.NewMultiaddr("/ip4/192.168.1.20/tcp/19710/ws/p2p/16Uiu2HAmUVVrJo1KMw4QwUANYF7Ws4mfcRqf9xHaaGP87GbMuY2f")                                                                                                              // Invalid  (it's useless)
	a8, _ := ma.NewMultiaddr("/dns4/node-02.gc-us-central1-a.status.prod.statusim.net/tcp/30303/p2p/16Uiu2HAmDQugwDHM3YeUp86iGjrUvbdw3JPRgikC7YoGBsT2ymMg/p2p-circuit/p2p/16Uiu2HAmUVVrJo1KMw4QwUANYF7Ws4mfcRqf9xHaaGP87GbMuY2f")     // VALID
	a9, _ := ma.NewMultiaddr("/dns4/node-02.gc-us-central1-a.status.prod.statusim.net/tcp/443/wss/p2p/16Uiu2HAmDQugwDHM3YeUp86iGjrUvbdw3JPRgikC7YoGBsT2ymMg/p2p-circuit/p2p/16Uiu2HAmUVVrJo1KMw4QwUANYF7Ws4mfcRqf9xHaaGP87GbMuY2f")   // VALID
	a10, _ := ma.NewMultiaddr("/dns4/node-01.gc-us-central1-a.wakuv2.test.statusim.net/tcp/8000/wss/p2p/16Uiu2HAmJb2e28qLXxT5kZxVUUoJt72EMzNGXB47Rxx5hw3q4YjS/p2p-circuit/p2p/16Uiu2HAmUVVrJo1KMw4QwUANYF7Ws4mfcRqf9xHaaGP87GbMuY2f") // VALID
	a11, _ := ma.NewMultiaddr("/dns4/node-01.gc-us-central1-a.wakuv2.test.statusim.net/tcp/30303/p2p/16Uiu2HAmJb2e28qLXxT5kZxVUUoJt72EMzNGXB47Rxx5hw3q4YjS/p2p-circuit/p2p/16Uiu2HAmUVVrJo1KMw4QwUANYF7Ws4mfcRqf9xHaaGP87GbMuY2f")    // VALID

	addrs := []ma.Multiaddr{a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11}

	w := &WakuNode{}
	extAddr, multiaddr, err := w.getENRAddresses([]ma.Multiaddr{a1, a2, a3, a4, a5, a6, a7, a8})

	a5NoP2P, _ := decapsulateP2P(a5)
	require.NoError(t, err)
	require.Equal(t, extAddr.IP, net.IPv4(192, 168, 0, 106))
	require.Equal(t, extAddr.Port, 60000)
	require.Equal(t, multiaddr[0].String(), a5NoP2P.String())
	require.Len(t, multiaddr, 1) // Should only have 1, without circuit relay

	a12, _ := ma.NewMultiaddr("/ip4/188.23.1.8/tcp/30303/p2p/16Uiu2HAmUVVrJo1KMw4QwUANYF7Ws4mfcRqf9xHaaGP87GbMuY2f") // VALID
	addrs = append(addrs, a12)
	extAddr, _, err = w.getENRAddresses(addrs)
	require.NoError(t, err)
	require.Equal(t, extAddr.IP, net.IPv4(188, 23, 1, 8))
	require.Equal(t, extAddr.Port, 30303)

	a8RelayNode, _ := decapsulateCircuitRelayAddr(a8)
	_, multiaddr, err = w.getENRAddresses([]ma.Multiaddr{a1, a8})
	require.NoError(t, err)
	require.Len(t, multiaddr, 1)
	require.Equal(t, multiaddr[0].String(), a8RelayNode.String()) // Should have included circuit-relay addr
}
