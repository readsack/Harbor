package routes

import (
	"context"
	"database/sql"
	"fmt"
	"harbor/main/db"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("JWT")
		if err != nil {
			http.Redirect(w, r, "/login", 303)
			return
		}

		tokenString := tokenCookie.Value
		if tokenString == "" {
			http.Redirect(w, r, "/login", 303)
			return
		}
		//fmt.Println(tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			secret := os.Getenv("SECRET")
			key := []byte(secret)
			return key, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			http.Redirect(w, r, "/login", 303)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && claims["key"] != nil {
			var key string = claims["key"].(string)
			//fmt.Println(key)
			u, _ := db.FindUserByKey(string(key))
			ctx := context.WithValue(r.Context(), "user", u)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", 303)
		}
	})
}

func createJWT(user_key string) string {
	var key []byte
	var t *jwt.Token
	var s string

	key = []byte(os.Getenv("SECRET"))
	t = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "my-auth-server",
			"key": user_key,
		})
	s, err := t.SignedString(key)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func signUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var username string = r.FormValue("username")
	var email string = r.FormValue("email")
	var pass string = r.FormValue("pass")
	_, err := db.FindUserByEmail(email)
	if err != sql.ErrNoRows {
		fmt.Fprintln(w, "Email Already Exists")
		return
	}
	usr := db.User{
		Username: username,
		Email:    email,
		Password: pass,
	}
	err = db.CreateUser(usr)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Write([]byte("Success!"))
}

func logIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var email string = r.FormValue("email")
	var pass string = r.FormValue("pass")
	fmt.Println(email, pass)
	u, err := db.FindUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintln(w, "No Such Account Exists")
		} else {
			fmt.Fprintln(w, "Internal Server Error")
		}
		return
	}
	if pass == u.Password {
		jwtToken := createJWT(u.Key)
		var c http.Cookie
		c.Name = "JWT"
		c.Value = jwtToken
		http.SetCookie(w, &c)
		fmt.Fprintln(w, "Logged In Successfully!")
	} else {
		fmt.Fprintln(w, "Failed To Login")
	}

}

func testingRoute(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello Testing!")
}

func SetupAuthRoutes() {
	http.HandleFunc("POST /signup", signUp)
	http.HandleFunc("POST /login", logIn)
	http.HandleFunc("/testroute", AuthMiddleware(testingRoute))
}
