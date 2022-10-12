package main

import (
	"github.com/adnsv/go-utils/version"
	"github.com/blang/semver/v4"
)

func makePR(s string, n uint64) []semver.PRVersion {
	return []semver.PRVersion{
		{VersionStr: s},
		{VersionNum: n, IsNum: true},
	}
}

func withPR(v version.Semantic, pr string, pn uint64) version.Semantic {
	v.Pre = makePR(pr, pn)
	return v
}

func withoutPR(v version.Semantic) version.Semantic {
	v.Pre = v.Pre[:0]
	return v
}

type action struct {
	desc         string
	ver          semver.Version
	showPRchoice bool
}

func collectActions(v version.Semantic) []action {
	v.Build = nil
	ret := []action{}

	if len(v.Pre) == 0 {
		if n := v; n.IncrementPatch() == nil {
			ret = append(ret, action{"increment patch|backwards compatible bug fixes", n, true})
		}
		if n := v; n.IncrementMinor() == nil {
			ret = append(ret, action{"increment minor|backwards compatible new functionality", n, true})
		}
		if n := v; n.IncrementMajor() == nil {
			ret = append(ret, action{"increment major|incompatible API changes", n, true})
		}
	} else {
		pr := v.Pre[0].VersionStr
		pn := uint64(0)
		if len(v.Pre) > 1 && v.Pre[1].IsNum {
			pn = v.Pre[1].VersionNum
		}
		ret = append(ret, action{"bump '" + pr + "'", withPR(v, pr, pn+1), false})
		if pr == "alpha" {
			ret = append(ret, action{"upgrade 'alpha' to 'beta'", withPR(v, "beta", 1), false})
			ret = append(ret, action{"upgrade 'alpha' to 'rc'", withPR(v, "rc", 1), false})
			// do not allow to go from alpha to release
		}
		if pr == "beta" {
			ret = append(ret, action{"upgrade 'beta' to 'rc'", withPR(v, "rc", 1), false})
			ret = append(ret, action{"make release", withoutPR(v), false})
		}
		if pr == "rc" {
			ret = append(ret, action{"make release", withoutPR(v), false})
		}
	}

	return ret
}
