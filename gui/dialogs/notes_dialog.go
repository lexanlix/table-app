package dialogs

import (
	"context"
	"image/color"
	"slices"
	"strconv"
	"time"

	"table-app/domain"
	"table-app/gui/iface"
	custom "table-app/gui/styles/colors"
	"table-app/internal/log"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/states"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/tree"
	"github.com/pkg/errors"
)

type NotesDialog struct {
	logger             log.Logger
	appBody            *core.Body
	dialogBody         *core.Body
	listFrame          *core.Frame
	noteCtrl           iface.NoteController
	noteList           *domain.NoteList
	deletedIds         []string
	hiddenFramesByYear map[int]map[string]*core.Frame
}

func NewNotesDialog(logger log.Logger, appBody *core.Body, noteCtrl iface.NoteController) (*NotesDialog, error) {
	dialogBody := core.NewBody("Notes").SetTitle("Заметки")
	dialogBody.Styler(func(s *styles.Style) {
		s.Align.Self = styles.Center
		s.Min.X.Dp(800)
		s.Min.Y.Dp(600)
		s.CenterAll()
	})

	dialogMainFrame := core.NewFrame(dialogBody)
	dialogMainFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.CenterAll()
	})

	titleFrame := core.NewFrame(dialogMainFrame)
	titleFrame.SetName("titleFrame")
	titleFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(400)
		s.CenterAll()
	})
	core.NewText(titleFrame).
		SetType(core.TextHeadlineMedium).
		SetText("Заметки")

	noteList, err := noteCtrl.GetListPtr(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, "get list pointer")
	}

	notesDialog := &NotesDialog{
		logger:             logger,
		appBody:            appBody,
		dialogBody:         dialogBody,
		noteCtrl:           noteCtrl,
		noteList:           noteList,
		hiddenFramesByYear: make(map[int]map[string]*core.Frame),
	}

	notesDialog.drawList(dialogMainFrame)
	notesDialog.drawControlField(dialogMainFrame)

	return notesDialog, nil
}

func (s *NotesDialog) drawList(dialogMainFrame *core.Frame) {
	s.listFrame = core.NewFrame(dialogMainFrame)
	s.listFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Min.X.Dp(800)
		s.Border.Width.SetAll(units.Dp(1))
		s.Gap.Zero()
		s.CenterAll()
	})

	s.listFrame.Maker(func(p *tree.Plan) {
		tree.AddAt(p, "head", func(noteFrame *core.Frame) {
			s.drawHeadRow(noteFrame)
		})

		var lastYear int
		for i := range *s.noteList {
			if (*s.noteList)[i].Deleted {
				continue
			}

			year := (*s.noteList)[i].UpdatedAt.Year()

			if year != time.Now().Year() {
				if (*s.noteList)[i].UpdatedAt.Year() != lastYear {
					tree.AddAt(p, "note_year_"+strconv.Itoa(year), func(noteFrame *core.Frame) {
						s.drawHidden(noteFrame, year)
					})
					lastYear = (*s.noteList)[i].UpdatedAt.Year()
				}
			}

			s.drawNoteRow(p, i)
		}
	})
}

func (s *NotesDialog) drawHeadRow(noteFrame *core.Frame) {
	noteFrame.Styler(func(s *styles.Style) {
		s.Align.Content = styles.Center
		s.Border.Width.SetAll(units.Dp(1))
		s.Min.X.Dp(1220)
		s.Gap.Zero()
		s.Gap.X.Dp(2)
	})

	core.NewSpace(noteFrame).Styler(func(s *styles.Style) {
		s.Min.X.Dp(5)
	})

	frameDate := core.NewFrame(noteFrame)
	frameDate.Styler(func(s *styles.Style) {
		s.Min.X.Dp(150)
		s.Min.Y.Dp(30)
		s.CenterAll()
	})
	core.NewText(frameDate).SetText("Дата обновления")
}

func (s *NotesDialog) drawHidden(noteFrame *core.Frame, year int) {
	noteFrame.Styler(func(s *styles.Style) {
		s.Align.Content = styles.Center
		s.Min.X.Dp(1220)
		s.Gap.Zero()
		s.Gap.X.Dp(2)
		s.Border.Width.SetAll(units.Dp(1))
		s.Align.Items = styles.Center
	})

	core.NewSpace(noteFrame).Styler(func(s *styles.Style) {
		s.Min.X.Dp(5)
	})

	// Поле даты обновления
	dateFrame := core.NewFrame(noteFrame)
	dateFrame.Styler(func(s *styles.Style) {
		s.Min.X.Dp(150)
		s.Max.Y.Dp(50)
		s.CenterAll()
	})

	dateText := core.NewText(dateFrame)
	dateText.SetText(strconv.Itoa(year))

	core.NewSpace(noteFrame).Styler(func(s *styles.Style) {
		s.Min.X.Dp(415)
	})

	// Кнопка развернуть заметки этого года
	expandButton := core.NewButton(noteFrame).
		SetType(core.ButtonElevated).
		SetIcon(icons.ArrowDropDown)

	expandButton.Styler(func(s *styles.Style) {
		s.Color = colors.Uniform(color.Black)
	})

	expandButton.OnClick(func(e events.Event) {
		hiddenFrames, ok := s.hiddenFramesByYear[year]
		if !ok {
			return
		}

		for _, frame := range hiddenFrames {
			frame.Styler(func(s *styles.Style) {
				if s.State.Is(states.Invisible) {
					s.SetState(false, states.Invisible)
				} else {
					s.SetState(true, states.Invisible)
				}
			})
		}

		if expandButton.Icon == icons.ArrowDropDown {
			expandButton.SetIcon(icons.ArrowDropUp)
		} else {
			expandButton.SetIcon(icons.ArrowDropDown)
		}

		s.listFrame.Update()
	})
}

func (s *NotesDialog) drawNoteRow(p *tree.Plan, idx int) {
	tree.AddAt(p, "note_"+strconv.Itoa(idx), func(noteFrame *core.Frame) {
		noteId := (*s.noteList)[idx].Id

		noteFrame.Styler(func(st *styles.Style) {
			st.Align.Content = styles.Center
			st.Min.X.Dp(1220)
			st.Gap.Zero()
			st.Gap.X.Dp(2)
			st.CenterAll()
			st.Border.Width.SetAll(units.Dp(1))

			year := (*s.noteList)[idx].UpdatedAt.Year()
			if year != time.Now().Year() {
				st.Background = custom.ColorSoftGrey

				framesById, ok := s.hiddenFramesByYear[year]
				if !ok {
					s.hiddenFramesByYear[year] = map[string]*core.Frame{
						noteId: noteFrame,
					}
				} else {
					framesById[noteId] = noteFrame
					s.hiddenFramesByYear[year] = framesById
				}

				st.SetState(true, states.Invisible)
			}
		})

		// Поле даты обновления
		dateFrame := core.NewFrame(noteFrame)
		dateFrame.Styler(func(s *styles.Style) {
			s.Min.X.Dp(150)
			s.Max.Y.Dp(30)
			s.CenterAll()
		})

		dateText := core.NewText(dateFrame)
		dateText.SetText((*s.noteList)[idx].UpdatedAt.Format("02.01.2006"))

		// Поле текста заметки
		fieldText := core.NewTextField(noteFrame)
		fieldText.Styler(func(s *styles.Style) {
			s.Min.X.Dp(950)
		})
		fieldText.SetText((*s.noteList)[idx].Text)
		fieldText.OnFocusLost(func(e events.Event) {
			note := (*s.noteList)[idx]
			note.Text = fieldText.Text()

			err := s.noteCtrl.UpdateNote(context.Background(), note)
			if err != nil {
				s.logger.Error(context.Background(), "update note error: "+err.Error())
				core.ErrorSnackbar(s.dialogBody, err, "Ошибка обновления заметки")
				return
			}

			fieldText.SetText(note.Text)
		})

		core.NewStretch(noteFrame)

		// Поле свитча для удаления
		isInSumSwitch := core.NewSwitch(noteFrame).
			SetType(core.SwitchCheckbox)

		isInSumSwitch.OnChange(func(e events.Event) {
			note := (*s.noteList)[idx]
			if isInSumSwitch.StateIs(states.Checked) {
				s.deletedIds = append(s.deletedIds, note.Id)
			} else {
				slices.DeleteFunc(s.deletedIds, func(s string) bool {
					return s == note.Id
				})
			}
		})
	})
}

func (s *NotesDialog) drawControlField(dialogMainFrame *core.Frame) {
	// Область кнопки "Добавить счет"
	upperFrame := core.NewFrame(dialogMainFrame)
	upperFrame.Styler(func(s *styles.Style) {
		s.CenterAll()
	})

	addButtonFrame := core.NewFrame(upperFrame)

	addButton := core.NewButton(addButtonFrame).
		SetType(core.ButtonElevated).
		SetIcon(icons.Add)

	addButton.Styler(func(s *styles.Style) {
		s.Background = colors.Scheme.Secondary.Container
		s.Color = colors.Scheme.Secondary.OnContainer
	})

	addButton.OnClick(func(e events.Event) {
		s.noteCtrl.AddNote(context.Background(), domain.Note{})
		s.listFrame.Update()
	})

	// Область с кнопками управления - "Назад" и "Сохранить и выйти"
	bottomFrame := core.NewFrame(dialogMainFrame)

	cancelButton := core.NewButton(bottomFrame).
		SetType(core.ButtonElevated).
		SetText("Назад")

	cancelButton.OnClick(func(e events.Event) {
		s.close()
	})

	core.NewStretch(bottomFrame)

	saveButton := core.NewButton(bottomFrame).
		SetType(core.ButtonElevated).
		SetText("Сохранить и выйти")

	saveButton.Styler(func(s *styles.Style) {
		s.Background = custom.ColorGreen
		s.Color = colors.Uniform(colors.White)
	})

	saveButton.OnClick(func(e events.Event) {
		err := s.noteCtrl.SaveAll(context.Background())
		if err != nil {
			s.logger.Error(context.Background(), "save accounts list: "+err.Error())
			core.ErrorSnackbar(s.dialogBody, err, "Ошибка сохранения")
			return
		}

		s.close()
	})

	core.NewStretch(bottomFrame)

	deleteButton := core.NewButton(bottomFrame).
		SetType(core.ButtonElevated).
		SetText("Удалить")

	deleteButton.Styler(func(s *styles.Style) {
		s.Color = custom.ColorDeleteRed
	})

	deleteButton.OnClick(func(e events.Event) {
		if len(s.deletedIds) == 0 {
			return
		}

		for _, noteId := range s.deletedIds {
			s.noteCtrl.DeleteNote(context.Background(), noteId)
		}

		core.MessageSnackbar(s.dialogBody, "Заметки удалены")
		s.listFrame.Update()
	})
}

func (s *NotesDialog) Run() {
	stage := s.dialogBody.NewDialog(s.appBody)

	stage.Pos.X = int(s.appBody.Geom.Size.Actual.Total.X/2) - 1050
	firstTableFrame := s.appBody.Child(1).AsTree().Child(0).AsTree().This.(*core.Frame)
	stage.Pos.Y = int(firstTableFrame.Geom.Size.Actual.Total.Y/2) - 200

	stage.Run()
}

func (s *NotesDialog) close() {
	s.dialogBody.Close()
}
