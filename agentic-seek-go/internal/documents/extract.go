package documents

import (
	"io"
	"os"

	"github.com/chchench/textract"
	"github.com/shakeel/pdf2txt"
)

func ExtractTextFromPDF(file io.Reader) (string, error) {
	// pdf2txt requires a file path, so we need to save the reader to a temporary file.
	tmpFile, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, file); err != nil {
		return "", err
	}
	tmpFile.Close()

	return pdf2txt.Get(tmpFile.Name())
}

func ExtractTextFromDOCX(file io.Reader) (string, error) {
	// textract requires a file path, so we need to save the reader to a temporary file.
	tmpFile, err := os.CreateTemp("", "upload-*.docx")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, file); err != nil {
		return "", err
	}
	tmpFile.Close()

	return textract.Extract(tmpFile.Name())
}
