package client

import (
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/context"
)

func (c *Client) GetCategoryByAlias(ctx *context.Context, alias string) (category *models.Category, err error) {
	var categoryResp models.CategoriesResponse
	if resp, err := c.R().
		SetPathParams(map[string]string{
			"alias": alias,
		}).
		SetQueryParamsFromValues(url.Values{
			"limit":    []string{cast.ToString(1)},
			"includes": []string{"Nodes"},
			"excludes": []string{"product_nodes", "current_node"},
		}).
		SetResult(&categoryResp).
		SetError(&categoryResp).
		Get("/categories/{alias}"); err != nil {
		logrus.Info(err)
		return nil, err
	} else if !resp.IsSuccess() {
		return nil, fmt.Errorf(resp.Status())
	} else if categories := categoryResp.Data.Categories; len(categories) == 0 {
		return nil, models.ErrNotFound
	} else {
		return categories[0], nil
	}
}

func (c *Client) StoreCategory(ctx *context.Context, category *models.Category, in ...*models.Node) (*models.Category, error) {
	var (
		categoryResp models.CategoriesResponse
		nodeIDs      = make(database.PIDs, 0, len(in))
	)
	for _, n := range in {
		nodeIDs = append(nodeIDs, n.ID)
	}
	if resp, err := c.R().
		SetBody(&models.CategoryRequest{
			AddedNodes: nodeIDs,
			Category:   *category,
		}).
		SetResult(&categoryResp).
		SetError(&categoryResp).
		Post("/categories"); err != nil {
		logrus.Info(err)
		return nil, err
	} else if !resp.IsSuccess() {
		return nil, fmt.Errorf("create categoriy failed: %v", resp.Status())
	} else if categories := categoryResp.Data.Categories; len(categories) == 0 {
		return nil, models.ErrNotFound
	} else {
		return categories[0], nil
	}
}
