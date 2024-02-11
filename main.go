package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Author   string `json:"author"`
}

var books []Book
var idCounter = 0

func main() {
	app := fiber.New()

	// Initialize sample data
	books = append(books, Book{ID: getNextID(), Title: "Book 1", Category: "Fiction", Author: "Author 1"})
	books = append(books, Book{ID: getNextID(), Title: "Book 2", Category: "Non-fiction", Author: "Author 2"})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello Mommy!")
	})

	// Define routes
	app.Get("/books", getBooks)
	app.Get("/books/:id", getBookByID)
	app.Post("/books", createBook)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)

	// Start server
	app.Listen(":8080")
}

func getNextID() int {
	idCounter++
	return idCounter
}

// Handler functions

func getBooks(c *fiber.Ctx) error {
	return c.JSON(books)
}

func getBookByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	for _, book := range books {
		if book.ID == id {
			return c.JSON(book)
		}
	}

	return c.Status(fiber.StatusNotFound).SendString("Book not found")
}

func createBook(c *fiber.Ctx) error {
	var newBook Book
	if err := c.BodyParser(&newBook); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	newBook.ID = getNextID()
	books = append(books, newBook)
	return c.JSON(newBook)
}

func updateBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	var updateBook Book
	if err := c.BodyParser(&updateBook); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	for i, book := range books {
		if book.ID == id {
			updateBook.ID = id
			books[i] = updateBook
			return c.JSON(updateBook)
		}
	}

	return c.Status(fiber.StatusNotFound).SendString("Book not found")
}

func deleteBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	for i, book := range books {
		if book.ID == id {
			// Remove the book from the slice
			books = append(books[:i], books[i+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}

	return c.Status(fiber.StatusNotFound).SendString("Book not found")
}
