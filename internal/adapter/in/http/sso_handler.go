package http

import (
	"fmt"
	"net/http"
	"net/url"
	"phone-tokens/internal/service/users"
)

type SSOHandler struct {
	users users.Service
}

// func (h *SSOHandler) Login(w http.ResponseWriter, r *http.Request) {
//
//	   claims, ok := r.Context().Value(userContextKey).(*UserClaims)
//	   if !ok {
//	       http.Redirect(w, r, "/login?return_to="+url.QueryEscape(r.URL.String()), http.StatusFound)
//	       return
//	   }
//
//	   service := r.URL.Query().Get("service")
//	   agentID := r.URL.Query().Get("agent_id")
//	   redirectURI := r.URL.Query().Get("redirect_uri")
//
//	   token, err := h.service.Issue(r.Context(), tokens.IssueInput{
//	       UserID: claims.UserID,
//	       AgentId: agentID,
//	       Permissions: []string{"send_sms"},
//	       TTLSeconds: 3600,
//	   })
//	   if err != nil {
//	       http.Error(w, "internal error", 500)
//	       return
//	   }
//
//	   redirect := fmt.Sprintf("%s?token=%s", redirectURI, token.Token)
//
//	   http.Redirect(w, r, redirect, http.StatusFound)
//	}
func (h *SSOHandler) Authorize(w http.ResponseWriter, r *http.Request) {

	service := r.URL.Query().Get("service")
	agentID := r.URL.Query().Get("agent_id")
	redirectURI := r.URL.Query().Get("redirect_uri")

	agent, err := h.users.GetAgentByID(r.Context(), agentID)
	if err != nil {
		http.Error(w, "unknown agent", 400)
		return
	}

	//if agent.RedirectURI != redirectURI {
	//    http.Error(w, "invalid redirect uri", 400)
	//    return
	//}
	fmt.Printf("Authorizing %s to %s\n", service, redirectURI)
	fmt.Println(agent)

	//TODO: implement frontend
	frontendURL := fmt.Sprintf(
		"https://frontend.example.com/sso?client_id=%s&redirect_uri=%s",
		agentID,
		url.QueryEscape(redirectURI),
	)

	http.Redirect(w, r, frontendURL, http.StatusFound)
}
