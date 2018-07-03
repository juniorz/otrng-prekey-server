package prekeyserver

import (
	"github.com/twstrike/ed448"
)

type ringSignature struct {
	c1 ed448.Scalar
	r1 ed448.Scalar
	c2 ed448.Scalar
	r2 ed448.Scalar
	c3 ed448.Scalar
	r3 ed448.Scalar
}

func generateZqKeypair(r WithRandom) (ed448.Scalar, ed448.Point) {
	sym := [symKeyLength]byte{}
	randomInto(r, sym[:])
	return deriveZqKeypair(sym)
}

func deriveZqKeypair(sym [symKeyLength]byte) (ed448.Scalar, ed448.Point) {
	kp := deriveEDDSAKeypair(sym)
	return kp.priv.k, kp.pub.k
}

func choose_T(Ai ed448.Point, isSecret uint32, Ri ed448.Point, Ti ed448.Point, ci ed448.Scalar) ed448.Point {
	chosen := ed448.PointScalarMul(Ai, ci)
	chosen.Add(Ri, chosen)

	return ed448.ConstantTimeSelectPoint(chosen, Ti, isSecret)
}

var base_point_bytes_dup = []byte{
	0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
	0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
	0x66, 0x66, 0x66, 0x66, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
	0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
	0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x00,
}

var prime_order_bytes_dup = []byte{
	0x3f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0x7c, 0xca, 0x23, 0xe9, 0xc4, 0x4e, 0xdb, 0x49,
	0xae, 0xd6, 0x36, 0x90, 0x21, 0x6c, 0xc2, 0x72, 0x8d, 0xc5, 0x8f, 0x55,
	0x23, 0x78, 0xc2, 0x92, 0xab, 0x58, 0x44, 0xf3,
}

func calculateC(A1, A2, A3, T1, T2, T3 ed448.Point, msg []byte) ed448.Scalar {
	h := kdfx_otrv4(0x1D, 64,
		base_point_bytes_dup,
		prime_order_bytes_dup,
		A1.DSAEncode(),
		A2.DSAEncode(),
		A3.DSAEncode(),
		T1.DSAEncode(),
		T2.DSAEncode(),
		T3.DSAEncode(),
		msg,
	)

	return ed448.NewScalar(h)
}

func calculateCI(c, ci ed448.Scalar, isSecret uint32, cj, ck ed448.Scalar) ed448.Scalar {
	ifSecret := ed448.NewScalar()
	ifSecret.Sub(c, cj)
	ifSecret.Sub(ifSecret, ck)

	return ed448.ConstantTimeSelectScalar(ci, ifSecret, isSecret)
}

func calculateRI(secret, ri ed448.Scalar, isSecret uint32, ci, ti ed448.Scalar) ed448.Scalar {
	ifSecret := ed448.NewScalar()
	ifSecret.Mul(ci, secret)
	ifSecret.Sub(ti, ifSecret)

	return ed448.ConstantTimeSelectScalar(ri, ifSecret, isSecret)
}

func generateSignature(wr WithRandom, secret *privateKey, pub *publicKey, A1, A2, A3 *publicKey, m []byte) (*ringSignature, error) {
	r := &ringSignature{}

	is_A1 := pub.k.EqualsMask(A1.k)
	is_A2 := pub.k.EqualsMask(A2.k)
	is_A3 := pub.k.EqualsMask(A3.k)

	t1, T1 := generateZqKeypair(wr)
	t2, T2 := generateZqKeypair(wr)
	t3, T3 := generateZqKeypair(wr)

	r1, R1 := generateZqKeypair(wr)
	r2, R2 := generateZqKeypair(wr)
	r3, R3 := generateZqKeypair(wr)

	c1, _ := generateZqKeypair(wr)
	c2, _ := generateZqKeypair(wr)
	c3, _ := generateZqKeypair(wr)

	chosen_T1 := choose_T(A1.k, is_A1, R1, T1, c1)
	chosen_T2 := choose_T(A2.k, is_A2, R2, T2, c2)
	chosen_T3 := choose_T(A3.k, is_A3, R3, T3, c3)

	c := calculateC(A1.k, A2.k, A3.k, chosen_T1, chosen_T2, chosen_T3, m)

	r.c1 = calculateCI(c, c1, is_A1, c2, c3)
	r.c2 = calculateCI(c, c2, is_A2, c1, c3)
	r.c3 = calculateCI(c, c3, is_A3, c1, c2)

	r.r1 = calculateRI(secret.k, r1, is_A1, r.c1, t1)
	r.r2 = calculateRI(secret.k, r2, is_A2, r.c2, t2)
	r.r3 = calculateRI(secret.k, r3, is_A3, r.c3, t3)

	return r, nil
}

// verify will actually do the cryptographic validation of the ring signature
func (r *ringSignature) verify() bool {
	// TODO: implement
	return false
}

func (r *ringSignature) deserialize(buf []byte) ([]byte, bool) {
	buf, r.c1, _ = deserializeScalar(buf)
	buf, r.r1, _ = deserializeScalar(buf)
	buf, r.c2, _ = deserializeScalar(buf)
	buf, r.r2, _ = deserializeScalar(buf)
	buf, r.c3, _ = deserializeScalar(buf)
	buf, r.r3, _ = deserializeScalar(buf)
	return buf, true
}

func (r *ringSignature) serialize() []byte {
	var out []byte
	out = append(out, serializeScalar(r.c1)...)
	out = append(out, serializeScalar(r.r1)...)
	out = append(out, serializeScalar(r.c2)...)
	out = append(out, serializeScalar(r.r2)...)
	out = append(out, serializeScalar(r.c3)...)
	out = append(out, serializeScalar(r.r3)...)
	return out
}

func serializeScalar(s ed448.Scalar) []byte {
	return s.Encode()
}

func deserializeScalar(buf []byte) ([]byte, ed448.Scalar, bool) {
	ts := ed448.NewScalar()
	ts.Decode(buf[0:56])
	return buf[56:], ts, true

}
