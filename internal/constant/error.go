package constant

import "github.com/fahmialfareza/deals-dating-app-backend/internal/domain"

var (
	// 400
	UserHasAlreadyExist     = domain.WrapError(400, "User has been already exist!")
	SwipeHasAlreadyExist    = domain.WrapError(400, "Swipe has been already exists!")
	PurchaseHasAlreadyExist = domain.WrapError(400, "You has been already purchased!")
	UserNameIsRequired      = domain.WrapError(400, "Name is required!")
	UserEmailIsRequired     = domain.WrapError(400, "Email is required!")
	UserPasswordIsRequired  = domain.WrapError(400, "Password is required!")
	UserImageIsRequired     = domain.WrapError(400, "Image is required!")
	UserBioIsRequired       = domain.WrapError(400, "Bio is required!")
	SwipedIDIsRequired      = domain.WrapError(400, "Choosen Person is required!")
	SwipedTypeIsRequired    = domain.WrapError(400, "Like/Pass is required!")
	CanNotSwipeYourSelf     = domain.WrapError(400, "Can not swipe yourself!")

	// 401
	Unauthorized    = domain.WrapError(401, "You must be logged in!")
	TokenNotValid   = domain.WrapError(401, "Token is not valid!")
	PasswordIsWrong = domain.WrapError(401, "Wrong password!")

	// 403
	YouNeedToBePremiumMember = domain.WrapError(403, "You need to be premium member!")

	// 404
	UserAccountCanNotBeFound = domain.WrapError(404, "Account not found!")
	SwipeCanNotBeFound       = domain.WrapError(404, "Swipe not found!")
)
