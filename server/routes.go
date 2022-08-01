package server

import (
	"Rms/handlers"
	"Rms/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	chi.Router
}

func SetupRoutes() *Server {
	router := chi.NewRouter()
	router.Route("/api", func(api chi.Router) {
		api.Post("/signup", handlers.SignUp)
		api.Post("/login", handlers.Login)

		api.Route("/rms", func(r chi.Router) {

			r.Use(middleware.AuthMiddleware)
			r.Route("/admin", func(admin chi.Router) {
				admin.Use(middleware.AdminMiddleware)
				admin.Get("/users", handlers.FetchAllUserRole)
				admin.Route("/{userId}", func(user chi.Router) {
					user.Put("/", handlers.ChangeUserRole)
				})

				admin.Route("/restaurant", func(restaurant chi.Router) {
					restaurant.Post("/", handlers.CreateRestaurant)
					restaurant.Get("/", handlers.FetchAllRestaurant)
					restaurant.Route("/{id}", func(restaurant chi.Router) {
						restaurant.Get("/showDish", handlers.FetchAllDish)
						restaurant.Put("/", handlers.UpdateRestaurant)
						restaurant.Delete("/", handlers.DeleteRestaurant)
						restaurant.Route("/dish", func(dish chi.Router) {
							dish.Post("/create", handlers.CreateDish)
							dish.Route("/{dishId}", func(dishes chi.Router) {
								dishes.Delete("/", handlers.DeleteDish)
								dishes.Put("/", handlers.UpdateDish)
							})
						})
					})

				})
			})
			r.Route("/sub-admin", func(subAdmin chi.Router) {
				subAdmin.Use(middleware.SubAdminMiddleware)
				subAdmin.Route("/restaurant", func(sub chi.Router) {
					sub.Post("/", handlers.CreateRestaurant)
					sub.Get("/", handlers.FetchAllRestaurantSubAdmin)
					sub.Route("/{id}", func(restaurant chi.Router) {
						restaurant.Delete("/", handlers.DeleteRestaurantSubAdmin)
						restaurant.Put("/", handlers.UpdateRestaurantSubAdmin)

						restaurant.Route("/dish", func(dish chi.Router) {
							dish.Post("/create", handlers.CreateDish)
							dish.Get("/", handlers.FetchAllDishSubAdmin)
							dish.Route("/{dishId}", func(dishes chi.Router) {

								dishes.Delete("/", handlers.DeleteDishSubAdmin)
								dishes.Put("/", handlers.UpdateDishSubAdmin)
							})
						})

					})

				})
			})
			r.Route("/location", func(l chi.Router) {
				l.Post("/create", handlers.SetLocation)
				l.Route("/{restaurantId}", func(d chi.Router) {
					d.Get("/distance", handlers.GetDistance)
				})
			})
			r.Route("/restaurant", func(r chi.Router) {
				r.Get("/", handlers.FetchAllRestaurant)
				r.Route("/{id}", func(d chi.Router) {
					d.Get("/", handlers.FetchAllDish)

				})
			})
			r.Get("/logout", handlers.SignOut)

		})
	})
	return &Server{router}

}
func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
