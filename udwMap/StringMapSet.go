package udwMap

type StringMapSet struct {
	data map[string]struct{}
}

func NewStringMapSet() StringMapSet {
	return StringMapSet{
		data: map[string]struct{}{},
	}
}

func (this StringMapSet) Contains(s string) bool {
	_, ok := this.data[s]
	return ok
}

func (this StringMapSet) Set(s string) {
	this.data[s] = struct{}{}
	return
}

func (this StringMapSet) Delete(s string) {
	delete(this.data, s)
}

func (this StringMapSet) InsertSlice(ss []string) {
	for _, s := range ss {
		this.data[s] = struct{}{}
	}
}

func (this StringMapSet) Len() int {
	return len(this.data)
}

func (this StringMapSet) GetStringSlice() []string {
	ss := make([]string, 0, this.Len())
	for k := range this.data {
		ss = append(ss, k)
	}
	return ss
}

func (this StringMapSet) GetInnerMap() map[string]struct{} {
	return this.data
}

func (this *StringMapSet) InsertMap(m map[string]struct{}) {
	for key := range m {
		this.data[key] = struct{}{}
	}
}

func (this StringMapSet) Clear() {
	for k := range this.data {
		delete(this.data, k)
	}
}

func (sms StringMapSet) IsZero() bool {
	return sms.data == nil
}
