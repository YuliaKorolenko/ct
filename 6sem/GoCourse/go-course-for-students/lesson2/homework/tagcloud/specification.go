package tagcloud

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	container []TagStat
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
// TODO: You decide whether this function should return a pointer or a value
func New() TagCloud {
	return TagCloud{}
}

// return index, if tag in cloud, else return -1
func find(slice []TagStat, tag string) int {
	for j, v := range slice {
		if v.Tag == tag {
			return j
		}
	}
	return -1
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
// TODO: You decide whether receiver should be a pointer or a value
func (cloud *TagCloud) AddTag(tag string) {
	i := find(cloud.container, tag)
	if i == -1 {
		cloud.container = append(cloud.container, TagStat{Tag: tag, OccurrenceCount: 1})
		return
	}
	cloud.container[i].OccurrenceCount++
	var newValue int = cloud.container[i].OccurrenceCount
	for i > 0 && cloud.container[i-1].OccurrenceCount < newValue {
		cloud.container[i-1], cloud.container[i] = cloud.container[i], cloud.container[i-1]
		i--
	}
}

func boarder(amount int, size int) int {
	if amount < size {
		return amount
	}
	return size
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
// TODO: You decide whether receiver should be a pointer or a value
func (cloud TagCloud) TopN(n int) []TagStat {
	boarder := boarder(n, len(cloud.container))
	return cloud.container[:boarder]
}
