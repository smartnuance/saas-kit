package dbmodels

import (
	"github.com/smartnuance/saas-kit/pkg/lib/paging"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (w whereHelperstring) Page(page paging.Page) qm.QueryMod {
	switch spec := page.(type) {
	case *paging.Paging_First:
		return qm.Limit(int(spec.GetPageSize()))
	case *paging.Paging_Previous:
		return qm.Expr(w.LT(spec.End), qm.Limit(int(spec.GetPageSize())))
	case *paging.Paging_Current:
		return qm.Expr(w.GT(spec.Start), w.LT(spec.End), qm.Limit(int(spec.GetPageSize())))
	case *paging.Paging_Next:
		return qm.Expr(w.GT(spec.Start), qm.Limit(int(spec.PageSize)))
	default:
		// identity function not modifying query
		return qm.QueryModFunc(func(q *queries.Query) {})
	}
}
