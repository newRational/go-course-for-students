package tagcloud

type TagCloud struct {
	// frequencies отображает OccurrenceCount на индекс, на который нужно
	// переместить элемент типа TagStat c данным OccurrenceCount
	frequencies map[int]int

	// tags отображает Tag на индекс, по которому расположен
	// элемент типа TagStat с данным Tag
	tags map[string]int

	// stats - непосредственное хранилище элементов типа TagStat
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
func New() *TagCloud {
	return &TagCloud{
		frequencies: map[int]int{},
		tags:        map[string]int{},
		stats:       []TagStat{},
	}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
// Вставка выполняется за время O(1)
func (t *TagCloud) AddTag(tag string) {
	ind, ok := t.tags[tag]
	if ok {
		t.updateTag(ind)
	} else {
		t.addNewTag(tag)
	}
}
func (t *TagCloud) addNewTag(tag string) {
	t.tags[tag] = len(t.stats)
	t.stats = append(t.stats, NewTagStat(tag))
}

func (t *TagCloud) updateTag(ind int) {
	t.stats[ind].OccurrenceCount++
	f := t.stats[ind].OccurrenceCount
	t.swap(ind, t.frequencies[f])
	t.frequencies[f]++
}

func (t *TagCloud) swap(ind1, ind2 int) {
	t.stats[ind1], t.stats[ind2] = t.stats[ind2], t.stats[ind1]
	t.tags[t.tag(ind1)], t.tags[t.tag(ind2)] = t.tags[t.tag(ind2)], t.tags[t.tag(ind1)]
}

// tag returns tag by index
func (t *TagCloud) tag(ind int) string {
	return t.stats[ind].Tag
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (t *TagCloud) TopN(n int) []TagStat {
	if n > len(t.stats) {
		return t.stats[:]
	} else {
		return t.stats[:n]
	}
}
