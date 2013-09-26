package main

import (
	"fmt"
	"os"
	"sort"
)

var (
	asc          = false
	desc         = true
	sortOrderMap = map[string]string{
		"mode": "asc",
		"time": "asc",
		"size": "asc",
		"name": "asc",
	}
)

type csvflag []string

func (s *csvflag) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *csvflag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type sortable struct {
	Infos *[]os.FileInfo
	Key   string
	Order bool
}

func xnor(a, b bool) bool { return !((a || b) && (!a || !b)) }

func (s sortable) Len() int { return len(*s.Infos) }
func (s sortable) Less(i, j int) bool {
	switch s.Key {
	case "mode":
		return xnor((*s.Infos)[i].Mode() > (*s.Infos)[j].Mode(), s.Order)
	case "time":
		return xnor((*s.Infos)[i].ModTime().After((*s.Infos)[j].ModTime()), s.Order)
	case "size":
		return xnor((*s.Infos)[i].Size() > (*s.Infos)[j].Size(), s.Order)
	default:
		return xnor((*s.Infos)[i].Name() > (*s.Infos)[j].Name(), s.Order)
	}
	return xnor((*s.Infos)[i].Name() > (*s.Infos)[j].Name(), s.Order)
}
func (s sortable) Swap(i, j int) { (*s.Infos)[i], (*s.Infos)[j] = (*s.Infos)[j], (*s.Infos)[i] }

func readDir(dirname string, sortby string, order bool) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Sort(sortable{&list, sortby, order})
	return list, nil
}
