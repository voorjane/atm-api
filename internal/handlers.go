package internal

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"nedorez-test/pkg"
	"strconv"
	"sync"
	"time"
)

var globerr error

type Request struct {
	Balance float64 `json:"balance"`
}

func Route(app *fiber.App, db *gorm.DB) {
	acc := app.Group("/accounts")
	acc.Post("/", NewAccount(db))
	id := acc.Group(":id")
	id.Post("/deposit", DepositHandler(db))
	id.Post("/withdraw", WithdrawHandler(db))
	id.Get("/balance", BalanceHandler(db))
}

func NewAccount(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req Request
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to parse JSON",
			})
		}
		if req.Balance < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "balance must be positive",
			})
		}
		db.Create(&pkg.Account{
			Balance: req.Balance,
		})
		var id int
		db.Table("accounts").Select("id").Order("id DESC").Limit(1).Scan(&id) // "SELECT id FROM accounts ORDER BY id DESC LIMIT 1"
		log.Printf("account created with id %d at: %v\n", id, time.Now())
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"response": "account created",
			"time":     time.Now(),
			"balance":  req.Balance,
			"id":       id,
		})
	}
}

func BalanceHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		param := c.Params("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "account not found",
			})
		}
		account := &pkg.Account{
			Id: -1,
		}
		db.Table("accounts").Where("id = ?", id).Scan(&account) // "SELECT id FROM accounts WHERE id = ?"
		if account.Id == -1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "account not found",
			})
		}
		balance := account.GetBalance()
		log.Printf("Requested balance for account id %d at %v", id, time.Now())
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"balance": balance,
			"id":      id,
		})
	}
}

func DepositHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req Request
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to parse JSON",
			})
		}

		param := c.Params("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "account not found",
			})
		}
		account := &pkg.Account{
			Id: -1,
		}
		db.Table("accounts").Where("id = ?", id).Scan(&account)
		if account.Id == -1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "account not found",
			})
		}

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			globerr = account.Deposit(req.Balance)
		}()
		wg.Wait()
		if globerr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": globerr.Error(),
			})
		}
		db.Table("accounts").Where("id = ?", id).Update("balance", account.Balance)
		log.Printf("Balance added for account id %d at %v", id, time.Now())
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"response": "success",
			"balance":  account.Balance,
			"id":       id,
		})
	}
}

func WithdrawHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req Request
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to parse JSON",
			})
		}

		param := c.Params("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "account not found",
			})
		}
		account := &pkg.Account{
			Id: -1,
		}
		db.Table("accounts").Where("id = ?", id).Scan(&account)
		if account.Id == -1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "account not found",
			})
		}

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			globerr = account.Withdraw(req.Balance)
		}()
		wg.Wait()
		if globerr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": globerr.Error(),
			})
		}
		db.Table("accounts").Where("id = ?", id).Update("balance", account.Balance)
		log.Printf("Withdraw made for account id %d at %v", id, time.Now())
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"response": "success",
			"balance":  account.Balance,
			"id":       id,
		})
	}
}
