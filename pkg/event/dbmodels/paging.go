package dbmodels

import (
	"github.com/smartnuance/saas-kit/pkg/lib/paging"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (w whereHelperstring) Page(page paging.Page) qm.QueryMod {
	switch spec := page.(type) {
	case *paging.FirstSpec:
		return qm.Limit(spec.PageSize)
	case *paging.PreviousSpec:
		return qm.Expr(w.LT(spec.EndIDExcl), qm.Limit(spec.PageSize))
	case *paging.FullSpec:
		return qm.Expr(w.GT(spec.StartIDIncl), w.LT(spec.EndIDIncl), qm.Limit(spec.PageSize))
	case *paging.NextSpec:
		return qm.Expr(w.GT(spec.StartIDExcl), qm.Limit(spec.PageSize))
	default:
		// identity function not modifying query
		return qm.QueryModFunc(func(q *queries.Query) {})
	}
}
