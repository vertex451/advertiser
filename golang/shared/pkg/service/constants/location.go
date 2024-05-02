package constants

import "strconv"

type Location int

func (l Location) String() string {
	return strconv.Itoa(int(l))
}

const (
	Ukraine Location = iota

	Dnipro
	Kharkiv
	Kyiv
	Lviv
	Odessa
)

var Locations = map[Location]string{
	Ukraine: "Україна",

	Dnipro:  "Дніпро",
	Kharkiv: "Харків",
	Kyiv:    "Київ",
	Lviv:    "Львів",
	Odessa:  "Одеса",
}
