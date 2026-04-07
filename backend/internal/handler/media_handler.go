package handler

import (
    "strconv"

    "github.com/gin-gonic/gin"

    "mezian/internal/middleware"
    "mezian/internal/service"
)

// MediaHandler handles media management routes.
type MediaHandler struct {
    mediaSvc *service.MediaService
}

// NewMediaHandler creates a new MediaHandler.
func NewMediaHandler(mediaSvc *service.MediaService) *MediaHandler {
    return &MediaHandler{mediaSvc: mediaSvc}
}

// UploadImage godoc
// POST /ads/:id/media
// Uploads an image for an ad (multipart/form-data, field "image").
func (h *MediaHandler) UploadImage(c *gin.Context) {
    adIDStr := c.Param("id")
    adID, err := strconv.ParseUint(adIDStr, 10, 64)
    if err != nil {
        respondBadRequest(c, "invalid ad identifier")
        return
    }

    userID := middleware.GetUserID(c)

    file, header, err := c.Request.FormFile("image")
    if err != nil {
        respondBadRequest(c, "fichier image requis (champ: image)")
        return
    }
    defer file.Close()

    media, err := h.mediaSvc.UploadImage(uint(adID), userID, file, header)
    if err != nil {
        respondError(c, err)
        return
    }

    respondCreated(c, media)
}

// addYouTubeRequest est le body de POST /ads/:id/media/youtube.
type addYouTubeRequest struct {
    URL string `json:"url" binding:"required"`
}

// AddYouTube godoc
// POST /ads/:id/media/youtube
// Associates a YouTube video with an ad.
func (h *MediaHandler) AddYouTube(c *gin.Context) {
    adIDStr := c.Param("id")
    adID, err := strconv.ParseUint(adIDStr, 10, 64)
    if err != nil {
        respondBadRequest(c, "invalid ad identifier")
        return
    }

    userID := middleware.GetUserID(c)

    var req addYouTubeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondBadRequest(c, err.Error())
        return
    }

    media, err := h.mediaSvc.AddYouTube(uint(adID), userID, req.URL)
    if err != nil {
        respondError(c, err)
        return
    }

    respondCreated(c, media)
}

// DeleteMedia godoc
// DELETE /media/:id
// Supprime un media (fichier + entrée DB).
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
    mediaIDStr := c.Param("id")
    mediaID, err := strconv.ParseUint(mediaIDStr, 10, 64)
    if err != nil {
        respondBadRequest(c, "invalid media identifier")
        return
    }

    userID := middleware.GetUserID(c)
    userRole := middleware.GetUserRole(c)

    if err := h.mediaSvc.DeleteMedia(uint(mediaID), userID, userRole); err != nil {
        respondError(c, err)
        return
    }

    c.JSON(200, gin.H{"message": "media deleted"})
}

// SetCover godoc
// PUT /media/:id/cover
// Sets a media item as the ad cover image.
func (h *MediaHandler) SetCover(c *gin.Context) {
    mediaIDStr := c.Param("id")
    mediaID, err := strconv.ParseUint(mediaIDStr, 10, 64)
    if err != nil {
        respondBadRequest(c, "invalid media identifier")
        return
    }

    userID := middleware.GetUserID(c)

    if err := h.mediaSvc.SetCover(uint(mediaID), userID); err != nil {
        respondError(c, err)
        return
    }

    c.JSON(200, gin.H{"message": "cover image updated"})
}

// updateOrderRequest est le body de PUT /media/:id/order.
type updateOrderRequest struct {
    SortOrder int `json:"sort_order"`
}

// UpdateOrder godoc
// PUT /media/:id/order
// Updates the sort_order of a media item.
func (h *MediaHandler) UpdateOrder(c *gin.Context) {
    mediaIDStr := c.Param("id")
    mediaID, err := strconv.ParseUint(mediaIDStr, 10, 64)
    if err != nil {
        respondBadRequest(c, "invalid media identifier")
        return
    }

    userID := middleware.GetUserID(c)

    var req updateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondBadRequest(c, err.Error())
        return
    }

    if err := h.mediaSvc.UpdateOrder(uint(mediaID), userID, req.SortOrder); err != nil {
        respondError(c, err)
        return
    }

    c.JSON(200, gin.H{"message": "order updated"})
}
