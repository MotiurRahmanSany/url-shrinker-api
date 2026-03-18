package handlers

import (
	"net/http"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/service"
)

type UrlHandler struct {
	urlService service.UrlService
}

func NewUrlHandler(urlService service.UrlService) *UrlHandler {
	return &UrlHandler{urlService: urlService}
}

func (h *UrlHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {

}

func (h *UrlHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {

}

func (h *UrlHandler) ListMyURLs(w http.ResponseWriter, r *http.Request) {

}

func (h *UrlHandler) GetURLDetails(w http.ResponseWriter, r *http.Request) {

}

func (h *UrlHandler) DeactivateURL(w http.ResponseWriter, r *http.Request) {

}

func (h *UrlHandler) UpdateURL(w http.ResponseWriter, r *http.Request) {

}
