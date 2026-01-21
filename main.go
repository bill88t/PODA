package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Product struct
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Brand string  `json:"brand"`
	Price float64 `json:"price"`
}

var productsDB = []Product{
	{ID: 1, Name: "Urban Runner", Brand: "PODA", Price: 79.99},
	{ID: 2, Name: "Classic Loafer", Brand: "PODA", Price: 120.50},
	{ID: 3, Name: "Summer Sandal", Brand: "PODA", Price: 45.00},
}

func getAllProducts(c *fiber.Ctx) error {
	return c.JSON(productsDB)
}

func getProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	for _, p := range productsDB {
		if p.ID == id {
			return c.JSON(p)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
}

const userStorePath = "./users.json"

// User struct
type User struct {
	ID           int       `json:"id"`
	Fname        string    `json:"fname"`
	Lname        string    `json:"lname"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Birthday     time.Time `json:"birthday"`
}

// UserResponse struct
type UserResponse struct {
	ID       int       `json:"id"`
	Fname    string    `json:"fname"`
	Lname    string    `json:"lname"`
	Email    string    `json:"email"`
	Birthday time.Time `json:"birthday"`
}

var (
	usersDB    = make(map[int]User)
	usersMutex = &sync.Mutex{}
	nextUserID = 1
)

func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func saveUsers() error {
	usersMutex.Lock()
	var userSlice []User
	for _, u := range usersDB {
		userSlice = append(userSlice, u)
	}
	usersMutex.Unlock()

	data, err := json.MarshalIndent(userSlice, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(userStorePath, data, 0o644)
}

func loadUsers() error {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	file, err := os.Open(userStorePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("users.json not found, starting with empty DB")
			return nil
		}
		return err
	}
	defer func() { _ = file.Close() }()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var userSlice []User
	if err := json.Unmarshal(data, &userSlice); err != nil {
		return err
	}

	for _, u := range userSlice {
		usersDB[u.ID] = u
		if u.ID >= nextUserID {
			nextUserID = u.ID + 1
		}
	}
	log.Printf("Loaded %d users", len(usersDB))
	return nil
}

// SignUp user function
func SignUp(c *fiber.Ctx) error {
	req := new(struct {
		Fname    string `json:"fname"`
		Lname    string `json:"lname"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Birthday string `json:"birthday"`
	})
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if req.Email == "" || req.Password == "" || req.Fname == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "First name, email, and password are required"})
	}

	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid birthday format"})
	}

	usersMutex.Lock()
	for _, u := range usersDB {
		if u.Email == req.Email {
			usersMutex.Unlock()
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
		}
	}

	newUser := User{
		ID:           nextUserID,
		Fname:        req.Fname,
		Lname:        req.Lname,
		Email:        req.Email,
		PasswordHash: hashPassword(req.Password),
		Birthday:     birthday,
	}
	usersDB[newUser.ID] = newUser
	nextUserID++
	usersMutex.Unlock()

	if err := saveUsers(); err != nil {
		log.Printf("Error saving users: %v", err)
	}

	return c.Status(fiber.StatusCreated).JSON(UserResponse{
		ID:       newUser.ID,
		Fname:    newUser.Fname,
		Lname:    newUser.Lname,
		Email:    newUser.Email,
		Birthday: newUser.Birthday,
	})
}

// Login function
func Login(c *fiber.Ctx) error {
	req := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	hashedPassword := hashPassword(req.Password)

	usersMutex.Lock()
	defer usersMutex.Unlock()

	for _, u := range usersDB {
		if u.Email == req.Email && u.PasswordHash == hashedPassword {
			return c.JSON(UserResponse{
				ID:       u.ID,
				Fname:    u.Fname,
				Lname:    u.Lname,
				Email:    u.Email,
				Birthday: u.Birthday,
			})
		}
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
}

func main() {
	if err := loadUsers(); err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}

	app := fiber.New()

	app.Use(limiter.New(limiter.Config{
		Max:               15,
		Expiration:        15 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://0.0.0.0:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/products", getAllProducts)
	v1.Get("/products/:id", getProductByID)
	v1.Post("/users/signup", SignUp)
	v1.Post("/users/login", Login)

	app.Static("/", "./frontend/dist")
	app.Get("*", func(c *fiber.Ctx) error { return c.SendFile("./frontend/dist/index.html") })

	port := os.Getenv("PORT")
	if port == "" {
		port = "5173"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
