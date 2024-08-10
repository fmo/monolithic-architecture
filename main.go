package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func processOrder(db *sql.DB, productID int, quantity int, amount float64, address string) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer a rollback in case anything fails
	defer tx.Rollback()

	// Step 1: Create Order
	orderID, err := createOrder(tx, "new")
	if err != nil {
		return fmt.Errorf("order creation failed: %w", err)
	}

	// Step 2: Process Payment
	if err := processPayment(tx, orderID, amount); err != nil {
		return fmt.Errorf("payment processing failed: %w", err)
	}

	// Step 3: Update Inventory
	if err := updateInventory(tx, productID, quantity); err != nil {
		return fmt.Errorf("inventory update failed: %w", err)
	}

	// Step 4: Deliver Product
	if err := deliverProduct(tx, orderID, address); err != nil {
		return fmt.Errorf("delivery failed: %w", err)
	}

	// Commit the transaction if all steps succeed
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Println("Order processed successfully!")
	return nil
}

func createOrder(tx *sql.Tx, status string) (int, error) {
	var orderID int
	query := `INSERT INTO orders (status) VALUES ($1) RETURNING order_id`
	err := tx.QueryRow(query, status).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}
	fmt.Printf("Order created with ID %d\n", orderID)
	return orderID, nil
}

func processPayment(tx *sql.Tx, orderID int, amount float64) error {
	// Simulate payment processing and update the payments table
	query := `INSERT INTO payments (order_id, amount) VALUES ($1, $2)`
	_, err := tx.Exec(query, orderID, amount)
	if err != nil {
		return fmt.Errorf("failed to process payment: %w", err)
	}
	fmt.Printf("Payment processed for order %d: $%.2f\n", orderID, amount)
	return nil
}

func updateInventory(tx *sql.Tx, productID int, quantity int) error {
	// Simulate inventory update by decrementing the product quantity
	query := `UPDATE inventory SET quantity = quantity - $1 WHERE product_id = $2 AND quantity >= $1`
	result, err := tx.Exec(query, quantity, productID)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("not enough inventory for product %d", productID)
	}

	fmt.Printf("Inventory updated for product %d: -%d units\n", productID, quantity)
	return nil
}

func deliverProduct(tx *sql.Tx, orderID int, address string) error {
	// Simulate the delivery process by inserting a record into the deliveries table
	query := `INSERT INTO deliveries (order_id, address, status) VALUES ($1, $2, $3)`
	_, err := tx.Exec(query, orderID, address, "in progress")
	if err != nil {
		return fmt.Errorf("failed to record delivery: %w", err)
	}

	// For this example, we'll simulate a delivery failure
	fmt.Printf("Attempting delivery for order %d to %s\n", orderID, address)
	return fmt.Errorf("delivery failed for order %d", orderID)
}

func main() {
	// Replace with your database connection details
	connStr := "user=user dbname=shop password=password sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	productID := 1
	quantity := 2
	amount := 100.00
	address := "123 Main St, Anytown, USA"

	err = processOrder(db, productID, quantity, amount, address)
	if err != nil {
		log.Fatalf("Order processing failed: %v", err)
	}
}
