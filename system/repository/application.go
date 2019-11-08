package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	ApplicationRepository interface {
		With(ctx context.Context, db *factory.DB) ApplicationRepository

		FindByID(id uint64) (*types.Application, error)
		Find(types.ApplicationFilter) (types.ApplicationSet, types.ApplicationFilter, error)

		Create(mod *types.Application) (*types.Application, error)
		Update(mod *types.Application) (*types.Application, error)

		DeleteByID(id uint64) error
	}

	application struct {
		*repository
	}
)

const (
	ErrApplicationNotFound = repositoryError("ApplicationNotFound")
)

// @todo migrate to same pattern as we have for users
func Application(ctx context.Context, db *factory.DB) ApplicationRepository {
	return (&application{}).With(ctx, db)
}

func (r *application) With(ctx context.Context, db *factory.DB) ApplicationRepository {
	return &application{
		repository: r.repository.With(ctx, db),
	}
}

func (r application) table() string {
	return "sys_application"
}

func (r application) columns() []string {
	return []string{
		"id",
		"rel_owner",
		"name",
		"enabled",
		"unify",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func (r application) query() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.table()).
		Where("deleted_at IS NULL")
}

func (r *application) FindByID(id uint64) (*types.Application, error) {
	return r.findOneBy("id", id)
}

func (r application) findOneBy(field string, value interface{}) (*types.Application, error) {
	var (
		app = &types.Application{}

		q = r.query().
			Where(squirrel.Eq{field: value})

		err = rh.FetchOne(r.db(), q, app)
	)

	if err != nil {
		return nil, err
	} else if app.ID == 0 {
		return nil, ErrApplicationNotFound
	}

	return app, nil
}

func (r *application) Find(filter types.ApplicationFilter) (set types.ApplicationSet, f types.ApplicationFilter, err error) {
	f = filter

	if f.Sort == "" {
		f.Sort = "id"
	}

	query := r.query()

	if f.IsReadable != nil {
		query = query.Where(f.IsReadable)
	}

	var orderBy []string
	if orderBy, err = rh.ParseOrder(f.Sort, r.columns()...); err != nil {
		return
	} else {
		query = query.OrderBy(orderBy...)
	}

	if f.Count, err = rh.Count(r.db(), query); err != nil || f.Count == 0 {
		return
	}

	return set, f, rh.FetchPaged(r.db(), query, f.Page, f.PerPage, &set)
}

func (r *application) Create(mod *types.Application) (*types.Application, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert(r.table(), mod)
}

func (r *application) Update(mod *types.Application) (*types.Application, error) {
	rh.SetCurrentTimeRounded(&mod.UpdatedAt)

	return mod, r.db().Replace(r.table(), mod)
}

func (r *application) DeleteByID(id uint64) error {
	return r.updateColumnByID(r.table(), "deleted_at", time.Now(), id)
}
