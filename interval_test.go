package main

import "testing"

func TestInterval(t *testing.T) {
	cases := []struct {
		a, b, want Inter
	}{
		{empty, empty, empty},
		{Inter{}, Inter{}, Inter{}},
		{Inter{1, 2}, Inter{1, 2}, Inter{1, 2}},
		{Inter{1, 2}, Inter{3, 4}, empty},
		{Inter{1, 3}, Inter{2, 4}, Inter{2, 3}},
		{Inter{1, 4}, Inter{2, 3}, Inter{2, 3}},
	}

	for _, c := range cases {
		if have := c.a.And(c.b); have != c.want {
			t.Errorf("and %v, %v: have %v, want %v", c.a, c.b, have, c.want)
		}
		if have := c.b.And(c.a); have != c.want {
			t.Errorf("and %v, %v: have %v, want %v", c.b, c.a, have, c.want)
		}
	}
}