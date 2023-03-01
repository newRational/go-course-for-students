package tagcloud

// TagCloud aggregates statistics about used tags
type TagCloud struct {

	// 1. Первое поле - tagNames map[string]*TagStat
	tagNames map[string]*TagStat

	// 2. Второе поле - tagStats []*TagStat
	tagStats []*TagStat

	// TODO: add fields if necessary
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
// TODO: You decide whether this function should return a pointer or a value
func New() TagCloud {
	// TODO: Implement this
	return TagCloud{}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
// TODO: You decide whether receiver should be a pointer or a value
func (TagCloud) AddTag(tag string) {
	// 1. 	Ищется тег в tagNames
	// 1a.	Тег есть в tagNames - идет переход на соответствующий элемент в tagStats
	// 1b.	Тега нет в tagNames - Добавляется новый тег в tagNames, затем идет
	// 		переход на соответствующий элемент в tagStats
	// 2.	В элементе в tagStats увеличивается на 1 OccurrenceCount
	// 3.	Производится "всплывание" данного элемента
	// TODO: Implement this
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
// TODO: You decide whether receiver should be a pointer or a value
func (TagCloud) TopN(n int) []TagStat {
	// Так как tagStats типа slice, то просто возвращаем slice из первых N элементов
	// TODO: Implement this
	return nil
}
