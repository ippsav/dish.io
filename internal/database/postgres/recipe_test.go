package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dish.io/internal/domain"
	_ "github.com/lib/pq"
	"testing"
)

func TestStore_CreateRecipe(t *testing.T) {
	db, err := sql.Open("postgres", "user=user dbname=dish-db password=longpassword port=7001 sslmode=disable")
	if err != nil {
		fmt.Printf("%v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx    context.Context
		recipe *domain.Recipe
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Recipe
		wantErr bool
	}{
		{
			name: "TEST 1",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx: ctx,
				recipe: &domain.Recipe{
					ID:          "",
					RecipeName:  "TEST",
					Description: "TEST RECIPE",
					Ingredients: []domain.Ingredient{
						{
							ID:          "",
							Name:        "Tomato",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheese",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "ADaw",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "dwadawd",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "dwadwqdqwsad",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "dwqdqsadwqa",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "random",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheeses",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatos",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheesess",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatoq",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheesef",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatoc",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheeser",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatop",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheesep",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatok",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheesen",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomaton",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheeseb",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatog",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheesex",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatox",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheesez",
							Quantity:    250,
							Measurement: "g",
						},
					},
					OwnerID: "4bf38f69-04c8-4791-9a64-4838446dd913",
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				DB: tt.fields.DB,
			}
			_, err := s.CreateRecipe(tt.args.ctx, tt.args.recipe)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRecipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStore_CreateRecipeSync(t *testing.T) {
	db, err := sql.Open("postgres", "user=user dbname=dish-db password=longpassword port=7001 sslmode=disable")
	if err != nil {
		fmt.Printf("%v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx    context.Context
		recipe *domain.Recipe
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Recipe
		wantErr bool
	}{
		{
			name: "TEST 1",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx: ctx,
				recipe: &domain.Recipe{
					ID:          "",
					RecipeName:  "TEST",
					Description: "TEST RECIPE",
					Ingredients: []domain.Ingredient{
						{
							ID:          "",
							Name:        "Tomato",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheese",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "ADaw",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "dwadawd",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "dwadwqdqwsad",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "dwqdqsadwqa",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "random",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheeses",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatos",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheesess",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatoq",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheesef",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatoc",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheeser",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatop",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheesep",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatok",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheesen",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomaton",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheeseb",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatog",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheesex",
							Quantity:    250,
							Measurement: "g",
						},
						{
							ID:          "",
							Name:        "Tomatox",
							Quantity:    2,
							Measurement: "unit",
						},
						{
							ID:          "",
							Name:        "Cheesez",
							Quantity:    250,
							Measurement: "g",
						},
					},
					OwnerID: "b0c87efc-5f65-444b-9cbb-2eae87911ff8",
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				DB: tt.fields.DB,
			}
			_, err := s.CreateRecipeSync(tt.args.ctx, tt.args.recipe)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRecipeSync() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
