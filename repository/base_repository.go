package repository

import (
	"errors"
	"iter"
	"slices"

	"gorm.io/gorm"
)

type BaseRepository[T any, ID any] struct {
	db *gorm.DB
}

var _ Repository[any, any] = (*BaseRepository[any, any])(nil)

func (b *BaseRepository[T, ID]) Count() (int64, error) {
	var count int64
	err := b.db.Model(new(T)).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (b *BaseRepository[T, ID]) Delete(entity *T) error {
	return b.db.Delete(entity).Error
}

func (b *BaseRepository[T, ID]) DeleteAll() error {
	return b.db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(new(T)).Error
}

func (b *BaseRepository[T, ID]) DeleteAllByID(ids iter.Seq[ID]) error {
	idsSlice := slices.Collect(ids)
	if len(idsSlice) == 0 {
		return nil
	}

	return b.db.Delete(new(T), idsSlice).Error
}

func (b *BaseRepository[T, ID]) DeleteByID(id ID) error {
	return b.db.Delete(new(T), id).Error
}

func (b *BaseRepository[T, ID]) ExistsByID(id ID) (bool, error) {
	var count int64
	err := b.db.Model(new(T)).
		Where("id = ?", id). // FIXME: can't assume the PK is always "id"
		Limit(1).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (b *BaseRepository[T, ID]) FindAll() iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		var entities []T
		if err := b.db.Find(&entities).Error; err != nil {
			yield(nil, err)
			return
		}

		for i := range entities {
			if !yield(&entities[i], nil) {
				return
			}
		}
	}
}

func (b *BaseRepository[T, ID]) FindAllByID(ids iter.Seq[ID]) iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		idsSlice := slices.Collect(ids)
		if len(idsSlice) == 0 {
			return
		}

		var entities []T
		if err := b.db.Find(&entities, idsSlice).Error; err != nil {
			yield(nil, err)
			return
		}

		// FIXME: entities are not guaranteed to be in the same order as idsSlice
		for i := range entities {
			if !yield(&entities[i], nil) {
				return
			}
		}
	}
}

func (b *BaseRepository[T, ID]) FindByID(id ID) (*T, error) {
	var user T
	if err := b.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (b *BaseRepository[T, ID]) Save(entity *T) (*T, error) {
	if err := b.db.Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

// TODO: think if it makes sense to have a batch save with transaction
func (b *BaseRepository[T, ID]) SaveAll(entities iter.Seq[*T]) iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		for e := range entities {
			savedEntity, err := b.Save(e)
			if err != nil {
				yield(nil, err)
				return
			}
			if !yield(savedEntity, nil) {
				return
			}
		}
	}
}
