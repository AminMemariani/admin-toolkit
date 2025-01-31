package product

import (
	"net/http"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/labstack/echo/v4"
	"github.com/sika365/admin-tools/context"
)

// Errors
var ()

type Rest interface {
	Find(ctx echo.Context) error
	ScanByImages(ctx echo.Context) error
	SyncByImages(ctx echo.Context) error
	SyncBySpreadSheets(ctx echo.Context) error
}

type rest struct {
	logic Logic
}

func newRest(h *simutils.HttpServer, logic Logic) (Rest, error) {
	r := &rest{
		logic: logic,
	}

	sg := h.
		// use prefix group
		PrefixGroup().
		// [prefixgroup path]/products
		Group("/products")
	{
		sg.GET("", r.Find)
		sg.POST("/scan", r.ScanByImages)
		sg.POST("/sync/images", r.SyncByImages)
		sg.POST("/sync/spreadsheets", r.SyncBySpreadSheets)
	}

	return r, nil
}

func (r *rest) Find(ectx echo.Context) error {
	return nil
}

func (r *rest) ScanByImages(ctx echo.Context) error {
	return nil
}

func (r *rest) SyncByImages(ectx echo.Context) error {
	var req SyncByImageRequest
	if ctx, ok := ectx.(*context.Context); !ok {
		return nil
	} else if err := ctx.Bind(&req); err != nil {
		return err
	} else if filters := ctx.QueryParams(); false {
		return nil
	} else if products, err := r.logic.SyncByImages(ctx, &req, filters); err != nil {
		return err
	} else {
		return simutils.ReplyTemplate(ctx, http.StatusOK, nil,
			&SyncResponse{Data: NewMapProducts(products...)},
			simutils.CreatePaginateTemplate(len(products), 0, len(products)),
		)
	}
}

func (r *rest) SyncBySpreadSheets(ectx echo.Context) error {
	var req SyncBySpreadSheetsRequest
	if ctx, ok := ectx.(*context.Context); !ok {
		return nil
	} else if err := ctx.Bind(&req); err != nil {
		return err
	} else if filters := ctx.QueryParams(); false {
		return nil
	} else if products, err := r.logic.SyncBySpreadSheets(ctx, &req, filters); err != nil {
		return err
	} else {
		return simutils.ReplyTemplate(ctx, http.StatusOK, nil,
			&SyncBySpreadSheetsResponse{Data: products},
			simutils.CreatePaginateTemplate(len(products.Nodes), 0, len(products.Nodes)),
		)
	}
}
