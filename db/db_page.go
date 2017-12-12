package db

import (
	"fmt"
	"strconv"
)

type Page struct {
	page, pageSize, pageCount, total int
}

func PageNew(page, pageSize int) *Page {
	return &Page{page: page, pageSize: pageSize}
}

func PageNewWithCtx(m map[string]interface{}) *Page {
	obj := &Page{}
	if page, ok := m["page"]; ok {
		val, _ := strconv.Atoi(page.(string))
		obj.page = val
	}
	if pageSize, ok := m["rows"]; ok {
		val, _ := strconv.Atoi(pageSize.(string))
		obj.pageSize = val
	}
	return obj
}

func (this *Page) GetPage() int {
	if this.page < 1 {
		this.page = 1
	}
	return this.page
}

func (this *Page) GetPageSize() int {
	return this.pageSize
}

func (this *Page) SetPageSize(pageSize int) *Page {
	this.pageSize = pageSize
	return this
}

func (this *Page) GetPageCount() int {
	if this.pageCount < 1 {
		return 1
	} else {
		return this.pageCount
	}
}

func (this *Page) GetTotal() int {
	return this.total
}

func (this *Page) GetPrev() int {
	if this.page-1 < 1 {
		return 1
	} else {
		return this.page - 1
	}
}

func (this *Page) GetNext() int {
	if this.page+1 > this.GetPageCount() {
		return this.GetPageCount()
	} else {
		return this.page + 1
	}
}

func (this *Page) ToLimit() string {
	return fmt.Sprintf("%d,%d", (this.page-1)*this.pageSize, this.pageSize)
}

func (this *Page) SetTotal(total int) *Page {
	this.total = total
	return this
}
