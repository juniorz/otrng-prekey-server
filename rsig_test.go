package prekeyserver

import (
	. "gopkg.in/check.v1"
)

func (s *GenericServerSuite) Test_generateSignature_generatesACorrectSignature(c *C) {
	msg := []byte("hi")
	p1 := deriveEDDSAKeypair([symKeyLength]byte{0x0A})
	p2 := deriveEDDSAKeypair([symKeyLength]byte{0x19})
	p3 := deriveEDDSAKeypair([symKeyLength]byte{0x26})
	wr := fixedRandBytes(
		// for t1
		[]byte{
			0x42, 0x07, 0x96, 0xa4, 0xac, 0x11, 0xc8, 0xb5,
			0x17, 0xd1, 0x58, 0xa3, 0xd4, 0xdb, 0xc9, 0xf2,
			0x7f, 0x38, 0xa4, 0x8e, 0x4e, 0xba, 0xa4, 0xae,
			0x67, 0xf3, 0x9c, 0x0d, 0x35, 0xe4, 0x09, 0xd7,
			0x52, 0xb2, 0x82, 0xd1, 0x92, 0x5d, 0xcc, 0xbe,
			0x7a, 0xb4, 0xdb, 0xa1, 0x98, 0x75, 0xf6, 0xc3,
			0x00, 0x93, 0xb6, 0xac, 0xec, 0x49, 0x23, 0x3e,
			0xd0,
		},
		// for t2
		[]byte{
			0xbe, 0x7b, 0x3a, 0xe3, 0xd3, 0x58, 0xcb, 0xe6,
			0x67, 0x78, 0xe5, 0xe6, 0xa5, 0xd0, 0x93, 0x2b,
			0xa9, 0x1f, 0x83, 0xfa, 0x86, 0x7a, 0x1b, 0xdd,
			0x9f, 0x8c, 0x6c, 0x59, 0xee, 0x95, 0x72, 0x46,
			0xde, 0x15, 0xba, 0x1c, 0x23, 0x5b, 0x38, 0xf1,
			0x75, 0x9d, 0x14, 0x2b, 0xd3, 0xb0, 0xa5, 0xe7,
			0xfd, 0xc8, 0x34, 0x79, 0xf2, 0x52, 0xdd, 0x13,
			0x29,
		},
		// for t3
		[]byte{
			0x91, 0xf6, 0x13, 0x5e, 0x79, 0x97, 0x4a, 0x1e,
			0x9d, 0x46, 0x0e, 0x87, 0xf5, 0xd5, 0xba, 0x88,
			0xec, 0x7f, 0x60, 0x49, 0x44, 0x8f, 0x0e, 0x64,
			0xb1, 0x8b, 0x22, 0x52, 0x8d, 0x4b, 0x23, 0x81,
			0xc3, 0x54, 0xcc, 0x0c, 0xb0, 0xbb, 0x1b, 0x88,
			0x68, 0x54, 0x21, 0x22, 0xbd, 0x7c, 0xea, 0xa8,
			0x14, 0x5f, 0x19, 0xda, 0xeb, 0x5b, 0x71, 0xf7,
			0x25,
		},
		// for r1
		[]byte{
			0xa9, 0x63, 0x8e, 0x3b, 0xec, 0x80, 0xc3, 0x34,
			0x7f, 0x61, 0xe3, 0x39, 0xb5, 0x86, 0xa8, 0xfa,
			0x87, 0x77, 0xf1, 0xee, 0x5c, 0xec, 0x3e, 0xbe,
			0x86, 0xde, 0x0c, 0x08, 0x95, 0x2b, 0x98, 0x7d,
			0xc1, 0x8f, 0x34, 0x97, 0x99, 0x51, 0x60, 0x87,
			0x51, 0xbb, 0x7f, 0x8f, 0x07, 0x5e, 0x82, 0xc0,
			0x53, 0x5c, 0x93, 0xc3, 0x16, 0xe6, 0x13, 0xee,
			0x2c,
		},
		// for r2
		[]byte{
			0xc1, 0x8d, 0xd5, 0x82, 0x87, 0x25, 0xc9, 0xef,
			0x0e, 0xc8, 0xe6, 0x17, 0x81, 0x19, 0x37, 0x08,
			0xde, 0x7e, 0xe3, 0x9c, 0xac, 0x2c, 0x3a, 0xee,
			0xa1, 0xf5, 0xfe, 0xab, 0x23, 0x47, 0x7c, 0x6c,
			0xe6, 0x0e, 0xdb, 0x6b, 0x3d, 0xd0, 0x0e, 0x83,
			0x84, 0x30, 0xe7, 0x7d, 0x60, 0xda, 0x2d, 0xf9,
			0x56, 0x0b, 0xb8, 0xf0, 0x60, 0x5c, 0xdd, 0xd0,
			0x2e,
		},
		// for r3
		[]byte{
			0x48, 0xf3, 0x9e, 0x77, 0x0a, 0x50, 0x68, 0xaa,
			0x7d, 0x5b, 0xc6, 0x3f, 0xe3, 0x5a, 0x59, 0x48,
			0xa5, 0x0d, 0x4c, 0x34, 0x7a, 0xf5, 0xf2, 0x53,
			0x49, 0x1e, 0xe8, 0x8e, 0x46, 0xa8, 0xb8, 0xf0,
			0x6a, 0x9c, 0x17, 0x81, 0x65, 0x40, 0x7a, 0x18,
			0x1e, 0xb6, 0x7e, 0x40, 0x9c, 0xc5, 0xd0, 0x8c,
			0x25, 0xdf, 0x7a, 0x8b, 0x58, 0xc4, 0x45, 0x19,
			0x0e,
		},
		// for c1
		[]byte{
			0x36, 0x8f, 0x5e, 0xd7, 0x69, 0x5e, 0x8d, 0x3c,
			0xae, 0x70, 0xdf, 0x64, 0x3b, 0x8d, 0xd3, 0xaf,
			0xe8, 0x59, 0xb6, 0x18, 0x91, 0xb6, 0x29, 0x7b,
			0x86, 0x7f, 0x26, 0x54, 0x06, 0x14, 0x98, 0x56,
			0x8e, 0xb2, 0xfb, 0x6f, 0xdd, 0xdf, 0x82, 0x26,
			0xeb, 0x63, 0x9d, 0xe9, 0x2e, 0x3d, 0x9c, 0x58,
			0x9a, 0xa0, 0xf3, 0xa2, 0xb1, 0x83, 0xd7, 0x9a,
			0x89,
		},
		// for c2
		[]byte{
			0x96, 0xd5, 0x74, 0xfe, 0xe3, 0xcf, 0x97, 0x16,
			0x85, 0x28, 0xf7, 0x92, 0xcc, 0x51, 0xfd, 0x3d,
			0x30, 0x4d, 0xf8, 0xed, 0x4c, 0x5d, 0x17, 0x24,
			0x5d, 0x14, 0xf3, 0xdb, 0x96, 0xcf, 0x0a, 0x48,
			0x88, 0x54, 0x4e, 0xa3, 0x0f, 0x5c, 0x9b, 0x6b,
			0xa2, 0xe2, 0x4e, 0x98, 0x85, 0x80, 0x9f, 0x75,
			0xb6, 0x77, 0x0e, 0xaf, 0xdd, 0xe4, 0x84, 0xf0,
			0x76,
		},
		// for c3
		[]byte{0x16, 0x58, 0xa8, 0xaa, 0x2e, 0x2c, 0x8e, 0x98,
			0xc9, 0xe7, 0xea, 0xa5, 0x7a, 0xea, 0x66, 0x85,
			0xcb, 0xf8, 0xb4, 0x10, 0xca, 0x6c, 0xfa, 0x0f,
			0xbb, 0x02, 0xa9, 0x18, 0xfb, 0x66, 0xec, 0x5c,
			0xdd, 0x5f, 0xd3, 0x6c, 0xb9, 0x81, 0xf1, 0xde,
			0xa8, 0x94, 0x7f, 0x88, 0x00, 0xb6, 0xa1, 0xd9,
			0x03, 0xd9, 0x4d, 0x1c, 0x18, 0x00, 0xfd, 0xae,
			0xf3,
		},
	)

	rsig, _ := generateSignature(wr, p1.priv, p1.pub, p1.pub, p2.pub, p3.pub, msg)

	c.Assert(rsig.c1.Encode(), DeepEquals, []byte{
		0x0b, 0x45, 0x28, 0x36, 0xa4, 0x1e, 0xcd, 0x11,
		0x1c, 0xd2, 0x59, 0x98, 0x15, 0x8b, 0xa6, 0xba,
		0xa8, 0x85, 0x08, 0xf9, 0x4a, 0x51, 0xbd, 0x95,
		0x80, 0x3c, 0x72, 0x74, 0x64, 0x98, 0x11, 0x2a,
		0x4e, 0xa7, 0x91, 0x8a, 0xe2, 0x32, 0x1e, 0xff,
		0xd1, 0xf2, 0x3b, 0xef, 0x9c, 0x12, 0x96, 0x36,
		0xc1, 0x0c, 0x5b, 0x40, 0x0f, 0x80, 0xc8, 0x2a,
	})
	c.Assert(rsig.r1.Encode(), DeepEquals, []byte{
		0x7f, 0xa2, 0xdc, 0x8b, 0x51, 0x3e, 0xf6, 0x87,
		0xd6, 0xbf, 0x24, 0x89, 0x20, 0xcf, 0x76, 0xc9,
		0x3e, 0x76, 0x05, 0x54, 0x9a, 0xf4, 0xef, 0xa7,
		0xe5, 0x1c, 0x3f, 0x92, 0xb8, 0x92, 0xf6, 0xce,
		0xe2, 0xe8, 0xaa, 0xb8, 0x62, 0x8b, 0x15, 0x8b,
		0x3c, 0xf7, 0x7d, 0x8b, 0x1f, 0xd0, 0x2d, 0x72,
		0xe7, 0x76, 0x16, 0x9d, 0x3f, 0xcc, 0x0d, 0x2d,
	})
	c.Assert(rsig.c2.Encode(), DeepEquals, []byte{
		0x6c, 0x99, 0x30, 0x35, 0x65, 0x11, 0x93, 0x24,
		0x3a, 0x4b, 0x9d, 0xbe, 0xb3, 0xbd, 0x0c, 0x7f,
		0xc8, 0x28, 0x61, 0x20, 0x57, 0x3a, 0x85, 0x05,
		0xa3, 0x15, 0xcb, 0x42, 0x4f, 0xab, 0xfb, 0xce,
		0x6b, 0x41, 0x26, 0x7c, 0xb3, 0x6a, 0x87, 0xfe,
		0x8e, 0x82, 0x36, 0x4f, 0xf8, 0x0f, 0x79, 0x45,
		0xb0, 0xf0, 0x45, 0xe7, 0xed, 0x9a, 0x17, 0x3b,
	})
	c.Assert(rsig.r2.Encode(), DeepEquals, []byte{
		0x1d, 0xd4, 0xf9, 0x79, 0xd8, 0x95, 0xca, 0xca,
		0x09, 0x66, 0x26, 0x90, 0x73, 0xf6, 0x4f, 0xbf,
		0x3e, 0x20, 0xc3, 0x61, 0x03, 0x34, 0x03, 0x3b,
		0x16, 0x35, 0x8d, 0xb7, 0x35, 0xf3, 0x77, 0x21,
		0xa1, 0x32, 0x0e, 0xc4, 0x27, 0xf1, 0x07, 0x56,
		0x5f, 0xa0, 0xba, 0x1c, 0xa0, 0x3f, 0xd3, 0x8f,
		0xd9, 0x91, 0x8b, 0xd3, 0xa2, 0xad, 0xfd, 0x3d,
	})
	c.Assert(rsig.c3.Encode(), DeepEquals, []byte{
		0x96, 0x20, 0x5a, 0xd3, 0x1f, 0x1d, 0x60, 0xbb,
		0xe4, 0x08, 0x27, 0x4d, 0xf2, 0x4e, 0x91, 0xa7,
		0xed, 0xdc, 0x19, 0x07, 0x2c, 0x73, 0x67, 0xf8,
		0x27, 0x7d, 0x83, 0x13, 0x5e, 0x49, 0x42, 0x92,
		0xfa, 0x92, 0x96, 0xc0, 0x14, 0xb8, 0x23, 0xed,
		0xff, 0x4f, 0xf9, 0xe0, 0x65, 0x06, 0x54, 0xf4,
		0x79, 0x73, 0xfe, 0x6e, 0x31, 0x70, 0x7f, 0x39,
	})
	c.Assert(rsig.r3.Encode(), DeepEquals, []byte{
		0x99, 0xca, 0xa2, 0x6c, 0xa5, 0xf9, 0x8d, 0x8c,
		0x03, 0x9c, 0x96, 0xab, 0x6b, 0x8b, 0xe7, 0xe9,
		0xe4, 0xa2, 0xb2, 0x8d, 0xbe, 0xfe, 0x71, 0x26,
		0x43, 0x5d, 0x00, 0xb5, 0xbb, 0x0a, 0xaa, 0x55,
		0x7b, 0xb2, 0x66, 0x70, 0xc5, 0xbb, 0xe7, 0x9c,
		0xd9, 0xe8, 0x83, 0x23, 0xf8, 0x38, 0x04, 0x81,
		0xda, 0x1b, 0x03, 0xa2, 0xd8, 0xc6, 0x7b, 0x3f,
	})
}
