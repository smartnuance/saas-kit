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
	isPage() struct{}
}

type Paging struct {
	Previous *PreviousSpec `json:"prev,omitempty"`
	Current  FullSpec      `json:"cur"`
	Next     *NextSpec     `json:"next,omitempty"`
}

type FirstSpec struct {
	PageSize int `json:"size"`
}

func (f *FirstSpec) isPage() struct{} { return struct{}{} }

type PreviousSpec struct {
	EndIDExcl string `json:"end"`
	PageSize  int    `json:"size"`
}

func (f *PreviousSpec) isPage() struct{} { return struct{}{} }

type FullSpec struct {
	// StartIDIncl is the start ID (inclusive), might be empty if page contains no elements
	StartIDIncl string `json:"start,omitempty"`
	// EndIDIncl is the end ID (inclusive), might be empty if page contains no elements
	EndIDIncl string `json:"end,omitempty"`
	PageSize  int    `json:"size"`
}

func (f *FullSpec) isPage() struct{} { return struct{}{} }

type NextSpec struct {
	StartIDExcl string `json:"start"`
	PageSize    int    `json:"size"`
}

func (f *NextSpec) isPage() struct{} { return struct{}{} }

func (p FullSpec) Next() NextSpec {
	return NextSpec{
		StartIDExcl: p.EndIDIncl,
		PageSize:    p.PageSize,
	}
}

func (p FullSpec) Prev() PreviousSpec {
	return PreviousSpec{
		EndIDExcl: p.StartIDIncl,
		PageSize:  p.PageSize,
	}
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

func AddQuery(u *url.URL, p Page) {
	u.Query().Del(PageStartQueryParam)
	u.Query().Del(PageEndQueryParam)
	u.Query().Del(PageSizeQueryParam)

	switch spec := p.(type) {
	case *FirstSpec:
		u.Query().Add(PageSizeQueryParam, strconv.Itoa(spec.PageSize))
	case *PreviousSpec:
		u.Query().Add(PageEndQueryParam, spec.EndIDExcl)
		u.Query().Add(PageSizeQueryParam, strconv.Itoa(spec.PageSize))
	case *FullSpec:
		u.Query().Add(PageStartQueryParam, spec.StartIDIncl)
		u.Query().Add(PageEndQueryParam, spec.EndIDIncl)
		u.Query().Add(PageSizeQueryParam, strconv.Itoa(spec.PageSize))
	case *NextSpec:
		u.Query().Add(PageStartQueryParam, spec.StartIDExcl)
		u.Query().Add(PageSizeQueryParam, strconv.Itoa(spec.PageSize))
	}
}
