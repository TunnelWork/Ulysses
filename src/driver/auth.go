package driver

import (
	"net/http"
	"net/mail"
	"strconv"

	"github.com/TunnelWork/Ulysses.Lib/api"
	"github.com/TunnelWork/Ulysses.Lib/auth"
	"github.com/TunnelWork/Ulysses.Lib/logging"
	"github.com/TunnelWork/Ulysses/src/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Auth struct {
	Affiliation Affiliation
	MFA         MFA
	User        User
}

type Affiliation struct { // abstract, static
}

func (Affiliation) GetByID(c *gin.Context, user *auth.User) {
	idstr := c.Query("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	// To inspect an affiliation, the user must be one of the following:
	// - Affiliation Account User that belongs to the corresponding affiliation
	// - Global Admin
	if (!user.Role.Includes(auth.AFFILIATION_ACCOUNT_USER) || user.AffiliationID != id) && !user.Role.Includes(auth.GLOBAL_ADMIN) {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}

	affiliation, err := auth.GetAffiliationByID(id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, affiliation))
}

func (Affiliation) Create(c *gin.Context, _ *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}

func (Affiliation) Update(c *gin.Context, _ *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}

func (Affiliation) ParentAffiliation(c *gin.Context, _ *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}

type MFA struct { // abstract, static
}

func (MFA) Registered(c *gin.Context, uid uint64) {
	mfaType := c.Query("type")
	// uid, err := utils.AuthorizationToUserID(c)
	// if err != nil {
	// 	logging.Debug("Can't get userID: %s", err.Error())
	// 	if err == sql.ErrNoRows {
	// 		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespUserNotFound)
	// 	} else {
	// 		utils.HandleError(c, err)
	// 	}
	// 	return
	// }
	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, auth.MFARegistered(mfaType, uid)))
}

func (MFA) InitSignUp(c *gin.Context, uid uint64) {
	var response interface{}
	var err error

	var form FormMfaInitSignUp = FormMfaInitSignUp{}
	err = c.ShouldBindBodyWith(&form, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	// Get User Info
	user, err := auth.GetUserByID(uid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	response, err = auth.MFAInitSignUp(form.Type, uid, user.Email)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, gin.H{
		"type":     form.Type,
		"response": response,
	}))
}

func (MFA) CompleteSignUp(c *gin.Context, uid uint64) {
	var form FormMfaCompleteSignUp = FormMfaCompleteSignUp{}

	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if err := auth.MFACompleteSignUp(form.Type, uid, form.Response); err != nil {
		utils.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, utils.RespOK)
}

func (MFA) NewChallenge(c *gin.Context, uid uint64) {
	var form FormMfaNewChallenge = FormMfaNewChallenge{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	var challengeMap map[string]interface{} = map[string]interface{}{}
	if auth.MFARegistered(form.Type, uid) {
		challenge, err := auth.MFANewChallenge(form.Type, uid)
		if err != nil {
			utils.HandleError(c, err)
			return
		} else {
			challengeMap[form.Type] = challenge
		}
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespMfaRequestInvalid)
		return
	}
	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, challengeMap))
}

func (MFA) SubmitChallenge(c *gin.Context, uid uint64) {
	var form FormMfaSubmitChallenge = FormMfaSubmitChallenge{
		Mfa: &MfaChallengeResponse{
			Type:     "",
			Response: map[string]string{},
		},
	}

	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		logging.Debug("A valid MFA response is expected, but not received.")
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespMfaResponseRequired)
		return
	}

	if err := auth.MFASubmitChallenge(form.Mfa.Type, uid, form.Mfa.Response); err != nil {
		logging.Debug("Failed to submit challenge: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespMfaResponseInvalid)
		return
	}
	c.JSON(http.StatusOK, utils.RespOK)
}

func (MFA) Remove(c *gin.Context, uid uint64) {
	var err error

	mfaType := c.Query("type")

	err = auth.MFARemove(mfaType, uid)
	if err != nil {
		utils.HandleError(c, err)
	} else {
		c.JSON(http.StatusOK, utils.RespOK)
	}
}

type User struct { // abstract, static
}

func (User) GetByID(c *gin.Context, currentUser *auth.User) {
	if currentUser == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespNotAuthorized)
		return
	}

	targetIDStr := c.Query("id")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil {
		utils.HandleError(c, err) // other unhandled error
		return
	}
	targetUser, err := auth.GetUserByID(targetID)
	if err != nil {
		utils.HandleError(c, err) // other unhandled error
		return
	}

	if currentUser.Role.Includes(auth.GLOBAL_ADMIN) { // Anyone
		logging.Debug("User(%d) as a Global Admin checked user(%d)", currentUser.ID(), targetID)
	} else if currentUser.Role.Includes(auth.AFFILIATION_ACCOUNT_USER) { // Same affiliation only
		if targetUser.AffiliationID != currentUser.AffiliationID {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied) // TODO: evaluate side channel attack vulnerability
			return
		}
	} else { // Self only
		if targetID != currentUser.ID() {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied) // TODO: evaluate side channel attack vulnerability
			return
		}
	}
	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, gin.H{
		"id":          targetUser.ID(),
		"email":       targetUser.Email,
		"role":        targetUser.Role,
		"affiliation": targetUser.AffiliationID,
	}))
}

func (User) List(c *gin.Context, currentUser *auth.User) {
	if currentUser == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespNotAuthorized)
		return
	}

	if currentUser.Role.Includes(auth.GLOBAL_ADMIN) { // Anyone
		logging.Debug("User(%d) as a Global Admin exported user list", currentUser.ID())
	} else { // Can't list
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied) // TODO: evaluate side channel attack vulnerability
		return
	}

	userList, err := auth.ListUserID()
	if err != nil {
		logging.Debug("Can't list users: %s", err.Error())
	}
	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, userList))
}

func (User) ListByAffiliation(c *gin.Context, currentUser *auth.User) {
	if currentUser == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespNotAuthorized)
		return
	}

	affiliationIDStr := c.Query("affiliation")
	affiliationID, err := strconv.ParseUint(affiliationIDStr, 10, 64)
	if err != nil {
		logging.Debug("Can't get affiliationID: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	if currentUser.Role.Includes(auth.GLOBAL_ADMIN) { // Anyone
		logging.Debug("User(%d) as a Global Admin exported user list from affiliation(%d)", currentUser.ID(), affiliationID)
	} else if currentUser.Role.Includes(auth.AFFILIATION_ACCOUNT_USER) { // Same affiliation only
		if currentUser.AffiliationID != affiliationID {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied) // TODO: evaluate side channel attack vulnerability
			return
		}
	} else { // Can't list
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied) // TODO: evaluate side channel attack vulnerability
		return
	}

	userList, err := auth.ListUserIDByAffiliationID(affiliationID)
	if err != nil {
		logging.Debug("Can't list users: %s", err.Error())
	}
	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, userList))
}

func (User) Create(c *gin.Context, currentUser *auth.User) {
	var form FormCreateUser = FormCreateUser{}
	err := c.ShouldBindBodyWith(&form, binding.JSON)
	if err != nil {
		logging.Error("Failed to bind JSON, err:", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}
	user := &auth.User{
		Email:         form.Email,
		PublicKey:     form.PublicKey,
		Role:          auth.Roles(auth.GLOBAL_EVALUATION_USER, auth.EXEMPT_MARKETING_CONTACT), // New user: in evaluation and not subbed to mailing list.
		AffiliationID: form.Affiliation,                                                       // Default: no affiliation
	}

	// Check: Affiliation User needs to be created by a user that:
	// - is an Affiliation Admin of the corresponding Affiliation
	// - or is a Global Admin
	if form.Affiliation != 0 {
		if currentUser == nil || ((!currentUser.Role.Includes(auth.AFFILIATION_ACCOUNT_ADMIN) || currentUser.AffiliationID != form.Affiliation) && !currentUser.Role.Includes(auth.GLOBAL_ADMIN)) {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
			return
		}
	}

	err = user.Create()
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	// TODO: Post-registration handling

	c.JSON(http.StatusOK, api.PayloadResponse(
		api.SUCCESS,
		gin.H{
			"id":          user.ID(),
			"email":       user.Email,
			"role":        user.Role,
			"affiliation": user.AffiliationID,
		},
	))
}

func (User) EmailExists(c *gin.Context) {
	emailAddr, err := mail.ParseAddress(c.Query("email"))
	if err != nil {
		logging.Debug("Failed to parse email address %s, err: %s", c.Query("email"), err.Error()) // won't log user error
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.RespInvalidEmail)
	}

	user := &auth.User{
		Email: emailAddr.String(),
	}

	exist, err := user.EmailExists()
	if err != nil {
		utils.HandleError(c, err)
	}

	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, exist))
}

// GLOBAL_ADMIN can update anyone
// AFFILIATION_ACCOUNT_ADMIN can update anyone in the same affiliation
// OTHERS can update themselves ONLY
// Note: Check the post form and only update field(s) that has been set
func (User) Update(c *gin.Context, currentUser *auth.User) {
	if currentUser == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespNotAuthorized)
		return
	}

	var form FormUpdateUser = FormUpdateUser{}
	err := c.ShouldBindBodyWith(&form, binding.JSON)
	if err != nil {
		logging.Error("Failed to bind JSON, err:", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.RespBadRequest)
		return
	}

	targetUser, err := auth.GetUserByID(form.ID)
	if err != nil {
		utils.HandleError(c, err) // other unhandled error
		return
	}

	if currentUser.Role.Includes(auth.GLOBAL_ADMIN) { // GLOBAL_ADMIN can update anyone
		// Update targetUser.
		// Everything is allowed to be updated.
		// *No validation*
		if form.Email != "" {
			targetUser.Email = form.Email
		}
		if form.Password != "" {
			targetUser.Password = utils.HashPassword(form.Password)
		}
		if form.Role != 0 { // A user needs at least 1 role. 0 (Roleless) is not allowed.
			targetUser.Role = auth.Role(form.Role)
		}
		if form.Affiliation != 0 {
			targetUser.AffiliationID = form.Affiliation
		}
	} else if currentUser.Role.Includes(auth.AFFILIATION_ACCOUNT_ADMIN) && currentUser.AffiliationID == targetUser.AffiliationID { // AFFILIATION_ACCOUNT_ADMIN can update anyone in the same affiliation
		if form.Email != "" {
			targetUser.Email = form.Email
		}
		if form.Password != "" {
			targetUser.Password = utils.HashPassword(form.Password)
		}
		if form.Role != 0 { // A user needs at least 1 role. 0 (Roleless) is not allowed.
			var tmpRole auth.Role = auth.Role(form.Role)
			// strip-off global roles from update request
			// global roles can't be manually updated
			tmpRole = tmpRole.RemoveRole(auth.GLOBAL_EVALUATION_USER)
			tmpRole = tmpRole.RemoveRole(auth.GLOBAL_PRODUCTION_USER)
			tmpRole = tmpRole.RemoveRole(auth.GLOBAL_INTERNAL_USER)
			tmpRole = tmpRole.RemoveRole(auth.GLOBAL_ADMIN)

			// preserve global roles of the targetUser
			if targetUser.Role.Includes(auth.GLOBAL_EVALUATION_USER) {
				tmpRole = tmpRole.AddRole(auth.GLOBAL_EVALUATION_USER)
			}
			if targetUser.Role.Includes(auth.GLOBAL_PRODUCTION_USER) {
				tmpRole = tmpRole.AddRole(auth.GLOBAL_PRODUCTION_USER)
			}
			if targetUser.Role.Includes(auth.GLOBAL_INTERNAL_USER) {
				tmpRole = tmpRole.AddRole(auth.GLOBAL_INTERNAL_USER)
			}
			if targetUser.Role.Includes(auth.GLOBAL_ADMIN) {
				tmpRole = tmpRole.AddRole(auth.GLOBAL_ADMIN)
			}
			targetUser.Role = tmpRole
		}
	} else if currentUser.ID() == targetUser.ID() { // OTHERS can update themselves ONLY
		// Can update only password
		if form.Password != "" {
			targetUser.Password = utils.HashPassword(form.Password)
		}
	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.RespAccessDenied)
		return
	}
	err = targetUser.Update()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, api.PayloadResponse(api.SUCCESS, gin.H{
		"id":          targetUser.ID(),
		"email":       targetUser.Email,
		"role":        targetUser.Role,
		"affiliation": targetUser.AffiliationID,
	}))
}

func (User) Wipe(c *gin.Context, _ *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}

func (User) CreateInfo(c *gin.Context, _ *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}

func (User) Info(c *gin.Context, _ *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}

func (User) UpdateInfo(c *gin.Context, _ *auth.User) {
	c.JSON(http.StatusNotImplemented, utils.RespNotImplemented)
}
