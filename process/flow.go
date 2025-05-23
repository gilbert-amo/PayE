package process

import (
	"bufio"
	"fmt"
	"github.com/gilbert-amo/PayE/types"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	products  []types.Product
	workers   []types.Worker
	processID = 1
	productID = 1
	workerID  = 1
)

func AssignWorkers(p *types.Process, supervisor, staff string) {
	p.Supervisor = supervisor
	p.StaffWorker = staff
}

// Duration calculates the time taken for the process
func Duration(p *types.Process) time.Duration {
	if p.Status == "Completed" {
		return p.EndTime.Sub(p.StartTime)
	} else if p.Status == "In Progress" {
		return time.Since(p.StartTime)
	}
	return 0
}

// SetQuantity records the quantity produced in this process
func SetQuantity(p *types.Process, quantity int) {
	p.Quantity = quantity
}

// IsComplete returns true if the process is completed
func IsComplete(p *types.Process) bool {
	return p.Status == "Completed"
}

// HasQualityPass returns true if quality check passed
func HasQualityPass(p *types.Process) bool {
	return p.QualityCheck
}

func InitializeSampleData() {
	workers = append(workers, types.Worker{
		ID:      workerID,
		Name:    "John Smith",
		Role:    "Supervisor",
		Contact: "john@company.com",
		Shift:   "Morning",
	})
	workerID++

	workers = append(workers, types.Worker{
		ID:      workerID,
		Name:    "Alice Johnson",
		Role:    "Staff",
		Contact: "alice@company.com",
		Shift:   "Morning",
	})
	workerID++

	workers = append(workers, types.Worker{
		ID:      workerID,
		Name:    "Bob Williams",
		Role:    "Supervisor",
		Contact: "bob@company.com",
		Shift:   "Evening",
	})
	workerID++

	workers = append(workers, types.Worker{
		ID:      workerID,
		Name:    "Carol Brown",
		Role:    "Staff",
		Contact: "carol@company.com",
		Shift:   "Evening",
	})
	workerID++
}

func DisplayMainMenu() {
	fmt.Println("\nMain Menu:")
	fmt.Println("1. Create New Product")
	fmt.Println("2. Add Process to Product")
	fmt.Println("3. Assign Workers to Process")
	fmt.Println("4. Update Process Status")
	fmt.Println("5. Record Production Quantity")
	fmt.Println("6. Add New Worker")
	fmt.Println("7. View All Products")
	fmt.Println("8. Generate Production Report")
	fmt.Println("9. Exit")
}

func GetInput(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func CreateNewProduct() {
	fmt.Println("\nCreate New Product")
	fmt.Println("------------------")

	name, err := GetInput("Enter product name: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	description, err := GetInput("Enter product description: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	targetDateStr, err := GetInput("Enter target completion date (YYYY-MM-DD): ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	targetDate, err := time.Parse("2006-01-02", targetDateStr)
	if err != nil {
		fmt.Println("Invalid date format. Please use YYYY-MM-DD.")
		return
	}

	newProduct := types.Product{
		ID:          productID,
		Name:        name,
		Description: description,
		StartDate:   time.Now(),
		TargetDate:  targetDate,
	}

	productID++

	// Ask for number of processes
	numProcessesStr, err := GetInput("Enter number of processes required: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	numProcesses, err := strconv.Atoi(numProcessesStr)
	if err != nil || numProcesses <= 0 {
		fmt.Println("Invalid number of processes. Must be a positive integer.")
		return
	}

	// Add processes
	for i := 0; i < numProcesses; i++ {
		fmt.Printf("\nProcess %d:\n", i+1)
		processName, err := GetInput("Enter process name: ")
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		processDesc, err := GetInput("Enter process description: ")
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		newProcess := types.Process{
			ID:          processID,
			Name:        processName,
			Description: processDesc,
			Status:      "Pending",
			StartTime:   time.Time{}, // zero time indicates not started
			EndTime:     time.Time{}, // zero time indicates not completed
		}
		processID++

		newProduct.Processes = append(newProduct.Processes, newProcess)
	}

	products = append(products, newProduct)
	fmt.Printf("\nProduct '%s' created successfully with %d processes.\n", name, numProcesses)
}

func AddProcessToProduct() {
	if len(products) == 0 {
		fmt.Println("No products available. Please create a product first.")
		return
	}

	fmt.Println("\nAdd Process to Product")
	fmt.Println("----------------------")

	ViewAllProducts()

	productIDStr, err := GetInput("Enter product ID to add process to: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		fmt.Println("Invalid product ID.")
		return
	}

	var selectedProduct *types.Product
	for i := range products {
		if products[i].ID == productID {
			selectedProduct = &products[i]
			break
		}
	}

	if selectedProduct == nil {
		fmt.Println("Product not found.")
		return
	}

	fmt.Printf("\nAdding process to product: %s\n", selectedProduct.Name)

	processName, err := GetInput("Enter process name: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	processDesc, err := GetInput("Enter process description: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	newProcess := types.Process{
		ID:          processID,
		Name:        processName,
		Description: processDesc,
		Status:      "Pending",
		StartTime:   time.Time{},
		EndTime:     time.Time{},
	}
	processID++

	selectedProduct.Processes = append(selectedProduct.Processes, newProcess)
	fmt.Printf("\nProcess '%s' added to product '%s' successfully.\n", processName, selectedProduct.Name)
}

func AssignWorkersToProcess() {
	if len(products) == 0 {
		fmt.Println("No products available. Please create a product first.")
		return
	}

	fmt.Println("\nAssign Workers to Process")
	fmt.Println("-------------------------")

	ViewAllProducts()

	productIDStr, err := GetInput("Enter product ID: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		fmt.Println("Invalid product ID.")
		return
	}

	var selectedProduct *types.Product
	for i := range products {
		if products[i].ID == productID {
			selectedProduct = &products[i]
			break
		}
	}

	if selectedProduct == nil {
		fmt.Println("Product not found.")
		return
	}

	if len(selectedProduct.Processes) == 0 {
		fmt.Println("No processes available for this product.")
		return
	}

	fmt.Printf("\nProcesses for product '%s':\n", selectedProduct.Name)
	for _, process := range selectedProduct.Processes {
		fmt.Printf("%d. %s (Status: %s)\n", process.ID, process.Name, process.Status)
	}

	processIDStr, err := GetInput("Enter process ID to assign workers: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	processID, err := strconv.Atoi(processIDStr)
	if err != nil {
		fmt.Println("Invalid process ID.")
		return
	}

	var selectedProcess *types.Process
	for i := range selectedProduct.Processes {
		if selectedProduct.Processes[i].ID == processID {
			selectedProcess = &selectedProduct.Processes[i]
			break
		}
	}

	if selectedProcess == nil {
		fmt.Println("Process not found.")
		return
	}

	// Display available supervisors
	fmt.Println("\nAvailable Supervisors:")
	for _, worker := range workers {
		if worker.Role == "Supervisor" {
			fmt.Printf("%d. %s (Shift: %s)\n", worker.ID, worker.Name, worker.Shift)
		}
	}

	supervisorIDStr, err := GetInput("Enter supervisor ID: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	supervisorID, err := strconv.Atoi(supervisorIDStr)
	if err != nil {
		fmt.Println("Invalid supervisor ID.")
		return
	}

	var selectedSupervisor *types.Worker
	for i := range workers {
		if workers[i].ID == supervisorID && workers[i].Role == "Supervisor" {
			selectedSupervisor = &workers[i]
			break
		}
	}

	if selectedSupervisor == nil {
		fmt.Println("Supervisor not found or not a supervisor.")
		return
	}

	// Display available staff
	fmt.Println("\nAvailable Staff Workers:")
	for _, worker := range workers {
		if worker.Role == "Staff" {
			fmt.Printf("%d. %s (Shift: %s)\n", worker.ID, worker.Name, worker.Shift)
		}
	}

	staffIDStr, err := GetInput("Enter staff worker ID: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	staffID, err := strconv.Atoi(staffIDStr)
	if err != nil {
		fmt.Println("Invalid staff worker ID.")
		return
	}

	var selectedStaff *types.Worker
	for i := range workers {
		if workers[i].ID == staffID && workers[i].Role == "Staff" {
			selectedStaff = &workers[i]
			break
		}
	}

	if selectedStaff == nil {
		fmt.Println("Staff worker not found or not a staff member.")
		return
	}

	// Assign workers to process
	selectedProcess.Supervisor = selectedSupervisor.Name
	selectedProcess.StaffWorker = selectedStaff.Name

	fmt.Printf("\nAssigned workers to process '%s':\n", selectedProcess.Name)
	fmt.Printf("Supervisor: %s\n", selectedProcess.Supervisor)
	fmt.Printf("Staff Worker: %s\n", selectedProcess.StaffWorker)
}

func UpdateProcessStatus() {
	if len(products) == 0 {
		fmt.Println("No products available. Please create a product first.")
		return
	}

	fmt.Println("\nUpdate Process Status")
	fmt.Println("---------------------")

	ViewAllProducts()

	productIDStr, err := GetInput("Enter product ID: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		fmt.Println("Invalid product ID.")
		return
	}

	var selectedProduct *types.Product
	for i := range products {
		if products[i].ID == productID {
			selectedProduct = &products[i]
			break
		}
	}

	if selectedProduct == nil {
		fmt.Println("Product not found.")
		return
	}

	if len(selectedProduct.Processes) == 0 {
		fmt.Println("No processes available for this product.")
		return
	}

	fmt.Printf("\nProcesses for product '%s':\n", selectedProduct.Name)
	for _, process := range selectedProduct.Processes {
		fmt.Printf("%d. %s (Status: %s)\n", process.ID, process.Name, process.Status)
	}

	processIDStr, err := GetInput("Enter process ID to update: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	processID, err := strconv.Atoi(processIDStr)
	if err != nil {
		fmt.Println("Invalid process ID.")
		return
	}

	var selectedProcess *types.Process
	for i := range selectedProduct.Processes {
		if selectedProduct.Processes[i].ID == processID {
			selectedProcess = &selectedProduct.Processes[i]
			break
		}
	}

	if selectedProcess == nil {
		fmt.Println("Process not found.")
		return
	}

	fmt.Println("\nCurrent Status:", selectedProcess.Status)
	fmt.Println("1. Mark as In Progress")
	fmt.Println("2. Mark as Completed")
	fmt.Println("3. Add Quality Check Note")
	fmt.Println("4. Cancel")

	choice, err := GetInput("Enter your choice: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	switch choice {
	case "1":
		if selectedProcess.Status == "Pending" {
			selectedProcess.Status = "In Progress"
			selectedProcess.StartTime = time.Now()
			fmt.Println("Process marked as In Progress.")
		} else {
			fmt.Println("Process can only be started from Pending status.")
		}
	case "2":
		if selectedProcess.Status == "In Progress" {
			selectedProcess.Status = "Completed"
			selectedProcess.EndTime = time.Now()

			// Ask for quality check
			qualityCheck, err := GetInput("Did the process pass quality check? (y/n): ")
			if err == nil && strings.ToLower(qualityCheck) == "y" {
				selectedProcess.QualityCheck = true
			} else {
				selectedProcess.QualityCheck = false
			}

			notes, err := GetInput("Enter any notes for this process: ")
			if err == nil {
				selectedProcess.Notes = notes
			}

			fmt.Println("Process marked as Completed with quality check recorded.")
		} else {
			fmt.Println("Process can only be completed from In Progress status.")
		}
	case "3":
		notes, err := GetInput("Enter quality check notes: ")
		if err == nil {
			selectedProcess.Notes = notes
			fmt.Println("Quality check notes added.")
		}
	case "4":
		fmt.Println("Operation cancelled.")
	default:
		fmt.Println("Invalid choice.")
	}
}

func RecordProductionQuantity() {
	if len(products) == 0 {
		fmt.Println("No products available. Please create a product first.")
		return
	}

	fmt.Println("\nRecord Production Quantity")
	fmt.Println("--------------------------")

	ViewAllProducts()

	productIDStr, err := GetInput("Enter product ID: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		fmt.Println("Invalid product ID.")
		return
	}

	var selectedProduct *types.Product
	for i := range products {
		if products[i].ID == productID {
			selectedProduct = &products[i]
			break
		}
	}

	if selectedProduct == nil {
		fmt.Println("Product not found.")
		return
	}

	quantityStr, err := GetInput("Enter quantity produced: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		fmt.Println("Invalid quantity. Must be a positive integer.")
		return
	}

	selectedProduct.TotalQuantity += quantity

	// Update quantities for completed processes
	for i := range selectedProduct.Processes {
		if selectedProduct.Processes[i].Status == "Completed" {
			selectedProduct.Processes[i].Quantity = quantity
		}
	}

	fmt.Printf("Recorded %d units for product '%s'. Total produced: %d\n",
		quantity, selectedProduct.Name, selectedProduct.TotalQuantity)
}

func AddNewWorker() {
	fmt.Println("\nAdd New Worker")
	fmt.Println("--------------")

	name, err := GetInput("Enter worker name: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	fmt.Println("Select role:")
	fmt.Println("1. Supervisor")
	fmt.Println("2. Staff Worker")
	roleChoice, err := GetInput("Enter choice: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	var role string
	switch roleChoice {
	case "1":
		role = "Supervisor"
	case "2":
		role = "Staff"
	default:
		fmt.Println("Invalid choice. Defaulting to Staff.")
		role = "Staff"
	}

	contact, err := GetInput("Enter contact information: ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	shift, err := GetInput("Enter shift (Morning/Evening/Night): ")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	newWorker := types.Worker{
		ID:      workerID,
		Name:    name,
		Role:    role,
		Contact: contact,
		Shift:   shift,
	}
	workerID++

	workers = append(workers, newWorker)
	fmt.Printf("\nWorker '%s' added successfully as %s.\n", name, role)
}

func ViewAllProducts() {
	if len(products) == 0 {
		fmt.Println("No products available.")
		return
	}

	fmt.Println("\nList of Products:")
	fmt.Println("-----------------")
	for _, product := range products {
		fmt.Printf("\nID: %d\nName: %s\nDescription: %s\n",
			product.ID, product.Name, product.Description)
		fmt.Printf("Start Date: %s\nTarget Date: %s\n",
			product.StartDate.Format("2006-01-02"), product.TargetDate.Format("2006-01-02"))
		fmt.Printf("Total Quantity Produced: %d\n", product.TotalQuantity)

		fmt.Println("\nProcesses:")
		for _, process := range product.Processes {
			fmt.Printf("- %s (ID: %d, Status: %s)\n", process.Name, process.ID, process.Status)
			if process.Supervisor != "" {
				fmt.Printf("  Supervisor: %s, Staff: %s\n", process.Supervisor, process.StaffWorker)
			}
			if process.Status == "Completed" {
				fmt.Printf("  Completed on: %s\n", process.EndTime.Format("2006-01-02 15:04"))
				fmt.Printf("  Quality Check: %t\n", process.QualityCheck)
				if process.Notes != "" {
					fmt.Printf("  Notes: %s\n", process.Notes)
				}
			}
		}
	}
}

func GenerateProductionReport() {
	if len(products) == 0 {
		fmt.Println("No products available to generate report.")
		return
	}

	fmt.Println("\nProduction Report")
	fmt.Println("-----------------")
	fmt.Println("Generated on:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("==========================================")

	for _, product := range products {
		fmt.Printf("\nProduct: %s (ID: %d)\n", product.Name, product.ID)
		fmt.Printf("Description: %s\n", product.Description)
		fmt.Printf("Start Date: %s | Target Date: %s\n",
			product.StartDate.Format("2006-01-02"), product.TargetDate.Format("2006-01-02"))
		fmt.Printf("Total Quantity Produced: %d\n", product.TotalQuantity)

		// Calculate progress
		completedProcesses := 0
		for _, process := range product.Processes {
			if process.Status == "Completed" {
				completedProcesses++
			}
		}
		progress := float64(completedProcesses) / float64(len(product.Processes)) * 100
		fmt.Printf("Progress: %.1f%% (%d/%d processes completed)\n",
			progress, completedProcesses, len(product.Processes))

		// Check if target date is approaching or passed
		daysRemaining := int(product.TargetDate.Sub(time.Now()).Hours() / 24)
		if daysRemaining < 0 {
			fmt.Printf("⚠️ Target date passed %d days ago!\n", -daysRemaining)
		} else if daysRemaining < 7 {
			fmt.Printf("⚠️ Target date in %d days! Urgent!\n", daysRemaining)
		} else if daysRemaining < 14 {
			fmt.Printf("⚠️ Target date in %d days. Approaching deadline.\n", daysRemaining)
		}

		fmt.Println("\nProcess Details:")
		for _, process := range product.Processes {
			statusIcon := "🟡" // yellow for in progress
			if process.Status == "Completed" {
				if process.QualityCheck {
					statusIcon = "🟢" // green for completed with quality pass
				} else {
					statusIcon = "🔴" // red for completed but quality fail
				}
			} else if process.Status == "Pending" {
				statusIcon = "⚪" // white for pending
			}

			fmt.Printf("%s %s (ID: %d)\n", statusIcon, process.Name, process.ID)
			fmt.Printf("  Status: %s\n", process.Status)
			if process.Supervisor != "" {
				fmt.Printf("  Supervisor: %s\n", process.Supervisor)
				fmt.Printf("  Staff Worker: %s\n", process.StaffWorker)
			}
			if process.Status == "In Progress" {
				fmt.Printf("  Started: %s\n", process.StartTime.Format("2006-01-02 15:04"))
				duration := time.Since(process.StartTime)
				fmt.Printf("  Duration: %.1f hours\n", duration.Hours())
			} else if process.Status == "Completed" {
				fmt.Printf("  Completed: %s\n", process.EndTime.Format("2006-01-02 15:04"))
				duration := process.EndTime.Sub(process.StartTime)
				fmt.Printf("  Duration: %.1f hours\n", duration.Hours())
				fmt.Printf("  Quality Check: %t\n", process.QualityCheck)
				if process.Notes != "" {
					fmt.Printf("  Notes: %s\n", process.Notes)
				}
			}
		}
		fmt.Println("------------------------------------------")
	}

	// Worker utilization report
	fmt.Println("\nWorker Utilization:")
	fmt.Println("-------------------")
	for _, worker := range workers {
		processCount := 0
		for _, product := range products {
			for _, process := range product.Processes {
				if process.Supervisor == worker.Name || process.StaffWorker == worker.Name {
					processCount++
				}
			}
		}
		fmt.Printf("%s (%s): Assigned to %d processes\n", worker.Name, worker.Role, processCount)
	}
}
