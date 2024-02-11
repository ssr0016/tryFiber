package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gofiber/fiber/v2"
)

var users []User
var secretKey = []byte("your-secret-key")

// User struct represents user data
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

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

	// Initialize sample users
	users = append(users, User{ID: 1, Username: "user1", Password: "password1"})
	users = append(users, User{ID: 2, Username: "user2", Password: "password2"})

	// Initialize sample data
	books = append(books, Book{ID: getNextID(), Title: "Book 1", Category: "Fiction", Author: "Author 1"})
	books = append(books, Book{ID: getNextID(), Title: "Book 2", Category: "Non-fiction", Author: "Author 2"})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello Mommy!")
	})

	// Define login routes
	app.Post("/login", login)
	app.Get("/protected", authenticate, protectedRoute)

	// Define book routes
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

// Login handler authenticates users
func login(c *fiber.Ctx) error {
	var reqUser User
	if err := c.BodyParser(&reqUser); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Find the user in the list
	for _, user := range users {
		if user.Username == reqUser.Username && user.Password == reqUser.Password {
			// Create token
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["username"] = user.Username
			claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

			// Sign the token with a secret key
			tokenString, err := token.SignedString([]byte(secretKey))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate token")
			}

			return c.JSON(fiber.Map{"token": tokenString})
		}
	}

	// User not found, return error message
	return c.Status(fiber.StatusUnauthorized).SendString("Invalid username or password")
}

// Middleware to authenticate requests using JWT token
func authenticate(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing authorization token")
	}

	// Extract JWT token from the "Authorization" header
	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid authorization header format")
	}
	token := parts[1]

	// Parse and validate token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Check token signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}

	// Check if token is valid
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token: username claim missing")
		}
		c.Locals("username", username)
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
}

// Protected route that requires authentication
func protectedRoute(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	return c.SendString("Welcome, " + username)
}
