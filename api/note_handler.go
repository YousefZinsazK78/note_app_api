package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (a *Api) HandleIndex(c *fiber.Ctx) error {
	notes, err := a.NoteStorer.GetNotes()
	if err != nil {
		return ErrNotFound()
	}
	return c.Render("index", fiber.Map{
		"notesVar": notes,
	})
}

func (a *Api) HandleCreate(c *fiber.Ctx) error {
	return c.Render("create", fiber.Map{})
}

func (a *Api) HandleCreatePost(c *fiber.Ctx) error {
	var note types.Note

	if err := c.BodyParser(&note); err != nil {
		return NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err := a.NoteStorer.InsertNote(&note)
	if err != nil {
		return NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Redirect("/")
}

func (a *Api) HandleDeletePost(c *fiber.Ctx) error {
	var note types.Note

	if err := c.BodyParser(&note); err != nil {
		return NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err := a.NoteStorer.DeleteNote(note.ID)
	if err != nil {
		return NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Redirect("/")
}

func (a *Api) HandleEdit(c *fiber.Ctx) error {
	var note types.Note

	if err := c.QueryParser(&note); err != nil {
		return NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Render("edit", fiber.Map{
		"ID":          note.ID,
		"Title":       note.Title,
		"Description": note.Description,
	})
}

func (a *Api) HandleEditPost(c *fiber.Ctx) error {
	var note types.Note

	if err := c.BodyParser(&note); err != nil {
		return NewError(fiber.StatusBadRequest, err.Error())
	}

	log.Println(note.ID)
	_, err := a.NoteStorer.UpdateNote(&note, note.ID)
	if err != nil {
		return NewError(fiber.StatusAccepted, err.Error())
	}

	return c.Redirect("/")
}

func (a *Api) HandleWelcome(c *fiber.Ctx) error {
	session_token := c.Cookies("session_token")
	if session_token == "" {
		return ErrUnAuthorized()
	}

	session, err := a.SessionStorer.GetSession(session_token)
	log.Println("get session from db : ", session)
	if err != nil {
		return NewError(fiber.StatusUnauthorized, err.Error())
		// return ErrUnAuthorized()
	}
	if session.IsExpired() {
		log.Println("session IsExpired : ", session.IsExpired())
		err := a.SessionStorer.DeleteSession(session_token)
		if err != nil {
			return ErrBadRequest()
		}
		return ErrUnAuthorized()
	}
	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("welcome %s", session.Username))
}

func (a *Api) HandleSignIn(c *fiber.Ctx) error {
	//get information from request body
	var user types.User
	if err := c.BodyParser(&user); err != nil {
		return ErrBadRequest()
	}
	//check username and password
	dbUser, err := a.UserStorer.GetUserByUsername(user.Username)
	if err != nil {
		return ErrBadRequest()
	}

	if dbUser.Username != user.Username && dbUser.Password != user.Password {
		return ErrUnAuthorized()
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	err = a.SessionStorer.InsertSession(user.Username, expiresAt, sessionToken)
	if err != nil {
		return NewError(fiber.StatusInternalServerError, "session token not inserted")
	}

	user_cookie := fiber.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	}
	c.Cookie(&user_cookie)

	return c.Status(fiber.StatusAccepted).SendString("user login successfully!")
}

func (a *Api) HandleSignUp(c *fiber.Ctx) error {
	//get information from request body
	var user types.User
	if err := c.BodyParser(&user); err != nil {
		return ErrBadRequest()
	}
	//store user in database
	if err := a.UserStorer.InsertUser(user); err != nil {
		return NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusAccepted).SendString("user successfully inserted!")
}
