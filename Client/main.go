package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Drug struct {
	SerialNumber string `json:"serialNumber"`
	DrugName     string `json:"drugName"`
	Manufacturer string `json:"manufacturer"`
}

type PrivateDrug struct {
	SerialNumber string `json:"serialNumber"` // Unique ID
	DrugName     string `json:"drugName"`     // Name of the drug
	Manufacturer string `json:"manufacturer"` // Manufacturer of the drug
	BatchNumber  string `json:"batchNumber"`  // Batch number
	ExpiryDate   string `json:"expiryDate"`   // Expiry date
	Quantity     int    `json:"quantity"`     // Quantity
	Status       string `json:"status"`       // Drug status (e.g., Manufactured, Shipped, Sold)
}

type DrugData struct {
	SerialNumber string `json:"serialNumber"`
	DrugName     string `json:"drugName"`
	Manufacturer string `json:"manufacturer"`
	Status       string `json:"status"` // e.g., "Manufactured", "Shipped", "Received", "Sold"
}

type PrivateDrugData struct {
	SerialNumber string `json:"serialNumber"` // Unique ID
	DrugName     string `json:"drugName"`     // Name of the drug
	Manufacturer string `json:"manufacturer"` // Manufacturer of the drug
	BatchNumber  string `json:"batchNumber"`  // Batch number
	ExpiryDate   string `json:"expiryDate"`   // Expiry date
	Quantity     int    `json:"quantity"`     // Quantity
	Status       string `json:"status"`       // Drug status (e.g., Manufactured, Shipped, Sold)
}

type HistoryQueryResult struct {
	Record    *Drug  `json:"record"`
	TxId      string `json:"txId"`
	Timestamp string `json:"timestamp"`
	IsDelete  bool   `json:"isDelete"`
}

func main() {
	router := gin.Default()

	var wg sync.WaitGroup
	wg.Add(1)
	go ChaincodeEventListener("manufacturer", "autochannel", "Project-Pharma", &wg)

	// router.Static("/public", "./public")
	// router.LoadHTMLGlob("templates/*")

	router.GET("/", func(ctx *gin.Context) {
		result := submitTxnFn("manufacturer", "autochannel", "Project-Pharma", "DrugContract", "query", make(map[string][]byte), "GetDrugHistory")

		var drug []Drug

		if len(result) > 0 {
			// Unmarshal the JSON array string into the cars slice
			if err := json.Unmarshal([]byte(result), &drug); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"title": "Auto App", "drugList": drug,
		})
	})

	// router.GET("/create/drug", func(ctx *gin.Context) {
	// 	var req Drug
	// 	result := submitTxnFn("manufacturer", "autochannel", "Project-Pharma", "DrugContract", "invoke", make(map[string][]byte), "CreateDrug")

	// 	var drug [] Drug

	// 	if len(result) > 0 {
	// 		// Unmarshal the JSON array string into the cars slice
	// 		if err := json.Unmarshal([]byte(result), &drug); err != nil {
	// 			fmt.Println("Error:", err)
	// 			return
	// 		}
	// 	}

	// 	ctx.JSON(http.StatusOK,gin.H{
	// 		"title": "Manufacturer Dashboard", "drugList": drug,
	// 	})
	// })

	router.POST("/api/create/drug", func(ctx *gin.Context) {
		var req Drug
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}

		result := submitTxnFn("manufacturer", "autochannel", "Project-Pharma", "DrugContract", "invoke", make(map[string][]byte), "CreateDrug", req.SerialNumber, req.DrugName, req.Manufacturer)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Created Drug", "result": result,
		})
	})

	// 	router.GET("/api/car/:id", func(ctx *gin.Context) {
	// 		carId := ctx.Param("id")

	// 		result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "query", make(map[string][]byte), "ReadCar", carId)

	// 		ctx.JSON(http.StatusOK, gin.H{"data": result})
	// 	})

	// 	router.GET("/api/order/match-car", func(ctx *gin.Context) {
	// 		carID := ctx.Query("carId")
	// 		result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "query", make(map[string][]byte), "GetMatchingOrders", carID)

	// 		// fmt.Printf("result %s", result)

	// 		var orders []OrderData

	// 		if len(result) > 0 {
	// 			// Unmarshal the JSON array string into the orders slice
	// 			if err := json.Unmarshal([]byte(result), &orders); err != nil {
	// 				fmt.Println("Error:", err)
	// 				return
	// 			}
	// 		}

	// 		ctx.HTML(http.StatusOK, "matchOrder.html", gin.H{
	// 			"title": "Matching Orders", "orderList": orders, "carId": carID,
	// 		})
	// 	})

	// 	router.POST("/api/car/match-order", func(ctx *gin.Context) {
	// 		var req Match
	// 		if err := ctx.BindJSON(&req); err != nil {
	// 			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
	// 			return
	// 		}

	// 		fmt.Printf("match  %s", req)
	// 		submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "invoke", make(map[string][]byte), "MatchOrder", req.CarId, req.OrderId)

	// 		ctx.JSON(http.StatusOK, req)
	// 	})

	// 	router.GET("/api/event", func(ctx *gin.Context) {
	// 		result := getEvents()
	// 		fmt.Println("result:", result)

	// 		ctx.JSON(http.StatusOK, gin.H{"carEvent": result})

	// 	})

	// 	router.GET("/dealer", func(ctx *gin.Context) {

	// 		ctx.HTML(http.StatusOK, "dealer.html", gin.H{
	// 			"title": "Dealer Dashboard",
	// 		})
	// 	})

	// 	//Get all orders
	// 	router.GET("/api/order/all", func(ctx *gin.Context) {

	// 		result := submitTxnFn("dealer", "autochannel", "KBA-Automobile", "OrderContract", "query", make(map[string][]byte), "GetAllOrders")

	// 		var orders []OrderData

	// 		if len(result) > 0 {
	// 			// Unmarshal the JSON array string into the orders slice
	// 			if err := json.Unmarshal([]byte(result), &orders); err != nil {
	// 				fmt.Println("Error:", err)
	// 				return
	// 			}
	// 		}

	// 		ctx.HTML(http.StatusOK, "orders.html", gin.H{
	// 			"title": "All Orders", "orderList": orders,
	// 		})
	// 	})

	// 	router.POST("/api/order", func(ctx *gin.Context) {
	// 		var req Order
	// 		if err := ctx.BindJSON(&req); err != nil {
	// 			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
	// 			return
	// 		}

	// 		fmt.Printf("order  %s", req)

	// 		privateData := map[string][]byte{
	// 			"make":       []byte(req.Make),
	// 			"model":      []byte(req.Model),
	// 			"color":      []byte(req.Color),
	// 			"dealerName": []byte(req.Dealer),
	// 		}

	// 		submitTxnFn("dealer", "autochannel", "KBA-Automobile", "OrderContract", "private", privateData, "CreateOrder", req.OrderId)

	// 		ctx.JSON(http.StatusOK, req)
	// 	})

	// 	router.GET("/api/order/:id", func(ctx *gin.Context) {
	// 		orderId := ctx.Param("id")

	// 		result := submitTxnFn("dealer", "autochannel", "KBA-Automobile", "OrderContract", "query", make(map[string][]byte), "ReadOrder", orderId)

	// 		ctx.JSON(http.StatusOK, gin.H{"data": result})
	// 	})

	// 	router.GET("/mvd", func(ctx *gin.Context) {
	// 		result := submitTxnFn("mvd", "autochannel", "KBA-Automobile", "CarContract", "query", make(map[string][]byte), "GetAllCars")

	// 		var cars []CarData

	// 		if len(result) > 0 {
	// 			// Unmarshal the JSON array string into the cars slice
	// 			if err := json.Unmarshal([]byte(result), &cars); err != nil {
	// 				fmt.Println("Error:", err)
	// 				return
	// 			}
	// 		}

	// 		ctx.HTML(http.StatusOK, "mvd.html", gin.H{
	// 			"title": "MVD Dashboard", "carList": cars,
	// 		})
	// 	})

	// 	router.GET("/api/car/history", func(ctx *gin.Context) {
	// 		carID := ctx.Query("carId")
	// 		result := submitTxnFn("mvd", "autochannel", "KBA-Automobile", "CarContract", "query", make(map[string][]byte), "GetCarHistory", carID)

	// 		// fmt.Printf("result %s", result)

	// 		var cars []CarHistory

	// 		if len(result) > 0 {
	// 			// Unmarshal the JSON array string into the orders slice
	// 			if err := json.Unmarshal([]byte(result), &cars); err != nil {
	// 				fmt.Println("Error:", err)
	// 				return
	// 			}
	// 		}

	// 		ctx.HTML(http.StatusOK, "history.html", gin.H{
	// 			"title": "Car History", "itemList": cars,
	// 		})
	// 	})

	// 	router.POST("/api/car/register", func(ctx *gin.Context) {
	// 		var req Register
	// 		if err := ctx.BindJSON(&req); err != nil {
	// 			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
	// 			return
	// 		}

	// 		fmt.Printf("car response %s", req)
	// 		submitTxnFn("mvd", "autochannel", "KBA-Automobile", "CarContract", "invoke", make(map[string][]byte), "RegisterCar", req.CarId, req.CarOwner, req.RegNumber)

	// 		ctx.JSON(http.StatusOK, req)
	// 	})

	router.Run("localhost:8080")
}

// package main

// import "fmt"

// func main() {
// 	result := submitTxnFn(
// 		"manufacturer",
// 		"autochannel",
// 		"Project-Pharma",
// 		"DrugContract",
// 		"invoke",
// 		make(map[string][]byte),
// 		"CreateDrug",
// 		"01",
// 		"Aspirin",
// 		"BharatPharma",
// 		"Manufactured",
// 	)
// 	// privateData := map[string][]byte{
// 	// 	"drugName":       []byte("PainKiller"),
// 	// 	"manufacturer":      []byte("PharmaInc"),
// 	// 	"batchNumber":      []byte("BATCH001"),
// 	// 	"expiryDate": []byte("2025-12-31"),
// 	// 	"quantity": []byte("100"),
// 	// 	"status": []byte("Manufactured"),
// 	// }
// 	//--transient '{"serialNumber":"SN12345","drugName":"Painkiller","manufacturer":"PharmaInc","batchNumber":"BATCH001","expiryDate":"2025-12-31","quantity":100,"status":"Manufactured"}'
// 	// result := submitTxnFn("dealer", "autochannel", "KBA-Automobile", "OrderContract", "private", privateData, "CreateOrder", "ORD-03")

// 	// result := submitTxnFn("dealer", "autochannel", "KBA-Automobile", "OrderContract", "query", make(map[string][]byte), "ReadOrder", "ORD-03")

// 	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "query", make(map[string][]byte), "GetAllCars")

// 	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "OrderContract", "query", make(map[string][]byte), "GetAllOrders")

// 	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "query", make(map[string][]byte), "GetMatchingOrders", "Car-06")

// 	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "invoke", make(map[string][]byte), "MatchOrder", "Car-06", "ORD-03")

// 	// result := submitTxnFn("mvd", "autochannel", "KBA-Automobile", "CarContract", "invoke", make(map[string][]byte), "RegisterCar", "Car-06", "Dani", "KL-01-CD-01")

// 	//result = submitTxnFn("manufacturer", "autochannel", "Project-Pharma", "DrugContract", "query", make(map[string][]byte), "ReadDrug", "01")

// 	fmt.Println(result)

//}
