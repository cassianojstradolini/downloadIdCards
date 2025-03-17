package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"main/data"
	"main/to_image"
	"main/to_pdf"
	"net/http"
)

func main() {
	addr := ":8081"
	if err := StartServer(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// Server represents the PDF HTTP server
type Server struct {
	router      chi.Router
	idCardsResp data.IdCardsResponseSchema // Store ID cards data
}

// NewServer creates a new PDF server
func NewServer() *Server {
	s := &Server{
		router: chi.NewRouter(),
		idCardsResp: data.IdCardsResponseSchema{ // Initialize ID cards data
			Data: []data.IdCard{
				data.MockImageIdCardFront,
				data.MockImageIdCardBack,
				data.MockIdCardFront,
				data.MockIdCardBack,
				data.MockHTMLIdCardFront,
				data.MockHTMLIdCardBack,
				data.MockHTMLIdCardBoth,
			},
		},
	}
	s.routes()
	return s
}

// routes sets up all the routes for the PDF server
func (s *Server) routes() {
	s.router.Get("/pdf/idcards", s.handleGetIDCardsPDF())
	s.router.Get("/image/idcards", s.handleGetIDCardsImage())
	s.router.Get("/template-extension/idcards", s.handleGetIDCardsTemplateExtension())
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// writeResponse is a helper function to write content to response with proper error handling
func writeResponse(w http.ResponseWriter, content []byte, fileName string, contentType string) {
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprint(len(content)))

	if _, err := io.Copy(w, bytes.NewReader(content)); err != nil {
		log.Printf("Error writing response: %v", err)
		// We can't change the status code at this point as headers are already sent
	}
}

func (s *Server) handleGetIDCardsTemplateExtension() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation not provided
	}
}

func (s *Server) handleGetIDCardsImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate the merged image
		response, err := to_image.MergeImages(context.Background(), s.idCardsResp)
		if err != nil {
			http.Error(w, "Failed to generate image", http.StatusInternalServerError)
			return
		}

		writeResponse(w, response.ImageContent, response.FileName, "image/png")
	}
}

// handleGetIDCardsPDF returns a handler function for generating PDF from ID cards
func (s *Server) handleGetIDCardsPDF() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate the PDF
		response, err := to_pdf.GeneratePDFFromIDCards(context.Background(), s.idCardsResp)
		if err != nil {
			http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
			return
		}

		writeResponse(w, response.PDFContent, response.FileName, "application/pdf")
	}
}

// StartServer starts the PDF server on the specified address
func StartServer(addr string) error {
	server := NewServer()
	fmt.Printf("Starting PDF server on %s\n", addr)
	return http.ListenAndServe(addr, server)
}
