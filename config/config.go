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
package config

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/gcs"
	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/rs/zerolog/log"
)

// Setup will be called once at the start of the program.
func Setup() error {
	return gcs.Setup()
}

// GET will be called in main.go for GET requests
func GET(ctx context.Context, output http.ResponseWriter, input *http.Request) {
	if stringUtils.IsEmpty(input.Header.Get("x-lpse-id")) {
		http.Error(output, "x-lpse-id header not found", http.StatusBadRequest)
		return
	}
	log.Info().Msgf("GET triggered with path: %q", input.URL.Path)

	requestHeadersJson, err := json.Marshal(input.Header)
	if err != nil {
		log.Error().Msgf("ERROR parsing header to json")
	}
	log.Info().Msgf("request header: %q", string(requestHeadersJson))

	if strings.Contains(input.URL.Path, "/public/") {
		gcs.Read(ctx, output, input, LoggingOnly)
	} else {
		gcs.ReadWithSignatureURL(ctx, output, input, LoggingOnly)
	}
	//gcs.ReadWithCache(ctx, output, input, CacheMedia, cacheGetter, LoggingOnly)
}

// HEAD will be called in main.go for HEAD requests
func HEAD(ctx context.Context, output http.ResponseWriter, input *http.Request) {
	gcs.ReadMetadata(ctx, output, input, LoggingOnly)
}

// func POST

// func DELETE

// OPTIONS will be called in main.go for OPTIONS requests
// func OPTIONS(ctx context.Context, output http.ResponseWriter, input *http.Request) {
// 	proxy.SendOptions(ctx, output, input, LoggingOnly)
// }
