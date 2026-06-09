package todo

// testSchema mirrors a typical four-section + done layout used across the
// domain tests.
func testSchema() Schema {
	return Schema{Sections: []Section{
		{Key: "in_progress", Title: "In Progress", Letter: "P"},
		{Key: "now", Title: "Now", Letter: "N"},
		{Key: "next", Title: "Next", Letter: "X"},
		{Key: "later", Title: "Later", Letter: "L"},
		{Key: "done", Title: "Done", Done: true},
	}}
}

// sample builds a small normalized list for tests.
func sample() *List {
	l := &List{Items: []Item{
		{Task: "a", Section: "now", Order: 0},
		{Task: "b", Section: "now", Order: 1},
		{Task: "c", Section: "next", Order: 0},
		{Task: "d", Section: "later", Order: 0},
	}}
	l.Normalize(testSchema())
	return l
}
