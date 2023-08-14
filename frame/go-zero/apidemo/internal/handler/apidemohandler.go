package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-lib/frame/go-zero/apidemo/internal/logic"
	"go-lib/frame/go-zero/apidemo/internal/svc"
	"go-lib/frame/go-zero/apidemo/internal/types"
)

func ApidemoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewApidemoLogic(r.Context(), svcCtx)
		resp, err := l.Apidemo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
