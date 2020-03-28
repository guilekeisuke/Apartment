package controllers

import (
	"api/usecase"
)

type ContentController struct {
	Interactor usecase.ContentInteractor
}
