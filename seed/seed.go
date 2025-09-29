package seed

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Runner struct {
	DB  *gorm.DB
	Env string // "prod" | "dev"
}

func (r *Runner) Run(ctx context.Context) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := seedRoles(tx); err != nil {
			err = errors.Wrap(err, "[seed.Run]: Error seeding roles")
			return err
		}
		if err := seedPermissions(tx); err != nil {
			err = errors.Wrap(err, "[seed.Run]: Error seeding permissions")
			return err
		}
		if err := seedRolePermissions(tx); err != nil {
			err = errors.Wrap(err, "[seed.Run]: Error seeding role permissions")
			return err
		}
		if err := seedAdminUser(tx); err != nil {
			err = errors.Wrap(err, "[seed.Run]: Error seeding admin user")
			return err
		}
		if err := seedAdminUserRoles(tx); err != nil {
			err = errors.Wrap(err, "[seed.Run]: Error seeding admin user roles")
			return err
		}
		if err := seedUsers(tx); err != nil {
			err = errors.Wrap(err, "[seed.Run]: Error seeding users")
			return err
		}
		if err := seedUserRoles(tx); err != nil {
			err = errors.Wrap(err, "[seed.Run]: Error seeding user roles")
			return err
		}

		if r.Env == "development" {
			if err := seedAreas(tx); err != nil {
				err = errors.Wrap(err, "[seed.Run]: Error seeding areas")
				return err
			}
			if err := seedTables(tx); err != nil {
				err = errors.Wrap(err, "[seed.Run]: Error seeding tables")
				return err
			}
			if err := seedCategories(tx); err != nil {
				err = errors.Wrap(err, "[seed.Run]: Error seeding categories")
				return err
			}
			if err := seedMenuItems(tx); err != nil {
				err = errors.Wrap(err, "[seed.Run]: Error seeding menu items")
				return err
			}
			if err := seedModifiers(tx); err != nil {
				err = errors.Wrap(err, "[seed.Run]: Error seeding modifiers")
				return err
			}
			if err := seedSampleOrders(tx); err != nil {
				err = errors.Wrap(err, "[seed.Run]: Error seeding sample orders")
				return err
			}
		}
		return nil
	})
}
