package tagcloud

type TagCloud struct {
	tags  map[string]int
	stats []TagStat
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

func NewTagStat(tag string) TagStat {
	return TagStat{Tag: tag, OccurrenceCount: 1}
}

// New should create a valid TagCloud instance
// TODO: You decide whether this function should return a pointer or a value
func New() *TagCloud {
	return &TagCloud{
		tags:  map[string]int{},
		stats: []TagStat{},
	}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
// TODO: You decide whether receiver should be a pointer or a value
func (t *TagCloud) AddTag(tag string) {
	tsInd, ok := t.tags[tag]
	if ok {
		t.updateTag(tsInd)
	} else {
		t.addNewTag(tag)
	}
}
func (t *TagCloud) addNewTag(tag string) {
	tsInd := len(t.stats)

	t.tags[tag] = tsInd

	t.stats = append(t.stats, NewTagStat(tag))
}

func (t *TagCloud) updateTag(tsInd int) {
	t.stats[tsInd].OccurrenceCount++

	newInd := tsInd - 1
	for t.validIndex(newInd) && t.occurrenceCountLess(newInd, tsInd) {
		newInd--
	}
	t.swap(tsInd, newInd+1)
}

func (t *TagCloud) swap(ind1, ind2 int) {
	t.stats[ind1], t.stats[ind2] = t.stats[ind2], t.stats[ind1]
	t.tags[t.tag(ind1)], t.tags[t.tag(ind2)] = t.tags[t.tag(ind2)], t.tags[t.tag(ind1)]
}
func (t *TagCloud) tag(ind int) string {
	return t.stats[ind].Tag
}
func (t *TagCloud) validIndex(ind int) bool {
	return -1 < ind && ind < len(t.stats)
}
func (t *TagCloud) occurrenceCountLess(ind1, ind2 int) bool {
	return t.stats[ind1].OccurrenceCount < t.stats[ind2].OccurrenceCount
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
// TODO: You decide whether receiver should be a pointer or a value
func (t *TagCloud) TopN(n int) []TagStat {
	if n > len(t.stats) {
		return t.stats[:]
	} else {
		return t.stats[:n]
	}
}
