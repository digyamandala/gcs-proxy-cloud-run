// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package gcs

import (
	"context"
	"net/http"
	"time"

	storage "cloud.google.com/go/storage"
	"github.com/DomZippilli/gcs-proxy-cloud-function/common"
	"github.com/DomZippilli/gcs-proxy-cloud-function/filter"
	"github.com/rs/zerolog/log"
)

// Read returns objects from a GCS bucket, mapping the URL to object names.
// Media caching is bypassed.
func UploadFile(ctx context.Context, response http.ResponseWriter,
	request *http.Request, pipeline filter.Pipeline) {
	objectName := common.NormalizePath(request.URL.Path)
	// GENERATE SIGNED_URL
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}
	url, signedErr := gcs.Bucket(bucket).SignedURL(objectName, opts)
	if signedErr != nil {
		log.Error().Msgf("Bucket(%q).SignedURL: %w", bucket, signedErr)
	}
	log.Info().Msgf("bucket: %q; signed_url: %q", bucket, url)
	log.Info().Msgf("redirecting to: %q", url)
	http.Redirect(response, request, url, http.StatusMovedPermanently)
}
