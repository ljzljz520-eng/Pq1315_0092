package catalog

import (
	"fmt"
	"strings"

	"aromaatelier/internal/model"
)

type AddressChange struct {
	Address string
	Reason  string
}

func ValidateAddressChange(change AddressChange) error {
	if strings.TrimSpace(change.Address) == "" {
		return fmt.Errorf("new address is required")
	}
	if len(strings.TrimSpace(change.Address)) < 8 {
		return fmt.Errorf("address is too short")
	}
	if strings.TrimSpace(change.Reason) == "" {
		return fmt.Errorf("reason is required")
	}
	return nil
}

func ApplyAddressChange(record *model.Record, change AddressChange) error {
	if err := ValidateAddressChange(change); err != nil {
		return err
	}
	if record.Status == model.StatusArchived {
		return fmt.Errorf("archived record cannot move")
	}
	record.PreviousAddr = record.Address
	record.Address = model.SanitizeText(change.Address)
	record.Notes = append(record.Notes, "address changed: "+model.SanitizeText(change.Reason))
	return nil
}
