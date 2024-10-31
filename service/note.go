package service

import (
	"context"

	"table-app/domain"
	"table-app/internal/log"
	"table-app/repository"

	"github.com/pkg/errors"
)

type NoteRepository interface {
	UpsertAll(ctx context.Context, notes []domain.Note) error
	GetAll(ctx context.Context) ([]domain.Note, error)
}

type Note struct {
	logger log.Logger
	repo   NoteRepository
	cache  *repository.NoteCache
}

func NewNote(logger log.Logger, repo NoteRepository, cache *repository.NoteCache) *Note {
	return &Note{
		logger: logger,
		repo:   repo,
		cache:  cache,
	}
}

func (s *Note) SaveAll(ctx context.Context) error {
	all := s.cache.ReadAll()

	err := s.repo.UpsertAll(ctx, all)
	if err != nil {
		return errors.WithMessage(err, "upsert all")
	}

	return nil
}

func (s *Note) GetListPtr(ctx context.Context) (*domain.NoteList, error) {
	noteData, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get all notes")
	}

	s.cache.SetData(noteData)

	return s.cache.GetListPtr(), nil
}

func (s *Note) Insert(ctx context.Context, newNote domain.Note) {
	s.cache.Insert(newNote)
}

func (s *Note) Update(ctx context.Context, note domain.Note) error {
	err := s.cache.UpdateNote(note)
	if err != nil {
		return errors.WithMessage(err, "update note")
	}

	return nil
}

func (s *Note) Delete(ctx context.Context, noteId string) {
	s.cache.Delete(noteId)
}
