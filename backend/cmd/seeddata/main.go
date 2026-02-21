package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gourmet-guide/backend/internal/seed"
)

func main() {
	outDir := flag.String("out", "seed/output", "output folder for generated restaurant data")
	count := flag.Int("count", 100, "number of AI-generated menu items")
	bucket := flag.String("gcs-bucket", os.Getenv("GCS_BUCKET"), "target GCS bucket for generated images")
	projectID := flag.String("project-id", os.Getenv("GOOGLE_CLOUD_PROJECT"), "Google Cloud project id for Firestore writes")
	collection := flag.String("firestore-collection", "restaurants", "Firestore collection for generated restaurants")
	skipGCSUpload := flag.Bool("skip-gcs-upload", false, "skip uploading generated images to GCS")
	skipFirestoreWrite := flag.Bool("skip-firestore-write", false, "skip writing generated restaurants to Firestore")
	flag.Parse()

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		panic("GOOGLE_API_KEY is required for Gemini seed generation")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	provider := seed.NewGeminiProvider(apiKey)
	restaurants, images, err := seed.GenerateRestaurantsFromAI(ctx, provider, *count)
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Join(*outDir, "images"), 0o755); err != nil {
		panic(err)
	}

	for _, image := range images {
		outputPath := filepath.Join(*outDir, "images", image.FileName)
		if err := os.WriteFile(outputPath, image.Data, 0o644); err != nil {
			panic(err)
		}
	}

	if !*skipGCSUpload {
		if *bucket == "" {
			panic("gcs bucket is required unless --skip-gcs-upload is set")
		}
		uploader, err := seed.NewGCSImageUploader(ctx, *bucket)
		if err != nil {
			panic(err)
		}
		defer uploader.Close()

		for rIdx := range restaurants {
			for mIdx := range restaurants[rIdx].MenuItems {
				menuItem := &restaurants[rIdx].MenuItems[mIdx]
				img, ok := images[menuItem.ID]
				if !ok {
					panic(fmt.Sprintf("image bytes missing for menu item %s", menuItem.ID))
				}
				remoteURL, err := uploader.UploadMenuImage(ctx, restaurants[rIdx].ID, menuItem.ID, img.FileName, img.Data)
				if err != nil {
					panic(err)
				}
				menuItem.ImageURL = remoteURL
			}
		}
	}

	if !*skipFirestoreWrite {
		if *projectID == "" {
			panic("project id is required unless --skip-firestore-write is set")
		}
		store, err := seed.NewFirestoreRestaurantStore(ctx, *projectID, *collection)
		if err != nil {
			panic(err)
		}
		defer store.Close()
		for _, restaurant := range restaurants {
			if err := store.SaveRestaurant(ctx, restaurant); err != nil {
				panic(err)
			}
		}
	}

	payload, err := json.MarshalIndent(restaurants, "", "  ")
	if err != nil {
		panic(err)
	}

	jsonPath := filepath.Join(*outDir, "restaurants.json")
	if err := os.WriteFile(jsonPath, payload, 0o644); err != nil {
		panic(err)
	}

	fmt.Printf("seeded %d restaurants, %d menu images. output=%s\n", len(restaurants), len(images), *outDir)
}
