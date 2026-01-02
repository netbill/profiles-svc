package domain_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/test"
)

func TestProfiles(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	test.CleanDb(t)

	ctx := context.Background()

	firstID := uuid.New()
	secondID := uuid.New()

	first, err := s.domain.profile.Create(ctx, firstID, "first")
	if err != nil {
		t.Fatalf("CreateProfile first: %v", err)
	}

	second, err := s.domain.profile.Create(ctx, secondID, "second")
	if err != nil {
		t.Fatalf("CreateProfile second: %v", err)
	}

	if first.UserID == second.UserID {
		t.Fatalf("expected different IDs, got same: %v", first.UserID)
	}

	first, err = s.domain.profile.GetProfileByID(ctx, firstID)
	if err != nil {
		t.Fatalf("GetProfileByAccountID first: %v", err)
	}

	if first.UserID != firstID {
		t.Fatalf("GetProfileByAccountID first: expected ID %v, got %v", firstID, first.UserID)
	}

	second, err = s.domain.profile.GetProfileByID(ctx, secondID)
	if err != nil {
		t.Fatalf("GetProfileByAccountID second: %v", err)
	}

	if second.UserID != secondID {
		t.Fatalf("GetProfileByAccountID second: expected ID %v, got %v", secondID, second.UserID)
	}

	avatar := "avatar"
	newFirst := "new_first"
	description := "description"

	first, err = s.domain.profile.UpdateProfile(ctx, firstID, domain2.Update{
		Avatar:      &avatar,
		Pseudonym:   &newFirst,
		Description: &description,
	})
	if err != nil {
		t.Fatalf("UpdateProfile first: %v", err)
	}

	if *first.Avatar != avatar {
		t.Fatalf("UpdateProfile first: expected avatar %s, got %s", avatar, *first.Avatar)
	}
	if *first.Pseudonym != newFirst {
		t.Fatalf("UpdateProfile first: expected pseudonym %s, got %s", newFirst, *first.Pseudonym)
	}
	if *first.Description != description {
		t.Fatalf("UpdateProfile first: expected description %s, got %s", description, *first.Description)
	}

	second, err = s.domain.profile.ResetProfile(ctx, secondID)
	if err != nil {
		t.Fatalf("ResetUserProfile second: %v", err)
	}
	if second.Avatar != nil {
		t.Fatalf("ResetUserProfile second: expected avatar nil, got %v", *second.Avatar)
	}
	if second.Pseudonym != nil {
		t.Fatalf("ResetUserProfile second: expected pseudonym nil, got %v", *second.Pseudonym)
	}
	if second.Description != nil {
		t.Fatalf("ResetUserProfile second: expected description nil, got %v", *second.Description)
	}

	first, err = s.domain.profile.ResetProfile(ctx, firstID)
	if err != nil {
		t.Fatalf("ResetUsername first: %v", err)
	}
	if first.Username == "first" {
		t.Fatalf("ResetUsername first: expected username not %s, got %s", "first", first.Username)
	}

	first, err = s.domain.profile.UpdateProfileOfficial(ctx, firstID, false)
	if err != nil {
		t.Fatalf("UpdateProfileOfficial first to false: %v", err)
	}

	if first.Official {
		t.Fatalf("UpdateProfileOfficial first to false: expected official false, got true")
	}

	first, err = s.domain.profile.UpdateProfileOfficial(ctx, firstID, true)
	if err != nil {
		t.Fatalf("UpdateProfileOfficial first to true: %v", err)
	}

	if !first.Official {
		t.Fatalf("UpdateProfileOfficial first to true: expected official true, got false")
	}

	list, err := s.domain.profile.Filter(ctx, profile.FilterParams{
		UserID: []uuid.UUID{firstID, secondID},
	}, 0, 10)
	if err != nil {
		t.Fatalf("FilterProfiles by IDs: %v", err)
	}
	if len(list.Data) != 2 {
		t.Fatalf("FilterProfiles by IDs: expected 2 profiles, got %d", len(list.Data))
	}

	first, err = s.domain.profile.UpdateProfileUsername(ctx, firstID, "first")
	if err != nil {
		t.Fatalf("UpdateProfileUsername first: %v", err)
	}
	first, err = s.domain.profile.UpdateProfile(ctx, firstID, domain2.Update{
		Avatar:      func() *string { s := "avatar"; return &s }(),
		Pseudonym:   func() *string { s := "new_first"; return &s }(),
		Description: func() *string { s := "first description"; return &s }(),
	})

	second, err = s.domain.profile.UpdateProfileUsername(ctx, secondID, "second")
	if err != nil {
		t.Fatalf("UpdateProfileUsername second: %v", err)
	}
	second, err = s.domain.profile.UpdateProfile(ctx, secondID, domain2.Update{
		Avatar:      func() *string { s := "avatar2"; return &s }(),
		Pseudonym:   func() *string { s := "new_second"; return &s }(),
		Description: func() *string { s := "second description"; return &s }(),
	})

	third, err := s.domain.profile.Create(ctx, uuid.New(), "third")
	if err != nil {
		t.Fatalf("CreateProfile third: %v", err)
	}
	third, err = s.domain.profile.UpdateProfile(ctx, third.UserID, domain2.Update{
		Avatar:      func() *string { s := "avatar3"; return &s }(),
		Pseudonym:   func() *string { s := "new_third"; return &s }(),
		Description: func() *string { s := "third description"; return &s }(),
	})
	if err != nil {
		t.Fatalf("UpdateProfile third: %v", err)
	}

	list, err = s.domain.profile.Filter(ctx, profile.FilterParams{}, 0, 10)
	if err != nil {
		t.Fatalf("FilterProfiles all: %v", err)
	}
	if len(list.Data) != 3 {
		t.Fatalf("FilterProfiles all: expected 3 profiles, got %d", len(list.Data))
	}
}
