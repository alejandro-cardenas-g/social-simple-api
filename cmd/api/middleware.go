package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/alejandro-cardenas-g/social/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicError(w, r, err)
				return
			}

			creds := strings.SplitN(string(decoded), ":", 2)

			username := app.config.auth.basic.user
			password := app.config.auth.basic.password

			fmt.Println(username, password, creds)

			if username == "" || password == "" {
				app.internalServerError(w, r, errors.New("basic auth is not configured"))
				return
			}

			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unauthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedError(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				app.unauthorizedError(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			token := parts[1]

			jwtToken, err := app.authenticator.ValidateToken(token)
			if err != nil {
				app.unauthorizedBasicError(w, r, err)
				return
			}

			claims := jwtToken.Claims.(jwt.MapClaims)

			userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)

			if err != nil {
				app.unauthorizedBasicError(w, r, err)
				return
			}

			ctx := r.Context()

			user, err := app.getUser(ctx, userID)

			if err != nil {
				app.unauthorizedBasicError(w, r, err)
				return
			}

			ctx = context.WithValue(ctx, userCtx, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (app *application) CheckPostOwnershipMiddleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := getUserFromCtx(r)
		post := getPostFromCtx(r)

		if post.UserId == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenError(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}
	return user.Role.Level >= role.Level, nil
}

func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {
	if !app.config.redisCfg.enabled {
		return app.store.Users.GetByID(ctx, userID)
	}

	cacheUser, err := app.cacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cacheUser == nil {
		dbUser, err := app.store.Users.GetByID(ctx, userID)

		if err != nil {
			return nil, err
		}

		if err := app.cacheStorage.Users.Set(ctx, dbUser); err != nil {
			return nil, err
		}
		return dbUser, nil
	}
	return cacheUser, nil
}
