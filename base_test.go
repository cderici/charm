// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package charm_test

import (
	"encoding/json"
	"strings"

	"github.com/juju/os/v2"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v3/arch"
	gc "gopkg.in/check.v1"

	"github.com/juju/charm/v9"
)

type baseSuite struct {
	testing.CleanupSuite
}

var _ = gc.Suite(&baseSuite{})

func (s *baseSuite) TestParseBase(c *gc.C) {
	tests := []struct {
		base       charm.Base
		str        string
		parsedBase charm.Base
		err        string
	}{
		{
			base:       charm.Base{Name: os.Ubuntu.String()},
			str:        "ubuntu",
			parsedBase: charm.Base{},
			err:        `invalid base string "ubuntu": channel not valid`,
		}, {
			base:       charm.Base{Name: os.Windows.String()},
			str:        "windows",
			parsedBase: charm.Base{},
			err:        `invalid base string "windows": channel not valid`,
		}, {
			base:       charm.Base{Name: "mythicalos"},
			str:        "mythicalos",
			parsedBase: charm.Base{},
			err:        `invalid base string "mythicalos": os "mythicalos" not valid`,
		}, {
			base:       charm.Base{Name: os.Ubuntu.String(), Channel: mustParseChannel("20.04/stable")},
			str:        "ubuntu/20.04/stable",
			parsedBase: charm.Base{Name: strings.ToLower(os.Ubuntu.String()), Channel: mustParseChannel("20.04/stable")},
		}, {
			base:       charm.Base{Name: os.Windows.String(), Channel: mustParseChannel("win10/stable")},
			str:        "windows/win10/stable",
			parsedBase: charm.Base{Name: strings.ToLower(os.Windows.String()), Channel: mustParseChannel("win10/stable")},
		}, {
			base:       charm.Base{Name: os.Ubuntu.String(), Channel: mustParseChannel("20.04/edge")},
			str:        "ubuntu/20.04/edge",
			parsedBase: charm.Base{Name: strings.ToLower(os.Ubuntu.String()), Channel: mustParseChannel("20.04/edge")},
		},
	}
	for i, v := range tests {
		str := v.base.String()
		comment := gc.Commentf("test %d", i)
		c.Check(str, gc.Equals, v.str, comment)
		s, err := charm.ParseBase(str)
		if v.err != "" {
			c.Check(err, gc.ErrorMatches, v.err, comment)
		} else {
			c.Check(err, jc.ErrorIsNil, comment)
		}
		c.Check(s, jc.DeepEquals, v.parsedBase, comment)
	}
}

func (s *baseSuite) TestParseBaseWithArchitectures(c *gc.C) {
	tests := []struct {
		base       charm.Base
		str        string
		baseString string
		archs      []string
		parsedBase charm.Base
		err        string
	}{
		{
			base:       charm.Base{Name: os.Ubuntu.String(), Architectures: []string{arch.AMD64}},
			baseString: "ubuntu",
			str:        "ubuntu on amd64",
			archs:      []string{"amd64"},
			parsedBase: charm.Base{},
			err:        `invalid base string "ubuntu" with architectures "amd64": channel not valid`,
		}, {
			base:       charm.Base{Name: os.Windows.String()},
			baseString: "windows",
			str:        "windows",
			parsedBase: charm.Base{},
			err:        `invalid base string "windows": channel not valid`,
		}, {
			base:       charm.Base{Name: "mythicalos"},
			baseString: "mythicalos",
			str:        "mythicalos",
			parsedBase: charm.Base{},
			err:        `invalid base string "mythicalos": os "mythicalos" not valid`,
		}, {
			base: charm.Base{
				Name:          os.Ubuntu.String(),
				Channel:       mustParseChannel("20.04/stable"),
				Architectures: []string{arch.AMD64, arch.PPC64EL},
			},
			baseString: "ubuntu/20.04/stable",
			archs:      []string{arch.AMD64, "ppc64"},
			str:        "ubuntu/20.04/stable on amd64, ppc64el",
			parsedBase: charm.Base{
				Name:          strings.ToLower(os.Ubuntu.String()),
				Channel:       mustParseChannel("20.04/stable"),
				Architectures: []string{arch.AMD64, arch.PPC64EL}},
		}, {
			base:       charm.Base{Name: os.Windows.String(), Channel: mustParseChannel("win10/stable")},
			baseString: "windows/win10/stable",
			archs:      []string{"testme"},
			str:        "windows/win10/stable",
			parsedBase: charm.Base{},
			err:        `invalid base string "windows/win10/stable" with architectures "testme": architecture "testme" not valid`,
		},
	}
	for i, v := range tests {
		str := v.base.String()
		comment := gc.Commentf("test %d", i)
		c.Check(str, gc.Equals, v.str, comment)
		s, err := charm.ParseBase(v.baseString, v.archs...)
		if v.err != "" {
			c.Check(err, gc.ErrorMatches, v.err, comment)
		} else {
			c.Check(err, jc.ErrorIsNil, comment)
		}
		c.Check(s, jc.DeepEquals, v.parsedBase, comment)
	}
}

func (s *baseSuite) TestJSONEncoding(c *gc.C) {
	sys := charm.Base{
		Name:    "ubuntu",
		Channel: mustParseChannel("20.04/stable"),
	}
	bytes, err := json.Marshal(sys)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(bytes), gc.Equals, `{"name":"ubuntu","channel":{"track":"20.04","risk":"stable"}}`)
	sys2 := charm.Base{}
	err = json.Unmarshal(bytes, &sys2)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(sys2, jc.DeepEquals, sys)
}

// MustParseChannel parses a given string or returns a panic.
// Used for unit tests.
func mustParseChannel(s string) charm.Channel {
	c, err := charm.ParseChannelNormalize(s)
	if err != nil {
		panic(err)
	}
	return c
}
