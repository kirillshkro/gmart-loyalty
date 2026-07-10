package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
	"github.com/kirillshkro/gmart-loyalty/pkg/utils"
)

// Проверяем в том ли порядке выдаются заказы
func (s *TestOrderSuite) Test_SortOrders() {
	userID, err := s.service.userFromCookie(s.uCookie)
	if err != nil {
		s.T().Error(err)
	}
	ctx := context.WithValue(context.TODO(), types.UserID, userID)
	//Создаем несколько заказов
	for range 5 {
		numOrder := utils.GenNumberOrder(8)
		reqOrder := httptest.NewRequestWithContext(ctx, http.MethodPost, "/api/user/orders", bytes.NewBufferString(numOrder))
		reqOrder.AddCookie(s.uCookie)
		rr := httptest.NewRecorder()
		s.service.SetOrderUser(rr, reqOrder)
	}
	//получить список заказов
	reqGetOrders := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/user/orders", nil)
	reqGetOrders.AddCookie(s.uCookie)
	rr := httptest.NewRecorder()
	s.service.OrdersByUser(rr, reqGetOrders)
	//проверить что заказы получены
	s.resp = rr.Result()
	var orders []model.Order
	if err = json.NewDecoder(s.resp.Body).Decode(&orders); err != nil {
		s.T().Errorf("Error decoding response body: %v", err)
		return
	}
	s.Assert().Positive(len(orders))

	//проверить что заказы в правильном порядке
	for i := range len(orders) - 1 {
		//сравнить каждую пару заказов на убывание даты создания
		if !orders[i].UploadedAt.After(orders[i+1].UploadedAt) {
			s.T().Errorf("Orders are not sorted by creation date")
		}
	}

	s.Assert().Equal(http.StatusOK, s.resp.StatusCode)
}

func (s *TestOrderSuite) Test_UnautorizedGetOrder() {
	userID, err := s.service.userFromCookie(s.uCookie)
	if err != nil {
		s.T().Error(err)
	}
	//Создаем тестовый заказ
	ctx := context.WithValue(context.TODO(), types.UserID, userID)
	numOrder := utils.GenNumberOrder(8)
	reqOrder := httptest.NewRequestWithContext(ctx, http.MethodPost, "/api/user/orders", bytes.NewBufferString(numOrder))
	reqOrder.AddCookie(s.uCookie)
	rr := httptest.NewRecorder()
	s.service.SetOrderUser(rr, reqOrder)
	//получить список заказов
	reqGetOrders := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/user/orders", nil)
	rr = httptest.NewRecorder()
	s.service.OrdersByUser(rr, reqGetOrders)

	s.resp = rr.Result()

	s.Assert().Equal(http.StatusUnauthorized, s.resp.StatusCode)
}

func (s *TestOrderSuite) Test_EmptyOrdersList() {
	//создаем юзера
	user := model.User{
		Login:    "emptyorderuser",
		Password: "orders",
	}

	body, _ := json.Marshal(user)
	uReq := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	s.service.Register(rr, uReq)
	resp := rr.Result()
	defer resp.Body.Close()
	uc := resp.Cookies()[0]
	userID, err := s.service.userFromCookie(uc)
	if err != nil {
		s.T().Error(err)
	}
	//теперь у юзера еще нет заказов, список пуст
	ctx := context.WithValue(context.TODO(), types.UserID, userID)
	//получить список заказов
	reqGetOrders := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/user/orders", nil)
	reqGetOrders.AddCookie(uc)
	rr = httptest.NewRecorder()
	s.service.OrdersByUser(rr, reqGetOrders)
	resp = rr.Result()
	s.Assert().Equal(http.StatusNoContent, resp.StatusCode)
}
