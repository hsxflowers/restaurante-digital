package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"

	process "github.com/hsxflowers/restaurante-digital/processing"
	restauranteDatabase "github.com/hsxflowers/restaurante-digital/processing/db"
)

func main() {
	var wg sync.WaitGroup
	var ctx = context.Background()

	var restauranteDb restauranteDatabase.RestauranteDatabase
	var err error

	db, err := sql.Open("postgres", "postgresql://admin:Ep7elyAZpA5dYxZMe3vOjIIJJ1dXF3XZ@dpg-cr7fi923esus7388q3c0-a.oregon-postgres.render.com/dbrestaurante")
	if err != nil {
		log.Fatal("Erro ao conectar com DB: ", err)
	}

	restauranteDb = restauranteDatabase.NewSQLStore(db)
	processInstance := process.NewProcess(&wg, restauranteDb, ctx)

	processInstance.StartWorkers()
	processInstance.DispatchPedidos(ctx)

	usuario := "1234"
	pedido := "abc"

	// aqui você pode cancelar um pedido
	process.CancelarPedido(ctx, pedido, restauranteDb)

	time.Sleep(15 * time.Second)
	// aqui você consegue pedir a conta do usuário desejado
	process.PedirConta(ctx, restauranteDb, usuario)

	// aqui você consegue deletar todos os pedidos do restaurante para testar melhor o código
	process.DeletarPedidos(ctx, restauranteDb)

	wg.Wait()
	fmt.Println("Todos os pedidos foram processados ou cancelados.")
}
