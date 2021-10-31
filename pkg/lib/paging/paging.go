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

type Paging struct {
	Previous *PreviousSpec `json:"prev,omitempty"`
	Current  FullSpec      `json:"cur"`
	Next     *NextSpec     `json:"next,omitempty"`
}

type FirstSpec struct {
	PageSize int `json:"size"`
}

func (f *FirstSpec) Query(query *url.Values) {
	pageQuery(query, f)
}

func (f *FirstSpec) Size() int { return f.PageSize }

type PreviousSpec struct {
	EndIDExcl string `json:"end"`
	PageSize  int    `json:"size"`
}

func (f *PreviousSpec) Size() int { return f.PageSize }

func (f *PreviousSpec) Query(query *url.Values) {
	pageQuery(query, f)
}

type FullSpec struct {
	// StartIDIncl is the start ID (inclusive), might be empty if page contains no elements
	StartIDIncl string `json:"start,omitempty"`
	// EndIDIncl is the end ID (inclusive), might be empty if page contains no elements
	EndIDIncl string `json:"end,omitempty"`
	PageSize  int    `json:"size"`
}

func (f *FullSpec) Size() int { return f.PageSize }

type NextSpec struct {
	StartIDExcl string `json:"start"`
	PageSize    int    `json:"size"`
}

func (f *NextSpec) Size() int { return f.PageSize }

func (f *NextSpec) Query(query *url.Values) {
	pageQuery(query, f)
}

func (p *FullSpec) Next() NextSpec {
	return NextSpec{
		StartIDExcl: p.EndIDIncl,
		PageSize:    p.PageSize,
	}
}

func (p *FullSpec) Prev() PreviousSpec {
	return PreviousSpec{
		EndIDExcl: p.StartIDIncl,
		PageSize:  p.PageSize,
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
				StartIDIncl: start,
				EndIDIncl:   end,
				PageSize:    DefaultPageSize,
			}
		} else {
			// only start provided
			return &NextSpec{
				StartIDExcl: start,
				PageSize:    DefaultPageSize,
			}
		}
	} else if end != "" {
		// only end provided
		return &PreviousSpec{
			EndIDExcl: end,
			PageSize:  DefaultPageSize,
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
		query.Add(PageEndQueryParam, spec.EndIDExcl)
		query.Add(PageSizeQueryParam, strconv.Itoa(spec.PageSize))
	case *FullSpec:
		query.Add(PageStartQueryParam, spec.StartIDIncl)
		query.Add(PageEndQueryParam, spec.EndIDIncl)
		query.Add(PageSizeQueryParam, strconv.Itoa(spec.PageSize))
	case *NextSpec:
		query.Add(PageStartQueryParam, spec.StartIDExcl)
		query.Add(PageSizeQueryParam, strconv.Itoa(spec.PageSize))
	}
}
