// Generated by tmpl
// https://github.com/benbjohnson/tmpl
//
// DO NOT EDIT!
// Source: batch_cursor.gen.go.tmpl

package storage

import (
	"errors"

	"github.com/influxdata/influxdb/tsdb"
)

// ********************
// Float BatchCursor

type floatFilterBatchCursor struct {
	tsdb.FloatBatchCursor
	cond expression
	m    *singleValue
	t    []int64
	v    []float64
	tb   []int64
	vb   []float64
}

func newFloatFilterBatchCursor(cond expression) *floatFilterBatchCursor {
	return &floatFilterBatchCursor{
		cond: cond,
		m:    &singleValue{},
		t:    make([]int64, tsdb.DefaultMaxPointsPerBlock),
		v:    make([]float64, tsdb.DefaultMaxPointsPerBlock),
	}
}

func (c *floatFilterBatchCursor) reset(cur tsdb.FloatBatchCursor) {
	c.FloatBatchCursor = cur
	c.tb, c.vb = nil, nil
}

func (c *floatFilterBatchCursor) Next() (key []int64, value []float64) {
	pos := 0
	var ks []int64
	var vs []float64

	if len(c.tb) > 0 {
		ks, vs = c.tb, c.vb
		c.tb, c.vb = nil, nil
	} else {
		ks, vs = c.FloatBatchCursor.Next()
	}

	for len(ks) > 0 {
		for i, v := range vs {
			c.m.v = v
			if c.cond.EvalBool(c.m) {
				c.t[pos], c.v[pos] = ks[i], v
				pos++
				if pos >= tsdb.DefaultMaxPointsPerBlock {
					c.tb, c.vb = ks[i+1:], vs[i+1:]
					return c.t[:pos], c.v[:pos]
				}
			}
		}
		ks, vs = c.FloatBatchCursor.Next()
	}

	return c.t[:pos], c.v[:pos]
}

type floatMultiShardBatchCursor struct {
	tsdb.FloatBatchCursor
	cursorContext
	filter *floatFilterBatchCursor
}

func (c *floatMultiShardBatchCursor) reset(cur tsdb.FloatBatchCursor, itrs tsdb.CursorIterators, cond expression) {
	if cond != nil {
		if c.filter == nil {
			c.filter = newFloatFilterBatchCursor(cond)
		}
		c.filter.reset(cur)
		cur = c.filter
	}

	c.FloatBatchCursor = cur
	c.itrs = itrs
	c.err = nil
}

func (c *floatMultiShardBatchCursor) Err() error { return c.err }

func (c *floatMultiShardBatchCursor) Next() (key []int64, value []float64) {
	for {
		ks, vs := c.FloatBatchCursor.Next()
		if len(ks) == 0 {
			if c.nextBatchCursor() {
				continue
			}
		}
		c.count += int64(len(ks))
		if c.count > c.limit {
			diff := c.count - c.limit
			c.count -= diff
			rem := int64(len(ks)) - diff
			ks = ks[:rem]
			vs = vs[:rem]
		}
		return ks, vs
	}
}

func (c *floatMultiShardBatchCursor) nextBatchCursor() bool {
	if len(c.itrs) == 0 {
		return false
	}

	c.FloatBatchCursor.Close()

	var itr tsdb.CursorIterator
	var cur tsdb.Cursor
	for cur == nil && len(c.itrs) > 0 {
		itr, c.itrs = c.itrs[0], c.itrs[1:]
		cur, _ = itr.Next(c.ctx, c.req)
	}

	var ok bool
	if cur != nil {
		var next tsdb.FloatBatchCursor
		next, ok = cur.(tsdb.FloatBatchCursor)
		if !ok {
			cur.Close()
			next = FloatEmptyBatchCursor
			c.itrs = nil
			c.err = errors.New("expected float cursor")
		} else {
			if c.filter != nil {
				c.filter.reset(next)
				next = c.filter
			}
		}
		c.FloatBatchCursor = next
	} else {
		c.FloatBatchCursor = FloatEmptyBatchCursor
	}

	return ok
}

type floatSumBatchCursor struct {
	tsdb.FloatBatchCursor
	ts [1]int64
	vs [1]float64
}

func (c *floatSumBatchCursor) Next() (key []int64, value []float64) {
	ks, vs := c.FloatBatchCursor.Next()
	if len(ks) == 0 {
		return nil, nil
	}

	ts := ks[0]
	var acc float64

	for {
		for _, v := range vs {
			acc += v
		}
		ks, vs = c.FloatBatchCursor.Next()
		if len(ks) == 0 {
			c.ts[0] = ts
			c.vs[0] = acc
			return c.ts[:], c.vs[:]
		}
	}
}

type integerFloatCountBatchCursor struct {
	tsdb.FloatBatchCursor
}

func (c *integerFloatCountBatchCursor) Next() (key []int64, value []int64) {
	ks, _ := c.FloatBatchCursor.Next()
	if len(ks) == 0 {
		return nil, nil
	}

	ts := ks[0]
	var acc int64
	for {
		acc += int64(len(ks))
		ks, _ = c.FloatBatchCursor.Next()
		if len(ks) == 0 {
			return []int64{ts}, []int64{acc}
		}
	}
}

type floatEmptyBatchCursor struct{}

var FloatEmptyBatchCursor tsdb.FloatBatchCursor = &floatEmptyBatchCursor{}

func (*floatEmptyBatchCursor) Err() error                           { return nil }
func (*floatEmptyBatchCursor) Close()                               {}
func (*floatEmptyBatchCursor) Next() (key []int64, value []float64) { return nil, nil }

// ********************
// Integer BatchCursor

type integerFilterBatchCursor struct {
	tsdb.IntegerBatchCursor
	cond expression
	m    *singleValue
	t    []int64
	v    []int64
	tb   []int64
	vb   []int64
}

func newIntegerFilterBatchCursor(cond expression) *integerFilterBatchCursor {
	return &integerFilterBatchCursor{
		cond: cond,
		m:    &singleValue{},
		t:    make([]int64, tsdb.DefaultMaxPointsPerBlock),
		v:    make([]int64, tsdb.DefaultMaxPointsPerBlock),
	}
}

func (c *integerFilterBatchCursor) reset(cur tsdb.IntegerBatchCursor) {
	c.IntegerBatchCursor = cur
	c.tb, c.vb = nil, nil
}

func (c *integerFilterBatchCursor) Next() (key []int64, value []int64) {
	pos := 0
	var ks []int64
	var vs []int64

	if len(c.tb) > 0 {
		ks, vs = c.tb, c.vb
		c.tb, c.vb = nil, nil
	} else {
		ks, vs = c.IntegerBatchCursor.Next()
	}

	for len(ks) > 0 {
		for i, v := range vs {
			c.m.v = v
			if c.cond.EvalBool(c.m) {
				c.t[pos], c.v[pos] = ks[i], v
				pos++
				if pos >= tsdb.DefaultMaxPointsPerBlock {
					c.tb, c.vb = ks[i+1:], vs[i+1:]
					return c.t[:pos], c.v[:pos]
				}
			}
		}
		ks, vs = c.IntegerBatchCursor.Next()
	}

	return c.t[:pos], c.v[:pos]
}

type integerMultiShardBatchCursor struct {
	tsdb.IntegerBatchCursor
	cursorContext
	filter *integerFilterBatchCursor
}

func (c *integerMultiShardBatchCursor) reset(cur tsdb.IntegerBatchCursor, itrs tsdb.CursorIterators, cond expression) {
	if cond != nil {
		if c.filter == nil {
			c.filter = newIntegerFilterBatchCursor(cond)
		}
		c.filter.reset(cur)
		cur = c.filter
	}

	c.IntegerBatchCursor = cur
	c.itrs = itrs
	c.err = nil
}

func (c *integerMultiShardBatchCursor) Err() error { return c.err }

func (c *integerMultiShardBatchCursor) Next() (key []int64, value []int64) {
	for {
		ks, vs := c.IntegerBatchCursor.Next()
		if len(ks) == 0 {
			if c.nextBatchCursor() {
				continue
			}
		}
		c.count += int64(len(ks))
		if c.count > c.limit {
			diff := c.count - c.limit
			c.count -= diff
			rem := int64(len(ks)) - diff
			ks = ks[:rem]
			vs = vs[:rem]
		}
		return ks, vs
	}
}

func (c *integerMultiShardBatchCursor) nextBatchCursor() bool {
	if len(c.itrs) == 0 {
		return false
	}

	c.IntegerBatchCursor.Close()

	var itr tsdb.CursorIterator
	var cur tsdb.Cursor
	for cur == nil && len(c.itrs) > 0 {
		itr, c.itrs = c.itrs[0], c.itrs[1:]
		cur, _ = itr.Next(c.ctx, c.req)
	}

	var ok bool
	if cur != nil {
		var next tsdb.IntegerBatchCursor
		next, ok = cur.(tsdb.IntegerBatchCursor)
		if !ok {
			cur.Close()
			next = IntegerEmptyBatchCursor
			c.itrs = nil
			c.err = errors.New("expected integer cursor")
		} else {
			if c.filter != nil {
				c.filter.reset(next)
				next = c.filter
			}
		}
		c.IntegerBatchCursor = next
	} else {
		c.IntegerBatchCursor = IntegerEmptyBatchCursor
	}

	return ok
}

type integerSumBatchCursor struct {
	tsdb.IntegerBatchCursor
	ts [1]int64
	vs [1]int64
}

func (c *integerSumBatchCursor) Next() (key []int64, value []int64) {
	ks, vs := c.IntegerBatchCursor.Next()
	if len(ks) == 0 {
		return nil, nil
	}

	ts := ks[0]
	var acc int64

	for {
		for _, v := range vs {
			acc += v
		}
		ks, vs = c.IntegerBatchCursor.Next()
		if len(ks) == 0 {
			c.ts[0] = ts
			c.vs[0] = acc
			return c.ts[:], c.vs[:]
		}
	}
}

type integerIntegerCountBatchCursor struct {
	tsdb.IntegerBatchCursor
}

func (c *integerIntegerCountBatchCursor) Next() (key []int64, value []int64) {
	ks, _ := c.IntegerBatchCursor.Next()
	if len(ks) == 0 {
		return nil, nil
	}

	ts := ks[0]
	var acc int64
	for {
		acc += int64(len(ks))
		ks, _ = c.IntegerBatchCursor.Next()
		if len(ks) == 0 {
			return []int64{ts}, []int64{acc}
		}
	}
}

type integerEmptyBatchCursor struct{}

var IntegerEmptyBatchCursor tsdb.IntegerBatchCursor = &integerEmptyBatchCursor{}

func (*integerEmptyBatchCursor) Err() error                         { return nil }
func (*integerEmptyBatchCursor) Close()                             {}
func (*integerEmptyBatchCursor) Next() (key []int64, value []int64) { return nil, nil }

// ********************
// Unsigned BatchCursor

type unsignedFilterBatchCursor struct {
	tsdb.UnsignedBatchCursor
	cond expression
	m    *singleValue
	t    []int64
	v    []uint64
	tb   []int64
	vb   []uint64
}

func newUnsignedFilterBatchCursor(cond expression) *unsignedFilterBatchCursor {
	return &unsignedFilterBatchCursor{
		cond: cond,
		m:    &singleValue{},
		t:    make([]int64, tsdb.DefaultMaxPointsPerBlock),
		v:    make([]uint64, tsdb.DefaultMaxPointsPerBlock),
	}
}

func (c *unsignedFilterBatchCursor) reset(cur tsdb.UnsignedBatchCursor) {
	c.UnsignedBatchCursor = cur
	c.tb, c.vb = nil, nil
}

func (c *unsignedFilterBatchCursor) Next() (key []int64, value []uint64) {
	pos := 0
	var ks []int64
	var vs []uint64

	if len(c.tb) > 0 {
		ks, vs = c.tb, c.vb
		c.tb, c.vb = nil, nil
	} else {
		ks, vs = c.UnsignedBatchCursor.Next()
	}

	for len(ks) > 0 {
		for i, v := range vs {
			c.m.v = v
			if c.cond.EvalBool(c.m) {
				c.t[pos], c.v[pos] = ks[i], v
				pos++
				if pos >= tsdb.DefaultMaxPointsPerBlock {
					c.tb, c.vb = ks[i+1:], vs[i+1:]
					return c.t[:pos], c.v[:pos]
				}
			}
		}
		ks, vs = c.UnsignedBatchCursor.Next()
	}

	return c.t[:pos], c.v[:pos]
}

type unsignedMultiShardBatchCursor struct {
	tsdb.UnsignedBatchCursor
	cursorContext
	filter *unsignedFilterBatchCursor
}

func (c *unsignedMultiShardBatchCursor) reset(cur tsdb.UnsignedBatchCursor, itrs tsdb.CursorIterators, cond expression) {
	if cond != nil {
		if c.filter == nil {
			c.filter = newUnsignedFilterBatchCursor(cond)
		}
		c.filter.reset(cur)
		cur = c.filter
	}

	c.UnsignedBatchCursor = cur
	c.itrs = itrs
	c.err = nil
}

func (c *unsignedMultiShardBatchCursor) Err() error { return c.err }

func (c *unsignedMultiShardBatchCursor) Next() (key []int64, value []uint64) {
	for {
		ks, vs := c.UnsignedBatchCursor.Next()
		if len(ks) == 0 {
			if c.nextBatchCursor() {
				continue
			}
		}
		c.count += int64(len(ks))
		if c.count > c.limit {
			diff := c.count - c.limit
			c.count -= diff
			rem := int64(len(ks)) - diff
			ks = ks[:rem]
			vs = vs[:rem]
		}
		return ks, vs
	}
}

func (c *unsignedMultiShardBatchCursor) nextBatchCursor() bool {
	if len(c.itrs) == 0 {
		return false
	}

	c.UnsignedBatchCursor.Close()

	var itr tsdb.CursorIterator
	var cur tsdb.Cursor
	for cur == nil && len(c.itrs) > 0 {
		itr, c.itrs = c.itrs[0], c.itrs[1:]
		cur, _ = itr.Next(c.ctx, c.req)
	}

	var ok bool
	if cur != nil {
		var next tsdb.UnsignedBatchCursor
		next, ok = cur.(tsdb.UnsignedBatchCursor)
		if !ok {
			cur.Close()
			next = UnsignedEmptyBatchCursor
			c.itrs = nil
			c.err = errors.New("expected unsigned cursor")
		} else {
			if c.filter != nil {
				c.filter.reset(next)
				next = c.filter
			}
		}
		c.UnsignedBatchCursor = next
	} else {
		c.UnsignedBatchCursor = UnsignedEmptyBatchCursor
	}

	return ok
}

type unsignedSumBatchCursor struct {
	tsdb.UnsignedBatchCursor
	ts [1]int64
	vs [1]uint64
}

func (c *unsignedSumBatchCursor) Next() (key []int64, value []uint64) {
	ks, vs := c.UnsignedBatchCursor.Next()
	if len(ks) == 0 {
		return nil, nil
	}

	ts := ks[0]
	var acc uint64

	for {
		for _, v := range vs {
			acc += v
		}
		ks, vs = c.UnsignedBatchCursor.Next()
		if len(ks) == 0 {
			c.ts[0] = ts
			c.vs[0] = acc
			return c.ts[:], c.vs[:]
		}
	}
}

type integerUnsignedCountBatchCursor struct {
	tsdb.UnsignedBatchCursor
}

func (c *integerUnsignedCountBatchCursor) Next() (key []int64, value []int64) {
	ks, _ := c.UnsignedBatchCursor.Next()
	if len(ks) == 0 {
		return nil, nil
	}

	ts := ks[0]
	var acc int64
	for {
		acc += int64(len(ks))
		ks, _ = c.UnsignedBatchCursor.Next()
		if len(ks) == 0 {
			return []int64{ts}, []int64{acc}
		}
	}
}

type unsignedEmptyBatchCursor struct{}

var UnsignedEmptyBatchCursor tsdb.UnsignedBatchCursor = &unsignedEmptyBatchCursor{}

func (*unsignedEmptyBatchCursor) Err() error                          { return nil }
func (*unsignedEmptyBatchCursor) Close()                              {}
func (*unsignedEmptyBatchCursor) Next() (key []int64, value []uint64) { return nil, nil }

// ********************
// String BatchCursor

type stringFilterBatchCursor struct {
	tsdb.StringBatchCursor
	cond expression
	m    *singleValue
	t    []int64
	v    []string
	tb   []int64
	vb   []string
}

func newStringFilterBatchCursor(cond expression) *stringFilterBatchCursor {
	return &stringFilterBatchCursor{
		cond: cond,
		m:    &singleValue{},
		t:    make([]int64, tsdb.DefaultMaxPointsPerBlock),
		v:    make([]string, tsdb.DefaultMaxPointsPerBlock),
	}
}

func (c *stringFilterBatchCursor) reset(cur tsdb.StringBatchCursor) {
	c.StringBatchCursor = cur
	c.tb, c.vb = nil, nil
}

func (c *stringFilterBatchCursor) Next() (key []int64, value []string) {
	pos := 0
	var ks []int64
	var vs []string

	if len(c.tb) > 0 {
		ks, vs = c.tb, c.vb
		c.tb, c.vb = nil, nil
	} else {
		ks, vs = c.StringBatchCursor.Next()
	}

	for len(ks) > 0 {
		for i, v := range vs {
			c.m.v = v
			if c.cond.EvalBool(c.m) {
				c.t[pos], c.v[pos] = ks[i], v
				pos++
				if pos >= tsdb.DefaultMaxPointsPerBlock {
					c.tb, c.vb = ks[i+1:], vs[i+1:]
					return c.t[:pos], c.v[:pos]
				}
			}
		}
		ks, vs = c.StringBatchCursor.Next()
	}

	return c.t[:pos], c.v[:pos]
}

type stringMultiShardBatchCursor struct {
	tsdb.StringBatchCursor
	cursorContext
	filter *stringFilterBatchCursor
}

func (c *stringMultiShardBatchCursor) reset(cur tsdb.StringBatchCursor, itrs tsdb.CursorIterators, cond expression) {
	if cond != nil {
		if c.filter == nil {
			c.filter = newStringFilterBatchCursor(cond)
		}
		c.filter.reset(cur)
		cur = c.filter
	}

	c.StringBatchCursor = cur
	c.itrs = itrs
	c.err = nil
}

func (c *stringMultiShardBatchCursor) Err() error { return c.err }

func (c *stringMultiShardBatchCursor) Next() (key []int64, value []string) {
	for {
		ks, vs := c.StringBatchCursor.Next()
		if len(ks) == 0 {
			if c.nextBatchCursor() {
				continue
			}
		}
		c.count += int64(len(ks))
		if c.count > c.limit {
			diff := c.count - c.limit
			c.count -= diff
			rem := int64(len(ks)) - diff
			ks = ks[:rem]
			vs = vs[:rem]
		}
		return ks, vs
	}
}

func (c *stringMultiShardBatchCursor) nextBatchCursor() bool {
	if len(c.itrs) == 0 {
		return false
	}

	c.StringBatchCursor.Close()

	var itr tsdb.CursorIterator
	var cur tsdb.Cursor
	for cur == nil && len(c.itrs) > 0 {
		itr, c.itrs = c.itrs[0], c.itrs[1:]
		cur, _ = itr.Next(c.ctx, c.req)
	}

	var ok bool
	if cur != nil {
		var next tsdb.StringBatchCursor
		next, ok = cur.(tsdb.StringBatchCursor)
		if !ok {
			cur.Close()
			next = StringEmptyBatchCursor
			c.itrs = nil
			c.err = errors.New("expected string cursor")
		} else {
			if c.filter != nil {
				c.filter.reset(next)
				next = c.filter
			}
		}
		c.StringBatchCursor = next
	} else {
		c.StringBatchCursor = StringEmptyBatchCursor
	}

	return ok
}

type integerStringCountBatchCursor struct {
	tsdb.StringBatchCursor
}

func (c *integerStringCountBatchCursor) Next() (key []int64, value []int64) {
	ks, _ := c.StringBatchCursor.Next()
	if len(ks) == 0 {
		return nil, nil
	}

	ts := ks[0]
	var acc int64
	for {
		acc += int64(len(ks))
		ks, _ = c.StringBatchCursor.Next()
		if len(ks) == 0 {
			return []int64{ts}, []int64{acc}
		}
	}
}

type stringEmptyBatchCursor struct{}

var StringEmptyBatchCursor tsdb.StringBatchCursor = &stringEmptyBatchCursor{}

func (*stringEmptyBatchCursor) Err() error                          { return nil }
func (*stringEmptyBatchCursor) Close()                              {}
func (*stringEmptyBatchCursor) Next() (key []int64, value []string) { return nil, nil }

// ********************
// Boolean BatchCursor

type booleanFilterBatchCursor struct {
	tsdb.BooleanBatchCursor
	cond expression
	m    *singleValue
	t    []int64
	v    []bool
	tb   []int64
	vb   []bool
}

func newBooleanFilterBatchCursor(cond expression) *booleanFilterBatchCursor {
	return &booleanFilterBatchCursor{
		cond: cond,
		m:    &singleValue{},
		t:    make([]int64, tsdb.DefaultMaxPointsPerBlock),
		v:    make([]bool, tsdb.DefaultMaxPointsPerBlock),
	}
}

func (c *booleanFilterBatchCursor) reset(cur tsdb.BooleanBatchCursor) {
	c.BooleanBatchCursor = cur
	c.tb, c.vb = nil, nil
}

func (c *booleanFilterBatchCursor) Next() (key []int64, value []bool) {
	pos := 0
	var ks []int64
	var vs []bool

	if len(c.tb) > 0 {
		ks, vs = c.tb, c.vb
		c.tb, c.vb = nil, nil
	} else {
		ks, vs = c.BooleanBatchCursor.Next()
	}

	for len(ks) > 0 {
		for i, v := range vs {
			c.m.v = v
			if c.cond.EvalBool(c.m) {
				c.t[pos], c.v[pos] = ks[i], v
				pos++
				if pos >= tsdb.DefaultMaxPointsPerBlock {
					c.tb, c.vb = ks[i+1:], vs[i+1:]
					return c.t[:pos], c.v[:pos]
				}
			}
		}
		ks, vs = c.BooleanBatchCursor.Next()
	}

	return c.t[:pos], c.v[:pos]
}

type booleanMultiShardBatchCursor struct {
	tsdb.BooleanBatchCursor
	cursorContext
	filter *booleanFilterBatchCursor
}

func (c *booleanMultiShardBatchCursor) reset(cur tsdb.BooleanBatchCursor, itrs tsdb.CursorIterators, cond expression) {
	if cond != nil {
		if c.filter == nil {
			c.filter = newBooleanFilterBatchCursor(cond)
		}
		c.filter.reset(cur)
		cur = c.filter
	}

	c.BooleanBatchCursor = cur
	c.itrs = itrs
	c.err = nil
}

func (c *booleanMultiShardBatchCursor) Err() error { return c.err }

func (c *booleanMultiShardBatchCursor) Next() (key []int64, value []bool) {
	for {
		ks, vs := c.BooleanBatchCursor.Next()
		if len(ks) == 0 {
			if c.nextBatchCursor() {
				continue
			}
		}
		c.count += int64(len(ks))
		if c.count > c.limit {
			diff := c.count - c.limit
			c.count -= diff
			rem := int64(len(ks)) - diff
			ks = ks[:rem]
			vs = vs[:rem]
		}
		return ks, vs
	}
}

func (c *booleanMultiShardBatchCursor) nextBatchCursor() bool {
	if len(c.itrs) == 0 {
		return false
	}

	c.BooleanBatchCursor.Close()

	var itr tsdb.CursorIterator
	var cur tsdb.Cursor
	for cur == nil && len(c.itrs) > 0 {
		itr, c.itrs = c.itrs[0], c.itrs[1:]
		cur, _ = itr.Next(c.ctx, c.req)
	}

	var ok bool
	if cur != nil {
		var next tsdb.BooleanBatchCursor
		next, ok = cur.(tsdb.BooleanBatchCursor)
		if !ok {
			cur.Close()
			next = BooleanEmptyBatchCursor
			c.itrs = nil
			c.err = errors.New("expected boolean cursor")
		} else {
			if c.filter != nil {
				c.filter.reset(next)
				next = c.filter
			}
		}
		c.BooleanBatchCursor = next
	} else {
		c.BooleanBatchCursor = BooleanEmptyBatchCursor
	}

	return ok
}

type integerBooleanCountBatchCursor struct {
	tsdb.BooleanBatchCursor
}

func (c *integerBooleanCountBatchCursor) Next() (key []int64, value []int64) {
	ks, _ := c.BooleanBatchCursor.Next()
	if len(ks) == 0 {
		return nil, nil
	}

	ts := ks[0]
	var acc int64
	for {
		acc += int64(len(ks))
		ks, _ = c.BooleanBatchCursor.Next()
		if len(ks) == 0 {
			return []int64{ts}, []int64{acc}
		}
	}
}

type booleanEmptyBatchCursor struct{}

var BooleanEmptyBatchCursor tsdb.BooleanBatchCursor = &booleanEmptyBatchCursor{}

func (*booleanEmptyBatchCursor) Err() error                        { return nil }
func (*booleanEmptyBatchCursor) Close()                            {}
func (*booleanEmptyBatchCursor) Next() (key []int64, value []bool) { return nil, nil }
