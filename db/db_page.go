package db

import (
	"fmt"

	"github.com/artwebs/aogo/web"
)

type Page struct {
	page, pageSize, pageCount, total int
}

func PageNew(page, pageSize int) *Page {
	return &Page{page: page, pageSize: pageSize}
}

func PageNewWithCtx(ctx *web.Context) *Page{
	obj := &Page{}
  if page := ctx.Form["page"]{
    obj.page = strconv.Atoi(page)
  }
  if pageSize := ctx.Form["rows"]{
    obj.pageSize = strconv.Atoi(pageSize)
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
