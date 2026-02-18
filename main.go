package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

type Expense struct {
	Item     string  `json:"item"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func main() {
	// 1. Iniciamos el motor de trazas (simplificado)
	tp := initTracer() // (Función que configura el envío a Tempo)
	defer tp.Shutdown(context.Background())

	r := gin.Default()
	r.Use(otelgin.Middleware("travel-budget-api"))

	// EL ENDPOINT ÚTIL
	r.POST("/expense", func(c *gin.Context) {
		var exp Expense
		if err := c.ShouldBindJSON(&exp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Simulamos que vamos a buscar el tipo de cambio a una API externa
		// Este bloque es el que "observaremos" con OpenTelemetry
		ctx, span := otel.Tracer("budget").Start(c.Request.Context(), "fetch-exchange-rate")
		time.Sleep(200 * time.Millisecond) // Simula latencia de red en Asia
		span.End()

		// Aquí iría el insert a CloudNativePG (Postgres)
		logGasto(exp) 

		c.JSON(http.StatusCreated, gin.H{
			"status": "Gasto registrado en la DB",
			"item": exp.Item,
			"converted_amount": exp.Amount * 1.25, // Simulación de conversión
		})
	})

	r.Run(":8080")
}

func logGasto(e Expense) {
    // Acá iría tu log para Loki
}