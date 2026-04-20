package order

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrOrderPaidImmutable = errors.New("paid orders are immutable")
	ErrUnavailableProduct = errors.New("unavailable product")
)

type ProductStatus string

const (
	ProductStatusPublished  ProductStatus = "published"
	ProductStatusHidden     ProductStatus = "hidden"
	ProductStatusArchived   ProductStatus = "archived"
	ProductStatusOutOfStock ProductStatus = "out_of_stock"
)

type Product struct {
	ID         string
	Name       string
	PriceCents int64
	Status     ProductStatus
	Stock      uint32
}

func (p Product) Available() bool {
	return p.Status == ProductStatusPublished && p.Stock > 0
}

type OrderStatus string

const (
	OrderStatusCreated OrderStatus = "created"
	OrderStatusPaid    OrderStatus = "paid"
)

type OrderItem struct {
	ProductID           string
	Quantity            uint32
	PriceSnapshotCents  int64
	ProductNameSnapshot string
}

type Order struct {
	ID     string
	UserID string
	Items  []OrderItem
	Status OrderStatus
}

type auditEvent struct {
	Name    string
	OrderID string
	At      time.Time
	Payload map[string]string
}

type idempotencyStore struct {
	mu      sync.Mutex
	results map[string]string
}

func newIdempotencyStore() *idempotencyStore {
	return &idempotencyStore{results: map[string]string{}}
}

func (s *idempotencyStore) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.results[key]
	return v, ok
}

func (s *idempotencyStore) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results[key] = value
}

type Service struct {
	mu          sync.Mutex
	orders      map[string]*Order
	idempotency *idempotencyStore
	auditLog    []auditEvent
	nextID      int
}

func NewService() *Service {
	return &Service{
		orders:      map[string]*Order{},
		idempotency: newIdempotencyStore(),
		auditLog:    make([]auditEvent, 0, 16),
		nextID:      1,
	}
}

func (s *Service) CreateOrder(userID, idempotencyKey string, products []Product) (*Order, error) {
	if reused, ok := s.idempotency.Get("create:" + idempotencyKey); ok {
		return s.orders[reused], nil
	}

	items := make([]OrderItem, 0, len(products))
	for _, product := range products {
		if !product.Available() {
			return nil, fmt.Errorf("%w: product_id=%s", ErrUnavailableProduct, product.ID)
		}
		items = append(items, OrderItem{
			ProductID:           product.ID,
			Quantity:            1,
			PriceSnapshotCents:  product.PriceCents,
			ProductNameSnapshot: product.Name,
		})
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	orderID := fmt.Sprintf("ord-%d", s.nextID)
	s.nextID++
	order := &Order{ID: orderID, UserID: userID, Items: items, Status: OrderStatusCreated}
	s.orders[orderID] = order
	s.idempotency.Put("create:"+idempotencyKey, orderID)
	s.appendAudit("order.created", orderID, map[string]string{"user_id": userID})
	return order, nil
}

func (s *Service) UpdateOrderItems(orderID string, items []OrderItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	order := s.orders[orderID]
	if order.Status == OrderStatusPaid {
		return ErrOrderPaidImmutable
	}
	order.Items = items
	s.appendAudit("order.updated", orderID, map[string]string{"items_updated": "true"})
	return nil
}

func (s *Service) ProcessPayment(orderID, idempotencyKey string) (OrderStatus, error) {
	if _, ok := s.idempotency.Get("payment:" + idempotencyKey); ok {
		return s.orders[orderID].Status, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	order := s.orders[orderID]
	order.Status = OrderStatusPaid
	s.idempotency.Put("payment:"+idempotencyKey, orderID)
	s.appendAudit("payment.processed", orderID, map[string]string{"status": string(order.Status)})
	return order.Status, nil
}

func (s *Service) PublishEvent(orderID, eventName, idempotencyKey string) bool {
	if _, ok := s.idempotency.Get("event:" + idempotencyKey); ok {
		return false
	}
	s.idempotency.Put("event:"+idempotencyKey, orderID)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.appendAudit(eventName, orderID, map[string]string{"published": "true"})
	return true
}

func (s *Service) AuditLog() []auditEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := make([]auditEvent, len(s.auditLog))
	copy(clone, s.auditLog)
	return clone
}

func (s *Service) appendAudit(name, orderID string, payload map[string]string) {
	s.auditLog = append(s.auditLog, auditEvent{Name: name, OrderID: orderID, Payload: payload, At: time.Now().UTC()})
}
