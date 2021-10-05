package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dish.io/internal/domain"
	"github.com/pkg/errors"
)

func (s *Store) CreateRecipe(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create transaction")
	}
	defer tx.Rollback()
	st, err := tx.PrepareContext(ctx, "INSERT INTO recipes(recipe_name,description,owner_id) VALUES($1,$2,$3) RETURNING id,created_at")
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare insert statement")
	}
	err = st.QueryRowContext(ctx, recipe.RecipeName, recipe.Description, recipe.OwnerID).Scan(&recipe.ID, &recipe.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "could not insert into recipes table")
	}
	err = createIngredient(ctx, tx, recipe)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "could not commit transaction")
	}
	fmt.Printf("%v\n", recipe)
	return recipe, nil
}

type Result struct {
	Index int
	Value string
}

type Source struct {
	Index int
	Value domain.Ingredient
}

func createIngredient(ctx context.Context, tx *sql.Tx, recipe *domain.Recipe) error {
	st, err := tx.PrepareContext(ctx, "INSERT INTO ingredients(name,qty,measurement,recipe_id) VALUES($1,$2,$3,$4) RETURNING id")
	if err != nil {
		return errors.Wrap(err, "could not prepare insert statement")
	}
	errorChan := make(chan error)
	srcChan := make(chan Source, len(recipe.Ingredients))
	resChan := make(chan Result, len(recipe.Ingredients))

	// Creating Workers
	numWorkers := 4

	for i := 0; i < numWorkers; i++ {
		go func(source <-chan Source, result chan Result, errorChan chan error) {
			src, ok := <-source
			if !ok {
				return
			}
			res := Result{
				Index: src.Index,
			}
			ingredient := src.Value
			err := st.QueryRowContext(ctx, ingredient.Name, ingredient.Quantity, ingredient.Measurement, recipe.ID).Scan(&res.Value)
			if err != nil {
				errorChan <- errors.Wrap(err, "could not insert into ingredients table")
				return
			}
			fmt.Println(res)
			result <- res
		}(srcChan, resChan, errorChan)
	}
	for i, ing := range recipe.Ingredients {
		srcChan <- Source{
			Index: i,
			Value: ing,
		}
	}
	go func() {
		for i := 0; i < len(recipe.Ingredients); i++ {
			res := <-resChan
			recipe.Ingredients[res.Index].ID = res.Value
		}
		close(errorChan)
		close(resChan)
		close(srcChan)
	}()
	err = <-errorChan
	if err != nil {
		return err
	}
	return nil
}
