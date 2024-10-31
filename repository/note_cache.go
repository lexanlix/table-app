package repository

import (
	"sync"
	"time"

	"table-app/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type NoteCache struct {
	noteIndexById map[string]int
	notes         domain.NoteList

	mutex sync.Mutex
}

func NewNoteCache() *NoteCache {
	return &NoteCache{
		mutex: sync.Mutex{},
	}
}

func (r *NoteCache) SetData(notesData []domain.Note) {
	noteList := domain.NoteList(notesData)
	noteList.SortByDateAsc()

	r.noteIndexById = make(map[string]int)
	r.notes = make([]domain.Note, 0)

	for _, note := range noteList {
		r.noteIndexById[note.Id] = len(r.notes)
		r.notes = append(r.notes, note)
	}
}

func (r *NoteCache) Insert(newNote domain.Note) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	newNote.Id = uuid.New().String()
	newNote.UpdatedAt = time.Now()

	r.noteIndexById[newNote.Id] = len(r.notes)
	r.notes = append(r.notes, newNote)
}

func (r *NoteCache) ReadAll() []domain.Note {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.notes
}

func (r *NoteCache) GetListPtr() *domain.NoteList {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return &r.notes
}

func (r *NoteCache) IsInCache(noteId string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, ok := r.noteIndexById[noteId]
	return ok
}

func (r *NoteCache) UpdateNote(note domain.Note) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	note.UpdatedAt = time.Now()

	idx, ok := r.noteIndexById[note.Id]
	if !ok {
		return errors.Errorf("note %s not found", note.Id)
	}

	r.notes[idx] = note

	r.notes.SortByDateAsc()

	return nil
}

func (r *NoteCache) Delete(noteId string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	idx, ok := r.noteIndexById[noteId]
	if !ok {
		return
	}

	note := r.notes[idx]
	note.Deleted = true
	r.notes[idx] = note
}
