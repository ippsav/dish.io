package postgres

import (
	"context"
	"fmt"
	"github.com/dish.io/internal/domain"
	"github.com/pkg/errors"
	"strings"
	"time"
)

// CreateRecipe : With goroutines / error = bad connection when having a lot of inserts
//func (s *Store) CreateRecipe(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error) {
//	now := time.Now()
//	tx, err := s.DB.BeginTx(ctx, nil)
//	if err != nil {
//		return nil, errors.Wrap(err, "could not create transaction")
//	}
//	defer tx.Rollback()
//
//	// Inserting the recipe first
//	st, err := tx.PrepareContext(ctx, "INSERT INTO recipes(recipe_name,description,owner_id) VALUES($1,$2,$3) RETURNING id,created_at")
//	if err != nil {
//		return nil, errors.Wrap(err, "could not prepare insert statement")
//	}
//	err = st.QueryRowContext(ctx, recipe.RecipeName, recipe.Description, recipe.OwnerID).Scan(&recipe.ID, &recipe.CreatedAt)
//	if err != nil {
//		return nil, errors.Wrap(err, "could not insert into recipes table")
//	}
//
//	//Inserting the ingredients of the recipe
//	err = createIngredient(ctx, tx, recipe)
//	if err != nil {
//		return nil, err
//	}
//	if err = tx.Commit(); err != nil {
//		return nil, errors.Wrap(err, "could not commit transaction")
//	}
//	fmt.Printf("%v\n", recipe)
//	fmt.Printf("it took %v",time.Since(now))
//	return recipe, nil
//}
//
//type Result struct {
//	Index int
//	Value string
//}
//
//type Source struct {
//	Index int
//	Value domain.Ingredient
//}
//
//func createIngredient(ctx context.Context, tx *sql.Tx, recipe *domain.Recipe) error {
//	st, err := tx.PrepareContext(ctx, "INSERT INTO ingredients(name,qty,measurement,recipe_id) VALUES($1,$2,$3,$4) RETURNING id")
//	if err != nil {
//		return errors.Wrap(err, "could not prepare insert statement")
//	}
//	errorChan := make(chan error)
//	srcChan := make(chan Source, len(recipe.Ingredients))
//	resChan := make(chan Result, len(recipe.Ingredients))
//
//	// Creating Workers
//	numWorkers := 4
//
//	for i := 0; i < numWorkers; i++ {
//		go func(source <-chan Source, result chan Result, errorChan chan error) {
//			for{
//				fmt.Println("entered")
//				src, ok := <-source
//				if !ok {
//					return
//				}
//				res := Result{
//					Index: src.Index,
//				}
//				ingredient := src.Value
//				err := st.QueryRowContext(ctx, ingredient.Name, ingredient.Quantity, ingredient.Measurement, recipe.ID).Scan(&res.Value)
//				if err != nil {
//					errorChan <- errors.Wrap(err, "could not insert into ingredients table")
//					return
//				}
//				fmt.Println(res)
//				result <- res
//				}
//		}(srcChan, resChan, errorChan)
//	}
//	for i, ing := range recipe.Ingredients {
//		srcChan <- Source{
//			Index: i,
//			Value: ing,
//		}
//	}
//	go func() {
//		for i := 0; i < len(recipe.Ingredients); i++ {
//			res := <-resChan
//			recipe.Ingredients[res.Index].ID = res.Value
//		}
//		close(errorChan)
//		close(resChan)
//		close(srcChan)
//	}()
//	err = <-errorChan
//	if err != nil {
//		return err
//	}
//	return nil
//}

// CreateRecipe : concat sql values into one query

func (s *Store) CreateRecipe(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error) {
	now := time.Now()
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create transaction")
	}
	defer tx.Rollback()

	// Inserting recipe
	st, err := tx.PrepareContext(ctx, "INSERT INTO recipes(recipe_name,description,owner_id) VALUES($1,$2,$3) RETURNING id,created_at")
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare insert statement")
	}
	err = st.QueryRowContext(ctx, recipe.RecipeName, recipe.Description, recipe.OwnerID).Scan(&recipe.ID, &recipe.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "could not insert into recipes table")
	}
	// Inserting the ingredients
	errChan := make(chan error)
	idChan := make(chan string)
	go func() {
		query := "INSERT INTO ingredients(name,qty,measurement,recipe_id) VALUES "
		vals := make([]interface{}, 0)
		index := 1
		for _, ing := range recipe.Ingredients {
			query += fmt.Sprintf("($%d,$%d,$%d,$%d),", index, index+1, index+2, index+3)
			vals = append(vals, ing.Name, ing.Quantity, ing.Measurement, recipe.ID)
			index += 4
		}
		query = strings.TrimSuffix(query, ",")
		query = fmt.Sprintf("%s RETURNING id", query)
		fmt.Println(query)
		st, err := tx.PrepareContext(ctx, query)
		if err != nil {
			errChan <- errors.Wrap(err, "could not prepare insert statement")
			return
		}
		rows, err := st.QueryContext(ctx, vals...)
		if err != nil {
			errChan <- errors.Wrap(err, "could not insert into ingredients table")
			return
		}
		for rows.Next() {
			var id string
			rows.Scan(&id)
			idChan <- id
		}
	}()
	go func() {
		for i, _ := range recipe.Ingredients {
			recipe.Ingredients[i].ID = <-idChan
		}
		close(errChan)
	}()
	err = <-errChan
	close(idChan)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	fmt.Printf("recipe first ing ID = %s,name =%s\n", recipe.Ingredients[0].ID, recipe.Ingredients[0].Name)
	fmt.Printf("it took %v\n", time.Since(now))
	return recipe, nil
}

func (s *Store) CreateRecipeSync(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error) {
	now := time.Now()
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
	st, err = tx.PrepareContext(ctx, "INSERT INTO ingredients(name,qty,measurement,recipe_id) VALUES($1,$2,$3,$4) RETURNING id")
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare insert statement")
	}
	for i, ing := range recipe.Ingredients {
		err := st.QueryRowContext(ctx, ing.Name, ing.Quantity, ing.Measurement, recipe.ID).Scan(&recipe.Ingredients[i].ID)
		if err != nil {
			return nil, errors.Wrap(err, "could not insert into ingredients table")
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "could not commit transaction")
	}
	fmt.Printf("it took %v", time.Since(now))
	return recipe, nil
}
