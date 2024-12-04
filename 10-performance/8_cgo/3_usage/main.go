package main

import (
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/single_threaded"
)

// Be sure to close pools/instances when you're done with them.
var pool pdfium.Pool
var instance pdfium.Pdfium

func init() {
	// Init the PDFium library and return the instance to open documents.
	pool = single_threaded.Init(single_threaded.Config{})

	var err error
	instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	filePath := "example.pdf"
	output := "example.pdf.png"
	err := renderPage(filePath, 1, output)
	if err != nil {
		log.Fatal(err)
	}
}

func renderPage(filePath string, page int, output string) error {
	// Load the PDF file into a byte array.
	pdfBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Open the PDF using PDFium (and claim a worker)
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return err
	}

	// Always close the document, this will release its resources.
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	// Render the page in DPI 200.
	pageRender, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
		DPI: 200, // The DPI to render the page in.
		Page: requests.Page{
			ByIndex: &requests.PageByIndex{
				Document: doc,
				Index:    0,
			},
		}, // The page to render, 0-indexed.
	})
	if err != nil {
		return err
	}

	// Write the output to a file.
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, pageRender.Image)
	if err != nil {
		return err
	}

	return nil
}
