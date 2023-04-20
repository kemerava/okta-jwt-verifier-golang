/*******************************************************************************
 * Copyright 2018 - Present Okta, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 ******************************************************************************/

package lestrratGoJwx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kemerava/okta-jwt-verifier-golang/adaptors"
	"github.com/kemerava/okta-jwt-verifier-golang/utils"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
)

func fetchJwkSet(jwkUri string, client *http.Client) (interface{}, error) {
	return jwk.Fetch(context.Background(), jwkUri, jwk.WithHTTPClient(client))
}

type LestrratGoJwx struct {
	JWKSet      jwk.Set
	Cache       func(func(string, *http.Client) (interface{}, error), time.Duration, time.Duration) (utils.Cacher, error)
	jwkSetCache utils.Cacher
	Timeout     time.Duration
	Cleanup     time.Duration
}

func (lgj *LestrratGoJwx) New() adaptors.Adaptor {
	if lgj.Cache == nil {
		lgj.Cache = utils.NewDefaultCache
	}

	return lgj
}

func (lgj *LestrratGoJwx) GetKey(jwkUri string) {
}

func (lgj *LestrratGoJwx) Decode(jwt string, jwkUri string, client *http.Client) (interface{}, error) {
	if lgj.jwkSetCache == nil {
		jwkSetCache, err := lgj.Cache(fetchJwkSet, lgj.Timeout, lgj.Cleanup)
		if err != nil {
			return nil, err
		}
		lgj.jwkSetCache = jwkSetCache
	}

	value, err := lgj.jwkSetCache.Get(jwkUri, client)
	if err != nil {
		return nil, err
	}

	jwkSet, ok := value.(jwk.Set)
	if !ok {
		return nil, fmt.Errorf("could not cast %v to jwk.Set", value)
	}

	token, err := jws.VerifySet([]byte(jwt), jwkSet)
	if err != nil {
		return nil, err
	}

	var claims interface{}
	if err := json.Unmarshal(token, &claims); err != nil {
		return nil, fmt.Errorf("could not unmarshal claims: %w", err)
	}

	return claims, nil
}
