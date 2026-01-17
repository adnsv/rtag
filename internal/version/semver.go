package version

import (
	"github.com/adnsv/go-utils/version"
	"github.com/blang/semver/v4"
)

func MakePR(s string, n uint64) []semver.PRVersion {
	return []semver.PRVersion{
		{VersionStr: s},
		{VersionNum: n, IsNum: true},
	}
}

func WithPR(v version.Semantic, pr string, pn uint64) version.Semantic {
	v.Pre = MakePR(pr, pn)
	return v
}

func WithoutPR(v version.Semantic) version.Semantic {
	v.Pre = v.Pre[:0]
	return v
}

type Action struct {
	Desc         string
	Ver          semver.Version
	ShowPRChoice bool
}

func CollectActions(v version.Semantic) []Action {
	v.Build = nil
	ret := []Action{}

	if len(v.Pre) == 0 {
		if n := v; n.IncrementPatch() == nil {
			ret = append(ret, Action{"increment patch|backwards compatible bug fixes", n, true})
		}
		if n := v; n.IncrementMinor() == nil {
			ret = append(ret, Action{"increment minor|backwards compatible new functionality", n, true})
		}
		if n := v; n.IncrementMajor() == nil {
			ret = append(ret, Action{"increment major|incompatible API changes", n, true})
		}
	} else {
		pr := v.Pre[0].VersionStr
		pn := uint64(0)
		if len(v.Pre) > 1 && v.Pre[1].IsNum {
			pn = v.Pre[1].VersionNum
		}
		ret = append(ret, Action{"bump '" + pr + "'", WithPR(v, pr, pn+1), false})
		if pr == "alpha" {
			ret = append(ret, Action{"upgrade 'alpha' to 'beta'", WithPR(v, "beta", 1), false})
			ret = append(ret, Action{"upgrade 'alpha' to 'rc'", WithPR(v, "rc", 1), false})
			ret = append(ret, Action{"make release", WithoutPR(v), false})
		}
		if pr == "beta" {
			ret = append(ret, Action{"upgrade 'beta' to 'rc'", WithPR(v, "rc", 1), false})
			ret = append(ret, Action{"make release", WithoutPR(v), false})
		}
		if pr == "rc" {
			ret = append(ret, Action{"make release", WithoutPR(v), false})
		}
	}

	return ret
}
