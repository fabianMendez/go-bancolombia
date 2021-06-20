package main

import (
	"math"
	"math/rand"
	"strconv"
)

type rsaKey struct {
	n     bigInteger
	e     int64
	d     string
	p     string
	q     string
	dmp1  string
	dmq1  string
	coeff string
}

func (k *rsaKey) setPublic(a, b string) {
	i, err := strconv.ParseInt(b, 16, 64)
	if err != nil {
		panic(err)
	}
	k.e = i
	k.n = newBigInteger(a, 16)
}

func (k *rsaKey) doPublic(x bigInteger) bigInteger {
	return x.modPowInt(k.e, k.n)
}

func (k *rsaKey) encrypt(text string) string {
	m := pkcs1pad2(text, (k.n.bitLength()+7)>>3)
	c := k.doPublic(m)
	h := c.toString(16)
	if len(h)&1 == 0 {
		return h
	}
	return "0" + h
}

type bigInteger struct {
	t   int
	s   int64
	DB  int
	DM  int64
	DV  int64
	FV  float64
	F1  int
	F2  int
	arr []int64
	// arr []int64
}

func (bi bigInteger) toString(b int) string {
	if bi.s < 0 {
		return "-" + bi.negate().toString(b)
	}

	var k int

	if b == 16 {
		k = 4
	} else if b == 8 {
		k = 3
	} else if b == 2 {
		k = 1
	} else if b == 32 {
		k = 5
	} else if b == 4 {
		k = 2
	} else {
		// return bi.toRadix(b)
		panic("not implemented")
	}

	km := int64((1 << k) - 1)
	var d int64
	m := false
	r := ""
	i := bi.t

	p := bi.DB - (i*bi.DB)%k

	if i > 0 {
		i--
		if p < bi.DB {
			d = bi.arr[i] >> p
			if d > 0 {
				m = true
				r = int2char(d)
			}
		}

		for i >= 0 {
			if p < k {
				d = (bi.arr[i] & ((1 << p) - 1)) << (k - p)
				i--
				p += bi.DB - k
				d |= bi.arr[i] >> (p)
			} else {
				p -= k
				d = (bi.arr[i] >> (p)) & km
				if p <= 0 {
					p += bi.DB
					i--
				}
			}
			if d > 0 {
				m = true
			}
			if m {
				r += int2char(d)
			}
		}
	}

	if m {
		return r
	}
	return "0"
}

func (bi bigInteger) isEven() bool {
	var x int64
	if bi.t > 0 {
		x = bi.arr[0] & 1
	} else {
		x = bi.s
	}
	return x == 0
}

type something interface {
	convert(bigInteger) bigInteger
	sqrTo(bigInteger, *bigInteger)
	mulTo(*bigInteger, bigInteger, *bigInteger)
	revert(x *bigInteger) *bigInteger
}

type classic struct {
	m bigInteger
}

func (cl classic) convert(_ bigInteger) bigInteger {
	return bigInteger{}
}

func (cl classic) sqrTo(x bigInteger, r *bigInteger) {
	x.squareTo(r)
	cl.reduce(r)
}

func (cl classic) reduce(x *bigInteger) {
	x.divRemTo(cl.m, nil, x)
}

func (cl classic) mulTo(x *bigInteger, y bigInteger, r *bigInteger) {
	x.multiplyTo(y, r)
	cl.reduce(r)
}

func (cl classic) revert(x *bigInteger) *bigInteger {
	r := nbi()
	x.copyTo(r)
	cl.reduce(r)
	return r
}

type montgomery struct {
	m   bigInteger
	mp  int64
	mpl int64
	mph int64
	um  int64
	mt2 int
}

func (mg montgomery) sqrTo(x bigInteger, r *bigInteger) {
	x.squareTo(r)
	mg.reduce(r)
}

func (mg montgomery) mulTo(x *bigInteger, y bigInteger, r *bigInteger) {
	x.multiplyTo(y, r)
	mg.reduce(r)
}
func (mg montgomery) reduce(x *bigInteger) {
	// fmt.Println("const m = parseInt('" + mg.m.toString(16) + "')")
	// fmt.Println("const z = new Montgomery(m)")
	// fmt.Println("const x = parseInt('" + x.toString(16) + "')")
	// fmt.Println("z.reduce(x)")

	for x.t <= mg.mt2 {
		// x.arr[x.t] = 0
		x.arr = append(x.arr, 0)
		x.t++
	}

	// fmt.Println(x.toString(16))

	for i := 0; i < mg.m.t; i++ {
		j := int64(x.arr[i] & 0x7fff)
		u0 := (j*mg.mpl + (((j*mg.mph + (x.arr[i]>>15)*mg.mpl) & mg.um) << 15)) & x.DM
		j = int64(i + mg.m.t)

		x.arr[j] += mg.m.am(0, u0, x, i, 0, mg.m.t)
		for x.arr[j] >= x.DV {
			x.arr[j] -= x.DV
			j++
			x.arr[j]++
		}
	}

	x.clamp()
	x.drShiftTo(mg.m.t, x)
	if x.compareTo(mg.m) >= 0 {
		x.subTo(mg.m, x)
	}
}

func (mg montgomery) revert(x *bigInteger) *bigInteger {
	r := nbi()
	x.copyTo(r)
	mg.reduce(r)
	return r
}

func (bi bigInteger) squareTo(r *bigInteger) {
	var x = bi.abs()
	r.t = 2 * x.t
	var i = r.t

	// for i--; i >= 0; i-- {
	// 	r.arr[i] = 0
	// }
	r.arr = make([]int64, i)

	for i = 0; i < x.t-1; i++ {
		var c = x.am(i, x.arr[i], r, 2*i, 0, 1)
		r.arr[i+x.t] += x.am(i+1, 2*x.arr[i], r, 2*i+1, c, x.t-i-1)
		if (r.arr[i+x.t]) >= x.DV {
			r.arr[i+x.t] -= x.DV
			r.arr[i+x.t+1] = 1
		}
	}
	if r.t > 0 {
		r.arr[r.t-1] += x.am(i, x.arr[i], r, 2*i, 0, 1)
	}
	r.s = 0
	r.clamp()
}
func (bi bigInteger) multiplyTo(a bigInteger, r *bigInteger) {
	x := bi.abs()
	y := a.abs()

	var i = x.t

	r.t = i + y.t
	for i--; i >= 0; i-- {
		r.arr[i] = 0
	}
	for i = 0; i < y.t; i++ {
		r.arr[i+x.t] = x.am(0, y.arr[i], r, i, 0, x.t)
	}
	r.s = 0
	r.clamp()
	if bi.s != a.s {
		bigIntegerZero.subTo(*r, r)
	}
}

func (bi bigInteger) dlShiftTo(n int, r *bigInteger) {
	r.arr = make([]int64, bi.t+n)

	for i := bi.t - 1; i >= 0; i-- {
		r.arr[i+n] = bi.arr[i]
	}

	for i := n - 1; i >= 0; i-- {
		r.arr[i] = 0
	}
	r.t = bi.t + n
	r.s = bi.s
}

func (bi bigInteger) lShiftTo(n int, r *bigInteger) {
	var bs = int64(n % bi.DB)
	var cbs = int64(bi.DB) - bs
	var bm = int64((1 << cbs) - 1)
	var ds = int(math.Floor(float64(n) / float64(bi.DB)))
	c := (bi.s << bs) & bi.DM
	// fmt.Printf("bm: %v, bs %v, c: %v, cbs: %v, ds: %v\n", bm, bs, c, cbs, ds)

	r.arr = make([]int64, bi.t+ds+1)

	for i := bi.t - 1; i >= 0; i-- {
		r.arr[i+ds+1] = (bi.arr[i] >> cbs) | c
		// fmt.Println((bi.arr[i]>>cbs)|c, r.arr[i+ds+1])
		c = (bi.arr[i] & bm) << bs
		// fmt.Println("c:", c)
	}
	// fmt.Println(r.toString(16))

	for i := ds - 1; i >= 0; i-- {
		r.arr[i] = 0
	}
	r.arr[ds] = c
	r.t = bi.t + ds + 1
	r.s = bi.s
	// fmt.Println(r.arr[ds])
	// fmt.Println(r.t)
	// fmt.Println(r.s)
	r.clamp()
}

func (bi bigInteger) subTo(a bigInteger, r *bigInteger) {
	i := 0
	c := int64(0)
	m := int(math.Min(float64(a.t), float64(bi.t)))

	for i < m {
		c += bi.arr[i] - a.arr[i]
		r.arr[i] = c & bi.DM
		i++
		c >>= bi.DB
	}

	if a.t < bi.t {
		c -= a.s
		for i < bi.t {
			c += bi.arr[i]
			r.arr[i] = c & bi.DM
			i++
			c >>= bi.DB
		}
		c += bi.s
	} else {
		c += bi.s
		for i < a.t {
			c -= a.arr[i]
			r.arr[i] = c & bi.DM
			i++
			c >>= bi.DB
		}
		c -= a.s
	}

	if c < 0 {
		r.s = -1
	} else {
		r.s = 0
	}
	if c < -1 {
		r.arr[i] = bi.DV + c
		i++
	} else if c > 0 {
		r.arr[i] = c
		i++
	}
	r.t = i
	r.clamp()
}

func (bi *bigInteger) negate() bigInteger {
	var r = nbi()
	bigIntegerZero.subTo(*bi, r)
	return *r
}

func (bi bigInteger) abs() bigInteger {
	if bi.s < 0 {
		return bi.negate()
	}
	return bi
}

func (bi *bigInteger) fromInt(x int64) {
	bi.t = 1
	if x < 0 {
		bi.s = -1
	} else {
		bi.s = 0
	}
	if x > 0 {
		// bi.arr[0] = x
		bi.arr = []int64{x}
	} else if x < -1 {
		// bi.arr[0] = x + bi.DV // ??? en el código original DV no existe
		bi.arr = []int64{x + bi.DV}
	} else {
		bi.t = 0
		bi.arr = nil
	}
}

func (bi bigInteger) copyTo(r *bigInteger) {
	r.arr = make([]int64, bi.t)
	for i := bi.t - 1; i >= 0; i-- {
		r.arr[i] = bi.arr[i]
	}
	r.t = bi.t
	r.s = bi.s
}

func (bi bigInteger) compareTo(a bigInteger) int64 {
	r := bi.s - a.s

	if r != 0 {
		return r
	}
	var i = bi.t
	r = int64(i - a.t)

	if r != 0 {
		return r
	}
	for i--; i >= 0; i-- {
		r = bi.arr[i] - a.arr[i]

		if r != 0 {
			return r
		}
	}

	return 0
}

func (bi *bigInteger) am(i int, x int64, w *bigInteger, j int, c int64, n int) int64 {
	// print := func(i int, x int64, w *bigInteger, j int, c int64, n int) {
	// 	fmt.Println("const y = parseBigInt('", bi.toString(16), "', 16)")
	// 	fmt.Println("const i =", i)
	// 	fmt.Println("const x =", x)
	// 	fmt.Println("const w = parseBigInt('" + w.toString(16) + "', 16)")
	// 	fmt.Println("const j =", j)
	// 	fmt.Println("const c =", c)
	// 	fmt.Println("const n =", n)
	// 	fmt.Println("console.log(y.am(i, x, w, j, c, n))")
	// }

	xl := x & 0x3fff
	xh := x >> 14

	for n--; n >= 0; n-- {
		// if len(w.arr) <= j {
		// print(i, x, w, j, c, n)
		// runtime.Breakpoint()
		// }
		var l = bi.arr[i] & 0x3fff
		var h = bi.arr[i] >> 14
		i++
		var m = xh*l + h*xl
		if len(w.arr) <= j {
			w.arr = append(w.arr, 0)
		}
		l = xl*l + ((m & 0x3fff) << 14) + w.arr[j] + c
		c = (l >> 28) + (m >> 14) + xh*h
		w.arr[j] = l & 0xfffffff
		j++
	}

	return c
}
func (bi bigInteger) drShiftTo(n int, r *bigInteger) {
	for i := n; i < bi.t; i++ {
		r.arr[i-n] = bi.arr[i]
	}
	r.t = int(math.Max(float64(bi.t-n), 0))
	r.s = bi.s
}

func (bi bigInteger) rShiftTo(n int, r *bigInteger) {
	r.s = bi.s
	var ds = int(math.Floor(float64(n / bi.DB)))
	if ds >= bi.t {
		r.t = 0
		return
	}
	var bs = n % bi.DB
	var cbs = bi.DB - bs
	var bm = int64((1 << bs) - 1)
	r.arr[0] = bi.arr[ds] >> bs
	for i := ds + 1; i < bi.t; i++ {
		r.arr[i-ds-1] |= (bi.arr[i] & bm) << cbs
		r.arr[i-ds] = bi.arr[i] >> bs
	}
	if bs > 0 {
		r.arr[bi.t-ds-1] |= (bi.s & bm) << cbs
	}
	r.t = bi.t - ds
	r.clamp()
}

func (bi bigInteger) divRemTo(m bigInteger, q, r *bigInteger) {
	var pm = m.abs()
	if pm.t <= 0 {
		return
	}

	var pt = bi.abs()
	if pt.t < pm.t {
		if q != nil {
			q.fromInt(0)
		}
		if r != nil {
			bi.copyTo(r)
		}
		return
	}

	if r == nil {
		r = nbi()
	}

	// fmt.Println(pm.toString(16))
	// fmt.Println(pt.toString(16))
	// fmt.Println(q)
	// fmt.Println(r.toString(16))
	// fmt.Println(bi.toString(16))
	// fmt.Println(m.toString(16))

	y := nbi()
	ts := bi.s
	ms := m.s

	nsh := bi.DB - nbits(pm.arr[pm.t-1])

	if nsh > 0 {
		// fmt.Println(pm.toString(16))
		// fmt.Println(nsh)
		// fmt.Println(y.toString(16))
		pm.lShiftTo(nsh, y)
		pt.lShiftTo(nsh, r)
	} else {
		pm.copyTo(y)
		pt.copyTo(r)
	}
	var ys = y.t
	var y0 = y.arr[ys-1]

	if y0 == 0 {
		return
	}

	// fmt.Println(y.toString(16))
	// fmt.Println(ts)
	// fmt.Println(ms)
	// fmt.Println(nsh)
	// fmt.Println(r.toString(16))
	// fmt.Println(ys)
	// fmt.Println(y0)

	var yt = y0 * (1 << bi.F1)
	if ys > 1 {
		yt += y.arr[ys-2] >> bi.F2
	}

	d1 := bi.FV / float64(yt)
	d2 := float64(int(1)<<bi.F1) / float64(yt)
	e := int64(1 << bi.F2)

	i := r.t
	j := i - ys
	var t *bigInteger
	if q == nil {
		t = nbi()
	} else {
		t = q
	}

	y.dlShiftTo(j, t)

	if r.compareTo(*t) >= 0 {
		r.arr[r.t] = 1
		r.t++
		r.subTo(*t, r)
	}

	// fmt.Println(yt)
	// fmt.Println(y0)
	// fmt.Println(bi.F1)
	// fmt.Println(ys)
	// fmt.Println(d1)
	// fmt.Println(d2)
	// fmt.Println(e)
	// fmt.Println(i)
	// fmt.Println(j)
	// fmt.Println(y.toString(16))
	// fmt.Println(y.toString(16))
	// fmt.Println(t.toString(16))
	// fmt.Println(r.toString(16))

	bigIntegerOne.dlShiftTo(ys, t)
	t.subTo(*y, y)
	for y.t < ys {
		y.arr[y.t] = 0
		y.t++
	}
	for j--; j >= 0; j-- {
		var qd int64
		i--
		if r.arr[i] == y0 {
			qd = bi.DM
		} else {
			qd = int64(math.Floor(float64(r.arr[i])*d1 + float64(r.arr[i-1]+e)*d2))
		}
		r.arr[i] += y.am(0, qd, r, j, 0, ys)
		if (r.arr[i]) < qd {
			y.dlShiftTo(j, t)
			r.subTo(*t, r)
			for qd--; r.arr[i] < qd; qd-- {
				r.subTo(*t, r)
			}
		}
	}

	// fmt.Println(t.toString(16))
	// fmt.Println(y.toString(16))
	// fmt.Println(r.toString(16))

	if q != nil {
		r.drShiftTo(ys, q)
		if ts != ms {
			bigIntegerZero.subTo(*q, q)
		}
	}
	r.t = ys
	r.clamp()
	if nsh > 0 {
		r.rShiftTo(nsh, r)
	}
	if ts < 0 {
		bigIntegerZero.subTo(*r, r)
	}
}

func (mg montgomery) convert(x bigInteger) bigInteger {
	var r = nbi()
	// fmt.Println(mg.m.t)
	x.abs().dlShiftTo(mg.m.t, r)

	r.divRemTo(mg.m, nil, r)
	if x.s < 0 && r.compareTo(*bigIntegerZero) > 0 {
		mg.m.subTo(*r, r)
	}
	return *r
}

func newMontgomery(m bigInteger) montgomery {
	mg := montgomery{m: m}
	mg.mp = m.invDigit()
	mg.mpl = mg.mp & 0x7fff
	mg.mph = mg.mp >> 15
	mg.um = (1 << (m.DB - 15)) - 1
	mg.mt2 = 2 * m.t
	return mg
}

func (bi bigInteger) modPowInt(e int64, m bigInteger) bigInteger {
	var z something

	if e < 256 || m.isEven() {
		// z = new Classic(m);
		z = classic{m: m}
	} else {
		// z = new Montgomery(m);
		z = newMontgomery(m)
	}

	// fmt.Println(z)

	return bi.exp(e, z)
}

// var bigIntegerOne bigInteger
// var bigIntegerOne bigInteger = bigInteger{
// 	t: 1, s: 0,
// 	DB: 28, DM: 268435455, DV: 268435456, F1: 24, F2: 4, FV: 4503599627370496,
// 	arr: []int64{1},
// }

// TODO: implement these
// var bigIntegerZero bigInteger
var bigIntegerOne *bigInteger = nbv(1)
var bigIntegerZero *bigInteger = nbv(0)

func nbi() *bigInteger {
	dbits := 28
	BI_FP := 52

	return &bigInteger{
		DB: dbits,
		DM: (1 << dbits) - 1,
		DV: 1 << dbits,
		FV: math.Pow(2, float64(BI_FP)),
		F1: BI_FP - dbits,
		F2: 2*dbits - BI_FP,
	}
}

func nbv(i int64) *bigInteger {
	var r = nbi()
	r.fromInt(i)
	return r
}

func (bi bigInteger) exp(e int64, z something) bigInteger {
	if e > 0xffffffff || e < 1 {
		return *bigIntegerOne
	}
	r := nbi()
	r2 := nbi()
	g := z.convert(bi)
	i := nbits(e) - 1
	g.copyTo(r)

	// fmt.Println(r.toString(16))
	// fmt.Println(r2.toString(16))
	// fmt.Println(bi.toString(16))
	// fmt.Println(g.toString(16))
	// fmt.Println(i)

	for i--; i >= 0; i-- {
		// fmt.Println(i)
		z.sqrTo(*r, r2)
		if (e & (1 << i)) > 0 {
			z.mulTo(r2, g, r)
		} else {
			var t = r
			r = r2
			r2 = t
		}
	}
	// fmt.Println(r.toString(16))
	return *z.revert(r)
	// return bigInteger{}
}

func (bi bigInteger) invDigit() int64 {
	if bi.t < 1 {
		return 0
	}

	var x = bi.arr[0]
	if (x & 1) == 0 {
		return 0
	}

	var y = x & 3
	y = (y * (2 - (x&0xf)*y)) & 0xf
	y = (y * (2 - (x&0xff)*y)) & 0xff
	y = (y * (2 - (((x & 0xffff) * y) & 0xffff))) & 0xffff
	y = (y * (2 - x*y%bi.DV)) % bi.DV
	if y > 0 {
		return bi.DV - y
	}
	return -y
}

func charCodeAt(s string, n int) rune {
	i := 0
	for _, r := range s {
		if i == n {
			return r
		}
		i++
	}
	return 0
}

var BI_RC map[rune]int

// var BI_RC []int

func intAt(s string, i int) int64 {
	if BI_RC == nil {
		BI_RC = map[rune]int{}

		var rr rune

		rr = '0'
		for vv := 0; vv <= 9; vv++ {
			BI_RC[rr] = vv
			rr++
		}
		rr = 'a'
		for vv := 10; vv < 36; vv++ {
			BI_RC[rr] = vv
			rr++
		}
		rr = 'A'
		for vv := 10; vv < 36; vv++ {
			BI_RC[rr] = vv
			rr++
		}
	}

	if c, found := BI_RC[charCodeAt(s, i)]; found {
		return int64(c)
	}
	return -1

	// c := s[i : i+1]
	// fmt.Println()
	// result, err := strconv.ParseInt(c, 10, 64)
	// if err != nil {
	// 	panic(err)
	// }
	// return result
}

func int2char(n int64) string {
	const BI_RM = "0123456789abcdefghijklmnopqrstuvwxyz"
	return string(BI_RM[n])
}

func newBigInteger(s string, b int) bigInteger {
	bigInt := nbi()
	arr := map[int]int64{}
	// bigInt.arr = make([]int64, 74)
	// bigInt.arr = make([]int64, len(s))
	// bigInt.arr = map[int]int64{}

	var k int
	if b == 16 {
		k = 4
	} else if b == 8 {
		k = 3
	} else if b == 256 {
		k = 8
	} else if b == 2 {
		k = 1
	} else if b == 32 {
		k = 5
	} else if b == 4 {
		k = 2
	} else {
		panic("not implemented")
		// this.fromRadix(s, b);
		// return;
	}

	// mi := false
	sh := 0

	for i := len(s) - 1; i >= 0; i-- {
		var x int64
		if k == 8 {
			x = int64(s[i] & 0xff)
		} else {
			x = intAt(s, i)
		}
		// fmt.Println(x)

		if x < 0 {
			if s[i] == '-' {
				panic("not implemented")
				// mi = true
			}
			continue
		}
		// mi = false

		if sh == 0 {
			arr[bigInt.t] = x
			bigInt.t++
		} else if sh+k > bigInt.DB {
			arr[bigInt.t-1] |= (x & ((1 << (bigInt.DB - sh)) - 1)) << sh
			arr[bigInt.t] = (x >> (bigInt.DB - sh))
			bigInt.t++
		} else {
			arr[bigInt.t-1] |= x << sh
		}
		// fmt.Println(bigInt.t)

		sh += k
		if sh >= bigInt.DB {
			sh -= bigInt.DB
		}
	}

	if k == 8 && (s[0]&0x80) != 0 {
		bigInt.s = -1
		if sh > 0 {
			arr[bigInt.t-1] |= ((1 << (bigInt.DB - sh)) - 1) << sh
		}
	}

	bigInt.arr = mapToSlice(arr)

	bigInt.clamp()

	// fmt.Println(arr)
	// bas := ""
	// ba := mapToSlice(arr)
	// ba := bigInt.arr
	// for _, x := range ba {
	// 	bas += strconv.Itoa(int(x))
	// }
	// fmt.Println(bas)

	// if mi {
	// 	BigInteger.ZERO.subTo(this, this)
	// }

	return *bigInt
}

func mapToSlice(m map[int]int64) []int64 {
	arr := make([]int64, len(m))
	for i, b := range m {
		arr[i] = b
	}
	return arr
}

func (b *bigInteger) clamp() {
	c := b.s & b.DM

	for b.t > 0 && b.arr[b.t-1] == c {
		b.t--
	}

}

func nbits(x int64) int {
	r := 1
	var t int64

	t = x >> 16
	if t != 0 {
		x = t
		r += 16
	}

	t = x >> 8
	if t != 0 {
		x = t
		r += 8
	}

	t = x >> 4
	if t != 0 {
		x = t
		r += 4
	}

	t = x >> 2
	if t != 0 {
		x = t
		r += 2
	}

	t = x >> 1
	if t != 0 {
		// x = t
		r += 1
	}

	return r
}

func (b bigInteger) bitLength() int {
	if b.t <= 0 {
		return 0
	}

	return b.DB*(b.t-1) + nbits(b.arr[b.t-1]^(b.s&b.DM))
}

func pkcs1pad2(s string, n int) bigInteger {
	if n < len(s)+11 {
		panic("Message too long for RSA")
	}

	// ba := []byte{}
	ba := make([]byte, n)
	i := len(s) - 1
	for i >= 0 && n > 0 {
		// fmt.Println(i, n)
		c := charCodeAt(s, i)
		i--
		if c < 128 {
			ba[n-1] = byte(c)
			n--
		} else if c > 127 && c < 2048 {
			ba[n-1] = byte((c & 63) | 128)
			n--
			ba[n-1] = byte((c >> 6) | 192)
			n--
		} else {
			ba[n-1] = byte((c & 63) | 128)
			n--
			ba[n-1] = byte(((c >> 6) & 63) | 128)
			n--
			ba[n-1] = byte((c >> 12) | 224)
			n--
		}
	}
	ba[n-1] = 0
	n--

	// fmt.Println(ba)

	rng := secureRandom{}
	x := make([]byte, 1)

	for n > 2 {
		x[0] = 0

		for x[0] == 0 {
			rng.nextBytes(x)
			// fmt.Println(x[0])
		}
		ba[n-1] = x[0]
		n--
	}
	ba[n-1] = 2
	n--
	ba[n-1] = 0
	n--

	// fmt.Println(len(ba))
	// fmt.Println(ba)
	// bas := ""
	// for _, x := range ba {
	// 	bas += strconv.Itoa(int(x))
	// }
	// fmt.Println(bas)

	return newBigInteger(string(ba), 256)
}

type secureRandom struct {
}

func (sr secureRandom) nextBytes(ba []byte) {
	for i := 0; i < len(ba); i++ {
		ba[i] = rngGetByte()
	}
}

type randomState struct {
	i byte
	j byte
	S []byte
}

func (rs *randomState) next() byte {
	var t byte
	rs.i = (rs.i + 1) & 255
	rs.j = (rs.j + rs.S[rs.i]) & 255
	t = rs.S[rs.i]
	rs.S[rs.i] = rs.S[rs.j]
	rs.S[rs.j] = t
	return rs.S[(t+rs.S[rs.i])&255]
}

func (rs *randomState) init(key []byte) {
	for i := 0; i < 256; i++ {
		rs.S[i] = byte(i)
	}

	var j byte
	for i := 0; i < 256; i++ {
		j = (j + rs.S[i] + key[i%len(key)]) & 255
		t := rs.S[i]
		rs.S[i] = rs.S[j]
		rs.S[j] = t
	}

	rs.i = 0
	rs.j = 0
	// fmt.Println(rs.i)
	// fmt.Println(rs.j)
	// fmt.Println(rs.S)
}

var rngState *randomState

func rngSeedTime() {
	// rng_seed_int(new Date().getTime());
	rngSeedInt(1623517380172)
}

var rng_pptr int
var rng_psize int
var rng_pool []byte

func initRngPool() {
	if rng_pool == nil {
		rng_pool = make([]byte, 256)
		rng_pptr = 0

		r := rand.New(rand.NewSource(99))

		for rng_pptr < rng_psize {
			t := math.Floor(65536 * r.Float64())
			// t := math.Floor(65536 * math.Random())
			rng_pool[rng_pptr] = byte(uint64(t) >> 8)
			rng_pptr++
			rng_pool[rng_pptr] = byte(uint64(t) & 255)
			rng_pptr++
		}
		rng_pptr = 0
		rngSeedTime()
	}
}

func rngSeedInt(x int) {
	rng_pool[rng_pptr] ^= byte(x & 255)
	rng_pptr++
	rng_pool[rng_pptr] ^= byte((x >> 8) & 255)
	rng_pptr++
	rng_pool[rng_pptr] ^= byte((x >> 16) & 255)
	rng_pptr++
	rng_pool[rng_pptr] ^= byte((x >> 24) & 255)
	rng_pptr++

	if rng_pptr >= rng_psize {
		rng_pptr -= rng_psize
	}
	// fmt.Println(rng_pool)
	// fmt.Println(rng_pptr)
}

func prng_newstate() *randomState {
	// Arcfour
	s := make([]byte, 256)
	return &randomState{
		S: s,
	}
}

func rngGetByte() byte {
	if rngState == nil {
		rngSeedTime()
		rngState = prng_newstate()
		rngState.init(rng_pool)

		for rng_pptr = 0; rng_pptr < len(rng_pool); rng_pptr++ {
			rng_pool[rng_pptr] = 0
		}

		rng_pptr = 0
	}

	return rngState.next()
}

func processPassword(password, t1Assertion string) string {
	key := rsaKey{}
	key.setPublic(
		"A6CA1BB4BD803E5704A071E8F7370FD68F2A42CAB574A765693F0F54796CB8AD2CF1B624005119FE651227F7992FF6A6D1979C9B72EA0EAD789F1CBADAB9851779CB8F5F82F40BC71C5C303A10298ED6DC5657E3401AE5720F06836F098366441AC30AB35F13FAB8B6CE81955A1181FCA0AD4EA471CC09C51EAE8EDA42E8C615F933483449CBC67883F407430CB856E4EEC1919BFDD38850CCF5837EC67D8CF802EC30836099592FCDB6CEF4D4AB8EC7F95229B6B262DC6F9A62BFD082CCF98D8FC73FADFA2CCBDDBD17126206E0EC41FE85ECDB9B7631A7EDEF193E4971ADA3E4AB3FFE05F5146907255AD29D0AFB91160C95E225514E1CD07E35BA157A44D1",
		"10001",
	)

	res2 := key.encrypt(password + "|" + t1Assertion)
	return hex2b64(res2)
}

const b64map = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
const b64pad = "="

func hex2b64(h string) string {
	var i int
	var ret = ""

	for i = 0; i+3 <= len(h); i += 3 {
		c, err := strconv.ParseInt(h[i:i+3], 16, 64)
		if err != nil {
			panic(err)
		}
		ret += string(b64map[c>>6]) + string(b64map[c&63])
	}

	if i+1 == len(h) {
		c, err := strconv.ParseInt(h[i:i+1], 16, 64)
		if err != nil {
			panic(err)
		}
		ret += string(b64map[c<<2])
	} else if i+2 == len(h) {
		c, err := strconv.ParseInt(h[i:i+2], 16, 64)
		if err != nil {
			panic(err)
		}
		ret += string(b64map[c>>2]) + string(b64map[(c&3)<<4])
	}

	for (len(ret) & 3) > 0 {
		ret += b64pad
	}

	return ret
}
