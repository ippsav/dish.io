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
					},
					OwnerID: "01fbabfe-30e8-4162-aa7f-52dffc0c3e6b",
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
