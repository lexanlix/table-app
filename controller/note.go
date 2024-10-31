package controller

import (
	"context"

	"table-app/domain"
	"table-app/internal/log"
)

type NoteService interface {
	SaveAll(ctx context.Context) error
	GetListPtr(ctx context.Context) (*domain.NoteList, error)
	Insert(ctx context.Context, newNote domain.Note)
	Update(ctx context.Context, note domain.Note) error
	Delete(ctx context.Context, noteId string)
}

type Note struct {
	logger  log.Logger
	service NoteService
}

func NewNote(logger log.Logger, service NoteService) Note {
	return Note{
		logger:  logger,
		service: service,
	}
}

func (c Note) SaveAll(ctx context.Context) error {
	c.logger.Debug(ctx, "save all notes")

	return c.service.SaveAll(ctx)
}

func (c Note) GetListPtr(ctx context.Context) (*domain.NoteList, error) {
	c.logger.Debug(ctx, "get notes list pointer")

	return c.service.GetListPtr(ctx)
}

func (c Note) AddNote(ctx context.Context, note domain.Note) {
	c.logger.Debug(ctx, "add new note", log.String("text", note.Text))

	c.service.Insert(ctx, note)
}

func (c Note) UpdateNote(ctx context.Context, note domain.Note) error {
	c.logger.Debug(ctx, "update note", log.String("text", note.Text))

	return c.service.Update(ctx, note)
}

func (c Note) DeleteNote(ctx context.Context, noteId string) {
	c.logger.Debug(ctx, "delete note", log.String("noteId", noteId))

	c.service.Delete(ctx, noteId)
}
