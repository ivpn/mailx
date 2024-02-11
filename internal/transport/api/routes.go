package api

func (h *Handler) SetupRoutes() {
	v1 := h.Server.Group("/v1")

	v1.Get("/recipient/:id", h.GetRecipient)
	v1.Get("/recipients/:userID", h.GetRecipients)
	v1.Post("/recipient", h.PostRecipient)
	v1.Put("/recipient", h.UpdateRecipient)
	v1.Delete("/recipient/:id", h.DeleteRecipient)
	v1.Get("/recipient/verify/:id/:verification", h.VerifyRecipient)
}
