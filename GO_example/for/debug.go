package main

import (
	"errors"
	"fmt"
	"strings"
)

type User struct {
	ID      int
	Name    string
	Email   string
	Age     int
	Balance int
	Active  bool
}

type Order struct {
	ID     int
	UserID int
	Amount int
}

func main() {
	users := generateUsers(10_000)

	orders := []Order{
		{ID: 101, UserID: 100, Amount: 500},
		{ID: 102, UserID: 7359, Amount: 1200},
		{ID: 103, UserID: 9000, Amount: 800},
	}

	userMap := make(map[int]User, len(users))

	totalBalance := 0

	for i, user := range users {
		// DEBUG 1 — WITH INDEX
		// Breakpoint condition: i == 5000
		currentIndex := i

		// DEBUG 2 — WITH ID
		// Breakpoint condition: user.ID == 7359
		currentUserID := user.ID

		// DEBUG 3 — WITH INVALID ITEM
		invalid := !isValidUser(user)

		// DEBUG 4 — WITH INDEX RANGE
		// Breakpoint condition: i >= 4995 && i <= 5005
		inDebugRange := i >= 4995 && i <= 5005

		// DEBUG 5 — CHECKPOINT EVERY 1000 ITEM
		// Breakpoint condition: isCheckpoint
		isCheckpoint := i%1000 == 0

		// DEBUG 6 — WHEN FUNCTION RETURN ERROR
		err := processUser(user)

		// DEBUG 8 — BEFORE AND AFTER THE ITEM BEEN UPDATE
		oldBalance := user.Balance
		newBalance := user.Balance

		if user.Active {
			newBalance += 100
		}

		// DEBUG 10 — MARKER VARIABLE
		hasInvalidEmail := !strings.Contains(user.Email, "@")
		hasNegativeBalance := user.Balance < 0
		isUnderAge := user.Age < 18

		isTarget := hasInvalidEmail ||
			hasNegativeBalance ||
			isUnderAge

		debugPoint := true

		if invalid {
			fmt.Printf(
				"INVALID index=%d id=%d user=%+v\n",
				i,
				user.ID,
				user,
			)
		}

		if err != nil {
			fmt.Printf(
				"ERROR index=%d id=%d err=%v\n",
				i,
				user.ID,
				err,
			)
		}

		if isCheckpoint {
			fmt.Printf(
				"CHECKPOINT index=%d id=%d total=%d\n",
				i,
				user.ID,
				totalBalance,
			)
		}

		if user.Active {
			totalBalance += newBalance
		}

		userMap[user.ID] = user

		// Giữ các biến tồn tại để dễ Watch trong debugger.
		_ = currentIndex
		_ = currentUserID
		_ = inDebugRange
		_ = isTarget
		_ = debugPoint
		_ = oldBalance
		_ = newBalance
	}

	// DEBUG 7 — 2 FOR LOOP
	for userIndex, user := range users {
		for orderIndex, order := range orders {
			if user.ID != order.UserID {
				continue
			}

			orderTotal := order.Amount + calculateFee(user)

			// Breakpoint condition:
			// user.ID == 7359 && order.ID == 102
			fmt.Printf(
				"ORDER userIndex=%d orderIndex=%d userID=%d orderID=%d total=%d\n",
				userIndex,
				orderIndex,
				user.ID,
				order.ID,
				orderTotal,
			)
		}
	}

	// DEBUG 9 — MAP LOOP
	for id, user := range userMap {
		mapValue := user.Balance * 2

		// Breakpoint condition: id == 7359
		if id == -1 {
			fmt.Println(mapValue)
		}
	}

	fmt.Println("Done. Total balance:", totalBalance)
}

func generateUsers(total int) []User {
	users := make([]User, total)

	for i := range users {
		id := i + 1

		users[i] = User{
			ID:      id,
			Name:    fmt.Sprintf("User %d", id),
			Email:   fmt.Sprintf("user%d@example.com", id),
			Age:     18 + id%50,
			Balance: id * 10,
			Active:  id%2 == 0,
		}
	}

	// Proactively generate error data.
	users[2499].Name = ""
	users[4999].Email = "invalid-email"
	users[7358].Balance = -500
	users[8999].Age = 15

	return users
}

func isValidUser(user User) bool {
	return user.ID > 0 &&
		user.Name != "" &&
		strings.Contains(user.Email, "@") &&
		user.Age >= 18 &&
		user.Balance >= 0
}

func processUser(user User) error {
	switch {
	case user.Name == "":
		return errors.New("name is empty")
	case !strings.Contains(user.Email, "@"):
		return errors.New("email is invalid")
	case user.Age < 18:
		return errors.New("user is under age")
	case user.Balance < 0:
		return errors.New("balance is negative")
	default:
		return nil
	}
}

func calculateFee(user User) int {
	if user.Active {
		return 10
	}

	return 20
}
