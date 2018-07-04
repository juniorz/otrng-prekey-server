package prekeyserver

import (
	"time"

	. "gopkg.in/check.v1"
)

// One: test check storage flow to check the simplest thing with DAKE
// Two: do retrieval flow with failure, to check simplest non-DAKE flow

func generateSitaClientProfile(longTerm *keypair) *clientProfile {
	sita := &clientProfile{}
	sita.identifier = 0xAABBCCDD
	sita.instanceTag = 0x1245ABCD
	sita.publicKey = longTerm.pub
	sita.versions = []byte{0x04}
	sita.expiration = time.Date(2028, 11, 5, 13, 46, 00, 13, time.UTC)
	sita.dsaKey = nil
	sita.transitionalSignature = nil
	// This eddsa signature is NOT correct, since we have no way of generating proper eddsa signatures at the moment.
	// This will all have to wait
	sita.sig = &eddsaSignature{s: [114]byte{0x01, 0x02, 0x03}}
	return sita
}

func generateSitaIPoint() *keypair {
	return deriveECDHKeypair([symKeyLength]byte{0x42, 0x11, 0xCC, 0x22, 0xDD, 0x11, 0xFF})
}

type testData struct {
	instanceTag   uint32
	longTerm      *keypair
	clientProfile *clientProfile
	i             *keypair
}

func generateSitaTestData() *testData {
	t := &testData{}
	t.instanceTag = 0x1245ABCD
	t.longTerm = deriveEDDSAKeypair([symKeyLength]byte{0x42, 0x00, 0x00, 0x55, 0x55, 0x00, 0x00, 0x55})
	t.clientProfile = generateSitaClientProfile(t.longTerm)
	t.i = generateSitaIPoint()
	return t
}

var sita = generateSitaTestData()

func (s *GenericServerSuite) Test_flow_CheckStorageNumber(c *C) {
	serverKey := deriveEDDSAKeypair([symKeyLength]byte{0x25, 0x25, 0x25, 0x25, 0x25, 0x25, 0x25, 0x25})
	gs := &GenericServer{
		identity:    "masterOfKeys.example.org",
		rand:        fixtureRand(),
		key:         serverKey,
		fingerprint: serverKey.pub.fingerprint(),
	}
	mh := &otrngMessageHandler{s: gs}

	d1 := generateDake1(sita.instanceTag, sita.clientProfile, sita.i.pub.k)

	r, e := mh.handleMessage("sita@example.org", d1.serialize())

	c.Assert(e, IsNil)

	d2 := dake2Message{}
	_, ok := d2.deserialize(r)

	c.Assert(ok, Equals, true)
	c.Assert(d2.instanceTag, Equals, uint32(0x1245ABCD))
	c.Assert(d2.serverIdentity, DeepEquals, []byte("masterOfKeys.example.org"))
	c.Assert(d2.serverFingerprint[:], DeepEquals, []byte{
		0x32, 0x7c, 0xd2, 0xfc, 0xcb, 0x3b, 0xd1, 0x1d,
		0x63, 0x6a, 0x33, 0x44, 0xd5, 0x4b, 0xc9, 0xd,
		0x8d, 0x7e, 0xf3, 0x38, 0x39, 0x1e, 0x9d, 0x21,
		0x1f, 0x66, 0x39, 0x61, 0xd0, 0xf7, 0xea, 0x4,
		0xf0, 0x12, 0xd0, 0x76, 0xe3, 0x5a, 0x9c, 0x7a,
		0xe7, 0x37, 0xfd, 0xd8, 0xab, 0x1e, 0x3e, 0xf1,
		0xbd, 0x66, 0x57, 0xa2, 0x71, 0x1, 0xf7, 0x4e,
	})
	c.Assert(d2.s.DSAEncode(), DeepEquals, []byte{
		0xd9, 0xe9, 0xed, 0x15, 0xf1, 0x57, 0x6f, 0x39,
		0x80, 0xa4, 0x57, 0xa0, 0x3c, 0xc5, 0x9, 0xec,
		0xa0, 0x13, 0x90, 0x57, 0xfc, 0xb, 0x33, 0x36,
		0x55, 0x17, 0xf, 0x7f, 0x34, 0x8e, 0xe1, 0x15,
		0x19, 0xdc, 0x86, 0x2f, 0x82, 0xb, 0x3a, 0xe,
		0x42, 0x9, 0xc3, 0xdb, 0xd0, 0x5b, 0x93, 0x19,
		0x2c, 0x39, 0x96, 0x2a, 0x51, 0xfe, 0x58, 0xf9,
		0x0})

	c.Assert(d2.sigma.c1.Encode(), DeepEquals, []byte{
		0x3e, 0x4f, 0x9a, 0xe1, 0x98, 0x28, 0x67, 0x86,
		0xf1, 0xba, 0x33, 0x60, 0x31, 0x54, 0x50, 0x49,
		0x5, 0xfa, 0xc0, 0x93, 0xf5, 0x5d, 0x64, 0xca,
		0x22, 0x8d, 0x27, 0x22, 0x6c, 0xf6, 0x59, 0xd9,
		0xb3, 0x31, 0x31, 0x73, 0x10, 0xb4, 0x6e, 0xc6,
		0x17, 0xba, 0x5f, 0x91, 0xdd, 0x31, 0xb5, 0x9,
		0x83, 0x1, 0x51, 0x7c, 0x8, 0x2e, 0x1c, 0x33})

	c.Assert(d2.sigma.r1.Encode(), DeepEquals, []byte{
		0x47, 0x71, 0x5b, 0x81, 0xa8, 0x56, 0x47, 0x16,
		0x5, 0x8f, 0x9a, 0x2e, 0x9b, 0x2c, 0x55, 0xc3,
		0xd7, 0x0, 0xd3, 0x26, 0x13, 0xf5, 0x93, 0xe4,
		0xf4, 0xcb, 0x98, 0xb7, 0xe7, 0x81, 0xd, 0x35,
		0xa7, 0xa5, 0x59, 0x74, 0x9b, 0x7d, 0x19, 0x63,
		0x20, 0x5c, 0x1, 0x3b, 0x79, 0x70, 0x35, 0x33,
		0xfa, 0x1f, 0x38, 0xe3, 0x81, 0x96, 0x78, 0x2e})

	c.Assert(d2.sigma.c2.Encode(), DeepEquals, []byte{
		0x9f, 0xf2, 0x1d, 0x1b, 0xa7, 0x11, 0x1, 0x31,
		0x19, 0xe7, 0xdd, 0x7, 0xcf, 0x3b, 0xbd, 0xf7,
		0xe9, 0x9e, 0x34, 0x81, 0x28, 0x3e, 0xfe, 0x96,
		0xa7, 0xc1, 0x93, 0xd3, 0xc9, 0x8d, 0x6, 0xba,
		0xdd, 0x9b, 0xcc, 0x44, 0x4f, 0x64, 0xe, 0xa,
		0xd2, 0xfe, 0x4c, 0x54, 0xf0, 0xcc, 0x9c, 0x9d,
		0x86, 0xe7, 0x6c, 0x2a, 0x85, 0x63, 0x91, 0xc})

	c.Assert(d2.sigma.r2.Encode(), DeepEquals, []byte{
		0xa4, 0x62, 0x9b, 0xd8, 0x83, 0x5e, 0x41, 0x7a,
		0x5b, 0x6d, 0xaa, 0xd0, 0xf7, 0xa9, 0x61, 0xa5,
		0xe, 0x66, 0xe3, 0xc4, 0x4c, 0xe3, 0xc1, 0xa9,
		0xe2, 0xe5, 0xfb, 0xaa, 0x9f, 0x2d, 0xc, 0xbf,
		0x18, 0xaf, 0x86, 0xe1, 0xb4, 0xa1, 0x83, 0x48,
		0x25, 0xca, 0x1e, 0x6, 0x6e, 0x82, 0x9b, 0x2f,
		0x22, 0xfc, 0xf, 0x80, 0x9d, 0x9c, 0x90, 0x1c})

	c.Assert(d2.sigma.c3.Encode(), DeepEquals, []byte{
		0x31, 0xb3, 0xc2, 0xa1, 0x10, 0x46, 0x2d, 0xd2,
		0x4a, 0x3c, 0x4d, 0x8c, 0x2c, 0xba, 0xd4, 0xe3,
		0x6e, 0x73, 0xfb, 0x8, 0x1f, 0x92, 0xb4, 0x88,
		0x85, 0x50, 0xd, 0xe4, 0x26, 0x9a, 0x3b, 0x86,
		0x94, 0x5a, 0xf3, 0x33, 0xb3, 0x95, 0x10, 0x6f,
		0x54, 0x6c, 0x14, 0xde, 0x51, 0x97, 0x14, 0x86,
		0x1f, 0xb0, 0x27, 0xdf, 0x57, 0x48, 0x7c, 0x3f})

	c.Assert(d2.sigma.r3.Encode(), DeepEquals, []byte{
		0x75, 0x8f, 0x53, 0x1c, 0x7b, 0x2f, 0x4, 0xbf,
		0x34, 0x16, 0xf0, 0x8e, 0x7, 0x19, 0x53, 0x9f,
		0x9c, 0xab, 0xcd, 0xab, 0xfa, 0x5f, 0x3a, 0xe3,
		0x55, 0xf5, 0x85, 0xbd, 0x3e, 0x46, 0x8b, 0xe,
		0xfb, 0x1a, 0xc, 0x1f, 0xa, 0xe3, 0x9e, 0x1e,
		0x93, 0x4a, 0x86, 0x95, 0x4c, 0x7, 0x0, 0xda,
		0xee, 0xd2, 0x8c, 0x4, 0xc0, 0x57, 0x71, 0x28})

	// Generate a DAKE3 + storage request
	// Send the DAKE3 to the message handler
	// Check that the returned storage information message is correct
}
