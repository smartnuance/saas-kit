package models

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

func (f *FirstSpec) Query(query *url.Values) {
	pageQuery(query, f)
}

func (f *FirstSpec) Size() int { return f.PageSize }

func (f *PreviousSpec) Size() int { return f.PageSize }

func (f *PreviousSpec) Query(query *url.Values) {
	pageQuery(query, f)
}

func (f *FullSpec) Size() int { return f.PageSize }

func (f *NextSpec) Size() int { return f.PageSize }

func (f *NextSpec) Query(query *url.Values) {
	pageQuery(query, f)
}

func (p *FullSpec) Next() NextSpec {
	return NextSpec{
		Start:    p.End,
		PageSize: p.PageSize,
	}
}

func (p *FullSpec) Prev() PreviousSpec {
	return PreviousSpec{
		End:      p.Start,
		PageSize: p.PageSize,
	}
}

func (f *FullSpec) Query(query *url.Values) {
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
			return &FullSpec{
				Start:    start,
				End:      end,
				PageSize: DefaultPageSize,
			}
		} else {
			// only start provided
			return &NextSpec{
				Start:    start,
				PageSize: DefaultPageSize,
			}
		}
	} else if end != "" {
		// only end provided
		return &PreviousSpec{
			End:      end,
			PageSize: DefaultPageSize,
		}
	}
	return &FirstSpec{
		PageSize: size,
	}
}

func pageQuery(query *url.Values, p Page) {
	query.Del(PageStartQueryParam)
	query.Del(PageEndQueryParam)
	query.Del(PageSizeQueryParam)

	switch spec := p.(type) {
	case *FirstSpec:
		query.Add(PageSizeQueryParam, strconv.Itoa(spec.PageSize))
	case *PreviousSpec:
		query.Add(PageEndQueryParam, spec.End)
		query.Add(PageSizeQueryParam, strconv.Itoa(spec.PageSize))
	case *FullSpec:
		query.Add(PageStartQueryParam, spec.Start)
		query.Add(PageEndQueryParam, spec.End)
		query.Add(PageSizeQueryParam, strconv.Itoa(spec.PageSize))
	case *NextSpec:
		query.Add(PageStartQueryParam, spec.Start)
		query.Add(PageSizeQueryParam, strconv.Itoa(spec.PageSize))
	}
}
