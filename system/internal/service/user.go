package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	internalAuth "github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/types"
)

const (
	ErrUserInvalidCredentials = serviceError("UserInvalidCredentials")
	ErrUserLocked             = serviceError("UserLocked")

	uuidLength = 36
)

type (
	user struct {
		db  *factory.DB
		ctx context.Context

		prm PermissionsService

		user repository.UserRepository
	}

	UserService interface {
		With(ctx context.Context) UserService

		FindByUsername(username string) (*types.User, error)
		FindByEmail(email string) (*types.User, error)
		FindByID(id uint64) (*types.User, error)
		FindByIDs(id ...uint64) (types.UserSet, error)
		Find(filter *types.UserFilter) (types.UserSet, error)

		Create(input *types.User) (*types.User, error)
		Update(mod *types.User) (*types.User, error)

		Delete(id uint64) error
		Suspend(id uint64) error
		Unsuspend(id uint64) error

		// ValidateCredentials(username, password string) (*types.User, error)
	}
)

func User(ctx context.Context) UserService {
	return (&user{}).With(ctx)
}

func (svc *user) With(ctx context.Context) UserService {
	db := repository.DB(ctx)

	return &user{
		db:   db,
		ctx:  ctx,
		prm:  Permissions(ctx),
		user: repository.User(ctx, db),
	}
}

func (svc *user) FindByID(id uint64) (*types.User, error) {
	return svc.user.FindByID(id)
}

func (svc *user) FindByIDs(ids ...uint64) (types.UserSet, error) {
	return svc.user.FindByIDs(ids...)
}

func (svc *user) FindByEmail(email string) (*types.User, error) {
	return svc.user.FindByEmail(email)
}

func (svc *user) FindByUsername(username string) (*types.User, error) {
	return svc.user.FindByUsername(username)
}

func (svc *user) Find(filter *types.UserFilter) (types.UserSet, error) {
	return svc.user.Find(filter)
}

func (svc *user) Create(input *types.User) (out *types.User, err error) {
	return out, svc.db.Transaction(func() error {
		if !svc.prm.CanCreateUser() {
			return errors.New("not allowed to create users")
		}

		if out, err = svc.user.Create(input); err != nil {
			return err
		}

		return nil
	})
}

func (svc *user) Update(mod *types.User) (u *types.User, err error) {
	return u, svc.db.Transaction(func() (err error) {
		if u, err = svc.user.FindByID(mod.ID); err != nil {
			return
		}

		if mod.ID != internalAuth.GetIdentityFromContext(svc.ctx).Identity() && !svc.prm.CanUpdateUser(u) {
			return errors.New("not allowed to update this user")
		}

		// Assign changed values
		u.Email = mod.Email
		u.Username = mod.Username
		u.Name = mod.Name
		u.Handle = mod.Handle
		u.Kind = mod.Kind

		if u, err = svc.user.Update(u); err != nil {
			return err
		}

		return nil
	})
}

func (svc *user) Delete(id uint64) error {
	return svc.db.Transaction(func() (err error) {
		var u *types.User
		if u, err = svc.user.FindByID(id); err != nil {
			return
		}

		if !svc.prm.CanDeleteUser(u) {
			return errors.New("not allowed to update this user")
		}

		return svc.user.DeleteByID(id)
	})
}

func (svc *user) Suspend(id uint64) error {
	return svc.db.Transaction(func() (err error) {
		var u *types.User
		if u, err = svc.user.FindByID(id); err != nil {
			return
		}

		if !svc.prm.CanSuspendUser(u) {
			return errors.New("not allowed to update this user")
		}

		return svc.user.SuspendByID(id)
	})
}

func (svc *user) Unsuspend(id uint64) error {
	return svc.db.Transaction(func() (err error) {
		var u *types.User
		if u, err = svc.user.FindByID(id); err != nil {
			return
		}

		if !svc.prm.CanUnsuspendUser(u) {
			return errors.New("not allowed to update this user")
		}

		return svc.user.UnsuspendByID(id)
	})
}
