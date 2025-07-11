package utils

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"server/dto"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

func RandomUserAvatar(fullname string) string {
	return fmt.Sprintf("https://api.dicebear.com/6.x/initials/svg?seed=%s", fullname)
}

func GenerateOTP(length int) string {
	digits := "0123456789"
	var sb strings.Builder

	for i := 0; i < length; i++ {
		sb.WriteByte(digits[rand.Intn(len(digits))])
	}

	return sb.String()
}

func GenerateSlug(input string) string {

	slug := strings.ToLower(input)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	suffix := strconv.Itoa(rand.Intn(1_000_000))
	slug = slug + "-" + leftPad(suffix, "0", 6)

	return slug
}

func leftPad(s string, pad string, length int) string {
	for len(s) < length {
		s = pad + s
	}
	return s
}

func GenerateInvoiceNumber(paymentID uuid.UUID) string {
	timestamp := time.Now().Format("20060102")
	shortID := paymentID.String()[:8]
	return fmt.Sprintf("INV/%s/%s", timestamp, shortID)
}

func GenerateTicketPDF(ticket *dto.UserTicketResponse) string {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "E-Ticket")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Ticket ID: %s", ticket.ID))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Event ID: %s", ticket.EventID))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Ticket Type: %s", ticket.TicketName))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("QR Code: %s", ticket.QRCode))

	dir := "generated_pdfs"
	_ = os.MkdirAll(dir, 0755)
	filename := filepath.Join(dir, fmt.Sprintf("ticket_%s.pdf", ticket.ID))
	if err := pdf.OutputFileAndClose(filename); err != nil {
		return ""
	}
	return fmt.Sprintf("/%s", filename)
}
