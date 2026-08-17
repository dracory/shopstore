package shopstore

import (
	"strings"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestNewDiscountDefaults(t *testing.T) {
	discount := NewDiscount()
	if discount == nil {
		t.Fatal("NewDiscount returned nil")
	}

	if discount.GetStatus() != DISCOUNT_STATUS_DRAFT {
		t.Fatalf("expected status %q, got %q", DISCOUNT_STATUS_DRAFT, discount.GetStatus())
	}

	if discount.GetType() != DISCOUNT_TYPE_PERCENT {
		t.Fatalf("expected type %q, got %q", DISCOUNT_TYPE_PERCENT, discount.GetType())
	}

	if discount.GetTitle() != "" {
		t.Fatalf("expected empty title, got %q", discount.GetTitle())
	}

	if discount.GetDescription() != "" {
		t.Fatalf("expected empty description, got %q", discount.GetDescription())
	}

	if discount.GetMemo() != "" {
		t.Fatalf("expected empty memo, got %q", discount.GetMemo())
	}

	if discount.GetAmount() != 0 {
		t.Fatalf("expected amount 0, got %f", discount.GetAmount())
	}

	if discount.GetCode() == "" {
		t.Fatal("expected generated code to be non-empty")
	}

	if discount.GetStartsAt() != NULL_DATETIME {
		t.Fatalf("expected starts at %q, got %q", NULL_DATETIME, discount.GetStartsAt())
	}

	if discount.GetEndsAt() != NULL_DATETIME {
		t.Fatalf("expected ends at %q, got %q", NULL_DATETIME, discount.GetEndsAt())
	}

	if discount.GetSoftDeletedAt() != MAX_DATETIME {
		t.Fatalf("expected soft deleted at %q, got %q", MAX_DATETIME, discount.GetSoftDeletedAt())
	}

	if discount.GetID() == "" {
		t.Fatal("expected generated ID to be non-empty")
	}

	if discount.GetCreatedAt() == "" {
		t.Fatal("expected created at to be set")
	}

	if discount.GetUpdatedAt() == "" {
		t.Fatal("expected updated at to be set")
	}

	metas, err := discount.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected no metas by default, got %v", metas)
	}

	if discount.GetMeta("missing") != "" {
		t.Fatal("expected GetMeta for missing key to return empty string")
	}
}

func TestNewDiscountCodeUsesCrockfordAlphabet(t *testing.T) {
	discount := NewDiscount()
	code := strings.ToUpper(discount.GetCode())

	const allowed = "BCDFGHJKLMNPQRSTVWXYZ23456789"
	for _, r := range code {
		if !strings.ContainsRune(allowed, r) {
			t.Fatalf("generated code %q contains disallowed character %q", code, r)
		}
	}
}

func TestDiscountDataTracking(t *testing.T) {
	discount := &Discount{}

	discount.SetTitle("Title").
		SetDescription("Desc").
		SetMemo("Memo").
		SetCode("CODE").
		SetStatus(DISCOUNT_STATUS_ACTIVE).
		SetType(DISCOUNT_TYPE_AMOUNT)

	if _, ok := discount.SetAmount(12.34).(*Discount); !ok {
		t.Fatal("expected SetAmount to return *Discount")
	}

	data := discount.Data()
	if data == nil {
		t.Fatal("expected Data to return initialized map")
	}

	expected := map[string]string{
		COLUMN_TITLE:       "Title",
		COLUMN_DESCRIPTION: "Desc",
		COLUMN_MEMO:        "Memo",
		COLUMN_CODE:        "CODE",
		COLUMN_STATUS:      DISCOUNT_STATUS_ACTIVE,
		COLUMN_TYPE:        DISCOUNT_TYPE_AMOUNT,
		COLUMN_AMOUNT:      "12.34",
	}

	for key, want := range expected {
		if got := data[key]; got != want {
			t.Fatalf("expected %s to be %q, got %q", key, want, got)
		}
	}

	changed := discount.DataChanged()
	for key, want := range expected {
		if got := changed[key]; got != want {
			t.Fatalf("expected DataChanged to track %s as %q, got %q", key, want, got)
		}
	}

	discount.MarkAsNotDirty()
	if len(discount.DataChanged()) != 0 {
		t.Fatal("expected DataChanged to be empty after MarkAsNotDirty")
	}

	discount.SetTitle("Updated")
	if discount.DataChanged()[COLUMN_TITLE] != "Updated" {
		t.Fatalf("expected DataChanged to track updated title, got %q", discount.DataChanged()[COLUMN_TITLE])
	}
}

func TestDiscountCarbonHelpers(t *testing.T) {
	discount := &Discount{}

	createdAt := "2025-01-01 00:00:00"
	updatedAt := "2025-02-02 12:34:56"
	startsAt := "2025-03-03 09:00:00"
	endsAt := "2025-04-04 18:30:00"
	softDeletedAt := "2026-05-05 00:00:00"

	if _, ok := discount.SetCreatedAt(createdAt).(*Discount); !ok {
		t.Fatal("expected SetCreatedAt to return *Discount")
	}
	if discount.GetCreatedAt() != createdAt {
		t.Fatalf("expected CreatedAt to be %q, got %q", createdAt, discount.GetCreatedAt())
	}
	if discount.GetCreatedAtCarbon().ToDateTimeString() != createdAt {
		t.Fatalf("expected CreatedAtCarbon to match input, got %q", discount.GetCreatedAtCarbon().ToDateTimeString())
	}

	if _, ok := discount.SetUpdatedAt(updatedAt).(*Discount); !ok {
		t.Fatal("expected SetUpdatedAt to return *Discount")
	}
	if discount.GetUpdatedAtCarbon().ToDateTimeString() != updatedAt {
		t.Fatalf("expected UpdatedAtCarbon to match input, got %q", discount.GetUpdatedAtCarbon().ToDateTimeString())
	}

	if _, ok := discount.SetStartsAt(startsAt).(*Discount); !ok {
		t.Fatal("expected SetStartsAt to return *Discount")
	}
	if discount.GetStartsAtCarbon().ToDateTimeString(carbon.UTC) != startsAt {
		t.Fatalf("expected StartsAtCarbon to match input, got %q", discount.GetStartsAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := discount.SetEndsAt(endsAt).(*Discount); !ok {
		t.Fatal("expected SetEndsAt to return *Discount")
	}
	if discount.GetEndsAtCarbon().ToDateTimeString(carbon.UTC) != endsAt {
		t.Fatalf("expected EndsAtCarbon to match input, got %q", discount.GetEndsAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := discount.SetSoftDeletedAt(softDeletedAt).(*Discount); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Discount")
	}
	if discount.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != softDeletedAt {
		t.Fatalf("expected SoftDeletedAtCarbon to match input, got %q", discount.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}

func TestDiscountMetasRoundTrip(t *testing.T) {
	discount := &Discount{}

	if err := discount.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err := discount.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if got := metas["alpha"]; got != "beta" {
		t.Fatalf("expected meta to be %q, got %q", "beta", got)
	}

	if got := discount.GetMeta("alpha"); got != "beta" {
		t.Fatalf("expected GetMeta helper to return %q, got %q", "beta", got)
	}
}

func TestDiscountMetasUpsertMergesValues(t *testing.T) {
	discount := &Discount{}

	if err := discount.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting initial metas: %v", err)
	}

	if err := discount.MetasUpsert(map[string]string{"alpha": "updated", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	metas, err := discount.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if metas["alpha"] != "updated" {
		t.Fatalf("expected alpha meta to be updated, got %q", metas["alpha"])
	}

	if metas["gamma"] != "delta" {
		t.Fatalf("expected gamma meta to be present, got %q", metas["gamma"])
	}
}

func TestDiscountMetaRemove(t *testing.T) {
	discount := &Discount{}

	if err := discount.SetMetas(map[string]string{"alpha": "beta", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := discount.MetaRemove("alpha"); err != nil {
		t.Fatalf("unexpected error removing meta: %v", err)
	}

	if discount.GetMeta("alpha") != "" {
		t.Fatal("expected removed meta to return empty string")
	}

	metas, err := discount.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if _, exists := metas["alpha"]; exists {
		t.Fatal("expected alpha meta to be removed from stored metas")
	}

	if metas["gamma"] != "delta" {
		t.Fatalf("expected gamma meta to remain, got %q", metas["gamma"])
	}
}

func TestDiscountMetasRemoveList(t *testing.T) {
	discount := &Discount{}

	if err := discount.SetMetas(map[string]string{"alpha": "beta", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := discount.MetasRemove([]string{"alpha", "gamma"}); err != nil {
		t.Fatalf("unexpected error removing metas: %v", err)
	}

	metas, err := discount.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected all metas to be removed, got %v", metas)
	}
}

func TestDiscountMetasInvalidJSON(t *testing.T) {
	discount := NewDiscountFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if _, err := discount.GetMetas(); err == nil {
		t.Fatal("expected error when parsing invalid metas JSON")
	}

	if got := discount.GetMeta("anything"); got != "" {
		t.Fatalf("expected GetMeta to return empty string on invalid JSON, got %q", got)
	}
}

func TestDiscountMetasHandlesNullJSON(t *testing.T) {
	discount := NewDiscountFromExistingData(map[string]string{
		COLUMN_METAS: "null",
	})

	metas, err := discount.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected empty metas map for null JSON, got %v", metas)
	}

	if err := discount.MetasUpsert(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	if got := discount.GetMeta("alpha"); got != "beta" {
		t.Fatalf("expected GetMeta helper to return %q, got %q", "beta", got)
	}
}

func TestDiscountSetMetaConvenience(t *testing.T) {
	discount := &Discount{}

	if err := discount.SetMeta("key", "value"); err != nil {
		t.Fatalf("unexpected error from SetMeta: %v", err)
	}

	if got := discount.GetMeta("key"); got != "value" {
		t.Fatalf("expected GetMeta to return %q, got %q", "value", got)
	}
}

func TestDiscountSetterChainingAndGetters(t *testing.T) {
	discount := &Discount{}

	if _, ok := discount.SetDescription("desc").(*Discount); !ok {
		t.Fatal("expected SetDescription to return *Discount")
	}
	if discount.GetDescription() != "desc" {
		t.Fatalf("expected GetDescription getter to return %q, got %q", "desc", discount.GetDescription())
	}

	if _, ok := discount.SetMemo("memo").(*Discount); !ok {
		t.Fatal("expected SetMemo to return *Discount")
	}
	if discount.GetMemo() != "memo" {
		t.Fatalf("expected GetMemo getter to return %q, got %q", "memo", discount.GetMemo())
	}

	if _, ok := discount.SetTitle("title").(*Discount); !ok {
		t.Fatal("expected SetTitle to return *Discount")
	}
	if discount.GetTitle() != "title" {
		t.Fatalf("expected GetTitle getter to return %q, got %q", "title", discount.GetTitle())
	}

	if _, ok := discount.SetStatus(DISCOUNT_STATUS_INACTIVE).(*Discount); !ok {
		t.Fatal("expected SetStatus to return *Discount")
	}
	if discount.GetStatus() != DISCOUNT_STATUS_INACTIVE {
		t.Fatalf("expected GetStatus getter to return %q, got %q", DISCOUNT_STATUS_INACTIVE, discount.GetStatus())
	}

	if _, ok := discount.SetType(DISCOUNT_TYPE_AMOUNT).(*Discount); !ok {
		t.Fatal("expected SetType to return *Discount")
	}
	if discount.GetType() != DISCOUNT_TYPE_AMOUNT {
		t.Fatalf("expected GetType getter to return %q, got %q", DISCOUNT_TYPE_AMOUNT, discount.GetType())
	}

	start := "2024-01-01 00:00:00"
	end := "2024-12-31 23:59:59"

	if _, ok := discount.SetStartsAt(start).(*Discount); !ok {
		t.Fatal("expected SetStartsAt to return *Discount")
	}
	if discount.GetStartsAt() != start {
		t.Fatalf("expected GetStartsAt getter to return %q, got %q", start, discount.GetStartsAt())
	}

	if _, ok := discount.SetEndsAt(end).(*Discount); !ok {
		t.Fatal("expected SetEndsAt to return *Discount")
	}
	if discount.GetEndsAt() != end {
		t.Fatalf("expected GetEndsAt getter to return %q, got %q", end, discount.GetEndsAt())
	}
}

func TestDiscountStatusPredicates(t *testing.T) {
	discount := &Discount{}

	if _, ok := discount.SetStatus(DISCOUNT_STATUS_ACTIVE).(*Discount); !ok {
		t.Fatal("expected SetStatus to return *Discount")
	}

	if !discount.IsActive() {
		t.Fatal("expected discount to be active")
	}

	if discount.IsDraft() {
		t.Fatal("expected discount not to be draft when active")
	}

	if discount.IsInactive() {
		t.Fatal("expected discount not to be inactive when active")
	}

	if _, ok := discount.SetStatus(DISCOUNT_STATUS_DRAFT).(*Discount); !ok {
		t.Fatal("expected SetStatus to return *Discount")
	}

	if !discount.IsDraft() {
		t.Fatal("expected discount to be draft")
	}

	if _, ok := discount.SetStatus(DISCOUNT_STATUS_INACTIVE).(*Discount); !ok {
		t.Fatal("expected SetStatus to return *Discount")
	}

	if !discount.IsInactive() {
		t.Fatal("expected discount to be inactive")
	}
}

func TestDiscountTemporalPredicates(t *testing.T) {
	discount := &Discount{}

	// No dates set - should not be started or ended
	if discount.IsStarted() {
		t.Fatal("expected discount with no start date not to be started")
	}

	if discount.IsEnded() {
		t.Fatal("expected discount with no end date not to be ended")
	}

	// Set start date in past
	past := carbon.Now(carbon.UTC).SubDay().ToDateTimeString(carbon.UTC)
	discount.SetStartsAt(past)

	if !discount.IsStarted() {
		t.Fatal("expected discount with past start date to be started")
	}

	// Set start date in future
	future := carbon.Now(carbon.UTC).AddDay().ToDateTimeString(carbon.UTC)
	discount.SetStartsAt(future)

	if discount.IsStarted() {
		t.Fatal("expected discount with future start date not to be started")
	}

	// Set end date in past
	discount.SetEndsAt(past)

	if !discount.IsEnded() {
		t.Fatal("expected discount with past end date to be ended")
	}

	// IsExpired should be alias for IsEnded
	if !discount.IsExpired() {
		t.Fatal("expected IsExpired to match IsEnded")
	}

	// Set end date in future
	discount.SetEndsAt(future)

	if discount.IsEnded() {
		t.Fatal("expected discount with future end date not to be ended")
	}
}

func TestDiscountIsValidNow(t *testing.T) {
	now := carbon.Now(carbon.UTC)
	past := now.SubDay().ToDateTimeString(carbon.UTC)
	future := now.AddDay().ToDateTimeString(carbon.UTC)

	// Not active
	d := NewDiscount().SetStatus(DISCOUNT_STATUS_DRAFT)
	if d.IsValidNow() {
		t.Fatal("expected non-active discount to not be valid")
	}

	// Active but not started
	d.SetStatus(DISCOUNT_STATUS_ACTIVE).SetStartsAt(future).SetEndsAt(future)
	if d.IsValidNow() {
		t.Fatal("expected future discount to not be valid")
	}

	// Active, started, not ended
	d.SetStartsAt(past).SetEndsAt(future)
	if !d.IsValidNow() {
		t.Fatal("expected valid discount to be valid now")
	}

	// Active, started, but ended
	d.SetEndsAt(past)
	if d.IsValidNow() {
		t.Fatal("expected ended discount to not be valid")
	}

	// Active, started, not ended, but max_uses exhausted
	d.SetEndsAt(future).SetMaxUses(1).SetMaxUsesCount(1)
	if d.IsValidNow() {
		t.Fatal("expected discount with exhausted max_uses to not be valid")
	}

	// Same discount, reset count — valid again
	d.SetMaxUsesCount(0)
	if !d.IsValidNow() {
		t.Fatal("expected discount with max_uses not exhausted to be valid")
	}
}

func TestDiscountNullDateTimeSentinelPredicates(t *testing.T) {
	discount := NewDiscount()

	if discount.GetStartsAt() != NULL_DATETIME {
		t.Fatalf("expected default starts_at to be %q, got %q", NULL_DATETIME, discount.GetStartsAt())
	}

	if discount.GetEndsAt() != NULL_DATETIME {
		t.Fatalf("expected default ends_at to be %q, got %q", NULL_DATETIME, discount.GetEndsAt())
	}

	if discount.IsStarted() {
		t.Fatal("expected default discount with NULL_DATETIME start not to be started")
	}

	if discount.IsEnded() {
		t.Fatal("expected default discount with NULL_DATETIME end not to be ended")
	}

	if discount.IsExpired() {
		t.Fatal("expected default discount with NULL_DATETIME end not to be expired")
	}

	if discount.IsValidNow() {
		t.Fatal("expected default draft discount with null sentinel dates not to be valid now")
	}

	// Active status with null sentinel dates should still not be valid now
	discount.SetStatus(DISCOUNT_STATUS_ACTIVE)
	if !discount.IsActive() {
		t.Fatal("expected discount to be active after status change")
	}

	if discount.IsValidNow() {
		t.Fatal("expected active discount with null sentinel dates not to be valid now")
	}

	if discount.IsStarted() || discount.IsEnded() {
		t.Fatal("expected null sentinel dates to mean neither started nor ended")
	}
}

func TestNewDiscountMaxUsesDefaults(t *testing.T) {
	discount := NewDiscount()

	if discount.GetMaxUses() != DEFAULT_MAX_USES {
		t.Fatalf("expected default max_uses to be %d (DEFAULT_MAX_USES), got %d", DEFAULT_MAX_USES, discount.GetMaxUses())
	}

	if discount.GetMaxUsesCount() != 0 {
		t.Fatalf("expected default max_uses_count to be 0, got %d", discount.GetMaxUsesCount())
	}

	if discount.GetMaxUsesPerCustomer() != DEFAULT_MAX_USES_PER_CUSTOMER {
		t.Fatalf("expected default max_uses_per_customer to be %d (DEFAULT_MAX_USES_PER_CUSTOMER), got %d", DEFAULT_MAX_USES_PER_CUSTOMER, discount.GetMaxUsesPerCustomer())
	}

	// Per-customer count map starts empty in metas (no key set)
	if v := discount.GetMeta(META_MAX_USES_PER_CUSTOMER_COUNT); v != "" {
		t.Fatalf("expected meta max_uses_per_customer_count to be empty, got %q", v)
	}
}

func TestDiscountMaxUsesSettersAndGetters(t *testing.T) {
	discount := &Discount{}

	if _, ok := discount.SetMaxUses(50).(*Discount); !ok {
		t.Fatal("expected SetMaxUses to return *Discount")
	}
	if discount.GetMaxUses() != 50 {
		t.Fatalf("expected max_uses to be 50, got %d", discount.GetMaxUses())
	}

	if _, ok := discount.SetMaxUsesCount(3).(*Discount); !ok {
		t.Fatal("expected SetMaxUsesCount to return *Discount")
	}
	if discount.GetMaxUsesCount() != 3 {
		t.Fatalf("expected max_uses_count to be 3, got %d", discount.GetMaxUsesCount())
	}

	if _, ok := discount.SetMaxUsesPerCustomer(2).(*Discount); !ok {
		t.Fatal("expected SetMaxUsesPerCustomer to return *Discount")
	}
	if discount.GetMaxUsesPerCustomer() != 2 {
		t.Fatalf("expected max_uses_per_customer to be 2, got %d", discount.GetMaxUsesPerCustomer())
	}
}

func TestDiscountMaxUsesDataTracking(t *testing.T) {
	discount := &Discount{}

	discount.SetMaxUses(100).
		SetMaxUsesCount(42).
		SetMaxUsesPerCustomer(5)

	data := discount.Data()

	expected := map[string]string{
		COLUMN_MAX_USES:              "100",
		COLUMN_MAX_USES_COUNT:        "42",
		COLUMN_MAX_USES_PER_CUSTOMER: "5",
	}

	for key, want := range expected {
		if got := data[key]; got != want {
			t.Fatalf("expected %s to be %q, got %q", key, want, got)
		}
	}

	changed := discount.DataChanged()
	for key, want := range expected {
		if got := changed[key]; got != want {
			t.Fatalf("expected DataChanged to track %s as %q, got %q", key, want, got)
		}
	}

	discount.MarkAsNotDirty()
	if len(discount.DataChanged()) != 0 {
		t.Fatal("expected DataChanged to be empty after MarkAsNotDirty")
	}

	discount.SetMaxUses(200)
	if discount.DataChanged()[COLUMN_MAX_USES] != "200" {
		t.Fatalf("expected DataChanged to track updated max_uses, got %q", discount.DataChanged()[COLUMN_MAX_USES])
	}
}

func TestDiscountIsMaxUsesReached(t *testing.T) {
	discount := &Discount{}

	// Default cap (DEFAULT_MAX_USES), count well below — not reached
	discount.SetMaxUses(DEFAULT_MAX_USES).SetMaxUsesCount(10)
	if discount.IsMaxUsesReached() {
		t.Fatal("expected discount with count well below DEFAULT_MAX_USES to not have reached limit")
	}

	// Default cap, count just below — not reached
	discount.SetMaxUses(DEFAULT_MAX_USES).SetMaxUsesCount(DEFAULT_MAX_USES - 1)
	if discount.IsMaxUsesReached() {
		t.Fatal("expected discount with count one below limit to not have reached limit")
	}

	// Default cap, count equals — reached
	discount.SetMaxUses(DEFAULT_MAX_USES).SetMaxUsesCount(DEFAULT_MAX_USES)
	if !discount.IsMaxUsesReached() {
		t.Fatal("expected discount with count equal to limit to have reached limit")
	}

	// Default cap, count exceeds — reached
	discount.SetMaxUses(DEFAULT_MAX_USES).SetMaxUsesCount(DEFAULT_MAX_USES + 1)
	if !discount.IsMaxUsesReached() {
		t.Fatal("expected discount with count exceeding limit to have reached limit")
	}

	// Finite cap (50), count below — not reached
	discount.SetMaxUses(50).SetMaxUsesCount(10)
	if discount.IsMaxUsesReached() {
		t.Fatal("expected discount with count below finite limit to not have reached limit")
	}

	// Finite cap (50), count equals — reached
	discount.SetMaxUses(50).SetMaxUsesCount(50)
	if !discount.IsMaxUsesReached() {
		t.Fatal("expected discount with count equal to finite limit to have reached limit")
	}

	// Single-use (max_uses=1), not yet redeemed — not reached
	discount.SetMaxUses(1).SetMaxUsesCount(0)
	if discount.IsMaxUsesReached() {
		t.Fatal("expected single-use discount with 0 redemptions to not have reached limit")
	}

	// Single-use (max_uses=1), redeemed once — reached
	discount.SetMaxUses(1).SetMaxUsesCount(1)
	if !discount.IsMaxUsesReached() {
		t.Fatal("expected single-use discount with 1 redemption to have reached limit")
	}
}

func TestDiscountIsMaxUsesPerCustomerReached(t *testing.T) {
	discount := &Discount{}
	discount.SetMetas(map[string]string{})

	// Default per-customer cap, customer with 0 uses — not reached
	discount.SetMaxUsesPerCustomer(DEFAULT_MAX_USES_PER_CUSTOMER)
	discount.IncrementMaxUsesPerCustomer("cust-1")
	// cust-1 now has 1 use, well below DEFAULT_MAX_USES_PER_CUSTOMER
	if discount.IsMaxUsesPerCustomerReached("cust-1") {
		t.Fatal("expected customer with 1 use to not have reached default per-customer limit")
	}

	// Finite per-customer cap (1), customer with 0 uses — not reached
	discount.SetMaxUsesPerCustomer(1)
	if discount.IsMaxUsesPerCustomerReached("cust-2") {
		t.Fatal("expected customer with 0 uses to not have reached per-customer limit of 1")
	}

	// Finite per-customer cap (1), increment to 1 — reached
	discount.IncrementMaxUsesPerCustomer("cust-2")
	if !discount.IsMaxUsesPerCustomerReached("cust-2") {
		t.Fatal("expected customer with 1 use to have reached per-customer limit of 1")
	}

	// N-per-customer (5), increment 4 times — not reached
	discount.SetMaxUsesPerCustomer(5)
	for i := 0; i < 4; i++ {
		discount.IncrementMaxUsesPerCustomer("cust-3")
	}
	if discount.IsMaxUsesPerCustomerReached("cust-3") {
		t.Fatal("expected customer with 4 uses to not have reached per-customer limit of 5")
	}

	// Increment to 5 — reached
	discount.IncrementMaxUsesPerCustomer("cust-3")
	if !discount.IsMaxUsesPerCustomerReached("cust-3") {
		t.Fatal("expected customer with 5 uses to have reached per-customer limit of 5")
	}

	// Different customers are independent
	if discount.IsMaxUsesPerCustomerReached("cust-4") {
		t.Fatal("expected new customer to not have reached limit")
	}

	// Empty customerID — not reached (safe default)
	if discount.IsMaxUsesPerCustomerReached("") {
		t.Fatal("expected empty customerID to not have reached limit")
	}
}

func TestDiscountGetMaxUsesPerCustomerLeft(t *testing.T) {
	discount := &Discount{}
	discount.SetMetas(map[string]string{})

	// Cap=5, 0 uses — 5 left
	discount.SetMaxUsesPerCustomer(5)
	if left := discount.GetMaxUsesPerCustomerLeft("cust-1"); left != 5 {
		t.Fatalf("expected 5 left, got %d", left)
	}

	// Cap=5, 2 uses — 3 left
	discount.IncrementMaxUsesPerCustomer("cust-1")
	discount.IncrementMaxUsesPerCustomer("cust-1")
	if left := discount.GetMaxUsesPerCustomerLeft("cust-1"); left != 3 {
		t.Fatalf("expected 3 left, got %d", left)
	}

	// Cap=5, 5 uses — 0 left
	for i := 0; i < 3; i++ {
		discount.IncrementMaxUsesPerCustomer("cust-1")
	}
	if left := discount.GetMaxUsesPerCustomerLeft("cust-1"); left != 0 {
		t.Fatalf("expected 0 left, got %d", left)
	}

	// Cap=5, 7 uses (over) — 0 left, not negative
	discount.IncrementMaxUsesPerCustomer("cust-1")
	discount.IncrementMaxUsesPerCustomer("cust-1")
	if left := discount.GetMaxUsesPerCustomerLeft("cust-1"); left != 0 {
		t.Fatalf("expected 0 left (clamped), got %d", left)
	}

	// Empty customerID — returns full cap
	if left := discount.GetMaxUsesPerCustomerLeft(""); left != 5 {
		t.Fatalf("expected full cap for empty customerID, got %d", left)
	}
}

func TestDiscountIncrementMaxUsesPerCustomer(t *testing.T) {
	discount := &Discount{}
	discount.SetMetas(map[string]string{})

	// Increment returns *Discount (chainable)
	if _, ok := discount.IncrementMaxUsesPerCustomer("cust-1").(*Discount); !ok {
		t.Fatal("expected IncrementMaxUsesPerCustomer to return *Discount")
	}

	// Count increments correctly
	if count := discount.GetMaxUsesPerCustomerCount("cust-1"); count != 1 {
		t.Fatalf("expected count 1 after one increment, got %d", count)
	}

	discount.IncrementMaxUsesPerCustomer("cust-1")
	if count := discount.GetMaxUsesPerCustomerCount("cust-1"); count != 2 {
		t.Fatalf("expected count 2 after two increments, got %d", count)
	}

	// Different customers are independent
	if count := discount.GetMaxUsesPerCustomerCount("cust-2"); count != 0 {
		t.Fatalf("expected cust-2 count to be 0, got %d", count)
	}

	// Empty customerID — no-op, returns self
	discount.IncrementMaxUsesPerCustomer("")
	// No panic, no change to any customer
	if count := discount.GetMaxUsesPerCustomerCount("cust-1"); count != 2 {
		t.Fatalf("expected count to remain 2 after empty-ID increment, got %d", count)
	}

	// Data tracking — metas column is marked dirty after increment
	changed := discount.DataChanged()
	if _, ok := changed[COLUMN_METAS]; !ok {
		t.Fatal("expected DataChanged to track metas after increment")
	}
}
