package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yousefzinsazk78/note_app_api/types"
)

func (a *Api) HandleCreateNote(c *fiber.Ctx) error {
	var note types.Note
	if err := c.BodyParser(&note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	res, err := a.NoteStorer.InsertNote(&note)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (a *Api) HandleUpdateNote(c *fiber.Ctx) error {
	var note *types.Note
	id := c.Params("id")
	resID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(err)
	}
	if err := c.BodyParser(&note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	log.Println(note.Title, note.Description, resID)
	res, err := a.NoteStorer.UpdateNote(note, resID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (a *Api) HandleNotes(c *fiber.Ctx) error {
	res, err := a.NoteStorer.GetNotes()
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(res)
}

func (a *Api) HandleNoteID(c *fiber.Ctx) error {
	id := c.Params("id")
	resID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(err)
	}
	res, err := a.NoteStorer.GetNoteByID(resID)
	if err != nil {
		return ErrNotFound()
	}
	if res.ID == 0 && res.Title == "" && res.Description == "" {
		return NewError(http.StatusNotFound, "not found!")
	}
	return c.Status(http.StatusOK).JSON(res)
}

func (a *Api) HandleDeleteNote(c *fiber.Ctx) error {
	id := c.Params("id")
	resID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(err)
	}
	res, err := a.NoteStorer.DeleteNote(resID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	if res == 0 {
		return NewError(http.StatusBadRequest, "bad request")
	}
	return c.Status(fiber.StatusOK).JSON(res)
}
