package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/some-things/ent-eval/ent"
	"github.com/some-things/ent-eval/ent/car"
	"github.com/some-things/ent-eval/ent/user"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// TODO: switch to mysql and read from env file here
	// client, err := ent.Open("mysql", "<user>:<pass>@tcp(<host>:<port>)/<database>?parseTime=True")
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()
	ctx := context.Background()

	// Run the auto migraiton tool
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Create a user
	if _, err = CreateUser(ctx, client); err != nil {
		log.Fatalf("failed creating user: %v", err)
	}

	// Query the user
	if _, err = QueryUser(ctx, client); err != nil {
		log.Fatalf("failed querying user: %v", err)
	}

	// Create cars
	a8m, err := CreateCars(ctx, client)
	if err != nil {
		log.Fatalf("failed creating cars: %v", err)
	}

	// Query cars
	if err := QueryCars(ctx, a8m); err != nil {
		log.Fatalf("failed querying cars: %v", err)
	}

	// Query car users
	if err := QueryCarUsers(ctx, a8m); err != nil {
		log.Fatalf("failed querying car users: %v", err)
	}
}

// Users
func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.Create().SetAge(30).SetName("a8m").Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.Name("a8m")).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}

	log.Println("user returned: ", u)
	return u, nil
}

// Cars
func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	// Create a new car with model "Tesla"
	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", tesla)

	// Create a new car with the model "Ford"
	ford, err := client.Car.Create().SetModel("Ford").SetRegisteredAt(time.Now()).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", ford)

	// Create a new user and add it to the 2 cars
	a8m, err := client.User.Create().SetAge(30).SetName("a8m").AddCars(tesla, ford).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", a8m)

	return a8m, nil
}

func QueryCars(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying cars: %w", err)
	}
	log.Println("cars returned: ", cars)

	// Filtering specific cars
	ford, err := a8m.QueryCars().Where(car.Model("Ford")).Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying cars: %w", err)
	}
	log.Println("a8m user's fords returned: ", ford)
	return nil
}

func QueryCarUsers(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying cars: %w", err)
	}
	// Query the inverse edge
	for _, c := range cars {
		owner, err := c.QueryOwner().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying car %q owner: %w", c.Model, err)
		}
		log.Printf("car %q owner: %q\n", c.Model, owner.Name)
	}
	return nil
}
