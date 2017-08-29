package main

import (
//"container/list"
)

type Array struct {
	Lst []int
}

func (a *Array) Add(el int) {
	a.Lst = append(a.Lst, el)
}

func (a Array) ForEach(f func(int) bool) {
	for i := 0; i < len(a.Lst); i++ {
		if !f(a.Lst[i]) {
			return
		}
	}
}

func (a *Array) Remove(el int) {
	for i := 0; i < len(a.Lst); i++ {
		if a.Lst[i] == el {
			a.Lst = append(a.Lst[:i], a.Lst[i+1:]...)
		}
	}
}

func NewArray() Array {
	return Array{}
}

//1000046*52 + 763407*44 + 10000460*20

/*
type Array struct {
	Lst *list.List
}

func (a Array) Add(el int) {
	a.Lst.PushBack(el)
}

func (a Array) ForEach(f func(int) bool) {
	for e := a.Lst.Front(); e != nil; e = e.Next() {
		if !f(e.Value.(int)) {
			break
		}
	}
}

func (a Array) Remove(el int) {
	for e := a.Lst.Front(); e != nil; e = e.Next() {
		if e.Value.(int) == el {
			a.Lst.Remove(e)
		}
	}
}

func NewArray() Array {
	return Array{
		Lst: list.New(),
	}
}
*/
