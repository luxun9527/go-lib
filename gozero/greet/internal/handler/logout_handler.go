package handler

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"greet/internal/logic"
	"greet/internal/response" // ①
	"greet/internal/svc"
	"greet/internal/types"
	"net/http"
)

func logoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(w, nil, err)
			return
		}

		l := logic.NewLogoutLogic(r.Context(), svcCtx)
		resp, err := l.Logout(&req)
		response.Response(w, resp, err) //②

	}
}
