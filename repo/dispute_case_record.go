package repo

import (
	"time"

	"github.com/OpenBazaar/openbazaar-go/pb"
)

var (
	ModeratorDisputeExpiry_firstInterval  = time.Duration(1) * time.Minute
	ModeratorDisputeExpiry_secondInterval = time.Duration(3) * time.Minute
	ModeratorDisputeExpiry_thirdInterval  = time.Duration(4) * time.Minute
	ModeratorDisputeExpiry_lastInterval   = time.Duration(5) * time.Minute
)

// DisputeCaseRecord is a one-to-one relationship with records in the
// SQL datastore
type DisputeCaseRecord struct {
	CaseID           string
	Timestamp        time.Time
	LastNotifiedAt   time.Time
	BuyerContract    *pb.RicardianContract
	VendorContract   *pb.RicardianContract
	IsBuyerInitiated bool
}

// BuildModeratorDisputeExpiryFirstNotification returns a Notification with ExpiresIn set for the First Interval
func (r *DisputeCaseRecord) BuildModeratorDisputeExpiryFirstNotification(createdAt time.Time) *Notification {
	return r.buildModeratorDisputeExpiry(ModeratorDisputeExpiry_firstInterval, createdAt)
}

// BuildModeratorDisputeExpirySecondNotification returns a Notification with ExpiresIn set for the Second Interval
func (r *DisputeCaseRecord) BuildModeratorDisputeExpirySecondNotification(createdAt time.Time) *Notification {
	return r.buildModeratorDisputeExpiry(ModeratorDisputeExpiry_secondInterval, createdAt)
}

// BuildModeratorDisputeExpiryThirdNotification returns a Notification with ExpiresIn set for the Third Interval
func (r *DisputeCaseRecord) BuildModeratorDisputeExpiryThirdNotification(createdAt time.Time) *Notification {
	return r.buildModeratorDisputeExpiry(ModeratorDisputeExpiry_thirdInterval, createdAt)
}

// BuildModeratorDisputeExpiryLastNotification returns a Notification with ExpiresIn set for the Last Interval
func (r *DisputeCaseRecord) BuildModeratorDisputeExpiryLastNotification(createdAt time.Time) *Notification {
	return r.buildModeratorDisputeExpiry(ModeratorDisputeExpiry_lastInterval, createdAt)
}

func (r *DisputeCaseRecord) buildModeratorDisputeExpiry(interval time.Duration, createdAt time.Time) *Notification {
	timeRemaining := ModeratorDisputeExpiry_lastInterval - interval
	notification := ModeratorDisputeExpiry{
		ID:        NewNotificationID(),
		Type:      NotifierTypeModeratorDisputeExpiry,
		CaseID:    r.CaseID,
		ExpiresIn: uint(timeRemaining.Seconds()),
		Thumbnail: Thumbnail{},
	}
	if r.IsBuyerInitiated {
		if len(r.BuyerContract.VendorListings) > 0 && len(r.BuyerContract.VendorListings[0].Item.Images) > 0 {
			notification.Thumbnail = Thumbnail{
				Tiny:  r.BuyerContract.VendorListings[0].Item.Images[0].Tiny,
				Small: r.BuyerContract.VendorListings[0].Item.Images[0].Small,
			}
		}
	} else {
		if len(r.VendorContract.VendorListings) > 0 && len(r.VendorContract.VendorListings[0].Item.Images) > 0 {
			notification.Thumbnail = Thumbnail{
				Tiny:  r.VendorContract.VendorListings[0].Item.Images[0].Tiny,
				Small: r.VendorContract.VendorListings[0].Item.Images[0].Small,
			}
		}
	}
	return NewNotification(notification, createdAt, false)
}
