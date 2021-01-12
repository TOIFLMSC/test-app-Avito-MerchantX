package apiserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/TOIFLMSC/test-app-Avito-MerchantX/internal/app/model"
	"github.com/TOIFLMSC/test-app-Avito-MerchantX/internal/app/store"
	u "github.com/TOIFLMSC/test-app-Avito-MerchantX/internal/app/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"})))
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/offer", s.addNewOffer()).Methods("POST")
	s.router.HandleFunc("/offer", s.getOfferList()).Methods("GET")
}

func (s *server) logRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})

		logger.Infof("started %s %s", r.Method, r.RequestURI)
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level

		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}

		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) addNewOffer() http.HandlerFunc {

	type request struct {
		Link   string `json:"link"`
		Seller string `json:"seller"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			u.Error(w, http.StatusBadRequest, err)
			return
		}

		err := u.DownloadFile(req.Link, w)
		if err != nil {
			response := u.Message(false, "Unable to download an excel file")
			u.Respond(w, response)
			return
		}

		file, err := excelize.OpenFile("storage/test.xlsx")
		if err != nil {
			response := u.Message(false, "Unable to open an excel file")
			u.Respond(w, response)
			return
		}

		var dataarray []string
		var convprice int
		var convquantity int
		var amountoferrors int = 0
		var amountofnewstrings int = 0
		var amountofupdatedstrings int = 0
		var amountofdeletedstrings int = 0

		rows, err := file.GetRows("List1")
		for idx, row := range rows {

			if idx < 1 {
				continue
			}

			for _, colCell := range row {
				dataarray = append(dataarray, colCell)
			}

			if flag, err := s.store.Offer().CheckForOffer(req.Seller, dataarray[0]); flag == true && err == store.ErrRecordNotFound {

				convprice, err = strconv.Atoi(dataarray[2])
				if err != nil {
					amountoferrors++
				} else if convprice <= 0 {
					amountoferrors++
					convprice = 0
				}

				convquantity, err = strconv.Atoi(dataarray[3])
				if err != nil {
					amountoferrors++
				} else if convquantity <= 0 {
					amountoferrors++
					convquantity = 0
				}

				offermodel := &model.Offer{
					Seller:    req.Seller,
					OfferID:   dataarray[0],
					Name:      dataarray[1],
					Price:     convprice,
					Quantity:  convquantity,
					Available: dataarray[4],
				}

				if err = s.store.Offer().NewOffer(offermodel); err != nil {
					u.Error(w, http.StatusUnprocessableEntity, err)
					return
				}

				amountofnewstrings++

			} else if flag, err := s.store.Offer().CheckForOffer(req.Seller, dataarray[0]); flag == false && err == nil {

				if dataarray[4] == "false" {

					err = s.store.Offer().DeleteOffer(req.Seller, dataarray[0])
					if err != nil {
						u.Error(w, http.StatusUnprocessableEntity, err)
						return
					}

					amountofdeletedstrings++

				} else if dataarray[4] == "true" {

					convprice, err = strconv.Atoi(dataarray[2])
					if err != nil {
						amountoferrors++
					} else if convprice <= 0 {
						amountoferrors++
						convprice = 0
					}

					convquantity, err = strconv.Atoi(dataarray[3])
					if err != nil {
						amountoferrors++
					} else if convquantity <= 0 {
						amountoferrors++
						convquantity = 0
					}

					offermodel := &model.Offer{
						Seller:    req.Seller,
						OfferID:   dataarray[0],
						Name:      dataarray[1],
						Price:     convprice,
						Quantity:  convquantity,
						Available: dataarray[4],
					}

					if err = s.store.Offer().UpdateOffer(offermodel); err != nil {
						u.Error(w, http.StatusUnprocessableEntity, err)
						return
					}

					amountofupdatedstrings++
				}
			}
			dataarray = dataarray[:0]
		}

		response := u.Message(true, "Product data is fully processed")
		response["created"] = amountofnewstrings
		response["updated"] = amountofupdatedstrings
		response["deleted"] = amountofdeletedstrings
		response["errors"] = amountoferrors
		u.Respond(w, response)

		return
	}
}

func (s *server) getOfferList() http.HandlerFunc {

	type request struct {
		Seller  string `json:"seller"`
		OfferID string `json:"offer_id"`
		Name    string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			response := u.Message(false, "Cannot search data, you may not have entered any search parametres, please enter at least one parameter.")
			u.Respond(w, response)
			return
		}

		switch {
		case req.Seller != "" && req.OfferID != "" && req.Name != "":
			offerarray, err := s.store.Offer().GetOfferWithAllSpecs(req.Seller, req.OfferID, req.Name)
			if err != nil {
				u.Error(w, http.StatusUnprocessableEntity, err)
				return
			}

			if len(*offerarray) == 0 {
				response := u.Message(false, "No results were found for your search parametres.")
				u.Respond(w, response)
				return
			}

			err = u.ConvertDataToExcel(offerarray, w)
			if err != nil {
				response := u.Message(false, "Unable to convert data to xlsx file")
				u.Respond(w, response)
				return
			}

			err = u.UploadFile(w)
			if err != nil {
				response := u.Message(false, "Unable to upload xlsx file")
				u.Respond(w, response)
				return
			}

			return

		case req.Seller == "" && req.OfferID != "" && req.Name != "":

			offerarray, err := s.store.Offer().GetOfferWithIDandName(req.OfferID, req.Name)
			if err != nil {
				u.Error(w, http.StatusUnprocessableEntity, err)
				return
			}

			if len(*offerarray) == 0 {
				response := u.Message(false, "No results were found for your search parametres.")
				u.Respond(w, response)
				return
			}

			err = u.ConvertDataToExcel(offerarray, w)
			if err != nil {
				response := u.Message(false, "Unable to convert data to xlsx file")
				u.Respond(w, response)
				return
			}

			err = u.UploadFile(w)
			if err != nil {
				response := u.Message(false, "Unable to upload xlsx file")
				u.Respond(w, response)
				return
			}

			return

		case req.Seller != "" && req.OfferID == "" && req.Name != "":
			offerarray, err := s.store.Offer().GetOfferWithSelandName(req.Seller, req.Name)
			if err != nil {
				u.Error(w, http.StatusUnprocessableEntity, err)
				return
			}

			if len(*offerarray) == 0 {
				response := u.Message(false, "No results were found for your search parametres.")
				u.Respond(w, response)
				return
			}

			err = u.ConvertDataToExcel(offerarray, w)
			if err != nil {
				response := u.Message(false, "Unable to convert data to xlsx file")
				u.Respond(w, response)
				return
			}

			err = u.UploadFile(w)
			if err != nil {
				response := u.Message(false, "Unable to upload xlsx file")
				u.Respond(w, response)
				return
			}

			return

		case req.Seller != "" && req.OfferID != "" && req.Name == "":
			offerarray, err := s.store.Offer().GetOfferWithSelandID(req.Seller, req.OfferID)
			if err != nil {
				u.Error(w, http.StatusUnprocessableEntity, err)
				return
			}

			if len(*offerarray) == 0 {
				response := u.Message(false, "No results were found for your search parametres.")
				u.Respond(w, response)
				return
			}

			err = u.ConvertDataToExcel(offerarray, w)
			if err != nil {
				response := u.Message(false, "Unable to convert data to xlsx file")
				u.Respond(w, response)
				return
			}

			err = u.UploadFile(w)
			if err != nil {
				response := u.Message(false, "Unable to upload xlsx file")
				u.Respond(w, response)
				return
			}

			return

		case req.Seller == "" && req.OfferID == "" && req.Name != "":
			offerarray, err := s.store.Offer().GetOfferWithName(req.Name)
			if err != nil {
				u.Error(w, http.StatusUnprocessableEntity, err)
				return
			}

			if len(*offerarray) == 0 {
				response := u.Message(false, "No results were found for your search parametres.")
				u.Respond(w, response)
				return
			}

			err = u.ConvertDataToExcel(offerarray, w)
			if err != nil {
				response := u.Message(false, "Unable to convert data to xlsx file")
				u.Respond(w, response)
				return
			}

			err = u.UploadFile(w)
			if err != nil {
				response := u.Message(false, "Unable to upload xlsx file")
				u.Respond(w, response)
				return
			}

			return

		case req.Seller == "" && req.OfferID != "" && req.Name == "":
			offerarray, err := s.store.Offer().GetOfferWithID(req.OfferID)
			if err != nil {
				u.Error(w, http.StatusUnprocessableEntity, err)
				return
			}

			if len(*offerarray) == 0 {
				response := u.Message(false, "No results were found for your search parametres.")
				u.Respond(w, response)
				return
			}

			err = u.ConvertDataToExcel(offerarray, w)
			if err != nil {
				response := u.Message(false, "Unable to convert data to xlsx file")
				u.Respond(w, response)
				return
			}

			err = u.UploadFile(w)
			if err != nil {
				response := u.Message(false, "Unable to upload xlsx file")
				u.Respond(w, response)
				return
			}

			return

		case req.Seller != "" && req.OfferID == "" && req.Name == "":
			offerarray, err := s.store.Offer().GetOfferWithSel(req.Seller)
			if err != nil {
				u.Error(w, http.StatusUnprocessableEntity, err)
				return
			}

			if len(*offerarray) == 0 {
				response := u.Message(false, "No results were found for your search parametres.")
				u.Respond(w, response)
				return
			}

			err = u.ConvertDataToExcel(offerarray, w)
			if err != nil {
				response := u.Message(false, "Unable to convert data to xlsx file")
				u.Respond(w, response)
				return
			}

			err = u.UploadFile(w)
			if err != nil {
				response := u.Message(false, "Unable to upload xlsx file")
				u.Respond(w, response)
				return
			}

			return
		}
	}
}
