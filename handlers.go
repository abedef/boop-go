package boop

import (
	"context"
	"errors"
	"io"
	"junk/boop-server/pgdb"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func HandleTasks(w *MyResponseWriter, r *http.Request, ctx context.Context, queries *pgdb.Queries) error {
	log.Printf("Handling %v request for %v", r.Method, r.URL.Path)
	if r.URL.Query().Get("From") != "+14164521467" {
		log.Print("Tried to communicate with endpoint without the \"From\" query parameter")
		return errors.New("you are not the chosen one")
	}
	boops, err := queries.GetBoopsTasks(ctx)
	if err != nil {
		log.Printf("Error fetching boops: %v", err)
		return err
	}
	return w.WriteJSON(transformBoops(boops))
}

func HandleEvents(w *MyResponseWriter, r *http.Request, ctx context.Context, queries *pgdb.Queries) error {
	log.Printf("Handling %v request for %v", r.Method, r.URL.Path)
	if r.URL.Query().Get("From") != "+14164521467" {
		log.Print("Tried to communicate with endpoint without the \"From\" query parameter")
		return errors.New("you are not among the chosen ones")
	}
	boops, err := queries.GetBoops(ctx)
	if err != nil {
		log.Printf("Error fetching boops: %v", err)
		return err
	}
	filtered_boops := filterBoops(containsEvent, boops)
	return w.WriteJSON(transformBoops(filtered_boops))
}

func HandleBeans(w *MyResponseWriter, r *http.Request, ctx context.Context, queries *pgdb.Queries) error {
	log.Printf("Handling %v request for %v", r.Method, r.URL.Path)
	if r.URL.Query().Get("From") != "+14164521467" {
		log.Print("Tried to communicate with endpoint without the \"From\" query parameter")
		return errors.New("you are not among the chosen ones")
	}
	boops, err := queries.GetBoops(ctx)
	if err != nil {
		log.Printf("Error fetching boops: %v", err)
		return err
	}
	filtered_boops := filterBoops(containsBean, boops)
	beans := []Bean{}
	for _, boop := range filtered_boops {
		beans = append(beans, extractBeans(boop)...)
	}
	simplifiedBeans := simplifyBeans(beans)
	return w.WriteJSON(BeanSummary{Totals: simplifiedBeans, Boops: transformBoops(filtered_boops)})
}

func HandleRoot(w *MyResponseWriter, r *http.Request, ctx context.Context, queries *pgdb.Queries) error {
	log.Printf("Handling %v request for %v", r.Method, r.URL.Path)
	// John's number:         +16473366972
	switch r.Method {
	case "GET":
		if r.URL.Query().Get("From") != "+14164521467" {
			log.Print("Tried to communicate with endpoint without the \"From\" query parameter")
			return errors.New("you are not among the chosen ones")
		}

		folder := r.URL.Query().Get("folder")

		var boops []pgdb.Boop
		var err error
		if folder != "" {
			log.Print(folder)
			if strings.HasPrefix(folder, "@") {
				log.Printf("Querying specific Boop %v", folder)
				boopIdString := strings.TrimPrefix(folder, "@")
				boopId, err := strconv.ParseInt(boopIdString, 10, 32)
				if err != nil {
					log.Printf("Error parsing boop @%v", boopIdString)
					return err
				}
				boop, err := queries.GetBoop(ctx, int32(boopId))
				if err != nil {
					log.Printf("Error fetching boop: %v", err)
					return err
				}
				boops = []pgdb.Boop{
					boop,
				}
			} else {
				log.Printf("Narrowing Boops down to those in folder %v", folder)
				boops, err = queries.GetBoopsFolder(ctx, folder)
				if err != nil {
					log.Printf("Error fetching boops: %v", err)
					return err
				}
			}
		} else {
			boops, err = queries.GetBoops(ctx)
			if err != nil {
				log.Printf("Error fetching boops: %v", err)
				return err
			}
		}

		return w.WriteJSON(transformBoops(boops))
	case "PATCH":
		err := r.ParseForm()
		if err != nil {
			log.Print(err)
			return err
		}
		if r.URL.Query().Get("From") != "+14164521467" {
			log.Print("Tried to communicate with endpoint without the \"From\" query parameter")
			return errors.New("you are not among the chosen ones")
		}

		id, err := strconv.ParseInt(r.URL.Query().Get("Id"), 10, 32)
		if err != nil {
			log.Println(err)
			return err
		}
		body := r.Form.Get("Body")
		if body == "" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}

			body = string(bodyBytes)
		}
		if body == "" {
			log.Print("error")
			return errors.New("empty body")
		}
		_, err = queries.UpdateBoop(ctx,
			pgdb.UpdateBoopParams{ID: int32(id), Text: body},
		)
		if err != nil {
			log.Println(err)
			return err
		}
		return w.WriteHeader(http.StatusOK)
	case "DELETE":
		if r.URL.Query().Get("From") != "+14164521467" {
			log.Print("Tried to communicate with endpoint without the \"From\" query parameter")
			return errors.New("you are not among the chosen ones")
		}

		id, err := strconv.ParseInt(r.URL.Query().Get("Id"), 10, 32)
		if err != nil {
			log.Println(err)
			return err
		}
		err = queries.DeleteBoop(ctx, int32(id))
		if err != nil {
			log.Println(err)
			return err
		}
		return w.WriteHeader(http.StatusOK)
	case "POST":
		err := r.ParseForm()
		if err != nil {
			log.Print(err)
			return err
		}
		if r.Form.Get("From") != "+14164521467" {
			log.Printf("Tried to communicate with endpoint using form parameter \"From\" = \"%v\"", r.Form.Get("From"))
			return errors.New("you are not among the chosen ones")
		}

		body := r.Form.Get("Body")
		if body == "" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}

			body = string(bodyBytes)
		}
		if body == "" {
			log.Print("error")
			return errors.New("empty body")
		}

		_, err = queries.CreateBoop(ctx, body)
		if err != nil {
			log.Printf("Error creating boop: %v", err)
			return err
		}
		return w.WriteHeader(http.StatusOK)
	default:
		return w.WriteHeader(http.StatusBadRequest)
	}
}
