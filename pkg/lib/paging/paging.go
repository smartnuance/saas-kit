package paging

import (
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	DefaultPageSize     = 10
	PageStartQueryParam = "s"
	PageEndQueryParam   = "e"
	PageSizeQueryParam  = "l"
)

type Page interface {
	Size() int
	Query(query *url.Values)
}

func (f *Paging_First) Query(query *url.Values) {
	pageQuery(query, f)
}

func (f *Paging_First) Size() int { return int(f.GetPageSize()) }

func (f *Paging_Previous) Size() int { return int(f.GetPageSize()) }

func (f *Paging_Previous) Query(query *url.Values) {
	pageQuery(query, f)
}

func (f *Paging_Current) Size() int { return int(f.GetPageSize()) }

func (f *Paging_Next) Size() int { return int(f.GetPageSize()) }

func (f *Paging_Next) Query(query *url.Values) {
	pageQuery(query, f)
}

func (p *Paging_Current) Next() Paging_Next {
	return Paging_Next{
		Start:    p.End,
		PageSize: p.PageSize,
	}
}

func (p *Paging_Current) Prev() Paging_Previous {
	return Paging_Previous{
		End:      p.Start,
		PageSize: p.PageSize,
	}
}

func (f *Paging_Current) Query(query *url.Values) {
	pageQuery(query, f)
}

func FromQuery(c *gin.Context) Page {
	start := c.Query(PageStartQueryParam)
	end := c.Query(PageEndQueryParam)
	size, err := strconv.Atoi(c.Query(PageSizeQueryParam))
	if err != nil || size <= 0 {
		size = DefaultPageSize
	}

	if start != "" {
		if end != "" {
			// both provided
			return &Paging_Current{
				Start:    start,
				End:      end,
				PageSize: int32(DefaultPageSize),
			}
		} else {
			// only start provided
			return &Paging_Next{
				Start:    start,
				PageSize: int32(DefaultPageSize),
			}
		}
	} else if end != "" {
		// only end provided
		return &Paging_Previous{
			End:      end,
			PageSize: int32(DefaultPageSize),
		}
	}
	return &Paging_First{
		PageSize: int32(size),
	}
}

func pageQuery(query *url.Values, p Page) {
	query.Del(PageStartQueryParam)
	query.Del(PageEndQueryParam)
	query.Del(PageSizeQueryParam)

	switch spec := p.(type) {
	case *Paging_First:
		query.Add(PageSizeQueryParam, strconv.Itoa(spec.Size()))
	case *Paging_Previous:
		query.Add(PageEndQueryParam, spec.End)
		query.Add(PageSizeQueryParam, strconv.Itoa(int(spec.GetPageSize())))
	case *Paging_Current:
		query.Add(PageStartQueryParam, spec.Start)
		query.Add(PageEndQueryParam, spec.End)
		query.Add(PageSizeQueryParam, strconv.Itoa(int(spec.GetPageSize())))
	case *Paging_Next:
		query.Add(PageStartQueryParam, spec.Start)
		query.Add(PageSizeQueryParam, strconv.Itoa(int(spec.GetPageSize())))
	}
}
