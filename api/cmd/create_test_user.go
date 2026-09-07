package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/repository"
	"ivpn.net/email/api/internal/utils"
)

// testUserSubscriptionTier must never contain model.Tier1, which Subscription.Active() treats as always-inactive.
const testUserSubscriptionTier = "IVPN Tier 2"
const testUserSubscriptionDays = 365

func runCreateTestUser(args []string) error {
	fs := flag.NewFlagSet("create-test-user", flag.ContinueOnError)
	email := fs.String("email", "", "Email address for the test user")
	pass := fs.String("pass", "", "Password for the test user")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *email == "" {
		fs.Usage()
		return fmt.Errorf("--email is required")
	}
	if *pass == "" {
		fs.Usage()
		return fmt.Errorf("--pass is required")
	}

	if os.Getenv("ALLOW_TEST_USER_CREATION") != "true" {
		return fmt.Errorf("create-test-user is disabled; set ALLOW_TEST_USER_CREATION=true to enable")
	}

	if err := utils.ValidateEmail(*email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := utils.ValidatePassword(*pass); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	db, err := repository.NewDB(dbConfigFromEnv())
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	ctx := context.Background()

	exists, err := db.CheckDuplicateRecipient(ctx, *email)
	if err != nil {
		return fmt.Errorf("checking duplicate recipient: %w", err)
	}
	if exists {
		return fmt.Errorf("a recipient with email %s already exists", *email)
	}

	user := model.User{
		Email:    *email,
		IsActive: true,
	}
	if err := user.SetPassword(*pass); err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	user, err = db.PostUser(ctx, user)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return fmt.Errorf("a user with email %s already exists", *email)
		}
		return fmt.Errorf("creating user: %w", err)
	}

	if err := db.PostSettings(ctx, model.Settings{UserID: user.ID}); err != nil {
		return fmt.Errorf("creating settings: %w", err)
	}

	sub := model.Subscription{
		UserID:      user.ID,
		ActiveUntil: time.Now().AddDate(0, 0, testUserSubscriptionDays),
		Tier:        testUserSubscriptionTier,
	}
	sub.ID = uuid.New().String()
	if err := db.PostSubscription(ctx, sub); err != nil {
		return fmt.Errorf("creating subscription: %w", err)
	}

	if _, err := db.PostRecipient(ctx, model.Recipient{
		UserID:   user.ID,
		Email:    *email,
		IsActive: true,
	}); err != nil {
		return fmt.Errorf("creating recipient: %w", err)
	}

	log.Printf("test user created: id=%s email=%s", user.ID, *email)
	return nil
}
