package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"main.go/models"
	"main.go/storage"
	"net/http"
	"os"
)

const port = ":8000"
const localhost ="127.0.0.1"

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupHandler(app *fiber.App) {
	api := app.Group("./api")
	api.Post("/create_books", r.CreateBook)
	api.Get("/get_books/:id", r.GetBookById)
	api.Get("/get_books", r.GetBooks)
	// api.Put("/update_book/:id", r.UpdateBooks)
	api.Delete("/delete_book/:id", r.DeleteBooks)
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}
	// mux.Headers().Set("Content-Types", "application/json")
	// json.NewEncoder(req).Encode(book)
	// Bind the request body to a struct

	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}
	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "unable to create book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book has been added"})
	return nil

}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := []models.Books{}

	err := r.DB.Find(&bookModels).Error
	if err != nil {
		context.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "resource not found"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "resource found",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) GetBookById(context *fiber.Ctx) error {
	bookModel := &models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book id fetched successfully",
		"data":    bookModel,
	})
	return nil
}

// func (r *Repository) UpdateBooks(context *fiber.Ctx)error{}

func (r *Repository) DeleteBooks(context *fiber.Ctx) error {
	bookModel := &models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil
	}
	err := r.DB.Delete(bookModel, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not delete book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book successfully deleted"})
	return nil
}

func main() {
	err := godotenv.Load(".env") // load configurations files
	handleErr(err, "unable to load configuration files")

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSL"),
	}

	db, err := storage.NewConnection(config)
	handleErr(err, "unable to load db")

	err = models.MigrateBooks(db)
	handleErr(err, "unable to migrate db")

	r := Repository{
		DB: db,
	}

	app := fiber.New()

	r.SetupHandler(app)

	app.Listen(port)
}

func handleErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s", msg)
	}
}
