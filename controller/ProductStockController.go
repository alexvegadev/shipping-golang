package controller

import (
	"log"
	"os"
	"shipping/core"
	"shipping/repository"
	"strconv"
)

/**
** Constants
**/
const INSERT_PRODUCT_STOCK = "INSERT INTO product_stock (product_id, deposit, location, quantity) VALUES ($1, $2, $3, $4)"
const UPDATE_PRODUCT_STOCK = "UPDATE product_stock SET quantity = $1 WHERE product_id = $2 and deposit = $3 and location = $4"
const DELETE_PRODUCT_STOCK = "DELETE FROM product_stock WHERE product_id = $1 and deposit = $2 and location = $3"
const SET_PRODUCT_STOCK = "SELECT * FROM product_stock WHERE product_id = $1 and deposit = $2 and location = $3"
const SET_PRODUCT_STOCK_BY_DEPOSIT_AND_LOCATION = "SELECT * FROM product_stock WHERE deposit = $1 and location = $2"
const SET_PRODUCT_STOCK_BY_DEPOSIT_AND_PRODUCT = "SELECT * FROM product_stock WHERE deposit = $1 and product_id = $2"
const GET_QUANTITY_BY_LOCATION = "SELECT SUM(quantity) FROM product_stock WHERE location = $1 GROUP BY location"

type ProductStockController struct {
}

func (p *ProductStockController) IsLimitReachedDeposit(newQuantity int, Location string, repository repository.PostgreRepository) bool {
	var client = repository.CreateClientDatabase()
	defer client.Close()
	row := client.QueryRow(GET_QUANTITY_BY_LOCATION, Location)
	var currentDepositCount int
	err := row.Scan(&currentDepositCount)
	limit, err := strconv.Atoi(os.Getenv("DEPOSIT_LIMIT"))
	if err != nil {
		log.Println("No rows were returned!", err)
		return false
	}
	return currentDepositCount+newQuantity >= limit
}

func (p *ProductStockController) IsLimitProductCountReached(ProductId string, Location string, Deposit string, repository repository.PostgreRepository) bool {
	products := p.GetDistinctByProduct(Location, Deposit, repository)
	return len(products) >= 3 && !contains(products, ProductId)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func (p *ProductStockController) GetDistinctByProduct(Location string, Deposit string, repository repository.PostgreRepository) []string {
	var client = repository.CreateClientDatabase()
	defer client.Close()
	rows, err := client.Query("SELECT DISTINCT product_id FROM product_stock where location = $1 AND deposit = $2", Location, Deposit)
	if err != nil {
		log.Println("No rows were returned!", err)
		return []string{}
	}
	var products []string
	for rows.Next() {
		var productId string
		err := rows.Scan(&productId)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		// append the productId
		products = append(products, productId)
	}
	return products
}

func (p *ProductStockController) Save(ps core.ProductStock, repository repository.PostgreRepository) core.ProductStock {
	var client = repository.CreateClientDatabase()
	defer client.Close()
	product, _ := p.GetProduct(ps.ProductId, ps.Deposit, ps.Location, repository)
	if product.ProductId != "" {
		ps.Quantity = ps.Quantity + product.Quantity
		p.Update(ps, repository)
	} else {
		p.Create(ps, repository)
	}
	log.Println("Inserted a new product")
	return ps
}

func (p *ProductStockController) Create(ps core.ProductStock, repository repository.PostgreRepository) core.ProductStock {
	var client = repository.CreateClientDatabase()
	defer client.Close()
	_, err := client.Exec(INSERT_PRODUCT_STOCK, ps.ProductId, ps.Deposit, ps.Location, ps.Quantity)
	if err != nil {
		return core.ProductStock{}
	}
	log.Println("Inserted a new product")
	return ps
}

func (p *ProductStockController) GetProductsByDepositAndLocation(Deposit string, Location string, repository repository.PostgreRepository) []core.ProductStock {
	var client = repository.CreateClientDatabase()
	defer client.Close()
	rows, err := client.Query(SET_PRODUCT_STOCK_BY_DEPOSIT_AND_LOCATION, Deposit, Location)
	if err != nil {
		log.Println("No rows were returned!", err)
		return []core.ProductStock{}
	}
	defer rows.Close()
	var products []core.ProductStock
	for rows.Next() {
		var product core.ProductStock

		// unmarshal
		err := rows.Scan(&product.Deposit, &product.Location, &product.ProductId, &product.Quantity)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		// append the product
		products = append(products, product)
	}
	return products
}

func (p *ProductStockController) GetProductsByDepositAndProduct(Deposit string, Location string, repository repository.PostgreRepository) []core.ProductStock {
	var client = repository.CreateClientDatabase()
	defer client.Close()
	rows, err := client.Query(SET_PRODUCT_STOCK_BY_DEPOSIT_AND_PRODUCT, Deposit, Location)
	if err != nil {
		log.Println("No rows were returned!", err)
		return []core.ProductStock{}
	}
	defer rows.Close()
	var products []core.ProductStock
	for rows.Next() {
		var product core.ProductStock

		// unmarshal
		err := rows.Scan(&product.Deposit, &product.Location, &product.ProductId, &product.Quantity)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		// append the product
		products = append(products, product)
	}
	return products
}

func (p *ProductStockController) GetProduct(ProductId string, Deposit string, Location string, repository repository.PostgreRepository) (core.ProductStock, error) {
	var client = repository.CreateClientDatabase()
	defer client.Close()
	row := client.QueryRow(SET_PRODUCT_STOCK, ProductId, Deposit, Location)
	var product core.ProductStock
	err := row.Scan(&product.Deposit, &product.Location, &product.ProductId, &product.Quantity)
	if err != nil {
		log.Println("No rows were returned!", err)
		return core.ProductStock{}, err
	}

	return product, nil
}

func (p *ProductStockController) Delete(ps core.ProductStock, repository repository.PostgreRepository) core.ProductStock {
	var client = repository.CreateClientDatabase()
	defer client.Close()
	_, err := client.Exec(DELETE_PRODUCT_STOCK, ps.ProductId, ps.Deposit, ps.Location)
	if err != nil {
		return core.ProductStock{}
	}
	log.Println("Deleted a product")
	return ps
}

func (p *ProductStockController) Update(ps core.ProductStock, repository repository.PostgreRepository) core.ProductStock {
	var client = repository.CreateClientDatabase()
	defer client.Close()
	_, err := client.Exec(UPDATE_PRODUCT_STOCK, ps.Quantity, ps.ProductId, ps.Deposit, ps.Location)
	if err != nil {
		return core.ProductStock{}
	}
	log.Println("Updated a product")
	return ps
}

func (p *ProductStockController) DrawOutProduct(ps core.ProductStock, repository repository.PostgreRepository) (core.ProductStock, error) {
	currentProduct, err := p.GetProduct(ps.ProductId, ps.Deposit, ps.Location, repository)
	if err != nil {
		return core.ProductStock{}, err
	}
	if currentProduct.Quantity >= ps.Quantity {
		currentProduct.Quantity = currentProduct.Quantity - ps.Quantity
		if currentProduct.Quantity == 0 {
			p.Delete(currentProduct, repository)
		} else {
			p.Update(currentProduct, repository)
		}
		return currentProduct, nil
	}
	return ps, nil
}
