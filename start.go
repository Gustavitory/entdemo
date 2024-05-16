package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"entdemo/ent"
	"entdemo/ent/car"
	"entdemo/ent/user"

	_ "github.com/go-sql-driver/mysql"
)

func QueryCarUsers(ctx context.Context, user *ent.User) error {
	cars, err := user.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}

	for _, c := range cars {
		owner, err := c.QueryOwner().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying car %q owner: %w", c.Model, err)
		}
		log.Printf("car %q owner: %q\n", c.Model, owner.Name)
	}
	return nil
}

func QueryCars(ctx context.Context, user *ent.User) error {
	cars, err := user.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	log.Println("returned cars: %w", cars)

	ford, err := user.QueryCars().
		Where(car.Model("Ford")).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	log.Println(ford)
	return nil
}

func CreateCar(ctx context.Context, client *ent.Client) (*ent.User, error) {
	tesla, err := client.Car.Create().SetModel("Tesla").SetRegisteredAt(time.Now()).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car:  %w", err)
	}
	log.Println("car was created: ", tesla)

	ford, err := client.Car.Create().SetModel("Ford").SetRegisteredAt(time.Now()).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car:  %w", err)
	}
	log.Println("car was created: ", ford)
	newUser, err := client.User.Create().SetAge(30).SetName("Andres").AddCars(tesla, ford).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", newUser)
	return newUser, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.Query().Where(user.Name("Gustavo")).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.Create().SetAge(27).SetName("Gustavo").Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creting user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func main() {
	client, err := ent.Open("mysql", "root:12345@tcp(localhost:3306)/entgotest?parseTime=True")
	if err != nil {
		log.Fatal("Failed opening connection psql: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatal("failed creating schema resources: %v", err)
	}
}
