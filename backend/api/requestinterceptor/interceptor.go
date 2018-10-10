package requestinterceptor

import (
	"nyota/backend/model"
	"net/http"
)

/*PrizmHandler will pass session info along with http request and writer*/
type PrizmHandler func(*model.SessionContext, http.ResponseWriter, *http.Request)

/*Interceptor acts as middle ware which will be called before invoking API*/
type Interceptor func(PrizmHandler) PrizmHandler
