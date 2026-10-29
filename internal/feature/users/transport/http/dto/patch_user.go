package user_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/starfrus/golang-todoapp/internal/core/domain"
	core_logger "github.com/starfrus/golang-todoapp/internal/core/logger"
	core_http_request "github.com/starfrus/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/starfrus/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/starfrus/golang-todoapp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"     swaggertype:"string" example:"Максим Максимович"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"  swaggertype:"string" example:"+71112223344"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("'FullName' can`t be NILL")
		}
		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("'FullName' must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.FullName.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("'PhoneNumber' must be between 10 and 15 symbols")
			}
			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("'PhoneNumber' must startswitch '+' symbols")
			}
		}
	}
	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser	godoc
// @Summary		Изменение пользователя
// @Description	Изменение информации об уже существующем в системе пользователе.
// @Description	### Логика обновления полей (Three-state logic):
// @Description	1. **Поле не передано**: `phone_number` игнорируется, значение в БД не меняется.
// @Description	2. **Явно передано значение**: `"phone_number": "+71112223344"` - устанавливает новый номер.
// @Description	3. **Передан null**: `"phone_number": null` - очищает поле в БД (set to NULL).
// @Description	4. **Ограничение**: поле `full_name` не может быть установлено в `null`.
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		id		path		integer					true	"ID изменяемого пользователя"
// @Param		request	body		PatchUserRequest		true	"Данные для обновления пользователя"
// @Success		200		{object}	PatchUserResponse							"Успешно изменённый пользователь"
// @Failure		400		{object}	core_http_response.ErrorResponse			"Bad request"
// @Failure		404		{object}	core_http_response.ErrorResponse			"User not found"
// @Failure		409		{object}	core_http_response.ErrorResponse			"Conflict"
// @Failure		500		{object}	core_http_response.ErrorResponse			"Internal server error"
// @Router		/users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID to path value",
		)
		return
	}

	var request PatchUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userId, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
