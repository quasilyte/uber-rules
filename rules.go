package gorules

import "github.com/quasilyte/go-ruleguard/dsl"

func ifacePtr(m dsl.Matcher) {
	m.Match(`*$x`).
		Where(m["x"].Type.Underlying().Is(`interface{ $*_ }`)).
		Report(`don't use pointers to an interface`)
}

func newMutex(m dsl.Matcher) {
	m.Match(`$mu := new(sync.Mutex); $mu.Lock()`).
		Report(`use zero mutex value instead, 'var $mu sync.Mutex'`).
		At(m["mu"])
}

func channelSize(m dsl.Matcher) {
	m.Match(`make(chan $_, $size)`).
		Where(m["size"].Value.Int() != 0 && m["size"].Value.Int() != 1).
		Report(`channels should have a size of one or be unbuffered`)
}

func uncheckedTypeAssert(m dsl.Matcher) {
	m.Match(
		`$_ := $_.($_)`,
		`$_ = $_.($_)`,
		`$_($*_, $_.($_), $*_)`,
		`$_{$*_, $_.($_), $*_}`,
		`$_{$*_, $_: $_.($_), $*_}`).
		Report(`avoid unchecked type assertions as they can panic`)
}

func unnecessaryElse(m dsl.Matcher) {
	m.Match(`var $v $_; if $cond { $v = $x } else { $v = $y }`).
		Where(m["y"].Pure).
		Report(`rewrite as '$v := $y; if $cond { $v = $x }'`)
}

func localVarDecl(m dsl.Matcher) {
	// TODO: cond for a local scope?
	// m.Match()
}
