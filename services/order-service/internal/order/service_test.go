package order

import (
	"errors"
	"testing"
)

func TestCreateOrderChecksAvailabilityAndSnapshots(t *testing.T) {
	svc := NewService()

	_, err := svc.CreateOrder("u1", "k1", []Product{{ID: "p1", Status: ProductStatusHidden, Stock: 1}})
	if !errors.Is(err, ErrUnavailableProduct) {
		t.Fatalf("expected unavailable product error, got %v", err)
	}

	created, err := svc.CreateOrder("u1", "k2", []Product{{ID: "p1", Name: "Phone", PriceCents: 1000, Status: ProductStatusPublished, Stock: 2}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.Items[0].PriceSnapshotCents != 1000 || created.Items[0].ProductNameSnapshot != "Phone" {
		t.Fatal("expected snapshots in order item")
	}
}

func TestOrderImmutableAfterPaid(t *testing.T) {
	svc := NewService()
	created, _ := svc.CreateOrder("u1", "create-1", []Product{{ID: "p1", Name: "Phone", PriceCents: 1000, Status: ProductStatusPublished, Stock: 2}})
	_, _ = svc.ProcessPayment(created.ID, "pay-1")

	err := svc.UpdateOrderItems(created.ID, []OrderItem{{ProductID: "p2", Quantity: 2}})
	if !errors.Is(err, ErrOrderPaidImmutable) {
		t.Fatalf("expected ErrOrderPaidImmutable, got %v", err)
	}
}

func TestIdempotencyKeys(t *testing.T) {
	svc := NewService()
	products := []Product{{ID: "p1", Name: "Phone", PriceCents: 1000, Status: ProductStatusPublished, Stock: 2}}

	first, _ := svc.CreateOrder("u1", "same-create", products)
	second, _ := svc.CreateOrder("u1", "same-create", products)
	if first.ID != second.ID {
		t.Fatal("create order idempotency failed")
	}

	status1, _ := svc.ProcessPayment(first.ID, "same-payment")
	status2, _ := svc.ProcessPayment(first.ID, "same-payment")
	if status1 != status2 {
		t.Fatal("payment idempotency failed")
	}

	if !svc.PublishEvent(first.ID, "order.paid", "event-1") {
		t.Fatal("first event publish should be true")
	}
	if svc.PublishEvent(first.ID, "order.paid", "event-1") {
		t.Fatal("duplicate event publish should be deduplicated")
	}
}

func TestAuditEventsRecorded(t *testing.T) {
	svc := NewService()
	created, _ := svc.CreateOrder("u1", "create-audit", []Product{{ID: "p1", Name: "Phone", PriceCents: 1000, Status: ProductStatusPublished, Stock: 2}})
	_, _ = svc.ProcessPayment(created.ID, "pay-audit")
	_ = svc.PublishEvent(created.ID, "order.paid", "event-audit")

	if len(svc.AuditLog()) < 3 {
		t.Fatalf("expected audit events to be recorded")
	}
}
