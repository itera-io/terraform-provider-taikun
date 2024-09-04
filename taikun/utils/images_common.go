package utils

import (
	tkcore "github.com/itera-io/taikungoclient/client"
)

func FlattenTaikunImages(rawImages ...tkcore.CommonStringBasedDropdownDto) []map[string]interface{} {

	images := make([]map[string]interface{}, len(rawImages))
	for i, rawImage := range rawImages {
		images[i] = map[string]interface{}{
			"id":   rawImage.GetId(),
			"name": rawImage.GetName(),
		}
	}
	return images
}
