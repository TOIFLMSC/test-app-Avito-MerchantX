package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/TOIFLMSC/test-app-Avito-MerchantX/internal/app/model"
)

// Message func
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond func
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Error func
func Error(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	Respond(w, map[string]interface{}{"error": err.Error()})
}

// DownloadFile func
func DownloadFile(url string, w http.ResponseWriter) error {
	f, err := os.Create("../storage/test.xlsx")
	if err != nil {
		response := Message(false, "Unable to create xlsx file")
		Respond(w, response)
		return err
	}
	defer f.Close()

	c := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := c.Get(url)
	if err != nil {
		response := Message(false, fmt.Sprintf("Error while downloading %q: %v", url, err))
		Respond(w, response)
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		response := Message(false, "Unable to copy bytes from downloaded excel file to target file")
		Respond(w, response)
		return err
	}

	return nil
}

// UploadFile func
func UploadFile(w http.ResponseWriter) error {
	Openfile, err := os.Open("../storage/exporttest.xlsx")
	defer Openfile.Close()
	if err != nil {
		response := Message(false, "Unable to open created xlsx file")
		Respond(w, response)
		return nil
	}

	FileHeader := make([]byte, 512)

	Openfile.Read(FileHeader)

	FileContentType := http.DetectContentType(FileHeader)

	FileStat, err := Openfile.Stat()
	if err != nil {
		response := Message(false, "Unable to get xlsx file stats")
		Respond(w, response)
		return nil
	}

	FileSize := strconv.FormatInt(FileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename=exporttest.xlsx")
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	Openfile.Seek(0, 0)
	io.Copy(w, Openfile)

	return nil
}

// ConvertDataToExcel func
func ConvertDataToExcel(offerarray *[]model.Offer, w http.ResponseWriter) error {
	categories := map[string]string{"A1": "Offer_ID", "B1": "Name", "C1": "Price", "D1": "Quantity", "E1": "Available"}

	values := make(map[string]string)
	for i, v := range *offerarray {
		var idx string = strconv.Itoa(i + 2)
		values["A"+idx] = v.OfferID
		values["B"+idx] = v.Name
		values["C"+idx] = strconv.Itoa(v.Price)
		values["D"+idx] = strconv.Itoa(v.Quantity)
		values["E"+idx] = v.Available
	}

	f := excelize.NewFile()

	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}
	for k, v := range values {
		f.SetCellValue("Sheet1", k, v)
	}

	err := f.SaveAs("../storage/exporttest.xlsx")
	if err != nil {
		response := Message(false, "Unable to save xlsx file")
		Respond(w, response)
		return err
	}

	return nil
}

// ConvertRowToArray func
func ConvertRowToArray(rows *sql.Rows) ([]model.Offer, error) {
	array := []model.Offer{}
	for rows.Next() {
		var seller string
		var offerid string
		var name string
		var price int
		var quantity int
		var available string
		err := rows.Scan(&seller, &offerid, &name, &price, &quantity, &available)
		if err != nil {
			return nil, err
		} else {
			offer := model.Offer{seller, offerid, name, price, quantity, available}
			array = append(array, offer)
		}
	}
	return array, nil
}
