package handler

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"greet/internal/logic"
	"greet/internal/response" // ①
	"greet/internal/svc"
	"greet/internal/types"
	"net/http"
)

func loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		response.Response(w, resp, err) //②

	}
}
