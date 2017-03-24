package client

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net"

	"github.com/Impyy/kci-demo/cmd/kci-demo/crypto"
	"github.com/Impyy/kci-demo/cmd/kci-demo/profile"
)

type Client struct {
	conn     *net.UDPConn
	prof     *profile.Profile
	lastBlob *[crypto.NonceSize]byte
	kci      bool
}

func New(prof *profile.Profile, kci bool) (*Client, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", "")
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
		prof: prof,
		kci:  kci,
	}, nil
}

func (c *Client) DecryptPacket(packet CryptoPacket) (Packet, error) {
	var sharedKey *[crypto.KeySize]byte
	if c.kci {
		sharedKey = crypto.PrecomputeKey(c.prof.PublicKey, c.prof.SecretKey)
	} else {
		sharedKey = crypto.PrecomputeKey(packet.PublicKey, c.prof.SecretKey)
	}

	fmt.Printf("shared key: %s\n", hex.EncodeToString(sharedKey[:]))

	bytes, err := crypto.Decrypt(packet.Data, sharedKey, packet.Nonce)
	if err != nil {
		return nil, err
	}

	var res Packet
	switch packet.Type {
	case packetIDHandshakeRequest:
		res = &HandshakeRequestPacket{}
	case packetIDHandshakeResponse:
		res = &HandshakeResponsePacket{}
	default:
		return nil, errors.New("wtf")
	}

	if err := res.UnmarshalBinary(bytes); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) EncryptPacket(packet Packet, publicKey *[crypto.KeySize]byte) (Marshalable, error) {
	bytes, err := packet.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var sharedKey *[crypto.KeySize]byte
	if c.kci {
		sharedKey = crypto.PrecomputeKey(c.prof.PublicKey, c.prof.SecretKey)
	} else {
		sharedKey = crypto.PrecomputeKey(publicKey, c.prof.SecretKey)
	}

	fmt.Printf("shared key: %s\n", hex.EncodeToString(sharedKey[:]))

	encryptedBytes, nonce, err := crypto.Encrypt(bytes, sharedKey)
	if err != nil {
		return nil, err
	}

	return &CryptoPacket{
		Type:      packet.ID(),
		PublicKey: c.prof.PublicKey,
		Nonce:     nonce,
		Data:      encryptedBytes,
	}, nil
}

func (c *Client) Send(packet *RawPacket) error {
	_, err := c.conn.WriteTo(packet.Data, packet.Addr)
	return err
}

func (c *Client) HandshakeResponse(addr *net.UDPAddr, publicKey *[crypto.KeySize]byte, blob *[crypto.NonceSize]byte) error {
	packet := HandshakeResponsePacket{Blob: blob}
	encryptedPacket, err := c.EncryptPacket(&packet, publicKey)
	if err != nil {
		return err
	}

	data, err := encryptedPacket.MarshalBinary()
	if err != nil {
		return err
	}

	return c.Send(&RawPacket{Data: data, Addr: addr})
}

func (c *Client) Handshake(addr *net.UDPAddr, publicKey *[crypto.KeySize]byte) error {
	blob, err := crypto.GenerateNonce()
	if err != nil {
		return err
	}
	c.lastBlob = blob

	packet := HandshakeRequestPacket{Blob: blob}
	encryptedPacket, err := c.EncryptPacket(&packet, publicKey)
	if err != nil {
		return err
	}

	data, err := encryptedPacket.MarshalBinary()
	if err != nil {
		return err
	}

	return c.Send(&RawPacket{Data: data, Addr: addr})
}

func (c *Client) Listen() error {
	for {
		buffer := make([]byte, 2048)
		read, sender, err := c.conn.ReadFromUDP(buffer)
		/*if isNetErrClosing(err) {
			close(t.stopChan)
			return nil
		} else*/if err != nil {
			fmt.Printf("udp read error: %s\n", err.Error())
			return err
		}
		if read < 1 {
			continue
		}

		packet := RawPacket{
			Addr: sender,
			Data: buffer[:read],
		}

		if err := c.handle(&packet); err != nil {
			fmt.Printf("error handling packet: %s\n", err)
		}
	}
}

func (c *Client) handle(rawPacket *RawPacket) error {
	cryptoPacket := CryptoPacket{}
	if err := cryptoPacket.UnmarshalBinary(rawPacket.Data); err != nil {
		return err
	}

	packet, err := c.DecryptPacket(cryptoPacket)
	if err != nil {
		return err
	}

	switch p := packet.(type) {
	case *HandshakeRequestPacket:
		fmt.Println("handshake request received")
		if err := c.HandshakeResponse(rawPacket.Addr, cryptoPacket.PublicKey, p.Blob); err != nil {
			fmt.Printf("unable to send handshake response: %s\n", err)
		}
	case *HandshakeResponsePacket:
		fmt.Println("handshake response received")
		if (c.lastBlob != nil) && bytes.Equal(p.Blob[:], c.lastBlob[:]) {
			fmt.Println("handshake confirmed")
			break
		}
		fmt.Println("bad handshake")
	default:
		return errors.New("wtf")
	}

	return nil
}

func (c *Client) Print() {
	fmt.Printf("listening on %s\n", c.conn.LocalAddr())
	fmt.Printf("public key: %s\n", hex.EncodeToString(c.prof.PublicKey[:]))
}

func (c *Client) Stop() error {
	return nil
}
