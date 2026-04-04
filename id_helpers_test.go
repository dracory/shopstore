package shopstore

import (
	"strings"
	"testing"

	"github.com/dracory/uid"
)

func TestGenerateShortIDReturnsLowercaseAndNonEmpty(t *testing.T) {
	generated := GenerateShortID()

	if generated == "" {
		t.Errorf("GenerateShortID() returned empty string")
	}
	if strings.ToLower(generated) != generated {
		t.Errorf("GenerateShortID() returned non-lowercase: %s", generated)
	}

	unshortened, err := uid.UnshortenCrockford(generated)
	if err != nil {
		t.Errorf("UnshortenCrockford failed: %v", err)
	}
	if unshortened == "" {
		t.Errorf("UnshortenCrockford returned empty string")
	}
}

func TestNormalizeIDTrimsWhitespaceAndLowercases(t *testing.T) {
	normalized := NormalizeID("  AbC123XyZ  ")
	if normalized != "abc123xyz" {
		t.Errorf("NormalizeID(\"  AbC123XyZ  \") = %s, want abc123xyz", normalized)
	}
}

func TestNormalizeIDEmptyString(t *testing.T) {
	result := NormalizeID("   ")
	if result != "" {
		t.Errorf("NormalizeID(\"   \") = %s, want empty string", result)
	}
}

func TestIsShortIDReturnsTrueForNineCharIDs(t *testing.T) {
	if !IsShortID("abc123xyz") {
		t.Errorf("IsShortID(\"abc123xyz\") should return true")
	}
}

func TestIsShortIDReturnsTrueForTwentyOneCharIDs(t *testing.T) {
	if !IsShortID("abcdefghijklmnopqrstu") {
		t.Errorf("IsShortID(\"abcdefghijklmnopqrstu\") should return true")
	}
}

func TestIsShortIDReturnsFalseForOtherLengths(t *testing.T) {
	if IsShortID("abcd") {
		t.Errorf("IsShortID(\"abcd\") should return false")
	}
}

func TestShortenIDReturnsNormalizedNineCharID(t *testing.T) {
	result := ShortenID("  ABC123XYZ  ")
	if result != "abc123xyz" {
		t.Errorf("ShortenID(\"  ABC123XYZ  \") = %s, want abc123xyz", result)
	}
}

func TestShortenIDShortensValidThirtyTwoCharID(t *testing.T) {
	longID := uid.HumanUid()
	if len(longID) != 32 {
		t.Errorf("HumanUid() length = %d, want 32", len(longID))
	}

	short := ShortenID(longID)
	if len(short) != 21 {
		t.Errorf("ShortenID() length = %d, want 21", len(short))
	}
	if strings.ToLower(short) != short {
		t.Errorf("ShortenID() returned non-lowercase: %s", short)
	}
}

func TestShortenIDReturnsNormalizedOriginalOnInvalidThirtyTwoCharID(t *testing.T) {
	invalidLong := "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ"

	short := ShortenID(invalidLong)
	if short != strings.ToLower(invalidLong) {
		t.Errorf("ShortenID() = %s, want %s", short, strings.ToLower(invalidLong))
	}
}

func TestShortenIDReturnsNormalizedValueForOtherLengths(t *testing.T) {
	result := ShortenID("  SOME-HANDLE ")
	if result != "some-handle" {
		t.Errorf("ShortenID(\"  SOME-HANDLE \") = %s, want some-handle", result)
	}
}

func TestUnshortenIDReturnsOriginalForNonShortID(t *testing.T) {
	result := UnshortenID("  SOME-HANDLE ")
	if result != "some-handle" {
		t.Errorf("UnshortenID(\"  SOME-HANDLE \") = %s, want some-handle", result)
	}
}

func TestUnshortenIDUnshortensValidTwentyOneCharID(t *testing.T) {
	longID := uid.HumanUid()
	short, err := uid.ShortenCrockford(longID)
	if err != nil {
		t.Fatalf("ShortenCrockford failed: %v", err)
	}
	if len(short) != 21 {
		t.Errorf("ShortenCrockford length = %d, want 21", len(short))
	}

	unshortened := UnshortenID(strings.ToUpper(short))
	if unshortened != longID {
		t.Errorf("UnshortenID() = %s, want %s", unshortened, longID)
	}
}

func TestUnshortenIDReturnsOriginalWhenCrockfordDecodeFails(t *testing.T) {
	invalidShort := "!!!!!!!!!"
	result := UnshortenID(invalidShort)
	if result != invalidShort {
		t.Errorf("UnshortenID() = %s, want %s", result, invalidShort)
	}
}

func TestUnshortenIDReturnsOriginalWhenLengthIsNotSupported(t *testing.T) {
	generated := GenerateShortID()
	if IsShortID(generated) {
		t.Errorf("IsShortID(GenerateShortID()) should return false for generated ID")
	}
	result := UnshortenID(strings.ToUpper(generated))
	if result != generated {
		t.Errorf("UnshortenID() = %s, want %s", result, generated)
	}
}

func TestIsSQLiteReturnsTrueForSQLiteNames(t *testing.T) {
	if !isSQLite("sqlite") {
		t.Errorf("isSQLite(\"sqlite\") should return true")
	}
	if !isSQLite("sqlite3") {
		t.Errorf("isSQLite(\"sqlite3\") should return true")
	}
	if !isSQLite("MySQLiteDriver") {
		t.Errorf("isSQLite(\"MySQLiteDriver\") should return true")
	}
}

func TestIsSQLiteReturnsFalseForNonSQLiteDrivers(t *testing.T) {
	if isSQLite("postgres") {
		t.Errorf("isSQLite(\"postgres\") should return false")
	}
	if isSQLite("mysql") {
		t.Errorf("isSQLite(\"mysql\") should return false")
	}
}
