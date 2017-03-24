package client

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/Impyy/kci-demo/cmd/kci-demo/crypto"
)

var (
	order = binary.BigEndian

	packetIDHandshakeRequest  byte = 0x00
	packetIDHandshakeResponse byte = 0x01
)

type RawPacket struct {
	Data []byte
	Addr *net.UDPAddr
}

type CryptoPacket struct {
	Type      byte
	PublicKey *[crypto.KeySize]byte
	Nonce     *[crypto.NonceSize]byte
	Data      []byte
}

type HandshakeRequestPacket struct {
	Blob *[crypto.NonceSize]byte
}

type HandshakeResponsePacket struct {
	Blob *[crypto.NonceSize]byte
}

type Marshalable interface {
	MarshalBinary() ([]byte, error)
	UnmarshalBinary([]byte) error
}

type Packet interface {
	ID() byte
	Marshalable
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (p *CryptoPacket) MarshalBinary() ([]byte, error) {
	buff := new(bytes.Buffer)

	err := binary.Write(buff, order, p.Type)
	if err != nil {
		return nil, err
	}

	_, err = buff.Write(p.PublicKey[:])
	if err != nil {
		return nil, err
	}

	_, err = buff.Write(p.Nonce[:])
	if err != nil {
		return nil, err
	}

	_, err = buff.Write(p.Data)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (p *CryptoPacket) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)

	err := binary.Read(reader, binary.BigEndian, &p.Type)
	if err != nil {
		return err
	}

	p.PublicKey = new([crypto.KeySize]byte)
	_, err = reader.Read(p.PublicKey[:])
	if err != nil {
		return err
	}

	p.Nonce = new([crypto.NonceSize]byte)
	_, err = reader.Read(p.Nonce[:])
	if err != nil {
		return err
	}

	p.Data = make([]byte, reader.Len())
	_, err = reader.Read(p.Data)
	return err
}

func (p HandshakeRequestPacket) ID() byte {
	return packetIDHandshakeRequest
}

func (p *HandshakeRequestPacket) MarshalBinary() ([]byte, error) {
	buff := new(bytes.Buffer)

	_, err := buff.Write(p.Blob[:])
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (p *HandshakeRequestPacket) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)

	p.Blob = new([crypto.NonceSize]byte)
	_, err := reader.Read(p.Blob[:])
	return err
}

func (p HandshakeResponsePacket) ID() byte {
	return packetIDHandshakeResponse
}

func (p *HandshakeResponsePacket) MarshalBinary() ([]byte, error) {
	p2 := HandshakeRequestPacket{Blob: p.Blob}
	return p2.MarshalBinary()
}

func (p *HandshakeResponsePacket) UnmarshalBinary(data []byte) error {
	p2 := HandshakeRequestPacket{}
	if err := p2.UnmarshalBinary(data); err != nil {
		return err
	}

	p.Blob = p2.Blob
	return nil
}
