package repository

import "iter"

type Repository[T any, ID any] interface {
	Save(entity *T) (*T, error)
	SaveAll(entities iter.Seq[*T]) iter.Seq2[*T, error]
	FindByID(id ID) (*T, error)
	ExistsByID(id ID) (bool, error)
	FindAll() iter.Seq2[*T, error]
	FindAllByID(ids iter.Seq[ID]) iter.Seq2[*T, error]
	Count() (int64, error)
	DeleteByID(id ID) error
	Delete(entity *T) error
	DeleteAllByID(ids iter.Seq[ID]) error
	DeleteAll() error
}
