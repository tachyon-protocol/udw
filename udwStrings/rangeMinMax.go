package udwStrings

type RangeMinMax struct {
	minKey           string
	minKeyNotInclude string
	maxKey           string
	maxKeyNotInclude string
}

func (kr *RangeMinMax) AddMinKeyInclude(minKey string) {
	if kr.minKey < minKey {
		kr.minKey = minKey
	}
}
func (kr *RangeMinMax) AddMinKeyNotInclude(minKeyNotInclude string) {
	if kr.minKeyNotInclude < minKeyNotInclude {
		kr.minKeyNotInclude = minKeyNotInclude
	}
}

func (kr *RangeMinMax) AddMaxKeyInclude(maxKey string) {
	if maxKey != "" && (kr.maxKey == "" || kr.maxKey > maxKey) {
		kr.maxKey = maxKey
	}
}

func (kr *RangeMinMax) AddMaxKeyNotInclude(maxKeyNotInclude string) {
	if maxKeyNotInclude != "" && (kr.maxKeyNotInclude == "" || kr.maxKeyNotInclude > maxKeyNotInclude) {
		kr.maxKeyNotInclude = maxKeyNotInclude
	}
}
func (kr *RangeMinMax) AddPrefix(prefix string) {
	if prefix == "" {
		return
	}
	kr.AddMinKeyInclude(prefix)
	kr.AddMaxKeyNotInclude(GetSmallestBiggerStringPrefix(prefix))
}

func (kr *RangeMinMax) ToMinKeyIncludeAndMaxKeyNotIncludeRange() RangeMinMax {
	return RangeMinMax{
		minKey:           kr.GetMinKeyInclude(),
		maxKeyNotInclude: kr.GetMaxKeyNotInclude(),
	}
}

func (kr *RangeMinMax) IsInRange(s string) bool {
	if (s >= kr.minKey) == false {
		return false
	}
	if (s > kr.minKeyNotInclude) == false {
		return false
	}
	if kr.maxKey != "" && (s <= kr.maxKey) == false {
		return false
	}
	if kr.maxKeyNotInclude != "" && (s < kr.maxKeyNotInclude) == false {
		return false
	}
	return true
}

func (kr *RangeMinMax) HasRange() bool {
	return kr.minKey != "" || kr.minKeyNotInclude != "" || kr.maxKey != "" || kr.maxKeyNotInclude != ""
}

func (kr *RangeMinMax) GetMinKeyInclude() string {
	minKey := kr.minKey
	if kr.minKeyNotInclude != "" {
		thisMinKey := stringSmallestBigger(kr.minKeyNotInclude)
		if minKey < thisMinKey {
			minKey = thisMinKey
		}
	}
	return minKey
}

func (kr *RangeMinMax) GetMaxKeyNotInclude() string {
	maxKeyNotInclude := kr.maxKeyNotInclude
	if kr.maxKey != "" {
		thisMaxKeyNotInclude := stringSmallestBigger(kr.maxKey)
		if maxKeyNotInclude == "" || maxKeyNotInclude > thisMaxKeyNotInclude {
			maxKeyNotInclude = thisMaxKeyNotInclude
		}
	}
	return maxKeyNotInclude
}

func stringSmallestBigger(s string) string {
	return string(append([]byte(s), 0))
}
