package udwType

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

type T struct {
	String1 string
	Map1    map[string]string
	Map2    map[string]*string
	Map3    map[string]T2
	Map4    map[string]map[string]string
	Map5    map[string][]string
	Slice1  []string
	Ptr1    *string
	Ptr2    *T2
	Array1  [5]string
}

type T2 struct {
	A string
	B string
}

func TestPtrType(ot *testing.T) {
	var data **string
	data = new(*string)
	m, err := NewContext(data)
	udwTest.Equal(err, nil)

	err = m.SaveByPath(Path{"ptr", "ptr"}, "")
	udwTest.Equal(err, nil)
	udwTest.Ok(data != nil)
	udwTest.Ok(*data != nil)
	udwTest.Equal(**data, "")
}

func TestStringType(ot *testing.T) {
	var data *string
	data = new(string)
	m, err := NewContext(data)
	udwTest.Equal(err, nil)

	err = m.SaveByPath(Path{"ptr"}, "123")
	udwTest.Equal(err, nil)
	udwTest.Ok(data != nil)
	udwTest.Equal(*data, "123")
}

func TestStructType(ot *testing.T) {
	data := &struct {
		A string
	}{}
	m, err := NewContext(data)
	udwTest.Equal(err, nil)

	err = m.SaveByPath(Path{"ptr", "A"}, "123")
	udwTest.Equal(err, nil)
	udwTest.Ok(data != nil)
	udwTest.Equal(data.A, "123")
}

func TestType(ot *testing.T) {
	data := &T{}
	m, err := NewContext(data)
	udwTest.Equal(err, nil)

	err = m.SaveByPath(Path{"ptr", "String1"}, "B")
	udwTest.Equal(err, nil)
	udwTest.Equal(data.String1, "B")

	m.SaveByPath(Path{"ptr", "Map1", "A"}, "1123")
	_, ok := data.Map1["A"]
	udwTest.Equal(ok, true)
	udwTest.Equal(data.Map1["A"], "1123")

	err = m.SaveByPath(Path{"ptr", "Map1", "A"}, "1124")
	udwTest.Equal(err, nil)
	udwTest.Equal(data.Map1["A"], "1124")

	err = m.DeleteByPath(Path{"ptr", "Map1", "A"})
	udwTest.Equal(err, nil)
	_, ok = data.Map1["A"]
	udwTest.Equal(ok, false)

	err = m.SaveByPath(Path{"ptr", "Map2", "B", "ptr"}, "1")
	udwTest.Equal(err, nil)
	rpString, ok := data.Map2["B"]
	udwTest.Equal(ok, true)
	udwTest.Equal(*rpString, "1")

	err = m.SaveByPath(Path{"ptr", "Map2", "B", "ptr"}, "2")
	udwTest.Equal(err, nil)
	udwTest.Equal(*rpString, "2")

	err = m.DeleteByPath(Path{"ptr", "Map2", "B", "ptr"})
	udwTest.Equal(err, nil)
	_, ok = data.Map2["B"]
	udwTest.Equal(ok, true)
	udwTest.Equal(data.Map2["B"], nil)

	err = m.DeleteByPath(Path{"ptr", "Map2", "B"})
	udwTest.Equal(err, nil)
	_, ok = data.Map2["B"]
	udwTest.Equal(ok, false)

	err = m.SaveByPath(Path{"ptr", "Map3", "C", "A"}, "1")
	udwTest.Equal(err, nil)
	udwTest.Equal(data.Map3["C"].A, "1")

	err = m.DeleteByPath(Path{"ptr", "Map3", "C"})
	udwTest.Equal(err, nil)
	udwTest.Ok(data.Map3 != nil)
	_, ok = data.Map3["C"]
	udwTest.Equal(ok, false)

	err = m.SaveByPath(Path{"ptr", "Map4", "D", "F"}, "1234")
	udwTest.Equal(err, nil)
	udwTest.Equal(data.Map4["D"]["F"], "1234")

	err = m.SaveByPath(Path{"ptr", "Map4", "D", "H"}, "12345")
	udwTest.Equal(err, nil)
	udwTest.Equal(data.Map4["D"]["H"], "12345")

	err = m.SaveByPath(Path{"ptr", "Map4", "D", "H"}, "12346")
	udwTest.Equal(err, nil)
	udwTest.Equal(data.Map4["D"]["H"], "12346")

	err = m.DeleteByPath(Path{"ptr", "Map4", "D", "F"})
	udwTest.Equal(err, nil)
	udwTest.Ok(data.Map4["D"] != nil)
	_, ok = data.Map4["D"]["F"]
	udwTest.Equal(ok, false)

	_, ok = data.Map4["D"]["H"]
	udwTest.Equal(ok, true)

	err = m.SaveByPath(Path{"ptr", "Map5", "D", ""}, "1234")
	udwTest.Equal(err, nil)
	udwTest.Equal(len(data.Map5["D"]), 1)
	udwTest.Equal(data.Map5["D"][0], "1234")

	err = m.DeleteByPath(Path{"ptr", "Map5", "D", "0"})
	udwTest.Equal(err, nil)
	udwTest.Equal(len(data.Map5["D"]), 0)

	err = m.SaveByPath(Path{"ptr", "Slice1", ""}, "1234")
	udwTest.Equal(err, nil)
	udwTest.Equal(len(data.Slice1), 1)
	udwTest.Equal(data.Slice1[0], "1234")

	err = m.SaveByPath(Path{"ptr", "Slice1", ""}, "12345")
	udwTest.Equal(err, nil)
	udwTest.Equal(data.Slice1[1], "12345")
	udwTest.Equal(len(data.Slice1), 2)

	err = m.DeleteByPath(Path{"ptr", "Slice1", "0"})
	udwTest.Equal(err, nil)
	udwTest.Equal(len(data.Slice1), 1)
	udwTest.Equal(data.Slice1[0], "12345")

	err = m.SaveByPath(Path{"ptr", "Ptr1", "ptr"}, "12345")
	udwTest.Equal(err, nil)
	udwTest.Equal(*data.Ptr1, "12345")

	err = m.SaveByPath(Path{"ptr", "Ptr2", "ptr"}, "")
	udwTest.Equal(err, nil)
	udwTest.Equal(data.Ptr2.A, "")

	err = m.SaveByPath(Path{"ptr", "Array1", "1"}, "12345")
	udwTest.Equal(err, nil)
	udwTest.Equal(data.Array1[1], "12345")
}
