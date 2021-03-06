package prekeyserver

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
)

type publicationMessage struct {
	prekeyMessages []*prekeyMessage
	clientProfile  *clientProfile
	prekeyProfile  *prekeyProfile
	mac            [macLength]byte
}

type storageInformationRequestMessage struct {
	mac [macLength]byte
}

type storageStatusMessage struct {
	instanceTag uint32
	number      uint32
	mac         [macLength]byte
}

type successMessage struct {
	instanceTag uint32
	mac         [macLength]byte
}

type failureMessage struct {
	instanceTag uint32
	mac         [macLength]byte
}

type ensembleRetrievalQueryMessage struct {
	instanceTag uint32
	identity    string
	versions    []byte
}

type ensembleRetrievalMessage struct {
	instanceTag uint32
	ensembles   []*prekeyEnsemble
}

type noPrekeyEnsemblesMessage struct {
	instanceTag uint32
	message     string
}

type message interface {
	serializable
	validate(string, *GenericServer) error
	respond(string, *GenericServer) (serializable, error)
}

func parseVersion(message []byte) uint16 {
	_, v, _ := extractShort(message)
	return v
}

func parseMessage(msg []byte) (message, uint8, error) {
	if len(msg) <= indexOfMessageType {
		return nil, 0, errors.New("message too short to be a valid message")
	}

	if v := parseVersion(msg); v != uint16(4) {
		return nil, 0, errors.New("invalid protocol version")
	}

	messageType := msg[indexOfMessageType]

	var r message
	switch messageType {
	case messageTypeDAKE1:
		r = &dake1Message{}
	case messageTypeDAKE3:
		r = &dake3Message{}
	case messageTypePublication:
		r = &publicationMessage{}
	case messageTypeStorageInformationRequest:
		r = &storageInformationRequestMessage{}
	case messageTypeEnsembleRetrievalQuery:
		r = &ensembleRetrievalQueryMessage{}
	default:
		return nil, 0, fmt.Errorf("unknown message type: 0x%x", messageType)
	}

	r.deserialize(msg)

	return r, messageType, nil
}

func generateStorageInformationRequestMessage(macKey []byte) *storageInformationRequestMessage {
	mac := kdfx(usageStorageInfoMAC, 64, macKey, []byte{messageTypeStorageInformationRequest})
	res := &storageInformationRequestMessage{}
	copy(res.mac[:], mac)
	return res
}

func (m *storageInformationRequestMessage) respond(from string, s *GenericServer) (serializable, error) {
	ses := s.session(from)
	num := s.storage().numberStored(from, ses.instanceTag())
	itag := ses.instanceTag()
	prekeyMacK := ses.macKey()
	statusMac := kdfx(usageStatusMAC, 64, prekeyMacK, []byte{messageTypeStorageStatusMessage}, serializeWord(itag), serializeWord(num))

	ret := &storageStatusMessage{
		instanceTag: itag,
		number:      num,
	}
	copy(ret.mac[:], statusMac)

	return ret, nil
}

func (m *storageInformationRequestMessage) validate(from string, s *GenericServer) error {
	prekeyMacK := s.session(from).macKey()
	tag := kdfx(usageStorageInfoMAC, 64, prekeyMacK, []byte{messageTypeStorageInformationRequest})
	if !bytes.Equal(tag, m.mac[:]) {
		return errors.New("incorrect MAC")
	}

	return nil
}

func (m *ensembleRetrievalQueryMessage) validate(from string, s *GenericServer) error {
	return nil
}

func (m *ensembleRetrievalQueryMessage) respond(from string, s *GenericServer) (serializable, error) {
	stor := s.storage()
	bundles := stor.retrieveFor(m.identity)
	if len(bundles) == 0 {
		return &noPrekeyEnsemblesMessage{
			instanceTag: m.instanceTag,
			message:     noPrekeyMessagesAvailableMessage,
		}, nil
	}
	return &ensembleRetrievalMessage{
		instanceTag: m.instanceTag,
		ensembles:   bundles,
	}, nil
}

func generateMACForPublicationMessage(cp *clientProfile, pp *prekeyProfile, pms []*prekeyMessage, macKey []byte) []byte {
	kpms := kdfx(usagePrekeyMessage, 64, serializePrekeyMessages(pms))
	kpps := kdfx(usagePrekeyProfile, 64, pp.serialize())
	k := []byte{byte(0)}
	kcp := []byte{}
	if cp != nil {
		k = []byte{1}
		kcp = kdfx(usageClientProfile, 64, cp.serialize())
	}

	ppLen := 0
	if pp != nil {
		ppLen = 1
	}

	fmt.Printf("\n\n\nGenerating SIG for publication message\n")
	fmt.Printf("prekeyMsgKDF = %s\n", base64.StdEncoding.EncodeToString(kpms))
	fmt.Printf("prekeyProfileKDF = %s\n", base64.StdEncoding.EncodeToString(kpps))
	fmt.Printf("clientProfileKDF = %s\n", base64.StdEncoding.EncodeToString(kcp))
	fmt.Printf("fullKDFBody = %s\n", base64.StdEncoding.EncodeToString(concat(macKey, []byte{messageTypePublication, byte(len(pms))}, kpms, k, kcp, []byte{byte(ppLen)}, kpps)))

	return kdfx(usagePreMAC, 64, concat(macKey, []byte{messageTypePublication, byte(len(pms))}, kpms, k, kcp, []byte{byte(ppLen)}, kpps))
}

func generatePublicationMessage(cp *clientProfile, pp *prekeyProfile, pms []*prekeyMessage, macKey []byte) *publicationMessage {
	mac := generateMACForPublicationMessage(cp, pp, pms, macKey)
	pm := &publicationMessage{
		prekeyMessages: pms,
		clientProfile:  cp,
		prekeyProfile:  pp,
	}
	copy(pm.mac[:], mac)
	return pm
}

func (m *publicationMessage) validate(from string, s *GenericServer) error {
	macKey := s.session(from).macKey()
	clientProfile := s.session(from).clientProfile()
	mac := generateMACForPublicationMessage(m.clientProfile, m.prekeyProfile, m.prekeyMessages, macKey)

	fmt.Printf("mac = %s\n", base64.StdEncoding.EncodeToString(mac))
	if !bytes.Equal(mac[:], m.mac[:]) {
		return errors.New("invalid mac for publication message")
	}

	tag := s.session(from).instanceTag()
	if m.clientProfile != nil && m.clientProfile.validate(tag) != nil {
		return errors.New("invalid client profile in publication message")
	}

	if m.prekeyProfile != nil && m.prekeyProfile.validate(tag, clientProfile.publicKey) != nil {
		return errors.New("invalid prekey profile in publication message")
	}

	for _, pm := range m.prekeyMessages {
		if pm.validate(tag) != nil {
			return errors.New("invalid prekey message in publication message")
		}
	}

	return nil
}

func generateSuccessMessage(macKey []byte, tag uint32) *successMessage {
	m := &successMessage{
		instanceTag: tag,
	}

	mac := kdfx(usageSuccessMAC, 64, appendWord(append(macKey, messageTypeSuccess), tag))
	copy(m.mac[:], mac)

	return m
}

func (m *publicationMessage) respond(from string, s *GenericServer) (serializable, error) {
	stor := s.storage()
	stor.storeClientProfile(from, m.clientProfile)
	stor.storePrekeyProfile(from, m.prekeyProfile)
	stor.storePrekeyMessages(from, m.prekeyMessages)

	macKey := s.session(from).macKey()
	instanceTag := s.session(from).instanceTag()

	s.sessionComplete(from)

	return generateSuccessMessage(macKey, instanceTag), nil
}
