package roster

import (
	"github.com/coyim/coyim/xmpp/data"
	g "gopkg.in/check.v1"
)

type GroupListSuite struct{}

var _ = g.Suite(&GroupListSuite{})

func (s *GroupListSuite) Test_TopLevelGroup_returnsATopLevelGroup(c *g.C) {
	result := TopLevelGroup()
	c.Check(result.GroupName, g.Equals, "")
	c.Check(result.fullGroupName, g.DeepEquals, []string{})
	c.Check(result.peers, g.DeepEquals, []*Peer{})
	c.Check(result.groups, g.DeepEquals, map[string]*Group{})
}

func tj(s string) data.JIDWithoutResource {
	return data.JIDNR(s)
}

func tjr(s string) data.JIDWithResource {
	return data.JIDR(s)
}

func (s *GroupListSuite) Test_Grouped_WillGroupPeersInAList(c *g.C) {
	l := New()
	p1 := &Peer{Jid: tj("somewhere"), Name: "something", Groups: toSet("hello"), Subscription: "from"}
	p2 := &Peer{Jid: tj("somewhere2"), Name: "something2", Groups: toSet("hello", "goodbye::foo::bar"), Subscription: "from"}
	p3 := &Peer{Jid: tj("somewhere3"), Name: "something3", Groups: toSet(), Subscription: "from"}
	l.AddOrMerge(p1)
	l.AddOrMerge(p2)
	l.AddOrMerge(p3)

	result := l.Grouped("::")

	c.Check(result.GroupName, g.Equals, "")
	c.Check(result.fullGroupName, g.DeepEquals, []string{})
	c.Check(result.peers, g.DeepEquals, []*Peer{p3})
	c.Check(len(result.groups), g.Equals, 2)
	c.Check(result.groups["goodbye"].GroupName, g.DeepEquals, "goodbye")
	c.Check(result.groups["goodbye"].fullGroupName, g.DeepEquals, []string{"goodbye"})
	c.Check(result.groups["goodbye"].peers, g.DeepEquals, []*Peer{})
	c.Check(len(result.groups["goodbye"].groups), g.Equals, 1)

	c.Check(result.groups["goodbye"].groups["foo"].GroupName, g.DeepEquals, "foo")
	c.Check(result.groups["goodbye"].groups["foo"].fullGroupName, g.DeepEquals, []string{"goodbye", "foo"})
	c.Check(result.groups["goodbye"].groups["foo"].peers, g.DeepEquals, []*Peer{})
	c.Check(len(result.groups["goodbye"].groups["foo"].groups), g.Equals, 1)

	c.Check(result.groups["goodbye"].groups["foo"].groups["bar"].GroupName, g.DeepEquals, "bar")
	c.Check(result.groups["goodbye"].groups["foo"].groups["bar"].fullGroupName, g.DeepEquals, []string{"goodbye", "foo", "bar"})
	c.Check(result.groups["goodbye"].groups["foo"].groups["bar"].peers, g.DeepEquals, []*Peer{p2})
	c.Check(len(result.groups["goodbye"].groups["foo"].groups["bar"].groups), g.Equals, 0)

	c.Check(result.groups["hello"].GroupName, g.DeepEquals, "hello")
	c.Check(result.groups["hello"].fullGroupName, g.DeepEquals, []string{"hello"})
	if result.groups["hello"].peers[0] == p1 {
		c.Check(result.groups["hello"].peers, g.DeepEquals, []*Peer{p1, p2})
	} else {
		c.Check(result.groups["hello"].peers, g.DeepEquals, []*Peer{p2, p1})
	}
	c.Check(len(result.groups["hello"].groups), g.Equals, 0)
}

func (s *GroupListSuite) Test_Groups_willReturnTheGroups(c *g.C) {
	l := New()
	p1 := &Peer{Jid: tj("somewhere"), Name: "something", Groups: toSet("hello"), Subscription: "from"}
	p2 := &Peer{Jid: tj("somewhere2"), Name: "something2", Groups: toSet("hello", "goodbye::foo::bar"), Subscription: "from"}
	p3 := &Peer{Jid: tj("somewhere3"), Name: "something3", Groups: toSet(), Subscription: "from"}
	l.AddOrMerge(p1)
	l.AddOrMerge(p2)
	l.AddOrMerge(p3)

	gr := l.Grouped("::")
	res := gr.Groups()

	c.Check(len(res), g.Equals, 2)
	c.Check(res[0].FullGroupName(), g.Equals, "goodbye")
	c.Check(res[1].FullGroupName(), g.Equals, "hello")
}
